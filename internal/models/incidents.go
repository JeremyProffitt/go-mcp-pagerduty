package models

import (
	"fmt"
	"strings"
	"time"
)

// Incident represents a PagerDuty incident
type Incident struct {
	ID                    string              `json:"id"`
	Type                  string              `json:"type,omitempty"`
	Summary               string              `json:"summary,omitempty"`
	Self                  string              `json:"self,omitempty"`
	HTMLURL               string              `json:"html_url,omitempty"`
	IncidentNumber        int                 `json:"incident_number,omitempty"`
	Title                 string              `json:"title,omitempty"`
	CreatedAt             string              `json:"created_at,omitempty"`
	UpdatedAt             string              `json:"updated_at,omitempty"`
	Status                string              `json:"status,omitempty"`
	IncidentKey           string              `json:"incident_key,omitempty"`
	Service               *ServiceReference   `json:"service,omitempty"`
	Assignments           []Assignment        `json:"assignments,omitempty"`
	Acknowledgements      []Acknowledgement   `json:"acknowledgements,omitempty"`
	LastStatusChangeAt    string              `json:"last_status_change_at,omitempty"`
	LastStatusChangeBy    *UserReference      `json:"last_status_change_by,omitempty"`
	FirstTriggerLogEntry  *LogEntryReference  `json:"first_trigger_log_entry,omitempty"`
	EscalationPolicy      *EscalationPolicyReference `json:"escalation_policy,omitempty"`
	Teams                 []TeamReference     `json:"teams,omitempty"`
	Priority              *PriorityReference  `json:"priority,omitempty"`
	Urgency               string              `json:"urgency,omitempty"`
	ResolveReason         *ResolveReason      `json:"resolve_reason,omitempty"`
	AlertCounts           *AlertCounts        `json:"alert_counts,omitempty"`
	Body                  *IncidentBody       `json:"body,omitempty"`
	IsMergeable           bool                `json:"is_mergeable,omitempty"`
	ConferenceBridge      *ConferenceBridge   `json:"conference_bridge,omitempty"`
}

// Assignment represents an incident assignment
type Assignment struct {
	At       string        `json:"at"`
	Assignee UserReference `json:"assignee"`
}

// Acknowledgement represents an incident acknowledgement
type Acknowledgement struct {
	At           string        `json:"at"`
	Acknowledger UserReference `json:"acknowledger"`
}

// LogEntryReference represents a reference to a log entry
type LogEntryReference struct {
	ID      string `json:"id"`
	Type    string `json:"type,omitempty"`
	Summary string `json:"summary,omitempty"`
	Self    string `json:"self,omitempty"`
	HTMLURL string `json:"html_url,omitempty"`
}

// ResolveReason represents the reason an incident was resolved
type ResolveReason struct {
	Type     string         `json:"type"`
	Incident *IncidentReference `json:"incident,omitempty"`
}

// AlertCounts represents alert counts for an incident
type AlertCounts struct {
	All       int `json:"all"`
	Triggered int `json:"triggered"`
	Resolved  int `json:"resolved"`
}

// IncidentBody represents the body of an incident
type IncidentBody struct {
	Type    string `json:"type"`
	Details string `json:"details,omitempty"`
}

// ConferenceBridge represents conference bridge info
type ConferenceBridge struct {
	ConferenceNumber string `json:"conference_number,omitempty"`
	ConferenceURL    string `json:"conference_url,omitempty"`
}

// IncidentQuery represents query parameters for listing incidents
type IncidentQuery struct {
	Statuses     []string `json:"statuses,omitempty"`
	DateRange    string   `json:"date_range,omitempty"`
	Since        string   `json:"since,omitempty"`
	Until        string   `json:"until,omitempty"`
	Urgencies    []string `json:"urgencies,omitempty"`
	ServiceIDs   []string `json:"service_ids,omitempty"`
	TeamIDs      []string `json:"team_ids,omitempty"`
	UserIDs      []string `json:"user_ids,omitempty"`
	TimeZone     string   `json:"time_zone,omitempty"`
	SortBy       string   `json:"sort_by,omitempty"`
	Includes     []string `json:"include,omitempty"`
	Limit        int      `json:"limit,omitempty"`
	RequestScope string   `json:"request_scope,omitempty"` // "all", "assigned", "teams"
}

