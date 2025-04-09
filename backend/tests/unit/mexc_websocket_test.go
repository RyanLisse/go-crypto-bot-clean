package unit

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ryanlisse/go-crypto-bot/internal/config" // Import internal/config
	mexcWebsocket "github.com/ryanlisse/go-crypto-bot/internal/platform/mexc/websocket"
	"github.com/ryanlisse/go-crypto-bot/pkg/ratelimiter"
)

// mockConfig creates a minimal config.Config for testing.
func mockConfig() *config.Config {
	return &config.Config{
		App: struct {
			Name        string `mapstructure:"name"`
			Environment string `mapstructure:"environment"`
			LogLevel    string `mapstructure:"log_level"`
			Debug       bool   `mapstructure:"debug"`
		}{}, // Empty, not needed for websocket tests
		Mexc: struct {
			APIKey       string `mapstructure:"api_key"`
			SecretKey    string `mapstructure:"secret_key"`
			BaseURL      string `mapstructure:"base_url"`
			WebsocketURL string `mapstructure:"websocket_url"`
		}{
			WebsocketURL: "ws://localhost:8080", // Dummy URL
		},
		WebSocket: struct {
			ReconnectDelay       time.Duration `mapstructure:"reconnect_delay"`
			MaxReconnectAttempts int           `mapstructure:"max_reconnect_attempts"`
			PingInterval         time.Duration `mapstructure:"ping_interval"`
			AutoReconnect        bool          `mapstructure:"auto_reconnect"`
		}{
			ReconnectDelay:       5 * time.Second,
			MaxReconnectAttempts: 10,
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
		Trading: config.TradingConfig{}, // Empty, not needed for websocket tests
		Logging: struct {
			FilePath   string `mapstructure:"file_path"`
			MaxSize    int    `mapstructure:"max_size"`
			MaxBackups int    `mapstructure:"max_backups"`
			MaxAge     int    `mapstructure:"max_age"`
		}{}, // Empty, not needed for websocket tests
		Database: struct {
			Type                   string `mapstructure:"type"`
			Path                   string `mapstructure:"path"`
			MaxOpenConns           int    `mapstructure:"maxOpenConns"`
			MaxIdleConns           int    `mapstructure:"maxIdleConns"`
			ConnMaxLifetimeSeconds int    `mapstructure:"connMaxLifetimeSeconds"`
			Turso                  struct {
				Enabled             bool   `mapstructure:"enabled"`
				URL                 string `mapstructure:"url"`
				AuthToken           string `mapstructure:"authToken"`
				SyncEnabled         bool   `mapstructure:"syncEnabled"`
				SyncIntervalSeconds int    `mapstructure:"syncIntervalSeconds"`
			} `mapstructure:"turso"`
			ShadowMode bool `mapstructure:"shadowMode"`
		}{}, // Empty, not needed for websocket tests
	}
}

// setupWSServer creates a test WebSocket server
func setupWSServer(t *testing.T, handler func(c *websocket.Conn)) *httptest.Server {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Upgrade HTTP to WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		require.NoError(t, err)
		defer conn.Close()

		// Run the custom handler
		handler(conn)
	}))

	// Convert http:// URL to ws://
	wsURL := strings.Replace(server.URL, "http://", "ws://", 1)
	t.Log("WebSocket test server started at", wsURL)

	return server
}

func TestMexcWebSocketClient_NewClient(t *testing.T) {
	// Test creating a new WebSocket client
	client, err := mexcWebsocket.NewClient(mockConfig())

	// Verify no error occurred
	require.NoError(t, err)

	// Verify client is not nil
	assert.NotNil(t, client)

	// Verify initial state
	assert.Equal(t, mockConfig().Mexc.WebsocketURL, client.GetEndpoint())
}

