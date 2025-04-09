package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"go-crypto-bot-clean/railway/internal/config"
	"go-crypto-bot-clean/railway/internal/handlers"
)

// Server represents the HTTP server
type Server struct {
	config     *config.Config
	logger     *zap.Logger
	httpServer *http.Server
	router     chi.Router
}

// New creates a new Server instance
func New(cfg *config.Config, logger *zap.Logger) *Server {
	// Create router
	router := chi.NewRouter()

	// Create server
	server := &Server{
		config: cfg,
		logger: logger,
		router: router,
	}

	// Set up middleware and routes
	server.setupMiddleware()
	server.setupRoutes()

	return server
}

// setupMiddleware configures the middleware stack
func (s *Server) setupMiddleware() {
	// Standard middleware
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Timeout(60 * time.Second))

	// CORS middleware
	s.router.Use(middleware.SetHeader("Access-Control-Allow-Origin", "*"))
	s.router.Use(middleware.SetHeader("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS"))
	s.router.Use(middleware.SetHeader("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token"))
}

// setupRoutes configures the API routes
func (s *Server) setupRoutes() {
	// Create handlers
	healthHandler := handlers.NewHealthHandler()

	// Public routes
	s.router.Get("/health", healthHandler.HealthCheck)
	
	// API routes
	s.router.Route("/api/v1", func(r chi.Router) {
		// Add API routes here
		r.Get("/status", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"status":"ok"}`))
		})
	})
}

// Start starts the HTTP server
func (s *Server) Start(port string) error {
	// Create HTTP server
	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: s.router,
	}

	// Start server
	s.logger.Info("Starting server", zap.String("address", s.httpServer.Addr))
	return s.httpServer.ListenAndServe()
}

// Stop gracefully shuts down the server
func (s *Server) Stop(ctx context.Context) error {
	if s.httpServer != nil {
		s.logger.Info("Shutting down server")
		return s.httpServer.Shutdown(ctx)
	}
	return nil
}
