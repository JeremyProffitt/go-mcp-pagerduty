package models

import "fmt"

// ChangeEvent represents a PagerDuty change event
type ChangeEvent struct {
	ID             string             `json:"id,omitempty"`
	Type           string             `json:"type,omitempty"`
	Self           string             `json:"self,omitempty"`
	Summary        string             `json:"summary,omitempty"`
	Source         string             `json:"source,omitempty"`
	RoutingKey     string             `json:"routing_key,omitempty"`
	Timestamp      string             `json:"timestamp,omitempty"`
	Integration    *IntegrationReference `json:"integration,omitempty"`
	Services       []ServiceReference `json:"services,omitempty"`
	Links          []ChangeEventLink  `json:"links,omitempty"`
	Images         []ChangeEventImage `json:"images,omitempty"`
	CustomDetails  map[string]interface{} `json:"custom_details,omitempty"`
}

// ChangeEventLink represents a link in a change event
type ChangeEventLink struct {
	Href string `json:"href"`
	Text string `json:"text,omitempty"`
}

// ChangeEventImage represents an image in a change event
type ChangeEventImage struct {
	Src  string `json:"src"`
	Href string `json:"href,omitempty"`
	Alt  string `json:"alt,omitempty"`
}

// ChangeEventQuery represents query parameters for listing change events
type ChangeEventQuery struct {
	Since       string   `json:"since,omitempty"`
	Until       string   `json:"until,omitempty"`
	TeamIDs     []string `json:"team_ids,omitempty"`
	ServiceIDs  []string `json:"service_ids,omitempty"`
	IntegrationIDs []string `json:"integration_ids,omitempty"`
	Limit       int      `json:"limit,omitempty"`
}

// ToParams converts the query to URL parameters
func (q *ChangeEventQuery) ToParams() map[string]string {
	params := make(map[string]string)
	if q.Since != "" {
		params["since"] = q.Since
	}
	if q.Until != "" {
		params["until"] = q.Until
	}
	if q.Limit > 0 {
		params["limit"] = fmt.Sprintf("%d", q.Limit)
	}
	return params
}

// ToArrayParams converts the query to URL parameters with arrays
func (q *ChangeEventQuery) ToArrayParams() map[string][]string {
	params := make(map[string][]string)
	if len(q.TeamIDs) > 0 {
		params["team_ids[]"] = q.TeamIDs
	}
	if len(q.ServiceIDs) > 0 {
		params["service_ids[]"] = q.ServiceIDs
	}
	if len(q.IntegrationIDs) > 0 {
		params["integration_ids[]"] = q.IntegrationIDs
	}
	for k, v := range q.ToParams() {
		params[k] = []string{v}
	}
	return params
}

// ChangeEventResponse is the API response wrapper
type ChangeEventResponse struct {
	ChangeEvent ChangeEvent `json:"change_event"`
}

// ChangeEventsResponse is the API response wrapper for multiple change events
type ChangeEventsResponse struct {
	ChangeEvents []ChangeEvent `json:"change_events"`
	Offset       int           `json:"offset"`
	Limit        int           `json:"limit"`
	More         bool          `json:"more"`
	Total        int           `json:"total"`
}
