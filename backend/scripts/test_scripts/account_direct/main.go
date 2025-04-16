package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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

func main() {
	Run()
}

// Run is the main entry point for this script
func Run() {
	// Setup logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	logger := log.With().Str("component", "mexc-account-direct").Logger()

	// Load environment variables from .env file
	// Try current directory first, then parent directory
	err := godotenv.Load()
	if err != nil {
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

	// Create HTTP client
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create timestamp for the request
	timestamp := time.Now().UnixMilli()

	// Create query parameters
	params := fmt.Sprintf("timestamp=%d", timestamp)

	// Generate signature
	h := hmac.New(sha256.New, []byte(apiSecret))
	h.Write([]byte(params))
	signature := hex.EncodeToString(h.Sum(nil))

	// Add signature to parameters
	endpoint := fmt.Sprintf("/api/v3/account?%s&signature=%s", params, signature)

	// Create request
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.mexc.com"+endpoint, nil)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create request")
	}

	// Add API key header
	logger.Debug().Str("X-MBX-APIKEY", apiKey).Msg("Setting API key header")
	req.Header.Set("X-MBX-APIKEY", apiKey)

	// Send request
	logger.Debug().Msg("Sending request to MEXC API")
	resp, err := httpClient.Do(req)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to send request")
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to read response body")
	}

	// Check response status
	logger.Debug().Int("status", resp.StatusCode).Msg("Received response from MEXC API")
	if resp.StatusCode != http.StatusOK {
		var errResp struct {
			Code    int    `json:"code"`
			Message string `json:"msg"`
		}
		if err := json.Unmarshal(body, &errResp); err != nil {
			logger.Fatal().Err(err).Int("status", resp.StatusCode).Str("body", string(body)).Msg("Failed to decode error response")
		}
		logger.Fatal().Int("code", errResp.Code).Str("message", errResp.Message).Str("body", string(body)).Msg("MEXC API error")
	}

	// Parse response
	var accountInfo struct {
		MakerCommission  int    `json:"makerCommission"`
		TakerCommission  int    `json:"takerCommission"`
		BuyerCommission  int    `json:"buyerCommission"`
		SellerCommission int    `json:"sellerCommission"`
		CanTrade         bool   `json:"canTrade"`
		CanWithdraw      bool   `json:"canWithdraw"`
		CanDeposit       bool   `json:"canDeposit"`
		UpdateTime       int64  `json:"updateTime"`
		AccountType      string `json:"accountType"`
		Balances         []struct {
			Asset  string `json:"asset"`
			Free   string `json:"free"`
			Locked string `json:"locked"`
		} `json:"balances"`
		Permissions []string `json:"permissions"`
	}

	if err := json.Unmarshal(body, &accountInfo); err != nil {
		logger.Fatal().Err(err).Msg("Failed to decode account response")
	}

	// Print account information
	logger.Info().
		Int("Number of balances", len(accountInfo.Balances)).
		Bool("Can Trade", accountInfo.CanTrade).
		Bool("Can Withdraw", accountInfo.CanWithdraw).
		Bool("Can Deposit", accountInfo.CanDeposit).
		Strs("Permissions", accountInfo.Permissions).
		Msg("Account information retrieved successfully")

	// Print balances
	fmt.Println("\n=== MEXC Account Balances ===")
	fmt.Printf("%-10s %-15s %-15s\n", "Asset", "Free", "Locked")
	fmt.Println("----------------------------------------")

	// Filter out zero balances
	nonZeroBalances := 0
	for _, balance := range accountInfo.Balances {
		if balance.Free != "0" || balance.Locked != "0" {
			fmt.Printf("%-10s %-15s %-15s\n", balance.Asset, balance.Free, balance.Locked)
			nonZeroBalances++
		}
	}

	if nonZeroBalances == 0 {
		fmt.Println("No non-zero balances found.")
	}

	// Save account info to file
	filename := "mexc_account_info.json"
	if err := SaveToFile(accountInfo, filename); err != nil {
		logger.Error().Err(err).Str("filename", filename).Msg("Failed to save account info to file")
	} else {
		logger.Info().Str("filename", filename).Msg("Account info saved to file")
	}
}

// SaveToFile saves data to a JSON file
func SaveToFile(data any, filename string) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, jsonData, 0644)
}
