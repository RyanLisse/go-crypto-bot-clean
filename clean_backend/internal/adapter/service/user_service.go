package service

import (
	"context"
	"fmt"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/apperror"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// UserService implements the UserServiceInterface
type UserService struct {
	logger      *zerolog.Logger
	config      *config.Config
	userRepo    port.UserRepository
	authService port.AuthServiceInterface
}

// NewUserService creates a new UserService
func NewUserService(
	logger *zerolog.Logger,
	config *config.Config,
	userRepo port.UserRepository,
	authService port.AuthServiceInterface,
) port.UserServiceInterface {
	return &UserService{
		logger:      logger,
		config:      config,
		userRepo:    userRepo,
		authService: authService,
	}
}

// GetUserByID retrieves a user by their ID
func (s *UserService) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	if userID == "" {
		return nil, apperror.NewBadRequest("User ID is required", nil, nil)
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		s.logger.Error().Err(err).Str("userID", userID).Msg("Failed to get user by ID")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetUserByEmail retrieves a user by their email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	if email == "" {
		return nil, apperror.NewBadRequest("Email is required", nil, nil)
	}

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		s.logger.Error().Err(err).Str("email", email).Msg("Failed to get user by email")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	if user == nil {
		return nil, apperror.NewBadRequest("User is required", nil, nil)
	}

	// Validate user
	if err := user.Validate(); err != nil {
		return nil, apperror.NewBadRequest("Invalid user data", nil, err)
	}

	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		return nil, apperror.NewConflict("User with email already exists", nil, nil)
	}

	// Set creation time
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Save user
	if err := s.userRepo.Save(ctx, user); err != nil {
		s.logger.Error().Err(err).Str("email", user.Email).Msg("Failed to create user")
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(ctx context.Context, user *model.User) (*model.User, error) {
	if user == nil {
		return nil, apperror.NewBadRequest("User is required", nil, nil)
	}

	// Validate user
	if err := user.Validate(); err != nil {
		return nil, apperror.NewBadRequest("Invalid user data", nil, err)
	}

	// Check if user exists
	existingUser, err := s.userRepo.GetByID(ctx, user.ID)
	if err != nil {
		s.logger.Error().Err(err).Str("userID", user.ID).Msg("Failed to get user for update")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Update fields
	existingUser.Name = user.Name
	existingUser.UpdatedAt = time.Now()

	// Save user
	if err := s.userRepo.Save(ctx, existingUser); err != nil {
		s.logger.Error().Err(err).Str("userID", user.ID).Msg("Failed to update user")
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return existingUser, nil
}

// DeleteUser deletes a user by their ID
func (s *UserService) DeleteUser(ctx context.Context, userID string) error {
	if userID == "" {
		return apperror.NewBadRequest("User ID is required", nil, nil)
	}

	// Check if user exists
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		s.logger.Error().Err(err).Str("userID", userID).Msg("Failed to get user for deletion")
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Delete user
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		s.logger.Error().Err(err).Str("userID", userID).Msg("Failed to delete user")
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// ListUsers lists all users
func (s *UserService) ListUsers(ctx context.Context) ([]*model.User, error) {
	users, err := s.userRepo.List(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to list users")
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}

// SetUserRole sets a role for a user
func (s *UserService) SetUserRole(ctx context.Context, userID string, role string) error {
	// This is a placeholder implementation
	// In a real implementation, this would call a role repository or auth service
	s.logger.Info().Str("userID", userID).Str("role", role).Msg("Setting user role")
	return nil
}

// RemoveUserRole removes a role from a user
func (s *UserService) RemoveUserRole(ctx context.Context, userID string, role string) error {
	// This is a placeholder implementation
	// In a real implementation, this would call a role repository or auth service
	s.logger.Info().Str("userID", userID).Str("role", role).Msg("Removing user role")
	return nil
}
