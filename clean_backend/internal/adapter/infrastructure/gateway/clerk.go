package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/apperror"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port/gateway"
	"github.com/rs/zerolog"
)

// Ensure ClerkGateway implements the gateway.ClerkGateway interface.
var _ gateway.ClerkGateway = (*ClerkGateway)(nil)

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

// ClerkGateway implements the gateway.ClerkGateway interface
type ClerkGateway struct {
	config     *config.Config
	logger     *zerolog.Logger
	httpClient *http.Client
	apiKey     string
	userCache  map[string]*model.User
	cacheMutex sync.RWMutex
}

// NewClerkGateway creates a new adapter for interacting with Clerk.
func NewClerkGateway(config *config.Config, logger *zerolog.Logger) gateway.ClerkGateway {
	// Check if Clerk secret key is set
	if config.Auth.ClerkSecretKey == "" {
		logger.Warn().Msg("CLERK_SECRET_KEY not set, authentication will not work properly")
	}

	return &ClerkGateway{
		config:     config,
		logger:     logger,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		apiKey:     config.Auth.ClerkSecretKey,
		userCache:  make(map[string]*model.User),
		cacheMutex: sync.RWMutex{},
	}
}

// VerifySession checks the validity of a session token and returns user details
func (g *ClerkGateway) VerifySession(ctx context.Context, sessionToken string) (*model.User, error) {
	g.logger.Debug().Msg("Verifying session token with Clerk")

	// Verify token and extract claims
	claims, err := g.verifyToken(ctx, sessionToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify session token: %w", err)
	}

	// Get user by ID from claims
	return g.GetUser(ctx, claims.Subject)
}

// GetUser retrieves user details by their Clerk User ID
func (g *ClerkGateway) GetUser(ctx context.Context, clerkUserID string) (*model.User, error) {
	g.logger.Debug().Str("user_id", clerkUserID).Msg("Getting user from Clerk")

	// Check cache first
	g.cacheMutex.RLock()
	if user, ok := g.userCache[clerkUserID]; ok {
		g.cacheMutex.RUnlock()
		return user, nil
	}
	g.cacheMutex.RUnlock()

	// Check if we're in development mode with mock auth enabled
	if g.config.Mock.AuthService {
		g.logger.Warn().Msg("Using mock authentication - returning mock user")
		user := &model.User{
			ID:        clerkUserID,
			Email:     "user@example.com",
			Name:      "Test User",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Cache user
		g.cacheMutex.Lock()
		g.userCache[clerkUserID] = user
		g.cacheMutex.Unlock()

		return user, nil
	}

	// Get user from Clerk API
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		"https://api.clerk.com/v1/users/"+clerkUserID,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authorization header with Clerk secret key
	req.Header.Add("Authorization", "Bearer "+g.apiKey)
	req.Header.Add("Content-Type", "application/json")

	// Send the request
	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to Clerk API: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		var errorResp struct {
			Errors []struct {
				Message string `json:"message"`
			} `json:"errors"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return nil, apperror.NewNotFound("User not found", clerkUserID, err)
		}

		var errMsg string
		if len(errorResp.Errors) > 0 {
			errMsg = errorResp.Errors[0].Message
		} else {
			errMsg = "User not found"
		}

		return nil, apperror.NewNotFound(errMsg, clerkUserID, nil)
	}

	// Parse the response
	var userResp struct {
		Data struct {
			ID             string `json:"id"`
			FirstName      string `json:"first_name"`
			LastName       string `json:"last_name"`
			EmailAddresses []struct {
				EmailAddress string `json:"email_address"`
				Verified     bool   `json:"verified"`
			} `json:"email_addresses"`
			CreatedAt int64 `json:"created_at"`
			UpdatedAt int64 `json:"updated_at"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		return nil, fmt.Errorf("failed to decode Clerk API response: %w", err)
	}

	// Extract email address (use the first verified one)
	var email string
	for _, emailAddr := range userResp.Data.EmailAddresses {
		if emailAddr.Verified {
			email = emailAddr.EmailAddress
			break
		}
	}

	// If no verified email found, use the first one
	if email == "" && len(userResp.Data.EmailAddresses) > 0 {
		email = userResp.Data.EmailAddresses[0].EmailAddress
	}

	// Construct full name
	name := userResp.Data.FirstName
	if userResp.Data.LastName != "" {
		if name != "" {
			name += " "
		}
		name += userResp.Data.LastName
	}

	// If name is empty, use email as name
	if name == "" {
		name = email
	}

	// Create user model
	user := &model.User{
		ID:        userResp.Data.ID,
		Email:     email,
		Name:      name,
		CreatedAt: time.Unix(userResp.Data.CreatedAt/1000, 0), // Convert from milliseconds
		UpdatedAt: time.Unix(userResp.Data.UpdatedAt/1000, 0), // Convert from milliseconds
	}

	// Cache user
	g.cacheMutex.Lock()
	g.userCache[clerkUserID] = user
	g.cacheMutex.Unlock()

	return user, nil
}

// verifyToken verifies a Clerk session token and returns the claims
func (g *ClerkGateway) verifyToken(ctx context.Context, token string) (*ClerkClaims, error) {
	if token == "" {
		return nil, apperror.NewUnauthorized("Empty token", nil)
	}

	// Check if we're in development mode with mock auth enabled
	if g.config.Mock.AuthService {
		g.logger.Warn().Msg("Using mock authentication - token verification bypassed")
		return &ClerkClaims{
			Subject: "user_123", // Mock user ID for development
			Exp:     time.Now().Add(time.Hour).Unix(),
		}, nil
	}

	// Verify token with Clerk API
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		"https://api.clerk.com/v1/sessions/verify",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add the session token as a query parameter
	q := req.URL.Query()
	q.Add("session_token", token)
	req.URL.RawQuery = q.Encode()

	// Add authorization header with Clerk secret key
	req.Header.Add("Authorization", "Bearer "+g.apiKey)
	req.Header.Add("Content-Type", "application/json")

	// Send the request
	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to Clerk API: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		var errorResp struct {
			Errors []struct {
				Message string `json:"message"`
			} `json:"errors"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return nil, apperror.NewUnauthorized("Invalid session token", nil)
		}

		var errMsg string
		if len(errorResp.Errors) > 0 {
			errMsg = errorResp.Errors[0].Message
		} else {
			errMsg = "Invalid session token"
		}

		return nil, apperror.NewUnauthorized(errMsg, nil)
	}

	// Parse the response
	var verifyResp struct {
		Data struct {
			ID        string `json:"id"`
			UserID    string `json:"user_id"`
			ExpiresAt int64  `json:"expires_at"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&verifyResp); err != nil {
		return nil, fmt.Errorf("failed to decode Clerk API response: %w", err)
	}

	// Check if the session is valid
	if verifyResp.Data.ID == "" || verifyResp.Data.UserID == "" {
		return nil, apperror.NewUnauthorized("Invalid session data", nil)
	}

	// Check if the session has expired
	if verifyResp.Data.ExpiresAt < time.Now().Unix() {
		return nil, apperror.NewUnauthorized("Session has expired", nil)
	}

	// Return the claims
	return &ClerkClaims{
		Subject: verifyResp.Data.UserID,
		Exp:     verifyResp.Data.ExpiresAt,
	}, nil
}
