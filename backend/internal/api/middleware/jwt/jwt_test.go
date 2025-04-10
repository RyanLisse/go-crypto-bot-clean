package jwt

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestJWTService_GenerateAccessToken(t *testing.T) {
	// Create JWT service with test configuration
	service := NewService(
		"test-access-secret",
		"test-refresh-secret",
		time.Hour,
		time.Hour*24*7,
		"test-issuer",
	)

	tests := []struct {
		name      string
		userID    string
		email     string
		roles     []string
		wantError bool
	}{
		{
			name:      "valid token generation",
			userID:    "user123",
			email:     "test@example.com",
			roles:     []string{"user", "admin"},
			wantError: false,
		},
		{
			name:      "empty user ID",
			userID:    "",
			email:     "test@example.com",
			roles:     []string{"user"},
			wantError: true,
		},
		{
			name:      "empty email",
			userID:    "user123",
			email:     "",
			roles:     []string{"user"},
			wantError: true,
		},
		{
			name:      "nil roles",
			userID:    "user123",
			email:     "test@example.com",
			roles:     nil,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, expiry, err := service.GenerateAccessToken(tt.userID, tt.email, tt.roles)

			if tt.wantError {
				assert.Error(t, err)
				assert.Empty(t, token)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, token)
			assert.True(t, expiry.After(time.Now()))

			// Validate the generated token
			claims, err := service.ValidateAccessToken(token)
			require.NoError(t, err)
			assert.Equal(t, tt.userID, claims.UserID)
			assert.Equal(t, tt.email, claims.Email)
			assert.Equal(t, tt.roles, claims.Roles)
			assert.Equal(t, "test-issuer", claims.Issuer)
		})
	}
}

func TestJWTService_GenerateRefreshToken(t *testing.T) {
	// Create JWT service with test configuration
	service := NewService(
		"test-access-secret",
		"test-refresh-secret",
		time.Hour,
		time.Hour*24*7,
		"test-issuer",
	)

	tests := []struct {
		name      string
		userID    string
		wantError bool
	}{
		{
			name:      "valid refresh token generation",
			userID:    "user123",
			wantError: false,
		},
		{
			name:      "empty user ID",
			userID:    "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, expiry, err := service.GenerateRefreshToken(tt.userID)

			if tt.wantError {
				assert.Error(t, err)
				assert.Empty(t, token)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, token)
			assert.True(t, expiry.After(time.Now()))

			// Validate the generated refresh token
			userID, err := service.ValidateRefreshToken(token)
			require.NoError(t, err)
			assert.Equal(t, tt.userID, userID)
		})
	}
}

func TestJWTService_ValidateAccessToken(t *testing.T) {
	// Create JWT service with test configuration
	service := NewService(
		"test-access-secret",
		"test-refresh-secret",
		time.Hour,
		time.Hour*24*7,
		"test-issuer",
	)

	// Generate a valid token for testing
	validToken, _, err := service.GenerateAccessToken("user123", "test@example.com", []string{"user"})
	require.NoError(t, err)

	tests := []struct {
		name      string
		token     string
		wantError bool
	}{
		{
			name:      "valid token",
			token:     validToken,
			wantError: false,
		},
		{
			name:      "empty token",
			token:     "",
			wantError: true,
		},
		{
			name:      "invalid token format",
			token:     "invalid.token.format",
			wantError: true,
		},
		{
			name:      "malformed token",
			token:     "malformed-token",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := service.ValidateAccessToken(tt.token)

			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, claims)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, claims)
			assert.Equal(t, "user123", claims.UserID)
			assert.Equal(t, "test@example.com", claims.Email)
			assert.Equal(t, []string{"user"}, claims.Roles)
		})
	}
}

func TestJWTService_ValidateRefreshToken(t *testing.T) {
	// Create JWT service with test configuration
	service := NewService(
		"test-access-secret",
		"test-refresh-secret",
		time.Hour,
		time.Hour*24*7,
		"test-issuer",
	)

	// Generate a valid refresh token for testing
	validToken, _, err := service.GenerateRefreshToken("user123")
	require.NoError(t, err)

	tests := []struct {
		name      string
		token     string
		wantError bool
	}{
		{
			name:      "valid refresh token",
			token:     validToken,
			wantError: false,
		},
		{
			name:      "empty token",
			token:     "",
			wantError: true,
		},
		{
			name:      "invalid token format",
			token:     "invalid.token.format",
			wantError: true,
		},
		{
			name:      "malformed token",
			token:     "malformed-token",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := service.ValidateRefreshToken(tt.token)

			if tt.wantError {
				assert.Error(t, err)
				assert.Empty(t, userID)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, "user123", userID)
		})
	}
}

