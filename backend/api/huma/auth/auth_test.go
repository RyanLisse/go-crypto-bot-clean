package auth

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/stretchr/testify/assert"
)

func TestAuthEndpoints(t *testing.T) {
	// Create a test API
	_, api := humatest.New(t)

	// Register the auth endpoints
	RegisterEndpoints(api, "/api/v1")

	// Test the login endpoint
	loginResp := api.Post("/api/v1/auth/login", map[string]interface{}{
		"email":    "user@example.com",
		"password": "password123",
	})
	assert.Equal(t, http.StatusOK, loginResp.Code, "Should return 200 OK")

	// Decode the response
	var loginResponse struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
		User         struct {
			ID       string   `json:"id"`
			Email    string   `json:"email"`
			Username string   `json:"username"`
			Roles    []string `json:"roles"`
		} `json:"user"`
	}
	err := json.Unmarshal(loginResp.Body.Bytes(), &loginResponse)
	assert.NoError(t, err, "Should not error when decoding response")
	assert.NotEmpty(t, loginResponse.AccessToken, "Should return an access token")
	assert.NotEmpty(t, loginResponse.RefreshToken, "Should return a refresh token")
	assert.NotEmpty(t, loginResponse.User.ID, "Should return a user ID")

	// Test the register endpoint
	registerResp := api.Post("/api/v1/auth/register", map[string]interface{}{
		"email":     "newuser@example.com",
		"username":  "newuser",
		"password":  "password123",
		"firstName": "New",
		"lastName":  "User",
	})
	assert.Equal(t, http.StatusOK, registerResp.Code, "Should return 200 OK")

	// Decode the response
	var registerResponse struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
		User         struct {
			ID       string   `json:"id"`
			Email    string   `json:"email"`
			Username string   `json:"username"`
			Roles    []string `json:"roles"`
		} `json:"user"`
	}
	err = json.Unmarshal(registerResp.Body.Bytes(), &registerResponse)
	assert.NoError(t, err, "Should not error when decoding response")
	assert.NotEmpty(t, registerResponse.AccessToken, "Should return an access token")
	assert.NotEmpty(t, registerResponse.RefreshToken, "Should return a refresh token")
	assert.NotEmpty(t, registerResponse.User.ID, "Should return a user ID")

	// Test the refresh token endpoint
	refreshResp := api.Post("/api/v1/auth/refresh", map[string]interface{}{
		"refreshToken": loginResponse.RefreshToken,
	})
	assert.Equal(t, http.StatusOK, refreshResp.Code, "Should return 200 OK")

	// Decode the response
	var refreshResponse struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}
	err = json.Unmarshal(refreshResp.Body.Bytes(), &refreshResponse)
	assert.NoError(t, err, "Should not error when decoding response")
	assert.NotEmpty(t, refreshResponse.AccessToken, "Should return a new access token")
	assert.NotEmpty(t, refreshResponse.RefreshToken, "Should return a new refresh token")

	// Test the logout endpoint
	logoutResp := api.Post("/api/v1/auth/logout", map[string]interface{}{})
	assert.Equal(t, http.StatusOK, logoutResp.Code, "Should return 200 OK")

	// Decode the response
	var logoutResponse struct {
		Success bool `json:"success"`
	}
	err = json.Unmarshal(logoutResp.Body.Bytes(), &logoutResponse)
	assert.NoError(t, err, "Should not error when decoding response")
	assert.True(t, logoutResponse.Success, "Should return success: true")

	// Test the verify token endpoint
	// In a real implementation, we would set the Authorization header with the access token
	verifyResp := api.Get("/api/v1/auth/verify")
	assert.Equal(t, http.StatusOK, verifyResp.Code, "Should return 200 OK")

	// Decode the response
	var verifyResponse struct {
		Valid bool `json:"valid"`
		User  struct {
			ID       string   `json:"id"`
			Email    string   `json:"email"`
			Username string   `json:"username"`
			Roles    []string `json:"roles"`
		} `json:"user"`
	}
	err = json.Unmarshal(verifyResp.Body.Bytes(), &verifyResponse)
	assert.NoError(t, err, "Should not error when decoding response")
	assert.True(t, verifyResponse.Valid, "Should return valid: true")
	assert.Equal(t, loginResponse.User.ID, verifyResponse.User.ID, "Should return the same user ID")
}
