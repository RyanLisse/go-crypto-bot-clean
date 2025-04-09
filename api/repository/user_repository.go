// Package repository provides database repositories for the API
package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go-crypto-bot-clean/api/models"
	"gorm.io/gorm"
)

// Common errors
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// UserRepository defines the interface for user operations
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *models.User) error
	
	// GetByID gets a user by ID
	GetByID(ctx context.Context, id string) (*models.User, error)
	
	// GetByEmail gets a user by email
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	
	// GetByUsername gets a user by username
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	
	// Update updates a user
	Update(ctx context.Context, user *models.User) error
	
	// Delete deletes a user
	Delete(ctx context.Context, id string) error
	
	// UpdateLastLogin updates the last login time for a user
	UpdateLastLogin(ctx context.Context, id string) error
	
	// AddRole adds a role to a user
	AddRole(ctx context.Context, userID string, role string) error
	
	// RemoveRole removes a role from a user
	RemoveRole(ctx context.Context, userID string, role string) error
	
	// GetRoles gets all roles for a user
	GetRoles(ctx context.Context, userID string) ([]string, error)
	
	// GetSettings gets the settings for a user
	GetSettings(ctx context.Context, userID string) (*models.UserSettings, error)
	
	// UpdateSettings updates the settings for a user
	UpdateSettings(ctx context.Context, settings *models.UserSettings) error
	
	// SaveRefreshToken saves a refresh token for a user
	SaveRefreshToken(ctx context.Context, token *models.RefreshToken) error
	
	// GetRefreshToken gets a refresh token by token string
	GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error)
	
	// RevokeRefreshToken revokes a refresh token
	RevokeRefreshToken(ctx context.Context, token string) error
	
	// RevokeAllRefreshTokens revokes all refresh tokens for a user
	RevokeAllRefreshTokens(ctx context.Context, userID string) error
}

// GormUserRepository implements UserRepository using GORM
type GormUserRepository struct {
	db *gorm.DB
}

// NewGormUserRepository creates a new GORM user repository
func NewGormUserRepository(db *gorm.DB) UserRepository {
	return &GormUserRepository{
		db: db,
	}
}

// Create creates a new user
func (r *GormUserRepository) Create(ctx context.Context, user *models.User) error {
	// Check if email already exists
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("email = ?", user.Email).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check email: %w", err)
	}
	if count > 0 {
		return ErrEmailAlreadyExists
	}

	// Check if username already exists
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("username = ?", user.Username).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check username: %w", err)
	}
	if count > 0 {
		return ErrUsernameAlreadyExists
	}

	// Create user
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByID gets a user by ID
func (r *GormUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Preload("Roles").Preload("Settings").First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return &user, nil
}

// GetByEmail gets a user by email
func (r *GormUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Preload("Roles").Preload("Settings").First(&user, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}

// GetByUsername gets a user by username
func (r *GormUserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Preload("Roles").Preload("Settings").First(&user, "username = ?", username).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	return &user, nil
}

// Update updates a user
func (r *GormUserRepository) Update(ctx context.Context, user *models.User) error {
	// Check if user exists
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", user.ID).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if count == 0 {
		return ErrUserNotFound
	}

	// Update user
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// Delete deletes a user
func (r *GormUserRepository) Delete(ctx context.Context, id string) error {
	// Check if user exists
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if count == 0 {
		return ErrUserNotFound
	}

	// Delete user (soft delete)
	if err := r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// UpdateLastLogin updates the last login time for a user
func (r *GormUserRepository) UpdateLastLogin(ctx context.Context, id string) error {
	now := time.Now()
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Update("last_login_at", now).Error; err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}
	return nil
}

// AddRole adds a role to a user
func (r *GormUserRepository) AddRole(ctx context.Context, userID string, role string) error {
	// Check if user exists
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if count == 0 {
		return ErrUserNotFound
	}

	// Check if role already exists
	if err := r.db.WithContext(ctx).Model(&models.UserRole{}).Where("user_id = ? AND role = ?", userID, role).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check role existence: %w", err)
	}
	if count > 0 {
		// Role already exists, nothing to do
		return nil
	}

	// Add role
	userRole := models.UserRole{
		UserID: userID,
		Role:   role,
	}
	if err := r.db.WithContext(ctx).Create(&userRole).Error; err != nil {
		return fmt.Errorf("failed to add role: %w", err)
	}

	return nil
}