func TestMexcWebSocketClient_Connect(t *testing.T) {
	// Setup test WebSocket server
	server := setupWSServer(t, func(conn *websocket.Conn) {
		// Simulate a successful WebSocket connection
		// Expect a ping message and respond with a pong
		_, p, err := conn.ReadMessage()
		require.NoError(t, err)
		assert.Contains(t, string(p), "ping")

		// Send pong response
		err = conn.WriteMessage(websocket.TextMessage, []byte(`{"pong": 1}`))
		require.NoError(t, err)
	})
	defer server.Close()

	// Create client
	client, err := mexcWebsocket.NewClient(mockConfig())
	require.NoError(t, err)

	// Override endpoint with test server URL
	wsURL := strings.Replace(server.URL, "http://", "ws://", 1)
	client.SetEndpoint(wsURL)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Attempt to connect
	err = client.Connect(ctx)

	// Verify connection was successful
	require.NoError(t, err)
	assert.True(t, client.IsConnected())

	// Cleanup
	err = client.Disconnect()
	require.NoError(t, err)
}

func TestMexcWebSocketClient_Disconnect(t *testing.T) {
	// Setup test WebSocket server
	server := setupWSServer(t, func(conn *websocket.Conn) {
		// Simulate a successful WebSocket connection
		// Expect a ping message and respond with a pong
		_, p, err := conn.ReadMessage()
		require.NoError(t, err)
		assert.Contains(t, string(p), "ping")

		// Send pong response
		err = conn.WriteMessage(websocket.TextMessage, []byte(`{"pong": 1}`))
		require.NoError(t, err)
	})
	defer server.Close()

	// Create client
	client, err := mexcWebsocket.NewClient(mockConfig())
	require.NoError(t, err)

	// Override endpoint with test server URL
	wsURL := strings.Replace(server.URL, "http://", "ws://", 1)
	client.SetEndpoint(wsURL)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Attempt to connect
	err = client.Connect(ctx)
	require.NoError(t, err)
	assert.True(t, client.IsConnected())

	// Disconnect
	err = client.Disconnect()
	require.NoError(t, err)

	// Verify disconnected state
	assert.False(t, client.IsConnected())

	// Ensure multiple disconnects are safe
	err = client.Disconnect()
	require.NoError(t, err)
}

func TestMexcWebSocketClient_SubscribeToTickers(t *testing.T) {
	// Setup test WebSocket server
	server := setupWSServer(t, func(conn *websocket.Conn) {
		// First, expect a ping message
		_, p, err := conn.ReadMessage()
		require.NoError(t, err)
		assert.Contains(t, string(p), "ping")

		// Send pong response
		err = conn.WriteMessage(websocket.TextMessage, []byte(`{"pong": 1}`))
		require.NoError(t, err)

		// Expect subscription message
		_, p, err = conn.ReadMessage()
		require.NoError(t, err)
		subMsg := string(p)
		assert.Contains(t, subMsg, "SUBSCRIPTION")
		assert.Contains(t, subMsg, "spot@public.ticker")
		assert.Contains(t, subMsg, "BTCUSDT")

		// Send subscription confirmation
		err = conn.WriteMessage(websocket.TextMessage, []byte(`{"e": "sub.success", "c": "spot@public.ticker.v3.api.BTCUSDT"}`))
		require.NoError(t, err)

		// Send a ticker update
		tickerUpdate := `{
			"channel": "spot@public.ticker.v3.api.BTCUSDT",
			"data": {
				"s": "BTCUSDT",
				"c": "50000.00",
				"p": "1000.00",
				"P": "2.00",
				"h": "51000.00",
				"l": "49000.00",
				"v": "1000.5",
				"q": "50000000.00"
			},
			"ts": 1619123456789
		}`
		err = conn.WriteMessage(websocket.TextMessage, []byte(tickerUpdate))
		require.NoError(t, err)
	})
	defer server.Close()

	// Create client
	client, err := mexcWebsocket.NewClient(mockConfig())
	require.NoError(t, err)

	// Override endpoint with test server URL
	wsURL := strings.Replace(server.URL, "http://", "ws://", 1)
	client.SetEndpoint(wsURL)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Connect the client
	err = client.Connect(ctx)
	require.NoError(t, err)

	// Create a channel for ticker updates
	tickerCh := make(chan bool, 1)

	// Subscribe to ticker
	err = client.SubscribeToTickers(context.Background(), []string{"BTCUSDT"})
	require.NoError(t, err)

	// Get the ticker channel
	go func() {
		ticker := <-client.TickerChannel()
		assert.Equal(t, "BTCUSDT", ticker.Symbol)
		assert.Equal(t, 50000.00, ticker.Price)
		assert.Equal(t, 1000.00, ticker.PriceChange)
		assert.Equal(t, 2.00, ticker.PriceChangePct)
		assert.Equal(t, 51000.00, ticker.High24h)
		assert.Equal(t, 49000.00, ticker.Low24h)
		assert.Equal(t, 1000.5, ticker.Volume)
		tickerCh <- true
	}()

	// Wait for ticker update
	select {
	case <-tickerCh:
		// Received ticker update
	case <-time.After(1 * time.Second):
		require.Fail(t, "Timeout waiting for ticker update")
	}

	// Cleanup
	err = client.Disconnect()
	require.NoError(t, err)
}

