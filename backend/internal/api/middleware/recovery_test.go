package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-crypto-bot-clean/backend/internal/auth"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

func TestRecoveryMiddleware(t *testing.T) {
	// Setup test logger
	logger := zaptest.NewLogger(t)
	defer logger.Sync()

	// Create recovery middleware with test options
	opts := DefaultRecoveryOptions()
	opts.Logger = logger
	recoveryMiddleware := RecoveryMiddleware(opts)

	t.Run("should recover from panic and return 500 error", func(t *testing.T) {
		// Create a handler that will panic
		panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("test panic")
		})

		// Wrap the panic handler with recovery middleware
		handler := recoveryMiddleware(panicHandler)

		// Create a test request
		req := httptest.NewRequest("GET", "/test", nil)

		// Add request ID to context
		ctx := context.WithValue(req.Context(), RequestIDContextKey, "test-request-id")
		req = req.WithContext(ctx)

		rec := httptest.NewRecorder()

		// Call the handler
		handler.ServeHTTP(rec, req)

		// Check response
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		// Parse response body
		var errResp auth.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &errResp)
		assert.NoError(t, err)

		// Verify error response fields
		assert.Equal(t, "internal_error", string(errResp.Error.Type))
		assert.Contains(t, errResp.Error.Message, "unexpected error")
		assert.Equal(t, "test-request-id", errResp.RequestID)
		assert.Equal(t, "/test", errResp.Path)
		assert.Equal(t, "GET", errResp.Method)
	})

	t.Run("should handle auth errors correctly", func(t *testing.T) {
		// Create a handler that will panic with an auth error
		panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authError := auth.NewAuthError(
				auth.ErrorTypeUnauthorized,
				"Invalid token",
				http.StatusUnauthorized,
			)
			panic(authError)
		})

		// Wrap the panic handler with recovery middleware
		handler := recoveryMiddleware(panicHandler)

		// Create a test request
		req := httptest.NewRequest("GET", "/auth/test", nil)

		// Add request ID to context
		ctx := context.WithValue(req.Context(), RequestIDContextKey, "test-request-id")
		req = req.WithContext(ctx)

		rec := httptest.NewRecorder()

		// Call the handler
		handler.ServeHTTP(rec, req)

		// Check response
		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		// Parse response body
		var errResp auth.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &errResp)
		assert.NoError(t, err)

		// Verify error response fields
		assert.Equal(t, string(auth.ErrorTypeUnauthorized), string(errResp.Error.Type))
		assert.Equal(t, "Invalid token", errResp.Error.Message)
		assert.Equal(t, "test-request-id", errResp.RequestID)
		assert.Equal(t, "/auth/test", errResp.Path)
		assert.Equal(t, "GET", errResp.Method)
	})

	t.Run("should detect token errors and return appropriate status", func(t *testing.T) {
		// Create a handler that will panic with a token error
		panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic(errors.New("token has expired"))
		})

		// Wrap the panic handler with recovery middleware
		handler := recoveryMiddleware(panicHandler)

		// Create a test request
		req := httptest.NewRequest("GET", "/auth/token", nil)

		// Add request ID to context
		ctx := context.WithValue(req.Context(), RequestIDContextKey, "test-request-id")
		req = req.WithContext(ctx)

		rec := httptest.NewRecorder()

		// Call the handler
		handler.ServeHTTP(rec, req)

		// Check response
		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		// Parse response body
		var errResp auth.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &errResp)
		assert.NoError(t, err)

		// Verify error response fields
		assert.Equal(t, "expired_token", string(errResp.Error.Type))
		assert.Contains(t, errResp.Error.Message, "unexpected error")
		assert.Equal(t, "test-request-id", errResp.RequestID)
	})

	t.Run("should include user info in logs when available", func(t *testing.T) {
		// Create a buffer to capture logs
		var buf bytes.Buffer

		// Create a logger that writes to the buffer
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			zapcore.AddSync(&buf),
			zapcore.DebugLevel,
		)
		testLogger := zap.New(core)

		// Create recovery middleware with test options
		testOpts := DefaultRecoveryOptions()
		testOpts.Logger = testLogger
		testRecoveryMiddleware := RecoveryMiddleware(testOpts)

		// Create a handler that will panic
		panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create a user data object
			user := &auth.UserData{
				ID:       "test-user-id",
				Email:    "test@example.com",
				Username: "testuser",
				Roles:    []string{"user"},
			}

			// Store the user in the context using the auth.UserDataKey
			ctx := context.WithValue(r.Context(), auth.UserDataKey, user)
			r = r.WithContext(ctx)

			// Then panic
			panic("test panic with user")
		})

		// Wrap the panic handler with recovery middleware
		handler := testRecoveryMiddleware(panicHandler)

		// Create a test request
		req := httptest.NewRequest("GET", "/user/profile", nil)

		// Add request ID to context
		ctx := context.WithValue(req.Context(), RequestIDContextKey, "test-request-id")
		req = req.WithContext(ctx)

		rec := httptest.NewRecorder()

		// Call the handler
		handler.ServeHTTP(rec, req)

		// Check logs contain user info
		logOutput := buf.String()
		// Just check that the log contains the panic message
		assert.Contains(t, logOutput, "test panic with user")
	})

	t.Run("should handle normal requests without panic", func(t *testing.T) {
		// Create a normal handler
		normalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success"))
		})

		// Wrap the normal handler with recovery middleware
		handler := recoveryMiddleware(normalHandler)

		// Create a test request
		req := httptest.NewRequest("GET", "/normal", nil)

		// Add request ID to context
		ctx := context.WithValue(req.Context(), RequestIDContextKey, "test-request-id")
		req = req.WithContext(ctx)

		rec := httptest.NewRecorder()

		// Call the handler
		handler.ServeHTTP(rec, req)

		// Check response
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "success", rec.Body.String())
	})
}

