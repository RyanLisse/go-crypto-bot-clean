package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Simple REST client that doesn't rely on specific packages
type MexcClient struct {
	apiKey    string
	apiSecret string
	baseURL   string
	client    *http.Client
	logger    zerolog.Logger
}

func NewMexcClient(apiKey, apiSecret string, logger zerolog.Logger) *MexcClient {
	return &MexcClient{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		baseURL:   "https://api.mexc.com",
		client:    &http.Client{Timeout: 10 * time.Second},
		logger:    logger,
	}
}

func main() {
	// Setup logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	logger := log.With().Str("component", "mexc-rest-script").Logger()

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		// Try parent directory
		err = godotenv.Load("../.env")
		if err != nil {
			logger.Fatal().Err(err).Msg("Error loading .env file")
		}
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
	client := NewMexcClient(apiKey, apiSecret, logger)
	logger.Info().Msg("MEXC REST client created")

	// Get exchange info (public endpoint that doesn't require authentication)
	logger.Debug().Msg("Fetching exchange information from MEXC")

	// Try to get exchange information
	resp, err := http.Get(client.baseURL + "/api/v3/exchangeInfo")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to get exchange information")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to read response body")
	}

	// Parse response as generic map
	var exchangeInfo map[string]interface{}
	if err := json.Unmarshal(body, &exchangeInfo); err != nil {
		logger.Fatal().Err(err).Msg("Failed to parse exchange information")
	}

	// Print exchange information summary
	if symbols, ok := exchangeInfo["symbols"].([]interface{}); ok {
		logger.Info().
			Int("Number of symbols", len(symbols)).
			Msg("Exchange information retrieved successfully")

		// Print first 5 symbols
		fmt.Println("\n=== MEXC Exchange Symbols (first 5) ===")
		count := 0
		for _, symbol := range symbols {
			if count >= 5 {
				break
			}
			if s, ok := symbol.(map[string]interface{}); ok {
				if symbolName, ok := s["symbol"].(string); ok {
					fmt.Println(symbolName)
					count++
				}
			}
		}
	} else {
		logger.Warn().Msg("Symbols field not found in exchange info response")
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
