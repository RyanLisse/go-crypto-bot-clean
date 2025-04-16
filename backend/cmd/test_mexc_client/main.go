package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/logger"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/pkg/platform/mexc"
)

// validateAPIKey checks if an API key contains any invalid characters for HTTP headers
func validateAPIKey(apiKey string) (bool, string) {
	// Check for control characters, spaces, or invalid header chars
	invalidChars := []string{"\r", "\n", "\t", " ", ",", ";", ":"}
	for _, char := range invalidChars {
		if strings.Contains(apiKey, char) {
			return false, fmt.Sprintf("API key contains invalid character: %q", char)
		}
	}

	// Check if the key is too long (unlikely but possible)
	if len(apiKey) > 500 {
		return false, "API key is too long (> 500 chars)"
	}

	return true, ""
}

func main() {
	// Initialize logger
	log := logger.NewLogger()
	log.Info().Msg("Testing MEXC client initialization")

	// Manually check environment variables
	apiKey := os.Getenv("MEXC_API_KEY")
	apiSecret := os.Getenv("MEXC_SECRET_KEY")
	encryptionKey := os.Getenv("MEXC_CRED_ENCRYPTION_KEY")

	if apiKey == "" {
		log.Error().Msg("MEXC_API_KEY environment variable is not set")
		os.Exit(1)
	}

	if apiSecret == "" {
		log.Error().Msg("MEXC_SECRET_KEY environment variable is not set")
		os.Exit(1)
	}

	if encryptionKey == "" {
		log.Error().Msg("MEXC_CRED_ENCRYPTION_KEY environment variable is not set")
		os.Exit(1)
	}

	// Validate API key format
	valid, reason := validateAPIKey(apiKey)
	if !valid {
		log.Error().Str("reason", reason).Msg("MEXC API key is invalid")
		// Try to fix the API key by trimming spaces
		apiKey = strings.TrimSpace(apiKey)
		log.Info().Msg("Trimmed spaces from API key")

		// Check again after trimming
		valid, reason = validateAPIKey(apiKey)
		if !valid {
			log.Error().Str("reason", reason).Msg("MEXC API key is still invalid after trimming")
			os.Exit(1)
		}
	}

	log.Info().
		Str("API Key (truncated)", apiKey[:5]+"..."+apiKey[len(apiKey)-4:]).
		Str("API Secret (truncated)", apiSecret[:5]+"..."+apiSecret[len(apiSecret)-4:]).
		Msg("Found MEXC credentials in environment variables")

	// Load config using the standard method
	log.Info().Msg("Loading configuration")
	cfg := config.LoadConfig(log)

	// Check if MEXC config is present
	if cfg.MEXC.APIKey == "" {
		log.Error().Msg("MEXC API Key is empty in the loaded configuration")
		os.Exit(1)
	}

	if cfg.MEXC.APISecret == "" {
		log.Error().Msg("MEXC API Secret is empty in the loaded configuration")
		os.Exit(1)
	}

	// Try direct HTTP call to test if the API key is valid
	log.Info().Msg("Testing API key directly with HTTP request")
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", "https://api.mexc.com/api/v3/exchangeInfo", nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create test request")
		os.Exit(1)
	}

	// MEXC API requires the APIKEY header, not X-MBX-APIKEY
	req.Header.Set("APIKEY", apiKey)
	log.Info().Str("header", "APIKEY").Str("value", apiKey).Msg("Set API key header")

	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send test request")
		os.Exit(1)
	}
	defer resp.Body.Close()

	log.Info().Int("status", resp.StatusCode).Msg("Test request complete")
	if resp.StatusCode != http.StatusOK {
		log.Error().Int("status_code", resp.StatusCode).Msg("Non-200 status code from direct API test")
	} else {
		log.Info().Msg("Direct API test successful")
	}

	log.Info().
		Str("API Key from config (truncated)", cfg.MEXC.APIKey[:5]+"..."+cfg.MEXC.APIKey[len(cfg.MEXC.APIKey)-4:]).
		Str("API Secret from config (truncated)", cfg.MEXC.APISecret[:5]+"..."+cfg.MEXC.APISecret[len(cfg.MEXC.APISecret)-4:]).
		Bool("UseTestnet", cfg.MEXC.UseTestnet).
		Int("Rate Limit (requests per minute)", cfg.MEXC.RateLimit.RequestsPerMinute).
		Msg("MEXC configuration loaded successfully")

	// Initialize MEXC client
	log.Info().Msg("Initializing MEXC client")
	mexcClient := mexc.NewClient(apiKey, apiSecret, log) // Use directly from environment

	// Test a simple API call
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Info().Msg("Testing MEXC client with GetExchangeInfo API call")
	exchangeInfo, err := mexcClient.GetExchangeInfo(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get exchange info")
		os.Exit(1)
	}

	log.Info().
		Int("Number of symbols", len(exchangeInfo.Symbols)).
		Msg("Successfully retrieved exchange info")

	// Print a few symbols
	log.Info().Msg("Sample symbols from exchange info:")
	for i := 0; i < 5 && i < len(exchangeInfo.Symbols); i++ {
		symbol := exchangeInfo.Symbols[i]
		log.Info().
			Str("Symbol", symbol.Symbol).
			Str("Base/Quote", fmt.Sprintf("%s/%s", symbol.BaseAsset, symbol.QuoteAsset)).
			Str("Status", symbol.Status).
			Msg("Symbol info")
	}

	log.Info().Msg("MEXC client test completed successfully")
}