func TestJWTService_TokenBlacklist(t *testing.T) {
	// Create JWT service with test configuration
	service := NewService(
		"test-access-secret",
		"test-refresh-secret",
		time.Hour,
		time.Hour*24*7,
		"test-issuer",
	)

	// Generate tokens for testing
	token1, _, err := service.GenerateAccessToken("user123", "test@example.com", []string{"user"})
	require.NoError(t, err)

	token2, _, err := service.GenerateAccessToken("user456", "test2@example.com", []string{"user"})
	require.NoError(t, err)

	// Test initial state
	assert.False(t, service.IsBlacklisted(token1))
	assert.False(t, service.IsBlacklisted(token2))

	// Blacklist token1
	service.BlacklistToken(token1)
	assert.True(t, service.IsBlacklisted(token1))
	assert.False(t, service.IsBlacklisted(token2))

	// Blacklist token2
	service.BlacklistToken(token2)
	assert.True(t, service.IsBlacklisted(token1))
	assert.True(t, service.IsBlacklisted(token2))

	// Test non-existent token
	assert.False(t, service.IsBlacklisted("non-existent-token"))
}

func TestJWTService_GetRefreshTTL(t *testing.T) {
	expectedTTL := time.Hour * 24 * 7
	service := NewService(
		"test-access-secret",
		"test-refresh-secret",
		time.Hour,
		expectedTTL,
		"test-issuer",
	)

	assert.Equal(t, expectedTTL, service.GetRefreshTTL())
}

func TestJWTService_ValidateAccessToken_EdgeCases(t *testing.T) {
	// Create JWT service with test configuration
	service := NewService(
		"test-access-secret",
		"test-refresh-secret",
		time.Hour,
		time.Hour*24*7,
		"test-issuer",
	)

	// Generate a valid token for base case
	validToken, _, err := service.GenerateAccessToken("user123", "test@example.com", []string{"user"})
	require.NoError(t, err)

	// Create a token with incorrect signing method
	wrongMethodToken := createTokenWithCustomSigningMethod(t, jwt.SigningMethodHS384, "test-access-secret")

	// Create an expired token
	expiredToken := createExpiredToken(t, service)

	// Create a token with future NBF
	futureToken := createFutureToken(t, service)

	// Create a token with wrong issuer
	wrongIssuerToken := createTokenWithWrongIssuer(t, service)

	// Create a token with missing claims
	missingClaimsToken := createTokenWithMissingClaims(t, service)

	tests := []struct {
		name      string
		token     string
		wantError bool
		errorMsg  string
	}{
		{
			name:      "valid token",
			token:     validToken,
			wantError: false,
		},
		{
			name:      "wrong signing method",
			token:     wrongMethodToken,
			wantError: true,
			errorMsg:  "unexpected signing method",
		},
		{
			name:      "expired token",
			token:     expiredToken,
			wantError: true,
			errorMsg:  "token is expired",
		},
		{
			name:      "future token",
			token:     futureToken,
			wantError: true,
			errorMsg:  "token not valid yet",
		},
		{
			name:      "wrong issuer",
			token:     wrongIssuerToken,
			wantError: true,
			errorMsg:  "invalid issuer",
		},
		{
			name:      "missing claims",
			token:     missingClaimsToken,
			wantError: true,
			errorMsg:  "missing required claims",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := service.ValidateAccessToken(tt.token)

			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, claims)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, claims)
		})
	}
}

