package trade

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/cenkalti/backoff/v4"
	"github.com/rs/zerolog"
	"golang.org/x/time/rate"
)

// RateLimitedExecutor implements the TradeExecutor interface with rate limiting and error handling
type RateLimitedExecutor struct {
	tradeService    port.TradeService
	logger          *zerolog.Logger
	rateLimiter     *rate.Limiter
	retryBackoff    backoff.BackOff
	maxRetries      int
	retryableErrors map[int]bool
	mutex           sync.RWMutex
	retryDelay      time.Duration
	// Track rate limit windows
	rateLimitResetTime time.Time
	rateLimitRemaining int
}

// RateLimitedExecutorConfig contains configuration for the rate limited executor
type RateLimitedExecutorConfig struct {
	RequestsPerSecond float64
	BurstSize         int
	MaxRetries        int
	InitialRetryDelay time.Duration
	MaxRetryDelay     time.Duration
}

// DefaultExecutorConfig returns a default configuration for the rate limited executor
func DefaultExecutorConfig() RateLimitedExecutorConfig {
	return RateLimitedExecutorConfig{
		RequestsPerSecond: 1.0,  // 1 request per second by default
		BurstSize:         3,    // Allow bursts of 3 requests
		MaxRetries:        5,    // Maximum number of retries
		InitialRetryDelay: 500 * time.Millisecond,
		MaxRetryDelay:     30 * time.Second,
	}
}

// NewRateLimitedExecutor creates a new rate limited executor
func NewRateLimitedExecutor(
	tradeService port.TradeService,
	logger *zerolog.Logger,
	config RateLimitedExecutorConfig,
) *RateLimitedExecutor {
	// Create exponential backoff strategy
	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.InitialInterval = config.InitialRetryDelay
	expBackoff.MaxInterval = config.MaxRetryDelay
	expBackoff.MaxElapsedTime = 0 // No maximum elapsed time, rely on context for timeouts

	// Initialize retryable error codes
	retryableErrors := map[int]bool{
		http.StatusTooManyRequests:     true, // 429 - Rate limit exceeded
		http.StatusInternalServerError: true, // 500 - Server error
		http.StatusBadGateway:          true, // 502 - Bad gateway
		http.StatusServiceUnavailable:  true, // 503 - Service unavailable
		http.StatusGatewayTimeout:      true, // 504 - Gateway timeout
	}

	return &RateLimitedExecutor{
		tradeService:    tradeService,
		logger:          logger,
		rateLimiter:     rate.NewLimiter(rate.Limit(config.RequestsPerSecond), config.BurstSize),
		retryBackoff:    expBackoff,
		maxRetries:      config.MaxRetries,
		retryableErrors: retryableErrors,
		retryDelay:      config.InitialRetryDelay,
	}
}

