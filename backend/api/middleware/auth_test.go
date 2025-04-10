package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-crypto-bot-clean/backend/api/middleware/jwt"

	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	// Create a JWT service
	jwtService := jwt.NewService(
		"test-access-secret",
		"test-refresh-secret",
		time.Hour,
		time.Hour*24*7,
		"test-issuer",
	)

	// Create an auth middleware
	authMiddleware := NewAuthMiddleware(jwtService)

	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := GetUserID(r)
		email := GetEmail(r)
		roles := GetRoles(r)

		assert.Equal(t, "user-123", userID, "User ID should match")
		assert.Equal(t, "user@example.com", email, "Email should match")
		assert.Equal(t, []string{"user", "admin"}, roles, "Roles should match")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Create a middleware chain
	handler := authMiddleware.Authenticate(testHandler)

	// Generate a valid token
	token, _, err := jwtService.GenerateAccessToken("user-123", "user@example.com", []string{"user", "admin"})
	assert.NoError(t, err, "Should not error when generating token")

	// Test with a valid token
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, "Should return 200 OK")
	assert.Equal(t, "OK", rec.Body.String(), "Should return OK")

	// Test with no Authorization header
	req = httptest.NewRequest("GET", "/", nil)
	rec = httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code, "Should return 401 Unauthorized")

	// Test with invalid Authorization header format
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	rec = httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code, "Should return 401 Unauthorized")

	// Test with invalid token
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	rec = httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code, "Should return 401 Unauthorized")
}

func TestRequireRole(t *testing.T) {
	// Create a JWT service
	jwtService := jwt.NewService(
		"test-access-secret",
		"test-refresh-secret",
		time.Hour,
		time.Hour*24*7,
		"test-issuer",
	)

	// Create an auth middleware
	authMiddleware := NewAuthMiddleware(jwtService)

	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Create a middleware chain
	handler := authMiddleware.Authenticate(authMiddleware.RequireRole("admin")(testHandler))

	// Generate a token with admin role
	adminToken, _, err := jwtService.GenerateAccessToken("user-123", "admin@example.com", []string{"admin"})
	assert.NoError(t, err, "Should not error when generating token")

	// Generate a token with user role
	userToken, _, err := jwtService.GenerateAccessToken("user-456", "user@example.com", []string{"user"})
	assert.NoError(t, err, "Should not error when generating token")

	// Test with admin role
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, "Should return 200 OK")
	assert.Equal(t, "OK", rec.Body.String(), "Should return OK")

	// Test with user role (should be forbidden)
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	rec = httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code, "Should return 403 Forbidden")
}

func TestRequireAdmin(t *testing.T) {
	// Create a JWT service
	jwtService := jwt.NewService(
		"test-access-secret",
		"test-refresh-secret",
		time.Hour,
		time.Hour*24*7,
		"test-issuer",
	)

	// Create an auth middleware
	authMiddleware := NewAuthMiddleware(jwtService)

	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Create a middleware chain
	handler := authMiddleware.Authenticate(authMiddleware.RequireAdmin()(testHandler))

	// Generate a token with admin role
	adminToken, _, err := jwtService.GenerateAccessToken("user-123", "admin@example.com", []string{"admin"})
	assert.NoError(t, err, "Should not error when generating token")

	// Generate a token with user role
	userToken, _, err := jwtService.GenerateAccessToken("user-456", "user@example.com", []string{"user"})
	assert.NoError(t, err, "Should not error when generating token")

	// Test with admin role
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, "Should return 200 OK")
	assert.Equal(t, "OK", rec.Body.String(), "Should return OK")

	// Test with user role (should be forbidden)
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	rec = httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code, "Should return 403 Forbidden")
}
