package models

import "fmt"

// User represents a PagerDuty user
type User struct {
	ID             string          `json:"id"`
	Type           string          `json:"type,omitempty"`
	Name           string          `json:"name"`
	Email          string          `json:"email"`
	TimeZone       string          `json:"time_zone,omitempty"`
	Color          string          `json:"color,omitempty"`
	Role           string          `json:"role,omitempty"`
	Description    string          `json:"description,omitempty"`
	InvitationSent bool            `json:"invitation_sent,omitempty"`
	JobTitle       string          `json:"job_title,omitempty"`
	Teams          []TeamReference `json:"teams,omitempty"`
	Self           string          `json:"self,omitempty"`
	HTMLURL        string          `json:"html_url,omitempty"`
	AvatarURL      string          `json:"avatar_url,omitempty"`
}

// UserQuery represents query parameters for listing users
type UserQuery struct {
	Query   string `json:"query,omitempty"`
	TeamIDs []string `json:"team_ids,omitempty"`
	Limit   int    `json:"limit,omitempty"`
}

// ToParams converts the query to URL parameters
func (q *UserQuery) ToParams() map[string]string {
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
func (q *UserQuery) ToArrayParams() map[string][]string {
	params := make(map[string][]string)
	if q.Query != "" {
		params["query"] = []string{q.Query}
	}
	if len(q.TeamIDs) > 0 {
		params["team_ids[]"] = q.TeamIDs
	}
	if q.Limit > 0 {
		params["limit"] = []string{fmt.Sprintf("%d", q.Limit)}
	}
	return params
}

// MCPContext holds the context for MCP operations
type MCPContext struct {
	User *User `json:"user,omitempty"`
}

// UserResponse is the API response wrapper for a single user
type UserResponse struct {
	User User `json:"user"`
}

// UsersResponse is the API response wrapper for multiple users
type UsersResponse struct {
	Users  []User `json:"users"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	More   bool   `json:"more"`
	Total  int    `json:"total"`
}
