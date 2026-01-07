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

// RegisterIncidentReadTools registers read-only incident tools
func RegisterIncidentReadTools(s *server.MCPServer, c *client.Client) {
	// list_incidents
	s.AddTool(mcp.NewTool("list_incidents",
		mcp.WithDescription("List incidents from PagerDuty with optional filtering. Use this to find active incidents (triggered/acknowledged), review incident history, or search for incidents affecting specific services or teams. For investigating a specific incident's history, use get_past_incidents instead."),
		mcp.WithTitleAnnotation("List Incidents"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithString("statuses", mcp.Description("Filter by incident status. Comma-separated values (e.g., 'triggered,acknowledged')"), mcp.Enum("triggered", "acknowledged", "resolved")),
		mcp.WithString("date_range", mcp.Description("Predefined date range filter"), mcp.Enum("all", "past_month", "past_week")),
		mcp.WithString("since", mcp.Description("Start date in ISO 8601 format (e.g., '2024-01-15T10:00:00Z'). Use with 'until' for custom date ranges.")),
		mcp.WithString("until", mcp.Description("End date in ISO 8601 format (e.g., '2024-01-15T18:00:00Z'). Use with 'since' for custom date ranges.")),
		mcp.WithString("urgencies", mcp.Description("Filter by urgency level. Comma-separated values (e.g., 'high,low')"), mcp.Enum("high", "low")),
		mcp.WithString("service_ids", mcp.Description("Filter by services. Comma-separated service IDs (e.g., 'PDSVC1,PDSVC2')")),
		mcp.WithString("team_ids", mcp.Description("Filter by teams. Comma-separated team IDs (e.g., 'PTEAM1,PTEAM2')")),
		mcp.WithString("user_ids", mcp.Description("Filter by assigned users. Comma-separated user IDs (e.g., 'PUSER1,PUSER2')")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results to return (default: 20)"), mcp.Min(1), mcp.Max(100)),
	), listIncidentsHandler(c))

	// get_incident
	s.AddTool(mcp.NewTool("get_incident",
		mcp.WithDescription("Get detailed information about a specific incident by ID, including status, assignments, urgency, and timestamps."),
		mcp.WithTitleAnnotation("Get Incident Details"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithString("incident_id", mcp.Required(), mcp.Description("The unique incident ID (e.g., 'PABC123')")),
	), getIncidentHandler(c))

	// get_outlier_incident
	s.AddTool(mcp.NewTool("get_outlier_incident",
		mcp.WithDescription("Analyze if an incident is an outlier compared to historical patterns. Returns machine learning-based analysis of whether this incident is unusual for the service."),
		mcp.WithTitleAnnotation("Get Outlier Analysis"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithString("incident_id", mcp.Required(), mcp.Description("The unique incident ID (e.g., 'PABC123')")),
		mcp.WithString("since", mcp.Description("Start date for historical analysis in ISO 8601 format (e.g., '2024-01-01T00:00:00Z')")),
	), getOutlierIncidentHandler(c))

	// get_past_incidents
	s.AddTool(mcp.NewTool("get_past_incidents",
		mcp.WithDescription("Find similar historical incidents that may help with troubleshooting. Uses machine learning to match incidents based on patterns in alerts and metadata. Different from get_related_incidents which finds concurrent incidents."),
		mcp.WithTitleAnnotation("Get Similar Past Incidents"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithString("incident_id", mcp.Required(), mcp.Description("The unique incident ID (e.g., 'PABC123')")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of past incidents to return (default: 5)"), mcp.Min(1), mcp.Max(100)),
	), getPastIncidentsHandler(c))

	// get_related_incidents
	s.AddTool(mcp.NewTool("get_related_incidents",
		mcp.WithDescription("Find incidents that may be related to a specific incident based on timing and service relationships. Useful for identifying widespread issues. Different from get_past_incidents which finds historical similar incidents."),
		mcp.WithTitleAnnotation("Get Related Incidents"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithString("incident_id", mcp.Required(), mcp.Description("The unique incident ID (e.g., 'PABC123')")),
	), getRelatedIncidentsHandler(c))

	// list_incident_notes
	s.AddTool(mcp.NewTool("list_incident_notes",
		mcp.WithDescription("List all notes/comments added to an incident. Notes contain investigation details, status updates, and resolution information added by responders."),
		mcp.WithTitleAnnotation("List Incident Notes"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithString("incident_id", mcp.Required(), mcp.Description("The unique incident ID (e.g., 'PABC123')")),
	), listIncidentNotesHandler(c))
}

// RegisterIncidentWriteTools registers write incident tools
func RegisterIncidentWriteTools(s *server.MCPServer, c *client.Client) {
	// create_incident
	s.AddTool(mcp.NewTool("create_incident",
		mcp.WithDescription("Create a new incident manually on a service. Use this to report issues that weren't automatically detected by monitoring. The incident will trigger notifications according to the service's escalation policy."),
		mcp.WithTitleAnnotation("Create Incident"),
		mcp.WithString("title", mcp.Required(), mcp.Description("A brief, descriptive title for the incident")),
		mcp.WithString("service_id", mcp.Required(), mcp.Description("The service ID where the incident will be created (e.g., 'PDSVC123')")),
		mcp.WithString("urgency", mcp.Description("Incident urgency level"), mcp.Enum("high", "low")),
		mcp.WithString("body", mcp.Description("Detailed description of the incident including symptoms, impact, and any relevant context")),
		mcp.WithString("incident_key", mcp.Description("Deduplication key to prevent duplicate incidents. Incidents with the same key on the same service will be grouped.")),
	), createIncidentHandler(c))

	// manage_incidents
	s.AddTool(mcp.NewTool("manage_incidents",
		mcp.WithDescription("Bulk update one or more incidents. Use to acknowledge incidents you're working on, resolve incidents that are fixed, change urgency, reassign to other users, or escalate to higher levels. Cannot change status to 'triggered' - use create_incident instead."),
		mcp.WithTitleAnnotation("Manage Incidents"),
		mcp.WithString("incident_ids", mcp.Required(), mcp.Description("Comma-separated incident IDs to update (e.g., 'PABC123,PDEF456')")),
		mcp.WithString("status", mcp.Description("New incident status"), mcp.Enum("acknowledged", "resolved")),
		mcp.WithString("urgency", mcp.Description("New urgency level"), mcp.Enum("high", "low")),
		mcp.WithString("assignee_id", mcp.Description("User ID to assign/reassign the incidents to (e.g., 'PUSER123')")),
		mcp.WithNumber("escalation_level", mcp.Description("Escalation level to set (escalates to users at that level in the escalation policy)"), mcp.Min(1)),
	), manageIncidentsHandler(c))

	// add_responders
	s.AddTool(mcp.NewTool("add_responders",
		mcp.WithDescription("Request additional responders to help with an incident. The specified users will receive notifications asking them to join the incident response."),
		mcp.WithTitleAnnotation("Add Responders"),
		mcp.WithString("incident_id", mcp.Required(), mcp.Description("The unique incident ID (e.g., 'PABC123')")),
		mcp.WithString("responder_ids", mcp.Required(), mcp.Description("Comma-separated user IDs to request as responders (e.g., 'PUSER1,PUSER2')")),
		mcp.WithString("message", mcp.Description("Optional message explaining why these responders are needed")),
	), addRespondersHandler(c))

	// add_note_to_incident
	s.AddTool(mcp.NewTool("add_note_to_incident",
		mcp.WithDescription("Add a note to document investigation progress, findings, or resolution details on an incident. Notes are visible to all responders and preserved in incident history."),
		mcp.WithTitleAnnotation("Add Incident Note"),
		mcp.WithString("incident_id", mcp.Required(), mcp.Description("The unique incident ID (e.g., 'PABC123')")),
		mcp.WithString("note", mcp.Required(), mcp.Description("The note content to add to the incident")),
	), addNoteToIncidentHandler(c))
}

func listIncidentsHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		params := make(map[string]string)

		if v, ok := getString(args, "statuses"); ok {
			params["statuses[]"] = v
		}
		if v, ok := getString(args, "date_range"); ok {
			params["date_range"] = v
		}
		if v, ok := getString(args, "since"); ok {
			params["since"] = v
		}
		if v, ok := getString(args, "until"); ok {
			params["until"] = v
		}
		if v, ok := getString(args, "urgencies"); ok {
			params["urgencies[]"] = v
		}
		if v, ok := getNumber(args, "limit"); ok {
			params["limit"] = fmt.Sprintf("%d", int(v))
		}

		var resp models.IncidentsResponse
		if err := c.GetJSON("/incidents", params, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		result := models.ListResponse[models.Incident]{Response: resp.Incidents}
		data, _ := json.Marshal(result)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func getIncidentHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		incidentID, ok := getString(args, "incident_id")
		if !ok {
			return mcp.NewToolResultError("incident_id is required"), nil
		}

		var resp models.IncidentResponse
		if err := c.GetJSON(fmt.Sprintf("/incidents/%s", incidentID), nil, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, _ := json.Marshal(resp.Incident)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func getOutlierIncidentHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		incidentID, ok := getString(args, "incident_id")
		if !ok {
			return mcp.NewToolResultError("incident_id is required"), nil
		}

		params := make(map[string]string)
		if v, ok := getString(args, "since"); ok {
			params["since"] = v
		}

		var resp models.OutlierIncidentResponse
		if err := c.GetJSON(fmt.Sprintf("/incidents/%s/outlier_incident", incidentID), params, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, _ := json.Marshal(resp)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func getPastIncidentsHandler(c *client.Client) server.ToolHandlerFunc {
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

		var resp models.PastIncidentsResponse
		if err := c.GetJSON(fmt.Sprintf("/incidents/%s/past_incidents", incidentID), params, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, _ := json.Marshal(resp)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func getRelatedIncidentsHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		incidentID, ok := getString(args, "incident_id")
		if !ok {
			return mcp.NewToolResultError("incident_id is required"), nil
		}

		var resp models.RelatedIncidentsResponse
		if err := c.GetJSON(fmt.Sprintf("/incidents/%s/related_incidents", incidentID), nil, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, _ := json.Marshal(resp)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func listIncidentNotesHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		incidentID, ok := getString(args, "incident_id")
		if !ok {
			return mcp.NewToolResultError("incident_id is required"), nil
		}

		var resp models.IncidentNotesResponse
		if err := c.GetJSON(fmt.Sprintf("/incidents/%s/notes", incidentID), nil, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		result := models.ListResponse[models.IncidentNote]{Response: resp.Notes}
		data, _ := json.Marshal(result)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func createIncidentHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		title, ok := getString(args, "title")
		if !ok {
			return mcp.NewToolResultError("title is required"), nil
		}

		serviceID, ok := getString(args, "service_id")
		if !ok {
			return mcp.NewToolResultError("service_id is required"), nil
		}

		incident := models.IncidentCreate{
			Type:  "incident",
			Title: title,
			Service: models.ServiceReference{
				ID:   serviceID,
				Type: "service_reference",
			},
		}

		if v, ok := getString(args, "urgency"); ok {
			incident.Urgency = v
		}
		if v, ok := getString(args, "body"); ok {
			incident.Body = &models.IncidentBody{
				Type:    "incident_body",
				Details: v,
			}
		}
		if v, ok := getString(args, "incident_key"); ok {
			incident.IncidentKey = v
		}

		req := models.IncidentCreateRequest{Incident: incident}

		var resp models.IncidentResponse
		if err := c.PostJSON("/incidents", req, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, _ := json.Marshal(resp.Incident)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func manageIncidentsHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		incidentIDsStr, ok := getString(args, "incident_ids")
		if !ok {
			return mcp.NewToolResultError("incident_ids is required"), nil
		}

		incidentIDs := splitAndTrim(incidentIDsStr)

		manageReq := models.IncidentManageRequest{
			IncidentIDs: incidentIDs,
		}

		if v, ok := getString(args, "status"); ok {
			manageReq.Status = v
		}
		if v, ok := getString(args, "urgency"); ok {
			manageReq.Urgency = v
		}
		if v, ok := getString(args, "assignee_id"); ok {
			manageReq.Assignment = &models.UserReference{ID: v}
		}
		if v, ok := getNumber(args, "escalation_level"); ok {
			manageReq.EscalationLevel = int(v)
		}

		payload := manageReq.ToAPIPayload()

		var resp models.IncidentsResponse
		if err := c.PutJSON("/incidents", payload, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		result := models.ListResponse[models.Incident]{Response: resp.Incidents}
		data, _ := json.Marshal(result)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func addRespondersHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		incidentID, ok := getString(args, "incident_id")
		if !ok {
			return mcp.NewToolResultError("incident_id is required"), nil
		}

		responderIDsStr, ok := getString(args, "responder_ids")
		if !ok {
			return mcp.NewToolResultError("responder_ids is required"), nil
		}

		responderIDs := splitAndTrim(responderIDsStr)
		targets := make([]models.ResponderRequestTarget, len(responderIDs))
		for i, id := range responderIDs {
			targets[i] = models.ResponderRequestTarget{
				Type: "user_reference",
				ID:   id,
			}
		}

		req := models.IncidentResponderRequest{
			Targets: targets,
		}

		if v, ok := getString(args, "message"); ok {
			req.Message = v
		}

		data, err := c.Post(fmt.Sprintf("/incidents/%s/responder_requests", incidentID), req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(string(data)), nil
	}
}

func addNoteToIncidentHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		incidentID, ok := getString(args, "incident_id")
		if !ok {
			return mcp.NewToolResultError("incident_id is required"), nil
		}

		note, ok := getString(args, "note")
		if !ok {
			return mcp.NewToolResultError("note is required"), nil
		}

		req := models.IncidentNoteCreateRequest{
			Note: models.NoteContent{Content: note},
		}

		var resp struct {
			Note models.IncidentNote `json:"note"`
		}
		if err := c.PostJSON(fmt.Sprintf("/incidents/%s/notes", incidentID), req, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, _ := json.Marshal(resp.Note)
		return mcp.NewToolResultText(string(data)), nil
	}
}
