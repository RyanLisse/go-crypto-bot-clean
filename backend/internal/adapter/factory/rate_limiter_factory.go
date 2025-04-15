package factory

import (
	"net/http"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/http/middleware"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/rs/zerolog"
)

// RateLimiterFactory creates rate limiters
type RateLimiterFactory struct {
	logger *zerolog.Logger
}

// NewRateLimiterFactory creates a new RateLimiterFactory
func NewRateLimiterFactory(logger *zerolog.Logger) *RateLimiterFactory {
	return &RateLimiterFactory{
		logger: logger,
	}
}

// CreateAdvancedRateLimiter creates an advanced rate limiter
func (f *RateLimiterFactory) CreateAdvancedRateLimiter(cfg *config.RateLimitConfig) *middleware.AdvancedRateLimiter {
	return middleware.NewAdvancedRateLimiter(cfg, f.logger)
}

// CreateAdvancedRateLimiterMiddleware creates an advanced rate limiter middleware
func (f *RateLimiterFactory) CreateAdvancedRateLimiterMiddleware(cfg *config.RateLimitConfig) func(next http.Handler) http.Handler {
	limiter := f.CreateAdvancedRateLimiter(cfg)
	return middleware.AdvancedRateLimiterMiddleware(limiter)
}
