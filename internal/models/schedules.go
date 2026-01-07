package models

import "fmt"

// Schedule represents a PagerDuty schedule
type Schedule struct {
	ID                   string              `json:"id,omitempty"`
	Type                 string              `json:"type,omitempty"`
	Summary              string              `json:"summary,omitempty"`
	Self                 string              `json:"self,omitempty"`
	HTMLURL              string              `json:"html_url,omitempty"`
	Name                 string              `json:"name"`
	Description          string              `json:"description,omitempty"`
	TimeZone             string              `json:"time_zone"`
	EscalationPolicies   []EscalationPolicyReference `json:"escalation_policies,omitempty"`
	Users                []UserReference     `json:"users,omitempty"`
	Teams                []TeamReference     `json:"teams,omitempty"`
	ScheduleLayers       []ScheduleLayer     `json:"schedule_layers,omitempty"`
	OverridesSubschedule *Subschedule        `json:"overrides_subschedule,omitempty"`
	FinalSchedule        *Subschedule        `json:"final_schedule,omitempty"`
}

// ScheduleLayer represents a layer in a schedule
type ScheduleLayer struct {
	ID                         string              `json:"id,omitempty"`
	Type                       string              `json:"type,omitempty"`
	Name                       string              `json:"name,omitempty"`
	Start                      string              `json:"start"`
	End                        string              `json:"end,omitempty"`
	RotationVirtualStart       string              `json:"rotation_virtual_start"`
	RotationTurnLengthSeconds  int                 `json:"rotation_turn_length_seconds"`
	Users                      []ScheduleLayerUser `json:"users"`
	Restrictions               []ScheduleLayerRestriction `json:"restrictions,omitempty"`
	RenderedScheduleEntries    []RenderedScheduleEntry    `json:"rendered_schedule_entries,omitempty"`
	RenderedCoveragePercentage float64             `json:"rendered_coverage_percentage,omitempty"`
}

// ScheduleLayerUser represents a user in a schedule layer
type ScheduleLayerUser struct {
	User UserReference `json:"user"`
}

// ScheduleLayerRestriction represents a restriction on a schedule layer
type ScheduleLayerRestriction struct {
	Type            string `json:"type"`
	StartTimeOfDay  string `json:"start_time_of_day"`
	DurationSeconds int    `json:"duration_seconds"`
	StartDayOfWeek  int    `json:"start_day_of_week,omitempty"`
}

// RenderedScheduleEntry represents a rendered schedule entry
type RenderedScheduleEntry struct {
	Start string        `json:"start"`
	End   string        `json:"end"`
	User  UserReference `json:"user"`
}

// Subschedule represents a subschedule
type Subschedule struct {
	Name                    string                  `json:"name,omitempty"`
	RenderedScheduleEntries []RenderedScheduleEntry `json:"rendered_schedule_entries,omitempty"`
	RenderedCoveragePercentage float64              `json:"rendered_coverage_percentage,omitempty"`
}

// ScheduleQuery represents query parameters for listing schedules
type ScheduleQuery struct {
	Query string `json:"query,omitempty"`
	Limit int    `json:"limit,omitempty"`
}

// ToParams converts the query to URL parameters
func (q *ScheduleQuery) ToParams() map[string]string {
	params := make(map[string]string)
	if q.Query != "" {
		params["query"] = q.Query
	}
	if q.Limit > 0 {
		params["limit"] = fmt.Sprintf("%d", q.Limit)
	}
	return params
}

// ScheduleCreateRequest represents a request to create a schedule
type ScheduleCreateRequest struct {
	Schedule ScheduleCreateData `json:"schedule"`
}

// ScheduleCreateData represents the data for creating a schedule
type ScheduleCreateData struct {
	Type           string                `json:"type"`
	Name           string                `json:"name"`
	Description    string                `json:"description,omitempty"`
	TimeZone       string                `json:"time_zone"`
	ScheduleLayers []ScheduleLayerCreate `json:"schedule_layers"`
}

// ScheduleLayerCreate represents a layer for creating a schedule
type ScheduleLayerCreate struct {
	Name                       string              `json:"name,omitempty"`
	Start                      string              `json:"start"`
	End                        string              `json:"end,omitempty"`
	RotationVirtualStart       string              `json:"rotation_virtual_start"`
	RotationTurnLengthSeconds  int                 `json:"rotation_turn_length_seconds"`
	Users                      []ScheduleLayerUser `json:"users"`
	Restrictions               []ScheduleLayerRestriction `json:"restrictions,omitempty"`
}

// ScheduleUpdateRequest represents a request to update a schedule
type ScheduleUpdateRequest struct {
	Schedule ScheduleUpdateData `json:"schedule"`
}

// ScheduleUpdateData represents the data for updating a schedule
type ScheduleUpdateData struct {
	Type           string                `json:"type"`
	Name           string                `json:"name,omitempty"`
	Description    string                `json:"description,omitempty"`
	TimeZone       string                `json:"time_zone,omitempty"`
	ScheduleLayers []ScheduleLayerCreate `json:"schedule_layers,omitempty"`
}

// ScheduleOverrideCreate represents data for creating a schedule override
type ScheduleOverrideCreate struct {
	Override OverrideData `json:"override"`
}

// OverrideData represents override data
type OverrideData struct {
	Start string        `json:"start"`
	End   string        `json:"end"`
	User  UserReference `json:"user"`
}

// ScheduleOverride represents a schedule override
type ScheduleOverride struct {
	ID    string        `json:"id"`
	Start string        `json:"start"`
	End   string        `json:"end"`
	User  UserReference `json:"user"`
}

// ScheduleResponse is the API response wrapper for a single schedule
type ScheduleResponse struct {
	Schedule Schedule `json:"schedule"`
}

// SchedulesResponse is the API response wrapper for multiple schedules
type SchedulesResponse struct {
	Schedules []Schedule `json:"schedules"`
	Offset    int        `json:"offset"`
	Limit     int        `json:"limit"`
	More      bool       `json:"more"`
	Total     int        `json:"total"`
}

// ScheduleUsersResponse is the API response wrapper for schedule users
type ScheduleUsersResponse struct {
	Users []User `json:"users"`
}

// ScheduleOverrideResponse is the API response wrapper for a schedule override
type ScheduleOverrideResponse struct {
	Override ScheduleOverride `json:"override"`
}
