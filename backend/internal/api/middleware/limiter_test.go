package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRateLimiterMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		rate           float64
		capacity       float64
		callCount      int
		expectedStatus []int
		identifier     string
	}{
		{
			name:           "allows requests under limit",
			rate:           10,
			capacity:       10,
			callCount:      5,
			expectedStatus: []int{200, 200, 200, 200, 200},
			identifier:     "test-client",
		},
		{
			name:           "rate limits when exceeded",
			rate:           2,
			capacity:       2,
			callCount:      5,
			expectedStatus: []int{200, 200, 429, 429, 429},
			identifier:     "test-client-2",
		},
		{
			name:           "handles empty identifier",
			rate:           1,
			capacity:       1,
			callCount:      2,
			expectedStatus: []int{200, 429},
			identifier:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock logger
			logger := &mockLogger{}

			// Create an identifier extractor function
			extractor := func(r *http.Request) string {
				return tt.identifier
			}

			// Create the middleware
			middleware := RateLimiterMiddleware(tt.rate, tt.capacity, extractor, logger)

			// Define a test handler
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Wrap with middleware
			handler := middleware(testHandler)

			// Make the specified number of requests
			for i := 0; i < tt.callCount; i++ {
				req := httptest.NewRequest(http.MethodGet, "/test", nil)
				w := httptest.NewRecorder()
				handler.ServeHTTP(w, req)

				// Check the response status
				assert.Equal(t, tt.expectedStatus[i], w.Code, "Request %d should have status %d, got %d", i+1, tt.expectedStatus[i], w.Code)

				// If rate limited, check the response body
				if tt.expectedStatus[i] == http.StatusTooManyRequests {
					assert.Contains(t, w.Body.String(), "rate_limited")
					assert.Contains(t, w.Body.String(), "Too many requests")
				}
			}

			// Verify logger was called as expected
			if tt.name == "rate limits when exceeded" || tt.name == "handles empty identifier" {
				assert.True(t, logger.errorCalled, "Error should have been called")
			}
		})
	}
}

func TestRateLimiterMiddleware_TokenRefill(t *testing.T) {
	// Create a mock logger
	logger := &mockLogger{}

	// Configure a rate limiter with 1 token per second
	rate := 1.0
	capacity := 1.0
	extractor := func(r *http.Request) string {
		return "test-client"
	}

	// Create the middleware
	middleware := RateLimiterMiddleware(rate, capacity, extractor, logger)

	// Define a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Wrap with middleware
	handler := middleware(testHandler)

	// First request should succeed
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Second immediate request should be rate limited
	req = httptest.NewRequest(http.MethodGet, "/test", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)

	// Wait for token to refill
	time.Sleep(1100 * time.Millisecond)

	// Third request after waiting should succeed
	req = httptest.NewRequest(http.MethodGet, "/test", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify logger was called
	assert.True(t, logger.errorCalled, "Error should have been called")
}
