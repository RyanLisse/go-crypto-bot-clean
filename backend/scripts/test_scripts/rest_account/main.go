package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/pkg/platform/mexc/rest"
)

func main() {
	// Setup logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	logger := log.With().Str("component", "mexc-rest-script").Logger()

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		logger.Fatal().Err(err).Msg("Error loading .env file")
	}

	// Get API credentials from environment variables
	apiKey := os.Getenv("MEXC_API_KEY")
	apiSecret := os.Getenv("MEXC_SECRET_KEY")

	if apiKey == "" || apiSecret == "" {
		logger.Fatal().Msg("MEXC_API_KEY and MEXC_SECRET_KEY environment variables must be set")
	}

	// Log the API key and secret (truncated for security)
	logger.Info().
		Str("API Key (truncated)", apiKey[:5]+"..."+apiKey[len(apiKey)-4:]).
		Str("API Secret (truncated)", apiSecret[:5]+"..."+apiSecret[len(apiSecret)-4:]).
		Msg("Using MEXC credentials")

	// Create MEXC REST client
	client := rest.NewClient(apiKey, apiSecret)
	logger.Info().Msg("MEXC REST client created")

	// Try to get account information
	account, err := client.GetAccount()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to get account information")
	}

	// Print account information
	logger.Info().
		Int("Number of balances", len(account.Wallet.Balances)).
		Strs("Permissions", account.Permissions).
		Msg("Account information retrieved successfully")

	// Print balances
	fmt.Println("\n=== MEXC Account Balances ===")
	fmt.Printf("%-10s %-15s %-15s %-15s\n", "Asset", "Free", "Locked", "Total")
	fmt.Println("--------------------------------------------------")

	// Filter out zero balances
	nonZeroBalances := 0
	for asset, balance := range account.Wallet.Balances {
		if balance.Free > 0 || balance.Locked > 0 {
			fmt.Printf("%-10s %-15f %-15f %-15f\n", asset, balance.Free, balance.Locked, balance.Total)
			nonZeroBalances++
		}
	}

	if nonZeroBalances == 0 {
		fmt.Println("No non-zero balances found.")
	}

	// Save account info to file
	filename := "mexc_account_info.json"
	if err := saveToFile(account, filename); err != nil {
		logger.Error().Err(err).Str("filename", filename).Msg("Failed to save account info to file")
	} else {
		logger.Info().Str("filename", filename).Msg("Account info saved to file")
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
