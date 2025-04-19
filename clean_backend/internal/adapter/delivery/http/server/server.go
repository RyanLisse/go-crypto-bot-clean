package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	// Import the new handler package
	handler "github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/adapter/delivery/http/handler"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/adapter/delivery/http/handler/health"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/adapter/delivery/http/middleware"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/factory" // Added factory import
)

// Server represents the HTTP server
type Server struct {
	router          chi.Router
	httpServer      *http.Server
	logger          *zerolog.Logger
	db              *gorm.DB
	config          *config.Config
	useCaseFactory  *factory.UseCaseFactory // Added UseCaseFactory
	authMiddleware  middleware.AuthMiddleware
	errorMiddleware *middleware.UnifiedErrorMiddleware
}

// NewServer creates a new HTTP server
func NewServer(db *gorm.DB, cfg *config.Config, logger *zerolog.Logger, useCaseFactory *factory.UseCaseFactory) *Server { // Added useCaseFactory param
	router := chi.NewRouter()

	// Create default middleware (will be replaced by DI container via SetMiddleware)
	errorMiddleware := middleware.NewUnifiedErrorMiddleware(logger)
	authMiddleware := middleware.NewTestMiddleware(logger)

	return &Server{
		router:          router,
		httpServer:      &http.Server{},
		logger:          logger,
		db:              db,
		config:          cfg,
		useCaseFactory:  useCaseFactory, // Store UseCaseFactory
		authMiddleware:  authMiddleware,
		errorMiddleware: errorMiddleware,
	}
}

// SetupRoutes sets up the routes for the server
func (s *Server) SetupRoutes() error {
	// Set up middleware
	s.router.Use(chimiddleware.RequestID) // Ensure Chi's RequestID is used for GetTraceID
	s.router.Use(chimiddleware.RealIP)
	// Use the logging integrated within UnifiedErrorMiddleware instead of Chi's logger
	// s.router.Use(chimiddleware.Logger)
	s.router.Use(s.errorMiddleware.Middleware()) // Use the UnifiedErrorMiddleware
	// Note: UnifiedErrorMiddleware includes Recoverer functionality
	// s.router.Use(chimiddleware.Recoverer)
	s.router.Use(chimiddleware.Timeout(60 * time.Second))

	// Set up CORS
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // TODO: Restrict in production
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Instantiate handlers using UseCaseFactory
	healthHandler := health.NewHandler(s.logger)
	rootTestHandler := handler.NewRootTestHandler(s.logger)
	statusHandler := handler.NewStatusHandler(s.logger)

	// Instantiate handlers needing use cases
	protectedUseCase := s.useCaseFactory.BuildProtectedUserUseCase()
	protectedHandler := handler.NewProtectedHandler(s.logger, protectedUseCase)

	// sniperShotUseCase := s.useCaseFactory.BuildSniperShotService()
	// sniperHandler := handler.NewSniperHandler(s.logger, sniperShotUseCase)

	// walletUseCase := s.useCaseFactory.BuildWalletUseCase()
	// walletHandler := handler.NewWalletHandler(s.logger, walletUseCase)

	// Register root level routes
	s.router.Get("/root-test", rootTestHandler.HandleRootTest)
	healthHandler.RegisterRoutes(s.router)

	// API routes v1
	s.router.Route("/api/v1", func(r chi.Router) {
		// Public routes
		r.Group(func(r chi.Router) {
			r.Get("/status", statusHandler.GetStatus)
			// r.Post("/sniper", sniperHandler.HandleExecuteSnipe)
			// Add Wallet related public routes (if any) using walletHandler
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(s.authMiddleware.Middleware())
			r.Use(s.authMiddleware.RequireAuthentication)
			r.Get("/protected", protectedHandler.GetProtected)
			// Add other protected routes needing use cases here
			// Example: Wallet routes
			// r.Get("/wallets", walletHandler.HandleGetUserWallets)
			// r.Post("/wallets/web3", walletHandler.HandleCreateWeb3Wallet)
		})

		s.logger.Info().Msg("API routes registered at /api/v1")
	})

	return nil
}

// Start starts the HTTP server
func (s *Server) Start(port int) error {
	// Kill any process using the port before starting
	if err := killProcessOnPort(port, s.logger); err != nil {
		s.logger.Warn().Err(err).Int("port", port).Msg("Failed to kill process on port")
	}

	addr := fmt.Sprintf(":%d", port)
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	s.logger.Info().Int("port", port).Msg("Starting HTTP server")

	// Try to start the server with retries
	var err error
	for i := 0; i < 3; i++ {
		err = s.httpServer.ListenAndServe()
		// If the error is not about the address being in use, return it immediately
		if err != nil && !strings.Contains(err.Error(), "address already in use") {
			return err
		}

		// If we got here, the error is about the address being in use
		if i < 2 { // Don't log on the last attempt
			s.logger.Warn().Int("port", port).Int("attempt", i+1).Msg("Port is still in use, retrying after a short delay")
			// Try to kill the process again
			if err := killProcessOnPort(port, s.logger); err != nil {
				s.logger.Warn().Err(err).Int("port", port).Msg("Failed to kill process on port")
			}
			// Wait a bit before retrying
			time.Sleep(2 * time.Second)
		}
	}

	return err // Return the last error
}

// Stop stops the HTTP server
func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info().Msg("Stopping HTTP server")
	return s.httpServer.Shutdown(ctx)
}

// SetMiddleware sets the middleware for the server
func (s *Server) SetMiddleware(authMiddleware middleware.AuthMiddleware, errorMiddleware *middleware.UnifiedErrorMiddleware) {
	s.authMiddleware = authMiddleware
	s.errorMiddleware = errorMiddleware
	s.logger.Info().Msg("Middleware set for server")
}
