package auth

import (
	"context"
)

// Authorizer defines the interface for authorizing requests
type Authorizer interface {
	Authorize(ctx context.Context, token string) (bool, error)
}

// MockAuthorizer is a mock implementation that always authorizes
type MockAuthorizer struct{}

// Authorize always returns true for the mock authorizer
func (m *MockAuthorizer) Authorize(ctx context.Context, token string) (bool, error) {
	return true, nil
}
