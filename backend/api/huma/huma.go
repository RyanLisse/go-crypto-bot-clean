// Package huma provides OpenAPI documentation for the API.
package huma

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"

	"go-crypto-bot-clean/backend/api/huma/auth"
	"go-crypto-bot-clean/backend/api/huma/strategy"
	"go-crypto-bot-clean/backend/api/huma/user"
	"go-crypto-bot-clean/backend/api/service"
)

// Config represents the configuration for the Huma API documentation.
type Config struct {
	Title       string
	Description string
	Version     string
	BasePath    string
}

// DefaultConfig returns a default configuration for the Huma API documentation.
func DefaultConfig() Config {
	return Config{
		Title:       "Crypto Trading Bot API",
		Description: "API for the cryptocurrency trading bot",
		Version:     "1.0.0",
		BasePath:    "/api/v1",
	}
}

// SetupHuma sets up the Huma API documentation.
func SetupHuma(router chi.Router, config Config, services *service.Provider) huma.API {
	// Create a new Huma API
	api := humachi.New(router, huma.DefaultConfig(config.Title, config.Version))

	// Register endpoints
	registerBacktestEndpointsWithService(api, config.BasePath, services)
	strategy.RegisterStrategyEndpoints(api, config.BasePath, services.StrategyService)
	auth.RegisterAuthEndpoints(api, config.BasePath, services.AuthService)
	user.RegisterUserEndpoints(api, config.BasePath, services.UserService)

	return api
}
