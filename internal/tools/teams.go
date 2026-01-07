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

// RegisterTeamReadTools registers read-only team tools
func RegisterTeamReadTools(s *server.MCPServer, c *client.Client) {
	// list_teams
	s.AddTool(mcp.NewTool("list_teams",
		mcp.WithDescription("List teams in PagerDuty. Teams are organizational units that group users together. Use to find team IDs for filtering services, escalation policies, or incidents."),
		mcp.WithTitleAnnotation("List Teams"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithString("query", mcp.Description("Filter teams by name (partial match supported)")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results to return"), mcp.Min(1), mcp.Max(100)),
	), listTeamsHandler(c))

	// get_team
	s.AddTool(mcp.NewTool("get_team",
		mcp.WithDescription("Get detailed information about a specific team including its description and settings."),
		mcp.WithTitleAnnotation("Get Team Details"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithString("team_id", mcp.Required(), mcp.Description("The unique team ID (e.g., 'PTEAM123')")),
	), getTeamHandler(c))

	// list_team_members
	s.AddTool(mcp.NewTool("list_team_members",
		mcp.WithDescription("List all users who are members of a specific team, including their roles (manager, responder, observer)."),
		mcp.WithTitleAnnotation("List Team Members"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithString("team_id", mcp.Required(), mcp.Description("The unique team ID (e.g., 'PTEAM123')")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results to return"), mcp.Min(1), mcp.Max(100)),
	), listTeamMembersHandler(c))
}

// RegisterTeamWriteTools registers write team tools
func RegisterTeamWriteTools(s *server.MCPServer, c *client.Client) {
	// create_team
	s.AddTool(mcp.NewTool("create_team",
		mcp.WithDescription("Create a new team to organize users. Teams can be associated with services, escalation policies, and used to filter incidents."),
		mcp.WithTitleAnnotation("Create Team"),
		mcp.WithString("name", mcp.Required(), mcp.Description("A descriptive name for the team (e.g., 'Platform Engineering', 'Customer Support')")),
		mcp.WithString("description", mcp.Description("Description of the team's purpose and responsibilities")),
	), createTeamHandler(c))

	// update_team
	s.AddTool(mcp.NewTool("update_team",
		mcp.WithDescription("Update an existing team's name or description."),
		mcp.WithTitleAnnotation("Update Team"),
		mcp.WithString("team_id", mcp.Required(), mcp.Description("The unique team ID to update (e.g., 'PTEAM123')")),
		mcp.WithString("name", mcp.Description("New team name")),
		mcp.WithString("description", mcp.Description("New team description")),
	), updateTeamHandler(c))

	// delete_team
	s.AddTool(mcp.NewTool("delete_team",
		mcp.WithDescription("WARNING: DESTRUCTIVE - Permanently delete a team. This removes the team from all associated services and escalation policies. This action cannot be undone."),
		mcp.WithTitleAnnotation("Delete Team"),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithString("team_id", mcp.Required(), mcp.Description("The unique team ID to delete (e.g., 'PTEAM123')")),
	), deleteTeamHandler(c))

	// add_team_member
	s.AddTool(mcp.NewTool("add_team_member",
		mcp.WithDescription("Add a user to a team with a specified role. Users can have different permissions based on their role within the team."),
		mcp.WithTitleAnnotation("Add Team Member"),
		mcp.WithString("team_id", mcp.Required(), mcp.Description("The unique team ID (e.g., 'PTEAM123')")),
		mcp.WithString("user_id", mcp.Required(), mcp.Description("The user ID to add to the team (e.g., 'PUSER123')")),
		mcp.WithString("role", mcp.Description("Member role within the team"), mcp.Enum("manager", "responder", "observer")),
	), addTeamMemberHandler(c))

	// remove_team_member
	s.AddTool(mcp.NewTool("remove_team_member",
		mcp.WithDescription("WARNING: DESTRUCTIVE - Remove a user from a team. The user will lose any team-specific permissions and may be removed from associated schedules and escalation policies."),
		mcp.WithTitleAnnotation("Remove Team Member"),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithString("team_id", mcp.Required(), mcp.Description("The unique team ID (e.g., 'PTEAM123')")),
		mcp.WithString("user_id", mcp.Required(), mcp.Description("The user ID to remove from the team (e.g., 'PUSER123')")),
	), removeTeamMemberHandler(c))
}

func listTeamsHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		params := make(map[string]string)

		if v, ok := getString(args, "query"); ok {
			params["query"] = v
		}
		if v, ok := getNumber(args, "limit"); ok {
			params["limit"] = fmt.Sprintf("%d", int(v))
		}

		var resp models.TeamsResponse
		if err := c.GetJSON("/teams", params, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		result := models.ListResponse[models.Team]{Response: resp.Teams}
		data, _ := json.Marshal(result)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func getTeamHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		teamID, ok := getString(args, "team_id")
		if !ok {
			return mcp.NewToolResultError("team_id is required"), nil
		}

		var resp models.TeamResponse
		if err := c.GetJSON(fmt.Sprintf("/teams/%s", teamID), nil, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, _ := json.Marshal(resp.Team)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func listTeamMembersHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		teamID, ok := getString(args, "team_id")
		if !ok {
			return mcp.NewToolResultError("team_id is required"), nil
		}

		params := make(map[string]string)
		if v, ok := getNumber(args, "limit"); ok {
			params["limit"] = fmt.Sprintf("%d", int(v))
		}

		var resp models.TeamMembersResponse
		if err := c.GetJSON(fmt.Sprintf("/teams/%s/members", teamID), params, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		result := models.ListResponse[models.TeamMember]{Response: resp.Members}
		data, _ := json.Marshal(result)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func createTeamHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		name, ok := getString(args, "name")
		if !ok {
			return mcp.NewToolResultError("name is required"), nil
		}

		team := models.TeamCreate{
			Type: "team",
			Name: name,
		}

		if v, ok := getString(args, "description"); ok {
			team.Description = v
		}

		req := models.TeamCreateRequest{Team: team}

		var resp models.TeamResponse
		if err := c.PostJSON("/teams", req, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, _ := json.Marshal(resp.Team)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func updateTeamHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		teamID, ok := getString(args, "team_id")
		if !ok {
			return mcp.NewToolResultError("team_id is required"), nil
		}

		team := models.TeamUpdate{
			Type: "team",
		}

		if v, ok := getString(args, "name"); ok {
			team.Name = v
		}
		if v, ok := getString(args, "description"); ok {
			team.Description = v
		}

		req := models.TeamUpdateRequest{Team: team}

		var resp models.TeamResponse
		if err := c.PutJSON(fmt.Sprintf("/teams/%s", teamID), req, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, _ := json.Marshal(resp.Team)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func deleteTeamHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		teamID, ok := getString(args, "team_id")
		if !ok {
			return mcp.NewToolResultError("team_id is required"), nil
		}

		if _, err := c.Delete(fmt.Sprintf("/teams/%s", teamID)); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Team %s deleted successfully", teamID)), nil
	}
}

func addTeamMemberHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		teamID, ok := getString(args, "team_id")
		if !ok {
			return mcp.NewToolResultError("team_id is required"), nil
		}

		userID, ok := getString(args, "user_id")
		if !ok {
			return mcp.NewToolResultError("user_id is required"), nil
		}

		member := models.TeamMemberAdd{}
		if v, ok := getString(args, "role"); ok {
			member.Role = v
		}

		if _, err := c.Put(fmt.Sprintf("/teams/%s/users/%s", teamID, userID), member); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("User %s added to team %s", userID, teamID)), nil
	}
}

func removeTeamMemberHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		teamID, ok := getString(args, "team_id")
		if !ok {
			return mcp.NewToolResultError("team_id is required"), nil
		}

		userID, ok := getString(args, "user_id")
		if !ok {
			return mcp.NewToolResultError("user_id is required"), nil
		}

		if _, err := c.Delete(fmt.Sprintf("/teams/%s/users/%s", teamID, userID)); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("User %s removed from team %s", userID, teamID)), nil
	}
}
