package factory

import (
	"net/http"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/http/middleware"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/rs/zerolog"
)

// SecurityFactory creates security-related components
type SecurityFactory struct {
	logger *zerolog.Logger
}

// NewSecurityFactory creates a new SecurityFactory
func NewSecurityFactory(logger *zerolog.Logger) *SecurityFactory {
	return &SecurityFactory{
		logger: logger,
	}
}

// CreateRateLimiter creates an advanced rate limiter
func (f *SecurityFactory) CreateRateLimiter(cfg *config.RateLimitConfig) *middleware.AdvancedRateLimiter {
	return middleware.NewAdvancedRateLimiter(cfg, f.logger)
}

// CreateRateLimiterMiddleware creates an advanced rate limiter middleware
func (f *SecurityFactory) CreateRateLimiterMiddleware(cfg *config.RateLimitConfig) func(next http.Handler) http.Handler {
	limiter := f.CreateRateLimiter(cfg)
	return middleware.AdvancedRateLimiterMiddleware(limiter)
}

// CreateCSRFMiddleware creates a CSRF middleware
func (f *SecurityFactory) CreateCSRFMiddleware(cfg *config.CSRFConfig) *middleware.CSRFMiddleware {
	return middleware.NewCSRFMiddleware(cfg, f.logger)
}

// CreateCSRFProtectionMiddleware creates a CSRF protection middleware
func (f *SecurityFactory) CreateCSRFProtectionMiddleware(cfg *config.CSRFConfig) func(next http.Handler) http.Handler {
	csrfMiddleware := f.CreateCSRFMiddleware(cfg)
	return csrfMiddleware.Middleware()
}

// CreateSecureHeadersMiddleware creates a secure headers middleware
func (f *SecurityFactory) CreateSecureHeadersMiddleware(cfg *config.SecureHeadersConfig) *middleware.SecureHeadersMiddleware {
	return middleware.NewSecureHeadersMiddleware(cfg, f.logger)
}

// CreateSecureHeadersHandler creates a secure headers handler
func (f *SecurityFactory) CreateSecureHeadersHandler(cfg *config.SecureHeadersConfig) func(next http.Handler) http.Handler {
	secureHeadersMiddleware := f.CreateSecureHeadersMiddleware(cfg)
	return secureHeadersMiddleware.Middleware()
}
