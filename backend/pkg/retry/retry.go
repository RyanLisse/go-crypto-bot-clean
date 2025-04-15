package retry

// MIGRATION: This package has been replaced by github.com/cenkalti/backoff/v4.
// Use backoff.Retry and related types for all retry logic.


import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"
)

// RetryableError is an error that indicates the operation should be retried
type RetryableError struct {
	Err error
}

func (e *RetryableError) Error() string {
	return fmt.Sprintf("retryable error: %v", e.Err)
}

func (e *RetryableError) Unwrap() error {
	return e.Err
}

// IsRetryable checks if an error is retryable
func IsRetryable(err error) bool {
	var retryableErr *RetryableError
	return errors.As(err, &retryableErr)
}

// NewRetryableError creates a new retryable error
func NewRetryableError(err error) *RetryableError {
	return &RetryableError{Err: err}
}

// BackoffStrategy defines a strategy for calculating backoff time between retries
type BackoffStrategy interface {
	NextBackoff(attempt int) time.Duration
}

// ConstantBackoff implements a constant backoff strategy
type ConstantBackoff struct {
	Interval time.Duration
}

func (b *ConstantBackoff) NextBackoff(_ int) time.Duration {
	return b.Interval
}

// ExponentialBackoff implements an exponential backoff strategy
type ExponentialBackoff struct {
	InitialInterval time.Duration
	MaxInterval     time.Duration
	Multiplier      float64
	RandomFactor    float64
}

func (b *ExponentialBackoff) NextBackoff(attempt int) time.Duration {
	// Calculate backoff time
	backoff := float64(b.InitialInterval) * math.Pow(b.Multiplier, float64(attempt))

	// Add jitter
	if b.RandomFactor > 0 {
		backoff = backoff * (1 + b.RandomFactor*(rand.Float64()*2-1))
	}

	// Ensure backoff doesn't exceed max interval
	if backoff > float64(b.MaxInterval) {
		backoff = float64(b.MaxInterval)
	}

	return time.Duration(backoff)
}

// RetryOptions defines options for retry operations
type RetryOptions struct {
	MaxAttempts     int
	BackoffStrategy BackoffStrategy
	RetryIf         func(error) bool
}

// DefaultRetryOptions returns default retry options
func DefaultRetryOptions() *RetryOptions {
	return &RetryOptions{
		MaxAttempts: 3,
		BackoffStrategy: &ExponentialBackoff{
			InitialInterval: 100 * time.Millisecond,
			MaxInterval:     10 * time.Second,
			Multiplier:      2.0,
			RandomFactor:    0.2,
		},
		RetryIf: IsRetryable,
	}
}

// Do executes the provided function with retries based on the provided options
func Do(ctx context.Context, fn func() error, opts *RetryOptions) error {
	if opts == nil {
		opts = DefaultRetryOptions()
	}

	var lastErr error

	for attempt := 0; attempt < opts.MaxAttempts; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}

		lastErr = err

		// Check if error is retryable
		if opts.RetryIf != nil && !opts.RetryIf(err) {
			return err
		}

		// Check if we've reached max attempts
		if attempt >= opts.MaxAttempts-1 {
			break
		}

		// Calculate backoff time
		backoff := opts.BackoffStrategy.NextBackoff(attempt)

		// Create a timer for backoff
		timer := time.NewTimer(backoff)
		defer timer.Stop()

		// Wait for backoff timer or context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			// Continue with the next attempt
		}
	}

	return fmt.Errorf("all %d attempts failed, last error: %w", opts.MaxAttempts, lastErr)
}
