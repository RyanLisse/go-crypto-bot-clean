package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/adapter/infrastructure/persistence/gorm/entity"
	domainmodel "github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
	domainport "github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// UserRepository implements the port.UserRepository interface using GORM
type UserRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *gorm.DB, logger *zerolog.Logger) *UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

// Save saves a user
func (r *UserRepository) Save(ctx context.Context, user *domainmodel.User) error {
	// Convert domain model to entity
	userEntity := &entity.UserEntity{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	// Save entity
	if err := r.db.WithContext(ctx).Save(userEntity).Error; err != nil {
		r.logger.Error().Err(err).Str("userID", user.ID).Msg("Failed to save user")
		return fmt.Errorf("failed to save user: %w", err)
	}

	return nil
}

// GetByID gets a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*domainmodel.User, error) {
	// Get entity
	var userEntity entity.UserEntity
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&userEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainmodel.ErrInvalidUserID
		}
		r.logger.Error().Err(err).Str("userID", id).Msg("Failed to get user")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Convert entity to domain model
	user := &domainmodel.User{
		ID:        userEntity.ID,
		Email:     userEntity.Email,
		Name:      userEntity.Name,
		CreatedAt: userEntity.CreatedAt,
		UpdatedAt: userEntity.UpdatedAt,
	}

	return user, nil
}

// GetByEmail gets a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domainmodel.User, error) {
	// Get entity
	var userEntity entity.UserEntity
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&userEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found with email: %s", email)
		}
		r.logger.Error().Err(err).Str("email", email).Msg("Failed to get user by email")
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	// Convert entity to domain model
	user := &domainmodel.User{
		ID:        userEntity.ID,
		Email:     userEntity.Email,
		Name:      userEntity.Name,
		CreatedAt: userEntity.CreatedAt,
		UpdatedAt: userEntity.UpdatedAt,
	}

	return user, nil
}

// Delete deletes a user by ID
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	// Delete entity
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.UserEntity{}).Error; err != nil {
		r.logger.Error().Err(err).Str("userID", id).Msg("Failed to delete user")
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// List lists all users
func (r *UserRepository) List(ctx context.Context) ([]*domainmodel.User, error) {
	// Get entities
	var userEntities []entity.UserEntity
	if err := r.db.WithContext(ctx).Find(&userEntities).Error; err != nil {
		r.logger.Error().Err(err).Msg("Failed to list users")
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	// Convert entities to domain models
	users := make([]*domainmodel.User, len(userEntities))
	for i, userEntity := range userEntities {
		users[i] = &domainmodel.User{
			ID:        userEntity.ID,
			Email:     userEntity.Email,
			Name:      userEntity.Name,
			CreatedAt: userEntity.CreatedAt,
			UpdatedAt: userEntity.UpdatedAt,
		}
	}

	return users, nil
}

// Ensure UserRepository implements port.UserRepository
var _ domainport.UserRepository = (*UserRepository)(nil)
