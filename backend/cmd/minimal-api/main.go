// Package main provides a minimal API server for Railway deployment testing
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"go-crypto-bot-clean/backend/internal/config"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// Load configuration
	cfg, err := config.LoadMinimalConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger := initLogger(cfg)
	defer logger.Sync()

	// Create router
	router := chi.NewRouter()

	// Add middleware
	router.Use(chimiddleware.Logger)
	router.Use(chimiddleware.Recoverer)
	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.RealIP)

	// Add health check endpoint
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

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
		safeConfig := map[string]interface{}{
			"app": map[string]interface{}{
				"name":        cfg.App.Name,
				"environment": cfg.App.Environment,
				"debug":       cfg.App.Debug,
			},
			"logging": map[string]interface{}{
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

// initLogger initializes the zap logger based on configuration
func initLogger(cfg *config.MinimalConfig) *zap.Logger {
	// Determine log level from config
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

	// Create logger config
	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(level),
		Development: cfg.App.Environment == "development",
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	// Create logger
	logger, err := config.Build()
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		os.Exit(1)
	}

	return logger
}
