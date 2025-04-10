package ratelimit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestFixedWindowRateLimiter(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	limiter := NewFixedWindowRateLimiter(3, time.Second, logger)

	// First 3 requests should be allowed
	assert.True(t, limiter.Allow("test"))
	assert.True(t, limiter.Allow("test"))
	assert.True(t, limiter.Allow("test"))

	// 4th request should be denied
	assert.False(t, limiter.Allow("test"))

	// Different key should be allowed
	assert.True(t, limiter.Allow("test2"))

	// Reset should allow more requests
	limiter.Reset("test")
	assert.True(t, limiter.Allow("test"))
}

func TestSlidingWindowRateLimiter(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	limiter := NewSlidingWindowRateLimiter(3, time.Second, logger)

	// First 3 requests should be allowed
	assert.True(t, limiter.Allow("test"))
	assert.True(t, limiter.Allow("test"))
	assert.True(t, limiter.Allow("test"))

	// 4th request should be denied
	assert.False(t, limiter.Allow("test"))

	// Different key should be allowed
	assert.True(t, limiter.Allow("test2"))

	// Reset should allow more requests
	limiter.Reset("test")
	assert.True(t, limiter.Allow("test"))
}

func TestTokenBucketRateLimiter(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	limiter := NewTokenBucketRateLimiter(3, 3, time.Second, logger)

	// First 3 requests should be allowed (burst)
	assert.True(t, limiter.Allow("test"))
	assert.True(t, limiter.Allow("test"))
	assert.True(t, limiter.Allow("test"))

	// Add small sleep to ensure time progresses slightly, although logic should handle this.
	time.Sleep(5 * time.Millisecond)
	// 4th request should be denied
	assert.False(t, limiter.Allow("test"), "4th request immediately after burst should be denied")

	// Different key should be allowed
	assert.True(t, limiter.Allow("test2"))

	// Reset should allow more requests
	limiter.Reset("test")
	assert.True(t, limiter.Allow("test"))

	// Test refill
	limiter = NewTokenBucketRateLimiter(60, 1, time.Minute, logger)
	assert.True(t, limiter.Allow("test3"))
	assert.False(t, limiter.Allow("test3"))

	// Wait for token to refill (should take 1 second)
	time.Sleep(1100 * time.Millisecond)
	// With rate 1/min, 1.1s is not enough to refill. Assertion should be False.
	assert.False(t, limiter.Allow("test3"), "Should not have refilled 1 token in ~1 second with 1/min rate")
}
