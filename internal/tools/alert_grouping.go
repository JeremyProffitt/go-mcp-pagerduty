package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jeremyproffitt/go-mcp-pagerduty/internal/client"
	"github.com/jeremyproffitt/go-mcp-pagerduty/internal/models"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterAlertGroupingReadTools registers read-only alert grouping tools
func RegisterAlertGroupingReadTools(s *server.MCPServer, c *client.Client) {
	// list_alert_grouping_settings
	s.AddTool(mcp.NewTool("list_alert_grouping_settings",
		mcp.WithDescription("List alert grouping settings. Alert grouping combines multiple related alerts into a single incident to reduce noise. Settings can be time-based, intelligent (ML-based), or content-based grouping."),
		mcp.WithTitleAnnotation("List Alert Grouping Settings"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithString("service_ids", mcp.Description("Filter by services. Comma-separated service IDs (e.g., 'PDSVC1,PDSVC2')")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results to return"), mcp.Min(1), mcp.Max(100)),
	), listAlertGroupingSettingsHandler(c))

	// get_alert_grouping_setting
	s.AddTool(mcp.NewTool("get_alert_grouping_setting",
		mcp.WithDescription("Get detailed information about a specific alert grouping configuration, including its type (time, intelligent, content_based) and associated services."),
		mcp.WithTitleAnnotation("Get Alert Grouping Setting"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithString("setting_id", mcp.Required(), mcp.Description("The unique alert grouping setting ID")),
	), getAlertGroupingSettingHandler(c))
}

// RegisterAlertGroupingWriteTools registers write alert grouping tools
func RegisterAlertGroupingWriteTools(s *server.MCPServer, c *client.Client) {
	// create_alert_grouping_setting
	s.AddTool(mcp.NewTool("create_alert_grouping_setting",
		mcp.WithDescription("Create a new alert grouping configuration for services. Alert grouping reduces noise by combining related alerts into single incidents. Choose 'time' for simple time windows, 'intelligent' for ML-based grouping, or 'content_based' for field matching."),
		mcp.WithTitleAnnotation("Create Alert Grouping Setting"),
		mcp.WithString("name", mcp.Required(), mcp.Description("A descriptive name for the alert grouping configuration")),
		mcp.WithString("service_ids", mcp.Required(), mcp.Description("Services to apply this grouping to. Comma-separated service IDs (e.g., 'PDSVC1,PDSVC2')")),
		mcp.WithString("type", mcp.Required(), mcp.Description("Alert grouping strategy"), mcp.Enum("time", "intelligent", "content_based")),
		mcp.WithNumber("timeout", mcp.Description("Time window in minutes for grouping alerts (only for 'time' type, default: 5)"), mcp.Min(1), mcp.Max(1440)),
	), createAlertGroupingSettingHandler(c))

	// update_alert_grouping_setting
	s.AddTool(mcp.NewTool("update_alert_grouping_setting",
		mcp.WithDescription("Update an existing alert grouping configuration. Can change the grouping strategy or timeout settings."),
		mcp.WithTitleAnnotation("Update Alert Grouping Setting"),
		mcp.WithString("setting_id", mcp.Required(), mcp.Description("The unique alert grouping setting ID to update")),
		mcp.WithString("name", mcp.Description("New name for the setting")),
		mcp.WithString("type", mcp.Description("New grouping strategy"), mcp.Enum("time", "intelligent", "content_based")),
		mcp.WithNumber("timeout", mcp.Description("New time window in minutes (only for 'time' type)"), mcp.Min(1), mcp.Max(1440)),
	), updateAlertGroupingSettingHandler(c))

	// delete_alert_grouping_setting
	s.AddTool(mcp.NewTool("delete_alert_grouping_setting",
		mcp.WithDescription("WARNING: DESTRUCTIVE - Permanently delete an alert grouping configuration. Affected services will revert to their default alert grouping behavior."),
		mcp.WithTitleAnnotation("Delete Alert Grouping Setting"),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithString("setting_id", mcp.Required(), mcp.Description("The unique alert grouping setting ID to delete")),
	), deleteAlertGroupingSettingHandler(c))
}

func listAlertGroupingSettingsHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		params := make(map[string]string)

		if v, ok := getString(args, "service_ids"); ok {
			params["service_ids[]"] = v
		}
		if v, ok := getNumber(args, "limit"); ok {
			params["limit"] = fmt.Sprintf("%d", int(v))
		}

		var resp models.AlertGroupingSettingsResponse
		if err := c.GetJSON("/alert_grouping_settings", params, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		result := models.ListResponse[models.AlertGroupingSetting]{Response: resp.AlertGroupingSettings}
		data, _ := json.Marshal(result)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func getAlertGroupingSettingHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		settingID, ok := getString(args, "setting_id")
		if !ok {
			return mcp.NewToolResultError("setting_id is required"), nil
		}

		var resp models.AlertGroupingSettingResponse
		if err := c.GetJSON(fmt.Sprintf("/alert_grouping_settings/%s", settingID), nil, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, _ := json.Marshal(resp.AlertGroupingSetting)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func createAlertGroupingSettingHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		name, ok := getString(args, "name")
		if !ok {
			return mcp.NewToolResultError("name is required"), nil
		}

		serviceIDsStr, ok := getString(args, "service_ids")
		if !ok {
			return mcp.NewToolResultError("service_ids is required"), nil
		}

		groupingType, ok := getString(args, "type")
		if !ok {
			return mcp.NewToolResultError("type is required"), nil
		}

		serviceIDs := splitAndTrim(serviceIDsStr)
		services := make([]models.ServiceReference, len(serviceIDs))
		for i, id := range serviceIDs {
			services[i] = models.ServiceReference{
				ID:   id,
				Type: "service_reference",
			}
		}

		config := models.AlertGroupingConfig{
			Type: groupingType,
		}

		if v, ok := getNumber(args, "timeout"); ok {
			config.Timeout = int(v)
		}

		setting := models.AlertGroupingSettingCreate{
			Type:     "alert_grouping_setting",
			Name:     name,
			Services: services,
			Config:   config,
		}

		req := models.AlertGroupingSettingCreateRequest{AlertGroupingSetting: setting}

		var resp models.AlertGroupingSettingResponse
		if err := c.PostJSON("/alert_grouping_settings", req, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, _ := json.Marshal(resp.AlertGroupingSetting)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func updateAlertGroupingSettingHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		settingID, ok := getString(args, "setting_id")
		if !ok {
			return mcp.NewToolResultError("setting_id is required"), nil
		}

		setting := models.AlertGroupingSettingUpdate{
			Type: "alert_grouping_setting",
		}

		if v, ok := getString(args, "name"); ok {
			setting.Name = v
		}
		if v, ok := getString(args, "type"); ok {
			setting.Config = &models.AlertGroupingConfig{Type: v}
		}
		if v, ok := getNumber(args, "timeout"); ok {
			if setting.Config == nil {
				setting.Config = &models.AlertGroupingConfig{}
			}
			setting.Config.Timeout = int(v)
		}

		req := models.AlertGroupingSettingUpdateRequest{AlertGroupingSetting: setting}

		var resp models.AlertGroupingSettingResponse
		if err := c.PutJSON(fmt.Sprintf("/alert_grouping_settings/%s", settingID), req, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, _ := json.Marshal(resp.AlertGroupingSetting)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func deleteAlertGroupingSettingHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		settingID, ok := getString(args, "setting_id")
		if !ok {
			return mcp.NewToolResultError("setting_id is required"), nil
		}

		if _, err := c.Delete(fmt.Sprintf("/alert_grouping_settings/%s", settingID)); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Alert grouping setting %s deleted successfully", settingID)), nil
	}
}
