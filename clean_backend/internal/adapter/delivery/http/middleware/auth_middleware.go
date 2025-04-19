package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/apperror"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// AuthMiddlewareImpl implements the AuthMiddleware interface
type AuthMiddlewareImpl struct {
	logger      *zerolog.Logger
	authService port.AuthServiceInterface
}

// NewAuthMiddleware creates a new AuthMiddleware
func NewAuthMiddleware(authService port.AuthServiceInterface, logger *zerolog.Logger) AuthMiddleware {
	return &AuthMiddlewareImpl{
		logger:      logger,
		authService: authService,
	}
}

// Middleware returns a middleware function that validates authentication
func (m *AuthMiddlewareImpl) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check for Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				// Check for Clerk-specific header
				authHeader = r.Header.Get("X-Clerk-Auth-Token")
				if authHeader == "" {
					m.logger.Debug().Msg("No authorization header present")
					next.ServeHTTP(w, r)
					return
				}
			}

			// Extract token
			sessionToken := authHeader
			if strings.HasPrefix(authHeader, "Bearer ") {
				sessionToken = strings.TrimPrefix(authHeader, "Bearer ")
			}

			// Get user from token
			user, err := m.authService.GetUserFromToken(r.Context(), sessionToken)
			if err != nil {
				m.logger.Error().Err(err).Msg("Failed to get user from token")
				apperror.WriteError(w, apperror.NewUnauthorized("Invalid authentication token", err))
				return
			}

			// Get user roles
			roles, err := m.authService.GetUserRoles(r.Context(), user.ID)
			if err != nil {
				m.logger.Error().Err(err).Msg("Failed to get user roles")
				roles = []string{"user"} // Default role
			}

			// Set user ID and roles in context
			ctx := context.WithValue(r.Context(), UserIDKey, user.ID)
			ctx = context.WithValue(ctx, RolesKey, roles)
			ctx = context.WithValue(ctx, UserKey, user)

			// Call the next handler with updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAuthentication is a middleware that requires authentication
func (m *AuthMiddlewareImpl) RequireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if user ID is in context
		userID, ok := GetUserIDFromContext(r.Context())
		if !ok || userID == "" {
			m.logger.Debug().Msg("Authentication required but user ID not found in context")
			apperror.WriteError(w, apperror.NewUnauthorized("Authentication required", nil))
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// RequireRole is a middleware that requires a specific role
func (m *AuthMiddlewareImpl) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if user has the required role
			roles, ok := GetRolesFromContext(r.Context())
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
