package middleware

import (
"context"
"fmt"
"net/http"
"strings"
"sync"
"time"

"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/adapter/delivery/http/response"
"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/config"
"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
domainport "github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port"
portgateway "github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port/gateway"

"github.com/rs/zerolog"
)

// ClerkMiddleware is the primary authentication middleware using Clerk
type ClerkMiddleware struct {
	logger       *zerolog.Logger
	authService  domainport.AuthServiceInterface
	clerkGateway portgateway.ClerkGateway
	config       *config.Config
	jwkCache     map[string]any // Cache for JWKs
	jwkCacheMu   sync.RWMutex   // Mutex for JWK cache
	jwkCacheExp  time.Time      // Expiration time for JWK cache
}

// NewClerkMiddleware creates a new ClerkMiddleware
func NewClerkMiddleware(authService domainport.AuthServiceInterface, clerkGateway portgateway.ClerkGateway, config *config.Config, logger *zerolog.Logger) AuthMiddleware {
	// Check if Clerk secret key is set
	if config.Auth.ClerkSecretKey == "" {
		logger.Warn().Msg("CLERK_SECRET_KEY not set, authentication will not work properly")
	}

	return &ClerkMiddleware{
		logger:       logger,
		authService:  authService,
		clerkGateway: clerkGateway,
		config:       config,
		jwkCache:     make(map[string]any),
		jwkCacheExp:  time.Now(),
	}
}

// SessionClaims represents the claims in a Clerk session token
type SessionClaims struct {
	Subject   string   `json:"sub"`
	SessionID string   `json:"sid"`
	UserID    string   `json:"userId"`
	Email     string   `json:"email"`
	Name      string   `json:"name"`
	Roles     []string `json:"roles"`
	Expiry    int64    `json:"exp"`
	IssuedAt  int64    `json:"iat"`
	NotBefore int64    `json:"nbf"`
}

// Middleware returns a middleware function that extracts and validates the auth token
func (m *ClerkMiddleware) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// Get the auth header
authHeader := r.Header.Get("Authorization")
if authHeader == "" {
// No auth header, continue without authentication
next.ServeHTTP(w, r)
return
}

// Extract the token
tokenParts := strings.Split(authHeader, " ")
if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
// Invalid auth header format, continue without authentication
next.ServeHTTP(w, r)
return
}

token := tokenParts[1]
if token == "" {
// Empty token, continue without authentication
next.ServeHTTP(w, r)
return
}

// Verify the token
claims, err := m.verifyToken(r.Context(), token)
			if err != nil {
				// Invalid token, continue without authentication
				m.logger.Debug().Err(err).Msg("Failed to verify token")
				next.ServeHTTP(w, r)
				return
			}

			// Get user from claims
			user, err := m.getUserFromClaims(r.Context(), claims)
			if err != nil {
				m.logger.Error().Err(err).Msg("Failed to get user from claims")
				next.ServeHTTP(w, r)
				return
			}

			// Get roles from claims
			roles := claims.Roles
			if len(roles) == 0 {
				// If no roles in claims, get from auth service
				roles, err = m.authService.GetUserRoles(r.Context(), claims.UserID)
				if err != nil {
					m.logger.Error().Err(err).Msg("Failed to get user roles")
					// Continue with empty roles
					roles = []string{"user"}
				}
			}

			// Set user ID, roles, and claims in context
			ctx := SetUserIDInContext(r.Context(), user.ID)
			ctx = SetUserInContext(ctx, user)
			ctx = SetRolesInContext(ctx, roles)
			ctx = context.WithValue(ctx, "clerk_claims", claims)

			// Call the next handler with the new context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAuthentication is a middleware that requires authentication
func (m *ClerkMiddleware) RequireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// Check if user ID is in context
userID, ok := GetUserIDFromContext(r.Context())
		if !ok || userID == "" {
			m.logger.Debug().Msg("Authentication required but user ID not found in context")
			response.WriteErrorJSON(w, http.StatusUnauthorized, "Authentication required", m.logger)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// RequireRole returns a middleware that requires the user to have the specified role
func (m *ClerkMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// Check if user ID is in context
userID, ok := GetUserIDFromContext(r.Context())
			if !ok || userID == "" {
				m.logger.Debug().Msg("Authentication required but user ID not found in context")
				response.WriteErrorJSON(w, http.StatusUnauthorized, "Authentication required", m.logger)
				return
			}

			// Check if user has the required role
			roles, ok := GetRolesFromContext(r.Context())
			if !ok {
				m.logger.Debug().Msg("Roles not found in context")
				response.WriteErrorJSON(w, http.StatusForbidden, "Forbidden", m.logger)
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
				m.logger.Debug().Str("role", role).Msg("User does not have the required role")
				response.WriteErrorJSON(w, http.StatusForbidden, "Forbidden", m.logger)
				return
			}

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}

// verifyToken verifies a Clerk session token and returns the claims
func (m *ClerkMiddleware) verifyToken(ctx context.Context, token string) (*SessionClaims, error) {
	// Use the Clerk gateway to verify the session token
	user, err := m.clerkGateway.VerifySession(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
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
	// Try to get user from Clerk gateway first
	user, err := m.clerkGateway.GetUser(ctx, claims.Subject)
	if err == nil {
		return user, nil
	}

	// If not found, try the auth service
	user, err = m.authService.GetUserByID(ctx, claims.Subject)
	if err == nil {
		return user, nil
	}

	// If still not found, create a new user from claims
	return model.NewUser(
claims.Subject,
claims.Email,
claims.Name,
), nil
}
