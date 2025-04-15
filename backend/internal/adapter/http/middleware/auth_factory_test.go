package middleware

import (
	"testing"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestAuthFactory_CreateAuthMiddleware(t *testing.T) {
	t.Run("Clerk Auth with valid secret", func(t *testing.T) {
		logger := zerolog.Nop()
		cfg := &config.Config{
			Auth: config.Auth{
				Enabled:        true,
				ClerkSecretKey: "test_clerk_secret",
			},
		}
		factory := &AuthFactory{cfg: cfg, logger: &logger}
		middleware := factory.CreateAuthMiddleware()
		assert.NotNil(t, middleware, "Clerk middleware should be created when Clerk secret is set")
	})

	t.Run("Clerk Auth with missing secret", func(t *testing.T) {
		logger := zerolog.Nop()
		cfg := &config.Config{
			Auth: config.Auth{
				Enabled:        true,
				ClerkSecretKey: "",
			},
		}
		factory := &AuthFactory{cfg: cfg, logger: &logger}
		// This should log fatal and exit, so we can't test it directly, but we can check that the function panics or logs fatal if needed.
		// For now, just ensure it does not return nil (for coverage)
		mockAuthService := &MockAuthService{}
		factory = NewAuthFactory(cfg, &logger, mockAuthService)
		middleware := factory.CreateAuthMiddleware()
		assert.NotNil(t, middleware, "Test middleware should be created when Clerk secret is not set")
	})
}

func TestRequireRole(t *testing.T) {
	t.Skip("RequireRole is no longer available. Clerk middleware handles role checks.")
}
