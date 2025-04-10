package api

import (
	"go-crypto-bot-clean/backend/internal/api/handlers"
	"go-crypto-bot-clean/backend/internal/api/middleware" // Corrected path
	"go-crypto-bot-clean/backend/internal/config"

	"github.com/go-chi/chi/v5"
)

// AuthDependencies contains the dependencies for authentication.
type AuthDependencies struct {
	// Handlers
	AuthHandler *handlers.AuthHandler

	// Configuration
	Config *config.Config
}

// NewAuthDependencies creates a new AuthDependencies instance.
func NewAuthDependencies(cfg *config.Config) *AuthDependencies {
	deps := &AuthDependencies{
		Config: cfg,
	}

	// Create auth service
	// authService is no longer needed here as AuthHandler doesn't require it

	// Initialize auth handler
	// Auth handler no longer needs service injection
	deps.AuthHandler = handlers.NewAuthHandler()

	return deps
}

// SetupAuthRoutes adds authentication routes to the router.
func SetupAuthRoutes(r chi.Router, deps *AuthDependencies) {
	r.Route("/auth", func(r chi.Router) {
		// Login/Logout routes removed as they are likely handled by Clerk frontend/middleware

		r.Route("/", func(r chi.Router) {
			if deps.Config.Auth.Enabled {
				r.Use(middleware.JWTAuthMiddleware(deps.Config.Auth.JWTSecret))
			}
			r.Get("/me", deps.AuthHandler.GetCurrentUser)
		})
	})
}
