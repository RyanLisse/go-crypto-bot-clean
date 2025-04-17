package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/pkg/platform/mexc"
	"github.com/rs/zerolog"
)

func main() {
	// Initialize logger
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	logger.Info().Msg("Starting market data test")

	// Check for API key in environment
	apiKey := os.Getenv("MEXC_API_KEY")
	secretKey := os.Getenv("MEXC_SECRET_KEY")
	baseURL := os.Getenv("MEXC_BASE_URL")

	if apiKey == "" || secretKey == "" {
		logger.Warn().Msg("MEXC_API_KEY or MEXC_SECRET_KEY not set. Using default configuration.")
		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			log.Fatalf("Failed to load configuration: %v", err)
		}
		apiKey = cfg.MEXC.APIKey
		secretKey = cfg.MEXC.APISecret
		baseURL = cfg.MEXC.BaseURL
	}

	if baseURL == "" {
		baseURL = "https://api.mexc.com"
	}

	logger.Info().Str("baseURL", baseURL).Msg("Using MEXC API")

	// Create MEXC client directly
	mexcClient := mexc.NewClient(apiKey, secretKey, &logger)

	// Test getting exchange info
	ctx := context.Background()
	exchangeInfo, err := mexcClient.GetExchangeInfo(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to get exchange info")
	}

	logger.Info().Int("symbolCount", len(exchangeInfo.Symbols)).Msg("Got exchange info")

	// Print first 5 symbols
	for i, symbol := range exchangeInfo.Symbols {
		if i >= 5 {
			break
		}
		logger.Info().
			Str("symbol", symbol.Symbol).
			Str("baseAsset", symbol.BaseAsset).
			Str("quoteAsset", symbol.QuoteAsset).
			Str("status", symbol.Status).
			Msg("Symbol info")
	}

	// Test getting ticker for BTC/USDT
	ticker, err := mexcClient.GetMarketData(ctx, "BTCUSDT")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to get ticker")
	}

	logger.Info().
		Str("symbol", ticker.Symbol).
		Float64("price", ticker.LastPrice).
		Float64("volume", ticker.Volume).
		Msg("Got ticker")

	// Test getting order book for BTC/USDT
	orderBook, err := mexcClient.GetOrderBook(ctx, "BTCUSDT", 5)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to get order book")
	}

	logger.Info().
		Str("symbol", orderBook.Symbol).
		Int("bidCount", len(orderBook.Bids)).
		Int("askCount", len(orderBook.Asks)).
		Msg("Got order book")

	// Print first bid and ask
	if len(orderBook.Bids) > 0 {
		logger.Info().
			Float64("price", orderBook.Bids[0].Price).
			Float64("quantity", orderBook.Bids[0].Quantity).
			Msg("Top bid")
	}

	if len(orderBook.Asks) > 0 {
		logger.Info().
			Float64("price", orderBook.Asks[0].Price).
			Float64("quantity", orderBook.Asks[0].Quantity).
			Msg("Top ask")
	}

	// Note: GetKlines requires authentication with a valid API key
	// Skipping candle test for now
	logger.Info().Msg("Skipping candle test as it requires authentication")

	// Pretty print the ticker as JSON
	tickerJSON, _ := json.MarshalIndent(ticker, "", "  ")
	fmt.Println("Ticker JSON:")
	fmt.Println(string(tickerJSON))

	logger.Info().Msg("Market data test completed successfully")
}
