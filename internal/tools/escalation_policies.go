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

// RegisterEscalationPolicyReadTools registers read-only escalation policy tools
func RegisterEscalationPolicyReadTools(s *server.MCPServer, c *client.Client) {
	// list_escalation_policies
	s.AddTool(mcp.NewTool("list_escalation_policies",
		mcp.WithDescription("List escalation policies in PagerDuty. Escalation policies define the order in which users and schedules are notified when an incident occurs. Each service must have an escalation policy. Use to find policy IDs for creating services or understanding notification chains."),
		mcp.WithTitleAnnotation("List Escalation Policies"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithString("query", mcp.Description("Filter policies by name (partial match supported)")),
		mcp.WithString("user_ids", mcp.Description("Filter by users in the policy. Comma-separated user IDs (e.g., 'PUSER1,PUSER2')")),
		mcp.WithString("team_ids", mcp.Description("Filter by associated teams. Comma-separated team IDs (e.g., 'PTEAM1,PTEAM2')")),
		mcp.WithString("sort_by", mcp.Description("Sort order for results"), mcp.Enum("name", "name:asc", "name:desc")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results to return"), mcp.Min(1), mcp.Max(100)),
	), listEscalationPoliciesHandler(c))

	// get_escalation_policy
	s.AddTool(mcp.NewTool("get_escalation_policy",
		mcp.WithDescription("Get detailed information about a specific escalation policy, including all escalation levels and targets (users, schedules) at each level."),
		mcp.WithTitleAnnotation("Get Escalation Policy Details"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithString("escalation_policy_id", mcp.Required(), mcp.Description("The unique escalation policy ID (e.g., 'PESCPOL123')")),
	), getEscalationPolicyHandler(c))
}

func listEscalationPoliciesHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		params := make(map[string]string)

		if v, ok := getString(args, "query"); ok {
			params["query"] = v
		}
		if v, ok := getString(args, "user_ids"); ok {
			params["user_ids[]"] = v
		}
		if v, ok := getString(args, "team_ids"); ok {
			params["team_ids[]"] = v
		}
		if v, ok := getString(args, "sort_by"); ok {
			params["sort_by"] = v
		}
		if v, ok := getNumber(args, "limit"); ok {
			params["limit"] = fmt.Sprintf("%d", int(v))
		}

		var resp models.EscalationPoliciesResponse
		if err := c.GetJSON("/escalation_policies", params, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		result := models.ListResponse[models.EscalationPolicy]{Response: resp.EscalationPolicies}
		data, _ := json.Marshal(result)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func getEscalationPolicyHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		policyID, ok := getString(args, "escalation_policy_id")
		if !ok {
			return mcp.NewToolResultError("escalation_policy_id is required"), nil
		}

		var resp models.EscalationPolicyResponse
		if err := c.GetJSON(fmt.Sprintf("/escalation_policies/%s", policyID), nil, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, _ := json.Marshal(resp.EscalationPolicy)
		return mcp.NewToolResultText(string(data)), nil
	}
}
