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

// UserIDKey is the context key for user ID
type UserIDKey struct{}

// RolesKey is the context key for user roles
type RolesKey struct{}

// UserKey is the context key for user
type UserKey struct{}

// AuthMiddleware defines the interface for authentication middleware
type AuthMiddleware interface {
	// Middleware returns a middleware function that validates authentication
	Middleware() func(http.Handler) http.Handler
	
	// RequireAuthentication is a middleware that requires authentication
	RequireAuthentication(next http.Handler) http.Handler
	
	// RequireRole is a middleware that requires a specific role
	RequireRole(role string) func(http.Handler) http.Handler
}

// AuthMiddlewareImpl is the primary authentication middleware
type AuthMiddlewareImpl struct {
	logger      *zerolog.Logger
	authService service.AuthServiceInterface
}

// NewAuthMiddleware creates a new AuthMiddlewareImpl
func NewAuthMiddleware(authService service.AuthServiceInterface, logger *zerolog.Logger) AuthMiddleware {
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
			ctx := context.WithValue(r.Context(), UserIDKey{}, user.ID)
			ctx = context.WithValue(ctx, RolesKey{}, roles)
			ctx = context.WithValue(ctx, UserKey{}, user)

			// Call the next handler with updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAuthentication is a middleware that requires authentication
func (m *AuthMiddlewareImpl) RequireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if user ID is in context
		userID, ok := r.Context().Value(UserIDKey{}).(string)
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
			roles, ok := r.Context().Value(RolesKey{}).([]string)
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

// GetUserFromContext gets the user from the context
func GetUserFromContext(ctx context.Context) (*model.User, bool) {
	user, ok := ctx.Value(UserKey{}).(*model.User)
	return user, ok
}

// GetUserIDFromContext gets the user ID from the context
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey{}).(string)
	return userID, ok
}

// GetRolesFromContext gets the user roles from the context
func GetRolesFromContext(ctx context.Context) ([]string, bool) {
	roles, ok := ctx.Value(RolesKey{}).([]string)
	return roles, ok
}

// TestAuthMiddleware is a middleware for testing authentication
type TestAuthMiddleware struct {
	logger *zerolog.Logger
}

// NewTestAuthMiddleware creates a new TestAuthMiddleware
func NewTestAuthMiddleware(logger *zerolog.Logger) AuthMiddleware {
	return &TestAuthMiddleware{
		logger: logger,
	}
}

// Middleware returns a middleware function that adds a test user to the context
func (m *TestAuthMiddleware) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set a test user ID in the context
			ctx := context.WithValue(r.Context(), UserIDKey{}, "test_user_id")
			ctx = context.WithValue(ctx, RolesKey{}, []string{"user", "admin"})
			ctx = context.WithValue(ctx, UserKey{}, &model.User{
				ID:    "test_user_id",
				Email: "test@example.com",
				Name:  "Test User",
			})

			// Call the next handler with updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAuthentication is a middleware that requires authentication
func (m *TestAuthMiddleware) RequireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if user ID is in context
		userID, ok := r.Context().Value(UserIDKey{}).(string)
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
func (m *TestAuthMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// In test mode, we'll just assume the user has the required role
			next.ServeHTTP(w, r)
		})
	}
}

// DisabledAuthMiddleware is a middleware that disables authentication
type DisabledAuthMiddleware struct {
	logger *zerolog.Logger
}

// NewDisabledAuthMiddleware creates a new DisabledAuthMiddleware
func NewDisabledAuthMiddleware(logger *zerolog.Logger) AuthMiddleware {
	return &DisabledAuthMiddleware{
		logger: logger,
	}
}

// Middleware returns a middleware function that does nothing
func (m *DisabledAuthMiddleware) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return next
	}
}

// RequireAuthentication is a middleware that does nothing
func (m *DisabledAuthMiddleware) RequireAuthentication(next http.Handler) http.Handler {
	return next
}

// RequireRole is a middleware that does nothing
func (m *DisabledAuthMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return next
	}
}
