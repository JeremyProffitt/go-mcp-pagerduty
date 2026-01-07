package models

import "fmt"

// Service represents a PagerDuty service
type Service struct {
	ID                     string                     `json:"id,omitempty"`
	Type                   string                     `json:"type,omitempty"`
	Name                   string                     `json:"name"`
	Description            string                     `json:"description,omitempty"`
	Status                 string                     `json:"status,omitempty"`
	EscalationPolicy       *EscalationPolicyReference `json:"escalation_policy,omitempty"`
	Teams                  []TeamReference            `json:"teams,omitempty"`
	Integrations           []IntegrationReference     `json:"integrations,omitempty"`
	IncidentUrgencyRule    *IncidentUrgencyRule       `json:"incident_urgency_rule,omitempty"`
	SupportHours           *SupportHours              `json:"support_hours,omitempty"`
	ScheduledActions       []ScheduledAction          `json:"scheduled_actions,omitempty"`
	AutoResolveTimeout     *int                       `json:"auto_resolve_timeout,omitempty"`
	AcknowledgementTimeout *int                       `json:"acknowledgement_timeout,omitempty"`
	AlertCreation          string                     `json:"alert_creation,omitempty"`
	AlertGrouping          string                     `json:"alert_grouping,omitempty"`
	AlertGroupingTimeout   *int                       `json:"alert_grouping_timeout,omitempty"`
	Self                   string                     `json:"self,omitempty"`
	HTMLURL                string                     `json:"html_url,omitempty"`
	CreatedAt              string                     `json:"created_at,omitempty"`
	UpdatedAt              string                     `json:"updated_at,omitempty"`
}

// IncidentUrgencyRule defines urgency rules for a service
type IncidentUrgencyRule struct {
	Type                string       `json:"type"`
	Urgency             string       `json:"urgency,omitempty"`
	DuringSupportHours  *UrgencyType `json:"during_support_hours,omitempty"`
	OutsideSupportHours *UrgencyType `json:"outside_support_hours,omitempty"`
}

// UrgencyType defines an urgency setting
type UrgencyType struct {
	Type    string `json:"type"`
	Urgency string `json:"urgency"`
}

// SupportHours defines support hours for a service
type SupportHours struct {
	Type      string `json:"type"`
	TimeZone  string `json:"time_zone"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	DaysOfWeek []int  `json:"days_of_week"`
}

// ScheduledAction defines a scheduled action for a service
type ScheduledAction struct {
	Type      string         `json:"type"`
	At        ScheduledAt    `json:"at"`
	ToUrgency string         `json:"to_urgency"`
}

// ScheduledAt defines when a scheduled action occurs
type ScheduledAt struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

// ServiceQuery represents query parameters for listing services
type ServiceQuery struct {
	Query    string   `json:"query,omitempty"`
	TeamIDs  []string `json:"team_ids,omitempty"`
	Limit    int      `json:"limit,omitempty"`
	Includes []string `json:"include,omitempty"`
}

// ToParams converts the query to URL parameters
func (q *ServiceQuery) ToParams() map[string]string {
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
func (q *ServiceQuery) ToArrayParams() map[string][]string {
	params := make(map[string][]string)
	if q.Query != "" {
		params["query"] = []string{q.Query}
	}
	if len(q.TeamIDs) > 0 {
		params["team_ids[]"] = q.TeamIDs
	}
	if len(q.Includes) > 0 {
		params["include[]"] = q.Includes
	}
	if q.Limit > 0 {
		params["limit"] = []string{fmt.Sprintf("%d", q.Limit)}
	}
	return params
}

// ServiceCreateRequest represents a request to create a service
type ServiceCreateRequest struct {
	Service ServiceCreate `json:"service"`
}

// ServiceCreate represents the data for creating a service
type ServiceCreate struct {
	Type             string                    `json:"type"`
	Name             string                    `json:"name"`
	Description      string                    `json:"description,omitempty"`
	EscalationPolicy EscalationPolicyReference `json:"escalation_policy"`
}

// ServiceUpdateRequest represents a request to update a service
type ServiceUpdateRequest struct {
	Service ServiceUpdate `json:"service"`
}

// ServiceUpdate represents the data for updating a service
type ServiceUpdate struct {
	Type             string                     `json:"type"`
	Name             string                     `json:"name,omitempty"`
	Description      string                     `json:"description,omitempty"`
	EscalationPolicy *EscalationPolicyReference `json:"escalation_policy,omitempty"`
}

// ServiceResponse is the API response wrapper for a single service
type ServiceResponse struct {
	Service Service `json:"service"`
}

// ServicesResponse is the API response wrapper for multiple services
type ServicesResponse struct {
	Services []Service `json:"services"`
	Offset   int       `json:"offset"`
	Limit    int       `json:"limit"`
	More     bool      `json:"more"`
	Total    int       `json:"total"`
}
