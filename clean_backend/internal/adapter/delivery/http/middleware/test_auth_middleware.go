package middleware

import (
	"context"
	"net/http"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
	"github.com/rs/zerolog"
)

// TestAuthMiddleware is a simple implementation of AuthMiddleware for testing
type TestAuthMiddleware struct {
	logger *zerolog.Logger
}

// NewTestAuthMiddleware creates a new TestAuthMiddleware
func NewTestAuthMiddleware(logger *zerolog.Logger) AuthMiddleware {
	return &TestAuthMiddleware{
		logger: logger,
	}
}

// Middleware returns a middleware function that validates authentication
func (m *TestAuthMiddleware) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// For testing, we'll always set a test user in the context
			user := &model.User{
				ID:    "test-user-id",
				Email: "test@example.com",
				Name:  "Test User",
			}

			// Set user ID and roles in context
			ctx := context.WithValue(r.Context(), UserIDKey, user.ID)
			ctx = context.WithValue(ctx, RolesKey, []string{"user", "admin"})
			ctx = context.WithValue(ctx, UserKey, user)

			// Call the next handler with updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAuthentication is a middleware that requires authentication
func (m *TestAuthMiddleware) RequireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// For testing, we'll always consider the user authenticated
		next.ServeHTTP(w, r)
	})
}

// RequireRole is a middleware that requires a specific role
func (m *TestAuthMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// For testing, we'll always consider the user has the required role
			next.ServeHTTP(w, r)
		})
	}
}
