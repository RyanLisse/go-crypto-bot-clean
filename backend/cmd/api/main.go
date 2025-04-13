// Package main provides the entry point for the API server.
package main

import (
	"log"
	"net/http"
	"os"

	"go.uber.org/zap"

	"go-crypto-bot-clean/backend/internal/api"
	"go-crypto-bot-clean/backend/internal/api/database"
	"go-crypto-bot-clean/backend/internal/api/repository"
	"go-crypto-bot-clean/backend/internal/api/service"
	internalAuth "go-crypto-bot-clean/backend/internal/auth" // Use internal/auth
	"go-crypto-bot-clean/backend/internal/config"

	// "go-crypto-bot-clean/backend/pkg/auth" // Removed pkg/auth import
	"go-crypto-bot-clean/backend/pkg/backtest"
	"go-crypto-bot-clean/backend/pkg/strategy"
)

func main() {
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	// Load configuration first
	cfg, err := config.LoadConfig(".") // Assuming config is in the current dir or handled by LoadConfig
	if err != nil {
		logger.Warn("Error loading config, using default development configuration", zap.Error(err))
		// Create a default development configuration
		cfg = &config.Config{}
		cfg.App.Environment = "development"
		cfg.App.Debug = true
		cfg.Database.Path = "./data/dev.db"
		cfg.Database.MaxIdleConns = 5
		cfg.Database.MaxOpenConns = 10
	}

	// Ensure environment is set
	if cfg.App.Environment == "" {
		cfg.App.Environment = "development"
		logger.Info("Environment not set, defaulting to development")
	}

	// Initialize database
	dbConfig := database.DefaultConfig()
	dbConfig.Path = cfg.Database.Path // Use config value
	dbConfig.Debug = cfg.App.Debug    // Use config value from App struct
	dbConfig.Logger = logger

	db, err := database.NewDatabase(dbConfig)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Run migrations
	migrationManager := database.NewMigrationManager(db, logger)
	if err := migrationManager.RunMigrations(); err != nil {
		logger.Fatal("Failed to run migrations", zap.Error(err))
	}

	// Seed database with initial data
	if getEnvBool("DB_SEED", false) { // Use environment variable check
		if err := migrationManager.SeedDatabase(); err != nil {
			logger.Fatal("Failed to seed database", zap.Error(err))
		}
	}

	// Initialize repositories
	userRepo := repository.NewGormUserRepository(db)
	strategyRepo := repository.NewGormStrategyRepository(db)
	backtestRepo := repository.NewGormBacktestRepository(db)

	// Initialize services and dependencies
	deps, err := api.NewDependencies(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to initialize dependencies", zap.Error(err))
	}

	// Initialize handlers
	deps.InitializeStatusHandler()
	deps.InitializePortfolioHandler()
	deps.InitializeNewCoinDependencies()
	deps.InitializeAnalyticsDependencies()

	// Create router using the consolidated router
	router := api.SetupConsolidatedRouter(deps)

	// Initialize services
	backtestService := backtest.NewService()
	strategyFactory := strategy.NewFactory()
	// Use internal/auth and config key
	var authProvider internalAuth.AuthProvider
	if cfg.Auth.Enabled {
		authProvider = internalAuth.NewService(cfg.Auth.ClerkSecretKey)
	} else {
		authProvider = internalAuth.NewDisabledService() // Use disabled service if auth is off
	}

	// JWT service and authentication middleware are initialized in the consolidated router

	// No need to register protected routes here, they are registered in the consolidated router

	// Create service provider (not used for now, but kept for future use)
	_ = service.NewProvider(
		&backtestService, // Convert to pointer to interface
		&strategyFactory, // Convert to pointer to interface
		authProvider,     // Pass the internal/auth provider (already an interface)
		userRepo,
		strategyRepo,
		backtestRepo,
	)

	// Huma API and health check endpoint are set up in the consolidated router

	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("Starting server on :%s...", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvBool gets a boolean environment variable or returns a default value
func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value == "true" || value == "1" || value == "yes"
}
