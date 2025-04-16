package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/http/controller"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/http/middleware"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm/repo"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/factory"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// Server represents the HTTP server
type Server struct {
	router     chi.Router
	httpServer *http.Server
	logger     *zerolog.Logger
	db         *gorm.DB
	config     *config.Config
	authConfig *config.AuthConfig
}

// NewServer creates a new HTTP server
func NewServer(db *gorm.DB, cfg *config.Config, authConfig *config.AuthConfig, logger *zerolog.Logger) *Server {
	router := chi.NewRouter()

	return &Server{
		router:     router,
		httpServer: &http.Server{},
		logger:     logger,
		db:         db,
		config:     cfg,
		authConfig: authConfig,
	}
}

// SetupRoutes sets up the routes for the server
func (s *Server) SetupRoutes() error {
	// Create wallet service/controller
	mexcClient := factory.NewMEXCClient(s.config, s.logger)
	walletRepo := repo.NewConsolidatedWalletRepository(s.db, s.logger)
	walletService := usecase.NewWalletService(walletRepo, mexcClient, s.logger)
	walletController := controller.NewWalletController(walletService, s.logger)

	// Create a mock middleware
	authMiddleware := middleware.NewMockMiddleware(s.logger)

	// Set up middleware
	s.router.Use(chimiddleware.RequestID)
	s.router.Use(chimiddleware.RealIP)
	s.router.Use(chimiddleware.Logger)
	s.router.Use(chimiddleware.Recoverer)
	s.router.Use(chimiddleware.Timeout(60 * time.Second))

	// Set up CORS
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Set up authentication middleware
	s.router.Use(authMiddleware.Middleware())

	// Register routes
	walletController.RegisterRoutes(s.router)

	return nil
}

// Start starts the HTTP server
func (s *Server) Start(port int) error {
	addr := fmt.Sprintf(":%d", port)
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	s.logger.Info().Int("port", port).Msg("Starting HTTP server")
	return s.httpServer.ListenAndServe()
}

// Stop stops the HTTP server
func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info().Msg("Stopping HTTP server")
	return s.httpServer.Shutdown(ctx)
}
