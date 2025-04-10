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

	"go-crypto-bot-clean/backend/internal/api"
	"go-crypto-bot-clean/backend/internal/api/websocket"
	"go-crypto-bot-clean/backend/internal/config"
	"go-crypto-bot-clean/backend/internal/repository/sqlite"
	"go-crypto-bot-clean/backend/internal/services/gemini"
	"go-crypto-bot-clean/backend/internal/services/reporting"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

// BotApp represents the main application
type BotApp struct {
	config  *config.Config
	logger  *zap.Logger
	db      *sql.DB
	server  *http.Server
	wsHub   *websocket.Hub
	handler http.Handler
	deps    *api.Dependencies
}

// NewBotApp creates a new BotApp
func NewBotApp(cfg *config.Config, logger *zap.Logger) *BotApp {
	return &BotApp{
		config: cfg,
		logger: logger,
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

	// Initialize dependencies
	deps, err := api.NewDependencies(a.config)
	if err != nil {
		return fmt.Errorf("failed to initialize dependencies: %w", err)
	}
	a.deps = deps

	// Setup consolidated Chi router with Huma integration
	a.handler = api.SetupConsolidatedRouter(deps)

	return nil
}

// Run starts the application
func (a *BotApp) Run(ctx context.Context) error {
	// Create server
	a.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", 8081),
		Handler: a.handler,
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
