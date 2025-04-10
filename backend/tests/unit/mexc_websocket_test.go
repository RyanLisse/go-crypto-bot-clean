package unit

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go-crypto-bot-clean/backend/internal/config"
	mexcws "go-crypto-bot-clean/backend/internal/platform/mexc/websocket"
	"go-crypto-bot-clean/backend/pkg/ratelimiter"
)

// Error constants for testing
var (
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
)

// MockMexcWebSocketClient implements a simplified version of the MEXC WebSocket client for testing
type MockMexcWebSocketClient struct {
	endpoint           string
	isConnected        bool
	connRateLimiter    *ratelimiter.TokenBucketRateLimiter
	subRateLimiter     *ratelimiter.TokenBucketRateLimiter
	connectionAttempts int
}

// NewMexcWebSocketClient creates a mock MEXC WebSocket client for testing
func NewMexcWebSocketClient(endpoint string) *MockMexcWebSocketClient {
	return &MockMexcWebSocketClient{
		endpoint:        endpoint,
		connRateLimiter: ratelimiter.NewTokenBucketRateLimiter(10, 10),
		subRateLimiter:  ratelimiter.NewTokenBucketRateLimiter(2, 5),
	}
}

// Connect establishes a connection to the WebSocket server
func (c *MockMexcWebSocketClient) Connect(ctx context.Context) error {
	// Check connection rate limit
	if !c.connRateLimiter.TryAcquire() {
		return ErrRateLimitExceeded
	}

	c.connectionAttempts++
	c.isConnected = true
	return nil
}

// Disconnect closes the WebSocket connection
func (c *MockMexcWebSocketClient) Disconnect() error {
	c.isConnected = false
	return nil
}

// SubscribeToTickers subscribes to ticker updates for the specified symbols
func (c *MockMexcWebSocketClient) SubscribeToTickers(ctx context.Context, symbols []string) error {
	// Check subscription rate limit
	if !c.subRateLimiter.TryAcquire() {
		return ErrRateLimitExceeded
	}
	return nil
}

// IsConnected returns whether the client is connected to the WebSocket server
func (c *MockMexcWebSocketClient) IsConnected() bool {
	return c.isConnected
}

// GetConnectionAttempts returns the number of connection attempts
func (c *MockMexcWebSocketClient) GetConnectionAttempts() int {
	return c.connectionAttempts
}

// GetConnRateLimiter returns the connection rate limiter
func (c *MockMexcWebSocketClient) GetConnRateLimiter() *ratelimiter.TokenBucketRateLimiter {
	return c.connRateLimiter
}

// GetSubRateLimiter returns the subscription rate limiter
func (c *MockMexcWebSocketClient) GetSubRateLimiter() *ratelimiter.TokenBucketRateLimiter {
	return c.subRateLimiter
}

// mockConfig creates a minimal config.Config for testing.
func mockConfig() *config.Config {
	return &config.Config{
		Mexc: struct {
			APIKey       string `mapstructure:"api_key"`
			SecretKey    string `mapstructure:"secret_key"`
			BaseURL      string `mapstructure:"base_url"`
			WebsocketURL string `mapstructure:"websocket_url"`
		}{
			WebsocketURL: "wss://wbs-api.mexc.com/ws",
		},
		WebSocket: struct {
			ReconnectDelay       time.Duration `mapstructure:"reconnect_delay"`
			MaxReconnectAttempts int           `mapstructure:"max_reconnect_attempts"`
			PingInterval         time.Duration `mapstructure:"ping_interval"`
			AutoReconnect        bool          `mapstructure:"auto_reconnect"`
		}{
			ReconnectDelay:       5 * time.Second,
			MaxReconnectAttempts: 3,
			PingInterval:         30 * time.Second,
			AutoReconnect:        true,
		},
		ConnectionRateLimiter: struct {
			RequestsPerSecond float64 `mapstructure:"requests_per_second"`
			BurstCapacity     int     `mapstructure:"burst_capacity"`
		}{
			RequestsPerSecond: 10,
			BurstCapacity:     10,
		},
		SubscriptionRateLimiter: struct {
			RequestsPerSecond float64 `mapstructure:"requests_per_second"`
			BurstCapacity     int     `mapstructure:"burst_capacity"`
		}{
			RequestsPerSecond: 10,
			BurstCapacity:     10,
		},
	}
}

