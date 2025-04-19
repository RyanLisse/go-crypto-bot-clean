package gateway

import (
	"context"
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

	// This is a placeholder implementation
	// In a real implementation, you would call the Clerk API to get user details
	// For now, we'll return a mock user
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

// verifyToken verifies a Clerk session token and returns the claims
func (g *ClerkGateway) verifyToken(ctx context.Context, token string) (*ClerkClaims, error) {
	// In a production environment, you should properly verify the token using JWT verification
	// For now, we'll implement a simplified version that extracts the subject claim

	// This is a placeholder implementation
	// In a real implementation, you would verify the token signature and expiration

	// For testing purposes, we'll accept any non-empty token and return a mock subject
	if token == "" {
		return nil, apperror.NewUnauthorized("Empty token", nil)
	}

	// Return mock claims for now
	return &ClerkClaims{
		Subject: "user_123", // This would normally be extracted from the verified token
		Exp:     time.Now().Add(time.Hour).Unix(),
	}, nil
}
