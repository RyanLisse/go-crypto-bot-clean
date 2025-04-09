package unit

import (
	"context"
	"log"
	"testing"

	"github.com/ryanlisse/go-crypto-bot/internal/config"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/service"
	"github.com/ryanlisse/go-crypto-bot/internal/platform/mexc"
	"github.com/stretchr/testify/assert"
)

// Simplified test config for the MEXC client
func createTestConfig() *config.Config {
	cfg := &config.Config{}

	// Set MEXC configuration
	cfg.Mexc.BaseURL = "https://api.mexc.com"
	cfg.Mexc.WebsocketURL = "wss://stream.mexc.com/ws"
	cfg.Mexc.APIKey = "test_api_key"
	cfg.Mexc.SecretKey = "test_secret_key"

	// Set rate limiter configurations
	cfg.ConnectionRateLimiter.RequestsPerSecond = 10
	cfg.ConnectionRateLimiter.BurstCapacity = 20

	cfg.SubscriptionRateLimiter.RequestsPerSecond = 5
	cfg.SubscriptionRateLimiter.BurstCapacity = 10

	return cfg
}

// TestNewClient tests the client creation
func TestNewClient(t *testing.T) {
	cfg := createTestConfig()
	logger := log.New(log.Writer(), "[TEST] ", log.LstdFlags)

	client, err := mexc.NewClient(cfg, mexc.WithLogger(logger))

	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Implements(t, (*service.ExchangeService)(nil), client)
}

// TestMexcClientUsage provides a concrete example of using the client
func TestMexcClientUsage(t *testing.T) {
	t.Skip("This is an example usage, not an actual test")

	// Create configuration
	cfg := &config.Config{}

	// Set MEXC configuration
	cfg.Mexc.BaseURL = "https://api.mexc.com"
	cfg.Mexc.WebsocketURL = "wss://stream.mexc.com/ws"
	cfg.Mexc.APIKey = "your_api_key"
	cfg.Mexc.SecretKey = "your_secret_key"

	// Set rate limiter configurations
	cfg.ConnectionRateLimiter.RequestsPerSecond = 20
	cfg.ConnectionRateLimiter.BurstCapacity = 40

	cfg.SubscriptionRateLimiter.RequestsPerSecond = 10
	cfg.SubscriptionRateLimiter.BurstCapacity = 20

	// Create a custom logger
	logger := log.New(log.Writer(), "[MEXC] ", log.LstdFlags)

	// Initialize the client
	client, err := mexc.NewClient(cfg, mexc.WithLogger(logger))
	if err != nil {
		t.Fatalf("Failed to create MEXC client: %v", err)
	}

	// Connect to WebSocket
	ctx := context.Background()
	err = client.Connect(ctx)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer client.Disconnect()

	// Create a channel to receive ticker updates
	tickerCh := make(chan *models.Ticker, 10)

	// Subscribe to tickers
	err = client.SubscribeToTickers(ctx, []string{"BTCUSDT", "ETHUSDT"}, tickerCh)
	if err != nil {
		t.Fatalf("Failed to subscribe to tickers: %v", err)
	}

	// Get current ticker data
	ticker, err := client.GetTicker(ctx, "BTCUSDT")
	if err != nil {
		t.Fatalf("Failed to get ticker: %v", err)
	}
	t.Logf("Current BTC price: %.2f", ticker.Price)

	// Get wallet balances
	wallet, err := client.GetWallet(ctx)
	if err != nil {
		t.Fatalf("Failed to get wallet: %v", err)
	}
	t.Logf("Wallet has %d assets", len(wallet.Balances))

	// Unsubscribe when done
	err = client.UnsubscribeFromTickers(ctx, []string{"BTCUSDT", "ETHUSDT"})
	if err != nil {
		t.Fatalf("Failed to unsubscribe from tickers: %v", err)
	}
}
