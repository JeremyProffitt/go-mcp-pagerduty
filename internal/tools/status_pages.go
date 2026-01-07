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

// RegisterStatusPageReadTools registers read-only status page tools
func RegisterStatusPageReadTools(s *server.MCPServer, c *client.Client) {
	// list_status_pages
	s.AddTool(mcp.NewTool("list_status_pages",
		mcp.WithDescription("List all public status pages. Status pages communicate service availability to external stakeholders and customers. They can display incidents, maintenance windows, and service health."),
		mcp.WithTitleAnnotation("List Status Pages"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results to return"), mcp.Min(1), mcp.Max(100)),
	), listStatusPagesHandler(c))

	// list_status_page_severities
	s.AddTool(mcp.NewTool("list_status_page_severities",
		mcp.WithDescription("List available severity levels for a status page. Severity options are configurable per status page and determine how critical a post appears to viewers."),
		mcp.WithTitleAnnotation("List Status Page Severities"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithString("status_page_id", mcp.Required(), mcp.Description("The unique status page ID")),
	), listStatusPageSeveritiesHandler(c))

	// list_status_page_impacts
	s.AddTool(mcp.NewTool("list_status_page_impacts",
		mcp.WithDescription("List available impact levels for a status page. Impact describes the scope of an incident (e.g., major outage, partial degradation)."),
		mcp.WithTitleAnnotation("List Status Page Impacts"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithString("status_page_id", mcp.Required(), mcp.Description("The unique status page ID")),
	), listStatusPageImpactsHandler(c))

	// list_status_page_statuses
	s.AddTool(mcp.NewTool("list_status_page_statuses",
		mcp.WithDescription("List available status values for a status page. Statuses track the lifecycle of an incident (e.g., investigating, identified, monitoring, resolved)."),
		mcp.WithTitleAnnotation("List Status Page Statuses"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithString("status_page_id", mcp.Required(), mcp.Description("The unique status page ID")),
	), listStatusPageStatusesHandler(c))

	// get_status_page_post
	s.AddTool(mcp.NewTool("get_status_page_post",
		mcp.WithDescription("Get detailed information about a specific status page post (incident or maintenance announcement), including its current status, severity, and timeline."),
		mcp.WithTitleAnnotation("Get Status Page Post"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithString("status_page_id", mcp.Required(), mcp.Description("The unique status page ID")),
		mcp.WithString("post_id", mcp.Required(), mcp.Description("The unique post ID")),
	), getStatusPagePostHandler(c))

	// list_status_page_post_updates
	s.AddTool(mcp.NewTool("list_status_page_post_updates",
		mcp.WithDescription("List all updates/timeline entries for a specific status page post. Updates show the progression of an incident from detection to resolution."),
		mcp.WithTitleAnnotation("List Status Page Post Updates"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithString("status_page_id", mcp.Required(), mcp.Description("The unique status page ID")),
		mcp.WithString("post_id", mcp.Required(), mcp.Description("The unique post ID")),
	), listStatusPagePostUpdatesHandler(c))
}

// RegisterStatusPageWriteTools registers write status page tools
func RegisterStatusPageWriteTools(s *server.MCPServer, c *client.Client) {
	// create_status_page_post
	s.AddTool(mcp.NewTool("create_status_page_post",
		mcp.WithDescription("Create a new incident or maintenance post on a public status page. This publicly announces an issue or planned maintenance to customers and stakeholders. Use list_status_page_severities and list_status_page_statuses to get valid IDs."),
		mcp.WithTitleAnnotation("Create Status Page Post"),
		mcp.WithString("status_page_id", mcp.Required(), mcp.Description("The unique status page ID")),
		mcp.WithString("post_type", mcp.Required(), mcp.Description("Type of status page post"), mcp.Enum("incident", "maintenance")),
		mcp.WithString("title", mcp.Required(), mcp.Description("Public-facing title describing the incident or maintenance")),
		mcp.WithString("status_id", mcp.Description("Initial status ID (get valid values from list_status_page_statuses)")),
		mcp.WithString("severity_id", mcp.Description("Severity ID (get valid values from list_status_page_severities)")),
		mcp.WithString("starts_at", mcp.Description("Start time in ISO 8601 format (e.g., '2024-01-15T09:00:00Z'). For maintenance, when it begins.")),
		mcp.WithString("ends_at", mcp.Description("End time in ISO 8601 format (e.g., '2024-01-15T11:00:00Z'). For maintenance, expected completion.")),
	), createStatusPagePostHandler(c))

	// create_status_page_post_update
	s.AddTool(mcp.NewTool("create_status_page_post_update",
		mcp.WithDescription("Add a public update to an existing status page post. Updates keep stakeholders informed about incident progress. Can optionally change status/severity and notify subscribers."),
		mcp.WithTitleAnnotation("Add Status Page Update"),
		mcp.WithString("status_page_id", mcp.Required(), mcp.Description("The unique status page ID")),
		mcp.WithString("post_id", mcp.Required(), mcp.Description("The unique post ID to update")),
		mcp.WithString("message", mcp.Required(), mcp.Description("Public-facing update message describing current state or actions taken")),
		mcp.WithString("status_id", mcp.Description("New status ID to transition to (get valid values from list_status_page_statuses)")),
		mcp.WithString("severity_id", mcp.Description("New severity ID if severity has changed")),
		mcp.WithBoolean("notify_subscribers", mcp.Description("Send notification to subscribers about this update (default: false)")),
	), createStatusPagePostUpdateHandler(c))
}

func listStatusPagesHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		params := make(map[string]string)

		if v, ok := getNumber(args, "limit"); ok {
			params["limit"] = fmt.Sprintf("%d", int(v))
		}

		var resp models.StatusPagesResponse
		if err := c.GetJSON("/status_pages", params, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		result := models.ListResponse[models.StatusPage]{Response: resp.StatusPages}
		data, _ := json.Marshal(result)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func listStatusPageSeveritiesHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		statusPageID, ok := getString(args, "status_page_id")
		if !ok {
			return mcp.NewToolResultError("status_page_id is required"), nil
		}

		var resp models.StatusPageSeveritiesResponse
		if err := c.GetJSON(fmt.Sprintf("/status_pages/%s/severities", statusPageID), nil, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		result := models.ListResponse[models.StatusPageSeverity]{Response: resp.Severities}
		data, _ := json.Marshal(result)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func listStatusPageImpactsHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		statusPageID, ok := getString(args, "status_page_id")
		if !ok {
			return mcp.NewToolResultError("status_page_id is required"), nil
		}

		var resp models.StatusPageImpactsResponse
		if err := c.GetJSON(fmt.Sprintf("/status_pages/%s/impacts", statusPageID), nil, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		result := models.ListResponse[models.StatusPageImpact]{Response: resp.Impacts}
		data, _ := json.Marshal(result)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func listStatusPageStatusesHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		statusPageID, ok := getString(args, "status_page_id")
		if !ok {
			return mcp.NewToolResultError("status_page_id is required"), nil
		}

		var resp models.StatusPageStatusesResponse
		if err := c.GetJSON(fmt.Sprintf("/status_pages/%s/statuses", statusPageID), nil, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		result := models.ListResponse[models.StatusPageStatus]{Response: resp.Statuses}
		data, _ := json.Marshal(result)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func getStatusPagePostHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		statusPageID, ok := getString(args, "status_page_id")
		if !ok {
			return mcp.NewToolResultError("status_page_id is required"), nil
		}

		postID, ok := getString(args, "post_id")
		if !ok {
			return mcp.NewToolResultError("post_id is required"), nil
		}

		var resp models.StatusPagePostResponse
		if err := c.GetJSON(fmt.Sprintf("/status_pages/%s/posts/%s", statusPageID, postID), nil, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, _ := json.Marshal(resp.Post)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func listStatusPagePostUpdatesHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		statusPageID, ok := getString(args, "status_page_id")
		if !ok {
			return mcp.NewToolResultError("status_page_id is required"), nil
		}

		postID, ok := getString(args, "post_id")
		if !ok {
			return mcp.NewToolResultError("post_id is required"), nil
		}

		var resp models.StatusPagePostUpdatesResponse
		if err := c.GetJSON(fmt.Sprintf("/status_pages/%s/posts/%s/post_updates", statusPageID, postID), nil, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		result := models.ListResponse[models.StatusPagePostUpdate]{Response: resp.PostUpdates}
		data, _ := json.Marshal(result)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func createStatusPagePostHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		statusPageID, ok := getString(args, "status_page_id")
		if !ok {
			return mcp.NewToolResultError("status_page_id is required"), nil
		}

		postType, ok := getString(args, "post_type")
		if !ok {
			return mcp.NewToolResultError("post_type is required"), nil
		}

		title, ok := getString(args, "title")
		if !ok {
			return mcp.NewToolResultError("title is required"), nil
		}

		post := models.StatusPagePostCreate{
			Type:     "status_page_post",
			PostType: postType,
			Title:    title,
		}

		if v, ok := getString(args, "status_id"); ok {
			post.Status = &models.StatusPageStatusReference{
				ID:   v,
				Type: "status_page_status_reference",
			}
		}
		if v, ok := getString(args, "severity_id"); ok {
			post.Severity = &models.StatusPageSeverityReference{
				ID:   v,
				Type: "status_page_severity_reference",
			}
		}
		if v, ok := getString(args, "starts_at"); ok {
			post.StartsAt = v
		}
		if v, ok := getString(args, "ends_at"); ok {
			post.EndsAt = v
		}

		req := models.StatusPagePostCreateRequestWrapper{Post: post}

		var resp models.StatusPagePostResponse
		if err := c.PostJSON(fmt.Sprintf("/status_pages/%s/posts", statusPageID), req, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, _ := json.Marshal(resp.Post)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func createStatusPagePostUpdateHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		statusPageID, ok := getString(args, "status_page_id")
		if !ok {
			return mcp.NewToolResultError("status_page_id is required"), nil
		}

		postID, ok := getString(args, "post_id")
		if !ok {
			return mcp.NewToolResultError("post_id is required"), nil
		}

		message, ok := getString(args, "message")
		if !ok {
			return mcp.NewToolResultError("message is required"), nil
		}

		update := models.StatusPagePostUpdateCreate{
			Type:    "status_page_post_update",
			Message: message,
		}

		if v, ok := getString(args, "status_id"); ok {
			update.Status = &models.StatusPageStatusReference{
				ID:   v,
				Type: "status_page_status_reference",
			}
		}
		if v, ok := getString(args, "severity_id"); ok {
			update.Severity = &models.StatusPageSeverityReference{
				ID:   v,
				Type: "status_page_severity_reference",
			}
		}
		if v, ok := getBool(args, "notify_subscribers"); ok {
			update.NotifySubscribers = v
		}

		req := models.StatusPagePostUpdateRequestWrapper{PostUpdate: update}

		var resp models.StatusPagePostUpdateResponse
		if err := c.PostJSON(fmt.Sprintf("/status_pages/%s/posts/%s/post_updates", statusPageID, postID), req, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, _ := json.Marshal(resp.PostUpdate)
		return mcp.NewToolResultText(string(data)), nil
	}
}
