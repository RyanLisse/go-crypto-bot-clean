package middleware

import (
	"context"
	"net/http"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/apperror"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
	"github.com/rs/zerolog"
)

// TestMiddleware is a middleware for testing authentication
type TestMiddleware struct {
	logger *zerolog.Logger
}

// NewTestMiddleware creates a new TestMiddleware
func NewTestMiddleware(logger *zerolog.Logger) AuthMiddleware {
	return &TestMiddleware{
		logger: logger,
	}
}

// Middleware returns a middleware function that adds a test user to the context
func (m *TestMiddleware) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set a test user ID in the context
			testUser := &model.User{
				ID:    "test_user_id",
				Email: "test@example.com",
				Name:  "Test User",
			}
			
			ctx := context.WithValue(r.Context(), UserIDKey, testUser.ID)
			ctx = context.WithValue(ctx, RolesKey, []string{"user", "admin"})
			ctx = context.WithValue(ctx, UserKey, testUser)

			// Call the next handler with updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAuthentication is a middleware that requires authentication
func (m *TestMiddleware) RequireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if user ID is in context
		userID, ok := GetUserIDFromContext(r.Context())
		if !ok || userID == "" {
			m.logger.Debug().Msg("Authentication required but user ID not found in context")
			appErr := apperror.NewUnauthorized("Authentication required", nil)
			apperror.WriteError(w, appErr)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// RequireRole is a middleware that requires a specific role
func (m *TestMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// In test mode, we'll always consider the user has the required role
			next.ServeHTTP(w, r)
		})
	}
}
