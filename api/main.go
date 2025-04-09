// Package main provides the entry point for the API server.
package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"go-crypto-bot-clean/api/huma"
	"go-crypto-bot-clean/api/middleware"
	"go-crypto-bot-clean/api/middleware/jwt"
	"go-crypto-bot-clean/api/service"
	"go-crypto-bot-clean/backend/internal/auth"
	"go-crypto-bot-clean/backend/internal/backtest"
	"go-crypto-bot-clean/backend/internal/domain/strategy"
)

func main() {
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
