package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/apperror"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
	domainport "github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port"

	"github.com/rs/zerolog"
)

// ClerkMiddleware is the primary authentication middleware using Clerk
type ClerkMiddleware struct {
	logger      *zerolog.Logger
	authService domainport.AuthServiceInterface
	config      *config.Config
	jwkCache    map[string]interface{} // Cache for JWKs
	jwkCacheMu  sync.RWMutex           // Mutex for JWK cache
	jwkCacheExp time.Time              // Expiration time for JWK cache
}

// NewClerkMiddleware creates a new ClerkMiddleware
func NewClerkMiddleware(authService domainport.AuthServiceInterface, config *config.Config, logger *zerolog.Logger) AuthMiddleware {
	// Check if Clerk secret key is set
	if config.Auth.ClerkSecretKey == "" {
		logger.Warn().Msg("CLERK_SECRET_KEY not set, authentication will not work properly")
	}

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

			// Verify token and get claims
			claims, err := m.verifyToken(r.Context(), sessionToken)
			if err != nil {
				m.logger.Error().Err(err).Msg("Failed to verify token")
				appErr := apperror.NewUnauthorized("Invalid authentication token", err)
				apperror.WriteError(w, appErr)
				return
			}

			// Get user from claims
			userID := claims.Subject
			user, err := m.getUserFromClaims(r.Context(), claims)
			if err != nil {
				m.logger.Error().Err(err).Str("user_id", userID).Msg("Failed to get user from claims")
				appErr := apperror.NewUnauthorized("Invalid user ID", err)
				apperror.WriteError(w, appErr)
				return
			}

			// Get user roles
			roles, err := m.getRolesFromClaims(r.Context(), claims)
			if err != nil {
				m.logger.Error().Err(err).Msg("Failed to get user roles, defaulting to user role")
				roles = []string{"user"} // Default role
			}

			// Set user ID, roles, and claims in context
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
	UserID          string   `json:"user_id"` // Changed from sub to avoid duplicate
	Email           string   `json:"email,omitempty"`
	Name            string   `json:"name,omitempty"`
	AuthorizedParty string   `json:"azp,omitempty"`
	Expiry          int64    `json:"exp"`
	IssuedAt        int64    `json:"iat"`
	NotBefore       int64    `json:"nbf"`
	Roles           []string `json:"roles,omitempty"`
}

// verifyToken verifies a Clerk session token and returns the claims
func (m *ClerkMiddleware) verifyToken(ctx context.Context, token string) (*SessionClaims, error) {
	// Use the auth service to validate the token
	valid, err := m.authService.ValidateToken(ctx, token)
	if err != nil || !valid {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}

	// Get user from token
	user, err := m.authService.GetUserFromToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get user from token: %w", err)
	}

	// Create session claims
	sessionClaims := &SessionClaims{
		Subject:   user.ID,
		SessionID: "session_" + user.ID, // Example
		UserID:    user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Expiry:    time.Now().Add(24 * time.Hour).Unix(), // Example
		IssuedAt:  time.Now().Unix(),
		NotBefore: time.Now().Unix(),
	}

	// Get roles
	roles, err := m.authService.GetUserRoles(ctx, user.ID)
	if err == nil && len(roles) > 0 {
		sessionClaims.Roles = roles
	}

	return sessionClaims, nil
}

// getUserFromClaims creates a user model from the claims
func (m *ClerkMiddleware) getUserFromClaims(ctx context.Context, claims *SessionClaims) (*model.User, error) {
	// Try to get user from auth service first
	user, err := m.authService.GetUserByID(ctx, claims.Subject)
	if err == nil {
		return user, nil
	}

	// If not found, create a new user from claims
	return model.NewUser(
		claims.Subject,
		claims.Email,
		claims.Name,
	), nil
}

// getRolesFromClaims extracts roles from the claims
func (m *ClerkMiddleware) getRolesFromClaims(ctx context.Context, claims *SessionClaims) ([]string, error) {
	// If roles are in the claims, use them
	if claims.Roles != nil && len(claims.Roles) > 0 {
		return claims.Roles, nil
	}

	// Otherwise, try to get roles from the auth service
	roles, err := m.authService.GetUserRoles(ctx, claims.Subject)
	if err != nil {
		return []string{"user"}, err
	}

	return roles, nil
}

// RequireAuthentication is a middleware that requires authentication
func (m *ClerkMiddleware) RequireAuthentication(next http.Handler) http.Handler {
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
func (m *ClerkMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if user has the required role
			roles, ok := GetRolesFromContext(r.Context())
			if !ok {
				m.logger.Debug().Msg("Roles not found in context")
				appErr := apperror.NewUnauthorized("Authentication required", nil)
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
				appErr := apperror.NewForbidden("Insufficient permissions", nil)
				apperror.WriteError(w, appErr)
				return
			}

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}
