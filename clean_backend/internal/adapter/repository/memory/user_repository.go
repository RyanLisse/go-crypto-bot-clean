package memory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/apperror"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// UserRepository is an in-memory implementation of the UserRepository interface
type UserRepository struct {
	logger *zerolog.Logger
	mu     sync.RWMutex
	users  map[string]*model.User
}

// NewUserRepository creates a new in-memory UserRepository
func NewUserRepository(logger *zerolog.Logger) port.UserRepository {
	return &UserRepository{
		logger: logger,
		users:  make(map[string]*model.User),
	}
}

// Save saves a user (creates or updates)
func (r *UserRepository) Save(ctx context.Context, user *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// If the user doesn't have an ID, generate one
	if user.ID == "" {
		user.ID = uuid.New().String()
		user.CreatedAt = time.Now()
	}

	// Update the updated_at timestamp
	user.UpdatedAt = time.Now()

	// Save the user
	r.users[user.ID] = user

	return nil
}

// GetByID retrieves a user by their ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return nil, apperror.ErrNotFound
	}

	return user, nil
}

// GetByEmail retrieves a user by their email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, apperror.ErrNotFound
}

// Delete deletes a user by their ID
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[id]; !ok {
		return fmt.Errorf("user not found: %w", apperror.ErrNotFound)
	}

	delete(r.users, id)

	return nil
}

// List lists all users
func (r *UserRepository) List(ctx context.Context) ([]*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*model.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	return users, nil
}
