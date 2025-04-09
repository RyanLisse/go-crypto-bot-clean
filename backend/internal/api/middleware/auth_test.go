package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		apiKey         string
		validKeys      map[string]struct{}
		expectedStatus int
		expectAbort    bool
	}{
		{
			name:           "valid API key",
			apiKey:         "valid-key",
			validKeys:      map[string]struct{}{"valid-key": {}},
			expectedStatus: http.StatusOK,
			expectAbort:    false,
		},
		{
			name:           "invalid API key",
			apiKey:         "invalid-key",
			validKeys:      map[string]struct{}{"valid-key": {}},
			expectedStatus: http.StatusUnauthorized,
			expectAbort:    true,
		},
		{
			name:           "missing API key",
			apiKey:         "",
			validKeys:      map[string]struct{}{"valid-key": {}},
			expectedStatus: http.StatusUnauthorized,
			expectAbort:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test router with our middleware
			router := gin.New()
			router.Use(AuthMiddleware(tt.validKeys))

			// Add a handler after the middleware
			router.GET("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			// Create a test request
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.apiKey != "" {
				req.Header.Set("X-API-Key", tt.apiKey)
			}
			w := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(w, req)

			// Check the response
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectAbort {
				assert.Contains(t, w.Body.String(), "unauthorized")
				assert.Contains(t, w.Body.String(), "Invalid or missing API key")
			} else {
				assert.Empty(t, w.Body.String())
			}
		})
	}
}
