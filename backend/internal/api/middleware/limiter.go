// Package middleware contains API middleware components.
package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/ryanlisse/go-crypto-bot/internal/api/dto/response"
	"github.com/ryanlisse/go-crypto-bot/pkg/ratelimiter"
)

// RateLimiterMiddleware limits requests per identifier (IP or API key).
//
//	@summary	Rate limiting middleware
//	@description	Limits requests per IP or API key using token bucket algorithm.
func RateLimiterMiddleware(rate, capacity float64, extractor func(*gin.Context) string, logger Logger) gin.HandlerFunc {
	var limiters sync.Map // map[string]*ratelimiter.TokenBucketRateLimiter

	return func(c *gin.Context) {
		id := extractor(c)
		if id == "" {
			id = "unknown"
		}

		limiterIface, _ := limiters.LoadOrStore(id, ratelimiter.NewTokenBucketRateLimiter(rate, capacity))
		limiter := limiterIface.(*ratelimiter.TokenBucketRateLimiter)

		if !limiter.TryAcquire() {
			logger.Error("rate limit exceeded for", id)
			c.AbortWithStatusJSON(http.StatusTooManyRequests, response.ErrorResponse{
				Code:    "rate_limited",
				Message: "Too many requests",
			})
			return
		}

		c.Next()
	}
}
