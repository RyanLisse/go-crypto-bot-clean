package middleware

import (
	"os"
	"testing"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestNewAuthFactory(t *testing.T) {
	// Create a logger
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create a mock auth service
	mockAuthService := new(MockAuthService)

	// Create a mock config
	cfg := &config.Config{
		ENV: "test",
	}

	// Create auth factory
	factory := NewAuthFactory(cfg, &logger, mockAuthService)

	// Verify that the factory is not nil
	assert.NotNil(t, factory)
}

func TestAuthFactory_CreateAuthMiddleware(t *testing.T) {
	// Create a logger
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create a mock auth service
	mockAuthService := new(MockAuthService)

	// Create a mock config
	cfg := &config.Config{
		ENV: "test",
	}

	// Create auth factory
	factory := NewAuthFactory(cfg, &logger, mockAuthService)

	// Create auth middleware
	middleware := factory.CreateAuthMiddleware()

	// Verify that the middleware is not nil
	assert.NotNil(t, middleware)
}