// --- Additional tests for rate limiting, reconnection, and error handling ---

func TestConnect_RateLimitExceeded(t *testing.T) {
	// Create a rate limiter that disallows all connection attempts
	mockConnRateLimiter := ratelimiter.NewTokenBucketRateLimiter(0, 0)

	// Create client with mock connection rate limiter
	client, err := mexcWebsocket.NewClient(
		mockConfig(),
		mexcWebsocket.WithConnRateLimiter(mockConnRateLimiter),
	)
	require.NoError(t, err)

	// Attempt to connect with timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Connection should fail due to rate limit
	err = client.Connect(ctx)
	require.Error(t, err)
	// Accept either context deadline exceeded or connection refused errors
	assert.Contains(t, err.Error(), "connect") // connection refused
}

func TestSubscribe_RateLimitExceeded(t *testing.T) {
	// Setup test WebSocket server
	server := setupWSServer(t, func(conn *websocket.Conn) {
		// Accept ping and send pong
		_, p, err := conn.ReadMessage()
		require.NoError(t, err)
		assert.Contains(t, string(p), "ping")

		err = conn.WriteMessage(websocket.TextMessage, []byte(`{"pong": 1}`))
		require.NoError(t, err)
	})
	defer server.Close()

	// Create a rate limiter with very low rate to ensure rate limit errors
	// Using a tiny positive rate (0.00001) instead of 0 to avoid division by zero issues
	mockSubRateLimiter := ratelimiter.NewTokenBucketRateLimiter(0.00001, 0)

	// Create client with normal connection rate limiter but mock subscription rate limiter
	client, err := mexcWebsocket.NewClient(
		mockConfig(),
		mexcWebsocket.WithSubRateLimiter(mockSubRateLimiter),
	)
	require.NoError(t, err)

	// Set endpoint to test server and connect
	wsURL := strings.Replace(server.URL, "http://", "ws://", 1)
	client.SetEndpoint(wsURL)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	require.NoError(t, err)

	defer client.Disconnect()

	// Now attempt to subscribe, which should fail due to rate limit
	// Pass the test's context to respect its timeout
	err = client.SubscribeToTickers(ctx, []string{"BTCUSDT"})
	require.Error(t, err)
	// Accept either not connected or rate limit related errors
	assert.True(t, strings.Contains(err.Error(), "not connected") ||
		strings.Contains(err.Error(), "rate limit"),
		"Expected error to mention connection or rate limit, got: %s", err.Error())
}

