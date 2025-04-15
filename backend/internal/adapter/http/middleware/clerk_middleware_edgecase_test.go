package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestClerkMiddleware_EdgeCases(t *testing.T) {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	middleware := NewClerkMiddleware("test_secret_key", &logger)

	t.Run("No Authorization Header", func(t *testing.T) {
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("no auth header"))
		})
		req := httptest.NewRequest("GET", "/test", nil)
		res := httptest.NewRecorder()
		handler := middleware.Middleware()(testHandler)
		handler.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, "no auth header", res.Body.String())
	})

	t.Run("Malformed Authorization Header", func(t *testing.T) {
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("Handler should not be called with malformed auth header")
		})
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer")
		res := httptest.NewRecorder()
		handler := middleware.Middleware()(testHandler)
		handler.ServeHTTP(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Code)
	})

	t.Run("Invalid Clerk Token", func(t *testing.T) {
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("Handler should not be called with invalid token")
		})
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer invalidtoken")
		res := httptest.NewRecorder()
		handler := middleware.Middleware()(testHandler)
		handler.ServeHTTP(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Code)
	})

	// For valid token edge cases, you would mock Clerk SDK responses
}
