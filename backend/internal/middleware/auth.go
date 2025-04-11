// Package middleware provides HTTP middleware functions
package middleware

import (
	"net/http"

	"go-crypto-bot-clean/backend/internal/auth"

	"go.uber.org/zap"
)

// AuthMiddleware creates a middleware that validates authentication
func AuthMiddleware(clerkAuth *auth.ClerkAuth, logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip authentication if disabled
			if clerkAuth == nil {
				next.ServeHTTP(w, r)
				return
			}

			// Use Clerk middleware
			clerkAuth.Middleware()(next).ServeHTTP(w, r)
		})
	}
}

// RequireAuthMiddleware creates a middleware that requires authentication
func RequireAuthMiddleware(clerkAuth *auth.ClerkAuth, logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip authentication if disabled
			if clerkAuth == nil {
				http.Error(w, "Authentication is required but disabled", http.StatusUnauthorized)
				return
			}

			// Use Clerk middleware
			clerkAuth.Middleware()(next).ServeHTTP(w, r)
		})
	}
}

// RoleMiddleware creates a middleware that validates user roles
func RoleMiddleware(clerkAuth *auth.ClerkAuth, requiredRoles []string, logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip role check if authentication is disabled
			if clerkAuth == nil {
				http.Error(w, "Role check is required but authentication is disabled", http.StatusUnauthorized)
				return
			}

			// Get user from context
			token, err := clerkAuth.GetUserFromContext(r.Context())
			if err != nil {
				logger.Error("Failed to get user from context", zap.Error(err))
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check if user has required roles
			// In a real implementation, this would check the user's roles
			// For now, we'll just allow all authenticated users
			_ = token

			// Continue to the next handler
			next.ServeHTTP(w, r)
		})
	}
}