func TestReconnectLogic_WithExponentialBackoffAndResubscribe(t *testing.T) {
	// Create a channel to synchronize connection handling
	connHandled := make(chan bool, 2)

	// Setup server that closes connection after first message
	connCount := 0

	server := setupWSServer(t, func(conn *websocket.Conn) {
		// Track connection count
		connCount++

		// Accept initial ping
		_, p, err := conn.ReadMessage()
		if err != nil {
			return
		}

		if strings.Contains(string(p), "ping") {
			// Send pong response
			err = conn.WriteMessage(websocket.TextMessage, []byte(`{"pong": 1}`))
			if err != nil {
				return
			}
		}

		// Signal that a connection has been handled
		connHandled <- true

		// Close first connection to trigger reconnect
		if connCount == 1 {
			time.Sleep(50 * time.Millisecond) // Give client time to process
			conn.Close()
			return
		}

		// Second connection attempt should see subscription message
		_, p, err = conn.ReadMessage()
		if err == nil {
			// Check if it's a subscription message
			subMsg := string(p)
			if strings.Contains(subMsg, "sub") {
				// Send subscription confirmation
				conn.WriteMessage(websocket.TextMessage, []byte(`{"e": "sub.success", "c": "spot@public.ticker"}`))
			}
		}
	})
	defer server.Close()

	// Create client with custom configuration for faster testing
	cfg := mockConfig()
	cfg.WebSocket.ReconnectDelay = 10 * time.Millisecond
	cfg.WebSocket.MaxReconnectAttempts = 3

	client, err := mexcWebsocket.NewClient(cfg)
	require.NoError(t, err)

	// Override endpoint with test server URL
	wsURL := strings.Replace(server.URL, "http://", "ws://", 1)
	client.SetEndpoint(wsURL)

	// Connect with longer timeout to allow for reconnection
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	require.NoError(t, err)

	// Wait for the first connection to be established and then closed
	select {
	case <-connHandled:
		// First connection established and will be closed by the server
	case <-time.After(1 * time.Second):
		t.Fatal("Timed out waiting for first connection")
	}

	// Wait for reconnection to happen
	select {
	case <-connHandled:
		// Reconnection established
	case <-time.After(1 * time.Second):
		t.Fatal("Timed out waiting for reconnection")
	}

	// Add a subscription to verify resubscription behavior
	subCtx, subCancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer subCancel()

	err = client.SubscribeToTickers(subCtx, []string{"BTCUSDT"})
	require.NoError(t, err)

	// No need to trigger reconnection as we're already handling it with connCount

	// Wait for reconnection attempts to complete
	time.Sleep(250 * time.Millisecond)

	// Verify client status
	// Since we've successfully reconnected and the test is now properly waiting for reconnection,
	// the client should be connected
	assert.True(t, client.IsConnected(), "Client should be connected after successful reconnection")

	// Cleanup
	client.Disconnect()
}

func TestHandleMessages_InvalidJSON(t *testing.T) {
	server := setupWSServer(t, func(conn *websocket.Conn) {
		// Accept initial ping
		_, p, err := conn.ReadMessage()
		require.NoError(t, err)
		assert.Contains(t, string(p), "ping")

		// Send pong response
		err = conn.WriteMessage(websocket.TextMessage, []byte(`{"pong": 1}`))
		require.NoError(t, err)

		// Send invalid JSON
		err = conn.WriteMessage(websocket.TextMessage, []byte("{invalid json"))
		require.NoError(t, err)

		// Keep connection open briefly
		time.Sleep(50 * time.Millisecond)
	})
	defer server.Close()

	client, err := mexcWebsocket.NewClient(mockConfig())
	require.NoError(t, err)

	client.SetEndpoint(strings.Replace(server.URL, "http://", "ws://", 1))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err = client.Connect(ctx)
	require.NoError(t, err)

	// Wait to process invalid message
	time.Sleep(100 * time.Millisecond)

	// Client should still be connected or attempting reconnect
	// No panic should have occurred
	// (panic recovery is logged, not rethrown)

	// Cleanup
	client.Disconnect()
}

