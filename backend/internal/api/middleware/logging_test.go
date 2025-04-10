package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggingMiddleware(t *testing.T) {
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

			// Define a test handler
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.addGinError {
					// Simulate error by writing 500 status
					http.Error(w, "test error", tt.status)
				} else {
					w.WriteHeader(tt.status)
				}
			})

			// Wrap with logging middleware
			handler := LoggingMiddleware(logger)(testHandler)

			// Make a request
			req := httptest.NewRequest(tt.method, tt.path, nil)
			if tt.addRequestID {
				req.Header.Set("X-Request-ID", "test-request-id")
			}
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			// Check the response status
			assert.Equal(t, tt.status, w.Code)

			// Verify logger was called correctly
			if tt.addGinError {
				assert.True(t, logger.errorCalled, "Error logger should have been called")
				assert.Contains(t, logger.errorArgs, "errors")
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
