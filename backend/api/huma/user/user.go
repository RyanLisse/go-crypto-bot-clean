// Package user provides the user management endpoints for the Huma API.
package user

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

// User represents a user in the system
type User struct {
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

// RegisterEndpoints registers the user management endpoints.
func RegisterEndpoints(api huma.API, basePath string) {
	// GET /user/profile
	huma.Register(api, huma.Operation{
		OperationID: "get-user-profile",
		Method:      http.MethodGet,
		Path:        basePath + "/user/profile",
		Summary:     "Get user profile",
		Description: "Returns the profile of the authenticated user",
		Tags:        []string{"User"},
	}, func(ctx context.Context, input *struct{}) (*struct {
		Body User
	}, error) {
		resp := &struct {
			Body User
		}{}

		// In a real implementation, we would get the user from the database
		// based on the authenticated user. For now, we'll just return a mock response.
		resp.Body = User{
			ID:        "user-123456",
			Email:     "user@example.com",
			Username:  "johndoe",
			FirstName: "John",
			LastName:  "Doe",
			Roles:     []string{"user"},
			CreatedAt: time.Now().AddDate(0, -1, 0),
			UpdatedAt: time.Now(),
		}

		return resp, nil
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
		Body User
	}, error) {
		resp := &struct {
			Body User
		}{}

		// In a real implementation, we would update the user in the database
		// based on the authenticated user. For now, we'll just return a mock response.
		resp.Body = User{
			ID:        "user-123456",
			Email:     "user@example.com",
			Username:  "updateduser",
			FirstName: "Updated",
			LastName:  "User",
			Roles:     []string{"user"},
			CreatedAt: time.Now().AddDate(0, -1, 0),
			UpdatedAt: time.Now(),
		}

		return resp, nil
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
		Body UserSettings
	}, error) {
		resp := &struct {
			Body UserSettings
		}{}

		// In a real implementation, we would get the user settings from the database
		// based on the authenticated user. For now, we'll just return a mock response.
		resp.Body = UserSettings{
			ID:                   "user-123456",
			Theme:                "light",
			Language:             "en",
			TimeZone:             "UTC",
			NotificationsEnabled: true,
			EmailNotifications:   true,
			PushNotifications:    true,
			DefaultCurrency:      "USD",
			UpdatedAt:            time.Now(),
		}

		return resp, nil
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
		Body UserSettings
	}, error) {
		resp := &struct {
			Body UserSettings
		}{}

		// In a real implementation, we would update the user settings in the database
		// based on the authenticated user. For now, we'll just return a mock response.
		notificationsEnabled := true
		if input.NotificationsEnabled == "false" {
			notificationsEnabled = false
		}

		emailNotifications := true
		if input.EmailNotifications == "false" {
			emailNotifications = false
		}

		pushNotifications := false
		if input.PushNotifications == "true" {
			pushNotifications = true
		}

		resp.Body = UserSettings{
			ID:                   "user-123456",
			Theme:                "dark",
			Language:             "en",
			TimeZone:             "America/New_York",
			NotificationsEnabled: notificationsEnabled,
			EmailNotifications:   emailNotifications,
			PushNotifications:    pushNotifications,
			DefaultCurrency:      "USD",
			UpdatedAt:            time.Now(),
		}

		return resp, nil
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
		resp := &struct {
			Body struct {
				Success   bool      `json:"success"`
				Timestamp time.Time `json:"timestamp"`
			}
		}{}

		// In a real implementation, we would validate the current password
		// and update the password in the database. For now, we'll just return a success response.
		resp.Body.Success = true
		resp.Body.Timestamp = time.Now()

		return resp, nil
	})
}
