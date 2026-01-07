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

// RegisterChangeEventReadTools registers read-only change event tools
func RegisterChangeEventReadTools(s *server.MCPServer, c *client.Client) {
	// list_change_events
	s.AddTool(mcp.NewTool("list_change_events",
		mcp.WithDescription("List change events (deployments, releases, config changes) across PagerDuty. Change events help correlate incidents with recent changes. Use this to investigate if a deployment caused an incident."),
		mcp.WithTitleAnnotation("List Change Events"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithString("since", mcp.Description("Start date in ISO 8601 format (e.g., '2024-01-15T00:00:00Z')")),
		mcp.WithString("until", mcp.Description("End date in ISO 8601 format (e.g., '2024-01-16T00:00:00Z')")),
		mcp.WithString("team_ids", mcp.Description("Filter by teams. Comma-separated team IDs (e.g., 'PTEAM1,PTEAM2')")),
		mcp.WithString("service_ids", mcp.Description("Filter by services. Comma-separated service IDs (e.g., 'PDSVC1,PDSVC2')")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results to return"), mcp.Min(1), mcp.Max(100)),
	), listChangeEventsHandler(c))

	// get_change_event
	s.AddTool(mcp.NewTool("get_change_event",
		mcp.WithDescription("Get detailed information about a specific change event, including its summary, source, and links to related resources."),
		mcp.WithTitleAnnotation("Get Change Event"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithString("change_event_id", mcp.Required(), mcp.Description("The unique change event ID")),
	), getChangeEventHandler(c))

	// list_service_change_events
	s.AddTool(mcp.NewTool("list_service_change_events",
		mcp.WithDescription("List change events for a specific service. Use this when investigating incidents on a service to see recent deployments or changes that might have caused the issue."),
		mcp.WithTitleAnnotation("List Service Change Events"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithString("service_id", mcp.Required(), mcp.Description("The unique service ID (e.g., 'PDSVC123')")),
		mcp.WithString("since", mcp.Description("Start date in ISO 8601 format (e.g., '2024-01-15T00:00:00Z')")),
		mcp.WithString("until", mcp.Description("End date in ISO 8601 format (e.g., '2024-01-16T00:00:00Z')")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results to return"), mcp.Min(1), mcp.Max(100)),
	), listServiceChangeEventsHandler(c))

	// list_incident_change_events
	s.AddTool(mcp.NewTool("list_incident_change_events",
		mcp.WithDescription("List change events that PagerDuty has automatically correlated with an incident. Shows deployments and changes that occurred around the time the incident started and may have caused it."),
		mcp.WithTitleAnnotation("List Related Change Events"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithString("incident_id", mcp.Required(), mcp.Description("The unique incident ID (e.g., 'PABC123')")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results to return"), mcp.Min(1), mcp.Max(100)),
	), listIncidentChangeEventsHandler(c))
}

func listChangeEventsHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		params := make(map[string]string)

		if v, ok := getString(args, "since"); ok {
			params["since"] = v
		}
		if v, ok := getString(args, "until"); ok {
			params["until"] = v
		}
		if v, ok := getString(args, "team_ids"); ok {
			params["team_ids[]"] = v
		}
		if v, ok := getString(args, "service_ids"); ok {
			params["service_ids[]"] = v
		}
		if v, ok := getNumber(args, "limit"); ok {
			params["limit"] = fmt.Sprintf("%d", int(v))
		}

		var resp models.ChangeEventsResponse
		if err := c.GetJSON("/change_events", params, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		result := models.ListResponse[models.ChangeEvent]{Response: resp.ChangeEvents}
		data, _ := json.Marshal(result)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func getChangeEventHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		changeEventID, ok := getString(args, "change_event_id")
		if !ok {
			return mcp.NewToolResultError("change_event_id is required"), nil
		}

		var resp models.ChangeEventResponse
		if err := c.GetJSON(fmt.Sprintf("/change_events/%s", changeEventID), nil, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, _ := json.Marshal(resp.ChangeEvent)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func listServiceChangeEventsHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		serviceID, ok := getString(args, "service_id")
		if !ok {
			return mcp.NewToolResultError("service_id is required"), nil
		}

		params := make(map[string]string)
		if v, ok := getString(args, "since"); ok {
			params["since"] = v
		}
		if v, ok := getString(args, "until"); ok {
			params["until"] = v
		}
		if v, ok := getNumber(args, "limit"); ok {
			params["limit"] = fmt.Sprintf("%d", int(v))
		}

		var resp models.ChangeEventsResponse
		if err := c.GetJSON(fmt.Sprintf("/services/%s/change_events", serviceID), params, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		result := models.ListResponse[models.ChangeEvent]{Response: resp.ChangeEvents}
		data, _ := json.Marshal(result)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func listIncidentChangeEventsHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		incidentID, ok := getString(args, "incident_id")
		if !ok {
			return mcp.NewToolResultError("incident_id is required"), nil
		}

		params := make(map[string]string)
		if v, ok := getNumber(args, "limit"); ok {
			params["limit"] = fmt.Sprintf("%d", int(v))
		}

		var resp models.ChangeEventsResponse
		if err := c.GetJSON(fmt.Sprintf("/incidents/%s/related_change_events", incidentID), params, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		result := models.ListResponse[models.ChangeEvent]{Response: resp.ChangeEvents}
		data, _ := json.Marshal(result)
		return mcp.NewToolResultText(string(data)), nil
	}
}
