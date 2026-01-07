package models

import "fmt"

// EscalationPolicy represents a PagerDuty escalation policy
type EscalationPolicy struct {
	ID               string              `json:"id,omitempty"`
	Type             string              `json:"type,omitempty"`
	Summary          string              `json:"summary,omitempty"`
	Self             string              `json:"self,omitempty"`
	HTMLURL          string              `json:"html_url,omitempty"`
	Name             string              `json:"name"`
	Description      string              `json:"description,omitempty"`
	NumLoops         int                 `json:"num_loops,omitempty"`
	OnCallHandoffNotifications string    `json:"on_call_handoff_notifications,omitempty"`
	EscalationRules  []EscalationRule    `json:"escalation_rules,omitempty"`
	Services         []ServiceReference  `json:"services,omitempty"`
	Teams            []TeamReference     `json:"teams,omitempty"`
}

// EscalationRule represents a rule in an escalation policy
type EscalationRule struct {
	ID                       string              `json:"id,omitempty"`
	EscalationDelayInMinutes int                 `json:"escalation_delay_in_minutes"`
	Targets                  []EscalationTarget  `json:"targets"`
}

// EscalationTarget represents a target in an escalation rule
type EscalationTarget struct {
	ID   string `json:"id"`
	Type string `json:"type"` // user_reference, schedule_reference
}

// EscalationPolicyQuery represents query parameters for listing escalation policies
type EscalationPolicyQuery struct {
	Query    string   `json:"query,omitempty"`
	UserIDs  []string `json:"user_ids,omitempty"`
	TeamIDs  []string `json:"team_ids,omitempty"`
	Includes []string `json:"include,omitempty"`
	SortBy   string   `json:"sort_by,omitempty"`
	Limit    int      `json:"limit,omitempty"`
}

// ToParams converts the query to URL parameters
func (q *EscalationPolicyQuery) ToParams() map[string]string {
	params := make(map[string]string)
	if q.Query != "" {
		params["query"] = q.Query
	}
	if q.SortBy != "" {
		params["sort_by"] = q.SortBy
	}
	if q.Limit > 0 {
		params["limit"] = fmt.Sprintf("%d", q.Limit)
	}
	return params
}

// ToArrayParams converts the query to URL parameters with arrays
func (q *EscalationPolicyQuery) ToArrayParams() map[string][]string {
	params := make(map[string][]string)
	if len(q.UserIDs) > 0 {
		params["user_ids[]"] = q.UserIDs
	}
	if len(q.TeamIDs) > 0 {
		params["team_ids[]"] = q.TeamIDs
	}
	if len(q.Includes) > 0 {
		params["include[]"] = q.Includes
	}
	// Merge single params
	for k, v := range q.ToParams() {
		params[k] = []string{v}
	}
	return params
}

// EscalationPolicyResponse is the API response wrapper for a single escalation policy
type EscalationPolicyResponse struct {
	EscalationPolicy EscalationPolicy `json:"escalation_policy"`
}

// EscalationPoliciesResponse is the API response wrapper for multiple escalation policies
type EscalationPoliciesResponse struct {
	EscalationPolicies []EscalationPolicy `json:"escalation_policies"`
	Offset             int                `json:"offset"`
	Limit              int                `json:"limit"`
	More               bool               `json:"more"`
	Total              int                `json:"total"`
}
