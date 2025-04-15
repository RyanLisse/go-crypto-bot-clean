package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/apperror"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
)

// TestAuthMiddleware is a middleware for testing authentication
type TestAuthMiddleware struct {
	logger *zerolog.Logger
	secret string
}

// NewTestAuthMiddleware creates a new TestAuthMiddleware
func NewTestAuthMiddleware(secret string, logger *zerolog.Logger) *TestAuthMiddleware {
	return &TestAuthMiddleware{
		logger: logger,
		secret: secret,
	}
}

// Middleware returns a middleware function that validates test authentication tokens
func (m *TestAuthMiddleware) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check for Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				m.logger.Debug().Msg("No authorization header present")
				next.ServeHTTP(w, r)
				return
			}

			// Extract token
			tokenString := authHeader
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}

			// Parse the token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Validate the signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(m.secret), nil
			})

			if err != nil {
				m.logger.Error().Err(err).Msg("Failed to parse token")
				apperror.WriteError(w, apperror.NewUnauthorized("Invalid authentication token", err))
				return
			}

			// Validate the token
			if !token.Valid {
				m.logger.Error().Msg("Invalid token")
				apperror.WriteError(w, apperror.NewUnauthorized("Invalid authentication token", nil))
				return
			}

			// Extract claims
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				m.logger.Error().Msg("Failed to extract claims")
				apperror.WriteError(w, apperror.NewUnauthorized("Invalid authentication token", nil))
				return
			}

			// Validate expiration
			if exp, expOk := claims["exp"].(float64); expOk {
				if time.Now().Unix() > int64(exp) {
					m.logger.Error().Msg("Token expired")
					apperror.WriteError(w, apperror.NewUnauthorized("Token expired", nil))
					return
				}
			}

			// Get user ID from claims
			userID, ok := claims["sub"].(string)
			if !ok || userID == "" {
				m.logger.Error().Msg("Missing user ID in claims")
				apperror.WriteError(w, apperror.NewUnauthorized("Invalid authentication token", nil))
				return
			}

			// Set user ID in context
			ctx := context.WithValue(r.Context(), UserIDKey, userID)

			// Extract roles from claims if available
			if rolesInterface, ok := claims["roles"].([]interface{}); ok {
				roleStrings := make([]string, 0, len(rolesInterface))
				for _, role := range rolesInterface {
					if roleStr, ok := role.(string); ok {
						roleStrings = append(roleStrings, roleStr)
					}
				}
				ctx = context.WithValue(ctx, RoleKey, roleStrings)
			} else {
				// Default role
				ctx = context.WithValue(ctx, RoleKey, []string{"user"})
			}

			// Store claims in context
			ctx = context.WithValue(ctx, "jwt_claims", claims)

			// Call the next handler with updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
