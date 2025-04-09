package middleware

import (
	"net/http"
	"sync"
	"time"

	"go-crypto-bot-clean/backend/internal/auth"
	"go-crypto-bot-clean/backend/internal/logging"

	"go.uber.org/zap"
)

// CircuitBreakerState represents the state of a circuit breaker
type CircuitBreakerState int

const (
	// CircuitClosed means the circuit is closed and requests flow normally
	CircuitClosed CircuitBreakerState = iota
	// CircuitOpen means the circuit is open and requests are rejected
	CircuitOpen
	// CircuitHalfOpen means the circuit is testing if it can be closed
	CircuitHalfOpen
)

// CircuitBreakerConfig holds configuration for the circuit breaker
type CircuitBreakerConfig struct {
	// FailureThreshold is the number of failures before opening the circuit
	FailureThreshold int
	// ResetTimeout is the time to wait before trying to close the circuit
	ResetTimeout time.Duration
	// MaxConcurrent is the maximum number of concurrent requests
	MaxConcurrent int
	// RequestTimeout is the timeout for requests
	RequestTimeout time.Duration
	// Logger is the logger to use
	Logger *zap.Logger
}

// DefaultCircuitBreakerConfig returns the default circuit breaker configuration
func DefaultCircuitBreakerConfig() CircuitBreakerConfig {
	return CircuitBreakerConfig{
		FailureThreshold: 5,
		ResetTimeout:     10 * time.Second,
		MaxConcurrent:    100,
		RequestTimeout:   5 * time.Second,
		Logger:           nil,
	}
}

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	config         CircuitBreakerConfig
	state          CircuitBreakerState
	failures       int
	lastFailure    time.Time
	mutex          sync.RWMutex
	activeSessions int
	logger         *zap.Logger
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(config CircuitBreakerConfig) *CircuitBreaker {
	if config.Logger == nil {
		config.Logger = logging.Logger.Logger
	}

	return &CircuitBreaker{
		config:         config,
		state:          CircuitClosed,
		failures:       0,
		lastFailure:    time.Time{},
		activeSessions: 0,
		logger:         config.Logger,
	}
}

// CircuitBreakerMiddleware adds circuit breaker functionality to a handler
func CircuitBreakerMiddleware(cb *CircuitBreaker) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if the circuit is open
			if !cb.AllowRequest() {
				// Circuit is open, reject the request
				cb.logger.Warn("Circuit breaker is open, rejecting request",
					zap.String("path", r.URL.Path),
					zap.String("method", r.Method),
					zap.String("remote_addr", r.RemoteAddr),
				)

				// Create error response
				errResp := auth.NewAuthError(
					auth.ErrorTypeServiceUnavailable,
					"Service temporarily unavailable",
					http.StatusServiceUnavailable,
				).WithDetails(map[string]interface{}{
					"circuit_state": "open",
					"retry_after":   int(cb.config.ResetTimeout.Seconds()),
				}).WithHelp("Please try again later").
					WithRequestID(GetRequestID(r.Context()))

				// Set response headers
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Retry-After", time.Now().Add(cb.config.ResetTimeout).Format(time.RFC1123))
				w.WriteHeader(http.StatusServiceUnavailable)

				// Write error response
				errResp.WriteJSON(w)
				return
			}

			// Increment active sessions
			cb.IncrementSessions()
			defer cb.DecrementSessions()

			// Create a channel to signal completion
			done := make(chan bool, 1)

			// Create a timeout context
			timer := time.NewTimer(cb.config.RequestTimeout)
			defer timer.Stop()

			// Handle the request in a goroutine
			go func() {
				defer func() {
					if err := recover(); err != nil {
						// Record failure
						cb.RecordFailure()

						// Re-panic to let the recovery middleware handle it
						panic(err)
					}
				}()

				// Call the next handler
				next.ServeHTTP(w, r)

				// Signal completion
				done <- true
			}()

			// Wait for completion or timeout
			select {
			case <-done:
				// Request completed successfully
				cb.RecordSuccess()
			case <-timer.C:
				// Request timed out
				cb.RecordFailure()

				// Create error response
				errResp := auth.NewAuthError(
					auth.ErrorTypeTimeout,
					"Request timed out",
					http.StatusGatewayTimeout,
				).WithDetails(map[string]interface{}{
					"timeout": cb.config.RequestTimeout.String(),
				}).WithHelp("Please try again later").
					WithRequestID(GetRequestID(r.Context()))

				// Set response headers
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusGatewayTimeout)

				// Write error response
				errResp.WriteJSON(w)
			}
		})
	}
}

// AllowRequest checks if a request should be allowed
func (cb *CircuitBreaker) AllowRequest() bool {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()

	switch cb.state {
	case CircuitClosed:
		return true
	case CircuitHalfOpen:
		// In half-open state, allow a limited number of requests
		return cb.activeSessions < 1
	case CircuitOpen:
		// Check if enough time has passed to try half-open
		if time.Since(cb.lastFailure) > cb.config.ResetTimeout {
			// Transition to half-open
			cb.mutex.RUnlock()
			cb.mutex.Lock()
			cb.state = CircuitHalfOpen
			cb.mutex.Unlock()
			cb.mutex.RLock()
			return true
		}
		return false
	default:
		return false
	}
}

// RecordSuccess records a successful request
func (cb *CircuitBreaker) RecordSuccess() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	// If in half-open state and a request succeeds, close the circuit
	if cb.state == CircuitHalfOpen {
		cb.state = CircuitClosed
		cb.failures = 0
		cb.logger.Info("Circuit breaker closed after successful request")
	}
}

// RecordFailure records a failed request
func (cb *CircuitBreaker) RecordFailure() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.failures++
	cb.lastFailure = time.Now()

	// If in half-open state, any failure opens the circuit
	if cb.state == CircuitHalfOpen {
		cb.state = CircuitOpen
		cb.logger.Warn("Circuit breaker opened after failure in half-open state",
			zap.Int("failures", cb.failures),
		)
		return
	}

	// If in closed state and failures exceed threshold, open the circuit
	if cb.state == CircuitClosed && cb.failures >= cb.config.FailureThreshold {
		cb.state = CircuitOpen
		cb.logger.Warn("Circuit breaker opened after exceeding failure threshold",
			zap.Int("failures", cb.failures),
			zap.Int("threshold", cb.config.FailureThreshold),
		)
	}
}

// IncrementSessions increments the active sessions counter
func (cb *CircuitBreaker) IncrementSessions() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	cb.activeSessions++
}

// DecrementSessions decrements the active sessions counter
func (cb *CircuitBreaker) DecrementSessions() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	cb.activeSessions--
	if cb.activeSessions < 0 {
		cb.activeSessions = 0
	}
}

// GetState returns the current state of the circuit breaker
func (cb *CircuitBreaker) GetState() CircuitBreakerState {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.state
}

// GetFailures returns the current failure count
func (cb *CircuitBreaker) GetFailures() int {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.failures
}

// Reset resets the circuit breaker to closed state
func (cb *CircuitBreaker) Reset() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	cb.state = CircuitClosed
	cb.failures = 0
	cb.lastFailure = time.Time{}
	cb.logger.Info("Circuit breaker manually reset to closed state")
}
