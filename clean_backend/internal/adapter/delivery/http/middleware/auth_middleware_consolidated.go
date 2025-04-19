package middleware

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	// Updated imports to clean_backend packages

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/apperror"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"

	// Import the domain port package, specifically the auth_service
	domainport "github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port"
	"github.com/rs/zerolog"
	// Note: Clerk SDK import needed for real implementation (Task 4)
)

// Remove duplicate ContextKey definition
/*
type ContextKey string
*/

// Remove duplicate context keys definitions
/*
const (
	// UserIDKey is the context key for user ID (Note: Using the one from middlewares.go)
	// UserIDKey ContextKey = "user_id"
	// RolesKey is the context key for user roles
	RolesKey ContextKey = "roles"
	// UserKey is the context key for user
	UserKey ContextKey = "user"
)
*/

// Remove duplicate AuthMiddleware interface definition and use the one from middlewares.go
/*
type AuthMiddleware interface {
	Middleware() func(http.Handler) http.Handler
	RequireAuthentication(next http.Handler) http.Handler
	RequireRole(role string) func(http.Handler) http.Handler
}
*/

// ClerkMiddleware is the primary authentication middleware using Clerk
type ClerkMiddleware struct {
	logger *zerolog.Logger
	// Updated to depend on domain port for auth service, use domainport alias
	authService domainport.AuthServiceInterface
	config      *config.Config
	jwkCache    map[string]interface{} // Cache for JWKs
	jwkCacheMu  sync.RWMutex           // Mutex for JWK cache
	jwkCacheExp time.Time              // Expiration time for JWK cache
}

// NewClerkMiddleware creates a new ClerkMiddleware
// Updated authService type to domainport.AuthServiceInterface
func NewClerkMiddleware(authService domainport.AuthServiceInterface, config *config.Config, logger *zerolog.Logger) AuthMiddleware {
	return &ClerkMiddleware{
		logger:      logger,
		authService: authService,
		config:      config,
		jwkCache:    make(map[string]interface{}),
		jwkCacheExp: time.Now(),
	}
}

