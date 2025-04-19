package apperror_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/apperror"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleError(t *testing.T) {
	// Setup logger for tests
	logger := zerolog.New(io.Discard)

	t.Run("Handles nil error", func(t *testing.T) {
		handler := apperror.HandleError(func(w http.ResponseWriter, r *http.Request) error {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success"))
			return nil
		}, &logger)

		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "success", w.Body.String())
	})

	t.Run("Handles AppError", func(t *testing.T) {
		handler := apperror.HandleError(func(w http.ResponseWriter, r *http.Request) error {
			return apperror.NewNotFound("User", "123", nil)
		}, &logger)

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("X-Request-ID", "test-trace-id")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]interface{}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))

		errorData := response["error"].(map[string]interface{})
		assert.Equal(t, "NOT_FOUND", errorData["code"])
		assert.Equal(t, "test-trace-id", errorData["trace_id"])
	})

	t.Run("Converts standard error to AppError", func(t *testing.T) {
		handler := apperror.HandleError(func(w http.ResponseWriter, r *http.Request) error {
			return errors.New("standard error")
		}, &logger)

		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))

		errorData := response["error"].(map[string]interface{})
		assert.Equal(t, "INTERNAL_ERROR", errorData["code"])
	})

	t.Run("Converts known error types correctly", func(t *testing.T) {
		handler := apperror.HandleError(func(w http.ResponseWriter, r *http.Request) error {
			return apperror.ErrNotFound
		}, &logger)

		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]interface{}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))

		errorData := response["error"].(map[string]interface{})
		assert.Equal(t, "NOT_FOUND", errorData["code"])
	})
}

func TestWithLogging(t *testing.T) {
	// Setup logger for tests
	logger := zerolog.New(io.Discard)

	t.Run("Logs request and calls handler", func(t *testing.T) {
		var handlerCalled bool
		handler := apperror.WithLogging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
			w.WriteHeader(http.StatusOK)
		}), &logger)

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("X-Request-ID", "test-trace-id")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.True(t, handlerCalled)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestRespondWithJSON(t *testing.T) {
	t.Run("Writes JSON response with data", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/test", nil)
		r.Header.Set("X-Request-ID", "test-trace-id")

		data := map[string]string{"name": "Test User"}
		err := apperror.RespondWithJSON(w, r, http.StatusOK, data)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]interface{}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))

		assert.True(t, response["success"].(bool))
		assert.Equal(t, "test-trace-id", response["trace_id"])
		
		responseData := response["data"].(map[string]interface{})
		assert.Equal(t, "Test User", responseData["name"])
	})
}

func TestRespondWithSuccess(t *testing.T) {
	t.Run("Writes success response with no data", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/test", nil)
		r.Header.Set("X-Request-ID", "test-trace-id")

		err := apperror.RespondWithSuccess(w, r, http.StatusNoContent)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]interface{}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))

		assert.True(t, response["success"].(bool))
		assert.Equal(t, "test-trace-id", response["trace_id"])
		_, hasData := response["data"]
		assert.False(t, hasData)
	})
}

func TestResponseHelpers(t *testing.T) {
	t.Run("RespondWithCreated sets correct status code", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/test", nil)
		
		data := map[string]string{"id": "123"}
		err := apperror.RespondWithCreated(w, r, data)
		
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, w.Code)
	})
	
	t.Run("RespondWithOK sets correct status code", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/test", nil)
		
		data := map[string]string{"name": "Test"}
		err := apperror.RespondWithOK(w, r, data)
		
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	
	t.Run("RespondWithNoContent sets correct status code", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("DELETE", "/test", nil)
		
		err := apperror.RespondWithNoContent(w, r)
		
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, w.Code)
	})
	
	t.Run("RespondWithAccepted sets correct status code", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/test", nil)
		
		data := map[string]string{"status": "processing"}
		err := apperror.RespondWithAccepted(w, r, data)
		
		assert.NoError(t, err)
		assert.Equal(t, http.StatusAccepted, w.Code)
	})
}
