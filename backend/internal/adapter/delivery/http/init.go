package http

import (
	"net/http"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/delivery/http/middleware"
	httpmiddleware "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/http/middleware"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

// NewRouter initializes the HTTP router with all middleware and base routes.
func NewRouter(cfg *config.Config, logger *zerolog.Logger) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(httpmiddleware.CORSMiddleware(cfg, logger))

	// Add standardized error handling middleware
	errorHandler := middleware.NewStandardizedErrorHandler(logger)
	r.Use(errorHandler.RecoverMiddleware())
	r.Use(errorHandler.LoggingMiddleware())
	r.Use(errorHandler.ErrorResponseMiddleware())

	// Add credential rate limiter middleware
	credentialRateLimiter := middleware.NewCredentialRateLimiter(logger)
	r.Use(credentialRateLimiter.Middleware())

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","version":"` + cfg.Version + `","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`))
	})

	// Root level test endpoint
	r.Get("/root-test", func(w http.ResponseWriter, r *http.Request) {
		logger.Info().Msg("Root level test endpoint called")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success":true,"message":"Root level test endpoint works!"}`))
	})

	return r
}
