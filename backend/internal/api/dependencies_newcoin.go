package api

import (
	"go-crypto-bot-clean/backend/internal/api/handlers"
	"go-crypto-bot-clean/backend/internal/core/newcoin"
	"go-crypto-bot-clean/backend/internal/platform/mexc/rest"

	"go.uber.org/zap"
)

// InitializeNewCoinDependencies initializes the NewCoin dependencies
func (d *Dependencies) InitializeNewCoinDependencies() {
	// Use the logger from the dependencies
	logger := d.logger

	// Check if we're in development mode
	if d.Config.App.Environment == "development" {
		logger.Info("Using mock NewCoin service for development mode")
		// Use our MockNewCoinService
		mockService := &mockNewCoinService{}
		d.NewCoinService = mockService
		d.NewCoinHandler = handlers.NewNewCoinsHandler(mockService)
		d.CoinHandler = handlers.NewCoinHandler(nil, mockService)
		return
	}

	// Use the repositories from the database dependencies
	if d.NewCoinRepository == nil {
		logger.Error("NewCoinRepository is nil, falling back to mock service")
		// Fall back to mock service
		mockService := &mockNewCoinService{}
		d.NewCoinService = mockService
		d.NewCoinHandler = handlers.NewNewCoinsHandler(mockService)
		d.CoinHandler = handlers.NewCoinHandler(nil, mockService)
		return
	}

	// Create MEXC client
	mexcClient, err := rest.NewClient(d.Config.Mexc.APIKey, d.Config.Mexc.SecretKey)
	if err != nil {
		logger.Error("Failed to create MEXC client", zap.Error(err))
		// Fall back to mock service
		mockService := &mockNewCoinService{}
		d.NewCoinService = mockService
		d.NewCoinHandler = handlers.NewNewCoinsHandler(mockService)
		d.CoinHandler = handlers.NewCoinHandler(nil, mockService)
		return
	}

	// Create NewCoin service using GORM repository
	newCoinService := newcoin.NewGORMNewCoinService(mexcClient, d.NewCoinRepository, logger)

	// Store the service in dependencies
	d.NewCoinService = newCoinService

	// Create NewCoin handler
	d.NewCoinHandler = handlers.NewNewCoinsHandler(newCoinService)

	// Create CoinHandler for market and tradable coin endpoints
	d.CoinHandler = handlers.NewCoinHandler(nil, newCoinService)
}

// We're using the mockNewCoinService from server.go
