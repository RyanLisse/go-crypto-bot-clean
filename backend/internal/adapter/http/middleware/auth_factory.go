package middleware

import (
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/service"
	"github.com/rs/zerolog"
)

// AuthFactory creates authentication middleware
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

// CreateAuthMiddleware creates the primary authentication middleware
func (f *AuthFactory) CreateAuthMiddleware() AuthMiddleware {
	// Check if we're in test mode
	if f.cfg.ENV == "test" {
		f.logger.Info().Msg("Using test authentication middleware")
		return NewTestAuthMiddleware(f.logger)
	}

	// Use standard authentication in production
	f.logger.Info().Msg("Using standard authentication middleware")
	return NewAuthMiddleware(f.authService, f.logger)
}
