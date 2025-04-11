// Package main provides a minimal API server for Railway deployment testing
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go-crypto-bot-clean/backend/internal/auth"
	"go-crypto-bot-clean/backend/internal/config"
	"go-crypto-bot-clean/backend/internal/database"
	"go-crypto-bot-clean/backend/internal/health"
	"go-crypto-bot-clean/backend/internal/middleware"
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

	// Initialize authentication if enabled
	var clerkAuth *auth.ClerkAuth
	if cfg.Auth.Enabled {
		clerkAuth, err = auth.FromMinimalConfig(cfg, logger)
		if err != nil {
			logger.Error("Failed to initialize authentication", zap.Error(err))
		} else {
			logger.Info("Authentication initialized",
				zap.Bool("enabled", cfg.Auth.Enabled),
				zap.String("domain", cfg.Auth.ClerkDomain),
			)
		}
	} else {
		logger.Info("Authentication is disabled")
	}

	// Initialize database if enabled
	var sqliteManager *database.SQLiteManager
	var tursoManager *database.TursoManager
	var syncManager *database.SyncManager
	var backupManager *database.BackupManager
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

	// Initialize SQLite database
	if cfg.Database.Enabled {
		// Create database manager
		sqliteManager = database.NewSQLiteManager(database.SQLiteConfig{
			Path:                   cfg.Database.Path,
			MaxOpenConns:           cfg.Database.MaxOpenConns,
			MaxIdleConns:           cfg.Database.MaxIdleConns,
			ConnMaxLifetimeSeconds: cfg.Database.ConnMaxLifetimeSeconds,
			Debug:                  cfg.App.Debug,
		}, logger)

		// Connect to database
		if err := sqliteManager.Connect(); err != nil {
			logger.Fatal("Failed to connect to database", zap.Error(err))
		}
		defer sqliteManager.Close()

		// Auto migrate models
		if err := sqliteManager.AutoMigrate(
			&models.SystemInfo{},
			&models.HealthCheck{},
			&models.LogEntry{},
		); err != nil {
			logger.Fatal("Failed to run auto migration", zap.Error(err))
		}

		// Create repository
		repo = repositories.NewMinimalRepository(sqliteManager.DB(), logger)

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

		// Initialize Turso if enabled
		if cfg.Database.Turso.Enabled {
			// Create Turso manager
			tursoManager = database.NewTursoManager(database.TursoConfig{
				Enabled:             cfg.Database.Turso.Enabled,
				URL:                 cfg.Database.Turso.URL,
				AuthToken:           cfg.Database.Turso.AuthToken,
				SyncEnabled:         cfg.Database.Turso.SyncEnabled,
				SyncIntervalSeconds: cfg.Database.Turso.SyncIntervalSeconds,
			}, logger)

			// Connect to Turso
			if err := tursoManager.Connect(context.Background()); err != nil {
				logger.Error("Failed to connect to Turso database", zap.Error(err))
			} else {
				defer tursoManager.Close()
				logger.Info("Connected to Turso database",
					zap.String("url", cfg.Database.Turso.URL),
					zap.Bool("syncEnabled", cfg.Database.Turso.SyncEnabled),
				)

				// Initialize sync manager if both SQLite and Turso are connected
				if cfg.Database.Turso.SyncEnabled {
					syncManager = database.NewSyncManager(database.SyncConfig{
						Enabled:             cfg.Database.Turso.SyncEnabled,
						SyncIntervalSeconds: cfg.Database.Turso.SyncIntervalSeconds,
						BatchSize:           cfg.Database.Turso.BatchSize,
						MaxRetries:          cfg.Database.Turso.MaxRetries,
						RetryDelaySeconds:   cfg.Database.Turso.RetryDelaySeconds,
					}, sqliteManager.DB(), tursoManager.DB(), logger)

					// Start sync
					if err := syncManager.StartSync(); err != nil {
						logger.Error("Failed to start database synchronization", zap.Error(err))
					} else {
						logger.Info("Database synchronization started")
					}

					// We'll add Turso component to health check later
				}
			}
		}

		logger.Info("Database initialized successfully",
			zap.String("path", cfg.Database.Path),
			zap.Bool("debug", cfg.App.Debug),
			zap.Bool("tursoEnabled", cfg.Database.Turso.Enabled),
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
	if cfg.Database.Enabled && sqliteManager != nil {
		healthCheck.AddComponent("database", health.StatusUp, "Database is connected")
	}

	// Add Turso component to health check if enabled
	if cfg.Database.Turso.Enabled && tursoManager != nil {
		healthCheck.AddComponent("turso", health.StatusUp, "Turso database is connected")
	}

	// Add sync component to health check if enabled
	if cfg.Database.Turso.SyncEnabled && syncManager != nil {
		healthCheck.AddComponent("sync", health.StatusUp, "Database synchronization is active")
	}

	// Add backup component to health check if enabled
	if backupManager != nil {
		healthCheck.AddComponent("backup", health.StatusUp, "Database backup system is active")
	}

	// Add authentication component to health check if enabled
	if cfg.Auth.Enabled && clerkAuth != nil {
		healthCheck.AddComponent("auth", health.StatusUp, "Authentication is enabled")
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
				"turso": map[string]any{
					"enabled":             cfg.Database.Turso.Enabled,
					"syncEnabled":         cfg.Database.Turso.SyncEnabled,
					"syncIntervalSeconds": cfg.Database.Turso.SyncIntervalSeconds,
				},
			},
			"auth": map[string]any{
				"enabled":     cfg.Auth.Enabled,
				"clerkDomain": cfg.Auth.ClerkDomain,
			},
		}

		json.NewEncoder(w).Encode(safeConfig)
	})

	// Add sync status endpoint if enabled
	if cfg.Database.Turso.SyncEnabled && syncManager != nil {
		router.Get("/sync/status", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(syncManager.GetSyncStatus())
		})
	}

	// Initialize backup manager
	backupManager, err = database.NewBackupManager(database.BackupConfig{
		Enabled:          true,
		BackupDir:        "./data/backups",
		BackupInterval:   24 * time.Hour,
		MaxBackups:       7,
		RetentionDays:    30,
		CompressBackups:  true,
		IncludeTimestamp: true,
	}, sqliteManager.DB(), tursoManager.DB(), logger)
	if err != nil {
		logger.Error("Failed to initialize backup manager", zap.Error(err))
	} else {
		// Add backup endpoint
		router.Post("/backup", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			err := backupManager.BackupDatabases(r.Context())
			if err != nil {
				logger.Error("Backup failed", zap.Error(err))
				http.Error(w, fmt.Sprintf("Backup failed: %v", err), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "success",
				"message": "Backup completed successfully",
			})
		})
	}

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

	// Add authentication endpoints if enabled
	if cfg.Auth.Enabled && clerkAuth != nil {
		// Create a protected router group
		protected := chi.NewRouter()
		protected.Use(middleware.RequireAuthMiddleware(clerkAuth, logger))

		// Add protected endpoints
		protected.Get("/me", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			// Get user from context
			token, err := clerkAuth.GetUserFromContext(r.Context())
			if err != nil {
				logger.Error("Failed to get user from context", zap.Error(err))
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Get user ID
			userID, err := clerkAuth.GetUserID(token)
			if err != nil {
				logger.Error("Failed to get user ID", zap.Error(err))
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Get user profile
			user, err := clerkAuth.GetUserProfile(r.Context(), userID)
			if err != nil {
				logger.Error("Failed to get user profile", zap.Error(err))
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(user)
		})

		// Mount protected router
		router.Mount("/api", protected)
	}

	// Start server
	port := cfg.App.Port
	logger.Info("Starting minimal server",
		zap.String("port", port),
		zap.String("environment", cfg.App.Environment),
		zap.String("log_level", cfg.App.LogLevel),
		zap.Bool("database_enabled", cfg.Database.Enabled),
		zap.Bool("auth_enabled", cfg.Auth.Enabled),
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
