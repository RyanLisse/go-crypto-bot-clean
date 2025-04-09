package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRecoveryMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		shouldPanic bool
	}{
		{
			name:        "handles panic",
			shouldPanic: true,
		},
		{
			name:        "normal request",
			shouldPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock logger
			logger := &mockLogger{}

			// Create a test router with the recovery middleware
			router := gin.New()
			router.Use(RecoveryMiddleware(logger))

			// Define a test handler that may panic
			router.GET("/test", func(c *gin.Context) {
				if tt.shouldPanic {
					panic("test panic")
				}
				c.Status(http.StatusOK)
			})

			// Make a request
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Check the response
			if tt.shouldPanic {
				// Should return 500 status code
				assert.Equal(t, http.StatusInternalServerError, w.Code)
				// Should contain error details
				assert.Contains(t, w.Body.String(), "internal_error")
				assert.Contains(t, w.Body.String(), "Internal server error")
				// Logger should have been called
				assert.True(t, logger.errorCalled, "Error logger should have been called")
				assert.Contains(t, logger.errorArgs, "panic recovered:")
				assert.Contains(t, logger.errorArgs, "test panic")
			} else {
				// Should return 200 status code
				assert.Equal(t, http.StatusOK, w.Code)
				// Logger should not have been called
				assert.False(t, logger.errorCalled, "Error logger should not have been called")
			}
		})
	}
}
