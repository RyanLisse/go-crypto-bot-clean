package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/apperror"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/service"
	"github.com/rs/zerolog"
)

// EnhancedClerkMiddleware handles Clerk authentication with database integration
type EnhancedClerkMiddleware struct {
	logger      *zerolog.Logger
	authService service.AuthServiceInterface
}

// NewEnhancedClerkMiddleware creates a new enhanced Clerk middleware
func NewEnhancedClerkMiddleware(authService service.AuthServiceInterface, logger *zerolog.Logger) *EnhancedClerkMiddleware {
	return &EnhancedClerkMiddleware{
		logger:      logger,
		authService: authService,
	}
}

// Middleware returns a middleware function that validates Clerk authentication
func (m *EnhancedClerkMiddleware) Middleware() func(http.Handler) http.Handler {
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
			ctx = context.WithValue(ctx, RoleKey, roles)
			ctx = context.WithValue(ctx, "user", user)

			// Call the next handler with updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAuthentication is a middleware that requires authentication
func (m *EnhancedClerkMiddleware) RequireAuthentication(next http.Handler) http.Handler {
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
func (m *EnhancedClerkMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
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

// GetUserFromContext extracts the user from the context
func GetUserFromContext(ctx context.Context) (*model.User, bool) {
	val := ctx.Value("user")
	if user, ok := val.(*model.User); ok && user != nil {
		return user, true
	}
	return nil, false
}
