package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestClerkMiddleware_NoAuthHeader(t *testing.T) {
	// Create a logger
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create middleware
	middleware := NewClerkMiddleware("test_secret_key", &logger)

	// Create a test handler that will be wrapped
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Should be called since we're not providing an auth header
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test passed"))
	})

	// Create a request without auth header
	req := httptest.NewRequest("GET", "/test", nil)
	res := httptest.NewRecorder()

	// Apply the middleware
	handler := middleware.Middleware()(testHandler)

	// Send the request
	handler.ServeHTTP(res, req)

	// Check the response
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "test passed", res.Body.String())
}

func TestClerkMiddleware_InvalidAuthFormat(t *testing.T) {
	// Create a logger
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create middleware
	middleware := NewClerkMiddleware("test_secret_key", &logger)

	// Create a test handler that should not be called
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called with invalid auth format")
	})

	// Create a request with invalid auth header (missing Bearer prefix)
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "invalid-format")
	res := httptest.NewRecorder()

	// Apply the middleware
	handler := middleware.Middleware()(testHandler)

	// Send the request
	handler.ServeHTTP(res, req)

	// Check that it returned unauthorized
	assert.Equal(t, http.StatusUnauthorized, res.Code)
}

// This test will always fail in a real environment because we can't
// verify a token without a real Clerk secret key and valid token
// This is essentially a structural test
func TestClerkMiddleware_WithAuthHeader(t *testing.T) {
	t.Skip("Skipping test that requires real Clerk credentials")

	// Create a logger
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create middleware - in a real test you'd use a valid key
	middleware := NewClerkMiddleware("test_secret_key", &logger)

	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// In a real test with a valid token, check that user context values are set
		userID := r.Context().Value(UserIDKey)
		if userID == nil {
			t.Error("UserID should be set in context")
		}
		w.WriteHeader(http.StatusOK)
	})

	// Create a request with auth header
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer some-valid-token") // would be a real token in actual test
	res := httptest.NewRecorder()

	// Apply the middleware
	handler := middleware.Middleware()(testHandler)

	// Send the request
	handler.ServeHTTP(res, req)
}
