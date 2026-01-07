package models

import "fmt"

const (
	DefaultPaginationLimit = 20
	MaxPaginationLimit     = 100
	MaxResults             = 1000
)

// ListResponse is a generic response wrapper for list operations
type ListResponse[T any] struct {
	Response []T `json:"response"`
}

// Summary returns a summary of the list response
func (r *ListResponse[T]) Summary() string {
	count := len(r.Response)
	summary := fmt.Sprintf("Returned %d record(s)", count)
	if count == MaxResults {
		summary += ". WARNING: The number of records equals the response limit. There may be more records not included in this response."
	}
	return summary
}

// QueryParams is an interface for models that can be converted to query parameters
type QueryParams interface {
	ToParams() map[string]string
}

// ArrayQueryParams is an interface for models that need array query parameters
type ArrayQueryParams interface {
	ToArrayParams() map[string][]string
}
