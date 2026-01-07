package models

import "fmt"

// Team represents a PagerDuty team
type Team struct {
	ID          string `json:"id"`
	Type        string `json:"type,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Self        string `json:"self,omitempty"`
	HTMLURL     string `json:"html_url,omitempty"`
	Summary     string `json:"summary,omitempty"`
}

// TeamQuery represents query parameters for listing teams
type TeamQuery struct {
	Query string `json:"query,omitempty"`
	Limit int    `json:"limit,omitempty"`
}

// ToParams converts the query to URL parameters
func (q *TeamQuery) ToParams() map[string]string {
	params := make(map[string]string)
	if q.Query != "" {
		params["query"] = q.Query
	}
	if q.Limit > 0 {
		params["limit"] = fmt.Sprintf("%d", q.Limit)
	}
	return params
}

// TeamCreateRequest represents a request to create a team
type TeamCreateRequest struct {
	Team TeamCreate `json:"team"`
}

// TeamCreate represents the data for creating a team
type TeamCreate struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// TeamUpdateRequest represents a request to update a team
type TeamUpdateRequest struct {
	Team TeamUpdate `json:"team"`
}

// TeamUpdate represents the data for updating a team
type TeamUpdate struct {
	Type        string `json:"type"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// TeamMemberAdd represents the data for adding a team member
type TeamMemberAdd struct {
	Role string `json:"role,omitempty"` // manager, responder, observer
}

// TeamMember represents a member of a team
type TeamMember struct {
	User UserReference `json:"user"`
	Role string        `json:"role"`
}

// TeamResponse is the API response wrapper for a single team
type TeamResponse struct {
	Team Team `json:"team"`
}

// TeamsResponse is the API response wrapper for multiple teams
type TeamsResponse struct {
	Teams  []Team `json:"teams"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	More   bool   `json:"more"`
	Total  int    `json:"total"`
}

// TeamMembersResponse is the API response wrapper for team members
type TeamMembersResponse struct {
	Members []TeamMember `json:"members"`
	Offset  int          `json:"offset"`
	Limit   int          `json:"limit"`
	More    bool         `json:"more"`
	Total   int          `json:"total"`
}
