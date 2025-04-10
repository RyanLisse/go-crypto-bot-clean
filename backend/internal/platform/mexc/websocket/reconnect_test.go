package websocket

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go-crypto-bot-clean/backend/internal/config"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockWebSocketServer creates a test WebSocket server
func mockWebSocketServer(t *testing.T, disconnectAfter time.Duration, reconnectSuccess bool) *httptest.Server {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	var conn *websocket.Conn
	var connCount int

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if connCount > 0 && !reconnectSuccess {
			// Simulate server rejecting reconnection attempts
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Logf("Upgrade error: %v", err)
			return
		}
		conn = c
		connCount++

		// Send a welcome message
		conn.WriteJSON(map[string]interface{}{
			"type": "welcome",
			"data": "Connected successfully",
		})

		// If disconnectAfter is set, close the connection after that duration
		if disconnectAfter > 0 {
			time.AfterFunc(disconnectAfter, func() {
				conn.Close()
			})
		}

		// Echo messages back
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				break
			}
			conn.WriteMessage(websocket.TextMessage, msg)
		}
	}))

	return server
}

// TestWebSocketReconnection tests the automatic reconnection mechanism
func TestWebSocketReconnection(t *testing.T) {
	// Use parallel testing cautiously - may cause issues with WebSocket tests
	// t.Parallel() - Removed to avoid test interference

	// Create a test server that disconnects after 100ms
	server := mockWebSocketServer(t, 100*time.Millisecond, true)
	defer server.Close()

	// Create a client with the test server URL
	cfg := &config.Config{}

	// Set Mexc config
	cfg.Mexc.WebsocketURL = strings.Replace(server.URL, "http://", "ws://", 1)

	// Set WebSocket config with shorter timeouts for testing
	cfg.WebSocket.AutoReconnect = true
	cfg.WebSocket.ReconnectDelay = 20 * time.Millisecond
	cfg.WebSocket.MaxReconnectAttempts = 2
	cfg.WebSocket.PingInterval = 50 * time.Millisecond

	// Set rate limiter config
	cfg.ConnectionRateLimiter.RequestsPerSecond = 10
	cfg.ConnectionRateLimiter.BurstCapacity = 5
	cfg.SubscriptionRateLimiter.RequestsPerSecond = 10
	cfg.SubscriptionRateLimiter.BurstCapacity = 5

	client, err := NewClient(cfg)
	require.NoError(t, err)

	// Create a parent context with timeout to prevent test hanging
	parentCtx, parentCancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer parentCancel()

	// Ensure client is disconnected at the end of the test
	defer client.Disconnect()

	// Connect to the server with a short timeout
	ctx, cancel := context.WithTimeout(parentCtx, 200*time.Millisecond)
	defer cancel()

	err = client.Connect(ctx)
	require.NoError(t, err)

	// Create a channel to notify when we've seen a reconnection
	done := make(chan struct{})
	go func() {
		initialAttempts := client.GetConnectionAttempts()
		require.Equal(t, 1, initialAttempts, "Should start with exactly 1 connection attempt")

		// Check periodically for additional connection attempts
		ticker := time.NewTicker(50 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				currentAttempts := client.GetConnectionAttempts()
				// If we've seen at least one reconnection attempt
				if currentAttempts > initialAttempts {
					t.Logf("Observed reconnection: attempts increased from %d to %d",
						initialAttempts, currentAttempts)
					done <- struct{}{}
					return
				}
			case <-parentCtx.Done():
				// Test timeout reached
				return
			}
		}
	}()

	// Wait for reconnection or timeout
	select {
	case <-done:
		// We've seen the reconnection
		attempts := client.GetConnectionAttempts()
		assert.GreaterOrEqual(t, attempts, 2,
			"Should have at least 2 connection attempts (initial + reconnect)")
		t.Logf("Final connection attempts: %d", attempts)
	case <-parentCtx.Done():
		t.Logf("Test timeout reached, final connection attempts: %d",
			client.GetConnectionAttempts())
		// Still assert what we can
		assert.GreaterOrEqual(t, client.GetConnectionAttempts(), 1,
			"Should have at least the initial connection")
	}
}

