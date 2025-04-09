// Package main provides the entry point for the API server.
package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"go-crypto-bot-clean/api/huma"
	"go-crypto-bot-clean/api/service"
	"go-crypto-bot-clean/backend/internal/auth"
	"go-crypto-bot-clean/backend/internal/backtest"
	"go-crypto-bot-clean/backend/internal/domain/strategy"
)

func main() {
	// Create router
	router := chi.NewRouter()

	// Add middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Initialize services
	backtestService := backtest.NewService()
	strategyFactory := strategy.NewStrategyFactory()
	authService := auth.NewDisabledService() // Use disabled auth service for development

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
	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
