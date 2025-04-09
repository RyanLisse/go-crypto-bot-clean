package unit

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"go-crypto-bot-clean/backend/internal/domain/models"
	"go-crypto-bot-clean/backend/internal/platform/mexc"
)

// TestMexcClientDeadlock specifically tests for deadlocks in the client
func TestMexcClientDeadlock(t *testing.T) {
	// Start mock websocket server
	mockServer := NewMockWsServer(t)
	defer mockServer.Close()

	// Create a test config with mock server URL
	cfg := createTestConfig()
	cfg.Mexc.WebsocketURL = mockServer.URL

	// Create client
	client, err := mexc.NewClient(cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Create a timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect
	err = client.Connect(ctx)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	// Ensure we disconnect even if the test fails
	defer func() {
		err := client.Disconnect()
		if err != nil {
			t.Logf("Error during disconnect: %v", err)
		}
	}()

	// Create a channel to receive tickers
	tickerCh := make(chan *models.Ticker, 10)

	// Subscribe to a ticker
	err = client.SubscribeToTickers(ctx, []string{"BTCUSDT"}, tickerCh)
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	// Send a mock ticker update
	go func() {
		time.Sleep(500 * time.Millisecond)
		tickerUpdate := map[string]interface{}{
			"channel": "spot@public.ticker.v3.api.BTCUSDT",
			"data": map[string]interface{}{
				"s": "BTCUSDT",
				"c": "45000.0",
				"h": "46000.0",
				"l": "44000.0",
				"v": "100.0",
				"q": "4500000.0",
				"p": "1000.0",
				"P": "2.2",
			},
			"ts": time.Now().UnixNano() / int64(time.Millisecond),
		}

		data, _ := json.Marshal(tickerUpdate)
		mockServer.WriteToClient(websocket.TextMessage, data)
	}()

	// Wait to receive ticker update
	select {
	case ticker := <-tickerCh:
		t.Logf("Received ticker: %+v", ticker)
	case <-time.After(3 * time.Second):
		t.Fatal("Timeout waiting for ticker update")
	}

	// Unsubscribe
	err = client.UnsubscribeFromTickers(ctx, []string{"BTCUSDT"})
	if err != nil {
		t.Fatalf("Failed to unsubscribe: %v", err)
	}

	// Test multiple connects/disconnects to check for goroutine leaks
	for i := 0; i < 3; i++ {
		err = client.Disconnect()
		if err != nil {
			t.Fatalf("Failed to disconnect: %v", err)
		}

		// Give time for goroutines to clean up
		time.Sleep(200 * time.Millisecond)

		err = client.Connect(ctx)
		if err != nil {
			t.Fatalf("Failed to reconnect: %v", err)
		}
	}

	// The test should complete within the timeout
	t.Log("Test completed without deadlock")
}