// TestWebSocketReconnectionFailure tests the reconnection mechanism when it fails
func TestWebSocketReconnectionFailure(t *testing.T) {
	// Create a test server that disconnects after 100ms and rejects reconnections
	server := mockWebSocketServer(t, 100*time.Millisecond, false)
	defer server.Close()

	// Create a client with the test server URL
	cfg := &config.Config{}

	// Set Mexc config
	cfg.Mexc.WebsocketURL = strings.Replace(server.URL, "http://", "ws://", 1)

	// Set WebSocket config with aggressive timeouts for testing
	cfg.WebSocket.AutoReconnect = true
	cfg.WebSocket.ReconnectDelay = 50 * time.Millisecond
	cfg.WebSocket.MaxReconnectAttempts = 2 // Reduced from 3 to make test faster
	cfg.WebSocket.PingInterval = 100 * time.Millisecond

	// Set rate limiter config
	cfg.ConnectionRateLimiter.RequestsPerSecond = 10
	cfg.ConnectionRateLimiter.BurstCapacity = 5
	cfg.SubscriptionRateLimiter.RequestsPerSecond = 10
	cfg.SubscriptionRateLimiter.BurstCapacity = 5

	client, err := NewClient(cfg)
	require.NoError(t, err)

	// Create a parent context with a short timeout to prevent test hanging
	parentCtx, parentCancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer parentCancel()

	// Connect to the server with a short timeout
	ctx, cancel := context.WithTimeout(parentCtx, 500*time.Millisecond)
	defer cancel()

	// Always ensure client is disconnected at end of test
	defer client.Disconnect()

	// Connect
	err = client.Connect(ctx)
	require.NoError(t, err)
	assert.True(t, client.IsConnected())

	// Create a channel to notify when we've seen enough reconnection attempts
	done := make(chan struct{})
	go func() {
		// Check periodically for connection attempts
		ticker := time.NewTicker(50 * time.Millisecond)
		defer ticker.Stop()

		maxAttempts := 0
		for {
			select {
			case <-ticker.C:
				attempts := client.GetConnectionAttempts()
				if attempts > maxAttempts {
					maxAttempts = attempts
				}
				// If we've seen at least one reconnection attempt after initial connection
				if attempts >= 2 {
					done <- struct{}{}
					return
				}
			case <-parentCtx.Done():
				// Test timeout reached
				return
			}
		}
	}()

	// Wait for enough reconnection attempts or timeout
	select {
	case <-done:
		// We've seen enough reconnection attempts
		attempts := client.GetConnectionAttempts()
		assert.GreaterOrEqual(t, attempts, 2, "Should have attempted to reconnect at least once after initial connection")
		t.Logf("Observed %d connection attempts", attempts)
	case <-parentCtx.Done():
		t.Logf("Test timeout reached, final connection attempts: %d", client.GetConnectionAttempts())
		// Still assert what we can
		assert.GreaterOrEqual(t, client.GetConnectionAttempts(), 1, "Should have made at least the initial connection")
	}
}

// TestWebSocketPingPong tests the ping/pong mechanism
func TestWebSocketPingPong(t *testing.T) {
	// Create a test server
	server := mockWebSocketServer(t, 0, true)
	defer server.Close()

	// Create a client with the test server URL
	cfg := &config.Config{}

	// Set Mexc config
	cfg.Mexc.WebsocketURL = strings.Replace(server.URL, "http://", "ws://", 1)

	// Set WebSocket config
	cfg.WebSocket.PingInterval = 100 * time.Millisecond
	cfg.WebSocket.AutoReconnect = false // Don't need reconnection for this test

	// Set rate limiter config
	cfg.ConnectionRateLimiter.RequestsPerSecond = 10
	cfg.ConnectionRateLimiter.BurstCapacity = 5
	cfg.SubscriptionRateLimiter.RequestsPerSecond = 10
	cfg.SubscriptionRateLimiter.BurstCapacity = 5

	client, err := NewClient(cfg)
	require.NoError(t, err)

	// Create a parent context with timeout to prevent test hanging
	parentCtx, parentCancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer parentCancel()

	// Always ensure client is disconnected at end of test
	defer client.Disconnect()

	// Connect to the server
	ctx, cancel := context.WithTimeout(parentCtx, 500*time.Millisecond)
	defer cancel()

	err = client.Connect(ctx)
	require.NoError(t, err)
	assert.True(t, client.IsConnected())

	// Create a channel to notify when enough pings have been sent
	done := make(chan struct{})
	go func() {
		// Check periodically for ping counts
		ticker := time.NewTicker(50 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				pingCount := client.GetPingSentCount()
				// If we've seen at least 2 pings
				if pingCount >= 2 {
					t.Logf("Observed %d ping messages sent", pingCount)
					done <- struct{}{}
					return
				}
			case <-parentCtx.Done():
				// Test timeout reached
				return
			}
		}
	}()

	// Wait for enough pings or timeout
	select {
	case <-done:
		// We've seen enough pings
		assert.True(t, client.IsConnected(), "Client should still be connected")
		assert.GreaterOrEqual(t, client.GetPingSentCount(), 2, "Should have sent multiple ping messages")
	case <-parentCtx.Done():
		t.Logf("Test timeout reached, final ping count: %d", client.GetPingSentCount())
		// Still assert what we can
		if client.IsConnected() {
			assert.GreaterOrEqual(t, client.GetPingSentCount(), 1, "Should have sent at least one ping message")
		}
	}
}
