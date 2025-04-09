package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// RateLimiter implements a token bucket rate limiter
type RateLimiter struct {
	tokens     float64
	capacity   float64
	rate       float64
	lastUpdate time.Time
	mu         sync.Mutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate float64, capacity float64) *RateLimiter {
	return &RateLimiter{
		tokens:     capacity,
		capacity:   capacity,
		rate:       rate,
		lastUpdate: time.Now(),
	}
}

// Allow checks if a request is allowed
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Update tokens based on time elapsed
	now := time.Now()
	elapsed := now.Sub(rl.lastUpdate).Seconds()
	rl.tokens = min(rl.capacity, rl.tokens+elapsed*rl.rate)
	rl.lastUpdate = now

	// Check if we have enough tokens
	if rl.tokens >= 1.0 {
		rl.tokens -= 1.0
		return true
	}

	return false
}

// RateLimitConfig contains configuration for rate limiting
type RateLimitConfig struct {
	// EnableRateLimit enables rate limiting
	EnableRateLimit bool

	// RequestsPerSecond is the number of requests allowed per second
	RequestsPerSecond float64

	// BurstSize is the maximum number of requests allowed in a burst
	BurstSize float64

	// Logger is the logger to use
	Logger *zap.Logger
}

// DefaultRateLimitConfig returns the default rate limit configuration
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		EnableRateLimit:    true,
		RequestsPerSecond:  10.0,
		BurstSize:          20.0,
		Logger:             zap.NewNop(),
	}
}

// RateLimitMiddleware returns a middleware that limits request rate
func RateLimitMiddleware(config RateLimitConfig) func(http.Handler) http.Handler {
	// Create a map of rate limiters by IP
	limiters := make(map[string]*RateLimiter)
	var limitersMu sync.Mutex

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !config.EnableRateLimit {
				next.ServeHTTP(w, r)
				return
			}

			// Get client IP
			ip := r.RemoteAddr
			if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
				ip = forwardedFor
			}

			// Get or create rate limiter for this IP
			limitersMu.Lock()
			limiter, ok := limiters[ip]
			if !ok {
				limiter = NewRateLimiter(config.RequestsPerSecond, config.BurstSize)
				limiters[ip] = limiter
			}
			limitersMu.Unlock()

			// Check if request is allowed
			if !limiter.Allow() {
				config.Logger.Warn("Rate limit exceeded",
					zap.String("ip", ip),
					zap.String("path", r.URL.Path),
				)
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}

// RegisterRateLimitMiddleware registers rate limit middleware with a Chi router
func RegisterRateLimitMiddleware(r chi.Router, logger *zap.Logger) {
	config := DefaultRateLimitConfig()
	config.Logger = logger

	// Add rate limit middleware
	r.Use(RateLimitMiddleware(config))
}
