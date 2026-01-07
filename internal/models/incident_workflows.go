package models

import "fmt"

// IncidentWorkflow represents a PagerDuty incident workflow
type IncidentWorkflow struct {
	ID          string          `json:"id,omitempty"`
	Type        string          `json:"type,omitempty"`
	Self        string          `json:"self,omitempty"`
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Team        *TeamReference  `json:"team,omitempty"`
	Steps       []Step          `json:"steps,omitempty"`
	CreatedAt   string          `json:"created_at,omitempty"`
	UpdatedAt   string          `json:"updated_at,omitempty"`
}

// Step represents a workflow step
type Step struct {
	ID            string              `json:"id,omitempty"`
	Name          string              `json:"name,omitempty"`
	Type          string              `json:"type,omitempty"`
	Configuration *ActionConfiguration `json:"configuration,omitempty"`
}

// ActionConfiguration represents step configuration
type ActionConfiguration struct {
	ActionID     string        `json:"action_id,omitempty"`
	Description  string        `json:"description,omitempty"`
	InlineStepsInputs []InlineStepInput `json:"inline_steps_inputs,omitempty"`
	Inputs       []ActionInput `json:"inputs,omitempty"`
	Outputs      []ActionOutput `json:"outputs,omitempty"`
}

// InlineStepInput represents an inline step input
type InlineStepInput struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ActionInput represents an action input
type ActionInput struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ActionOutput represents an action output
type ActionOutput struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// IncidentWorkflowQuery represents query parameters for listing workflows
type IncidentWorkflowQuery struct {
	Query    string `json:"query,omitempty"`
	Limit    int    `json:"limit,omitempty"`
	Includes []string `json:"include,omitempty"`
}

// ToParams converts the query to URL parameters
func (q *IncidentWorkflowQuery) ToParams() map[string]string {
	params := make(map[string]string)
	if q.Query != "" {
		params["query"] = q.Query
	}
	if q.Limit > 0 {
		params["limit"] = fmt.Sprintf("%d", q.Limit)
	}
	return params
}

// ToArrayParams converts the query to URL parameters with arrays
func (q *IncidentWorkflowQuery) ToArrayParams() map[string][]string {
	params := make(map[string][]string)
	if len(q.Includes) > 0 {
		params["include[]"] = q.Includes
	}
	for k, v := range q.ToParams() {
		params[k] = []string{v}
	}
	return params
}

// IncidentWorkflowInstance represents a workflow instance
type IncidentWorkflowInstance struct {
	ID        string            `json:"id,omitempty"`
	Type      string            `json:"type,omitempty"`
	Workflow  *IncidentWorkflow `json:"workflow,omitempty"`
	Incident  *IncidentReference `json:"incident,omitempty"`
	Status    string            `json:"status,omitempty"`
	CreatedAt string            `json:"created_at,omitempty"`
	UpdatedAt string            `json:"updated_at,omitempty"`
}

// IncidentWorkflowInstanceRequest represents a request to start a workflow
type IncidentWorkflowInstanceRequest struct {
	IncidentWorkflowInstance IncidentWorkflowInstanceCreate `json:"incident_workflow_instance"`
}

// IncidentWorkflowInstanceCreate represents data to create a workflow instance
type IncidentWorkflowInstanceCreate struct {
	Incident IncidentReference `json:"incident"`
	Workflow WorkflowReference `json:"workflow,omitempty"`
}

// WorkflowReference represents a workflow reference
type WorkflowReference struct {
	ID   string `json:"id"`
	Type string `json:"type,omitempty"`
}

// IncidentWorkflowResponse is the API response wrapper
type IncidentWorkflowResponse struct {
	IncidentWorkflow IncidentWorkflow `json:"incident_workflow"`
}

// IncidentWorkflowsResponse is the API response wrapper for multiple workflows
type IncidentWorkflowsResponse struct {
	IncidentWorkflows []IncidentWorkflow `json:"incident_workflows"`
	Offset            int                `json:"offset"`
	Limit             int                `json:"limit"`
	More              bool               `json:"more"`
	Total             int                `json:"total"`
}

// IncidentWorkflowInstanceResponse is the API response for workflow instance
type IncidentWorkflowInstanceResponse struct {
	IncidentWorkflowInstance IncidentWorkflowInstance `json:"incident_workflow_instance"`
}
