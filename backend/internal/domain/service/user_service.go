package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
)

// UserServiceInterface defines the interface for user-related operations
type UserServiceInterface interface {
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	CreateUser(ctx context.Context, id, email, name string) (*model.User, error)
	UpdateUser(ctx context.Context, id, name string) (*model.User, error)
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context) ([]*model.User, error)
	EnsureUserExists(ctx context.Context, id, email, name string) (*model.User, error)
}

// UserService handles user-related operations
type UserService struct {
	userRepo port.UserRepository
}

// NewUserService creates a new UserService
func NewUserService(userRepo port.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetUserByID gets a user by ID
func (s *UserService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	if id == "" {
		return nil, errors.New("user ID is required")
	}

	return s.userRepo.GetByID(ctx, id)
}

// GetUserByEmail gets a user by email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}

	return s.userRepo.GetByEmail(ctx, email)
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, id, email, name string) (*model.User, error) {
	if id == "" {
		return nil, errors.New("user ID is required")
	}

	if email == "" {
		return nil, errors.New("email is required")
	}

	// Check if user already exists
	existingUser, err := s.userRepo.GetByID(ctx, id)
	if err == nil && existingUser != nil {
		return existingUser, nil
	}

	// Create new user
	user := &model.User{
		ID:        id,
		Email:     email,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save user
	if err := s.userRepo.Save(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	return user, nil
}

// UpdateUser updates a user
func (s *UserService) UpdateUser(ctx context.Context, id, name string) (*model.User, error) {
	if id == "" {
		return nil, errors.New("user ID is required")
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Update user
	if name != "" {
		user.Name = name
	}
	user.UpdatedAt = time.Now()

	// Save user
	if err := s.userRepo.Save(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	return user, nil
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("user ID is required")
	}

	return s.userRepo.Delete(ctx, id)
}

// ListUsers lists all users
func (s *UserService) ListUsers(ctx context.Context) ([]*model.User, error) {
	return s.userRepo.List(ctx)
}

// EnsureUserExists ensures that a user exists in the database
// If the user doesn't exist, it creates a new user
func (s *UserService) EnsureUserExists(ctx context.Context, id, email, name string) (*model.User, error) {
	// Try to get user by ID
	user, err := s.userRepo.GetByID(ctx, id)
	if err == nil && user != nil {
		// User exists, update if needed
		if (email != "" && user.Email != email) || (name != "" && user.Name != name) {
			if email != "" {
				user.Email = email
			}
			if name != "" {
				user.Name = name
			}
			user.UpdatedAt = time.Now()
			if err := s.userRepo.Save(ctx, user); err != nil {
				return nil, fmt.Errorf("failed to update existing user: %w", err)
			}
		}
		return user, nil
	}

	// User doesn't exist, create new user
	return s.CreateUser(ctx, id, email, name)
}
