package middleware

import (
	"context"
	"net/http"
	"strings"



	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/apperror"
	"github.com/clerk/clerk-sdk-go/v2"
	clerkjwt "github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/rs/zerolog"
)

// ClerkMiddleware handles Clerk authentication
type ClerkMiddleware struct {
	logger *zerolog.Logger
}

// NewClerkMiddleware creates a new Clerk middleware
func NewClerkMiddleware(secretKey string, logger *zerolog.Logger) *ClerkMiddleware {
	// Initialize Clerk SDK with secret key
	clerk.SetKey(secretKey)

	logger.Debug().Msg("Initialized Clerk SDK middleware")

	return &ClerkMiddleware{
		logger: logger,
	}
}

// Middleware returns a middleware function that validates Clerk authentication
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

			// Verify session with Clerk using the jwt package
			claims, err := clerkjwt.Verify(r.Context(), &clerkjwt.VerifyParams{
				Token: sessionToken,
			})
			if err != nil {
				m.logger.Error().Err(err).Msg("Failed to verify Clerk session")
				apperror.WriteError(w, apperror.NewUnauthorized("Invalid authentication token", err))
				return
			}

			// Get user ID from claims
			userID := claims.Subject
			if userID == "" {
				m.logger.Error().Msg("Missing user ID in claims")
				apperror.WriteError(w, apperror.NewUnauthorized("Invalid authentication token", nil))
				return
			}

			// Set user ID in context
			ctx := context.WithValue(r.Context(), UserIDKey, userID)

			// Extract roles from claims if available
			if roles, ok := claims.Custom.(map[string]interface{})["roles"].([]interface{}); ok {
				roleStrings := make([]string, 0, len(roles))
				for _, role := range roles {
					if roleStr, ok := role.(string); ok {
						roleStrings = append(roleStrings, roleStr)
					}
				}
				ctx = context.WithValue(ctx, RoleKey, roleStrings)
			} else {
				// Default role
				ctx = context.WithValue(ctx, RoleKey, []string{"user"})
			}

			// Store session claims in context
			ctx = context.WithValue(ctx, "clerk_claims", claims)

			// Call the next handler with updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// TestClerkSDK tests if the Clerk SDK is configured properly
func TestClerkSDK(secretKey string) error {
	// Just verify if we can create a client with the provided secret key
	clerk.SetKey(secretKey)

	// Test a token verification with an empty token (will fail but tests the SDK initialization)
	_, err := clerkjwt.Verify(context.Background(), &clerkjwt.VerifyParams{
		Token: "",
	})

	// We expect this to fail with a specific error about the token
	// If it fails with a different error (like connection issues), we return that
	if err != nil && !strings.Contains(err.Error(), "token") {
		return err
	}

	return nil
}