// RemoveRole removes a role from a user
func (r *GormUserRepository) RemoveRole(ctx context.Context, userID string, role string) error {
	// Check if user exists
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if count == 0 {
		return ErrUserNotFound
	}

	// Remove role
	if err := r.db.WithContext(ctx).Where("user_id = ? AND role = ?", userID, role).Delete(&models.UserRole{}).Error; err != nil {
		return fmt.Errorf("failed to remove role: %w", err)
	}

	return nil
}

// GetRoles gets all roles for a user
func (r *GormUserRepository) GetRoles(ctx context.Context, userID string) ([]string, error) {
	// Check if user exists
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Count(&count).Error; err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if count == 0 {
		return nil, ErrUserNotFound
	}

	// Get roles
	var userRoles []models.UserRole
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&userRoles).Error; err != nil {
		return nil, fmt.Errorf("failed to get roles: %w", err)
	}

	// Extract role names
	roles := make([]string, len(userRoles))
	for i, userRole := range userRoles {
		roles[i] = userRole.Role
	}

	return roles, nil
}

// GetSettings gets the settings for a user
func (r *GormUserRepository) GetSettings(ctx context.Context, userID string) (*models.UserSettings, error) {
	// Check if user exists
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Count(&count).Error; err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if count == 0 {
		return nil, ErrUserNotFound
	}

	// Get settings
	var settings models.UserSettings
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&settings).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create default settings
			settings = models.UserSettings{
				UserID:              userID,
				Theme:               "light",
				Language:            "en",
				TimeZone:            "UTC",
				NotificationsEnabled: true,
				EmailNotifications:   true,
				PushNotifications:    false,
				DefaultCurrency:      "USD",
			}
			if err := r.db.WithContext(ctx).Create(&settings).Error; err != nil {
				return nil, fmt.Errorf("failed to create default settings: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to get settings: %w", err)
		}
	}

	return &settings, nil
}

// UpdateSettings updates the settings for a user
func (r *GormUserRepository) UpdateSettings(ctx context.Context, settings *models.UserSettings) error {
	// Check if user exists
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", settings.UserID).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if count == 0 {
		return ErrUserNotFound
	}

	// Check if settings exist
	if err := r.db.WithContext(ctx).Model(&models.UserSettings{}).Where("user_id = ?", settings.UserID).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check settings existence: %w", err)
	}

	if count == 0 {
		// Create settings
		if err := r.db.WithContext(ctx).Create(settings).Error; err != nil {
			return fmt.Errorf("failed to create settings: %w", err)
		}
	} else {
		// Update settings
		if err := r.db.WithContext(ctx).Save(settings).Error; err != nil {
			return fmt.Errorf("failed to update settings: %w", err)
		}
	}

	return nil
}

// SaveRefreshToken saves a refresh token for a user
func (r *GormUserRepository) SaveRefreshToken(ctx context.Context, token *models.RefreshToken) error {
	// Check if user exists
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", token.UserID).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if count == 0 {
		return ErrUserNotFound
	}

	// Save token
	if err := r.db.WithContext(ctx).Create(token).Error; err != nil {
		return fmt.Errorf("failed to save refresh token: %w", err)
	}

	return nil
}

// GetRefreshToken gets a refresh token by token string
func (r *GormUserRepository) GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	if err := r.db.WithContext(ctx).Where("token = ? AND revoked_at IS NULL", token).First(&refreshToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	// Check if token is expired
	if refreshToken.ExpiresAt.Before(time.Now()) {
		return nil, ErrInvalidCredentials
	}

	return &refreshToken, nil
}

// RevokeRefreshToken revokes a refresh token
func (r *GormUserRepository) RevokeRefreshToken(ctx context.Context, token string) error {
	now := time.Now()
	if err := r.db.WithContext(ctx).Model(&models.RefreshToken{}).Where("token = ?", token).Update("revoked_at", now).Error; err != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}
	return nil
}

// RevokeAllRefreshTokens revokes all refresh tokens for a user
func (r *GormUserRepository) RevokeAllRefreshTokens(ctx context.Context, userID string) error {
	now := time.Now()
	if err := r.db.WithContext(ctx).Model(&models.RefreshToken{}).Where("user_id = ? AND revoked_at IS NULL", userID).Update("revoked_at", now).Error; err != nil {
		return fmt.Errorf("failed to revoke all refresh tokens: %w", err)
	}
	return nil
}
