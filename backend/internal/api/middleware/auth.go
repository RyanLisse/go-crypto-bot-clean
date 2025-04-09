// Package middleware provides middleware functions for the API.
package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"go-crypto-bot-clean/backend/internal/api/middleware/jwt"
)

// contextKey is a type for context keys
type contextKey string

// Context keys
const (
	UserIDKey contextKey = "userID"
	EmailKey  contextKey = "email"
	RolesKey  contextKey = "roles"
)

// AuthMiddleware provides authentication middleware
type AuthMiddleware struct {
	jwtService *jwt.Service
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(jwtService *jwt.Service) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

// Authenticate authenticates a request using JWT
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// Check if the Authorization header has the correct format
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			http.Error(w, "Token is required", http.StatusUnauthorized)
			return
		}

		// Validate the token
		claims, err := m.jwtService.ValidateAccessToken(tokenString)
		if err != nil {
			if errors.Is(err, jwt.ErrExpiredToken) {
				http.Error(w, "Token has expired", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Check if the token is blacklisted
		if m.jwtService.IsBlacklisted(tokenString) {
			http.Error(w, "Token has been revoked", http.StatusUnauthorized)
			return
		}

		// Add the user information to the request context
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, EmailKey, claims.Email)
		ctx = context.WithValue(ctx, RolesKey, claims.Roles)

		// Call the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole requires the user to have a specific role
func (m *AuthMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the roles from the context
			roles, ok := r.Context().Value(RolesKey).([]string)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check if the user has the required role
			hasRole := false
			for _, r := range roles {
				if r == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}

// RequireAdmin requires the user to have the admin role
func (m *AuthMiddleware) RequireAdmin() func(http.Handler) http.Handler {
	return m.RequireRole("admin")
}

// GetUserID gets the user ID from the request context
func GetUserID(r *http.Request) string {
	userID, ok := r.Context().Value(UserIDKey).(string)
	if !ok {
		return ""
	}
	return userID
}

// GetEmail gets the email from the request context
func GetEmail(r *http.Request) string {
	email, ok := r.Context().Value(EmailKey).(string)
	if !ok {
		return ""
	}
	return email
}

// GetRoles gets the roles from the request context
func GetRoles(r *http.Request) []string {
	roles, ok := r.Context().Value(RolesKey).([]string)
	if !ok {
		return nil
	}
	return roles
}

// HasRole checks if the user has a specific role
func HasRole(r *http.Request, role string) bool {
	roles := GetRoles(r)
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

// RegisterAuthMiddleware registers the authentication middleware with a router
func RegisterAuthMiddleware(r chi.Router, jwtService *jwt.Service) {
	authMiddleware := NewAuthMiddleware(jwtService)
	
	// Apply authentication middleware to all routes under /api/v1
	// except for the authentication routes
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.Authenticate)
		
		// Protected routes go here
		// For example:
		// r.Mount("/api/v1/user", userRouter)
		// r.Mount("/api/v1/strategy", strategyRouter)
		// r.Mount("/api/v1/backtest", backtestRouter)
	})
	
	// Public routes (no authentication required)
	// For example:
	// r.Mount("/api/v1/auth", authRouter)
}
