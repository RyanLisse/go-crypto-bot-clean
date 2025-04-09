package api

import (
	"log"

	"github.com/ryanlisse/go-crypto-bot/internal/api/handlers"
	"github.com/ryanlisse/go-crypto-bot/internal/core/newcoin"
	"github.com/ryanlisse/go-crypto-bot/internal/platform/mexc/rest"
	"go.uber.org/zap"
)

// InitializeNewCoinDependencies initializes the NewCoin dependencies
func (d *Dependencies) InitializeNewCoinDependencies() {
	// Create logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Printf("Failed to create logger: %v", err)
		// Fall back to mock service
		mockService := &mockNewCoinService{}
		d.NewCoinHandler = handlers.NewNewCoinsHandler(mockService)
		d.CoinHandler = handlers.NewCoinHandler(nil, mockService)
		return
	}

	// Use the repositories from the database dependencies
	if d.NewCoinRepository == nil {
		logger.Error("NewCoinRepository is nil, falling back to mock service")
		// Fall back to mock service
		mockService := &mockNewCoinService{}
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
		d.NewCoinHandler = handlers.NewNewCoinsHandler(mockService)
		d.CoinHandler = handlers.NewCoinHandler(nil, mockService)
		return
	}

	// Create NewCoin service using GORM repository
	newCoinService := newcoin.NewGORMNewCoinService(mexcClient, d.NewCoinRepository, logger)

	// Create NewCoin handler
	d.NewCoinHandler = handlers.NewNewCoinsHandler(newCoinService)

	// Create CoinHandler for market and tradable coin endpoints
	d.CoinHandler = handlers.NewCoinHandler(nil, newCoinService)
}

// We're using the mockNewCoinService from server.go
