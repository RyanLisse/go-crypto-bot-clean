// Package huma provides OpenAPI documentation for the API.
package huma

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"

	"go-crypto-bot-clean/api/huma/auth"
	"go-crypto-bot-clean/api/huma/strategy"
	"go-crypto-bot-clean/api/huma/user"
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
func SetupHuma(router chi.Router, config Config) huma.API {
	// Create a new Huma API
	api := humachi.New(router, huma.DefaultConfig(config.Title, config.Version))

	// Register endpoints
	registerBacktestEndpoints(api, config.BasePath)
	strategy.RegisterEndpoints(api, config.BasePath)
	auth.RegisterEndpoints(api, config.BasePath)
	user.RegisterEndpoints(api, config.BasePath)

	return api
}
