package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestCircuitBreaker(t *testing.T) {
	// Setup test logger
	logger := zaptest.NewLogger(t)

	// Create circuit breaker with test configuration
	config := DefaultCircuitBreakerConfig()
	config.FailureThreshold = 3
	config.ResetTimeout = 100 * time.Millisecond
	config.Logger = logger
	cb := NewCircuitBreaker(config)

	// Create circuit breaker middleware
	middleware := CircuitBreakerMiddleware(cb)

	t.Run("should allow requests when circuit is closed", func(t *testing.T) {
		// Create a test handler
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success"))
		})

		// Wrap the handler with circuit breaker middleware
		wrappedHandler := middleware(handler)

		// Create a test request
		req := httptest.NewRequest("GET", "/test", nil)
		rec := httptest.NewRecorder()

		// Call the handler
		wrappedHandler.ServeHTTP(rec, req)

		// Check response
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "success", rec.Body.String())
	})

	t.Run("should open circuit after failure threshold", func(t *testing.T) {
		// Create a handler that will panic
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("test panic")
		})

		// Wrap with recovery middleware first, then circuit breaker
		recoveryMiddleware := RecoveryMiddleware(DefaultRecoveryOptions())
		wrappedHandler := middleware(recoveryMiddleware(handler))

		// Trigger failures to open the circuit
		for i := 0; i < config.FailureThreshold; i++ {
			req := httptest.NewRequest("GET", "/test", nil)
			rec := httptest.NewRecorder()
			wrappedHandler.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}

		// Circuit should now be open
		assert.Equal(t, CircuitOpen, cb.GetState())

		// Next request should be rejected with 503
		req := httptest.NewRequest("GET", "/test", nil)
		rec := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
	})

	t.Run("should transition to half-open after reset timeout", func(t *testing.T) {
		// Reset the circuit breaker
		cb.Reset()

		// Create a handler that will panic
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("test panic")
		})

		// Wrap with recovery middleware first, then circuit breaker
		recoveryMiddleware := RecoveryMiddleware(DefaultRecoveryOptions())
		wrappedHandler := middleware(recoveryMiddleware(handler))

		// Trigger failures to open the circuit
		for i := 0; i < config.FailureThreshold; i++ {
			req := httptest.NewRequest("GET", "/test", nil)
			rec := httptest.NewRecorder()
			wrappedHandler.ServeHTTP(rec, req)
		}

		// Circuit should now be open
		assert.Equal(t, CircuitOpen, cb.GetState())

		// Wait for reset timeout
		time.Sleep(config.ResetTimeout + 10*time.Millisecond)

		// Next request should be allowed (half-open state)
		req := httptest.NewRequest("GET", "/test", nil)
		rec := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(rec, req)
		
		// Should still return error but circuit should be in half-open state
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, CircuitOpen, cb.GetState()) // Will be open again after the failure
	})

	t.Run("should close circuit after successful request in half-open state", func(t *testing.T) {
		// Reset the circuit breaker
		cb.Reset()

		// Manually set to half-open state
		cb.RecordFailure()
		cb.RecordFailure()
		cb.RecordFailure()
		assert.Equal(t, CircuitOpen, cb.GetState())
		
		// Wait for reset timeout
		time.Sleep(config.ResetTimeout + 10*time.Millisecond)

		// Create a successful handler
		successHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success"))
		})

		// Wrap with circuit breaker
		wrappedSuccessHandler := middleware(successHandler)

		// Make a successful request
		req := httptest.NewRequest("GET", "/test", nil)
		rec := httptest.NewRecorder()
		wrappedSuccessHandler.ServeHTTP(rec, req)

		// Should return success and circuit should be closed
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, CircuitClosed, cb.GetState())
	})
}
