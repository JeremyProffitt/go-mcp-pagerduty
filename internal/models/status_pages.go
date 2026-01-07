package models

import "fmt"

// StatusPage represents a PagerDuty status page
type StatusPage struct {
	ID          string `json:"id,omitempty"`
	Type        string `json:"type,omitempty"`
	Self        string `json:"self,omitempty"`
	Name        string `json:"name"`
	URL         string `json:"url,omitempty"`
	Description string `json:"description,omitempty"`
	PublishedAt string `json:"published_at,omitempty"`
	StatusPageType string `json:"status_page_type,omitempty"`
}

// StatusPageQuery represents query parameters for listing status pages
type StatusPageQuery struct {
	Limit int `json:"limit,omitempty"`
}

// ToParams converts the query to URL parameters
func (q *StatusPageQuery) ToParams() map[string]string {
	params := make(map[string]string)
	if q.Limit > 0 {
		params["limit"] = fmt.Sprintf("%d", q.Limit)
	}
	return params
}

// StatusPageReference represents a reference to a status page
type StatusPageReference struct {
	ID   string `json:"id"`
	Type string `json:"type,omitempty"`
}

// StatusPageServiceReference represents a reference to a service on a status page
type StatusPageServiceReference struct {
	ID   string `json:"id"`
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
}

