package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// validateAPIKey checks if an API key contains any invalid characters for HTTP headers
func validateAPIKey(apiKey string) (bool, string) {
	// Check for control characters, spaces, or invalid header chars
	invalidChars := []string{"\r", "\n", "\t", " ", ",", ";", ":", "\""}
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
	// Setup logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	logger := log.With().Str("component", "fix-mexc-api-key").Logger()

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		logger.Fatal().Err(err).Msg("Error loading .env file")
	}

	// Get API credentials from environment variables
	apiKey := os.Getenv("MEXC_API_KEY")
	apiSecret := os.Getenv("MEXC_SECRET_KEY")

	if apiKey == "" {
		logger.Fatal().Msg("MEXC_API_KEY environment variable is not set")
	}

	if apiSecret == "" {
		logger.Fatal().Msg("MEXC_SECRET_KEY environment variable is not set")
	}

	// Log the original API key and secret (truncated for security)
	logger.Info().
		Str("Original API Key (truncated)", apiKey[:5]+"..."+apiKey[len(apiKey)-4:]).
		Str("Original API Secret (truncated)", apiSecret[:5]+"..."+apiSecret[len(apiSecret)-4:]).
		Msg("Original MEXC credentials")

	// Validate API key format
	valid, reason := validateAPIKey(apiKey)
	if !valid {
		logger.Error().Str("reason", reason).Msg("MEXC API key is invalid")
		
		// Try to fix the API key by trimming spaces and removing quotes
		fixedApiKey := strings.TrimSpace(apiKey)
		fixedApiKey = strings.Trim(fixedApiKey, "\"")
		
		logger.Info().
			Str("Original", apiKey).
			Str("Fixed", fixedApiKey).
			Msg("Fixed API key")
		
		// Check again after fixing
		valid, reason = validateAPIKey(fixedApiKey)
		if !valid {
			logger.Error().Str("reason", reason).Msg("MEXC API key is still invalid after fixing")
		} else {
			logger.Info().Msg("API key is now valid")
			
			// Update .env file
			envContent, err := os.ReadFile(".env")
			if err != nil {
				logger.Fatal().Err(err).Msg("Failed to read .env file")
			}
			
			lines := strings.Split(string(envContent), "\n")
			for i, line := range lines {
				if strings.HasPrefix(line, "MEXC_API_KEY=") {
					lines[i] = "MEXC_API_KEY=" + fixedApiKey
					logger.Info().Msg("Updated MEXC_API_KEY in .env file")
				}
			}
			
			err = os.WriteFile(".env", []byte(strings.Join(lines, "\n")), 0644)
			if err != nil {
				logger.Fatal().Err(err).Msg("Failed to write .env file")
			}
			
			logger.Info().Msg("Successfully updated .env file")
		}
	} else {
		logger.Info().Msg("API key is valid")
	}

	// Validate API secret format
	valid, reason = validateAPIKey(apiSecret)
	if !valid {
		logger.Error().Str("reason", reason).Msg("MEXC API secret is invalid")
		
		// Try to fix the API secret by trimming spaces and removing quotes
		fixedApiSecret := strings.TrimSpace(apiSecret)
		fixedApiSecret = strings.Trim(fixedApiSecret, "\"")
		
		logger.Info().
			Str("Original", apiSecret).
			Str("Fixed", fixedApiSecret).
			Msg("Fixed API secret")
		
		// Check again after fixing
		valid, reason = validateAPIKey(fixedApiSecret)
		if !valid {
			logger.Error().Str("reason", reason).Msg("MEXC API secret is still invalid after fixing")
		} else {
			logger.Info().Msg("API secret is now valid")
			
			// Update .env file
			envContent, err := os.ReadFile(".env")
			if err != nil {
				logger.Fatal().Err(err).Msg("Failed to read .env file")
			}
			
			lines := strings.Split(string(envContent), "\n")
			for i, line := range lines {
				if strings.HasPrefix(line, "MEXC_SECRET_KEY=") {
					lines[i] = "MEXC_SECRET_KEY=" + fixedApiSecret
					logger.Info().Msg("Updated MEXC_SECRET_KEY in .env file")
				}
			}
			
			err = os.WriteFile(".env", []byte(strings.Join(lines, "\n")), 0644)
			if err != nil {
				logger.Fatal().Err(err).Msg("Failed to write .env file")
			}
			
			logger.Info().Msg("Successfully updated .env file")
		}
	} else {
		logger.Info().Msg("API secret is valid")
	}
}
