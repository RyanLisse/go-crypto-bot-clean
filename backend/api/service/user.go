package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go-crypto-bot-clean/backend/api/repository"

	"golang.org/x/crypto/bcrypt"
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

// ChangePasswordResponse represents the response from changing a password with enhanced error handling
type ChangePasswordResponse struct {
	Success   bool      `json:"success"`
	Timestamp time.Time `json:"timestamp"`
	Error     string    `json:"error,omitempty"` // Add error field for detailed messages
}

// GetUserProfile gets the profile of the authenticated user with validation
func (s *UserService) GetUserProfile(ctx context.Context, userID string) (*UserProfile, error) {
	if userID == "" {
		return nil, errors.New("user ID is required") // HTTP 400: Bad Request
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err) // HTTP 404: Not Found
	}

	roles, err := s.userRepo.GetRoles(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles: %w", err) // HTTP 500: Internal Server Error
	}

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

// UpdateUserProfile updates the profile of the authenticated user with validation
func (s *UserService) UpdateUserProfile(ctx context.Context, userID string, req *UpdateProfileRequest) (*UserProfile, error) {
	if userID == "" {
		return nil, errors.New("user ID is required") // HTTP 400
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err) // HTTP 404
	}

	if req.Username != "" && len(req.Username) < 3 {
		return nil, errors.New("username must be at least 3 characters") // HTTP 400
	}
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err) // HTTP 500
	}

	roles, err := s.userRepo.GetRoles(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles: %w", err) // HTTP 500
	}

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

// GetUserSettings gets the settings of the authenticated user with validation
func (s *UserService) GetUserSettings(ctx context.Context, userID string) (*UserSettings, error) {
	if userID == "" {
		return nil, errors.New("user ID is required") // HTTP 400
	}

	settings, err := s.userRepo.GetSettings(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("settings not found: %w", err) // HTTP 404
	}

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

// UpdateUserSettings updates the settings of the authenticated user with validation
func (s *UserService) UpdateUserSettings(ctx context.Context, userID string, req *UpdateSettingsRequest) (*UserSettings, error) {
	if userID == "" {
		return nil, errors.New("user ID is required") // HTTP 400
	}

	currentSettings, err := s.userRepo.GetSettings(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("settings not found: %w", err) // HTTP 404
	}

	if req.Theme != "" && (req.Theme != "light" && req.Theme != "dark") {
		return nil, errors.New("invalid theme value") // HTTP 400
	}
	if req.Language != "" && len(req.Language) != 2 {
		return nil, errors.New("language must be a 2-letter code") // HTTP 400
	}

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

	currentSettings.Theme = req.Theme
	currentSettings.Language = req.Language
	currentSettings.TimeZone = req.TimeZone
	currentSettings.DefaultCurrency = req.DefaultCurrency
	currentSettings.NotificationsEnabled = notificationsEnabled
	currentSettings.EmailNotifications = emailNotifications
	currentSettings.PushNotifications = pushNotifications

	err = s.userRepo.UpdateSettings(ctx, currentSettings)
	if err != nil {
		return nil, fmt.Errorf("failed to update settings: %w", err) // HTTP 500
	}

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

// ChangePassword changes the password of the authenticated user with proper validation and hashing
func (s *UserService) ChangePassword(ctx context.Context, userID string, req *ChangePasswordRequest) (*ChangePasswordResponse, error) {
	if userID == "" {
		return &ChangePasswordResponse{Success: false, Error: "user ID is required", Timestamp: time.Now()}, nil
	}

	if req.CurrentPassword == "" || req.NewPassword == "" {
		return &ChangePasswordResponse{Success: false, Error: "current password and new password are required", Timestamp: time.Now()}, nil
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return &ChangePasswordResponse{Success: false, Error: fmt.Sprintf("user not found: %v", err), Timestamp: time.Now()}, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword))
	if err != nil {
		return &ChangePasswordResponse{Success: false, Error: "invalid current password", Timestamp: time.Now()}, nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return &ChangePasswordResponse{Success: false, Error: fmt.Sprintf("failed to hash password: %v", err), Timestamp: time.Now()}, nil
	}
	user.PasswordHash = string(hashedPassword)

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return &ChangePasswordResponse{Success: false, Error: fmt.Sprintf("failed to update user: %v", err), Timestamp: time.Now()}, nil
	}

	return &ChangePasswordResponse{Success: true, Timestamp: time.Now()}, nil
}