// StatusPageSeverity represents a severity level
type StatusPageSeverity struct {
	ID          string `json:"id"`
	Type        string `json:"type,omitempty"`
	Self        string `json:"self,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// StatusPageSeverityQuery represents query parameters for severities
type StatusPageSeverityQuery struct {
	Limit int `json:"limit,omitempty"`
}

// ToParams converts the query to URL parameters
func (q *StatusPageSeverityQuery) ToParams() map[string]string {
	params := make(map[string]string)
	if q.Limit > 0 {
		params["limit"] = fmt.Sprintf("%d", q.Limit)
	}
	return params
}

// StatusPageSeverityReference represents a reference to a severity
type StatusPageSeverityReference struct {
	ID   string `json:"id"`
	Type string `json:"type,omitempty"`
}

// StatusPageImpact represents an impact level
type StatusPageImpact struct {
	ID          string `json:"id"`
	Type        string `json:"type,omitempty"`
	Self        string `json:"self,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// StatusPageImpactQuery represents query parameters for impacts
type StatusPageImpactQuery struct {
	Limit int `json:"limit,omitempty"`
}

// ToParams converts the query to URL parameters
func (q *StatusPageImpactQuery) ToParams() map[string]string {
	params := make(map[string]string)
	if q.Limit > 0 {
		params["limit"] = fmt.Sprintf("%d", q.Limit)
	}
	return params
}

// StatusPageImpactReference represents a reference to an impact
type StatusPageImpactReference struct {
	ID   string `json:"id"`
	Type string `json:"type,omitempty"`
}

// StatusPageStatus represents a status
type StatusPageStatus struct {
	ID          string `json:"id"`
	Type        string `json:"type,omitempty"`
	Self        string `json:"self,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// StatusPageStatusQuery represents query parameters for statuses
type StatusPageStatusQuery struct {
	Limit int `json:"limit,omitempty"`
}

// ToParams converts the query to URL parameters
func (q *StatusPageStatusQuery) ToParams() map[string]string {
	params := make(map[string]string)
	if q.Limit > 0 {
		params["limit"] = fmt.Sprintf("%d", q.Limit)
	}
	return params
}

// StatusPageStatusReference represents a reference to a status
type StatusPageStatusReference struct {
	ID   string `json:"id"`
	Type string `json:"type,omitempty"`
}

// StatusPagePost represents a post on a status page
type StatusPagePost struct {
	ID           string                      `json:"id,omitempty"`
	Type         string                      `json:"type,omitempty"`
	Self         string                      `json:"self,omitempty"`
	PostType     string                      `json:"post_type"` // incident, maintenance
	Title        string                      `json:"title"`
	StartsAt     string                      `json:"starts_at,omitempty"`
	EndsAt       string                      `json:"ends_at,omitempty"`
	Status       *StatusPageStatusReference  `json:"status,omitempty"`
	Severity     *StatusPageSeverityReference `json:"severity,omitempty"`
	ImpactedServices []StatusPageServiceReference `json:"impacted_services,omitempty"`
	Updates      []StatusPagePostUpdate      `json:"updates,omitempty"`
	StatusPage   *StatusPageReference        `json:"status_page,omitempty"`
	CreatedAt    string                      `json:"created_at,omitempty"`
	UpdatedAt    string                      `json:"updated_at,omitempty"`
}

// StatusPagePostQuery represents query parameters for posts
type StatusPagePostQuery struct {
	PostType string `json:"post_type,omitempty"`
	Limit    int    `json:"limit,omitempty"`
}

// ToParams converts the query to URL parameters
func (q *StatusPagePostQuery) ToParams() map[string]string {
	params := make(map[string]string)
	if q.PostType != "" {
		params["post_type"] = q.PostType
	}
	if q.Limit > 0 {
		params["limit"] = fmt.Sprintf("%d", q.Limit)
	}
	return params
}

// StatusPagePostReference represents a reference to a post
type StatusPagePostReference struct {
	ID   string `json:"id"`
	Type string `json:"type,omitempty"`
}

// StatusPagePostUpdate represents an update to a post
type StatusPagePostUpdate struct {
	ID           string                       `json:"id,omitempty"`
	Type         string                       `json:"type,omitempty"`
	Self         string                       `json:"self,omitempty"`
	Message      string                       `json:"message"`
	Status       *StatusPageStatusReference   `json:"status,omitempty"`
	Severity     *StatusPageSeverityReference `json:"severity,omitempty"`
	ImpactedServices []StatusPagePostUpdateImpact `json:"impacted_services,omitempty"`
	NotifySubscribers bool                     `json:"notify_subscribers,omitempty"`
	ReportedAt   string                       `json:"reported_at,omitempty"`
	CreatedAt    string                       `json:"created_at,omitempty"`
	UpdatedAt    string                       `json:"updated_at,omitempty"`
}

// StatusPagePostUpdateImpact represents impact on a service in an update
type StatusPagePostUpdateImpact struct {
	ID     string                    `json:"id"`
	Type   string                    `json:"type,omitempty"`
	Impact *StatusPageImpactReference `json:"impact,omitempty"`
}

// StatusPagePostUpdateQuery represents query parameters for post updates
type StatusPagePostUpdateQuery struct {
	Limit int `json:"limit,omitempty"`
}

// ToParams converts the query to URL parameters
func (q *StatusPagePostUpdateQuery) ToParams() map[string]string {
	params := make(map[string]string)
	if q.Limit > 0 {
		params["limit"] = fmt.Sprintf("%d", q.Limit)
	}
	return params
}

// StatusPagePostCreateRequest represents a request to create a post
type StatusPagePostCreateRequest struct {
	Post StatusPagePostCreate `json:"post"`
}

// StatusPagePostCreate represents data to create a post
type StatusPagePostCreate struct {
	Type         string                       `json:"type"`
	PostType     string                       `json:"post_type"` // incident, maintenance
	Title        string                       `json:"title"`
	StartsAt     string                       `json:"starts_at,omitempty"`
	EndsAt       string                       `json:"ends_at,omitempty"`
	Status       *StatusPageStatusReference   `json:"status,omitempty"`
	Severity     *StatusPageSeverityReference `json:"severity,omitempty"`
	ImpactedServices []StatusPageServiceReference `json:"impacted_services,omitempty"`
}

// StatusPagePostCreateRequestWrapper wraps the create request
type StatusPagePostCreateRequestWrapper struct {
	Post StatusPagePostCreate `json:"post"`
}

// StatusPagePostUpdateRequest represents a request to add an update
type StatusPagePostUpdateRequest struct {
	PostUpdate StatusPagePostUpdateCreate `json:"post_update"`
}

// StatusPagePostUpdateCreate represents data to create an update
type StatusPagePostUpdateCreate struct {
	Type             string                       `json:"type"`
	Message          string                       `json:"message"`
	Status           *StatusPageStatusReference   `json:"status,omitempty"`
	Severity         *StatusPageSeverityReference `json:"severity,omitempty"`
	ImpactedServices []StatusPagePostUpdateImpact `json:"impacted_services,omitempty"`
	NotifySubscribers bool                        `json:"notify_subscribers,omitempty"`
	ReportedAt       string                       `json:"reported_at,omitempty"`
}

// StatusPagePostUpdateRequestWrapper wraps the update request
type StatusPagePostUpdateRequestWrapper struct {
	PostUpdate StatusPagePostUpdateCreate `json:"post_update"`
}

// StatusPageResponse is the API response wrapper
type StatusPageResponse struct {
	StatusPage StatusPage `json:"status_page"`
}

// StatusPagesResponse is the API response wrapper for multiple status pages
type StatusPagesResponse struct {
	StatusPages []StatusPage `json:"status_pages"`
	Offset      int          `json:"offset"`
	Limit       int          `json:"limit"`
	More        bool         `json:"more"`
	Total       int          `json:"total"`
}

// StatusPageSeveritiesResponse is the API response for severities
type StatusPageSeveritiesResponse struct {
	Severities []StatusPageSeverity `json:"severities"`
}

// StatusPageImpactsResponse is the API response for impacts
type StatusPageImpactsResponse struct {
	Impacts []StatusPageImpact `json:"impacts"`
}

// StatusPageStatusesResponse is the API response for statuses
type StatusPageStatusesResponse struct {
	Statuses []StatusPageStatus `json:"statuses"`
}

// StatusPagePostResponse is the API response for a post
type StatusPagePostResponse struct {
	Post StatusPagePost `json:"post"`
}

// StatusPagePostsResponse is the API response for multiple posts
type StatusPagePostsResponse struct {
	Posts  []StatusPagePost `json:"posts"`
	Offset int              `json:"offset"`
	Limit  int              `json:"limit"`
	More   bool             `json:"more"`
	Total  int              `json:"total"`
}

// StatusPagePostUpdatesResponse is the API response for post updates
type StatusPagePostUpdatesResponse struct {
	PostUpdates []StatusPagePostUpdate `json:"post_updates"`
}

// StatusPagePostUpdateResponse is the API response for a post update
type StatusPagePostUpdateResponse struct {
	PostUpdate StatusPagePostUpdate `json:"post_update"`
}
