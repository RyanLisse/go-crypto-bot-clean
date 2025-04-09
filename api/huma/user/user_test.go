package user

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/stretchr/testify/assert"
)

func TestUserEndpoints(t *testing.T) {
	// Create a test API
	_, api := humatest.New(t)

	// Register the user endpoints
	RegisterEndpoints(api, "/api/v1")

	// Test the get user profile endpoint
	profileResp := api.Get("/api/v1/user/profile")
	assert.Equal(t, http.StatusOK, profileResp.Code, "Should return 200 OK")

	// Decode the response
	var profileResponse struct {
		ID        string   `json:"id"`
		Email     string   `json:"email"`
		Username  string   `json:"username"`
		FirstName string   `json:"firstName"`
		LastName  string   `json:"lastName"`
		Roles     []string `json:"roles"`
	}
	err := json.Unmarshal(profileResp.Body.Bytes(), &profileResponse)
	assert.NoError(t, err, "Should not error when decoding response")
	assert.NotEmpty(t, profileResponse.ID, "Should return a user ID")
	assert.NotEmpty(t, profileResponse.Email, "Should return an email")
	assert.NotEmpty(t, profileResponse.Username, "Should return a username")

	// Test the update user profile endpoint
	updateProfileResp := api.Put("/api/v1/user/profile", map[string]interface{}{
		"username":  "updateduser",
		"firstName": "Updated",
		"lastName":  "User",
	})
	assert.Equal(t, http.StatusOK, updateProfileResp.Code, "Should return 200 OK")

	// Decode the response
	var updateProfileResponse struct {
		ID        string   `json:"id"`
		Email     string   `json:"email"`
		Username  string   `json:"username"`
		FirstName string   `json:"firstName"`
		LastName  string   `json:"lastName"`
		Roles     []string `json:"roles"`
	}
	err = json.Unmarshal(updateProfileResp.Body.Bytes(), &updateProfileResponse)
	assert.NoError(t, err, "Should not error when decoding response")
	assert.Equal(t, "updateduser", updateProfileResponse.Username, "Should update the username")
	assert.Equal(t, "Updated", updateProfileResponse.FirstName, "Should update the first name")
	assert.Equal(t, "User", updateProfileResponse.LastName, "Should update the last name")

	// Test the get user settings endpoint
	settingsResp := api.Get("/api/v1/user/settings")
	assert.Equal(t, http.StatusOK, settingsResp.Code, "Should return 200 OK")

	// Decode the response
	var settingsResponse struct {
		Theme                string `json:"theme"`
		Language             string `json:"language"`
		TimeZone             string `json:"timeZone"`
		NotificationsEnabled bool   `json:"notificationsEnabled"`
		EmailNotifications   bool   `json:"emailNotifications"`
		PushNotifications    bool   `json:"pushNotifications"`
		DefaultCurrency      string `json:"defaultCurrency"`
	}
	err = json.Unmarshal(settingsResp.Body.Bytes(), &settingsResponse)
	assert.NoError(t, err, "Should not error when decoding response")
	assert.NotEmpty(t, settingsResponse.Theme, "Should return a theme")
	assert.NotEmpty(t, settingsResponse.Language, "Should return a language")
	assert.NotEmpty(t, settingsResponse.TimeZone, "Should return a time zone")
	assert.NotEmpty(t, settingsResponse.DefaultCurrency, "Should return a default currency")

	// Test the update user settings endpoint
	updateSettingsResp := api.Put("/api/v1/user/settings", map[string]interface{}{
		"theme":                "dark",
		"language":             "en",
		"timeZone":             "America/New_York",
		"notificationsEnabled": "true",
		"emailNotifications":   "true",
		"pushNotifications":    "false",
		"defaultCurrency":      "USD",
	})
	assert.Equal(t, http.StatusOK, updateSettingsResp.Code, "Should return 200 OK")

	// Decode the response
	var updateSettingsResponse struct {
		Theme                string `json:"theme"`
		Language             string `json:"language"`
		TimeZone             string `json:"timeZone"`
		NotificationsEnabled bool   `json:"notificationsEnabled"`
		EmailNotifications   bool   `json:"emailNotifications"`
		PushNotifications    bool   `json:"pushNotifications"`
		DefaultCurrency      string `json:"defaultCurrency"`
	}
	err = json.Unmarshal(updateSettingsResp.Body.Bytes(), &updateSettingsResponse)
	assert.NoError(t, err, "Should not error when decoding response")
	assert.Equal(t, "dark", updateSettingsResponse.Theme, "Should update the theme")
	assert.Equal(t, "en", updateSettingsResponse.Language, "Should update the language")
	assert.Equal(t, "America/New_York", updateSettingsResponse.TimeZone, "Should update the time zone")
	assert.Equal(t, "USD", updateSettingsResponse.DefaultCurrency, "Should update the default currency")
	assert.True(t, updateSettingsResponse.NotificationsEnabled, "Should update notifications enabled")
	assert.True(t, updateSettingsResponse.EmailNotifications, "Should update email notifications")
	assert.False(t, updateSettingsResponse.PushNotifications, "Should update push notifications")

	// Test the change password endpoint
	changePasswordResp := api.Post("/api/v1/user/password", map[string]interface{}{
		"currentPassword": "password123",
		"newPassword":     "newpassword123",
	})
	assert.Equal(t, http.StatusOK, changePasswordResp.Code, "Should return 200 OK")

	// Decode the response
	var changePasswordResponse struct {
		Success bool `json:"success"`
	}
	err = json.Unmarshal(changePasswordResp.Body.Bytes(), &changePasswordResponse)
	assert.NoError(t, err, "Should not error when decoding response")
	assert.True(t, changePasswordResponse.Success, "Should return success: true")
}
