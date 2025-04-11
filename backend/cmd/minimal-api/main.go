// Package main provides a minimal API server for Railway deployment testing
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"go-crypto-bot-clean/backend/internal/config"
	"go-crypto-bot-clean/backend/internal/health"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// Initialize logger
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := loggerConfig.Build()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Determine environment
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}

	var environment config.Environment
	switch env {
	case "production":
		environment = config.EnvironmentProduction
	case "staging":
		environment = config.EnvironmentStaging
	default:
		environment = config.EnvironmentDevelopment
	}

	// Initialize configuration manager
	configManager := config.NewManager(logger, environment)

	// Load minimal configuration
	cfg, err := configManager.LoadMinimalConfig()
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Enable configuration reloading
	if err := configManager.EnableReload(); err != nil {
		logger.Warn("Failed to enable configuration reloading", zap.Error(err))
	} else {
		logger.Info("Configuration reloading enabled")
	}

	// Configure logger based on config
	configureLogger(logger, cfg)

	// Create router
	router := chi.NewRouter()

	// Add middleware
	router.Use(chimiddleware.Logger)
	router.Use(chimiddleware.Recoverer)
	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.RealIP)

	// Initialize health check
	healthCheck := health.NewHealthCheck("0.1.0", logger)

	// Add system component to health check
	healthCheck.AddComponent("system", health.StatusUp, "System is running")

	// Add health check endpoints
	router.Get("/health", healthCheck.SimpleHandler())
	router.Get("/health/detailed", healthCheck.Handler())

	// Add version endpoint
	router.Get("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"version": "0.1.0", "name": "Go Crypto Bot"}`))
	})

	// Add root endpoint
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Go Crypto Bot API"))
	})

	// Add config endpoint
	router.Get("/config", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Create a safe version of the config (without sensitive data)
		safeConfig := map[string]any{
			"app": map[string]any{
				"name":        cfg.App.Name,
				"environment": cfg.App.Environment,
				"debug":       cfg.App.Debug,
			},
			"logging": map[string]any{
				"filePath":   cfg.Logging.FilePath,
				"maxSize":    cfg.Logging.MaxSize,
				"maxBackups": cfg.Logging.MaxBackups,
				"maxAge":     cfg.Logging.MaxAge,
			},
		}

		json.NewEncoder(w).Encode(safeConfig)
	})

	// Start server
	port := cfg.App.Port
	logger.Info("Starting minimal server",
		zap.String("port", port),
		zap.String("environment", cfg.App.Environment),
		zap.String("log_level", cfg.App.LogLevel),
	)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		logger.Fatal("Server failed to start", zap.Error(err))
	}
}

// configureLogger configures the logger based on configuration
func configureLogger(logger *zap.Logger, cfg *config.MinimalConfig) {
	// Set log level from configuration
	var level zapcore.Level
	switch cfg.App.LogLevel {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// Log the current configuration
	logger.Info("Logger configured",
		zap.String("level", level.String()),
		zap.String("environment", cfg.App.Environment),
		zap.Bool("debug", cfg.App.Debug),
	)
}
