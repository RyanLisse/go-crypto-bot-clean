package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLoggingMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		path         string
		method       string
		status       int
		addRequestID bool
		addGinError  bool
	}{
		{
			name:         "successful request",
			path:         "/api/test",
			method:       http.MethodGet,
			status:       http.StatusOK,
			addRequestID: true,
			addGinError:  false,
		},
		{
			name:         "error request",
			path:         "/api/error",
			method:       http.MethodPost,
			status:       http.StatusInternalServerError,
			addRequestID: true,
			addGinError:  true,
		},
		{
			name:         "without request ID",
			path:         "/api/test",
			method:       http.MethodGet,
			status:       http.StatusOK,
			addRequestID: false,
			addGinError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock logger
			logger := &mockLogger{}

			// Create a test router with the logging middleware
			router := gin.New()
			router.Use(LoggingMiddleware(logger))

			// Define a test handler
			router.Handle(tt.method, tt.path, func(c *gin.Context) {
				if tt.addGinError {
					c.Error(errors.New("test error"))
				}
				c.Status(tt.status)
			})

			// Make a request
			req := httptest.NewRequest(tt.method, tt.path, nil)
			if tt.addRequestID {
				req.Header.Set("X-Request-ID", "test-request-id")
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Check the response status
			assert.Equal(t, tt.status, w.Code)

			// Verify logger was called correctly
			if tt.addGinError {
				assert.True(t, logger.errorCalled, "Error logger should have been called")
				assert.Contains(t, logger.errorArgs, "errors")
				// Gin formats errors with a prefix like "Error #01: " and a newline
				errorFound := false
				for _, arg := range logger.errorArgs {
					if s, ok := arg.(string); ok && s != "errors" {
						if assert.Contains(t, s, "test error") {
							errorFound = true
							break
						}
					}
				}
				assert.True(t, errorFound, "Error message should contain 'test error'")
			} else {
				assert.True(t, logger.infoCalled, "Info logger should have been called")
			}

			// Check logged fields
			var logArgs []interface{}
			if tt.addGinError {
				logArgs = logger.errorArgs
			} else {
				logArgs = logger.infoArgs
			}

			// Find method in logged fields
			methodFound := false
			pathFound := false
			statusFound := false
			requestIDFound := false

			for i := 0; i < len(logArgs); i += 2 {
				if i+1 < len(logArgs) {
					if logArgs[i] == "method" {
						methodFound = true
						assert.Equal(t, tt.method, logArgs[i+1])
					}
					if logArgs[i] == "path" {
						pathFound = true
						assert.Equal(t, tt.path, logArgs[i+1])
					}
					if logArgs[i] == "status" {
						statusFound = true
						assert.Equal(t, tt.status, logArgs[i+1])
					}
					if logArgs[i] == "request_id" {
						requestIDFound = true
						if tt.addRequestID {
							assert.Equal(t, "test-request-id", logArgs[i+1])
						} else {
							assert.Equal(t, "", logArgs[i+1])
						}
					}
				}
			}

			assert.True(t, methodFound, "Method should be logged")
			assert.True(t, pathFound, "Path should be logged")
			assert.True(t, statusFound, "Status should be logged")
			assert.True(t, requestIDFound, "Request ID should be logged")
		})
	}
}
