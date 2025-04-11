package rest

import (
	"context"
	"os"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestMEXCAPIKeys(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	// Get API keys from environment variables
	apiKey := os.Getenv("MEXC_API_KEY")
	secretKey := os.Getenv("MEXC_SECRET_KEY")

	if apiKey == "" || secretKey == "" {
		t.Skip("MEXC_API_KEY and MEXC_SECRET_KEY environment variables must be set")
	}

	// Create MEXC client
	client, err := NewClient(apiKey, secretKey, WithLogger(logger))
	if err != nil {
		t.Fatalf("Failed to create MEXC client: %v", err)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test API key validation
	t.Run("Validate API Keys", func(t *testing.T) {
		valid, err := client.ValidateKeys(ctx)
		if err != nil {
			t.Fatalf("Error validating keys: %v", err)
		}
		if !valid {
			t.Error("API keys are invalid")
		}
	})

	// Test wallet information retrieval
	t.Run("Get Wallet Information", func(t *testing.T) {
		wallet, err := client.GetWallet(ctx)
		if err != nil {
			t.Fatalf("Error getting wallet: %v", err)
		}
		if wallet == nil {
			t.Error("Wallet information is nil")
		}
		if len(wallet.Balances) == 0 {
			t.Log("Warning: Wallet has no balances")
		}
	})
}
