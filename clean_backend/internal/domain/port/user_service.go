package port

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
)

// UserServiceInterface defines the interface for user-related business logic
type UserServiceInterface interface {
	// GetUserByID retrieves a user by their ID
	GetUserByID(ctx context.Context, userID string) (*model.User, error)

	// GetUserByEmail retrieves a user by their email
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)

	// CreateUser creates a new user
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)

	// UpdateUser updates an existing user
	UpdateUser(ctx context.Context, user *model.User) (*model.User, error)

	// DeleteUser deletes a user by their ID
	DeleteUser(ctx context.Context, userID string) error

	// ListUsers lists all users
	ListUsers(ctx context.Context) ([]*model.User, error)

	// SetUserRole sets a role for a user
	SetUserRole(ctx context.Context, userID string, role string) error

	// RemoveUserRole removes a role from a user
	RemoveUserRole(ctx context.Context, userID string, role string) error
}
