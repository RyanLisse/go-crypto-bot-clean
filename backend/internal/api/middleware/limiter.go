// Package middleware contains API middleware components.
package middleware

import (
	"encoding/json"
	"net/http"
	"sync"

	"go-crypto-bot-clean/backend/internal/api/dto/response"
	"go-crypto-bot-clean/backend/pkg/ratelimiter"
)

// RateLimiterMiddleware limits requests per identifier (IP or API key).
//
//	@summary	Rate limiting middleware
//	@description	Limits requests per IP or API key using token bucket algorithm.
func RateLimiterMiddleware(rate, capacity float64, extractor func(*http.Request) string, logger Logger) func(http.Handler) http.Handler {
	var limiters sync.Map // map[string]*ratelimiter.TokenBucketRateLimiter

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := extractor(r)
			if id == "" {
				id = "unknown"
			}

			limiterIface, _ := limiters.LoadOrStore(id, ratelimiter.NewTokenBucketRateLimiter(rate, capacity))
			limiter := limiterIface.(*ratelimiter.TokenBucketRateLimiter)

			if !limiter.TryAcquire() {
				logger.Error("rate limit exceeded for", id)
				w.WriteHeader(http.StatusTooManyRequests)
				json.NewEncoder(w).Encode(response.ErrorResponse{
					Code:    "rate_limited",
					Message: "Too many requests",
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
