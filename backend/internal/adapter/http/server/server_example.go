package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/http/middleware"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"
)

// ExampleServer demonstrates how to set up a server with the unified error middleware and the example error controller.
type ExampleServer struct {
	router     chi.Router
	httpServer *http.Server
	logger     *zerolog.Logger
	config     *config.Config
}

// NewExampleServer creates a new example server
func NewExampleServer(cfg *config.Config, logger *zerolog.Logger) *ExampleServer {
	router := chi.NewRouter()

	return &ExampleServer{
		router:     router,
		httpServer: &http.Server{},
		logger:     logger,
		config:     cfg,
	}
}

// SetupRoutes sets up the routes for the example server
func (s *ExampleServer) SetupRoutes() error {
	// Create handlers
	// Commented out for now as we don't have the example controller
	// errorExampleHandler := example.NewErrorExampleController(s.logger)

	// Set up the unified error middleware
	errorMiddleware := middleware.NewUnifiedErrorMiddleware(s.logger)

	// Set up standard middleware
	s.router.Use(chimiddleware.RequestID)
	s.router.Use(chimiddleware.RealIP)
	s.router.Use(chimiddleware.Logger)
	s.router.Use(chimiddleware.StripSlashes)

	// Use our custom error middleware (should be early in the chain)
	s.router.Use(errorMiddleware.Middleware())

	// Other middleware
	s.router.Use(chimiddleware.Timeout(60 * time.Second))

	// Set up CORS
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Request-ID"},
		ExposedHeaders:   []string{"Link", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Register example routes
	// Commented out for now as we don't have the example controller
	// errorExampleHandler.RegisterRoutes(s.router)

	// Add a simple health check endpoint
	s.router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok","time":"%s"}`, time.Now().Format(time.RFC3339))
	})

	return nil
}

// Start starts the example server
func (s *ExampleServer) Start(port int) error {
	addr := fmt.Sprintf(":%d", port)
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	s.logger.Info().Int("port", port).Msg("Starting example server with unified error handling")
	return s.httpServer.ListenAndServe()
}

// Stop stops the example server
func (s *ExampleServer) Stop(ctx context.Context) error {
	s.logger.Info().Msg("Stopping example server")
	return s.httpServer.Shutdown(ctx)
}
