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

// RegisterOncallReadTools registers read-only on-call tools
func RegisterOncallReadTools(s *server.MCPServer, c *client.Client) {
	// list_oncalls
	s.AddTool(mcp.NewTool("list_oncalls",
		mcp.WithDescription("List current and upcoming on-call entries. Returns who is on-call right now or during a specified time range. Use 'earliest=true' to get just the current on-call person for each schedule. This is the primary tool for finding who to contact for an incident."),
		mcp.WithTitleAnnotation("List On-Calls"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithString("time_zone", mcp.Description("IANA time zone for returned times (e.g., 'America/New_York', 'UTC')")),
		mcp.WithString("since", mcp.Description("Start of date range in ISO 8601 format (e.g., '2024-01-15T00:00:00Z'). Defaults to now.")),
		mcp.WithString("until", mcp.Description("End of date range in ISO 8601 format (e.g., '2024-01-22T00:00:00Z'). Defaults to now.")),
		mcp.WithBoolean("earliest", mcp.Description("If true, return only the earliest/current on-call entry for each schedule. Useful for finding who is on-call right now.")),
		mcp.WithString("schedule_ids", mcp.Description("Filter by schedules. Comma-separated schedule IDs (e.g., 'PSCHED1,PSCHED2')")),
		mcp.WithString("user_ids", mcp.Description("Filter by users. Comma-separated user IDs (e.g., 'PUSER1,PUSER2')")),
		mcp.WithString("escalation_policy_ids", mcp.Description("Filter by escalation policies. Comma-separated policy IDs (e.g., 'PESCPOL1,PESCPOL2')")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results to return"), mcp.Min(1), mcp.Max(100)),
	), listOncallsHandler(c))
}

func listOncallsHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		params := make(map[string]string)

		if v, ok := getString(args, "time_zone"); ok {
			params["time_zone"] = v
		}
		if v, ok := getString(args, "since"); ok {
			params["since"] = v
		}
		if v, ok := getString(args, "until"); ok {
			params["until"] = v
		}
		if v, ok := getBool(args, "earliest"); ok && v {
			params["earliest"] = "true"
		}
		if v, ok := getString(args, "schedule_ids"); ok {
			params["schedule_ids[]"] = v
		}
		if v, ok := getString(args, "user_ids"); ok {
			params["user_ids[]"] = v
		}
		if v, ok := getString(args, "escalation_policy_ids"); ok {
			params["escalation_policy_ids[]"] = v
		}
		if v, ok := getNumber(args, "limit"); ok {
			params["limit"] = fmt.Sprintf("%d", int(v))
		}

		var resp models.OncallsResponse
		if err := c.GetJSON("/oncalls", params, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		result := models.ListResponse[models.Oncall]{Response: resp.Oncalls}
		data, _ := json.Marshal(result)
		return mcp.NewToolResultText(string(data)), nil
	}
}
