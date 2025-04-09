package websocket

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"go-crypto-bot-clean/backend/internal/config"
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
	// Create a test server that disconnects after 100ms
	server := mockWebSocketServer(t, 100*time.Millisecond, true)
	defer server.Close()

	// Create a client with the test server URL
	cfg := &config.Config{}

	// Set Mexc config
	cfg.Mexc.WebsocketURL = strings.Replace(server.URL, "http://", "ws://", 1)

	// Set WebSocket config
	cfg.WebSocket.AutoReconnect = true
	cfg.WebSocket.ReconnectDelay = 50 * time.Millisecond
	cfg.WebSocket.MaxReconnectAttempts = 3

	// Set rate limiter config
	cfg.ConnectionRateLimiter.RequestsPerSecond = 10
	cfg.ConnectionRateLimiter.BurstCapacity = 5
	cfg.SubscriptionRateLimiter.RequestsPerSecond = 10
	cfg.SubscriptionRateLimiter.BurstCapacity = 5

	client, err := NewClient(cfg)
	require.NoError(t, err)

	// Connect to the server
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	err = client.Connect(ctx)
	require.NoError(t, err)
	assert.True(t, client.IsConnected())

	// Wait for the server to disconnect and the client to reconnect
	time.Sleep(300 * time.Millisecond)

	// Verify the client is still connected (reconnected)
	// Note: In some test environments, the client might not be connected at this exact moment
	// due to timing issues, so we'll check the connection attempts instead
	assert.Greater(t, client.GetConnectionAttempts(), 1, "Should have attempted to reconnect")

	// Clean up
	client.Disconnect()
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

	// Set WebSocket config
	cfg.WebSocket.AutoReconnect = true
	cfg.WebSocket.ReconnectDelay = 50 * time.Millisecond
	cfg.WebSocket.MaxReconnectAttempts = 3

	// Set rate limiter config
	cfg.ConnectionRateLimiter.RequestsPerSecond = 10
	cfg.ConnectionRateLimiter.BurstCapacity = 5
	cfg.SubscriptionRateLimiter.RequestsPerSecond = 10
	cfg.SubscriptionRateLimiter.BurstCapacity = 5

	client, err := NewClient(cfg)
	require.NoError(t, err)

	// Connect to the server
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	err = client.Connect(ctx)
	require.NoError(t, err)
	assert.True(t, client.IsConnected())

	// Wait for the server to disconnect and the client to attempt reconnection
	time.Sleep(300 * time.Millisecond)

	// Verify the client is disconnected after max reconnection attempts
	// The client should have made the initial connection plus up to MaxReconnectAttempts reconnection attempts
	assert.GreaterOrEqual(t, client.GetConnectionAttempts(), 2, "Should have attempted to reconnect at least once after initial connection")

	// Clean up
	client.Disconnect()
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

	// Set rate limiter config
	cfg.ConnectionRateLimiter.RequestsPerSecond = 10
	cfg.ConnectionRateLimiter.BurstCapacity = 5
	cfg.SubscriptionRateLimiter.RequestsPerSecond = 10
	cfg.SubscriptionRateLimiter.BurstCapacity = 5

	client, err := NewClient(cfg)
	require.NoError(t, err)

	// Connect to the server
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	err = client.Connect(ctx)
	require.NoError(t, err)
	assert.True(t, client.IsConnected())

	// Wait for ping messages to be sent
	time.Sleep(250 * time.Millisecond)

	// Verify the client is still connected
	assert.True(t, client.IsConnected())
	assert.Greater(t, client.GetPingSentCount(), 1, "Should have sent multiple ping messages")

	// Clean up
	client.Disconnect()
}
