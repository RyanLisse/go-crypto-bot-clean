package server

import (
	"net/http"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

// HumaRouter wraps a Chi router with Huma API integration
type HumaRouter struct {
	chi    *chi.Mux
	logger *zerolog.Logger
	config *config.Config
}

// NewHumaRouter creates a new HumaRouter
func NewHumaRouter(cfg *config.Config, authConfig *config.AuthConfig, logger *zerolog.Logger) *HumaRouter {
	// Create Chi router
	r := chi.NewRouter()

	// Use Chi middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// For now, we'll just use Chi without Huma
	return &HumaRouter{
		chi:    r,
		logger: logger,
		config: cfg,
	}
}

// ServeHTTP implements the http.Handler interface
func (r *HumaRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.chi.ServeHTTP(w, req)
}

// API returns the Chi router for now
func (r *HumaRouter) API() interface{} {
	return r.chi
}

// RegisterHealthCheck registers a health check endpoint
func (r *HumaRouter) RegisterHealthCheck() {
	r.chi.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`))
	})
}

// RegisterOpenAPI registers OpenAPI documentation
func (r *HumaRouter) RegisterOpenAPI() {
	// For now, we'll just log that we're registering OpenAPI documentation
	r.logger.Info().Msg("Registering OpenAPI documentation")

	// We'll implement the actual OpenAPI documentation later
}
