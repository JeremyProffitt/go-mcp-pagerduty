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

// RegisterEventOrchestrationReadTools registers read-only event orchestration tools
func RegisterEventOrchestrationReadTools(s *server.MCPServer, c *client.Client) {
	// list_event_orchestrations
	s.AddTool(mcp.NewTool("list_event_orchestrations",
		mcp.WithDescription("List event orchestrations (also called Event Rules). Event orchestrations process incoming events and route them to services based on rules. They can transform, enrich, suppress, or deduplicate events before creating incidents."),
		mcp.WithTitleAnnotation("List Event Orchestrations"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results to return"), mcp.Min(1), mcp.Max(100)),
	), listEventOrchestrationsHandler(c))

	// get_event_orchestration
	s.AddTool(mcp.NewTool("get_event_orchestration",
		mcp.WithDescription("Get basic information about a specific event orchestration, including its integration URL for receiving events."),
		mcp.WithTitleAnnotation("Get Event Orchestration"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithString("orchestration_id", mcp.Required(), mcp.Description("The unique orchestration ID (e.g., 'E1A2B3C')")),
	), getEventOrchestrationHandler(c))

	// get_event_orchestration_router
	s.AddTool(mcp.NewTool("get_event_orchestration_router",
		mcp.WithDescription("Get the router rules for an event orchestration. Router rules determine which service an event is routed to based on conditions matching event fields."),
		mcp.WithTitleAnnotation("Get Orchestration Router Rules"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithString("orchestration_id", mcp.Required(), mcp.Description("The unique orchestration ID (e.g., 'E1A2B3C')")),
	), getEventOrchestrationRouterHandler(c))

	// get_event_orchestration_global
	s.AddTool(mcp.NewTool("get_event_orchestration_global",
		mcp.WithDescription("Get the global orchestration rules that apply before routing. Global rules can suppress, deduplicate, or transform events before they reach services."),
		mcp.WithTitleAnnotation("Get Global Orchestration Rules"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithString("orchestration_id", mcp.Required(), mcp.Description("The unique orchestration ID (e.g., 'E1A2B3C')")),
	), getEventOrchestrationGlobalHandler(c))

	// get_event_orchestration_service
	s.AddTool(mcp.NewTool("get_event_orchestration_service",
		mcp.WithDescription("Get the service-level orchestration rules for a specific service. These rules process events after routing and can set severity, add notes, or trigger automations."),
		mcp.WithTitleAnnotation("Get Service Orchestration Rules"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithString("service_id", mcp.Required(), mcp.Description("The unique service ID (e.g., 'PDSVC123')")),
	), getEventOrchestrationServiceHandler(c))
}

// RegisterEventOrchestrationWriteTools registers write event orchestration tools
func RegisterEventOrchestrationWriteTools(s *server.MCPServer, c *client.Client) {
	// update_event_orchestration_router
	s.AddTool(mcp.NewTool("update_event_orchestration_router",
		mcp.WithDescription("Replace the entire router configuration for an event orchestration. This completely overwrites existing rules. For adding a single rule, use append_event_orchestration_router_rule instead."),
		mcp.WithTitleAnnotation("Update Orchestration Router"),
		mcp.WithString("orchestration_id", mcp.Required(), mcp.Description("The unique orchestration ID (e.g., 'E1A2B3C')")),
		mcp.WithString("config", mcp.Required(), mcp.Description("Complete router configuration as JSON. Must include 'orchestration_path' with 'sets' and 'catch_all' fields.")),
	), updateEventOrchestrationRouterHandler(c))

	// append_event_orchestration_router_rule
	s.AddTool(mcp.NewTool("append_event_orchestration_router_rule",
		mcp.WithDescription("Add a new routing rule to an event orchestration without modifying existing rules. The rule will be appended to the first rule set. Use this for safely adding new routing logic."),
		mcp.WithTitleAnnotation("Add Router Rule"),
		mcp.WithString("orchestration_id", mcp.Required(), mcp.Description("The unique orchestration ID (e.g., 'E1A2B3C')")),
		mcp.WithString("label", mcp.Description("Human-readable label for the rule (e.g., 'Route database alerts')")),
		mcp.WithString("conditions", mcp.Description("JSON array of conditions. Each condition has 'expression' (JEXL format, e.g., 'event.source matches \"database\"')")),
		mcp.WithString("route_to", mcp.Required(), mcp.Description("The service ID to route matching events to (e.g., 'PDSVC123')")),
	), appendEventOrchestrationRouterRuleHandler(c))
}

func listEventOrchestrationsHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		params := make(map[string]string)

		if v, ok := getNumber(args, "limit"); ok {
			params["limit"] = fmt.Sprintf("%d", int(v))
		}

		var resp models.EventOrchestrationsResponse
		if err := c.GetJSON("/event_orchestrations", params, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		result := models.ListResponse[models.EventOrchestration]{Response: resp.Orchestrations}
		data, _ := json.Marshal(result)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func getEventOrchestrationHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		orchestrationID, ok := getString(args, "orchestration_id")
		if !ok {
			return mcp.NewToolResultError("orchestration_id is required"), nil
		}

		var resp models.EventOrchestrationResponse
		if err := c.GetJSON(fmt.Sprintf("/event_orchestrations/%s", orchestrationID), nil, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, _ := json.Marshal(resp.Orchestration)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func getEventOrchestrationRouterHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		orchestrationID, ok := getString(args, "orchestration_id")
		if !ok {
			return mcp.NewToolResultError("orchestration_id is required"), nil
		}

		var resp models.EventOrchestrationRouterResponse
		if err := c.GetJSON(fmt.Sprintf("/event_orchestrations/%s/router", orchestrationID), nil, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, _ := json.Marshal(resp.OrchestrationPath)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func getEventOrchestrationGlobalHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		orchestrationID, ok := getString(args, "orchestration_id")
		if !ok {
			return mcp.NewToolResultError("orchestration_id is required"), nil
		}

		var resp models.EventOrchestrationGlobalResponse
		if err := c.GetJSON(fmt.Sprintf("/event_orchestrations/%s/global", orchestrationID), nil, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, _ := json.Marshal(resp.OrchestrationPath)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func getEventOrchestrationServiceHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		serviceID, ok := getString(args, "service_id")
		if !ok {
			return mcp.NewToolResultError("service_id is required"), nil
		}

		var resp models.EventOrchestrationServiceResponse
		if err := c.GetJSON(fmt.Sprintf("/event_orchestrations/services/%s", serviceID), nil, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, _ := json.Marshal(resp.OrchestrationPath)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func updateEventOrchestrationRouterHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		orchestrationID, ok := getString(args, "orchestration_id")
		if !ok {
			return mcp.NewToolResultError("orchestration_id is required"), nil
		}

		configStr, ok := getString(args, "config")
		if !ok {
			return mcp.NewToolResultError("config is required"), nil
		}

		var config models.EventOrchestrationRouterUpdateRequest
		if err := json.Unmarshal([]byte(configStr), &config); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("invalid config JSON: %v", err)), nil
		}

		var resp models.EventOrchestrationRouterResponse
		if err := c.PutJSON(fmt.Sprintf("/event_orchestrations/%s/router", orchestrationID), config, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, _ := json.Marshal(resp.OrchestrationPath)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func appendEventOrchestrationRouterRuleHandler(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getArgs(request)
		orchestrationID, ok := getString(args, "orchestration_id")
		if !ok {
			return mcp.NewToolResultError("orchestration_id is required"), nil
		}

		routeTo, ok := getString(args, "route_to")
		if !ok {
			return mcp.NewToolResultError("route_to is required"), nil
		}

		// First, get the current router config
		var currentResp models.EventOrchestrationRouterResponse
		if err := c.GetJSON(fmt.Sprintf("/event_orchestrations/%s/router", orchestrationID), nil, &currentResp); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to get current router: %v", err)), nil
		}

		// Create the new rule
		newRule := models.EventOrchestrationRule{
			Actions: models.EventOrchestrationRuleActions{
				RouteTo: routeTo,
			},
		}

		if v, ok := getString(args, "label"); ok {
			newRule.Label = v
		}

		if v, ok := getString(args, "conditions"); ok {
			var conditions []models.EventOrchestrationRuleCondition
			if err := json.Unmarshal([]byte(v), &conditions); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("invalid conditions JSON: %v", err)), nil
			}
			newRule.Conditions = conditions
		}

		// Append the new rule to the first set
		if len(currentResp.OrchestrationPath.Sets) > 0 {
			currentResp.OrchestrationPath.Sets[0].Rules = append(currentResp.OrchestrationPath.Sets[0].Rules, newRule)
		}

		// Update the router
		updateReq := models.EventOrchestrationRouterUpdateRequest{
			OrchestrationPath: models.EventOrchestrationPath{
				Sets:     currentResp.OrchestrationPath.Sets,
				CatchAll: currentResp.OrchestrationPath.CatchAll,
			},
		}

		var resp models.EventOrchestrationRouterResponse
		if err := c.PutJSON(fmt.Sprintf("/event_orchestrations/%s/router", orchestrationID), updateReq, &resp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, _ := json.Marshal(resp.OrchestrationPath)
		return mcp.NewToolResultText(string(data)), nil
	}
}
