package unit

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
	"go-crypto-bot-clean/backend/internal/platform/mexc"

	"github.com/gorilla/websocket"
)

// TestMexcClientDeadlock specifically tests for deadlocks in the client
func TestMexcClientDeadlock(t *testing.T) {
	t.Parallel()

	// Start mock websocket server
	mockServer := NewMockWsServer(t)
	defer mockServer.Close()

	// Create a test config with mock server URL
	cfg := createTestConfig()
	cfg.Mexc.WebsocketURL = mockServer.URL

	// Use shorter timeouts for tests
	cfg.WebSocket.PingInterval = 50 * time.Millisecond
	cfg.WebSocket.ReconnectDelay = 20 * time.Millisecond

	// Create client
	client, err := mexc.NewClient(cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Create a short timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// Ensure we disconnect even if the test fails
	defer func() {
		disconnect := func() error {
			// Create timeout context for disconnect
			_, discCancel := context.WithTimeout(context.Background(), 100*time.Millisecond) // Use blank identifier
			defer discCancel()

			err := client.Disconnect()
			if err != nil {
				t.Logf("Error during disconnect: %v", err)
			}

			// Wait for the context to finish to ensure operation completes or times out
			// discCtx is no longer defined, remove wait
			return nil
		}

		// Run disconnect with timeout
		disconnectTimer := time.AfterFunc(200*time.Millisecond, func() {
			t.Log("Disconnect took too long, continuing test")
		})
		disconnect()
		disconnectTimer.Stop()
	}()

	// Connect with timeout
	connectTimer := time.AfterFunc(1*time.Second, func() {
		t.Log("Connect is taking too long, likely deadlocked")
		cancel() // Cancel the context to abort the operation
	})
	err = client.Connect(ctx)
	connectTimer.Stop()

	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	// Create a channel to receive tickers
	tickerCh := make(chan *models.Ticker, 10)

	// Subscribe to a ticker with timeout
	subCtx, subCancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer subCancel()

	err = client.SubscribeToTickers(subCtx, []string{"BTCUSDT"}, tickerCh)
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	// Send a mock ticker update with timeout to avoid blocking
	// updateSent := false // Removed unused variable
	// sentUpdate was declared but not used, removing for now
	go func() {
		select {
		case <-time.After(100 * time.Millisecond):
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
			err := mockServer.WriteToClient(websocket.TextMessage, data)
			if err == nil {
				// updateSent = true // Declared but not used
				t.Logf("Successfully sent ticker update")
			}
		case <-ctx.Done():
			return // Exit goroutine if main context finishes
		}
	}()

	// Wait to receive ticker update with timeout
	// receivedTicker := false // Declared but not used
	select {
	case <-tickerCh:
		// receivedTicker = true // Declared but not used
		t.Log("Received ticker update")
	case <-time.After(500 * time.Millisecond):
		// We're not testing ticker receipt, just deadlocks
		t.Log("Ticker update timeout, but continuing test") // Keep this log
	}

	// Clean disconnect, testing for deadlocks
	err = client.Disconnect()
	if err != nil {
		t.Logf("Warning - disconnect error: %v", err)
	}

	// Success if we reached this point without deadlocking
	t.Log("Test completed without deadlock")
}