// Middleware returns a middleware function that validates authentication
func (m *ClerkMiddleware) Middleware() func(http.Handler) http.Handler {
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

			// Check if token is empty
			if sessionToken == "" {
				m.logger.Debug().Msg("Empty token in authorization header")
				next.ServeHTTP(w, r)
				return
			}

			// Verify token and get claims (verifyToken needs review in Task 4)
			claims, err := m.verifyToken(r.Context(), sessionToken)
			if err != nil {
				m.logger.Error().Err(err).Msg("Failed to verify token")
				// Use clean_backend apperror and response writer
				appErr := apperror.NewUnauthorized("Invalid authentication token", err)
				apperror.WriteError(w, appErr)
				return
			}

			// Get user from claims
			userID := claims.Subject
			// Use authService from domain port
			user, err := m.authService.GetUserByID(r.Context(), userID)
			if err != nil {
				m.logger.Error().Err(err).Str("user_id", userID).Msg("Failed to get user from ID")
				// Use clean_backend apperror and response writer
				appErr := apperror.NewUnauthorized("Invalid user ID", err)
				apperror.WriteError(w, appErr)
				return
			}

			// Get user roles (use authService from domain port)
			roles, err := m.authService.GetUserRoles(r.Context(), user.ID)
			if err != nil {
				m.logger.Error().Err(err).Msg("Failed to get user roles, defaulting to user role")
				roles = []string{"user"} // Default role
			}

			// Set user ID, roles, and claims in context (using UserIDKey from middlewares.go)
			ctx := context.WithValue(r.Context(), UserIDKey, user.ID)
			ctx = context.WithValue(ctx, RolesKey, roles)
			ctx = context.WithValue(ctx, UserKey, user)
			ctx = context.WithValue(ctx, "clerk_claims", claims)

			// Call the next handler with updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// SessionClaims represents the claims in a Clerk session token
type SessionClaims struct {
	Subject         string   `json:"sub"`
	SessionID       string   `json:"sid"`
	UserID          string   `json:"user_id"` // Changed from sub to user_id to avoid duplicate
	Email           string   `json:"email,omitempty"`
	Name            string   `json:"name,omitempty"`
	AuthorizedParty string   `json:"azp,omitempty"`
	Expiry          int64    `json:"exp"`
	IssuedAt        int64    `json:"iat"`
	NotBefore       int64    `json:"nbf"`
	Roles           []string `json:"roles,omitempty"`
}

// verifyToken verifies a Clerk session token and returns the claims
// TODO: Implement real Clerk JWKS verification here (Task 4)
func (m *ClerkMiddleware) verifyToken(ctx context.Context, token string) (*SessionClaims, error) {
	// In a real implementation, this would verify the token with Clerk's JWKS
	// For now, we'll use the authService to get the user from the token (placeholder)
	user, err := m.authService.GetUserFromToken(ctx, token) // This call might need adjustment based on port
	if err != nil {
		return nil, apperror.NewUnauthorized("Token verification failed", err) // Use clean_backend apperror
	}

	// Create claims from user (example structure)
	claims := &SessionClaims{
		Subject:   user.ID,
		SessionID: "session_" + user.ID, // Example
		UserID:    user.ID,
		Email:     user.Email,                            // Assuming User model has Email
		Name:      user.Name,                             // Assuming User model has Name
		Expiry:    time.Now().Add(24 * time.Hour).Unix(), // Example expiry
		IssuedAt:  time.Now().Unix(),
		NotBefore: time.Now().Unix(),
	}

	return claims, nil
}

// RequireAuthentication is a middleware that requires authentication
func (m *ClerkMiddleware) RequireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if user ID is in context (using GetUserIDFromContext from middlewares.go)
		userID, ok := GetUserIDFromContext(r.Context())
		if !ok || userID == "" {
			m.logger.Debug().Msg("Authentication required but user ID not found in context")
			// Use clean_backend apperror and response writer
			appErr := apperror.NewUnauthorized("Authentication required", nil)
			apperror.WriteError(w, appErr)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// RequireRole is a middleware that requires a specific role
func (m *ClerkMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if user has the required role (using GetRolesFromContext from middlewares.go)
			roles, ok := GetRolesFromContext(r.Context())
			if !ok {
				m.logger.Debug().Msg("Roles not found in context")
				// Use clean_backend apperror and response writer
				appErr := apperror.NewUnauthorized("Authentication required (roles not found)", nil) // Use Unauthorized or Forbidden?
				apperror.WriteError(w, appErr)
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
				// Use clean_backend apperror and response writer
				appErr := apperror.NewForbidden("Insufficient permissions", nil)
				apperror.WriteError(w, appErr)
				return
			}

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}

// TestMiddleware is a middleware for testing authentication
// Keep for now, will likely be removed when factory is integrated
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
			// Set a test user ID in the context (using UserIDKey from middlewares.go)
			ctx := context.WithValue(r.Context(), UserIDKey, "test_user_id")
			ctx = context.WithValue(ctx, RolesKey, []string{"user", "admin"})
			// Use clean_backend model.User
			ctx = context.WithValue(ctx, UserKey, &model.User{
				ID:    "test_user_id",
				Email: "test@example.com",
				Name:  "Test User",
				// Add other fields if necessary based on clean_backend model.User
			})

			// Call the next handler with updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAuthentication is a middleware that requires authentication
func (m *TestMiddleware) RequireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if user ID is in context (using GetUserIDFromContext from middlewares.go)
		userID, ok := GetUserIDFromContext(r.Context())
		if !ok || userID == "" {
			m.logger.Debug().Msg("Authentication required but user ID not found in context")
			// Use clean_backend apperror and response writer
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
			// Check if user has the required role (using GetRolesFromContext from middlewares.go)
			roles, ok := GetRolesFromContext(r.Context())
			if !ok {
				m.logger.Debug().Msg("Roles not found in context")
				// Use clean_backend apperror and response writer
				appErr := apperror.NewUnauthorized("Authentication required (roles not found)", nil) // Use Unauthorized or Forbidden?
				apperror.WriteError(w, appErr)
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
				// Use clean_backend apperror and response writer
				appErr := apperror.NewForbidden("Insufficient permissions", nil)
				apperror.WriteError(w, appErr)
				return
			}

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}

// DisabledMiddleware is a middleware that bypasses authentication for all requests
// Keep for now, will likely be removed when factory is integrated
type DisabledMiddleware struct {
	logger *zerolog.Logger
}

// NewDisabledMiddleware creates a new DisabledMiddleware
func NewDisabledMiddleware(logger *zerolog.Logger) AuthMiddleware {
	return &DisabledMiddleware{
		logger: logger,
	}
}

// Middleware returns a middleware function that does nothing
func (m *DisabledMiddleware) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Authentication is disabled, just pass through
			m.logger.Debug().Msg("Authentication disabled, bypassing checks")
			next.ServeHTTP(w, r)
		})
	}
}

// RequireAuthentication is a middleware that requires authentication (disabled)
func (m *DisabledMiddleware) RequireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Authentication is disabled, just pass through
		next.ServeHTTP(w, r)
	})
}

// RequireRole is a middleware that requires a specific role (disabled)
func (m *DisabledMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Authentication is disabled, just pass through
			next.ServeHTTP(w, r)
		})
	}
}

// Remove duplicate helper functions - these should be in middlewares.go
/*
// GetUserFromContext retrieves the *model.User from the context
// Use clean_backend model.User
func GetUserFromContext(ctx context.Context) (*model.User, bool) {
	user, ok := ctx.Value(UserKey).(*model.User)
	return user, ok
}

// GetUserIDFromContext retrieves the user ID string from the context
// Use clean_backend UserIDKey from middlewares.go
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}

// GetRolesFromContext retrieves the user roles slice from the context
func GetRolesFromContext(ctx context.Context) ([]string, bool) {
	roles, ok := ctx.Value(RolesKey).([]string)
	return roles, ok
}
*/
