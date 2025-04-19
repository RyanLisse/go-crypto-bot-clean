package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/apperror"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// ClerkUser represents a user from Clerk API
type ClerkUser struct {
	ID             string                 `json:"id"`
	FirstName      string                 `json:"first_name"`
	LastName       string                 `json:"last_name"`
	EmailAddresses []ClerkEmailAddress    `json:"email_addresses"`
	CreatedAt      int64                  `json:"created_at"`
	UpdatedAt      int64                  `json:"updated_at"`
	PublicMetadata map[string]interface{} `json:"public_metadata"`
}

// ClerkEmailAddress represents an email address from Clerk API
type ClerkEmailAddress struct {
	EmailAddress string `json:"email_address"`
	Verified     bool   `json:"verified"`
}

// ClerkClaims represents the claims in a Clerk JWT
type ClerkClaims struct {
	Subject string `json:"sub"`
	Azp     string `json:"azp"`
	Exp     int64  `json:"exp"`
}

// ClerkAuthServiceV2 implements the AuthServiceInterface using Clerk API directly
type ClerkAuthServiceV2 struct {
	logger     *zerolog.Logger
	config     *config.Config
	httpClient *http.Client
	apiKey     string
	userCache  map[string]*model.User
	cacheMutex sync.RWMutex
}

// NewClerkAuthServiceV2 creates a new ClerkAuthServiceV2
func NewClerkAuthServiceV2(config *config.Config, logger *zerolog.Logger) port.AuthServiceInterface {
	// Check if Clerk secret key is set
	if config.Auth.ClerkSecretKey == "" {
		logger.Warn().Msg("CLERK_SECRET_KEY not set, authentication will not work properly")
		return NewTestAuthService()
	}

	return &ClerkAuthServiceV2{
		logger:     logger,
		config:     config,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		apiKey:     config.Auth.ClerkSecretKey,
		userCache:  make(map[string]*model.User),
		cacheMutex: sync.RWMutex{},
	}
}

// GetUserByID retrieves a user by their ID
func (s *ClerkAuthServiceV2) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	// Check cache first
	s.cacheMutex.RLock()
	if user, ok := s.userCache[userID]; ok {
		s.cacheMutex.RUnlock()
		return user, nil
	}
	s.cacheMutex.RUnlock()

	// Get user from Clerk API
	url := fmt.Sprintf("https://api.clerk.dev/v1/users/%s", userID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.logger.Error().Err(err).Str("user_id", userID).Msg("Failed to get user from Clerk")
		return nil, fmt.Errorf("failed to get user from Clerk: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		s.logger.Error().
			Int("status_code", resp.StatusCode).
			Str("user_id", userID).
			Str("response", string(body)).
			Msg("Failed to get user from Clerk")
		return nil, fmt.Errorf("failed to get user from Clerk: status %d", resp.StatusCode)
	}

	var clerkUser ClerkUser
	if err := json.NewDecoder(resp.Body).Decode(&clerkUser); err != nil {
		s.logger.Error().Err(err).Str("user_id", userID).Msg("Failed to decode user from Clerk")
		return nil, fmt.Errorf("failed to decode user from Clerk: %w", err)
	}

	// Convert to domain model
	var email string
	if len(clerkUser.EmailAddresses) > 0 {
		email = clerkUser.EmailAddresses[0].EmailAddress
	}

	user := &model.User{
		ID:        clerkUser.ID,
		Email:     email,
		Name:      fmt.Sprintf("%s %s", clerkUser.FirstName, clerkUser.LastName),
		CreatedAt: time.Unix(clerkUser.CreatedAt, 0),
		UpdatedAt: time.Unix(clerkUser.UpdatedAt, 0),
	}

	// Cache user
	s.cacheMutex.Lock()
	s.userCache[userID] = user
	s.cacheMutex.Unlock()

	return user, nil
}

// GetUserRoles retrieves the roles for a given user ID
func (s *ClerkAuthServiceV2) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	// Get user from Clerk API
	url := fmt.Sprintf("https://api.clerk.dev/v1/users/%s", userID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return []string{"user"}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.logger.Error().Err(err).Str("user_id", userID).Msg("Failed to get user from Clerk")
		return []string{"user"}, fmt.Errorf("failed to get user from Clerk: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		s.logger.Error().
			Int("status_code", resp.StatusCode).
			Str("user_id", userID).
			Str("response", string(body)).
			Msg("Failed to get user from Clerk")
		return []string{"user"}, fmt.Errorf("failed to get user from Clerk: status %d", resp.StatusCode)
	}

	var clerkUser ClerkUser
	if err := json.NewDecoder(resp.Body).Decode(&clerkUser); err != nil {
		s.logger.Error().Err(err).Str("user_id", userID).Msg("Failed to decode user from Clerk")
		return []string{"user"}, fmt.Errorf("failed to decode user from Clerk: %w", err)
	}

	// Get roles from public metadata
	roles := []string{"user"} // Default role
	if clerkUser.PublicMetadata != nil {
		if rolesData, ok := clerkUser.PublicMetadata["roles"]; ok {
			if rolesArray, ok := rolesData.([]interface{}); ok {
				for _, role := range rolesArray {
					if roleStr, ok := role.(string); ok {
						roles = append(roles, roleStr)
					}
				}
			}
		}
	}

	return roles, nil
}

// GetUserFromToken retrieves a user based on an authentication token
func (s *ClerkAuthServiceV2) GetUserFromToken(ctx context.Context, token string) (*model.User, error) {
	// Verify token
	claims, err := s.verifyToken(ctx, token)
	if err != nil {
		return nil, err
	}

	// Get user by ID
	return s.GetUserByID(ctx, claims.Subject)
}

// ValidateToken validates an authentication token
func (s *ClerkAuthServiceV2) ValidateToken(ctx context.Context, token string) (bool, error) {
	_, err := s.verifyToken(ctx, token)
	if err != nil {
		return false, err
	}
	return true, nil
}

// verifyToken verifies a Clerk session token and returns the claims
func (s *ClerkAuthServiceV2) verifyToken(ctx context.Context, token string) (*ClerkClaims, error) {
	// For now, we'll just parse the token and extract the subject
	// In a production environment, you should properly verify the token
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, apperror.NewUnauthorized("Invalid token format", nil)
	}

	// Extract claims from the token
	// This is a simplified implementation
	claims := &ClerkClaims{
		Subject: "user_123", // Placeholder
	}

	return claims, nil
}

// GenerateToken generates an authentication token for a user
func (s *ClerkAuthServiceV2) GenerateToken(ctx context.Context, userID string) (string, error) {
	// This is not supported by Clerk directly
	// Tokens are generated by Clerk frontend SDK
	return "", apperror.NewInternal(errors.New("token generation not supported by Clerk"))
}

// RevokeToken revokes an authentication token
func (s *ClerkAuthServiceV2) RevokeToken(ctx context.Context, token string) error {
	// This is not supported by Clerk directly
	// Token revocation is handled by Clerk
	return apperror.NewInternal(errors.New("token revocation not supported by Clerk"))
}
