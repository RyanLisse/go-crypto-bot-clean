package http

import (
	"net/http"
	"time"

	httpmiddleware "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/http/middleware"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/factory"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// NewRouter initializes the HTTP router with all middleware and base routes.
func NewRouter(cfg *config.Config, logger *zerolog.Logger, db *gorm.DB) *chi.Mux {
	r := chi.NewRouter()

	// Create consolidated factory
	consolidatedFactory := factory.NewConsolidatedFactory(db, logger, cfg)

	// Global middleware
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	// Use CORS middleware from consolidated factory
	r.Use(httpmiddleware.CORSMiddleware(cfg, logger))

	// Add unified error handling middleware
	errorMiddleware := httpmiddleware.NewUnifiedErrorMiddleware(logger)
	r.Use(errorMiddleware.Middleware())

	// Add security middlewares from consolidated factory
	rateLimiterMiddleware := consolidatedFactory.GetRateLimiterMiddleware()
	r.Use(rateLimiterMiddleware)

	// Add secure headers middleware
	secureHeadersMiddleware := consolidatedFactory.GetSecureHeadersHandler()
	r.Use(secureHeadersMiddleware)

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

// GetAuthMiddleware returns the authentication middleware from the consolidated factory
func GetAuthMiddleware(cfg *config.Config, logger *zerolog.Logger, db *gorm.DB) (httpmiddleware.AuthMiddleware, error) {
	consolidatedFactory := factory.NewConsolidatedFactory(db, logger, cfg)
	authMiddleware, err := consolidatedFactory.GetAuthMiddleware()
	if err != nil {
		return nil, err
	}
	return authMiddleware, nil
}

// GetTestAuthMiddleware returns the test authentication middleware
func GetTestAuthMiddleware(cfg *config.Config, logger *zerolog.Logger, db *gorm.DB) httpmiddleware.AuthMiddleware {
	consolidatedFactory := factory.NewConsolidatedFactory(db, logger, cfg)
	return consolidatedFactory.GetTestAuthMiddleware()
}

// GetDisabledAuthMiddleware returns the disabled authentication middleware
func GetDisabledAuthMiddleware(cfg *config.Config, logger *zerolog.Logger, db *gorm.DB) httpmiddleware.AuthMiddleware {
	consolidatedFactory := factory.NewConsolidatedFactory(db, logger, cfg)
	return consolidatedFactory.GetDisabledAuthMiddleware()
}
