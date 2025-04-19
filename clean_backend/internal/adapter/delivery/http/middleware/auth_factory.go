package middleware

import (
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// AuthType represents the type of authentication middleware to use
type AuthType string

const (
	// AuthTypeClerk uses Clerk for authentication
	AuthTypeClerk AuthType = "clerk"
	// AuthTypeTest uses test authentication (for testing)
	AuthTypeTest AuthType = "test"
	// AuthTypeDisabled disables authentication
	AuthTypeDisabled AuthType = "disabled"
)

// AuthFactory creates authentication middleware
type AuthFactory struct {
	logger      *zerolog.Logger
	config      *config.Config
	authService port.AuthServiceInterface
}

// NewAuthFactory creates a new AuthFactory
func NewAuthFactory(authService port.AuthServiceInterface, config *config.Config, logger *zerolog.Logger) *AuthFactory {
	return &AuthFactory{
		logger:      logger,
		config:      config,
		authService: authService,
	}
}

// CreateMiddleware creates an authentication middleware based on the configuration
func (f *AuthFactory) CreateMiddleware(authType AuthType) AuthMiddleware {
	switch authType {
	case AuthTypeClerk:
		return NewClerkMiddleware(f.authService, f.config, f.logger)
	case AuthTypeTest:
		return NewTestMiddleware(f.logger)
	case AuthTypeDisabled:
		return NewDisabledMiddleware(f.logger)
	default:
		// Default to Clerk if not specified
		return NewClerkMiddleware(f.authService, f.config, f.logger)
	}
}

// CreateDefaultMiddleware creates the default authentication middleware based on the environment
func (f *AuthFactory) CreateDefaultMiddleware() AuthMiddleware {
	// Use environment to determine which middleware to use
	env := f.config.Server.Host // Just using a placeholder since we don't have APP_ENV

	// Check if auth is disabled
	if f.config.Auth.Disabled {
		return f.CreateMiddleware(AuthTypeDisabled)
	}

	// Check if we're in test mode
	if env == "test" {
		return f.CreateMiddleware(AuthTypeTest)
	}

	// Default to Clerk
	return f.CreateMiddleware(AuthTypeClerk)
}
