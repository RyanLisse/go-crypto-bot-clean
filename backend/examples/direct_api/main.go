package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/factory"
	"github.com/rs/zerolog"
)

func main() {
	// Initialize logger
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	logger.Info().Msg("Starting direct API test")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Create MEXC client
	var mexcClient port.MEXCClient
	marketFactory := factory.NewMarketFactory(cfg, &logger, nil)
	mexcClient = marketFactory.CreateMEXCClient()

	// Test getting ticker for BTCUSDT
	ctx := context.Background()
	ticker, err := mexcClient.GetMarketData(ctx, "BTCUSDT")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to get ticker")
	}

	// Print ticker as JSON
	tickerJSON, err := json.MarshalIndent(ticker, "", "  ")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to marshal ticker to JSON")
	}
	fmt.Println("Ticker JSON:")
	fmt.Println(string(tickerJSON))

	// Test getting order book for BTCUSDT
	orderBook, err := mexcClient.GetOrderBook(ctx, "BTCUSDT", 5)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to get order book")
	}

	// Print order book as JSON
	orderBookJSON, err := json.MarshalIndent(orderBook, "", "  ")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to marshal order book to JSON")
	}
	fmt.Println("\nOrder Book JSON:")
	fmt.Println(string(orderBookJSON))

	// Test getting exchange info
	exchangeInfo, err := mexcClient.GetExchangeInfo(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to get exchange info")
	}

	// Print first 5 symbols from exchange info
	fmt.Println("\nFirst 5 symbols from exchange info:")
	for i, symbol := range exchangeInfo.Symbols {
		if i >= 5 {
			break
		}
		fmt.Printf("%d. %s (%s/%s)\n", i+1, symbol.Symbol, symbol.BaseAsset, symbol.QuoteAsset)
	}

	logger.Info().Msg("Direct API test completed successfully")
}
