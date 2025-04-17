package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/factory"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// It's okay if .env doesn't exist, we'll just use environment variables
		fmt.Println("Warning: .env file not found, using environment variables")
	}

	// Initialize logger
	log := logger.NewLogger()
	log.Info().Msg("Starting crypto trading bot")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Initialize DB connection
	db := gorm.NewDB(cfg, log)

	// Run database migrations
	if err := gorm.AutoMigrateModels(db, log); err != nil {
		log.Fatal().Err(err).Msg("Failed to run database migrations")
	}

	// Create trade service factory
	tradeServiceFactory := factory.NewTradingServiceFactory(cfg, log, db)

	// Create trade factory
	tradeFactory := factory.NewTradeFactory(cfg, log, db)

	// Create dependencies for trade service
	marketFactory := factory.NewMarketFactory(cfg, log, db)
	mexcClient := marketFactory.CreateMEXCClient()
	marketDataService := marketFactory.CreateMarketDataService()
	_, symbolRepo := marketFactory.CreateMarketRepository()
	orderRepo := tradeFactory.CreateOrderRepository()
	// Create trade service
	tradeService := tradeFactory.CreateTradeService(mexcClient, marketDataService, symbolRepo, orderRepo)

	// Create trade executor factory
	tradeExecutorFactory := factory.NewTradeExecutorFactory(cfg, log)

	// Create trade executor
	tradeExecutor := tradeExecutorFactory.CreateTradeExecutor(tradeService)

	// Create trade history factory
	tradeHistoryFactory := factory.NewTradeHistoryFactory(cfg, log, db)

	// Create trade history repository
	tradeHistoryRepo := tradeHistoryFactory.CreateTradeHistoryRepository()

	// Create trading service
	tradingService, err := tradeServiceFactory.CreateTradingService(
		tradeExecutor,
		tradeHistoryRepo,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create trading service")
	}

	// Start the trading service
	if err := tradingService.Start(); err != nil {
		log.Fatal().Err(err).Msg("Failed to start trading service")
	}
	log.Info().Msg("Trading service started")

	// Wait for signal to shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Info().Msg("Shutdown signal received")

	// Create a context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Stop the trading service
	if err := tradingService.Stop(); err != nil {
		log.Error().Err(err).Msg("Error stopping trading service")
	}

	// Wait for context to be done (either timeout or clean shutdown)
	<-ctx.Done()
	if ctx.Err() == context.DeadlineExceeded {
		log.Warn().Msg("Shutdown timed out")
	} else {
		log.Info().Msg("Shutdown completed gracefully")
	}
}
