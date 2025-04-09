package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-crypto-bot-clean/backend/internal/auth"

	"github.com/stretchr/testify/assert"
)

func TestRecoveryMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		handler        http.HandlerFunc
		expectedStatus int
		expectedType   string
	}{
		{
			name: "handles panic with error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				panic("test panic")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedType:   "internal_error",
		},
		{
			name: "handles panic with auth error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				panic(auth.ErrUnauthorized)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedType:   "unauthorized",
		},
		{
			name: "normal request",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			expectedStatus: http.StatusOK,
			expectedType:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test handler with recovery middleware
			handler := RecoveryMiddleware(tt.handler)

			// Create test request
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()

			// Add request ID to context for testing
			ctx := req.Context()
			ctx = context.WithValue(ctx, RequestIDContextKey, "test-request-id")
			req = req.WithContext(ctx)

			// Execute request
			handler.ServeHTTP(w, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// For non-panic cases, we're done
			if tt.expectedType == "" {
				return
			}

			// For panic cases, verify the error response
			var response auth.ErrorResponse
			err := json.NewDecoder(w.Body).Decode(&response)
			assert.NoError(t, err)

			// Verify error details
			assert.Equal(t, tt.expectedType, string(response.Error.Type))
			assert.Equal(t, "test-request-id", response.Error.RequestID)
			assert.NotNil(t, response.Error.Details)

			// Verify metadata
			details, ok := response.Error.Details.(map[string]interface{})
			assert.True(t, ok)
			assert.Contains(t, details, "stack_trace")
			assert.Contains(t, details, "request")

			// Verify request info
			reqInfo, ok := details["request"].(map[string]interface{})
			assert.True(t, ok)
			assert.Equal(t, "/test", reqInfo["path"])
			assert.Equal(t, "GET", reqInfo["method"])
		})
	}
}
