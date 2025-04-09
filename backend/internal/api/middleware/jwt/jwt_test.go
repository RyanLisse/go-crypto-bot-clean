package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJWTService(t *testing.T) {
	// Create a new JWT service
	service := NewService(
		"test-access-secret",
		"test-refresh-secret",
		time.Hour,
		time.Hour*24*7,
		"test-issuer",
	)

	// Test generating an access token
	userID := "user-123"
	email := "user@example.com"
	roles := []string{"user", "admin"}

	accessToken, expiresAt, err := service.GenerateAccessToken(userID, email, roles)
	assert.NoError(t, err, "Should not error when generating access token")
	assert.NotEmpty(t, accessToken, "Access token should not be empty")
	assert.True(t, expiresAt.After(time.Now()), "Expiry time should be in the future")

	// Test validating the access token
	claims, err := service.ValidateAccessToken(accessToken)
	assert.NoError(t, err, "Should not error when validating access token")
	assert.Equal(t, userID, claims.UserID, "User ID should match")
	assert.Equal(t, email, claims.Email, "Email should match")
	assert.Equal(t, roles, claims.Roles, "Roles should match")
	assert.Equal(t, AccessToken, claims.Type, "Token type should be access")

	// Test generating a refresh token
	refreshToken, refreshExpiresAt, err := service.GenerateRefreshToken(userID)
	assert.NoError(t, err, "Should not error when generating refresh token")
	assert.NotEmpty(t, refreshToken, "Refresh token should not be empty")
	assert.True(t, refreshExpiresAt.After(time.Now()), "Refresh expiry time should be in the future")

	// Test validating the refresh token
	refreshClaims, err := service.ValidateRefreshToken(refreshToken)
	assert.NoError(t, err, "Should not error when validating refresh token")
	assert.Equal(t, userID, refreshClaims.UserID, "User ID should match")
	assert.Equal(t, RefreshToken, refreshClaims.Type, "Token type should be refresh")

	// Test validating an access token as a refresh token (should fail)
	_, err = service.ValidateRefreshToken(accessToken)
	assert.Error(t, err, "Should error when validating access token as refresh token")

	// Test validating a refresh token as an access token (should fail)
	_, err = service.ValidateAccessToken(refreshToken)
	assert.Error(t, err, "Should error when validating refresh token as access token")

	// Test validating an invalid token
	_, err = service.ValidateAccessToken("invalid-token")
	assert.Error(t, err, "Should error when validating invalid token")
}