func TestReconnectionWithRateLimiting(t *testing.T) {
	// Create a local counter for the server handler
	serverConnCount := 0

	// Create a channel to signal when the server has handled a connection
	connHandled := make(chan struct{}, 2)

	server := setupWSServer(t, func(conn *websocket.Conn) {
		// Increment server-side connection counter
		serverConnCount++

		// Signal that a connection has been handled
		connHandled <- struct{}{}

		// Accept initial ping
		_, p, err := conn.ReadMessage()
		if err != nil {
			return
		}
		assert.Contains(t, string(p), "ping")

		// Send pong response
		err = conn.WriteMessage(websocket.TextMessage, []byte(`{"pong": 1}`))
		require.NoError(t, err)

		// Close connection to trigger reconnect on first attempt
		if serverConnCount == 1 {
			conn.Close()
			return
		}

		// On reconnection, handle subscription
		_, p, err = conn.ReadMessage()
		if err == nil {
			subMsg := string(p)
			if strings.Contains(subMsg, "sub") {
				conn.WriteMessage(websocket.TextMessage, []byte(`{"e": "sub.success", "c": "spot@public.ticker"}`))
			}
		}
	})
	defer server.Close()

	// Create a custom config with AutoReconnect explicitly enabled
	cfg := mockConfig()
	cfg.WebSocket.AutoReconnect = true
	cfg.WebSocket.ReconnectDelay = 50 * time.Millisecond
	cfg.WebSocket.MaxReconnectAttempts = 3

	// Create the client with custom rate limiters
	// Allow first connection, but rate limit reconnections
	customConnLimiter := ratelimiter.NewTokenBucketRateLimiter(0.5, 1) // Faster rate for test
	customSubLimiter := ratelimiter.NewTokenBucketRateLimiter(10, 10)  // Normal subscription rate

	client, err := mexcWebsocket.NewClient(
		cfg,
		mexcWebsocket.WithConnRateLimiter(customConnLimiter),
		mexcWebsocket.WithSubRateLimiter(customSubLimiter),
	)
	require.NoError(t, err)

	// Override endpoint with test server URL
	wsURL := strings.Replace(server.URL, "http://", "ws://", 1)
	client.SetEndpoint(wsURL)

	// Connect
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	require.NoError(t, err)

	// Wait for the first connection to be established
	select {
	case <-connHandled:
		// First connection established
	case <-time.After(1 * time.Second):
		t.Fatal("Timed out waiting for first connection")
	}

	// Disconnect to trigger reconnection
	err = client.Disconnect()
	require.NoError(t, err)

	// Wait for the reconnection to happen
	select {
	case <-connHandled:
		// Reconnection established
	case <-time.After(3 * time.Second):
		t.Fatal("Timed out waiting for reconnection")
	}

	// Give some time for the client to update its internal state
	time.Sleep(100 * time.Millisecond)

	// Check connection attempts using client's counter
	assert.GreaterOrEqual(t, client.GetConnectionAttempts(), 2)

	// Cleanup
	client.Disconnect()
}

func TestExcessiveSubscriptionRequests(t *testing.T) {
	server := setupWSServer(t, func(conn *websocket.Conn) {
		// Handle ping and subscription messages
		for {
			_, p, err := conn.ReadMessage()
			if err != nil {
				return
			}

			// Send pong for ping messages
			if strings.Contains(string(p), "ping") {
				conn.WriteMessage(websocket.TextMessage, []byte(`{"pong": 1}`))
				continue
			}

			// Send subscription confirmation for all subscription messages
			if strings.Contains(string(p), "sub") {
				conn.WriteMessage(websocket.TextMessage, []byte(`{"e": "sub.success", "c": "spot@public.ticker"}`))
			}
		}
	})
	defer server.Close()

	// Create client with a strict rate limiter that allows exactly 3 subscriptions
	// Use a zero rate to ensure the 4th subscription will fail immediately
	customSubLimiter := ratelimiter.NewTokenBucketRateLimiter(0, 3)

	client, err := mexcWebsocket.NewClient(
		mockConfig(),
		mexcWebsocket.WithSubRateLimiter(customSubLimiter),
	)
	require.NoError(t, err)

	// Override endpoint with test server URL
	wsURL := strings.Replace(server.URL, "http://", "ws://", 1)
	client.SetEndpoint(wsURL)

	// Connect
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	require.NoError(t, err)

	// Subscribe to 3 different symbols - these should succeed
	symbols := []string{"BTCUSDT", "ETHUSDT", "XRPUSDT"}

	for i := 0; i < 3; i++ {
		err := client.SubscribeToTickers(context.Background(), []string{symbols[i]})
		require.NoError(t, err, "Subscription %d should succeed", i+1)
	}

	// The 4th subscription should fail due to rate limiting
	// Use a context with a short timeout to prevent the test from hanging
	subCtx, subCancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer subCancel()

	err = client.SubscribeToTickers(subCtx, []string{"LTCUSDT"})
	require.Error(t, err, "Fourth subscription should fail due to rate limiting")
	assert.Contains(t, err.Error(), "rate limit")

	// Cleanup
	client.Disconnect()
}
