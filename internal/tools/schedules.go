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

// RegisterScheduleReadTools registers read-only schedule tools
func RegisterScheduleReadTools(s *server.MCPServer, c *client.Client) {
	// list_schedules
	s.AddTool(mcp.NewTool("list_schedules",
		mcp.WithDescription("List on-call schedules in PagerDuty. Schedules define rotation patterns for who is on-call at any given time. Use to find schedule IDs for filtering on-calls or understanding coverage."),
		mcp.WithTitleAnnotation("List Schedules"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithString("query", mcp.Description("Filter schedules by name (partial match supported)")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results to return"), mcp.Min(1), mcp.Max(100)),
	), listSchedulesHandler(c))

	// get_schedule
	s.AddTool(mcp.NewTool("get_schedule",
		mcp.WithDescription("Get detailed information about a specific on-call schedule, including rotation layers and rendered on-call periods for a given time range."),
		mcp.WithTitleAnnotation("Get Schedule Details"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithString("schedule_id", mcp.Required(), mcp.Description("The unique schedule ID (e.g., 'PSCHED123')")),
		mcp.WithString("since", mcp.Description("Start of date range in ISO 8601 format (e.g., '2024-01-15T00:00:00Z'). Used to render on-call periods.")),
		mcp.WithString("until", mcp.Description("End of date range in ISO 8601 format (e.g., '2024-01-22T00:00:00Z'). Used to render on-call periods.")),
	), getScheduleHandler(c))

	// list_schedule_users
	s.AddTool(mcp.NewTool("list_schedule_users",
		mcp.WithDescription("List all users who are part of a schedule's rotation within a given time range. Use to see who is or will be on-call."),
		mcp.WithTitleAnnotation("List Schedule Users"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithString("schedule_id", mcp.Required(), mcp.Description("The unique schedule ID (e.g., 'PSCHED123')")),
		mcp.WithString("since", mcp.Description("Start of date range in ISO 8601 format (e.g., '2024-01-15T00:00:00Z')")),
		mcp.WithString("until", mcp.Description("End of date range in ISO 8601 format (e.g., '2024-01-22T00:00:00Z')")),
	), listScheduleUsersHandler(c))
}

// RegisterScheduleWriteTools registers write schedule tools
func RegisterScheduleWriteTools(s *server.MCPServer, c *client.Client) {
	// create_schedule
	s.AddTool(mcp.NewTool("create_schedule",
		mcp.WithDescription("Create a new on-call schedule. Schedules define rotation patterns for on-call coverage. Note: This creates an empty schedule - rotation layers need to be added separately."),
		mcp.WithTitleAnnotation("Create Schedule"),
		mcp.WithString("name", mcp.Required(), mcp.Description("A descriptive name for the schedule (e.g., 'Primary On-Call', 'Weekend Coverage')")),
		mcp.WithString("time_zone", mcp.Required(), mcp.Description("IANA time zone identifier (e.g., 'America/New_York', 'Europe/London', 'UTC')")),
		mcp.WithString("description", mcp.Description("Description of the schedule's purpose and coverage")),
	), createScheduleHandler(c))

	// create_schedule_override
	s.AddTool(mcp.NewTool("create_schedule_override",
		mcp.WithDescription("Create a temporary override on a schedule. Use for vacation coverage, shift swaps, or any temporary change to the normal rotation. The override takes precedence over the regular schedule during the specified time window."),
		mcp.WithTitleAnnotation("Create Schedule Override"),
		mcp.WithString("schedule_id", mcp.Required(), mcp.Description("The unique schedule ID (e.g., 'PSCHED123')")),
		mcp.WithString("user_id", mcp.Required(), mcp.Description("The user ID who will be on-call during the override (e.g., 'PUSER123')")),
		mcp.WithString("start", mcp.Required(), mcp.Description("Override start time in ISO 8601 format (e.g., '2024-01-15T09:00:00Z')")),
		mcp.WithString("end", mcp.Required(), mcp.Description("Override end time in ISO 8601 format (e.g., '2024-01-15T17:00:00Z')")),
	), createScheduleOverrideHandler(c))

	// update_schedule
	s.AddTool(mcp.NewTool("update_schedule",
		mcp.WithDescription("Update an existing schedule's metadata (name, description, time zone). Does not modify rotation layers."),
		mcp.WithTitleAnnotation("Update Schedule"),
		mcp.WithString("schedule_id", mcp.Required(), mcp.Description("The unique schedule ID to update (e.g., 'PSCHED123')")),
		mcp.WithString("name", mcp.Description("New schedule name")),
		mcp.WithString("description", mcp.Description("New schedule description")),
		mcp.WithString("time_zone", mcp.Description("New IANA time zone identifier (e.g., 'America/New_York')")),
	), updateScheduleHandler(c))
}

func listSchedulesHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		params := make(map[string]string)

		if v, ok := getString(args, "query"); ok {
			params["query"] = v
		}
		if v, ok := getNumber(args, "limit"); ok {
			params["limit"] = fmt.Sprintf("%d", int(v))
		}

		var resp models.SchedulesResponse
		if err := c.GetJSON("/schedules", params, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		result := models.ListResponse[models.Schedule]{Response: resp.Schedules}
		data, _ := json.Marshal(result)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func getScheduleHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		scheduleID, ok := getString(args, "schedule_id")
		if !ok {
			return mcp.NewToolResultError("schedule_id is required"), nil
		}

		params := make(map[string]string)
		if v, ok := getString(args, "since"); ok {
			params["since"] = v
		}
		if v, ok := getString(args, "until"); ok {
			params["until"] = v
		}

		var resp models.ScheduleResponse
		if err := c.GetJSON(fmt.Sprintf("/schedules/%s", scheduleID), params, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, _ := json.Marshal(resp.Schedule)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func listScheduleUsersHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		scheduleID, ok := getString(args, "schedule_id")
		if !ok {
			return mcp.NewToolResultError("schedule_id is required"), nil
		}

		params := make(map[string]string)
		if v, ok := getString(args, "since"); ok {
			params["since"] = v
		}
		if v, ok := getString(args, "until"); ok {
			params["until"] = v
		}

		var resp models.ScheduleUsersResponse
		if err := c.GetJSON(fmt.Sprintf("/schedules/%s/users", scheduleID), params, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		result := models.ListResponse[models.User]{Response: resp.Users}
		data, _ := json.Marshal(result)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func createScheduleHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		name, ok := getString(args, "name")
		if !ok {
			return mcp.NewToolResultError("name is required"), nil
		}

		timeZone, ok := getString(args, "time_zone")
		if !ok {
			return mcp.NewToolResultError("time_zone is required"), nil
		}

		schedule := models.ScheduleCreateData{
			Type:           "schedule",
			Name:           name,
			TimeZone:       timeZone,
			ScheduleLayers: []models.ScheduleLayerCreate{},
		}

		if v, ok := getString(args, "description"); ok {
			schedule.Description = v
		}

		req := models.ScheduleCreateRequest{Schedule: schedule}

		var resp models.ScheduleResponse
		if err := c.PostJSON("/schedules", req, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, _ := json.Marshal(resp.Schedule)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func createScheduleOverrideHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		scheduleID, ok := getString(args, "schedule_id")
		if !ok {
			return mcp.NewToolResultError("schedule_id is required"), nil
		}

		userID, ok := getString(args, "user_id")
		if !ok {
			return mcp.NewToolResultError("user_id is required"), nil
		}

		start, ok := getString(args, "start")
		if !ok {
			return mcp.NewToolResultError("start is required"), nil
		}

		end, ok := getString(args, "end")
		if !ok {
			return mcp.NewToolResultError("end is required"), nil
		}

		override := models.ScheduleOverrideCreate{
			Override: models.OverrideData{
				Start: start,
				End:   end,
				User: models.UserReference{
					ID:   userID,
					Type: "user_reference",
				},
			},
		}

		var resp models.ScheduleOverrideResponse
		if err := c.PostJSON(fmt.Sprintf("/schedules/%s/overrides", scheduleID), override, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, _ := json.Marshal(resp.Override)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func updateScheduleHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		scheduleID, ok := getString(args, "schedule_id")
		if !ok {
			return mcp.NewToolResultError("schedule_id is required"), nil
		}

		schedule := models.ScheduleUpdateData{
			Type: "schedule",
		}

		if v, ok := getString(args, "name"); ok {
			schedule.Name = v
		}
		if v, ok := getString(args, "description"); ok {
			schedule.Description = v
		}
		if v, ok := getString(args, "time_zone"); ok {
			schedule.TimeZone = v
		}

		req := models.ScheduleUpdateRequest{Schedule: schedule}

		var resp models.ScheduleResponse
		if err := c.PutJSON(fmt.Sprintf("/schedules/%s", scheduleID), req, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, _ := json.Marshal(resp.Schedule)
		return mcp.NewToolResultText(string(data)), nil
	}
}