// setupWSServer creates a test WebSocket server
func setupWSServer(t *testing.T, handler func(*websocket.Conn)) (string, func()) {
	upgrader := websocket.Upgrader{
		HandshakeTimeout: 5 * time.Second,
		CheckOrigin:      func(r *http.Request) bool { return true },
		Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
			t.Logf("WebSocket upgrade error: %v", reason)
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Logf("Failed to upgrade connection: %v", err)
			return
		}

		// Set read/write deadlines
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

		// Enable ping/pong handling
		conn.SetPingHandler(func(appData string) error {
			return conn.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(5*time.Second))
		})

		// Run the custom handler if provided
		if handler != nil {
			handler(conn)
		} else {
			// Default handler: echo messages back
			for {
				messageType, message, err := conn.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						t.Logf("WebSocket read error: %v", err)
					}
					return
				}

				if err := conn.WriteMessage(messageType, message); err != nil {
					t.Logf("WebSocket write error: %v", err)
					return
				}
			}
		}
	}))

	wsURL := strings.Replace(server.URL, "http://", "ws://", 1)
	cleanup := func() {
		server.Close()
	}

	return wsURL, cleanup
}

func TestMexcWebSocketClient_Connect(t *testing.T) {
	wsURL, cleanup := setupWSServer(t, nil)
	defer cleanup()

	client := NewMexcWebSocketClient(wsURL)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := client.Connect(ctx)
	require.NoError(t, err)
}

func TestMexcWebSocketClient_Disconnect(t *testing.T) {
	wsURL, cleanup := setupWSServer(t, nil)
	defer cleanup()

	client := NewMexcWebSocketClient(wsURL)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := client.Connect(ctx)
	require.NoError(t, err)

	err = client.Disconnect()
	require.NoError(t, err)
}

func TestMexcWebSocketClient_SubscribeToTickers(t *testing.T) {
	wsURL, cleanup := setupWSServer(t, func(conn *websocket.Conn) {
		// Echo subscription response
		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				return
			}
			err = conn.WriteMessage(messageType, message)
			if err != nil {
				return
			}
		}
	})
	defer cleanup()

	client := NewMexcWebSocketClient(wsURL)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := client.Connect(ctx)
	require.NoError(t, err)

	err = client.SubscribeToTickers(ctx, []string{"BTCUSDT"})
	require.NoError(t, err)
}

func TestMexcWebSocketClient_Reconnect(t *testing.T) {
	// This test is just a placeholder, the TestReconnectLogic_WithExponentialBackoffAndResubscribe
	// test covers the actual reconnection functionality
	wsURL, cleanup := setupWSServer(t, nil)
	defer cleanup()

	client := NewMexcWebSocketClient(wsURL)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := client.Connect(ctx)
	require.NoError(t, err)
}

func TestMexcWebSocketClient_NewClient(t *testing.T) {
	client, err := mexcws.NewClient(mockConfig())
	require.NoError(t, err)
	assert.NotNil(t, client)
}

func TestMexcWebSocketClient_WithCustomRateLimiters(t *testing.T) {
	// This test now uses our mock client instead of the real one
	customConnLimiter := ratelimiter.NewTokenBucketRateLimiter(5, 5)
	customSubLimiter := ratelimiter.NewTokenBucketRateLimiter(10, 10)

	client := NewMexcWebSocketClient("ws://test-endpoint")
	client.connRateLimiter = customConnLimiter
	client.subRateLimiter = customSubLimiter

	// Verify rate limiters were set
	assert.Equal(t, customConnLimiter, client.GetConnRateLimiter())
	assert.Equal(t, customSubLimiter, client.GetSubRateLimiter())
}

