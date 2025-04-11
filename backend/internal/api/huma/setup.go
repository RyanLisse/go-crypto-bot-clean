package huma

import (
	"go-crypto-bot-clean/backend/internal/api/huma/auth"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// SetupAPI initializes and configures the Huma API
func SetupAPI(router *chi.Mux) huma.API {
	// Create a new Huma API instance using the Chi router
	api := humachi.New(router, huma.DefaultConfig("Crypto Bot API", "1.0.0"))

	// Add common middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Base path for API endpoints
	const basePath = "/api/v1"

	// Register all endpoints
	auth.RegisterEndpoints(api, basePath)
	RegisterAccountEndpoints(api, basePath)

	return api
}
