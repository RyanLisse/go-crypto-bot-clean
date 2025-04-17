package main

import (
	"context"
	"log"
	"os"
	"time"

	gormadapter "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/factory"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

func init() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// It's okay if .env doesn't exist, we'll just use environment variables
		log.Println("Warning: .env file not found, using environment variables")
	}
}

func main() {
	// Initialize logger
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	logger.Info().Msg("Starting symbol sync")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database using GORM adapter
	db, err := gormadapter.NewDBConnection(cfg, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
	}

	// Create market factory and MEXC client
	marketFactory := factory.NewMarketFactory(cfg, &logger, db)
	mexcClient := marketFactory.CreateMEXCClient()

	// Get exchange info
	ctx := context.Background()
	exchangeInfo, err := mexcClient.GetExchangeInfo(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to get exchange info")
	}

	logger.Info().Int("symbolCount", len(exchangeInfo.Symbols)).Msg("Got exchange info")

	// Create repositories using factory
	repoFactory := factory.NewRepositoryFactory(db, &logger, cfg)
	symbolRepo := repoFactory.CreateSymbolRepository()

	// Save symbols to database
	for _, symbol := range exchangeInfo.Symbols {
		// Convert to domain model
		domainSymbol := &model.Symbol{
			Symbol:            symbol.Symbol,
			BaseAsset:         symbol.BaseAsset,
			QuoteAsset:        symbol.QuoteAsset,
			Status:            model.SymbolStatus(symbol.Status),
			Exchange:          "mexc",
			MinPrice:          0, // TODO: Get from symbol info if available
			MaxPrice:          0, // TODO: Get from symbol info if available
			PricePrecision:    symbol.PricePrecision,
			MinQuantity:       0, // TODO: Get from symbol info if available
			MaxQuantity:       0, // TODO: Get from symbol info if available
			QuantityPrecision: symbol.QuantityPrecision,
			AllowedOrderTypes: []string{"LIMIT", "MARKET"},
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}

		// Save to database using canonical model.Symbol
		err := symbolRepo.Create(ctx, domainSymbol)
		if err != nil {
			logger.Error().Err(err).Str("symbol", symbol.Symbol).Msg("Failed to save symbol")
			continue
		}

		logger.Info().Str("symbol", symbol.Symbol).Msg("Saved symbol")
	}

	logger.Info().Msg("Symbol sync completed")
}
