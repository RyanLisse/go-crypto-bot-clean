package port

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
)

// UserRepository defines the interface for user persistence
type UserRepository interface {
	Save(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*model.User, error)
}
