package models

import "fmt"

// EventOrchestration represents a PagerDuty event orchestration
type EventOrchestration struct {
	ID           string                        `json:"id,omitempty"`
	Type         string                        `json:"type,omitempty"`
	Self         string                        `json:"self,omitempty"`
	Name         string                        `json:"name"`
	Description  string                        `json:"description,omitempty"`
	Team         *TeamReference                `json:"team,omitempty"`
	Integrations []EventOrchestrationIntegration `json:"integrations,omitempty"`
	Routes       int                           `json:"routes,omitempty"`
	CreatedAt    string                        `json:"created_at,omitempty"`
	CreatedBy    *UserReference                `json:"created_by,omitempty"`
	UpdatedAt    string                        `json:"updated_at,omitempty"`
	UpdatedBy    *UserReference                `json:"updated_by,omitempty"`
}

// EventOrchestrationIntegration represents an integration in an orchestration
type EventOrchestrationIntegration struct {
	ID         string `json:"id"`
	Parameters *IntegrationParameters `json:"parameters,omitempty"`
}

// IntegrationParameters represents parameters for an integration
type IntegrationParameters struct {
	RoutingKey string `json:"routing_key,omitempty"`
	Type       string `json:"type,omitempty"`
}

// EventOrchestrationQuery represents query parameters for listing orchestrations
type EventOrchestrationQuery struct {
	Limit int `json:"limit,omitempty"`
}

// ToParams converts the query to URL parameters
func (q *EventOrchestrationQuery) ToParams() map[string]string {
	params := make(map[string]string)
	if q.Limit > 0 {
		params["limit"] = fmt.Sprintf("%d", q.Limit)
	}
	return params
}

// EventOrchestrationRouter represents the router configuration
type EventOrchestrationRouter struct {
	ID        string                  `json:"id,omitempty"`
	Type      string                  `json:"type,omitempty"`
	Self      string                  `json:"self,omitempty"`
	Parent    *EventOrchestration     `json:"parent,omitempty"`
	Sets      []EventOrchestrationRuleSet `json:"sets,omitempty"`
	CatchAll  *EventOrchestrationCatchAll `json:"catch_all,omitempty"`
}

// EventOrchestrationRuleSet represents a set of rules
type EventOrchestrationRuleSet struct {
	ID    string                    `json:"id"`
	Rules []EventOrchestrationRule  `json:"rules,omitempty"`
}

// EventOrchestrationRule represents a routing rule
type EventOrchestrationRule struct {
	ID         string                         `json:"id,omitempty"`
	Label      string                         `json:"label,omitempty"`
	Conditions []EventOrchestrationRuleCondition `json:"conditions,omitempty"`
	Actions    EventOrchestrationRuleActions  `json:"actions"`
	Disabled   bool                           `json:"disabled,omitempty"`
}

// EventOrchestrationRuleCondition represents a rule condition
type EventOrchestrationRuleCondition struct {
	Expression string `json:"expression"`
}

// EventOrchestrationRuleActions represents rule actions
type EventOrchestrationRuleActions struct {
	RouteTo              string                 `json:"route_to,omitempty"`
	Severity             string                 `json:"severity,omitempty"`
	EventAction          string                 `json:"event_action,omitempty"`
	Variables            []OrchestrationVariable `json:"variables,omitempty"`
	Extractions          []OrchestrationExtraction `json:"extractions,omitempty"`
	DropEvent            bool                   `json:"drop_event,omitempty"`
	Suppress             bool                   `json:"suppress,omitempty"`
	Suspend              *int                   `json:"suspend,omitempty"`
	Priority             string                 `json:"priority,omitempty"`
	Annotate             string                 `json:"annotate,omitempty"`
	PagerDutyAutomationActions []AutomationAction `json:"pagerduty_automation_actions,omitempty"`
	AutomationActions    []AutomationAction     `json:"automation_actions,omitempty"`
	IncidentCustomFieldUpdates []CustomFieldUpdate `json:"incident_custom_field_updates,omitempty"`
}

