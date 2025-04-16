package middleware

import (
	"context"
	"net/http"

	"github.com/rs/zerolog"
)

// MockMiddleware is a simple middleware that does nothing
type MockMiddleware struct {
	logger *zerolog.Logger
}

// NewMockMiddleware creates a new mock middleware
func NewMockMiddleware(logger *zerolog.Logger) *MockMiddleware {
	return &MockMiddleware{
		logger: logger,
	}
}

// Middleware returns a middleware function
func (m *MockMiddleware) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Add user information to the request context
			ctx := context.WithValue(r.Context(), "userID", "user123")
			ctx = context.WithValue(ctx, "roles", []string{"user"})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAuthentication is a middleware that requires authentication
func (m *MockMiddleware) RequireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Always allow access in mock middleware
		next.ServeHTTP(w, r)
	})
}

// RequireRole is a middleware that requires a specific role
func (m *MockMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Always allow access in mock middleware
			next.ServeHTTP(w, r)
		})
	}
}

// CreateBasicMiddlewareAdapter creates a basic middleware adapter
func CreateBasicMiddlewareAdapter(logger *zerolog.Logger) *MockMiddleware {
	return NewMockMiddleware(logger)
}
