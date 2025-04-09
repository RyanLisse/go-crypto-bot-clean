// Package user provides the user management endpoints for the Huma API.
package user

import (
	"context"
	"net/http"
	"time"

	"go-crypto-bot-clean/backend/internal/api/service"

	"github.com/danielgtaylor/huma/v2"
)

// RegisterUserEndpoints registers the user management endpoints with service implementation.
func RegisterUserEndpoints(api huma.API, basePath string, userService *service.UserService) {
	// GET /user/profile
	huma.Register(api, huma.Operation{
		OperationID: "get-user-profile",
		Method:      http.MethodGet,
		Path:        basePath + "/user/profile",
		Summary:     "Get user profile",
		Description: "Returns the profile of the authenticated user",
		Tags:        []string{"User"},
	}, func(ctx context.Context, input *struct{}) (*struct {
		Body service.UserProfile
	}, error) {
		// Get user ID from context
		// In a real implementation, we would get the user ID from the authenticated user
		userID := "user-123456"

		// Get user profile
		profile, err := userService.GetUserProfile(ctx, userID)
		if err != nil {
			return nil, err
		}

		// Return result
		return &struct {
			Body service.UserProfile
		}{
			Body: *profile,
		}, nil
	})

	// PUT /user/profile
	huma.Register(api, huma.Operation{
		OperationID: "update-user-profile",
		Method:      http.MethodPut,
		Path:        basePath + "/user/profile",
		Summary:     "Update user profile",
		Description: "Updates the profile of the authenticated user",
		Tags:        []string{"User"},
	}, func(ctx context.Context, input *struct {
		Username  string `json:"username,omitempty"`
		FirstName string `json:"firstName,omitempty"`
		LastName  string `json:"lastName,omitempty"`
	}) (*struct {
		Body service.UserProfile
	}, error) {
		// Get user ID from context
		// In a real implementation, we would get the user ID from the authenticated user
		userID := "user-123456"

		// Convert API request to service request
		req := &service.UpdateProfileRequest{
			Username:  input.Username,
			FirstName: input.FirstName,
			LastName:  input.LastName,
		}

		// Update user profile
		profile, err := userService.UpdateUserProfile(ctx, userID, req)
		if err != nil {
			return nil, err
		}

		// Return result
		return &struct {
			Body service.UserProfile
		}{
			Body: *profile,
		}, nil
	})

	// GET /user/settings
	huma.Register(api, huma.Operation{
		OperationID: "get-user-settings",
		Method:      http.MethodGet,
		Path:        basePath + "/user/settings",
		Summary:     "Get user settings",
		Description: "Returns the settings of the authenticated user",
		Tags:        []string{"User"},
	}, func(ctx context.Context, input *struct{}) (*struct {
		Body service.UserSettings
	}, error) {
		// Get user ID from context
		// In a real implementation, we would get the user ID from the authenticated user
		userID := "user-123456"

		// Get user settings
		settings, err := userService.GetUserSettings(ctx, userID)
		if err != nil {
			return nil, err
		}

		// Return result
		return &struct {
			Body service.UserSettings
		}{
			Body: *settings,
		}, nil
	})

	// PUT /user/settings
	huma.Register(api, huma.Operation{
		OperationID: "update-user-settings",
		Method:      http.MethodPut,
		Path:        basePath + "/user/settings",
		Summary:     "Update user settings",
		Description: "Updates the settings of the authenticated user",
		Tags:        []string{"User"},
	}, func(ctx context.Context, input *struct {
		Theme                string `json:"theme,omitempty"`
		Language             string `json:"language,omitempty"`
		TimeZone             string `json:"timeZone,omitempty"`
		NotificationsEnabled string `json:"notificationsEnabled,omitempty"`
		EmailNotifications   string `json:"emailNotifications,omitempty"`
		PushNotifications    string `json:"pushNotifications,omitempty"`
		DefaultCurrency      string `json:"defaultCurrency,omitempty"`
	}) (*struct {
		Body service.UserSettings
	}, error) {
		// Get user ID from context
		// In a real implementation, we would get the user ID from the authenticated user
		userID := "user-123456"

		// Convert API request to service request
		req := &service.UpdateSettingsRequest{
			Theme:                input.Theme,
			Language:             input.Language,
			TimeZone:             input.TimeZone,
			NotificationsEnabled: input.NotificationsEnabled,
			EmailNotifications:   input.EmailNotifications,
			PushNotifications:    input.PushNotifications,
			DefaultCurrency:      input.DefaultCurrency,
		}

		// Update user settings
		settings, err := userService.UpdateUserSettings(ctx, userID, req)
		if err != nil {
			return nil, err
		}

		// Return result
		return &struct {
			Body service.UserSettings
		}{
			Body: *settings,
		}, nil
	})

	// POST /user/password
	huma.Register(api, huma.Operation{
		OperationID: "change-password",
		Method:      http.MethodPost,
		Path:        basePath + "/user/password",
		Summary:     "Change password",
		Description: "Changes the password of the authenticated user",
		Tags:        []string{"User"},
	}, func(ctx context.Context, input *struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}) (*struct {
		Body struct {
			Success   bool      `json:"success"`
			Timestamp time.Time `json:"timestamp"`
		}
	}, error) {
		// Get user ID from context
		// In a real implementation, we would get the user ID from the authenticated user
		userID := "user-123456"

		// Convert API request to service request
		req := &service.ChangePasswordRequest{
			CurrentPassword: input.CurrentPassword,
			NewPassword:     input.NewPassword,
		}

		// Change password
		result, err := userService.ChangePassword(ctx, userID, req)
		if err != nil {
			return nil, err
		}

		// Return result
		resp := &struct {
			Body struct {
				Success   bool      `json:"success"`
				Timestamp time.Time `json:"timestamp"`
			}
		}{}
		resp.Body.Success = result.Success
		resp.Body.Timestamp = result.Timestamp

		return resp, nil
	})
}
