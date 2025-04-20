package apperror_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/apperror"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorHandling(t *testing.T) {
	t.Run("RespondWithError handles AppError correctly", func(t *testing.T) {
		// Arrange
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/test", nil)
		r.Header.Set("X-Request-ID", "test-trace-id")
		err := apperror.NewNotFound("user", "123", nil)

		// Act
		apperror.RespondWithError(w, r, err)

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

	t.Run("RespondWithError handles standard error", func(t *testing.T) {
		// Arrange
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/test", nil)
		r.Header.Set("X-Request-ID", "test-trace-id")
		err := errors.New("standard error")

		// Act
		apperror.RespondWithError(w, r, err)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]interface{}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))

		errorData, ok := response["error"].(map[string]interface{})
		require.True(t, ok, "Expected 'error' key in response")
		assert.Equal(t, "INTERNAL_ERROR", errorData["code"])
		assert.Equal(t, "Internal server error", errorData["message"])
		assert.Equal(t, "test-trace-id", errorData["trace_id"])
	})

	t.Run("WriteValidationError formats validation errors", func(t *testing.T) {
		// Arrange
		w := httptest.NewRecorder()
		traceID := "test-trace-id"

		// Act
		apperror.WriteValidationError(w, "email", "Email is invalid", traceID)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))

		errorData, ok := response["error"].(map[string]interface{})
		require.True(t, ok, "Expected 'error' key in response")
		assert.Equal(t, "INVALID_INPUT", errorData["code"])
		assert.Equal(t, "Validation error", errorData["message"])
		assert.Equal(t, traceID, errorData["trace_id"])

		fieldErrors, ok := errorData["field_errors"].(map[string]interface{})
		require.True(t, ok, "Expected 'field_errors' in error data")
		assert.Equal(t, "Email is invalid", fieldErrors["email"])
	})

	t.Run("WriteValidationErrors handles multiple errors", func(t *testing.T) {
		// Arrange
		w := httptest.NewRecorder()
		traceID := "test-trace-id"
		fieldErrors := map[string]string{
			"email":    "Email is invalid",
			"username": "Username is too short",
			"password": "Password is required",
		}

		// Act
		apperror.WriteValidationErrors(w, fieldErrors, traceID)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))

		errorData, ok := response["error"].(map[string]interface{})
		require.True(t, ok, "Expected 'error' key in response")
		assert.Equal(t, "INVALID_INPUT", errorData["code"])
		assert.Equal(t, "Validation errors", errorData["message"])
		assert.Equal(t, traceID, errorData["trace_id"])

		respFieldErrors, ok := errorData["field_errors"].(map[string]interface{})
		require.True(t, ok, "Expected 'field_errors' in error data")
		assert.Equal(t, 3, len(respFieldErrors))
		assert.Equal(t, "Email is invalid", respFieldErrors["email"])
		assert.Equal(t, "Username is too short", respFieldErrors["username"])
		assert.Equal(t, "Password is required", respFieldErrors["password"])
	})

	t.Run("WrapError preserves error type", func(t *testing.T) {
		// Arrange
		originalErr := apperror.NewNotFound("user", "123", nil)

		// Act
		wrappedErr := apperror.WrapError(originalErr, "Failed to retrieve user")

		// Assert
		var appErr *apperror.AppError
		assert.True(t, errors.As(wrappedErr, &appErr))
		assert.Equal(t, http.StatusNotFound, appErr.StatusCode)
		assert.Equal(t, "NOT_FOUND", appErr.Code)
		assert.Equal(t, "Failed to retrieve user: user with identifier 123 not found", appErr.Message)
		assert.True(t, apperror.IsNotFound(wrappedErr))
	})

	t.Run("Error type checking works", func(t *testing.T) {
		// Arrange
		notFoundErr := apperror.NewNotFound("user", "123", nil)
		unauthorizedErr := apperror.NewUnauthorized("Invalid token", nil)
		internalErr := apperror.NewInternal(errors.New("db error"))
		standardErr := errors.New("standard error")

		// Act & Assert
		assert.True(t, apperror.IsNotFound(notFoundErr))
		assert.False(t, apperror.IsNotFound(unauthorizedErr))

		assert.True(t, apperror.IsUnauthorized(unauthorizedErr))
		assert.False(t, apperror.IsUnauthorized(notFoundErr))

		assert.True(t, apperror.IsInternal(internalErr))
		assert.False(t, apperror.IsInternal(standardErr))

		assert.Equal(t, http.StatusNotFound, apperror.GetStatusCode(notFoundErr))
		assert.Equal(t, http.StatusInternalServerError, apperror.GetStatusCode(standardErr))
	})

	t.Run("ErrorContext works", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		customHandler := func(w http.ResponseWriter, err error, traceID string) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Custom-Header", "custom-value")
			w.WriteHeader(http.StatusTeapot)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"custom_error": err.Error(),
				"trace":        traceID,
			})
		}

		// Act
		ctxWithHandler := apperror.WithErrorHandler(ctx, customHandler)
		handler := apperror.GetErrorHandler(ctxWithHandler)

		// Test the handler
		w := httptest.NewRecorder()
		err := errors.New("test error")
		handler(w, err, "test-trace-id")

		// Assert
		assert.Equal(t, http.StatusTeapot, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		assert.Equal(t, "custom-value", w.Header().Get("X-Custom-Header"))

		var response map[string]interface{}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
		assert.Equal(t, "test error", response["custom_error"])
		assert.Equal(t, "test-trace-id", response["trace"])
	})

	t.Run("DefaultErrorHandler is used when no handler in context", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		// Act
		handler := apperror.GetErrorHandler(ctx)

		// Test the handler
		w := httptest.NewRecorder()
		err := apperror.NewNotFound("user", "123", nil)
		handler(w, err, "test-trace-id")

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
}
