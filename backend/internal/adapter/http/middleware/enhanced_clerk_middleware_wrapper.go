package middleware

import (
	"net/http"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/apperror"
	"github.com/rs/zerolog"
)

// CreateBasicMiddlewareAdapter creates an EnhancedClerkMiddleware from a basic ClerkMiddleware
func CreateBasicMiddlewareAdapter(basicMiddleware *ClerkMiddleware, logger *zerolog.Logger) *EnhancedClerkMiddleware {
	return &EnhancedClerkMiddleware{
		logger:      logger,
		authService: nil, // Not used in the adapter
	}
}

// BasicAuthMiddleware is a middleware that uses the basic ClerkMiddleware
type BasicAuthMiddleware struct {
	basicMiddleware *ClerkMiddleware
	logger          *zerolog.Logger
}

// NewBasicAuthMiddleware creates a new BasicAuthMiddleware
func NewBasicAuthMiddleware(basicMiddleware *ClerkMiddleware, logger *zerolog.Logger) *BasicAuthMiddleware {
	return &BasicAuthMiddleware{
		basicMiddleware: basicMiddleware,
		logger:          logger,
	}
}

// Middleware returns a middleware function that validates Clerk authentication
func (m *BasicAuthMiddleware) Middleware() func(http.Handler) http.Handler {
	return m.basicMiddleware.Middleware()
}

// RequireAuthentication is a middleware that requires authentication
func (m *BasicAuthMiddleware) RequireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if user ID is in context
		_, ok := GetUserIDFromContext(r.Context())
		if !ok {
			m.logger.Debug().Msg("Authentication required but user ID not found in context")
			apperror.WriteError(w, apperror.NewUnauthorized("Authentication required", nil))
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// RequireRole is a middleware that requires a specific role
func (m *BasicAuthMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if user has the required role
			roles, ok := r.Context().Value(RoleKey).([]string)
			if !ok {
				m.logger.Debug().Msg("Roles not found in context")
				apperror.WriteError(w, apperror.NewUnauthorized("Authentication required", nil))
				return
			}

			hasRole := false
			for _, r := range roles {
				if r == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				m.logger.Debug().Str("requiredRole", role).Strs("userRoles", roles).Msg("User does not have required role")
				apperror.WriteError(w, apperror.NewForbidden("Insufficient permissions", nil))
				return
			}

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}
