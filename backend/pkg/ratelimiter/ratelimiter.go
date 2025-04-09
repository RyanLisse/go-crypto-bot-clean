package ratelimiter

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"
)

// TokenBucketRateLimiter implements a token bucket rate limiting algorithm
type TokenBucketRateLimiter struct {
	// Maximum number of tokens in the bucket
	capacity float64
	// Current number of tokens
	tokens float64
	// Rate of token replenishment per second
	rate float64
	// Last time tokens were updated
	lastUpdate time.Time
	// Mutex to ensure thread-safety
	mu sync.Mutex
}

// NewTokenBucketRateLimiter creates a new rate limiter
func NewTokenBucketRateLimiter(rate float64, capacity float64) *TokenBucketRateLimiter {
	return &TokenBucketRateLimiter{
		capacity:   capacity,
		tokens:     capacity,
		rate:       rate,
		lastUpdate: time.Now(),
	}
}

// updateTokens updates the number of tokens based on elapsed time
func (rl *TokenBucketRateLimiter) updateTokens() {
	now := time.Now()
	elapsed := now.Sub(rl.lastUpdate).Seconds()
	rl.tokens += rl.rate * elapsed

	// Cap tokens at capacity
	if rl.tokens > rl.capacity {
		rl.tokens = rl.capacity
	}

	rl.lastUpdate = now
}

// Wait blocks until a token is available
func (rl *TokenBucketRateLimiter) Wait(ctx context.Context) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Update tokens
	rl.updateTokens()

	// If no tokens available, wait
	for rl.tokens < 1 {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Calculate wait time - handle division by zero and very small rates
			var waitTime time.Duration
			if rl.rate <= 1e-9 { // Consider rates smaller than 1 token per ~31 years as zero
				return fmt.Errorf("rate limit exceeded: effective rate is zero (%v tokens/sec)", rl.rate)
			} else {
				needed := 1.0 - rl.tokens
				// Calculate wait time in seconds, ensure it's positive
				secondsToWait := math.Max(0, needed/rl.rate)
				waitTime = time.Duration(secondsToWait * float64(time.Second))
				
				// Add a small buffer (e.g., 1ms) to prevent potential infinite loops due to floating-point inaccuracies
				if waitTime == 0 && needed > 0 {
					waitTime = time.Millisecond
				}
			}

			// Unlock mutex while waiting
			rl.mu.Unlock()
			// Use a timer that respects context cancellation
			timer := time.NewTimer(waitTime)

			// Wait with context
			select {
			case <-timer.C:
				// Reacquire mutex
				rl.mu.Lock()

				// Update tokens again after waiting
				rl.updateTokens()

				// Re-check if a token is available in the next loop iteration
			case <-ctx.Done():
				// Stop the timer if context is cancelled
				if !timer.Stop() {
					<-timer.C // Drain the channel if Stop() returned false
				}
				// Reacquire mutex before returning
				rl.mu.Lock()
				return ctx.Err()
			}
		}
	}

	// Consume a token
	rl.tokens--

	return nil
}

// TryAcquire attempts to acquire a token without blocking
func (rl *TokenBucketRateLimiter) TryAcquire() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Update tokens
	rl.updateTokens()

	// Check if token is available
	if rl.tokens >= 1 {
		rl.tokens--
		return true
	}

	return false
}

// GetTokens returns the current number of tokens
func (rl *TokenBucketRateLimiter) GetTokens() float64 {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	rl.updateTokens()
	return rl.tokens
}

// GetRate returns the rate at which tokens are replenished
func (rl *TokenBucketRateLimiter) GetRate() float64 {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	return rl.rate
}