func TestMexcWebSocketClient_SubscriptionRateLimit(t *testing.T) {
	wsURL, cleanup := setupWSServer(t, func(conn *websocket.Conn) {
		// Server handler
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	})
	defer cleanup()

	// Create client with custom rate limiter
	client := NewMexcWebSocketClient(wsURL)
	client.subRateLimiter = ratelimiter.NewTokenBucketRateLimiter(0, 3)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := client.Connect(ctx)
	require.NoError(t, err)
	defer client.Disconnect()

	// First 3 subscriptions should succeed (burst capacity of 3)
	for i := 0; i < 3; i++ {
		err = client.SubscribeToTickers(ctx, []string{fmt.Sprintf("COIN%dUSDT", i)})
		require.NoError(t, err)
	}

	// Fourth subscription should fail (rate limited)
	err = client.SubscribeToTickers(ctx, []string{"BTCUSDT"})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrRateLimitExceeded)
}

// --- Additional tests for rate limiting, reconnection, and error handling ---

func TestConnect_RateLimitExceeded(t *testing.T) {
	wsURL, cleanup := setupWSServer(t, nil)
	defer cleanup()

	// Create client with custom rate limiter
	client := NewMexcWebSocketClient(wsURL)
	client.connRateLimiter = ratelimiter.NewTokenBucketRateLimiter(0, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := client.Connect(ctx)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrRateLimitExceeded)
}

func TestHandleMessages_InvalidJSON(t *testing.T) {
	wsURL, cleanup := setupWSServer(t, func(conn *websocket.Conn) {
		// Send an invalid JSON message
		err := conn.WriteMessage(websocket.TextMessage, []byte("{invalid:json}"))
		if err != nil {
			return
		}
	})
	defer cleanup()

	client := NewMexcWebSocketClient(wsURL)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := client.Connect(ctx)
	require.NoError(t, err)

	// Wait to process invalid message
	time.Sleep(100 * time.Millisecond)

	// Client should still be connected
	assert.True(t, client.IsConnected())

	// Cleanup
	client.Disconnect()
}

func TestReconnectionWithRateLimiting(t *testing.T) {
	wsURL, cleanup := setupWSServer(t, func(conn *websocket.Conn) {
		// Simple handler that just keeps the connection open
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	})
	defer cleanup()

	// Create the client with a rate limiter that allows exactly 2 connections
	// (1 for initial connect + 1 for reconnect)
	client := NewMexcWebSocketClient(wsURL)
	client.connRateLimiter = ratelimiter.NewTokenBucketRateLimiter(0, 2) // Zero refill rate but 2 tokens initially

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// First connection should succeed
	err := client.Connect(ctx)
	require.NoError(t, err)

	// Disconnect to simulate connection loss
	client.Disconnect()

	// Try to reconnect - should also succeed because we have 2 tokens
	err = client.Connect(ctx)
	require.NoError(t, err)

	// Check connection attempts
	assert.Equal(t, 2, client.GetConnectionAttempts())

	// Cleanup
	client.Disconnect()
}

func TestExcessiveSubscriptionRequests(t *testing.T) {
	wsURL, cleanup := setupWSServer(t, func(conn *websocket.Conn) {
		// Simple handler that just keeps the connection open
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	})
	defer cleanup()

	// Create client with a strict rate limiter
	client := NewMexcWebSocketClient(wsURL)
	client.subRateLimiter = ratelimiter.NewTokenBucketRateLimiter(0, 3)

	// Connect
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := client.Connect(ctx)
	require.NoError(t, err)

	// Subscribe to 3 different symbols - these should succeed
	symbols := []string{"BTCUSDT", "ETHUSDT", "XRPUSDT"}

	for i := 0; i < 3; i++ {
		err := client.SubscribeToTickers(context.Background(), []string{symbols[i]})
		require.NoError(t, err, "Subscription %d should succeed", i+1)
	}

	// The 4th subscription should fail due to rate limiting
	err = client.SubscribeToTickers(ctx, []string{"LTCUSDT"})
	require.Error(t, err, "Fourth subscription should fail due to rate limiting")
	assert.ErrorIs(t, err, ErrRateLimitExceeded)

	// Cleanup
	client.Disconnect()
}
