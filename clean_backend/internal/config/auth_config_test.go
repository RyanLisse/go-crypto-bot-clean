package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAuthConfig(t *testing.T) {
	// Save original environment variables
	originalClerkSecretKey := os.Getenv("CLERK_SECRET_KEY")
	originalUseEnhancedAuth := os.Getenv("USE_ENHANCED_AUTH")

	// Restore environment variables after test
	defer func() {
		os.Setenv("CLERK_SECRET_KEY", originalClerkSecretKey)
		os.Setenv("USE_ENHANCED_AUTH", originalUseEnhancedAuth)
	}()

	t.Run("With Required Environment Variables", func(t *testing.T) {
		// Set environment variables
		os.Setenv("CLERK_SECRET_KEY", "test-secret-key")
		os.Setenv("USE_ENHANCED_AUTH", "true")

		// Create config
		config, err := NewAuthConfig()

		// Check result
		assert.NoError(t, err)
		assert.NotNil(t, config)
		assert.Equal(t, "test-secret-key", config.GetClerkSecretKey())
		assert.True(t, config.ShouldUseEnhanced())
	})

	t.Run("Without Clerk Secret Key", func(t *testing.T) {
		// Unset environment variables
		os.Unsetenv("CLERK_SECRET_KEY")
		os.Setenv("USE_ENHANCED_AUTH", "true")

		// Create config
		config, err := NewAuthConfig()

		// Check result
		assert.Error(t, err)
		assert.Nil(t, config)
	})

	t.Run("Without Enhanced Auth Flag", func(t *testing.T) {
		// Set environment variables
		os.Setenv("CLERK_SECRET_KEY", "test-secret-key")
		os.Unsetenv("USE_ENHANCED_AUTH")

		// Create config
		config, err := NewAuthConfig()

		// Check result
		assert.NoError(t, err)
		assert.NotNil(t, config)
		assert.Equal(t, "test-secret-key", config.GetClerkSecretKey())
		assert.False(t, config.ShouldUseEnhanced())
	})
}
