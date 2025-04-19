package middleware_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/adapter/delivery/http/middleware"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/apperror"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnifiedErrorMiddleware(t *testing.T) {
	// Set up logger for tests
	logger := zerolog.New(io.Discard)

	t.Run("Adds request ID when missing", func(t *testing.T) {
		// Arrange
		middleware := middleware.NewUnifiedErrorMiddleware(&logger)
		handler := middleware.Middleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Just echo the request ID
			requestID := r.Header.Get("X-Request-ID")
			w.Write([]byte(requestID))
		}))

		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		// Act
		handler.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		requestID := w.Body.String()
		assert.NotEmpty(t, requestID, "Request ID should be generated")
		assert.Equal(t, requestID, w.Header().Get("X-Request-ID"))
	})

	t.Run("Uses existing request ID", func(t *testing.T) {
		// Arrange
		middleware := middleware.NewUnifiedErrorMiddleware(&logger)
		handler := middleware.Middleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Just echo the request ID
			requestID := r.Header.Get("X-Request-ID")
			w.Write([]byte(requestID))
		}))

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("X-Request-ID", "test-trace-id")
		w := httptest.NewRecorder()

		// Act
		handler.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "test-trace-id", w.Body.String())
		assert.Equal(t, "test-trace-id", w.Header().Get("X-Request-ID"))
	})

	t.Run("Captures and handles errors", func(t *testing.T) {
		// Arrange
		middleware := middleware.NewUnifiedErrorMiddleware(&logger)
		handler := middleware.Middleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Trigger an error using the context handler
			err := apperror.NewNotFound("user", "123", nil)
			apperror.RespondWithError(w, r, err)
		}))

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("X-Request-ID", "test-trace-id")
		w := httptest.NewRecorder()

		// Act
		handler.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]interface{}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))

		errorData, ok := response["error"].(map[string]interface{})
		require.True(t, ok, "Expected 'error' key in response")
		assert.Equal(t, "NOT_FOUND", errorData["code"])
		assert.Equal(t, "user with identifier 123 not found", errorData["message"])
		assert.Equal(t, "test-trace-id", errorData["trace_id"])
	})

	t.Run("Recovers from panics", func(t *testing.T) {
		// Arrange
		middleware := middleware.NewUnifiedErrorMiddleware(&logger)
		handler := middleware.Middleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Trigger a panic
			panic("test panic")
		}))

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("X-Request-ID", "test-trace-id")
		w := httptest.NewRecorder()

		// Act
		handler.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]interface{}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))

		errorData, ok := response["error"].(map[string]interface{})
		require.True(t, ok, "Expected 'error' key in response")
		assert.Equal(t, "INTERNAL_ERROR", errorData["code"])
		assert.True(t, strings.Contains(errorData["message"].(string), "Internal server error"))
		assert.Equal(t, "test-trace-id", errorData["trace_id"])
	})

	t.Run("Recovers from error panics", func(t *testing.T) {
		// Arrange
		middleware := middleware.NewUnifiedErrorMiddleware(&logger)
		handler := middleware.Middleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Trigger a panic with an error
			panic(errors.New("error panic"))
		}))

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("X-Request-ID", "test-trace-id")
		w := httptest.NewRecorder()

		// Act
		handler.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Recovers from app error panics", func(t *testing.T) {
		// Arrange
		middleware := middleware.NewUnifiedErrorMiddleware(&logger)
		handler := middleware.Middleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Trigger a panic with an AppError
			panic(apperror.NewNotFound("resource", "id", nil))
		}))

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("X-Request-ID", "test-trace-id")
		w := httptest.NewRecorder()

		// Act
		handler.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]interface{}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))

		errorData, ok := response["error"].(map[string]interface{})
		require.True(t, ok, "Expected 'error' key in response")
		assert.Equal(t, "NOT_FOUND", errorData["code"])
	})

	t.Run("Captures response status code", func(t *testing.T) {
		// Arrange
		middleware := middleware.NewUnifiedErrorMiddleware(&logger)
		handler := middleware.Middleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set a non-default status code
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("created"))
		}))

		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		// Act
		handler.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, "created", w.Body.String())
	})
}
