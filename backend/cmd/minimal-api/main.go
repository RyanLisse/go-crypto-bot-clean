// Package main provides a minimal API server for Railway deployment testing
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"go-crypto-bot-clean/backend/internal/config"
	"go-crypto-bot-clean/backend/internal/database"
	"go-crypto-bot-clean/backend/internal/health"
	"go-crypto-bot-clean/backend/internal/models"
	"go-crypto-bot-clean/backend/internal/repositories"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
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

	// Initialize database if enabled
	var dbManager *database.SQLiteManager
	var repo *repositories.MinimalRepository

	// Force database enablement for testing
	cfg.Database.Enabled = true

	// Set default values if not provided
	if cfg.Database.Path == "" {
		cfg.Database.Path = "./data/minimal.db"
	}

	if cfg.Database.MaxOpenConns <= 0 {
		cfg.Database.MaxOpenConns = 10
	}

	if cfg.Database.MaxIdleConns <= 0 {
		cfg.Database.MaxIdleConns = 5
	}

	if cfg.Database.ConnMaxLifetimeSeconds <= 0 {
		cfg.Database.ConnMaxLifetimeSeconds = 300
	}

	if cfg.Database.Enabled {
		// Create database manager
		dbManager = database.NewSQLiteManager(database.SQLiteConfig{
			Path:                   cfg.Database.Path,
			MaxOpenConns:           cfg.Database.MaxOpenConns,
			MaxIdleConns:           cfg.Database.MaxIdleConns,
			ConnMaxLifetimeSeconds: cfg.Database.ConnMaxLifetimeSeconds,
			Debug:                  cfg.App.Debug,
		}, logger)

		// Connect to database
		if err := dbManager.Connect(); err != nil {
			logger.Fatal("Failed to connect to database", zap.Error(err))
		}
		defer dbManager.Close()

		// Auto migrate models
		if err := dbManager.AutoMigrate(
			&models.SystemInfo{},
			&models.HealthCheck{},
			&models.LogEntry{},
		); err != nil {
			logger.Fatal("Failed to run auto migration", zap.Error(err))
		}

		// Create repository
		repo = repositories.NewMinimalRepository(dbManager.DB(), logger)

		// Save system info
		systemInfo := &models.SystemInfo{
			Name:        cfg.App.Name,
			Version:     "0.1.0",
			Environment: cfg.App.Environment,
			StartTime:   time.Now(),
		}
		if err := repo.SaveSystemInfo(systemInfo); err != nil {
			logger.Error("Failed to save system info", zap.Error(err))
		}

		logger.Info("Database initialized successfully",
			zap.String("path", cfg.Database.Path),
			zap.Bool("debug", cfg.App.Debug),
		)
	} else {
		logger.Info("Database is disabled")
	}

	// Create router
	router := chi.NewRouter()

	// Add middleware
	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.RealIP)
	router.Use(chimiddleware.Logger)
	router.Use(chimiddleware.Recoverer)

	// Initialize health check
	healthCheck := health.NewHealthCheck("0.1.0", logger)

	// Add system component to health check
	healthCheck.AddComponent("system", health.StatusUp, "System is running")

	// Add database component to health check if enabled
	if cfg.Database.Enabled && dbManager != nil {
		healthCheck.AddComponent("database", health.StatusUp, "Database is connected")
	}

	// Add health check endpoints
	router.Get("/health", healthCheck.SimpleHandler())
	router.Get("/health/detailed", healthCheck.Handler())

	// Add version endpoint
	router.Get("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"version":     "0.1.0",
			"environment": cfg.App.Environment,
		})
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
			"database": map[string]any{
				"enabled": cfg.Database.Enabled,
				"path":    cfg.Database.Path,
			},
		}

		json.NewEncoder(w).Encode(safeConfig)
	})

	// Add database endpoints if enabled
	if cfg.Database.Enabled && repo != nil {
		// Add system info endpoint
		router.Get("/system", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			info, err := repo.GetSystemInfo()
			if err != nil {
				logger.Error("Failed to get system info", zap.Error(err))
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			if info == nil {
				http.Error(w, "System info not found", http.StatusNotFound)
				return
			}

			// Update uptime
			info.Uptime = int64(time.Since(info.StartTime).Seconds())

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(info)
		})

		// Add health checks endpoint
		router.Get("/health/history", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			limit := 10 // Default limit
			checks, err := repo.GetHealthChecks(limit)
			if err != nil {
				logger.Error("Failed to get health checks", zap.Error(err))
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(checks)
		})

		// Add log entry endpoint
		router.Post("/logs", func(w http.ResponseWriter, r *http.Request) {
			var entry models.LogEntry
			if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			// Set ID and timestamp
			entry.ID = uuid.New()
			entry.Timestamp = time.Now()

			if err := repo.SaveLogEntry(&entry); err != nil {
				logger.Error("Failed to save log entry", zap.Error(err))
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(entry)
		})

		// Add logs endpoint
		router.Get("/logs", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			limit := 10 // Default limit
			level := r.URL.Query().Get("level")

			entries, err := repo.GetLogEntries(limit, level)
			if err != nil {
				logger.Error("Failed to get log entries", zap.Error(err))
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(entries)
		})
	}

	// Start server
	port := cfg.App.Port
	logger.Info("Starting minimal server",
		zap.String("port", port),
		zap.String("environment", cfg.App.Environment),
		zap.String("log_level", cfg.App.LogLevel),
		zap.Bool("database_enabled", cfg.Database.Enabled),
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
