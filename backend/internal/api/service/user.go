package service

import (
	"context"
	"errors"
	"time"

	"go-crypto-bot-clean/backend/internal/api/repository"
)

// UserService provides user management functionality for the API
type UserService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// UserProfile represents a user profile
type UserProfile struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	FirstName string    `json:"firstName,omitempty"`
	LastName  string    `json:"lastName,omitempty"`
	Roles     []string  `json:"roles"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// UserSettings represents user settings
type UserSettings struct {
	ID                   string    `json:"id"`
	Theme                string    `json:"theme"`
	Language             string    `json:"language"`
	TimeZone             string    `json:"timeZone"`
	NotificationsEnabled bool      `json:"notificationsEnabled"`
	EmailNotifications   bool      `json:"emailNotifications"`
	PushNotifications    bool      `json:"pushNotifications"`
	DefaultCurrency      string    `json:"defaultCurrency"`
	UpdatedAt            time.Time `json:"updatedAt"`
}

// UpdateProfileRequest represents a request to update a user profile
type UpdateProfileRequest struct {
	Username  string `json:"username,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
}

// UpdateSettingsRequest represents a request to update user settings
type UpdateSettingsRequest struct {
	Theme                string `json:"theme,omitempty"`
	Language             string `json:"language,omitempty"`
	TimeZone             string `json:"timeZone,omitempty"`
	NotificationsEnabled string `json:"notificationsEnabled,omitempty"`
	EmailNotifications   string `json:"emailNotifications,omitempty"`
	PushNotifications    string `json:"pushNotifications,omitempty"`
	DefaultCurrency      string `json:"defaultCurrency,omitempty"`
}

// ChangePasswordRequest represents a request to change a user's password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

// ChangePasswordResponse represents the response from changing a password
type ChangePasswordResponse struct {
	Success   bool      `json:"success"`
	Timestamp time.Time `json:"timestamp"`
}

// GetUserProfile gets the profile of the authenticated user
func (s *UserService) GetUserProfile(ctx context.Context, userID string) (*UserProfile, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	// Get user from repository
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get user roles
	roles, err := s.userRepo.GetRoles(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Convert to UserProfile
	return &UserProfile{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Roles:     roles,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// UpdateUserProfile updates the profile of the authenticated user
func (s *UserService) UpdateUserProfile(ctx context.Context, userID string, req *UpdateProfileRequest) (*UserProfile, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	// Get user from repository
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Update user fields
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}

	// Save user
	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	// Get user roles
	roles, err := s.userRepo.GetRoles(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Convert to UserProfile
	return &UserProfile{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Roles:     roles,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// GetUserSettings gets the settings of the authenticated user
func (s *UserService) GetUserSettings(ctx context.Context, userID string) (*UserSettings, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	// Get user settings from repository
	settings, err := s.userRepo.GetSettings(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Convert to UserSettings
	return &UserSettings{
		ID:                   userID,
		Theme:                settings.Theme,
		Language:             settings.Language,
		TimeZone:             settings.TimeZone,
		NotificationsEnabled: settings.NotificationsEnabled,
		EmailNotifications:   settings.EmailNotifications,
		PushNotifications:    settings.PushNotifications,
		DefaultCurrency:      settings.DefaultCurrency,
		UpdatedAt:            settings.UpdatedAt,
	}, nil
}

// UpdateUserSettings updates the settings of the authenticated user
func (s *UserService) UpdateUserSettings(ctx context.Context, userID string, req *UpdateSettingsRequest) (*UserSettings, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	// Get current settings
	currentSettings, err := s.userRepo.GetSettings(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Parse boolean values from strings
	notificationsEnabled := currentSettings.NotificationsEnabled
	if req.NotificationsEnabled == "true" {
		notificationsEnabled = true
	} else if req.NotificationsEnabled == "false" {
		notificationsEnabled = false
	}

	emailNotifications := currentSettings.EmailNotifications
	if req.EmailNotifications == "true" {
		emailNotifications = true
	} else if req.EmailNotifications == "false" {
		emailNotifications = false
	}

	pushNotifications := currentSettings.PushNotifications
	if req.PushNotifications == "true" {
		pushNotifications = true
	} else if req.PushNotifications == "false" {
		pushNotifications = false
	}

	// Update settings
	if req.Theme != "" {
		currentSettings.Theme = req.Theme
	}
	if req.Language != "" {
		currentSettings.Language = req.Language
	}
	if req.TimeZone != "" {
		currentSettings.TimeZone = req.TimeZone
	}
	if req.DefaultCurrency != "" {
		currentSettings.DefaultCurrency = req.DefaultCurrency
	}
	currentSettings.NotificationsEnabled = notificationsEnabled
	currentSettings.EmailNotifications = emailNotifications
	currentSettings.PushNotifications = pushNotifications

	// Save settings
	err = s.userRepo.UpdateSettings(ctx, currentSettings)
	if err != nil {
		return nil, err
	}

	// Convert to UserSettings
	return &UserSettings{
		ID:                   userID,
		Theme:                currentSettings.Theme,
		Language:             currentSettings.Language,
		TimeZone:             currentSettings.TimeZone,
		NotificationsEnabled: currentSettings.NotificationsEnabled,
		EmailNotifications:   currentSettings.EmailNotifications,
		PushNotifications:    currentSettings.PushNotifications,
		DefaultCurrency:      currentSettings.DefaultCurrency,
		UpdatedAt:            currentSettings.UpdatedAt,
	}, nil
}

// ChangePassword changes the password of the authenticated user
func (s *UserService) ChangePassword(ctx context.Context, userID string, req *ChangePasswordRequest) (*ChangePasswordResponse, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	if req.CurrentPassword == "" || req.NewPassword == "" {
		return nil, errors.New("current password and new password are required")
	}

	// Get user from repository
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// TODO: Validate current password using bcrypt
	// For now, we'll just update the password without validation

	// TODO: Hash new password using bcrypt
	// For now, we'll just store the password as is
	user.PasswordHash = req.NewPassword

	// Save user
	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return &ChangePasswordResponse{
		Success:   true,
		Timestamp: time.Now(),
	}, nil
}
