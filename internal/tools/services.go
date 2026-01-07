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

// RegisterServiceReadTools registers read-only service tools
func RegisterServiceReadTools(s *server.MCPServer, c *client.Client) {
	// list_services
	s.AddTool(mcp.NewTool("list_services",
		mcp.WithDescription("List services (monitored applications/components) in PagerDuty. Services are the entities that receive alerts and generate incidents. Use to find service IDs for filtering incidents or understanding what's being monitored."),
		mcp.WithTitleAnnotation("List Services"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithString("query", mcp.Description("Filter services by name (partial match supported)")),
		mcp.WithString("team_ids", mcp.Description("Filter by owning teams. Comma-separated team IDs (e.g., 'PTEAM1,PTEAM2')")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results to return"), mcp.Min(1), mcp.Max(100)),
	), listServicesHandler(c))

	// get_service
	s.AddTool(mcp.NewTool("get_service",
		mcp.WithDescription("Get detailed information about a specific service including its escalation policy, integrations, and configuration settings."),
		mcp.WithTitleAnnotation("Get Service Details"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithString("service_id", mcp.Required(), mcp.Description("The unique service ID (e.g., 'PDSVC123')")),
	), getServiceHandler(c))
}

// RegisterServiceWriteTools registers write service tools
func RegisterServiceWriteTools(s *server.MCPServer, c *client.Client) {
	// create_service
	s.AddTool(mcp.NewTool("create_service",
		mcp.WithDescription("Create a new service to represent a monitored application or component. Services receive alerts from integrations and generate incidents based on their configuration. An escalation policy is required to define who gets notified."),
		mcp.WithTitleAnnotation("Create Service"),
		mcp.WithString("name", mcp.Required(), mcp.Description("A descriptive name for the service (e.g., 'Production API', 'Payment Gateway')")),
		mcp.WithString("escalation_policy_id", mcp.Required(), mcp.Description("The escalation policy ID that defines notification rules (e.g., 'PESCPOL123')")),
		mcp.WithString("description", mcp.Description("Detailed description of what this service monitors and its business impact")),
	), createServiceHandler(c))

	// update_service
	s.AddTool(mcp.NewTool("update_service",
		mcp.WithDescription("Update an existing service's configuration. Use to rename services, update descriptions, or change the escalation policy."),
		mcp.WithTitleAnnotation("Update Service"),
		mcp.WithString("service_id", mcp.Required(), mcp.Description("The unique service ID to update (e.g., 'PDSVC123')")),
		mcp.WithString("name", mcp.Description("New service name")),
		mcp.WithString("description", mcp.Description("New service description")),
		mcp.WithString("escalation_policy_id", mcp.Description("New escalation policy ID to assign (e.g., 'PESCPOL123')")),
	), updateServiceHandler(c))
}

func listServicesHandler(c *client.Client) server.ToolHandlerFunc {
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

		var resp models.ServicesResponse
		if err := c.GetJSON("/services", params, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		result := models.ListResponse[models.Service]{Response: resp.Services}
		data, _ := json.Marshal(result)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func getServiceHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		serviceID, ok := getString(args, "service_id")
		if !ok {
			return mcp.NewToolResultError("service_id is required"), nil
		}

		var resp models.ServiceResponse
		if err := c.GetJSON(fmt.Sprintf("/services/%s", serviceID), nil, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, _ := json.Marshal(resp.Service)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func createServiceHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		name, ok := getString(args, "name")
		if !ok {
			return mcp.NewToolResultError("name is required"), nil
		}

		escalationPolicyID, ok := getString(args, "escalation_policy_id")
		if !ok {
			return mcp.NewToolResultError("escalation_policy_id is required"), nil
		}

		service := models.ServiceCreate{
			Type: "service",
			Name: name,
			EscalationPolicy: models.EscalationPolicyReference{
				ID:   escalationPolicyID,
				Type: "escalation_policy_reference",
			},
		}

		if v, ok := getString(args, "description"); ok {
			service.Description = v
		}

		req := models.ServiceCreateRequest{Service: service}

		var resp models.ServiceResponse
		if err := c.PostJSON("/services", req, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, _ := json.Marshal(resp.Service)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func updateServiceHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		serviceID, ok := getString(args, "service_id")
		if !ok {
			return mcp.NewToolResultError("service_id is required"), nil
		}

		service := models.ServiceUpdate{
			Type: "service",
		}

		if v, ok := getString(args, "name"); ok {
			service.Name = v
		}
		if v, ok := getString(args, "description"); ok {
			service.Description = v
		}
		if v, ok := getString(args, "escalation_policy_id"); ok {
			service.EscalationPolicy = &models.EscalationPolicyReference{
				ID:   v,
				Type: "escalation_policy_reference",
			}
		}

		req := models.ServiceUpdateRequest{Service: service}

		var resp models.ServiceResponse
		if err := c.PutJSON(fmt.Sprintf("/services/%s", serviceID), req, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, _ := json.Marshal(resp.Service)
		return mcp.NewToolResultText(string(data)), nil
	}
}
