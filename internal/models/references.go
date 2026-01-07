package models

// UserReference represents a reference to a user
type UserReference struct {
	ID      string `json:"id"`
	Type    string `json:"type,omitempty"`
	Summary string `json:"summary,omitempty"`
	Self    string `json:"self,omitempty"`
	HTMLURL string `json:"html_url,omitempty"`
}

// TeamReference represents a reference to a team
type TeamReference struct {
	ID      string `json:"id"`
	Type    string `json:"type,omitempty"`
	Summary string `json:"summary,omitempty"`
	Self    string `json:"self,omitempty"`
	HTMLURL string `json:"html_url,omitempty"`
}

// ServiceReference represents a reference to a service
type ServiceReference struct {
	ID      string `json:"id"`
	Type    string `json:"type,omitempty"`
	Summary string `json:"summary,omitempty"`
	Self    string `json:"self,omitempty"`
	HTMLURL string `json:"html_url,omitempty"`
}

// ScheduleReference represents a reference to a schedule
type ScheduleReference struct {
	ID      string `json:"id"`
	Type    string `json:"type,omitempty"`
	Summary string `json:"summary,omitempty"`
	Self    string `json:"self,omitempty"`
	HTMLURL string `json:"html_url,omitempty"`
}

// EscalationPolicyReference represents a reference to an escalation policy
type EscalationPolicyReference struct {
	ID      string `json:"id"`
	Type    string `json:"type,omitempty"`
	Summary string `json:"summary,omitempty"`
	Self    string `json:"self,omitempty"`
	HTMLURL string `json:"html_url,omitempty"`
}

// IncidentReference represents a reference to an incident
type IncidentReference struct {
	ID      string `json:"id"`
	Type    string `json:"type,omitempty"`
	Summary string `json:"summary,omitempty"`
	Self    string `json:"self,omitempty"`
	HTMLURL string `json:"html_url,omitempty"`
}

// IntegrationReference represents a reference to an integration
type IntegrationReference struct {
	ID      string `json:"id"`
	Type    string `json:"type,omitempty"`
	Summary string `json:"summary,omitempty"`
	Self    string `json:"self,omitempty"`
	HTMLURL string `json:"html_url,omitempty"`
}

// PriorityReference represents a reference to a priority
type PriorityReference struct {
	ID      string `json:"id"`
	Type    string `json:"type,omitempty"`
	Summary string `json:"summary,omitempty"`
	Self    string `json:"self,omitempty"`
	HTMLURL string `json:"html_url,omitempty"`
}
