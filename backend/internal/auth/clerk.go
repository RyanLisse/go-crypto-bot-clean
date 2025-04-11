// Package auth provides authentication and authorization functionality
package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go-crypto-bot-clean/backend/internal/config"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"go.uber.org/zap"
)

// ClerkConfig represents the configuration for Clerk authentication
type ClerkConfig struct {
	Enabled       bool
	SecretKey     string
	JWKSEndpoint  string
	Issuer        string
	TokenLifetime time.Duration
}

// ClerkAuth handles authentication using Clerk
type ClerkAuth struct {
	config ClerkConfig
	logger *zap.Logger
	keySet jwk.Set
}

// NewClerkAuth creates a new Clerk authentication handler
func NewClerkAuth(config ClerkConfig, logger *zap.Logger) (*ClerkAuth, error) {
	if !config.Enabled {
		return &ClerkAuth{
			config: config,
			logger: logger,
		}, nil
	}

	if config.SecretKey == "" {
		return nil, errors.New("clerk secret key is required")
	}

	if config.JWKSEndpoint == "" {
		config.JWKSEndpoint = "https://clerk.your-app.com/.well-known/jwks.json"
	}

	// Fetch JWK set
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	keySet, err := jwk.Fetch(ctx, config.JWKSEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWK set: %w", err)
	}

	return &ClerkAuth{
		config: config,
		logger: logger,
		keySet: keySet,
	}, nil
}

// Middleware creates a middleware that validates Clerk JWT tokens
func (a *ClerkAuth) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !a.config.Enabled {
				next.ServeHTTP(w, r)
				return
			}

			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header is required", http.StatusUnauthorized)
				return
			}

			// Check if the header has the Bearer prefix
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			// Extract the token
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// Parse and validate the token
			token, err := a.ValidateToken(r.Context(), tokenString)
			if err != nil {
				a.logger.Error("Failed to validate token", zap.Error(err))
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Add the token to the request context
			ctx := context.WithValue(r.Context(), ContextKeyUser, token)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ValidateToken validates a JWT token
func (a *ClerkAuth) ValidateToken(ctx context.Context, tokenString string) (jwt.Token, error) {
	// Parse and validate the token
	token, err := jwt.Parse(
		[]byte(tokenString),
		jwt.WithKeySet(a.keySet),
		jwt.WithValidate(true),
		jwt.WithIssuer(a.config.Issuer),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	return token, nil
}

// GetUserFromContext extracts the user from the request context
func (a *ClerkAuth) GetUserFromContext(ctx context.Context) (jwt.Token, error) {
	user, ok := ctx.Value(ContextKeyUser).(jwt.Token)
	if !ok {
		return nil, errors.New("user not found in context")
	}
	return user, nil
}

// GetUserID extracts the user ID from the token
func (a *ClerkAuth) GetUserID(token jwt.Token) (string, error) {
	sub, ok := token.Get("sub")
	if !ok {
		return "", errors.New("subject not found in token")
	}

	subStr, ok := sub.(string)
	if !ok {
		return "", errors.New("subject is not a string")
	}

	return subStr, nil
}

// ContextKey is a type for context keys
type ContextKey string

// Context keys
const (
	ContextKeyUser ContextKey = "user"
)

// User represents a user in the system
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetUserProfile fetches the user profile from Clerk
func (a *ClerkAuth) GetUserProfile(ctx context.Context, userID string) (*User, error) {
	if !a.config.Enabled {
		return nil, errors.New("clerk is not enabled")
	}

	if a.config.SecretKey == "" {
		return nil, errors.New("clerk secret key is required")
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create request
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		fmt.Sprintf("https://api.clerk.dev/v1/users/%s", userID),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.config.SecretKey))
	req.Header.Add("Content-Type", "application/json")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse response
	var result struct {
		ID        string    `json:"id"`
		Email     string    `json:"email_address"`
		FirstName string    `json:"first_name"`
		LastName  string    `json:"last_name"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Convert to User
	user := &User{
		ID:        result.ID,
		Email:     result.Email,
		FirstName: result.FirstName,
		LastName:  result.LastName,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}

	return user, nil
}

// FromConfig creates a ClerkAuth from configuration
func FromConfig(cfg *config.Config, logger *zap.Logger) (*ClerkAuth, error) {
	return NewClerkAuth(ClerkConfig{
		Enabled:       cfg.Auth.Enabled,
		SecretKey:     cfg.Auth.ClerkSecretKey,
		JWKSEndpoint:  fmt.Sprintf("https://%s/.well-known/jwks.json", cfg.Auth.ClerkDomain),
		Issuer:        fmt.Sprintf("https://%s", cfg.Auth.ClerkDomain),
		TokenLifetime: time.Duration(cfg.Auth.JWTExpiry) * time.Hour,
	}, logger)
}

// FromMinimalConfig creates a ClerkAuth from minimal configuration
func FromMinimalConfig(cfg *config.MinimalConfig, logger *zap.Logger) (*ClerkAuth, error) {
	return NewClerkAuth(ClerkConfig{
		Enabled:       cfg.Auth.Enabled,
		SecretKey:     cfg.Auth.ClerkSecretKey,
		JWKSEndpoint:  fmt.Sprintf("https://%s/.well-known/jwks.json", cfg.Auth.ClerkDomain),
		Issuer:        fmt.Sprintf("https://%s", cfg.Auth.ClerkDomain),
		TokenLifetime: time.Duration(cfg.Auth.JWTExpiry) * time.Hour,
	}, logger)
}
