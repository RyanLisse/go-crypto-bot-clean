package gateway

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model" // Adjust if user model is elsewhere
)

// ClerkGateway defines the port for interacting with the Clerk authentication service.
type ClerkGateway interface {
	// VerifySession checks the validity of a session token and returns user details.
	VerifySession(ctx context.Context, sessionToken string) (*model.User, error) // Assuming a User model exists

	// GetUser retrieves user details by their Clerk User ID.
	GetUser(ctx context.Context, clerkUserID string) (*model.User, error)

	// TODO: Add other Clerk interactions as needed (e.g., updating user metadata, organization checks).
}
