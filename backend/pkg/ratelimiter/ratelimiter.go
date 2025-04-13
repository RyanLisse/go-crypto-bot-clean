package ratelimiter

import (
	"sync"
	"time"
)

// TokenBucket represents a token bucket rate limiter
type TokenBucket struct {
	rate       float64    // tokens per second
	capacity   float64    // bucket capacity
	tokens     float64    // current number of tokens
	lastUpdate time.Time  // last time tokens were added
	mu         sync.Mutex // mutex for concurrent access
}

// NewTokenBucket creates a new token bucket rate limiter
func NewTokenBucket(rate float64, capacity float64) *TokenBucket {
	return &TokenBucket{
		rate:       rate,
		capacity:   capacity,
		tokens:     capacity,
		lastUpdate: time.Now(),
	}
}

// Allow returns true if a request is allowed, and consumes a token
func (t *TokenBucket) Allow() bool {
	return t.AllowN(1)
}

// AllowN returns true if n requests are allowed, and consumes n tokens
func (t *TokenBucket) AllowN(n float64) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(t.lastUpdate).Seconds()
	t.lastUpdate = now

	// Add tokens based on elapsed time
	t.tokens += elapsed * t.rate
	if t.tokens > t.capacity {
		t.tokens = t.capacity
	}

	// Check if we have enough tokens
	if t.tokens < n {
		return false
	}

	// Consume tokens
	t.tokens -= n
	return true
}

// WaitN waits until n tokens are available and then consumes them
func (t *TokenBucket) WaitN(n float64) {
	for {
		t.mu.Lock()
		now := time.Now()
		elapsed := now.Sub(t.lastUpdate).Seconds()
		t.lastUpdate = now

		// Add tokens based on elapsed time
		t.tokens += elapsed * t.rate
		if t.tokens > t.capacity {
			t.tokens = t.capacity
		}

		// If we have enough tokens, consume them and return
		if t.tokens >= n {
			t.tokens -= n
			t.mu.Unlock()
			return
		}

		// Calculate how long to wait for the required tokens
		waitTime := time.Duration((n - t.tokens) / t.rate * float64(time.Second))
		t.mu.Unlock()

		// Wait for tokens to replenish
		time.Sleep(waitTime)
	}
}

// Wait waits until a token is available and then consumes it
func (t *TokenBucket) Wait() {
	t.WaitN(1)
}

// GetTokens returns the current number of tokens in the bucket
func (t *TokenBucket) GetTokens() float64 {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(t.lastUpdate).Seconds()

	// Add tokens based on elapsed time
	tokens := t.tokens + elapsed*t.rate
	if tokens > t.capacity {
		tokens = t.capacity
	}

	return tokens
}

// RateLimiterMap manages multiple rate limiters by key
type RateLimiterMap struct {
	limiters map[string]*TokenBucket
	rate     float64
	capacity float64
	mu       sync.Mutex
}

// NewRateLimiterMap creates a new rate limiter map
func NewRateLimiterMap(rate float64, capacity float64) *RateLimiterMap {
	return &RateLimiterMap{
		limiters: make(map[string]*TokenBucket),
		rate:     rate,
		capacity: capacity,
	}
}

// Get returns a rate limiter for the given key
func (r *RateLimiterMap) Get(key string) *TokenBucket {
	r.mu.Lock()
	defer r.mu.Unlock()

	limiter, exists := r.limiters[key]
	if !exists {
		limiter = NewTokenBucket(r.rate, r.capacity)
		r.limiters[key] = limiter
	}

	return limiter
}

// Allow returns true if a request is allowed for the given key
func (r *RateLimiterMap) Allow(key string) bool {
	return r.Get(key).Allow()
}

// Wait waits until a token is available for the given key
func (r *RateLimiterMap) Wait(key string) {
	r.Get(key).Wait()
}
