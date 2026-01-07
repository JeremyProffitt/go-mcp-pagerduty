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

// RegisterUserReadTools registers read-only user tools
func RegisterUserReadTools(s *server.MCPServer, c *client.Client) {
	// get_user_data
	s.AddTool(mcp.NewTool("get_user_data",
		mcp.WithDescription("Get the current authenticated user's information. This returns details about the user whose API token is being used, including their ID, name, email, and role. Call this first to scope subsequent requests by user ID."),
		mcp.WithTitleAnnotation("Get Current User"),
		mcp.WithReadOnlyHintAnnotation(true),
	), getUserDataHandler(c))

	// list_users
	s.AddTool(mcp.NewTool("list_users",
		mcp.WithDescription("List users in the PagerDuty account. Use to find user IDs for assignments, team membership, or filtering incidents."),
		mcp.WithTitleAnnotation("List Users"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithString("query", mcp.Description("Filter users by name or email address (partial match supported)")),
		mcp.WithString("team_ids", mcp.Description("Filter by team membership. Comma-separated team IDs (e.g., 'PTEAM1,PTEAM2')")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results to return"), mcp.Min(1), mcp.Max(100)),
	), listUsersHandler(c))
}

func getUserDataHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var resp models.UserResponse
		if err := c.GetJSON("/users/me", nil, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, _ := json.Marshal(resp.User)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func listUsersHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		params := make(map[string]string)

		if v, ok := getString(args, "query"); ok {
			params["query"] = v
		}
		if v, ok := getString(args, "team_ids"); ok {
			params["team_ids[]"] = v
		}
		if v, ok := getNumber(args, "limit"); ok {
			params["limit"] = fmt.Sprintf("%d", int(v))
		}

		var resp models.UsersResponse
		if err := c.GetJSON("/users", params, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		result := models.ListResponse[models.User]{Response: resp.Users}
		data, _ := json.Marshal(result)
		return mcp.NewToolResultText(string(data)), nil
	}
}
