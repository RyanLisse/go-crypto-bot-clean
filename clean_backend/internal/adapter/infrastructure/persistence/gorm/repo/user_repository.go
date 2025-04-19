package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/adapter/infrastructure/persistence/gorm/entity"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// GormUserRepository implements the port.UserRepository interface using GORM
type GormUserRepository struct {
	BaseRepository
}

// NewGormUserRepository creates a new GormUserRepository
func NewGormUserRepository(db *gorm.DB, logger *zerolog.Logger) port.UserRepository {
	return &GormUserRepository{
		BaseRepository: NewBaseRepository(db, logger),
	}
}

// Save saves a user
func (r *GormUserRepository) Save(ctx context.Context, user *model.User) error {
	// Convert domain model to entity
	userEntity := &entity.UserEntity{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	// Save entity
	if err := r.GetDB(ctx).Save(userEntity).Error; err != nil {
		r.logger.Error().Err(err).Str("userID", user.ID).Msg("Failed to save user")
		return fmt.Errorf("failed to save user: %w", err)
	}

	return nil
}

// GetByID gets a user by ID
func (r *GormUserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	// Get entity
	var userEntity entity.UserEntity
	if err := r.GetDB(ctx).Where("id = ?", id).First(&userEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found with ID: %s", id)
		}
		r.logger.Error().Err(err).Str("userID", id).Msg("Failed to get user")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Convert entity to domain model
	user := &model.User{
		ID:        userEntity.ID,
		Email:     userEntity.Email,
		Name:      userEntity.Name,
		CreatedAt: userEntity.CreatedAt,
		UpdatedAt: userEntity.UpdatedAt,
	}

	return user, nil
}

// GetByEmail gets a user by email
func (r *GormUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	// Get entity
	var userEntity entity.UserEntity
	if err := r.GetDB(ctx).Where("email = ?", email).First(&userEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found with email: %s", email)
		}
		r.logger.Error().Err(err).Str("email", email).Msg("Failed to get user by email")
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	// Convert entity to domain model
	user := &model.User{
		ID:        userEntity.ID,
		Email:     userEntity.Email,
		Name:      userEntity.Name,
		CreatedAt: userEntity.CreatedAt,
		UpdatedAt: userEntity.UpdatedAt,
	}

	return user, nil
}

// Delete deletes a user by ID
func (r *GormUserRepository) Delete(ctx context.Context, id string) error {
	// Delete entity
	if err := r.GetDB(ctx).Where("id = ?", id).Delete(&entity.UserEntity{}).Error; err != nil {
		r.logger.Error().Err(err).Str("userID", id).Msg("Failed to delete user")
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// List lists all users
func (r *GormUserRepository) List(ctx context.Context) ([]*model.User, error) {
	// Get entities
	var userEntities []entity.UserEntity
	if err := r.GetDB(ctx).Find(&userEntities).Error; err != nil {
		r.logger.Error().Err(err).Msg("Failed to list users")
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	// Convert entities to domain models
	users := make([]*model.User, len(userEntities))
	for i, userEntity := range userEntities {
		users[i] = &model.User{
			ID:        userEntity.ID,
			Email:     userEntity.Email,
			Name:      userEntity.Name,
			CreatedAt: userEntity.CreatedAt,
			UpdatedAt: userEntity.UpdatedAt,
		}
	}

	return users, nil
}

// Ensure GormUserRepository implements port.UserRepository
var _ port.UserRepository = (*GormUserRepository)(nil)