// ToParams converts the query to URL parameters
func (q *IncidentQuery) ToParams() map[string]string {
	params := make(map[string]string)
	if q.DateRange != "" {
		params["date_range"] = q.DateRange
	}
	if q.Since != "" {
		params["since"] = q.Since
	}
	if q.Until != "" {
		params["until"] = q.Until
	}
	if q.TimeZone != "" {
		params["time_zone"] = q.TimeZone
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
func (q *IncidentQuery) ToArrayParams() map[string][]string {
	params := make(map[string][]string)
	if len(q.Statuses) > 0 {
		params["statuses[]"] = q.Statuses
	}
	if len(q.Urgencies) > 0 {
		params["urgencies[]"] = q.Urgencies
	}
	if len(q.ServiceIDs) > 0 {
		params["service_ids[]"] = q.ServiceIDs
	}
	if len(q.TeamIDs) > 0 {
		params["team_ids[]"] = q.TeamIDs
	}
	if len(q.UserIDs) > 0 {
		params["user_ids[]"] = q.UserIDs
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

// IncidentCreateRequest represents a request to create an incident
type IncidentCreateRequest struct {
	Incident IncidentCreate `json:"incident"`
}

// IncidentCreate represents the data for creating an incident
type IncidentCreate struct {
	Type             string            `json:"type"`
	Title            string            `json:"title"`
	Service          ServiceReference  `json:"service"`
	Priority         *PriorityReference `json:"priority,omitempty"`
	Urgency          string            `json:"urgency,omitempty"`
	Body             *IncidentBody     `json:"body,omitempty"`
	IncidentKey      string            `json:"incident_key,omitempty"`
	Assignments      []Assignment      `json:"assignments,omitempty"`
	EscalationPolicy *EscalationPolicyReference `json:"escalation_policy,omitempty"`
	ConferenceBridge *ConferenceBridge `json:"conference_bridge,omitempty"`
}

// IncidentManageRequest represents a request to manage incidents
type IncidentManageRequest struct {
	IncidentIDs     []string       `json:"incident_ids"`
	Status          string         `json:"status,omitempty"`
	Urgency         string         `json:"urgency,omitempty"`
	Assignment      *UserReference `json:"assignment,omitempty"`
	EscalationLevel int            `json:"escalation_level,omitempty"`
}

// ToAPIPayload converts the manage request to the API payload format
func (r *IncidentManageRequest) ToAPIPayload() map[string]interface{} {
	incidents := make([]map[string]interface{}, len(r.IncidentIDs))
	for i, id := range r.IncidentIDs {
		incident := map[string]interface{}{
			"type": "incident_reference",
			"id":   id,
		}
		if r.Status != "" {
			incident["status"] = r.Status
		}
		if r.Urgency != "" {
			incident["urgency"] = r.Urgency
		}
		if r.EscalationLevel > 0 {
			incident["escalation_level"] = r.EscalationLevel
		}
		if r.Assignment != nil {
			incident["assignments"] = []map[string]interface{}{
				{
					"at": time.Now().Format(time.RFC3339),
					"assignee": map[string]interface{}{
						"type": "user_reference",
						"id":   r.Assignment.ID,
					},
				},
			}
		}
		incidents[i] = incident
	}
	return map[string]interface{}{"incidents": incidents}
}

// IncidentResponderRequest represents a request to add responders
type IncidentResponderRequest struct {
	RequesterID string                    `json:"requester_id,omitempty"`
	Message     string                    `json:"message,omitempty"`
	Targets     []ResponderRequestTarget  `json:"responder_request_targets"`
}

// ResponderRequestTarget represents a target for responder request
type ResponderRequestTarget struct {
	Type string `json:"responder_request_target_type"`
	ID   string `json:"id"`
}

// IncidentResponderRequestResponse represents the response from adding responders
type IncidentResponderRequestResponse struct {
	ID           string    `json:"id"`
	Incident     Incident  `json:"incident"`
	Requester    User      `json:"requester"`
	RequestedAt  string    `json:"requested_at"`
	Message      string    `json:"message,omitempty"`
}

// IncidentNote represents a note on an incident
type IncidentNote struct {
	ID        string        `json:"id"`
	User      UserReference `json:"user"`
	Content   string        `json:"content"`
	CreatedAt string        `json:"created_at"`
}

// IncidentNoteCreateRequest represents a request to create a note
type IncidentNoteCreateRequest struct {
	Note NoteContent `json:"note"`
}

// NoteContent represents note content
type NoteContent struct {
	Content string `json:"content"`
}

// OutlierIncidentQuery represents query parameters for outlier incidents
type OutlierIncidentQuery struct {
	Since              string `json:"since,omitempty"`
	AdditionalDetails  []string `json:"additional_details,omitempty"`
}

// ToParams converts the query to URL parameters
func (q *OutlierIncidentQuery) ToParams() map[string]string {
	params := make(map[string]string)
	if q.Since != "" {
		params["since"] = q.Since
	}
	if len(q.AdditionalDetails) > 0 {
		params["additional_details[]"] = strings.Join(q.AdditionalDetails, ",")
	}
	return params
}

// OutlierIncidentResponse represents the outlier incident response
type OutlierIncidentResponse struct {
	OutlierIncident *OutlierIncident `json:"outlier_incident,omitempty"`
}

// OutlierIncident represents outlier incident data
type OutlierIncident struct {
	Incident  Incident `json:"incident"`
	IsOutlier bool     `json:"is_outlier"`
}

// PastIncidentsQuery represents query parameters for past incidents
type PastIncidentsQuery struct {
	Limit int  `json:"limit,omitempty"`
	Total bool `json:"total,omitempty"`
}

// ToParams converts the query to URL parameters
func (q *PastIncidentsQuery) ToParams() map[string]string {
	params := make(map[string]string)
	if q.Limit > 0 {
		params["limit"] = fmt.Sprintf("%d", q.Limit)
	}
	if q.Total {
		params["total"] = "true"
	}
	return params
}

// PastIncidentsResponse represents the past incidents response
type PastIncidentsResponse struct {
	PastIncidents []PastIncident `json:"past_incidents"`
	Total         int            `json:"total,omitempty"`
}

// PastIncident represents a past incident with similarity
type PastIncident struct {
	Incident Incident `json:"incident"`
	Score    float64  `json:"score"`
}

// RelatedIncidentsQuery represents query parameters for related incidents
type RelatedIncidentsQuery struct {
	AdditionalDetails []string `json:"additional_details,omitempty"`
}

// ToParams converts the query to URL parameters
func (q *RelatedIncidentsQuery) ToParams() map[string]string {
	params := make(map[string]string)
	if len(q.AdditionalDetails) > 0 {
		params["additional_details[]"] = strings.Join(q.AdditionalDetails, ",")
	}
	return params
}

// RelatedIncidentsResponse represents the related incidents response
type RelatedIncidentsResponse struct {
	RelatedIncidents []RelatedIncident `json:"related_incidents"`
}

// RelatedIncident represents a related incident
type RelatedIncident struct {
	Incident      Incident          `json:"incident"`
	Relationships []RelationshipType `json:"relationships"`
}

// RelationshipType represents a relationship type
type RelationshipType struct {
	Type string `json:"type"`
}

// IncidentResponse is the API response wrapper for a single incident
type IncidentResponse struct {
	Incident Incident `json:"incident"`
}

// IncidentsResponse is the API response wrapper for multiple incidents
type IncidentsResponse struct {
	Incidents []Incident `json:"incidents"`
	Offset    int        `json:"offset"`
	Limit     int        `json:"limit"`
	More      bool       `json:"more"`
	Total     int        `json:"total"`
}

// IncidentNotesResponse is the API response wrapper for incident notes
type IncidentNotesResponse struct {
	Notes []IncidentNote `json:"notes"`
}
