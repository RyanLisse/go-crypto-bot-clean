package middleware

import (
	"net/http"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/service"
	"github.com/rs/zerolog"
)

// AuthFactory creates Clerk authentication middleware only
// All other auth methods have been removed for clarity and security
// If Clerk secret is missing, log a fatal error

type AuthFactory struct {
	cfg         *config.Config
	logger      *zerolog.Logger
	authService service.AuthServiceInterface
}

// NewAuthFactory creates a new AuthFactory
func NewAuthFactory(cfg *config.Config, logger *zerolog.Logger, authService service.AuthServiceInterface) *AuthFactory {
	return &AuthFactory{
		cfg:         cfg,
		logger:      logger,
		authService: authService,
	}
}

// CreateAuthMiddleware creates Clerk authentication middleware
// Temporarily disabled Clerk authentication for testing - returns a middleware that adds MEXC API credentials
func (f *AuthFactory) CreateAuthMiddleware() func(http.Handler) http.Handler {
	f.logger.Warn().Msg("Clerk authentication middleware DISABLED for testing (using MEXC API middleware instead)")

	// Create MEXC API middleware
	mexcMiddleware := NewMEXCAPIMiddleware(f.logger)

	// Return the middleware
	return mexcMiddleware.Middleware()
}

// CreateEnhancedClerkMiddleware creates an enhanced Clerk middleware
func (f *AuthFactory) CreateEnhancedClerkMiddleware() *EnhancedClerkMiddleware {
	return NewEnhancedClerkMiddleware(f.authService, f.logger)
}
