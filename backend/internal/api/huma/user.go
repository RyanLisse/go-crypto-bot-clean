package huma

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

// UserProfileResponse represents a user profile
type UserProfileResponse struct {
	Body struct {
		ID        string    `json:"id" doc:"User ID" example:"user-123456"`
		Email     string    `json:"email" doc:"Email address" example:"user@example.com"`
		Username  string    `json:"username" doc:"Username" example:"johndoe"`
		FirstName string    `json:"firstName,omitempty" doc:"First name" example:"John"`
		LastName  string    `json:"lastName,omitempty" doc:"Last name" example:"Doe"`
		Roles     []string  `json:"roles" doc:"User roles" example:"[\"user\", \"admin\"]"`
		CreatedAt time.Time `json:"createdAt" doc:"Time when the user was created" example:"2023-01-01T00:00:00Z"`
		UpdatedAt time.Time `json:"updatedAt" doc:"Time when the user was last updated" example:"2023-01-02T00:00:00Z"`
	}
}

// UpdateProfileRequest represents a request to update a user profile
type UpdateProfileRequest struct {
	Body struct {
		Username  string `json:"username,omitempty" doc:"Username" example:"johndoe"`
		FirstName string `json:"firstName,omitempty" doc:"First name" example:"John"`
		LastName  string `json:"lastName,omitempty" doc:"Last name" example:"Doe"`
	}
}

// UserSettingsResponse represents user settings
type UserSettingsResponse struct {
	Body struct {
		ID                 string    `json:"id" doc:"User ID" example:"user-123456"`
		Theme              string    `json:"theme" doc:"UI theme" example:"dark" enum:"light,dark,system"`
		Language           string    `json:"language" doc:"UI language" example:"en" enum:"en,es,fr,de,zh"`
		TimeZone           string    `json:"timeZone" doc:"Time zone" example:"America/New_York"`
		NotificationsEnabled bool      `json:"notificationsEnabled" doc:"Whether notifications are enabled" example:"true"`
		EmailNotifications  bool      `json:"emailNotifications" doc:"Whether email notifications are enabled" example:"true"`
		PushNotifications   bool      `json:"pushNotifications" doc:"Whether push notifications are enabled" example:"true"`
		DefaultCurrency     string    `json:"defaultCurrency" doc:"Default currency" example:"USD" enum:"USD,EUR,GBP,JPY,BTC,ETH"`
		UpdatedAt           time.Time `json:"updatedAt" doc:"Time when the settings were last updated" example:"2023-01-02T00:00:00Z"`
	}
}

// UpdateSettingsRequest represents a request to update user settings
type UpdateSettingsRequest struct {
	Body struct {
		Theme               string `json:"theme,omitempty" doc:"UI theme" example:"dark" enum:"light,dark,system"`
		Language            string `json:"language,omitempty" doc:"UI language" example:"en" enum:"en,es,fr,de,zh"`
		TimeZone            string `json:"timeZone,omitempty" doc:"Time zone" example:"America/New_York"`
		NotificationsEnabled *bool  `json:"notificationsEnabled,omitempty" doc:"Whether notifications are enabled" example:"true"`
		EmailNotifications   *bool  `json:"emailNotifications,omitempty" doc:"Whether email notifications are enabled" example:"true"`
		PushNotifications    *bool  `json:"pushNotifications,omitempty" doc:"Whether push notifications are enabled" example:"true"`
		DefaultCurrency      string `json:"defaultCurrency,omitempty" doc:"Default currency" example:"USD" enum:"USD,EUR,GBP,JPY,BTC,ETH"`
	}
}

// ChangePasswordRequest represents a request to change a user's password
type ChangePasswordRequest struct {
	Body struct {
		CurrentPassword string `json:"currentPassword" doc:"Current password" example:"password123" binding:"required"`
		NewPassword     string `json:"newPassword" doc:"New password" example:"newpassword123" binding:"required"`
	}
}

// registerUserEndpoints registers the user endpoints.
func registerUserEndpoints(api huma.API, basePath string) {
	// GET /user/profile
	huma.Register(api, huma.Operation{
		OperationID: "get-user-profile",
		Method:      http.MethodGet,
		Path:        basePath + "/user/profile",
		Summary:     "Get user profile",
		Description: "Returns the profile of the authenticated user",
		Tags:        []string{"User"},
	}, func(ctx context.Context, input *struct{}) (*UserProfileResponse, error) {
		// This is just a placeholder for documentation purposes
		return nil, nil
	})

	// PUT /user/profile
	huma.Register(api, huma.Operation{
		OperationID: "update-user-profile",
		Method:      http.MethodPut,
		Path:        basePath + "/user/profile",
		Summary:     "Update user profile",
		Description: "Updates the profile of the authenticated user",
		Tags:        []string{"User"},
	}, func(ctx context.Context, input *UpdateProfileRequest) (*UserProfileResponse, error) {
		// This is just a placeholder for documentation purposes
		return nil, nil
	})

	// GET /user/settings
	huma.Register(api, huma.Operation{
		OperationID: "get-user-settings",
		Method:      http.MethodGet,
		Path:        basePath + "/user/settings",
		Summary:     "Get user settings",
		Description: "Returns the settings of the authenticated user",
		Tags:        []string{"User"},
	}, func(ctx context.Context, input *struct{}) (*UserSettingsResponse, error) {
		// This is just a placeholder for documentation purposes
		return nil, nil
	})

	// PUT /user/settings
	huma.Register(api, huma.Operation{
		OperationID: "update-user-settings",
		Method:      http.MethodPut,
		Path:        basePath + "/user/settings",
		Summary:     "Update user settings",
		Description: "Updates the settings of the authenticated user",
		Tags:        []string{"User"},
	}, func(ctx context.Context, input *UpdateSettingsRequest) (*UserSettingsResponse, error) {
		// This is just a placeholder for documentation purposes
		return nil, nil
	})

	// POST /user/password
	huma.Register(api, huma.Operation{
		OperationID: "change-password",
		Method:      http.MethodPost,
		Path:        basePath + "/user/password",
		Summary:     "Change password",
		Description: "Changes the password of the authenticated user",
		Tags:        []string{"User"},
	}, func(ctx context.Context, input *ChangePasswordRequest) (*struct {
		Body struct {
			Success   bool      `json:"success" doc:"Whether the password change was successful" example:"true"`
			Timestamp time.Time `json:"timestamp" doc:"Timestamp of the password change" example:"2023-02-02T10:00:00Z"`
		}
	}, error) {
		// This is just a placeholder for documentation purposes
		return nil, nil
	})
}
