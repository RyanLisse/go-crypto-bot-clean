package port

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
)

// UserRepository defines the interface for user repository
type UserRepository interface {
	// Save saves a user
	Save(ctx context.Context, user *model.User) error

	// GetByID gets a user by ID
	GetByID(ctx context.Context, id string) (*model.User, error)

	// GetByEmail gets a user by email
	GetByEmail(ctx context.Context, email string) (*model.User, error)

	// Delete deletes a user by ID
	Delete(ctx context.Context, id string) error

	// List lists all users
	List(ctx context.Context) ([]*model.User, error)
}
