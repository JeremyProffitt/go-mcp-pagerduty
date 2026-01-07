package models

import "fmt"

// AlertGroupingSetting represents an alert grouping setting
type AlertGroupingSetting struct {
	ID          string            `json:"id,omitempty"`
	Type        string            `json:"type,omitempty"`
	Self        string            `json:"self,omitempty"`
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	Services    []ServiceReference `json:"services,omitempty"`
	Config      *AlertGroupingConfig `json:"config,omitempty"`
	CreatedAt   string            `json:"created_at,omitempty"`
	UpdatedAt   string            `json:"updated_at,omitempty"`
}

// AlertGroupingConfig represents the configuration for alert grouping
type AlertGroupingConfig struct {
	Type       string                   `json:"type"` // time, intelligent, content_based
	Timeout    int                      `json:"timeout,omitempty"`
	Aggregate  string                   `json:"aggregate,omitempty"`
	Fields     []string                 `json:"fields,omitempty"`
	TimeWindow int                      `json:"time_window,omitempty"`
}

// TimeGroupingConfig represents time-based grouping config
type TimeGroupingConfig struct {
	Timeout int `json:"timeout"` // minutes
}

// IntelligentGroupingConfig represents intelligent grouping config
type IntelligentGroupingConfig struct {
	TimeWindow int    `json:"time_window,omitempty"` // minutes
	Fields     []string `json:"fields,omitempty"`
}

// ContentBasedConfig represents content-based grouping config
type ContentBasedConfig struct {
	Aggregate  string   `json:"aggregate"` // all, any
	Fields     []string `json:"fields"`
	TimeWindow int      `json:"time_window,omitempty"`
}

// ContentBasedIntelligentConfig represents content-based intelligent config
type ContentBasedIntelligentConfig struct {
	Aggregate  string   `json:"aggregate"` // all, any
	Fields     []string `json:"fields"`
	TimeWindow int      `json:"time_window,omitempty"`
}

// AlertGroupingSettingQuery represents query parameters
type AlertGroupingSettingQuery struct {
	ServiceIDs []string `json:"service_ids,omitempty"`
	Limit      int      `json:"limit,omitempty"`
}

// ToParams converts the query to URL parameters
func (q *AlertGroupingSettingQuery) ToParams() map[string]string {
	params := make(map[string]string)
	if q.Limit > 0 {
		params["limit"] = fmt.Sprintf("%d", q.Limit)
	}
	return params
}

// ToArrayParams converts the query to URL parameters with arrays
func (q *AlertGroupingSettingQuery) ToArrayParams() map[string][]string {
	params := make(map[string][]string)
	if len(q.ServiceIDs) > 0 {
		params["service_ids[]"] = q.ServiceIDs
	}
	for k, v := range q.ToParams() {
		params[k] = []string{v}
	}
	return params
}

// AlertGroupingSettingCreateRequest represents a request to create a setting
type AlertGroupingSettingCreateRequest struct {
	AlertGroupingSetting AlertGroupingSettingCreate `json:"alert_grouping_setting"`
}

// AlertGroupingSettingCreate represents data to create a setting
type AlertGroupingSettingCreate struct {
	Type     string             `json:"type"`
	Name     string             `json:"name"`
	Services []ServiceReference `json:"services"`
	Config   AlertGroupingConfig `json:"config"`
}

// AlertGroupingSettingUpdateRequest represents a request to update a setting
type AlertGroupingSettingUpdateRequest struct {
	AlertGroupingSetting AlertGroupingSettingUpdate `json:"alert_grouping_setting"`
}

// AlertGroupingSettingUpdate represents data to update a setting
type AlertGroupingSettingUpdate struct {
	Type     string              `json:"type,omitempty"`
	Name     string              `json:"name,omitempty"`
	Services []ServiceReference  `json:"services,omitempty"`
	Config   *AlertGroupingConfig `json:"config,omitempty"`
}

// AlertGroupingSettingResponse is the API response wrapper
type AlertGroupingSettingResponse struct {
	AlertGroupingSetting AlertGroupingSetting `json:"alert_grouping_setting"`
}

// AlertGroupingSettingsResponse is the API response wrapper for multiple settings
type AlertGroupingSettingsResponse struct {
	AlertGroupingSettings []AlertGroupingSetting `json:"alert_grouping_settings"`
	Offset                int                    `json:"offset"`
	Limit                 int                    `json:"limit"`
	More                  bool                   `json:"more"`
	Total                 int                    `json:"total"`
}
