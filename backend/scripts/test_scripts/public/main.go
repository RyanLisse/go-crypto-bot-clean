package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/pkg/platform/mexc"
)

func main() {
	// Setup logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	logger := log.With().Str("component", "mexc-public-script").Logger()

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		logger.Fatal().Err(err).Msg("Error loading .env file")
	}

	// Create MEXC client without API key
	client := mexc.NewClient("", "", &logger)
	logger.Info().Msg("MEXC client created (public only)")

	// Get exchange info (public API, no authentication required)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Debug().Msg("Fetching exchange information from MEXC")
	exchangeInfo, err := client.GetExchangeInfo(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to get exchange information")
	}

	// Print exchange information
	logger.Info().
		Int("Number of symbols", len(exchangeInfo.Symbols)).
		Msg("Exchange information retrieved successfully")

	// Print a few symbols
	fmt.Println("\n=== MEXC Exchange Symbols ===")
	fmt.Printf("%-10s %-10s %-10s %-10s\n", "Symbol", "Base", "Quote", "Status")
	fmt.Println("----------------------------------------")

	// Print first 10 symbols
	for i := 0; i < 10 && i < len(exchangeInfo.Symbols); i++ {
		symbol := exchangeInfo.Symbols[i]
		fmt.Printf("%-10s %-10s %-10s %-10s\n", symbol.Symbol, symbol.BaseAsset, symbol.QuoteAsset, symbol.Status)
	}

	// Save exchange info to file
	filename := "mexc_exchange_info.json"
	if err := saveToFile(exchangeInfo, filename); err != nil {
		logger.Error().Err(err).Str("filename", filename).Msg("Failed to save exchange info to file")
	} else {
		logger.Info().Str("filename", filename).Msg("Exchange info saved to file")
	}
}

// saveToFile saves data to a JSON file
func saveToFile(data any, filename string) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, jsonData, 0644)
}