func TestJWTService_TokenBlacklist_Enhanced(t *testing.T) {
	// Create JWT service with test configuration
	service := NewService(
		"test-access-secret",
		"test-refresh-secret",
		time.Hour,
		time.Hour*24*7,
		"test-issuer",
	)

	// Generate tokens for testing
	validToken, _, err := service.GenerateAccessToken("user123", "test@example.com", []string{"user"})
	require.NoError(t, err)

	// Create an expired token
	expiredToken := createExpiredToken(t, service)

	t.Run("blacklist operations", func(t *testing.T) {
		// Test blacklisting valid token
		assert.False(t, service.IsBlacklisted(validToken))
		service.BlacklistToken(validToken)
		assert.True(t, service.IsBlacklisted(validToken))

		// Test blacklisting expired token
		assert.False(t, service.IsBlacklisted(expiredToken))
		service.BlacklistToken(expiredToken)
		assert.True(t, service.IsBlacklisted(expiredToken))

		// Test blacklisting invalid token
		invalidToken := "invalid.token.format"
		service.BlacklistToken(invalidToken)
		assert.True(t, service.IsBlacklisted(invalidToken))

		// Test non-existent token
		assert.False(t, service.IsBlacklisted("non-existent-token"))
	})

	t.Run("concurrent blacklist access", func(t *testing.T) {
		tokens := make([]string, 100)
		for i := 0; i < 100; i++ {
			token, _, err := service.GenerateAccessToken(
				fmt.Sprintf("user%d", i),
				fmt.Sprintf("test%d@example.com", i),
				[]string{"user"},
			)
			require.NoError(t, err)
			tokens[i] = token
		}

		var wg sync.WaitGroup
		wg.Add(2)

		// Goroutine to blacklist tokens
		go func() {
			defer wg.Done()
			for _, token := range tokens[:50] {
				service.BlacklistToken(token)
			}
		}()

		// Goroutine to check blacklisted tokens
		go func() {
			defer wg.Done()
			for _, token := range tokens[50:] {
				service.IsBlacklisted(token)
			}
		}()

		wg.Wait()

		// Verify final state
		for _, token := range tokens[:50] {
			assert.True(t, service.IsBlacklisted(token))
		}
	})
}

// Helper functions for creating test tokens

func createTokenWithCustomSigningMethod(t *testing.T, method jwt.SigningMethod, secret string) string {
	token := jwt.New(method)
	token.Claims = jwt.MapClaims{
		"user_id": "user123",
		"email":   "test@example.com",
		"roles":   []string{"user"},
		"exp":     time.Now().Add(time.Hour).Unix(),
		"iat":     time.Now().Unix(),
		"iss":     "test-issuer",
	}
	tokenString, err := token.SignedString([]byte(secret))
	require.NoError(t, err)
	return tokenString
}

func createExpiredToken(t *testing.T, service *Service) string {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = jwt.MapClaims{
		"user_id": "user123",
		"email":   "test@example.com",
		"roles":   []string{"user"},
		"exp":     time.Now().Add(-time.Hour).Unix(),
		"iat":     time.Now().Add(-2 * time.Hour).Unix(),
		"iss":     "test-issuer",
	}
	tokenString, err := token.SignedString([]byte("test-access-secret"))
	require.NoError(t, err)
	return tokenString
}

func createFutureToken(t *testing.T, service *Service) string {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = jwt.MapClaims{
		"user_id": "user123",
		"email":   "test@example.com",
		"roles":   []string{"user"},
		"exp":     time.Now().Add(2 * time.Hour).Unix(),
		"iat":     time.Now().Add(time.Hour).Unix(),
		"nbf":     time.Now().Add(time.Hour).Unix(),
		"iss":     "test-issuer",
	}
	tokenString, err := token.SignedString([]byte("test-access-secret"))
	require.NoError(t, err)
	return tokenString
}

func createTokenWithWrongIssuer(t *testing.T, service *Service) string {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = jwt.MapClaims{
		"user_id": "user123",
		"email":   "test@example.com",
		"roles":   []string{"user"},
		"exp":     time.Now().Add(time.Hour).Unix(),
		"iat":     time.Now().Unix(),
		"iss":     "wrong-issuer",
	}
	tokenString, err := token.SignedString([]byte("test-access-secret"))
	require.NoError(t, err)
	return tokenString
}

func createTokenWithMissingClaims(t *testing.T, service *Service) string {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = jwt.MapClaims{
		// Missing user_id and email
		"roles": []string{"user"},
		"exp":   time.Now().Add(time.Hour).Unix(),
		"iat":   time.Now().Unix(),
		"iss":   "test-issuer",
	}
	tokenString, err := token.SignedString([]byte("test-access-secret"))
	require.NoError(t, err)
	return tokenString
}