// ExecuteOrder places an order with error handling and rate limiting
func (e *RateLimitedExecutor) ExecuteOrder(ctx context.Context, request *model.OrderRequest) (*model.OrderResponse, error) {
	var response *model.OrderResponse
	var err error

	operation := func() error {
		// Check if we need to wait for rate limit reset
		if e.shouldWaitForRateLimit() {
			waitTime := time.Until(e.rateLimitResetTime)
			e.logger.Info().
				Dur("waitTime", waitTime).
				Time("resetTime", e.rateLimitResetTime).
				Msg("Waiting for rate limit reset")
			
			select {
			case <-time.After(waitTime):
				// Continue after waiting
			case <-ctx.Done():
				return backoff.Permanent(ctx.Err())
			}
		}

		// Apply rate limiting
		if err := e.rateLimiter.Wait(ctx); err != nil {
			return backoff.Permanent(fmt.Errorf("rate limiter wait failed: %w", err))
		}

		// Execute the order
		resp, orderErr := e.tradeService.PlaceOrder(ctx, request)
		if orderErr != nil {
			// Check if this is a rate limit error
			if e.isRateLimitError(orderErr) {
				e.updateRateLimitInfo(orderErr)
				return orderErr // Will be retried
			}

			// Check if this is another retryable error
			if e.isRetryableError(orderErr) {
				e.logger.Warn().
					Err(orderErr).
					Str("symbol", request.Symbol).
					Str("side", string(request.Side)).
					Msg("Retryable error placing order")
				return orderErr // Will be retried
			}

			// Non-retryable error
			e.logger.Error().
				Err(orderErr).
				Str("symbol", request.Symbol).
				Str("side", string(request.Side)).
				Msg("Non-retryable error placing order")
			return backoff.Permanent(orderErr)
		}

		// Success
		response = resp
		return nil
	}

	// Execute with backoff
	err = backoff.Retry(operation, backoff.WithContext(e.retryBackoff, ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to execute order after retries: %w", err)
	}

	return response, nil
}

// CancelOrderWithRetry attempts to cancel an order with retries
func (e *RateLimitedExecutor) CancelOrderWithRetry(ctx context.Context, symbol, orderID string) error {
	operation := func() error {
		// Check if we need to wait for rate limit reset
		if e.shouldWaitForRateLimit() {
			waitTime := time.Until(e.rateLimitResetTime)
			e.logger.Info().
				Dur("waitTime", waitTime).
				Time("resetTime", e.rateLimitResetTime).
				Msg("Waiting for rate limit reset")
			
			select {
			case <-time.After(waitTime):
				// Continue after waiting
			case <-ctx.Done():
				return backoff.Permanent(ctx.Err())
			}
		}

		// Apply rate limiting
		if err := e.rateLimiter.Wait(ctx); err != nil {
			return backoff.Permanent(fmt.Errorf("rate limiter wait failed: %w", err))
		}

		// Cancel the order
		err := e.tradeService.CancelOrder(ctx, symbol, orderID)
		if err != nil {
			// Check if this is a rate limit error
			if e.isRateLimitError(err) {
				e.updateRateLimitInfo(err)
				return err // Will be retried
			}

			// Check if this is another retryable error
			if e.isRetryableError(err) {
				e.logger.Warn().
					Err(err).
					Str("symbol", symbol).
					Str("orderID", orderID).
					Msg("Retryable error canceling order")
				return err // Will be retried
			}

			// Non-retryable error
			e.logger.Error().
				Err(err).
				Str("symbol", symbol).
				Str("orderID", orderID).
				Msg("Non-retryable error canceling order")
			return backoff.Permanent(err)
		}

		return nil
	}

	// Execute with backoff
	return backoff.Retry(operation, backoff.WithContext(e.retryBackoff, ctx))
}

// GetOrderStatusWithRetry attempts to get order status with retries
func (e *RateLimitedExecutor) GetOrderStatusWithRetry(ctx context.Context, symbol, orderID string) (*model.Order, error) {
	var order *model.Order

	operation := func() error {
		// Check if we need to wait for rate limit reset
		if e.shouldWaitForRateLimit() {
			waitTime := time.Until(e.rateLimitResetTime)
			e.logger.Info().
				Dur("waitTime", waitTime).
				Time("resetTime", e.rateLimitResetTime).
				Msg("Waiting for rate limit reset")
			
			select {
			case <-time.After(waitTime):
				// Continue after waiting
			case <-ctx.Done():
				return backoff.Permanent(ctx.Err())
			}
		}

		// Apply rate limiting
		if err := e.rateLimiter.Wait(ctx); err != nil {
			return backoff.Permanent(fmt.Errorf("rate limiter wait failed: %w", err))
		}

		// Get order status
		o, err := e.tradeService.GetOrderStatus(ctx, symbol, orderID)
		if err != nil {
			// Check if this is a rate limit error
			if e.isRateLimitError(err) {
				e.updateRateLimitInfo(err)
				return err // Will be retried
			}

			// Check if this is another retryable error
			if e.isRetryableError(err) {
				e.logger.Warn().
					Err(err).
					Str("symbol", symbol).
					Str("orderID", orderID).
					Msg("Retryable error getting order status")
				return err // Will be retried
			}

			// Non-retryable error
			e.logger.Error().
				Err(err).
				Str("symbol", symbol).
				Str("orderID", orderID).
				Msg("Non-retryable error getting order status")
			return backoff.Permanent(err)
		}

		order = o
		return nil
	}

	// Execute with backoff
	err := backoff.Retry(operation, backoff.WithContext(e.retryBackoff, ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to get order status after retries: %w", err)
	}

	return order, nil
}

// isRateLimitError checks if an error is a rate limit error
func (e *RateLimitedExecutor) isRateLimitError(err error) bool {
	// Check for HTTP 429 status code
	var httpErr *model.HTTPError
	if errors.As(err, &httpErr) {
		return httpErr.StatusCode == http.StatusTooManyRequests
	}
	
	// Check for rate limit error message
	return errors.Is(err, model.ErrRateLimitExceeded)
}

// isRetryableError checks if an error is retryable
func (e *RateLimitedExecutor) isRetryableError(err error) bool {
	// Check for HTTP errors with retryable status codes
	var httpErr *model.HTTPError
	if errors.As(err, &httpErr) {
		return e.retryableErrors[httpErr.StatusCode]
	}
	
	// Check for network errors
	if errors.Is(err, model.ErrNetworkFailure) || 
	   errors.Is(err, model.ErrTimeout) ||
	   errors.Is(err, model.ErrConnectionReset) {
		return true
	}
	
	return false
}

// updateRateLimitInfo updates rate limit information from response headers
func (e *RateLimitedExecutor) updateRateLimitInfo(err error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	
	var httpErr *model.HTTPError
	if !errors.As(err, &httpErr) || httpErr.Headers == nil {
		// Default backoff if we can't extract headers
		e.rateLimitResetTime = time.Now().Add(e.retryDelay)
		e.rateLimitRemaining = 0
		return
	}
	
	// Extract rate limit headers
	if resetStr, ok := httpErr.Headers["X-RateLimit-Reset"]; ok {
		if resetTime, err := strconv.ParseInt(resetStr, 10, 64); err == nil {
			e.rateLimitResetTime = time.Unix(resetTime, 0)
		}
	} else {
		// Default to current time + retry delay if header not found
		e.rateLimitResetTime = time.Now().Add(e.retryDelay)
	}
	
	if remainingStr, ok := httpErr.Headers["X-RateLimit-Remaining"]; ok {
		if remaining, err := strconv.Atoi(remainingStr); err == nil {
			e.rateLimitRemaining = remaining
		}
	} else {
		e.rateLimitRemaining = 0
	}
	
	e.logger.Info().
		Time("resetTime", e.rateLimitResetTime).
		Int("remaining", e.rateLimitRemaining).
		Msg("Updated rate limit information")
}

// shouldWaitForRateLimit checks if we should wait for rate limit reset
func (e *RateLimitedExecutor) shouldWaitForRateLimit() bool {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	
	return e.rateLimitRemaining <= 0 && time.Now().Before(e.rateLimitResetTime)
}

// Ensure RateLimitedExecutor implements port.TradeExecutor
var _ port.TradeExecutor = (*RateLimitedExecutor)(nil)
