package app

import (
	"context"
	"fmt"
	"log"

	"go.uber.org/zap"

	"go-crypto-bot-clean/railway/internal/config"
	"go-crypto-bot-clean/railway/internal/server"
)

// Application represents the main application
type Application struct {
	config *config.Config
	logger *zap.Logger
	server *server.Server
}

// New creates a new Application instance
func New(cfg *config.Config) (*Application, error) {
	// Initialize logger
	logger, err := createLogger(cfg.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	// Create the application
	app := &Application{
		config: cfg,
		logger: logger,
	}

	// Initialize server
	app.server = server.New(cfg, logger)

	return app, nil
}

// Start starts the application server
func (a *Application) Start(port string) error {
	a.logger.Info("Starting application server", zap.String("port", port))
	return a.server.Start(port)
}

// Stop gracefully shuts down the application
func (a *Application) Stop(ctx context.Context) error {
	a.logger.Info("Stopping application server")
	return a.server.Stop(ctx)
}

// createLogger creates a new zap logger
func createLogger(level string) (*zap.Logger, error) {
	// Parse log level
	var zapLevel zap.AtomicLevel
	if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
		log.Printf("Invalid log level %q, defaulting to info", level)
		zapLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	// Create logger configuration
	config := zap.Config{
		Level:       zapLevel,
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
		EncoderConfig: zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	// Build logger
	return config.Build()
}
