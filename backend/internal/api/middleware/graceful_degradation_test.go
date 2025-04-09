package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestGracefulDegradation(t *testing.T) {
	// Setup test logger
	logger := zaptest.NewLogger(t)

	// Create graceful degradation with test configuration
	config := DefaultGracefulDegradationConfig()
	config.ReadOnlyThreshold = 3
	config.MaintenanceThreshold = 6
	config.RecoveryInterval = 100 * time.Millisecond
	config.Logger = logger
	gd := NewGracefulDegradation(config)

	// Create graceful degradation middleware
	middleware := GracefulDegradationMiddleware(gd)

	t.Run("should allow all requests in normal mode", func(t *testing.T) {
		// Create a test handler
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success"))
		})

		// Wrap the handler with graceful degradation middleware
		wrappedHandler := middleware(handler)

		// Test GET request
		req := httptest.NewRequest("GET", "/test", nil)
		rec := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		// Test POST request
		req = httptest.NewRequest("POST", "/test", nil)
		rec = httptest.NewRecorder()
		wrappedHandler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		// Test PUT request
		req = httptest.NewRequest("PUT", "/test", nil)
		rec = httptest.NewRecorder()
		wrappedHandler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("should degrade to read-only mode after threshold errors", func(t *testing.T) {
		// Reset the graceful degradation
		gd.SetMode(NormalMode)
		gd.ResetErrorCount()

		// Record errors to reach read-only threshold
		for i := 0; i < config.ReadOnlyThreshold; i++ {
			gd.RecordError()
		}

		// Should now be in read-only mode
		assert.Equal(t, ReadOnlyMode, gd.GetMode())

		// Create a test handler
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success"))
		})

		// Wrap the handler with graceful degradation middleware
		wrappedHandler := middleware(handler)

		// Test GET request (should be allowed)
		req := httptest.NewRequest("GET", "/test", nil)
		rec := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		// Test POST request (should be rejected)
		req = httptest.NewRequest("POST", "/test", nil)
		rec = httptest.NewRecorder()
		wrappedHandler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
	})

	t.Run("should degrade to maintenance mode after threshold errors", func(t *testing.T) {
		// Reset the graceful degradation
		gd.SetMode(NormalMode)
		gd.ResetErrorCount()

		// Record errors to reach maintenance threshold
		for i := 0; i < config.MaintenanceThreshold; i++ {
			gd.RecordError()
		}

		// Should now be in maintenance mode
		assert.Equal(t, MaintenanceMode, gd.GetMode())

		// Create a test handler
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success"))
		})

		// Wrap the handler with graceful degradation middleware
		wrappedHandler := middleware(handler)

		// Test GET request (should be rejected)
		req := httptest.NewRequest("GET", "/test", nil)
		rec := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusServiceUnavailable, rec.Code)

		// Test critical path (should be allowed)
		req = httptest.NewRequest("GET", "/api/v1/health", nil)
		rec = httptest.NewRecorder()
		wrappedHandler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("should recover after recovery interval", func(t *testing.T) {
		// Set to maintenance mode
		gd.SetMode(MaintenanceMode)

		// Wait for recovery interval
		time.Sleep(config.RecoveryInterval + 10*time.Millisecond)

		// Create a test handler
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success"))
		})

		// Wrap the handler with graceful degradation middleware
		wrappedHandler := middleware(handler)

		// Make a request to trigger recovery check
		req := httptest.NewRequest("GET", "/test", nil)
		rec := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(rec, req)

		// Should now be in read-only mode
		assert.Equal(t, ReadOnlyMode, gd.GetMode())

		// Wait for another recovery interval
		time.Sleep(config.RecoveryInterval + 10*time.Millisecond)

		// Make another request to trigger recovery check
		req = httptest.NewRequest("GET", "/test", nil)
		rec = httptest.NewRecorder()
		wrappedHandler.ServeHTTP(rec, req)

		// Should now be in normal mode
		assert.Equal(t, NormalMode, gd.GetMode())
	})
}
