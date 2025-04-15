package http

import (
	"context"
	"net/http"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/http/middleware"
	"github.com/rs/zerolog"
)

// SimpleAuthMiddleware is a simple implementation of the AuthMiddleware interface for testing
type SimpleAuthMiddleware struct {
	logger *zerolog.Logger
}

// NewSimpleAuthMiddleware creates a new SimpleAuthMiddleware
func NewSimpleAuthMiddleware(logger *zerolog.Logger) middleware.AuthMiddleware {
	return &SimpleAuthMiddleware{
		logger: logger,
	}
}

// RequireAuthentication is a middleware that requires authentication
func (m *SimpleAuthMiddleware) RequireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// For testing, we'll just set a dummy user ID in the context
		ctx := context.WithValue(r.Context(), middleware.UserIDKey, "test_user_id")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
