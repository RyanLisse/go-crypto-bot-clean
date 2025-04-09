package unit_test

import (
	"context"
	"testing"
	"time"

	"go-crypto-bot-clean/backend/pkg/ratelimiter"
	"github.com/stretchr/testify/assert"
)

func TestTokenBucketRateLimiter_Wait(t *testing.T) {
	// Create a rate limiter with 1 token per second and capacity of 2
	rl := ratelimiter.NewTokenBucketRateLimiter(1, 2)

	// First two calls should succeed immediately
	ctx := context.Background()
	err1 := rl.Wait(ctx)
	assert.NoError(t, err1)

	err2 := rl.Wait(ctx)
	assert.NoError(t, err2)

	// Third call should block and take approximately 1 second
	start := time.Now()
	err3 := rl.Wait(ctx)
	duration := time.Since(start)

	assert.NoError(t, err3)
	assert.True(t, duration >= time.Second, "Wait should take at least 1 second")
	assert.True(t, duration < 1500*time.Millisecond, "Wait should not take much longer than 1 second")
}

func TestTokenBucketRateLimiter_TryAcquire(t *testing.T) {
	// Create a rate limiter with 1 token per second and capacity of 2
	rl := ratelimiter.NewTokenBucketRateLimiter(1, 2)

	// First two calls should succeed
	assert.True(t, rl.TryAcquire(), "First token should be acquired")
	assert.True(t, rl.TryAcquire(), "Second token should be acquired")

	// Third call should fail
	assert.False(t, rl.TryAcquire(), "Third token should not be acquired")
}

func TestTokenBucketRateLimiter_GetTokens(t *testing.T) {
	// Create a rate limiter with 1 token per second and capacity of 2
	rl := ratelimiter.NewTokenBucketRateLimiter(1, 2)

	// Initial tokens should be 2
	assert.InDelta(t, float64(2), rl.GetTokens(), 0.001, "Initial tokens should be 2")

	// Acquire two tokens
	rl.TryAcquire()
	rl.TryAcquire()

	// Tokens should now be very close to 0
	assert.InDelta(t, float64(0), rl.GetTokens(), 0.001, "Tokens should be 0 after acquiring 2")

	// Wait for 1 second to allow token replenishment
	time.Sleep(1100 * time.Millisecond)

	// Tokens should now be close to 1
	tokens := rl.GetTokens()
	assert.InDelta(t, 1.0, tokens, 0.2, "Tokens should be close to 1 after 1 second")
}

func TestTokenBucketRateLimiter_ContextCancellation(t *testing.T) {
	// Create a rate limiter with 1 token per second and capacity of 2
	rl := ratelimiter.NewTokenBucketRateLimiter(1, 2)

	// Acquire two tokens first
	rl.TryAcquire()
	rl.TryAcquire()

	// Create a context with a short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Wait should be cancelled
	err := rl.Wait(ctx)
	assert.Error(t, err, "Wait should be cancelled by context")
	assert.Equal(t, context.DeadlineExceeded, err, "Error should be context deadline exceeded")
}
