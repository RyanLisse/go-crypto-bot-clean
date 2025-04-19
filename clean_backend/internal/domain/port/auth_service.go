package port

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
)

// AuthServiceInterface defines the interface for authentication-related domain services.
type AuthServiceInterface interface {
	// GetUserByID retrieves a user by their ID.
	GetUserByID(ctx context.Context, userID string) (*model.User, error)

	// GetUserRoles retrieves the roles for a given user ID.
	GetUserRoles(ctx context.Context, userID string) ([]string, error)

	// GetUserFromToken retrieves a user based on an authentication token.
	// Note: This method might be refactored in Task 4 to align with actual Clerk token verification.
	GetUserFromToken(ctx context.Context, token string) (*model.User, error)

	// ValidateToken validates an authentication token
	ValidateToken(ctx context.Context, token string) (bool, error)

	// GenerateToken generates an authentication token for a user
	GenerateToken(ctx context.Context, userID string) (string, error)

	// RevokeToken revokes an authentication token
	RevokeToken(ctx context.Context, token string) error
}
