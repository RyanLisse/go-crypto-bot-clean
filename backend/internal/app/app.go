package app

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"go-crypto-bot-clean/backend/internal/api"
	"go-crypto-bot-clean/backend/internal/api/handlers"
	"go-crypto-bot-clean/backend/internal/api/websocket"
	"go-crypto-bot-clean/backend/internal/config"
	"go-crypto-bot-clean/backend/internal/core/status"
	"go-crypto-bot-clean/backend/internal/domain/repository"
	"go-crypto-bot-clean/backend/internal/repository/sqlite"
	"go-crypto-bot-clean/backend/internal/services/gemini"
	"go-crypto-bot-clean/backend/internal/services/reporting"
	"go.uber.org/zap"
)

// BotApp represents the main application
type BotApp struct {
	config             *config.Config
	logger             *zap.Logger
	db                 *sql.DB
	router             *gin.Engine
	server             *http.Server
	wsHub              *websocket.Hub
	statusHandler      *handlers.StatusHandler
	reportHandler      *handlers.ReportHandler
	wsHandler          *handlers.WebSocketHandler
	accountHandler     *handlers.AccountHandler
	tradeAnalyticsRepo repository.TradeAnalyticsRepository
	balanceHistoryRepo repository.BalanceHistoryRepository
}

// NewBotApp creates a new BotApp
func NewBotApp(cfg *config.Config, logger *zap.Logger) *BotApp {
	return &BotApp{
		config: cfg,
		logger: logger,
		router: gin.Default(),
	}
}

// Initialize initializes the application
func (a *BotApp) Initialize(ctx context.Context) error {
	// Connect to database
	db, err := sql.Open("sqlite3", a.config.Database.Path)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	a.db = db

	// Initialize WebSocket hub
	a.wsHub = websocket.NewHub()
	go a.wsHub.Run()

	// Initialize handlers and services
	if err := a.initializeHandlers(ctx); err != nil {
		return err
	}

	// Set up API routes
	a.setupRoutes()

	return nil
}

// Run starts the application
func (a *BotApp) Run(ctx context.Context) error {
	// Create server
	a.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", 8081),
		Handler: a.router,
	}

	// Channel to listen for errors coming from the server
	serverErrors := make(chan error, 1)

	// Start the server
	go func() {
		a.logger.Info("Starting server", zap.Int("port", 8081))
		serverErrors <- a.server.ListenAndServe()
	}()

	// Channel to listen for an interrupt or terminate signal
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Block until we receive a signal or an error
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case <-shutdown:
		a.logger.Info("Shutting down server")

		// Create a deadline for graceful shutdown
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		// Shutdown the server
		if err := a.server.Shutdown(ctx); err != nil {
			a.server.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}

// Close closes the application
func (a *BotApp) Close() error {
	if a.db != nil {
		if err := a.db.Close(); err != nil {
			return fmt.Errorf("failed to close database: %w", err)
		}
	}
	return nil
}

// initializeHandlers initializes the handlers
func (a *BotApp) initializeHandlers(ctx context.Context) error {
	// Initialize status handler
	statusProvider := status.NewMockStatusProvider()
	statusService := status.NewStatusService(statusProvider, "1.0.0")
	a.statusHandler = handlers.NewStatusHandler(statusService)

	// Initialize WebSocket handler
	a.wsHandler = handlers.NewWebSocketHandler(a.wsHub, a.logger)

	// Initialize account handler
	a.accountHandler = handlers.NewAccountHandler(nil, nil, a.logger)

	// Initialize report handler
	reportGenerator, err := a.setupReportGenerator(ctx)
	if err != nil {
		return fmt.Errorf("failed to set up report generator: %w", err)
	}
	a.reportHandler = handlers.NewReportHandler(reportGenerator, a.logger)

	return nil
}

// setupReportGenerator sets up the report generator
func (a *BotApp) setupReportGenerator(ctx context.Context) (*reporting.ReportGenerator, error) {
	// Create Gemini client
	geminiClient := gemini.NewGeminiClient(a.config.Gemini.APIKey)

	// Create report repository
	reportRepo := sqlite.NewReportRepository(a.db, a.logger)
	if err := reportRepo.Initialize(ctx); err != nil {
		return nil, fmt.Errorf("failed to initialize report repository: %w", err)
	}

	// Create report generator
	reportGenerator := reporting.NewReportGenerator(
		geminiClient,
		reportRepo,
		5*time.Minute,
		a.logger,
	)

	return reportGenerator, nil
}

// setupRoutes sets up the API routes
func (a *BotApp) setupRoutes() {
	// Set up routes
	api.SetupRoutes(
		a.router,
		a.statusHandler,
		a.reportHandler,
		a.wsHandler,
		a.accountHandler,
		a.logger,
	)
}
