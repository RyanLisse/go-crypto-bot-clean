package config

import (
	"errors"
	"os"
)

// AuthConfig contains authentication configuration
type AuthConfig struct {
	ClerkSecretKey string
	UseEnhanced    bool
}

// NewAuthConfig creates a new AuthConfig
func NewAuthConfig() (*AuthConfig, error) {
	// Get Clerk secret key from environment
	clerkSecretKey := os.Getenv("CLERK_SECRET_KEY")
	if clerkSecretKey == "" {
		return nil, errors.New("CLERK_SECRET_KEY environment variable is required")
	}

	// Check if enhanced authentication should be used
	useEnhanced := os.Getenv("USE_ENHANCED_AUTH") == "true"

	return &AuthConfig{
		ClerkSecretKey: clerkSecretKey,
		UseEnhanced:    useEnhanced,
	}, nil
}

// GetClerkSecretKey returns the Clerk secret key
func (c *AuthConfig) GetClerkSecretKey() string {
	return c.ClerkSecretKey
}

// ShouldUseEnhanced returns whether enhanced authentication should be used
func (c *AuthConfig) ShouldUseEnhanced() bool {
	return c.UseEnhanced
}