// OrchestrationVariable represents a variable extraction
type OrchestrationVariable struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	Type  string `json:"type"`
	Value string `json:"value,omitempty"`
}

// OrchestrationExtraction represents a field extraction
type OrchestrationExtraction struct {
	Target   string `json:"target"`
	Template string `json:"template,omitempty"`
	Regex    string `json:"regex,omitempty"`
	Source   string `json:"source,omitempty"`
}

// AutomationAction represents an automation action
type AutomationAction struct {
	ActionID string `json:"action_id"`
}

// CustomFieldUpdate represents a custom field update
type CustomFieldUpdate struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

// EventOrchestrationCatchAll represents catch-all configuration
type EventOrchestrationCatchAll struct {
	Actions EventOrchestrationRuleActions `json:"actions"`
}

// EventOrchestrationGlobal represents global orchestration configuration
type EventOrchestrationGlobal struct {
	ID       string                  `json:"id,omitempty"`
	Type     string                  `json:"type,omitempty"`
	Self     string                  `json:"self,omitempty"`
	Parent   *EventOrchestration     `json:"parent,omitempty"`
	Sets     []EventOrchestrationRuleSet `json:"sets,omitempty"`
	CatchAll *EventOrchestrationCatchAll `json:"catch_all,omitempty"`
}

// EventOrchestrationService represents service orchestration configuration
type EventOrchestrationService struct {
	ID       string                  `json:"id,omitempty"`
	Type     string                  `json:"type,omitempty"`
	Self     string                  `json:"self,omitempty"`
	Parent   *ServiceReference       `json:"parent,omitempty"`
	Sets     []EventOrchestrationRuleSet `json:"sets,omitempty"`
	CatchAll *EventOrchestrationCatchAll `json:"catch_all,omitempty"`
}

// EventOrchestrationRouterUpdateRequest represents a request to update router
type EventOrchestrationRouterUpdateRequest struct {
	OrchestrationPath EventOrchestrationPath `json:"orchestration_path"`
}

// EventOrchestrationPath represents orchestration path for updates
type EventOrchestrationPath struct {
	Type     string                  `json:"type,omitempty"`
	Sets     []EventOrchestrationRuleSet `json:"sets,omitempty"`
	CatchAll *EventOrchestrationCatchAll `json:"catch_all,omitempty"`
}

// EventOrchestrationRuleCreateRequest represents a request to add a rule
type EventOrchestrationRuleCreateRequest struct {
	Label      string                         `json:"label,omitempty"`
	Conditions []EventOrchestrationRuleCondition `json:"conditions,omitempty"`
	Actions    EventOrchestrationRuleActions  `json:"actions"`
	Disabled   bool                           `json:"disabled,omitempty"`
}

// EventOrchestrationResponse is the API response wrapper
type EventOrchestrationResponse struct {
	Orchestration EventOrchestration `json:"orchestration"`
}

// EventOrchestrationsResponse is the API response wrapper for multiple orchestrations
type EventOrchestrationsResponse struct {
	Orchestrations []EventOrchestration `json:"orchestrations"`
	Offset         int                  `json:"offset"`
	Limit          int                  `json:"limit"`
	More           bool                 `json:"more"`
	Total          int                  `json:"total"`
}

// EventOrchestrationRouterResponse is the API response for router
type EventOrchestrationRouterResponse struct {
	OrchestrationPath EventOrchestrationRouter `json:"orchestration_path"`
}

// EventOrchestrationGlobalResponse is the API response for global config
type EventOrchestrationGlobalResponse struct {
	OrchestrationPath EventOrchestrationGlobal `json:"orchestration_path"`
}

// EventOrchestrationServiceResponse is the API response for service config
type EventOrchestrationServiceResponse struct {
	OrchestrationPath EventOrchestrationService `json:"orchestration_path"`
}
