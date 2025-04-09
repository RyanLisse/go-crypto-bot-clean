// Package main provides the entry point for the API server.
package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"go-crypto-bot-clean/api/database"
	"go-crypto-bot-clean/api/huma"
	"go-crypto-bot-clean/api/middleware"
	"go-crypto-bot-clean/api/middleware/jwt"
	"go-crypto-bot-clean/api/repository"
	"go-crypto-bot-clean/api/service"
	"go-crypto-bot-clean/backend/internal/auth"
	"go-crypto-bot-clean/backend/internal/backtest"
	"go-crypto-bot-clean/backend/internal/domain/strategy"
)

func main() {
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	// Initialize database
	dbConfig := database.DefaultConfig()
	dbConfig.Path = getEnv("DB_PATH", dbConfig.Path)
	dbConfig.Debug = getEnvBool("DB_DEBUG", false)
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
	if getEnvBool("DB_SEED", false) {
		if err := migrationManager.SeedDatabase(); err != nil {
			logger.Fatal("Failed to seed database", zap.Error(err))
		}
	}

	// Initialize repositories
	userRepo := repository.NewGormUserRepository(db)
	strategyRepo := repository.NewGormStrategyRepository(db)
	backtestRepo := repository.NewGormBacktestRepository(db)

	// Create router
	router := chi.NewRouter()

	// Add middleware
	router.Use(chimiddleware.Logger)
	router.Use(chimiddleware.Recoverer)
	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.RealIP)

	// Initialize services
	backtestService := backtest.NewService()
	strategyFactory := strategy.NewStrategyFactory()
	authService := auth.NewDisabledService() // Use disabled auth service for development

	// Initialize JWT service
	accessSecret := getEnv("JWT_ACCESS_SECRET", "default-access-secret")
	refreshSecret := getEnv("JWT_REFRESH_SECRET", "default-refresh-secret")
	accessTTL := time.Hour
	refreshTTL := time.Hour * 24 * 7 // 7 days
	issuer := getEnv("JWT_ISSUER", "go-crypto-bot")

	jwtService := jwt.NewService(accessSecret, refreshSecret, accessTTL, refreshTTL, issuer)

	// Create authentication middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	// Register protected routes
	router.Group(func(r chi.Router) {
		// Apply authentication middleware to protected routes
		r.Use(authMiddleware.Authenticate)

		// Protected routes go here
		// For example:
		// r.Mount("/api/v1/user", userRouter)
		// r.Mount("/api/v1/strategy", strategyRouter)
		// r.Mount("/api/v1/backtest", backtestRouter)
	})

	// Create service provider
	serviceProvider := service.NewProvider(
		backtestService,
		strategyFactory,
		authService,
		userRepo,
		strategyRepo,
		backtestRepo,
	)

	// Setup Huma API
	config := huma.DefaultConfig()
	huma.SetupHuma(router, config, serviceProvider)

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
