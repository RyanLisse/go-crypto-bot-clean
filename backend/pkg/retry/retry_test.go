package retry

// MIGRATION: Tests removed. If needed, write tests for usage of github.com/cenkalti/backoff/v4.


import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRetryableError(t *testing.T) {
	// Test creating and checking retryable errors
	baseErr := errors.New("base error")
	retryableErr := NewRetryableError(baseErr)

	// Test Error method
	assert.Contains(t, retryableErr.Error(), baseErr.Error())

	// Test Unwrap method
	assert.Equal(t, baseErr, retryableErr.Unwrap())

	// Test IsRetryable function
	assert.True(t, IsRetryable(retryableErr))
	assert.False(t, IsRetryable(baseErr))
}

func TestConstantBackoff(t *testing.T) {
	backoff := &ConstantBackoff{
		Interval: 100 * time.Millisecond,
	}

	// Test that backoff returns the same interval regardless of attempt
	assert.Equal(t, 100*time.Millisecond, backoff.NextBackoff(0))
	assert.Equal(t, 100*time.Millisecond, backoff.NextBackoff(1))
	assert.Equal(t, 100*time.Millisecond, backoff.NextBackoff(5))
}

func TestExponentialBackoff(t *testing.T) {
	backoff := &ExponentialBackoff{
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     10 * time.Second,
		Multiplier:      2.0,
		RandomFactor:    0.0, // No randomness for testing
	}

	// Test backoff intervals with multiplier=2
	assert.Equal(t, 100*time.Millisecond, backoff.NextBackoff(0))
	assert.Equal(t, 200*time.Millisecond, backoff.NextBackoff(1))
	assert.Equal(t, 400*time.Millisecond, backoff.NextBackoff(2))
	assert.Equal(t, 800*time.Millisecond, backoff.NextBackoff(3))

	// Test max interval cap
	assert.Equal(t, 10*time.Second, backoff.NextBackoff(10)) // Should exceed max and be capped
}

func TestRetryDoWithSuccess(t *testing.T) {
	// Test retry operation that succeeds
	attempts := 0
	err := Do(context.Background(), func() error {
		attempts++
		return nil // Success
	}, nil) // Use default options

	assert.NoError(t, err)
	assert.Equal(t, 1, attempts) // Should only attempt once
}

func TestRetryDoWithPermanentError(t *testing.T) {
	// Test retry operation with a non-retryable error
	attempts := 0
	permanentErr := errors.New("permanent error")
	err := Do(context.Background(), func() error {
		attempts++
		return permanentErr // Non-retryable
	}, nil) // Use default options

	assert.Error(t, err)
	assert.Equal(t, 1, attempts) // Should only attempt once
}

func TestRetryDoWithRetryableError(t *testing.T) {
	// Test retry operation with a retryable error
	attempts := 0
	retryableErr := NewRetryableError(errors.New("retryable error"))
	err := Do(context.Background(), func() error {
		attempts++
		if attempts < 3 {
			return retryableErr // Retryable error for first 2 attempts
		}
		return nil // Success on 3rd attempt
	}, &RetryOptions{
		MaxAttempts: 3,
		BackoffStrategy: &ConstantBackoff{
			Interval: 1 * time.Millisecond, // Fast for testing
		},
		RetryIf: IsRetryable,
	})

	assert.NoError(t, err)
	assert.Equal(t, 3, attempts) // Should attempt 3 times total
}

func TestRetryDoWithMaxAttemptsExceeded(t *testing.T) {
	// Test retry operation that reaches max attempts without success
	attempts := 0
	retryableErr := NewRetryableError(errors.New("retryable error"))
	err := Do(context.Background(), func() error {
		attempts++
		return retryableErr // Always return retryable error
	}, &RetryOptions{
		MaxAttempts: 3,
		BackoffStrategy: &ConstantBackoff{
			Interval: 1 * time.Millisecond, // Fast for testing
		},
		RetryIf: IsRetryable,
	})

	assert.Error(t, err)
	assert.Equal(t, 3, attempts) // Should attempt 3 times total
	assert.Contains(t, err.Error(), "all 3 attempts failed")
}

func TestRetryDoWithContextCancellation(t *testing.T) {
	// Test retry operation with context cancellation
	attempts := 0
	ctx, cancel := context.WithCancel(context.Background())
	retryableErr := NewRetryableError(errors.New("retryable error"))

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel() // Cancel after a short delay
	}()

	err := Do(ctx, func() error {
		attempts++
		time.Sleep(10 * time.Millisecond) // Ensure at least one attempt
		return retryableErr               // Always return retryable error
	}, &RetryOptions{
		MaxAttempts: 10, // More than we'll actually do due to cancellation
		BackoffStrategy: &ConstantBackoff{
			Interval: 100 * time.Millisecond, // Long enough to ensure context cancellation
		},
		RetryIf: IsRetryable,
	})

	assert.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
	assert.True(t, attempts > 0, "Should have attempted at least once")
}

func TestCustomRetryIfFunction(t *testing.T) {
	// Test custom retryIf function
	customErr := errors.New("custom error")
	customRetryIf := func(err error) bool {
		return errors.Is(err, customErr)
	}

	attempts := 0
	err := Do(context.Background(), func() error {
		attempts++
		if attempts < 3 {
			return customErr // Custom retryable error
		}
		return nil // Success
	}, &RetryOptions{
		MaxAttempts: 3,
		BackoffStrategy: &ConstantBackoff{
			Interval: 1 * time.Millisecond,
		},
		RetryIf: customRetryIf,
	})

	assert.NoError(t, err)
	assert.Equal(t, 3, attempts)
}
