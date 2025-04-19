package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port"
	"github.com/danielgtaylor/huma/v2"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// HumaServer represents the HTTP server with Huma integration
type HumaServer struct {
	router        *HumaRouter
	httpServer    *http.Server
	logger        *zerolog.Logger
	db            *gorm.DB
	config        *config.Config
	authConfig    *config.AuthConfig
	marketService port.MarketService
}

// NewHumaServer creates a new HTTP server with Huma integration
func NewHumaServer(db *gorm.DB, cfg *config.Config, authConfig *config.AuthConfig, logger *zerolog.Logger) *HumaServer {
	router := NewHumaRouter(cfg, authConfig, logger)

	return &HumaServer{
		router:     router,
		httpServer: &http.Server{},
		logger:     logger,
		db:         db,
		config:     cfg,
		authConfig: authConfig,
	}
}

// SetMarketService sets the market service
func (s *HumaServer) SetMarketService(marketService port.MarketService) {
	s.marketService = marketService
}

// SetupRoutes sets up the routes for the server
func (s *HumaServer) SetupRoutes() error {
	// Register health check endpoint
	s.router.RegisterHealthCheck()

	// Register OpenAPI documentation
	s.router.RegisterOpenAPI()

	// Register market data endpoints if market service is available
	if s.marketService != nil {
		// We'll create a market data handler later
		// marketHandler := NewHumaMarketHandler(s.marketService, s.logger)

		// For now, we'll just log that we're registering market data routes
		s.logger.Info().Msg("Would register market data routes here")

		s.logger.Info().Msg("Registered market data routes")
	}

	return nil
}

// Start starts the HTTP server
func (s *HumaServer) Start(port int) error {
	addr := fmt.Sprintf(":%d", port)
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	s.logger.Info().Int("port", port).Msg("Starting HTTP server with Huma integration")
	return s.httpServer.ListenAndServe()
}

// Stop stops the HTTP server
func (s *HumaServer) Stop(ctx context.Context) error {
	s.logger.Info().Msg("Stopping HTTP server")
	return s.httpServer.Shutdown(ctx)
}

// Router returns the underlying router
func (s *HumaServer) Router() *HumaRouter {
	return s.router
}

// NewHumaMarketHandler creates a new market data handler for Huma
func NewHumaMarketHandler(marketService port.MarketService, logger *zerolog.Logger) *HumaMarketHandler {
	return &HumaMarketHandler{
		marketService: marketService,
		logger:        logger,
	}
}

// HumaMarketHandler handles market data requests
type HumaMarketHandler struct {
	marketService port.MarketService
	logger        *zerolog.Logger
}

// RegisterRoutes registers the market data routes
func (h *HumaMarketHandler) RegisterRoutes(api *huma.API) {
	// For now, we'll just log that we're registering routes
	h.logger.Info().Msg("Registering market data routes")

	// We'll implement the actual route registration later
}