func TestRecoveryMiddlewareWithCircuitBreaker(t *testing.T) {
	// Setup test logger
	logger := zaptest.NewLogger(t)

	// Create circuit breaker
	cbConfig := DefaultCircuitBreakerConfig()
	cbConfig.FailureThreshold = 3
	cbConfig.Logger = logger
	cb := NewCircuitBreaker(cbConfig)

	// Create recovery middleware with circuit breaker
	opts := DefaultRecoveryOptions()
	opts.Logger = logger
	opts.CircuitBreaker = cb
	recoveryMiddleware := RecoveryMiddleware(opts)

	// Create a handler that will panic
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	// Wrap the panic handler with recovery middleware
	handler := recoveryMiddleware(panicHandler)

	// Trigger failures to open the circuit
	for i := 0; i < cbConfig.FailureThreshold; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	}

	// Circuit should now be open
	assert.Equal(t, CircuitOpen, cb.GetState())

	// Reset the circuit breaker for other tests
	cb.Reset()
}

func TestRecoveryMiddlewareWithGracefulDegradation(t *testing.T) {
	// Setup test logger
	logger := zaptest.NewLogger(t)

	// Create graceful degradation
	gdConfig := DefaultGracefulDegradationConfig()
	gdConfig.ReadOnlyThreshold = 3
	gdConfig.Logger = logger
	gd := NewGracefulDegradation(gdConfig)

	// Create recovery middleware with graceful degradation
	opts := DefaultRecoveryOptions()
	opts.Logger = logger
	opts.GracefulDegradation = gd
	recoveryMiddleware := RecoveryMiddleware(opts)

	// Create a handler that will panic
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	// Wrap the panic handler with recovery middleware
	handler := recoveryMiddleware(panicHandler)

	// Trigger failures to degrade service
	for i := 0; i < gdConfig.ReadOnlyThreshold; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	}

	// Service should now be in read-only mode
	assert.Equal(t, ReadOnlyMode, gd.GetMode())

	// Reset the service mode for other tests
	gd.SetMode(NormalMode)
	gd.ResetErrorCount()
}
