package models

import "fmt"

// Oncall represents an on-call entry
type Oncall struct {
	EscalationPolicy EscalationPolicyReference `json:"escalation_policy"`
	EscalationLevel  int                       `json:"escalation_level"`
	Schedule         *ScheduleReference        `json:"schedule,omitempty"`
	User             UserReference             `json:"user"`
	Start            string                    `json:"start,omitempty"`
	End              string                    `json:"end,omitempty"`
}

// OncallQuery represents query parameters for listing on-calls
type OncallQuery struct {
	TimeZone            string   `json:"time_zone,omitempty"`
	Since               string   `json:"since,omitempty"`
	Until               string   `json:"until,omitempty"`
	Earliest            bool     `json:"earliest,omitempty"`
	ScheduleIDs         []string `json:"schedule_ids,omitempty"`
	UserIDs             []string `json:"user_ids,omitempty"`
	EscalationPolicyIDs []string `json:"escalation_policy_ids,omitempty"`
	Includes            []string `json:"include,omitempty"`
	Limit               int      `json:"limit,omitempty"`
}

// ToParams converts the query to URL parameters
func (q *OncallQuery) ToParams() map[string]string {
	params := make(map[string]string)
	if q.TimeZone != "" {
		params["time_zone"] = q.TimeZone
	}
	if q.Since != "" {
		params["since"] = q.Since
	}
	if q.Until != "" {
		params["until"] = q.Until
	}
	if q.Earliest {
		params["earliest"] = "true"
	}
	if q.Limit > 0 {
		params["limit"] = fmt.Sprintf("%d", q.Limit)
	}
	return params
}

// ToArrayParams converts the query to URL parameters with arrays
func (q *OncallQuery) ToArrayParams() map[string][]string {
	params := make(map[string][]string)
	if len(q.ScheduleIDs) > 0 {
		params["schedule_ids[]"] = q.ScheduleIDs
	}
	if len(q.UserIDs) > 0 {
		params["user_ids[]"] = q.UserIDs
	}
	if len(q.EscalationPolicyIDs) > 0 {
		params["escalation_policy_ids[]"] = q.EscalationPolicyIDs
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

// OncallsResponse is the API response wrapper for on-calls
type OncallsResponse struct {
	Oncalls []Oncall `json:"oncalls"`
	Offset  int      `json:"offset"`
	Limit   int      `json:"limit"`
	More    bool     `json:"more"`
	Total   int      `json:"total"`
}
