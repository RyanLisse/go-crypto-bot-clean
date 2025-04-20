package response_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/delivery/http/response"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/infrastructure/errorhandler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnifiedResponse(t *testing.T) {
	t.Run("NewSuccessResponse creates correct structure", func(t *testing.T) {
		// Arrange
		data := map[string]string{"key": "value"}
		
		// Act
		resp := response.NewSuccessResponse(data)
		
		// Assert
		assert.True(t, resp.Success)
		assert.Equal(t, data, resp.Data)
		assert.Nil(t, resp.Error)
		assert.NotEmpty(t, resp.Timestamp)
		assert.Equal(t, response.APIVersion, resp.Version)
	})

	t.Run("NewErrorResponse creates correct structure", func(t *testing.T) {
		// Act
		resp := response.NewErrorResponse("NOT_FOUND", "Resource not found")
		
		// Assert
		assert.False(t, resp.Success)
		assert.Nil(t, resp.Data)
		require.NotNil(t, resp.Error)
		assert.Equal(t, "NOT_FOUND", resp.Error.Code)
		assert.Equal(t, "Resource not found", resp.Error.Message)
		assert.NotEmpty(t, resp.Timestamp)
		assert.Equal(t, response.APIVersion, resp.Version)
	})

	t.Run("NewErrorResponseWithDetails creates correct structure", func(t *testing.T) {
		// Arrange
		details := map[string]string{"field": "error"}
		
		// Act
		resp := response.NewErrorResponseWithDetails("VALIDATION_ERROR", "Validation failed", details)
		
		// Assert
		assert.False(t, resp.Success)
		require.NotNil(t, resp.Error)
		assert.Equal(t, "VALIDATION_ERROR", resp.Error.Code)
		assert.Equal(t, "Validation failed", resp.Error.Message)
		assert.Equal(t, details, resp.Error.Details)
		assert.NotEmpty(t, resp.Timestamp)
		assert.Equal(t, response.APIVersion, resp.Version)
	})

	t.Run("NewErrorResponseWithTraceID creates correct structure", func(t *testing.T) {
		// Act
		resp := response.NewErrorResponseWithTraceID("NOT_FOUND", "Resource not found", "trace-123")
		
		// Assert
		assert.False(t, resp.Success)
		require.NotNil(t, resp.Error)
		assert.Equal(t, "NOT_FOUND", resp.Error.Code)
		assert.Equal(t, "Resource not found", resp.Error.Message)
		assert.Equal(t, "trace-123", resp.Error.TraceID)
		assert.NotEmpty(t, resp.Timestamp)
		assert.Equal(t, response.APIVersion, resp.Version)
	})

	t.Run("NewErrorResponseFromAppError creates correct structure", func(t *testing.T) {
		// Arrange
		appErr := errorhandler.NewNotFound("user", "123", nil)
		appErr.TraceID = "trace-123"
		
		// Act
		resp := response.NewErrorResponseFromAppError(appErr)
		
		// Assert
		assert.False(t, resp.Success)
		require.NotNil(t, resp.Error)
		assert.Equal(t, "NOT_FOUND", resp.Error.Code)
		assert.Equal(t, "user with identifier 123 not found", resp.Error.Message)
		assert.Equal(t, "trace-123", resp.Error.TraceID)
		assert.NotEmpty(t, resp.Timestamp)
		assert.Equal(t, response.APIVersion, resp.Version)
	})
}

func TestResponseWriters(t *testing.T) {
	t.Run("WriteUnifiedJSON writes JSON response", func(t *testing.T) {
		// Arrange
		w := httptest.NewRecorder()
		data := map[string]string{"key": "value"}
		
		// Act
		response.WriteUnifiedJSON(w, http.StatusOK, data)
		
		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		
		var result map[string]string
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
		assert.Equal(t, "value", result["key"])
	})

	t.Run("WriteUnifiedSuccess writes success response", func(t *testing.T) {
		// Arrange
		w := httptest.NewRecorder()
		data := map[string]string{"key": "value"}
		
		// Act
		response.WriteUnifiedSuccess(w, http.StatusOK, data)
		
		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		
		var result response.UnifiedResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
		assert.True(t, result.Success)
		
		dataMap, ok := result.Data.(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, "value", dataMap["key"])
		assert.NotEmpty(t, result.Timestamp)
		assert.Equal(t, response.APIVersion, result.Version)
	})

	t.Run("WriteUnifiedError writes error response", func(t *testing.T) {
		// Arrange
		w := httptest.NewRecorder()
		
		// Act
		response.WriteUnifiedError(w, http.StatusNotFound, "NOT_FOUND", "Resource not found")
		
		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		
		var result response.UnifiedResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
		assert.False(t, result.Success)
		
		require.NotNil(t, result.Error)
		assert.Equal(t, "NOT_FOUND", result.Error.Code)
		assert.Equal(t, "Resource not found", result.Error.Message)
		assert.NotEmpty(t, result.Timestamp)
		assert.Equal(t, response.APIVersion, result.Version)
	})

	t.Run("WriteUnifiedAppError writes app error response", func(t *testing.T) {
		// Arrange
		w := httptest.NewRecorder()
		appErr := errorhandler.NewNotFound("user", "123", nil)
		appErr.TraceID = "trace-123"
		
		// Act
		response.WriteUnifiedAppError(w, appErr)
		
		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		
		var result response.UnifiedResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
		assert.False(t, result.Success)
		
		require.NotNil(t, result.Error)
		assert.Equal(t, "NOT_FOUND", result.Error.Code)
		assert.Equal(t, "user with identifier 123 not found", result.Error.Message)
		assert.Equal(t, "trace-123", result.Error.TraceID)
		assert.NotEmpty(t, result.Timestamp)
		assert.Equal(t, response.APIVersion, result.Version)
	})
}

func TestTimestampFormat(t *testing.T) {
	t.Run("Timestamp is in RFC3339 format", func(t *testing.T) {
		// Arrange
		resp := response.NewSuccessResponse(nil)
		
		// Act
		timestamp := resp.Timestamp
		
		// Assert
		_, err := time.Parse(time.RFC3339, timestamp)
		assert.NoError(t, err, "Timestamp should be in RFC3339 format")
	})
}
