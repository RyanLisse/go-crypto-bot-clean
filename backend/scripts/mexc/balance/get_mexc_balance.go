package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/pkg/platform/mexc"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Setup logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	logger := log.With().Str("component", "mexc-balance-script").Logger()

	// Get API credentials from environment variables
	apiKey := os.Getenv("MEXC_API_KEY")
	apiSecret := os.Getenv("MEXC_SECRET_KEY")

	if apiKey == "" || apiSecret == "" {
		logger.Fatal().Msg("MEXC_API_KEY and MEXC_SECRET_KEY environment variables must be set")
	}

	// Create MEXC client
	client := mexc.NewClient(apiKey, apiSecret, &logger)
	logger.Info().Msg("MEXC client created")

	// Note: We're using the standard MEXC client, not the REST client directly

	// Get account information
	ctx := context.Background()
	wallet, err := client.GetAccount(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to get account information")
	}

	// Print wallet information
	logger.Info().
		Str("userID", wallet.UserID).
		Str("exchange", wallet.Exchange).
		Float64("totalUSDValue", wallet.TotalUSDValue).
		Time("lastUpdated", wallet.LastUpdated).
		Msg("Account information retrieved")

	// Print balances
	fmt.Println("\n=== ACCOUNT BALANCES ===")
	fmt.Printf("%-10s %-15s %-15s %-15s %-15s\n", "ASSET", "FREE", "LOCKED", "TOTAL", "USD VALUE")
	fmt.Println("----------------------------------------------------------------------")

	// Sort balances by USD value (descending)
	for asset, balance := range wallet.Balances {
		fmt.Printf("%-10s %-15.8f %-15.8f %-15.8f %-15.2f\n",
			asset,
			balance.Free,
			balance.Locked,
			balance.Total,
			balance.USDValue)
	}

	// Save the wallet data to a file for future reference
	saveWalletToFile(wallet, "mexc_balance.json")
	logger.Info().Msg("Wallet data saved to mexc_balance.json")
}

// saveWalletToFile saves the wallet data to a JSON file
func saveWalletToFile(wallet interface{}, filename string) error {
	data, err := json.MarshalIndent(wallet, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}
