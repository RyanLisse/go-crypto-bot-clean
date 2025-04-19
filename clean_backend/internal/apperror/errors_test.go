package apperror_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/apperror"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppError(t *testing.T) {
	t.Run("Error returns message with wrapped error", func(t *testing.T) {
		err := apperror.NewBadRequest("Invalid input", nil, errors.New("validation failed"))
		assert.Equal(t, "Invalid input: validation failed", err.Error())
	})

	t.Run("Error returns message without wrapped error", func(t *testing.T) {
		err := apperror.NewBadRequest("Invalid input", nil, nil)
		assert.Equal(t, "Invalid input", err.Error())
	})

	t.Run("Unwrap returns wrapped error", func(t *testing.T) {
		wrappedErr := errors.New("validation failed")
		err := apperror.NewBadRequest("Invalid input", nil, wrappedErr)
		assert.Equal(t, wrappedErr, errors.Unwrap(err))
	})

	t.Run("GetStatusCode returns StatusCode when set", func(t *testing.T) {
		err := &apperror.AppError{
			Code:       http.StatusBadRequest,
			StatusCode: http.StatusTeapot,
		}
		assert.Equal(t, http.StatusTeapot, err.GetStatusCode())
	})

	t.Run("GetStatusCode returns Code when StatusCode not set", func(t *testing.T) {
		err := &apperror.AppError{
			Code:       http.StatusBadRequest,
			StatusCode: 0,
		}
		assert.Equal(t, http.StatusBadRequest, err.GetStatusCode())
	})

	t.Run("GetErrorCode returns ErrorCode when set", func(t *testing.T) {
		err := &apperror.AppError{
			Code:      http.StatusBadRequest,
			ErrorCode: apperror.ErrCodeValidationError,
		}
		assert.Equal(t, "VALIDATION_ERROR", err.GetErrorCode())
	})

	t.Run("GetErrorCode returns mapped code when ErrorCode not set", func(t *testing.T) {
		err := &apperror.AppError{
			Code:      http.StatusBadRequest,
			ErrorCode: "",
		}
		assert.Equal(t, "BAD_REQUEST", err.GetErrorCode())
	})
}

func TestErrorCreationFunctions(t *testing.T) {
	t.Run("NewNotFound creates correct error", func(t *testing.T) {
		err := apperror.NewNotFound("User", "123", nil)
		assert.Equal(t, http.StatusNotFound, err.GetStatusCode())
		assert.Equal(t, "NOT_FOUND", err.GetErrorCode())
		assert.Contains(t, err.Message, "User with identifier 123 not found")
	})

	t.Run("NewBadRequest creates correct error", func(t *testing.T) {
		err := apperror.NewBadRequest("Invalid input", nil, nil)
		assert.Equal(t, http.StatusBadRequest, err.GetStatusCode())
		assert.Equal(t, "BAD_REQUEST", err.GetErrorCode())
		assert.Equal(t, "Invalid input", err.Message)
	})

	t.Run("NewUnauthorized creates correct error", func(t *testing.T) {
		err := apperror.NewUnauthorized("Authentication required", nil)
		assert.Equal(t, http.StatusUnauthorized, err.GetStatusCode())
		assert.Equal(t, "UNAUTHORIZED", err.GetErrorCode())
		assert.Equal(t, "Authentication required", err.Message)
	})

	t.Run("NewForbidden creates correct error", func(t *testing.T) {
		err := apperror.NewForbidden("Access denied", nil)
		assert.Equal(t, http.StatusForbidden, err.GetStatusCode())
		assert.Equal(t, "FORBIDDEN", err.GetErrorCode())
		assert.Equal(t, "Access denied", err.Message)
	})

	t.Run("NewInternal creates correct error", func(t *testing.T) {
		err := apperror.NewInternal(errors.New("database error"))
		assert.Equal(t, http.StatusInternalServerError, err.GetStatusCode())
		assert.Equal(t, "INTERNAL_ERROR", err.GetErrorCode())
		assert.Equal(t, "Internal server error", err.Message)
	})

	t.Run("NewValidationError creates correct error", func(t *testing.T) {
		details := map[string]string{"email": "Invalid format"}
		err := apperror.NewValidationError("Validation failed", details, nil)
		assert.Equal(t, http.StatusUnprocessableEntity, err.GetStatusCode())
		assert.Equal(t, "VALIDATION_ERROR", err.GetErrorCode())
		assert.Equal(t, "Validation failed", err.Message)
		assert.Equal(t, details, err.Details)
	})

	t.Run("NewRateLimit creates correct error", func(t *testing.T) {
		err := apperror.NewRateLimit("Too many requests", nil)
		assert.Equal(t, http.StatusTooManyRequests, err.GetStatusCode())
		assert.Equal(t, "RATE_LIMIT_EXCEEDED", err.GetErrorCode())
		assert.Equal(t, "Too many requests", err.Message)
	})
}

func TestWriteError(t *testing.T) {
	t.Run("WriteError writes correct response", func(t *testing.T) {
		w := httptest.NewRecorder()
		err := apperror.NewBadRequest("Invalid input", map[string]string{"field": "error"}, nil)
		
		apperror.WriteError(w, err)
		
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		
		var response map[string]interface{}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
		
		assert.False(t, response["success"].(bool))
		errorData := response["error"].(map[string]interface{})
		assert.Equal(t, "BAD_REQUEST", errorData["code"])
		assert.Equal(t, "Invalid input", errorData["message"])
		
		details := errorData["details"].(map[string]interface{})
		assert.Equal(t, "error", details["field"])
	})
}

func TestRespondWithError(t *testing.T) {
	t.Run("RespondWithError includes trace ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/test", nil)
		r.Header.Set("X-Request-ID", "test-trace-id")
		
		err := apperror.NewNotFound("User", "123", nil)
		apperror.RespondWithError(w, r, err)
		
		assert.Equal(t, http.StatusNotFound, w.Code)
		
		var response map[string]interface{}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
		
		errorData := response["error"].(map[string]interface{})
		assert.Equal(t, "test-trace-id", errorData["trace_id"])
	})
}

func TestFromError(t *testing.T) {
	t.Run("FromError returns AppError unchanged", func(t *testing.T) {
		original := apperror.NewBadRequest("Invalid input", nil, nil)
		result := apperror.FromError(original)
		assert.Same(t, original, result)
	})
	
	t.Run("FromError converts standard error to AppError", func(t *testing.T) {
		result := apperror.FromError(errors.New("standard error"))
		assert.Equal(t, http.StatusInternalServerError, result.GetStatusCode())
		assert.Equal(t, "INTERNAL_ERROR", result.GetErrorCode())
	})
	
	t.Run("FromError converts known error types correctly", func(t *testing.T) {
		result := apperror.FromError(apperror.ErrNotFound)
		assert.Equal(t, http.StatusNotFound, result.GetStatusCode())
		assert.Equal(t, "NOT_FOUND", result.GetErrorCode())
		
		result = apperror.FromError(apperror.ErrUnauthorized)
		assert.Equal(t, http.StatusUnauthorized, result.GetStatusCode())
		assert.Equal(t, "UNAUTHORIZED", result.GetErrorCode())
	})
}
