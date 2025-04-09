// Package auth provides public interfaces for authentication services
package auth

import (
	"context"

	"go-crypto-bot-clean/backend/internal/auth"
)

// Service is a public interface for the authentication service
type Service interface {
	// Authenticate authenticates a user with the given credentials
	Authenticate(ctx context.Context, username, password string) (string, error)

	// ValidateToken validates a JWT token
	ValidateToken(token string) (string, error)

	// RefreshToken refreshes a JWT token
	RefreshToken(token string) (string, error)
}

// serviceAdapter adapts the internal auth service to the public interface
type serviceAdapter struct {
	internalService *auth.Service
}

// Authenticate adapts the internal Authenticate method to the public interface
func (a *serviceAdapter) Authenticate(ctx context.Context, username, password string) (string, error) {
	// This is a simplified implementation - in a real app, you'd create an HTTP request with credentials
	// and use the internal service to authenticate it
	return "token", nil // Placeholder implementation
}

// ValidateToken adapts the internal ValidateToken method to the public interface
func (a *serviceAdapter) ValidateToken(token string) (string, error) {
	// Placeholder implementation
	return "user_id", nil
}

// RefreshToken adapts the internal RefreshToken method to the public interface
func (a *serviceAdapter) RefreshToken(token string) (string, error) {
	// Placeholder implementation
	return "new_token", nil
}

// NewService creates a new authentication service
func NewService(secretKey string) Service {
	internalService := auth.NewService(secretKey)
	return &serviceAdapter{
		internalService: internalService,
	}
}
