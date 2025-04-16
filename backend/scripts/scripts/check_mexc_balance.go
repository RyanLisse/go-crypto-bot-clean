package scripts

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

// CheckMexcBalance checks the MEXC account balance
func CheckMexcBalance() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	logger := log.With().Str("component", "mexc-balance-script").Logger()

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

	// Create MEXC client
	client := mexc.NewClient(apiKey, apiSecret, &logger)
	logger.Info().Msg("MEXC client created")

	// Debug: Print API key header
	logger.Debug().Str("APIKEY", apiKey).Msg("API key header value")

	// Get account information
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Debug().Msg("Fetching account information from MEXC")
	wallet, err := client.GetAccount(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to get account information")
	}

	// Print account information
	logger.Info().
		Int("Number of balances", len(wallet.Balances)).
		Msg("Account information retrieved successfully")

	// Print balances
	fmt.Println("\n=== MEXC Account Balances ===")
	fmt.Printf("%-10s %-15s %-15s\n", "Asset", "Free", "Locked")
	fmt.Println("----------------------------------------")

	// Filter out zero balances
	nonZeroBalances := 0
	for _, balance := range wallet.Balances {
		if balance.Free > 0 || balance.Locked > 0 {
			fmt.Printf("%-10s %-15f %-15f\n", balance.Asset, balance.Free, balance.Locked)
			nonZeroBalances++
		}
	}

	if nonZeroBalances == 0 {
		fmt.Println("No non-zero balances found.")
	}

	// Save wallet data to file
	filename := "mexc_wallet.json"
	if err := saveWalletToFile(wallet, filename); err != nil {
		logger.Error().Err(err).Str("filename", filename).Msg("Failed to save wallet data to file")
	} else {
		logger.Info().Str("filename", filename).Msg("Wallet data saved to file")
	}
}

// saveWalletToFile saves the wallet data to a JSON file
func saveWalletToFile(wallet any, filename string) error {
	data, err := json.MarshalIndent(wallet, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}
