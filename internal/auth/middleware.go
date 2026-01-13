package auth

import (
	"context"
	"net/http"
)

// ContextKey is a type for context keys used by the auth package
type ContextKey string

const (
	// PagerDutyTokenKey is the context key for storing the PagerDuty token
	PagerDutyTokenKey ContextKey = "pagerduty_token"
)

// Middleware creates an HTTP middleware that requires authorization
func Middleware(authorizer Authorizer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip auth for health endpoint
			if r.URL.Path == "/health" {
				next.ServeHTTP(w, r)
				return
			}

			// Get Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error":"Authorization header required"}`, http.StatusUnauthorized)
				return
			}

			// Authorize the token
			authorized, err := authorizer.Authorize(r.Context(), authHeader)
			if err != nil {
				http.Error(w, `{"error":"Authorization failed"}`, http.StatusInternalServerError)
				return
			}

			if !authorized {
				http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
				return
			}

			// Check for X-PagerDuty-Token header and add to context
			ctx := r.Context()
			if pdToken := r.Header.Get("X-PagerDuty-Token"); pdToken != "" {
				ctx = context.WithValue(ctx, PagerDutyTokenKey, pdToken)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetPagerDutyToken retrieves the PagerDuty token from context if present
func GetPagerDutyToken(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(PagerDutyTokenKey).(string)
	return token, ok
}
