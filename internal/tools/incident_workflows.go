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

// RegisterIncidentWorkflowReadTools registers read-only incident workflow tools
func RegisterIncidentWorkflowReadTools(s *server.MCPServer, c *client.Client) {
	// list_incident_workflows
	s.AddTool(mcp.NewTool("list_incident_workflows",
		mcp.WithDescription("List incident workflows available in PagerDuty. Incident workflows are automated sequences of actions that can be triggered on incidents, such as creating Slack channels, sending notifications, or running diagnostics."),
		mcp.WithTitleAnnotation("List Incident Workflows"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithString("query", mcp.Description("Filter workflows by name (partial match supported)")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results to return"), mcp.Min(1), mcp.Max(100)),
	), listIncidentWorkflowsHandler(c))

	// get_incident_workflow
	s.AddTool(mcp.NewTool("get_incident_workflow",
		mcp.WithDescription("Get detailed information about a specific incident workflow, including its trigger conditions and configured actions."),
		mcp.WithTitleAnnotation("Get Incident Workflow"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithString("workflow_id", mcp.Required(), mcp.Description("The unique workflow ID (e.g., 'PWFLOW123')")),
	), getIncidentWorkflowHandler(c))
}

// RegisterIncidentWorkflowWriteTools registers write incident workflow tools
func RegisterIncidentWorkflowWriteTools(s *server.MCPServer, c *client.Client) {
	// start_incident_workflow
	s.AddTool(mcp.NewTool("start_incident_workflow",
		mcp.WithDescription("Manually trigger an incident workflow on a specific incident. The workflow will execute its configured actions (e.g., create war room, notify stakeholders, run diagnostics). Workflows can also trigger automatically based on incident conditions."),
		mcp.WithTitleAnnotation("Start Incident Workflow"),
		mcp.WithString("workflow_id", mcp.Required(), mcp.Description("The unique workflow ID to execute (e.g., 'PWFLOW123')")),
		mcp.WithString("incident_id", mcp.Required(), mcp.Description("The incident ID to run the workflow on (e.g., 'PABC123')")),
	), startIncidentWorkflowHandler(c))
}

func listIncidentWorkflowsHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		params := make(map[string]string)

		if v, ok := getString(args, "query"); ok {
			params["query"] = v
		}
		if v, ok := getNumber(args, "limit"); ok {
			params["limit"] = fmt.Sprintf("%d", int(v))
		}

		var resp models.IncidentWorkflowsResponse
		if err := c.GetJSON("/incident_workflows", params, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		result := models.ListResponse[models.IncidentWorkflow]{Response: resp.IncidentWorkflows}
		data, _ := json.Marshal(result)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func getIncidentWorkflowHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		workflowID, ok := getString(args, "workflow_id")
		if !ok {
			return mcp.NewToolResultError("workflow_id is required"), nil
		}

		var resp models.IncidentWorkflowResponse
		if err := c.GetJSON(fmt.Sprintf("/incident_workflows/%s", workflowID), nil, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, _ := json.Marshal(resp.IncidentWorkflow)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func startIncidentWorkflowHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		workflowID, ok := getString(args, "workflow_id")
		if !ok {
			return mcp.NewToolResultError("workflow_id is required"), nil
		}

		incidentID, ok := getString(args, "incident_id")
		if !ok {
			return mcp.NewToolResultError("incident_id is required"), nil
		}

		req := models.IncidentWorkflowInstanceRequest{
			IncidentWorkflowInstance: models.IncidentWorkflowInstanceCreate{
				Incident: models.IncidentReference{
					ID:   incidentID,
					Type: "incident_reference",
				},
				Workflow: models.WorkflowReference{
					ID:   workflowID,
					Type: "incident_workflow_reference",
				},
			},
		}

		var resp models.IncidentWorkflowInstanceResponse
		if err := c.PostJSON(fmt.Sprintf("/incident_workflows/%s/instances", workflowID), req, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, _ := json.Marshal(resp.IncidentWorkflowInstance)
		return mcp.NewToolResultText(string(data)), nil
	}
}
