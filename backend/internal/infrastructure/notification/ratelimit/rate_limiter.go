package ratelimit

import (
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

// Strategy defines the rate limiting strategy
type Strategy string

const (
	// FixedWindow is a fixed window rate limiting strategy
	FixedWindow Strategy = "fixed_window"
	// SlidingWindow is a sliding window rate limiting strategy
	SlidingWindow Strategy = "sliding_window"
	// TokenBucket is a token bucket rate limiting strategy
	TokenBucket Strategy = "token_bucket"
)

// RateLimiter defines the interface for rate limiters
type RateLimiter interface {
	// Allow checks if a request is allowed
	Allow(key string) bool
	// Reset resets the rate limiter for a key
	Reset(key string)
}

// Config holds configuration for rate limiters
type Config struct {
	Strategy Strategy
	Limit    int           // Maximum number of requests
	Window   time.Duration // Time window for rate limiting
	Burst    int           // Burst size for token bucket
}

// NewRateLimiter creates a new rate limiter based on the strategy
func NewRateLimiter(config Config, logger *zap.Logger) (RateLimiter, error) {
	if logger == nil {
		logger, _ = zap.NewProduction()
	}

	switch config.Strategy {
	case FixedWindow:
		return NewFixedWindowRateLimiter(config.Limit, config.Window, logger), nil
	case SlidingWindow:
		return NewSlidingWindowRateLimiter(config.Limit, config.Window, logger), nil
	case TokenBucket:
		return NewTokenBucketRateLimiter(config.Limit, config.Burst, config.Window, logger), nil
	default:
		return nil, fmt.Errorf("unknown rate limiting strategy: %s", config.Strategy)
	}
}

// FixedWindowRateLimiter implements a fixed window rate limiter
type FixedWindowRateLimiter struct {
	limit    int
	window   time.Duration
	counters map[string]*fixedWindowCounter
	mu       sync.Mutex
	logger   *zap.Logger
}

type fixedWindowCounter struct {
	count     int
	startTime time.Time
}

// NewFixedWindowRateLimiter creates a new fixed window rate limiter
func NewFixedWindowRateLimiter(limit int, window time.Duration, logger *zap.Logger) *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{
		limit:    limit,
		window:   window,
		counters: make(map[string]*fixedWindowCounter),
		logger:   logger,
	}
}

// Allow checks if a request is allowed
func (r *FixedWindowRateLimiter) Allow(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	counter, ok := r.counters[key]
	if !ok || now.Sub(counter.startTime) >= r.window {
		// New window
		r.counters[key] = &fixedWindowCounter{
			count:     1,
			startTime: now,
		}
		return true
	}

	// Existing window
	if counter.count < r.limit {
		counter.count++
		return true
	}

	r.logger.Debug("Rate limit exceeded",
		zap.String("key", key),
		zap.Int("limit", r.limit),
		zap.Duration("window", r.window),
		zap.Int("count", counter.count),
	)
	return false
}

// Reset resets the rate limiter for a key
func (r *FixedWindowRateLimiter) Reset(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.counters, key)
}

// SlidingWindowRateLimiter implements a sliding window rate limiter
type SlidingWindowRateLimiter struct {
	limit    int
	window   time.Duration
	requests map[string][]time.Time
	mu       sync.Mutex
	logger   *zap.Logger
}

// NewSlidingWindowRateLimiter creates a new sliding window rate limiter
func NewSlidingWindowRateLimiter(limit int, window time.Duration, logger *zap.Logger) *SlidingWindowRateLimiter {
	return &SlidingWindowRateLimiter{
		limit:    limit,
		window:   window,
		requests: make(map[string][]time.Time),
		logger:   logger,
	}
}

// Allow checks if a request is allowed
func (r *SlidingWindowRateLimiter) Allow(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-r.window)

	// Get requests for this key
	times, ok := r.requests[key]
	if !ok {
		r.requests[key] = []time.Time{now}
		return true
	}

	// Filter out requests outside the window
	validTimes := make([]time.Time, 0, len(times))
	for _, t := range times {
		if t.After(windowStart) {
			validTimes = append(validTimes, t)
		}
	}

	// Check if we're under the limit
	if len(validTimes) < r.limit {
		r.requests[key] = append(validTimes, now)
		return true
	}

	r.requests[key] = validTimes
	r.logger.Debug("Rate limit exceeded",
		zap.String("key", key),
		zap.Int("limit", r.limit),
		zap.Duration("window", r.window),
		zap.Int("count", len(validTimes)),
	)
	return false
}

// Reset resets the rate limiter for a key
func (r *SlidingWindowRateLimiter) Reset(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.requests, key)
}

// TokenBucketRateLimiter implements a token bucket rate limiter
type TokenBucketRateLimiter struct {
	rate      float64 // Tokens per second
	burst     int     // Maximum bucket size
	buckets   map[string]*tokenBucket
	mu        sync.Mutex
	logger    *zap.Logger
}

type tokenBucket struct {
	tokens    float64
	lastRefill time.Time
}

// NewTokenBucketRateLimiter creates a new token bucket rate limiter
func NewTokenBucketRateLimiter(rate int, burst int, per time.Duration, logger *zap.Logger) *TokenBucketRateLimiter {
	// Calculate tokens per second
	tokensPerSecond := float64(rate) / per.Seconds()

	return &TokenBucketRateLimiter{
		rate:      tokensPerSecond,
		burst:     burst,
		buckets:   make(map[string]*tokenBucket),
		logger:    logger,
	}
}

// Allow checks if a request is allowed
func (r *TokenBucketRateLimiter) Allow(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	bucket, ok := r.buckets[key]
	if !ok {
		// New bucket, start with full tokens
		r.buckets[key] = &tokenBucket{
			tokens:    float64(r.burst),
			lastRefill: now,
		}
		return true
	}

	// Refill tokens based on time elapsed
	elapsed := now.Sub(bucket.lastRefill).Seconds()
	bucket.tokens = min(float64(r.burst), bucket.tokens+elapsed*r.rate)
	bucket.lastRefill = now

	// Check if we have enough tokens
	if bucket.tokens >= 1.0 {
		bucket.tokens -= 1.0
		return true
	}

	r.logger.Debug("Rate limit exceeded",
		zap.String("key", key),
		zap.Float64("rate", r.rate),
		zap.Int("burst", r.burst),
		zap.Float64("tokens", bucket.tokens),
	)
	return false
}

// Reset resets the rate limiter for a key
func (r *TokenBucketRateLimiter) Reset(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.buckets, key)
}

// min returns the minimum of two float64 values
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
