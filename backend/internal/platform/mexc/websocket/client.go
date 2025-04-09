package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/ryanlisse/go-crypto-bot/internal/config"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
	"github.com/ryanlisse/go-crypto-bot/pkg/ratelimiter"
	"go.uber.org/zap"
)

const (
	// Default WebSocket endpoints
	DefaultEndpoint = "wss://stream.mexc.com/ws"

	// WebSocket connection parameters
	defaultReadTimeout  = 45 * time.Second
	defaultWriteTimeout = 10 * time.Second
	defaultPingInterval = 30 * time.Second

	// Reconnection parameters
	defaultReconnectDelay    = 5 * time.Second
	defaultMaxReconnectDelay = 60 * time.Second
	defaultMaxReconnects     = 10
)

// Client represents a MEXC WebSocket client
type Client struct {
	cfg             *config.Config
	connRateLimiter *ratelimiter.TokenBucketRateLimiter
	subRateLimiter  *ratelimiter.TokenBucketRateLimiter

	// Context for cancellation
	ctx    context.Context
	cancel context.CancelFunc

	// WebSocket connection
	conn          *websocket.Conn
	connMutex     sync.RWMutex // Mutex for connection operations
	endpoint      string
	dialer        *websocket.Dialer
	apiKey        string
	secretKey     string
	listenKey     string
	authenticated bool
	connected     bool

	// Connection management
	stopCh         chan struct{}
	reconnecting   bool
	reconnectMutex sync.Mutex
	mu             sync.Mutex // General mutex for client operations
	pingTicker     *time.Ticker
	pingStopCh     chan struct{} // Channel to stop ping ticker goroutine
	pongCh         chan struct{} // Channel to receive pong responses
	lastPongTime   time.Time     // Time of last pong received
	pongTimeout    time.Duration // How long to wait for a pong before considering connection dead

	// For testing purposes
	connectionAttempts int
	connAttemptsMutex  sync.Mutex
	pingSentCount      int
	pingMutex          sync.Mutex

	// Data channels
	tickerCh chan *models.Ticker

	// Track subscriptions for resubscription after reconnect
	subscribedSymbols []string

	// Account updates handler
	accountHandler *AccountHandler

	// Reconnect handler
	reconnectHandler func() error

	// Logger
	logger *zap.Logger
}

// NewClient creates a new MEXC WebSocket client with options
func NewClient(cfg *config.Config, options ...ClientOption) (*Client, error) {
	connLimiter := ratelimiter.NewTokenBucketRateLimiter(
		cfg.ConnectionRateLimiter.RequestsPerSecond,
		float64(cfg.ConnectionRateLimiter.BurstCapacity),
	)
	subLimiter := ratelimiter.NewTokenBucketRateLimiter(
		cfg.SubscriptionRateLimiter.RequestsPerSecond,
		float64(cfg.SubscriptionRateLimiter.BurstCapacity),
	)

	// Set defaults for new configuration options
	if cfg.WebSocket.AutoReconnect == false {
		// Default to auto-reconnect unless explicitly disabled
		cfg.WebSocket.AutoReconnect = true
	}

	// Create a background context that will live for the client's lifetime
	ctx, cancel := context.WithCancel(context.Background())

	// Create logger
	logger, _ := zap.NewProduction()

	// Set default ping interval if not configured
	if cfg.WebSocket.PingInterval <= 0 {
		cfg.WebSocket.PingInterval = defaultPingInterval
	}

	// Set default pong timeout
	pongTimeout := cfg.WebSocket.PingInterval * 2
	if pongTimeout < 10*time.Second {
		pongTimeout = 10 * time.Second // Minimum 10 seconds
	}

	client := &Client{
		cfg:             cfg,
		connRateLimiter: connLimiter,
		subRateLimiter:  subLimiter,
		endpoint:        cfg.Mexc.WebsocketURL,
		apiKey:          cfg.Mexc.APIKey,
		secretKey:       cfg.Mexc.SecretKey,
		dialer:          websocket.DefaultDialer,
		stopCh:          make(chan struct{}),
		ctx:             ctx,
		cancel:          cancel,
		tickerCh:        make(chan *models.Ticker, 100),
		logger:          logger,
		pingSentCount:   0,
		pongCh:          make(chan struct{}, 10),
		lastPongTime:    time.Now(),
		pongTimeout:     pongTimeout,
	}

	// Apply options
	for _, option := range options {
		option(client)
	}

	// Create account handler with the same logger
	client.accountHandler = NewAccountHandler(client, client.logger)

	return client, nil
}

// GetAPIKey returns the API key
func (c *Client) GetAPIKey() string {
	return c.apiKey
}

// GetSecretKey returns the secret key
func (c *Client) GetSecretKey() string {
	return c.secretKey
}

// GetListenKey returns the current listen key
func (c *Client) GetListenKey() string {
	c.connMutex.RLock()
	defer c.connMutex.RUnlock()
	return c.listenKey
}

// Authenticate authenticates the WebSocket connection using the listen key
func (c *Client) Authenticate(ctx context.Context) error {
	// Check if we have a listen key
	if c.listenKey == "" {
		return errors.New("no listen key available for authentication")
	}

	// Check if we're connected
	c.connMutex.RLock()
	if c.conn == nil || !c.connected {
		c.connMutex.RUnlock()
		return errors.New("not connected to WebSocket")
	}
	c.connMutex.RUnlock()

	// Create authentication message
	authMsg := map[string]interface{}{
		"method": "SUBSCRIBE",
		"params": []string{"spot@private"}, // Private channel subscription
		"id":     time.Now().UnixNano(),
		"key":    c.listenKey,
	}

	// Send authentication message
	c.connMutex.Lock()
	defer c.connMutex.Unlock()

	c.conn.SetWriteDeadline(time.Now().Add(defaultWriteTimeout))
	if err := c.conn.WriteJSON(authMsg); err != nil {
		return fmt.Errorf("failed to send authentication message: %w", err)
	}

	// Mark as authenticated
	c.authenticated = true
	return nil
}

// GetEndpoint returns the current WebSocket endpoint
func (c *Client) GetEndpoint() string {
	return c.endpoint
}

// SetEndpoint allows changing the WebSocket endpoint
func (c *Client) SetEndpoint(endpoint string) {
	c.endpoint = endpoint
}

// SetReconnectDelay sets the delay between reconnection attempts
func (c *Client) SetReconnectDelay(delay time.Duration) {
	// Update the config instead of a direct field
	c.cfg.WebSocket.ReconnectDelay = delay
}

// SetReconnectHandler sets a custom handler for reconnection events
func (c *Client) SetReconnectHandler(handler func() error) {
	c.reconnectHandler = handler
}

// Connect establishes a WebSocket connection
// Constants for WebSocket operations
const (
	writeWait = 10 * time.Second // Time allowed to write a message to the peer
	pongWait = 60 * time.Second  // Time allowed to read the next pong message from the peer
	pingPeriod = 30 * time.Second // Send pings to peer with this period (must be less than pongWait)
)

// Connect establishes a WebSocket connection
func (c *Client) Connect(ctx context.Context) error {
	c.connMutex.Lock()
	defer c.connMutex.Unlock()

	// Check if already connected
	if c.connected {
		return nil
	}

	// Increment connection attempts (used for testing)
	c.connAttemptsMutex.Lock()
	c.connectionAttempts++
	c.connAttemptsMutex.Unlock()

	// Check if the rate limiter allows us to connect
	if err := c.connRateLimiter.Wait(ctx); err != nil {
		c.logger.Error("Connection rate limited",
			zap.Error(err),
			zap.Float64("rate", c.connRateLimiter.GetRate()))
		return fmt.Errorf("connection rate limited: %w", err)
	}

	// Set up WebSocket connection
	endpoint := c.endpoint
	c.logger.Info("Connecting to WebSocket", zap.String("endpoint", endpoint))

	// Create dialer if not exists
	if c.dialer == nil {
		c.dialer = websocket.DefaultDialer
	}

	// Connect to WebSocket
	conn, _, err := c.dialer.DialContext(ctx, endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to WebSocket: %w", err)
	}

	// Set ping handler
	conn.SetPingHandler(func(message string) error {
		c.logger.Debug("Received ping", zap.String("message", message))
		err := conn.WriteControl(websocket.PongMessage, []byte(message), time.Now().Add(writeWait))
		if err != nil {
			c.logger.Warn("Failed to send pong response", zap.Error(err))
		}
		return nil
	})

	// Set pong handler
	conn.SetPongHandler(func(string) error {
		c.mu.Lock()
		c.lastPongTime = time.Now()
		c.mu.Unlock()
		return nil
	})

	// Store the connection
	c.conn = conn
	c.connected = true

	// Start ping ticker
	c.startPingTicker()

	// Send initial ping to verify connection
	c.sendPing()

	// Start message handling goroutine
	go c.handleMessages()

	return nil
}

// Disconnect closes the WebSocket connection
func (c *Client) Disconnect() error {
	c.connMutex.Lock()
	defer c.connMutex.Unlock()

	if !c.connected || c.conn == nil {
		return nil
	}

	// Stop ping ticker if running
	c.stopPingTicker()

	// Close connection
	c.connected = false
	err := c.conn.Close()
	if err != nil {
		c.logger.Warn("Error closing WebSocket connection", zap.Error(err))
	}

	// Cancel context if we own it
	if c.cancel != nil {
		c.cancel()
		c.cancel = nil
	}

	c.conn = nil
	c.logger.Info("WebSocket disconnected")
	return err
}

// SubscribeToTickers subscribes to ticker updates for given symbols
func (c *Client) SubscribeToTickers(ctx context.Context, symbols []string) error {
	// Use the provided context for rate limiting
	if err := c.subRateLimiter.Wait(ctx); err != nil {
		return fmt.Errorf("subscription rate limit exceeded: %w", err)
	}

	c.connMutex.Lock()
	defer c.connMutex.Unlock()

	if c.conn == nil || !c.connected {
		return fmt.Errorf("not connected to WebSocket")
	}

	// Prepare subscription message
	subMsg := map[string]interface{}{
		"method": "SUBSCRIPTION",
		"params": []string{},
	}

	// Add each symbol to subscription
	for _, symbol := range symbols {
		subMsg["params"] = append(subMsg["params"].([]string), fmt.Sprintf("spot@public.deals.v3.api@%s", symbol))
	}

	// Marshal and send subscription message
	data, err := json.Marshal(subMsg)
	if err != nil {
		return fmt.Errorf("failed to create subscription message: %w", err)
	}

	if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return fmt.Errorf("failed to send subscription message: %w", err)
	}

	// Track subscriptions for resubscription after reconnect
	c.subscribedSymbols = symbols

	return nil
}

// UnsubscribeFromTickers unsubscribes from ticker updates for given symbols
func (c *Client) UnsubscribeFromTickers(symbols []string) error {
	c.connMutex.Lock()
	defer c.connMutex.Unlock()

	if c.conn == nil || !c.connected {
		return fmt.Errorf("not connected to WebSocket")
	}

	// Prepare unsubscription message
	unsubMsg := map[string]interface{}{
		"method": "UNSUBSCRIPTION",
		"params": []string{},
	}

	// Add each symbol to unsubscription
	for _, symbol := range symbols {
		unsubMsg["params"] = append(unsubMsg["params"].([]string), fmt.Sprintf("spot@public.ticker.v3.api.%s", symbol))
	}

	// Marshal and send unsubscription message
	data, err := json.Marshal(unsubMsg)
	if err != nil {
		return fmt.Errorf("failed to create unsubscription message: %w", err)
	}

	if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return fmt.Errorf("failed to send unsubscription message: %w", err)
	}

	return nil
}

// TickerChannel returns the channel for receiving ticker updates
func (c *Client) TickerChannel() <-chan *models.Ticker {
	return c.tickerCh
}

// GetConnectionAttempts returns the number of connection attempts
// This is primarily used for testing
func (c *Client) GetConnectionAttempts() int {
	c.connAttemptsMutex.Lock()
	defer c.connAttemptsMutex.Unlock()
	return c.connectionAttempts
}

// IsConnected is implemented in account_methods.go

// handleMessages processes incoming WebSocket messages
func (c *Client) handleMessages() {
	defer func() {
		if r := recover(); r != nil {
			// Log the panic
			c.logger.Error("Recovered from panic in handleMessages",
				zap.Any("panic", r))
		}
	}()

	c.logger.Info("Starting WebSocket message handler")

	for {
		select {
		case <-c.stopCh:
			c.logger.Info("Message handler stopped: stop channel closed")
			return
		case <-c.ctx.Done():
			// Context cancelled (e.g., parent context of Connect was cancelled)
			c.logger.Info("Message handler stopped: context cancelled",
				zap.Error(c.ctx.Err()))
			return
		default:
			// Check if connection is still open
			c.connMutex.RLock()
			if c.conn == nil {
				c.connMutex.RUnlock()
				// Connection closed, attempt to reconnect
				c.logger.Warn("WebSocket connection is nil, attempting to reconnect")
				if err := c.reconnect(); err != nil {
					// If reconnect fails and it's not a connection issue, don't retry
					if !errors.Is(err, context.Canceled) {
						c.logger.Error("Failed to reconnect", zap.Error(err))
					}
					return // Exit the message handler if reconnect fails
				}
				continue
			}
			c.connMutex.RUnlock()

			// Set read deadline
			c.connMutex.Lock()
			if c.conn != nil {
				c.conn.SetReadDeadline(time.Now().Add(defaultReadTimeout))
			}
			c.connMutex.Unlock()

			// Read message
			c.connMutex.RLock()
			conn := c.conn // Store locally to avoid race conditions
			c.connMutex.RUnlock()

			if conn == nil {
				// Connection was closed between checks
				continue
			}

			_, message, err := conn.ReadMessage()
			if err != nil {
				// Handle connection errors
				if websocket.IsUnexpectedCloseError(err,
					websocket.CloseGoingAway,
					websocket.CloseAbnormalClosure) {
					c.logger.Error("WebSocket read error", zap.Error(err))
				} else {
					c.logger.Warn("WebSocket connection closed", zap.Error(err))
				}

				// Trigger reconnection
				go c.reconnect()
				return
			}

			// Process the message
			c.processMessage(message)
		}
	}
}

// processMessage handles different types of messages from the WebSocket
func (c *Client) processMessage(message []byte) {
	// Log raw message for debugging (at debug level)
	if c.logger.Core().Enabled(zap.DebugLevel) {
		c.logger.Debug("Received WebSocket message",
			zap.String("message", string(message)))
	}

	// Try to parse the message as JSON
	var msg map[string]interface{}
	if err := json.Unmarshal(message, &msg); err != nil {
		c.logger.Error("Error unmarshaling message",
			zap.Error(err),
			zap.String("message", string(message)))
		return
	}

	// Handle ping messages
	if _, ok := msg["ping"]; ok {
		c.logger.Debug("Received ping message")
		c.handlePing(msg)
		return
	}

	// Handle pong messages
	if _, ok := msg["pong"]; ok {
		c.logger.Debug("Received pong message")
		// Signal pong received
		select {
		case c.pongCh <- struct{}{}:
			// Successfully sent pong notification
		default:
			// Channel full, but we still updated lastPongTime
		}
		c.lastPongTime = time.Now()
		return
	}

	// Check for account updates
	if c.accountHandler != nil {
		if processed := c.accountHandler.ProcessMessage(msg); processed {
			c.logger.Debug("Processed account update message")
			return
		}
	}

	// Handle ticker updates
	if channel, ok := msg["channel"].(string); ok {
		// Check if it's a ticker channel
		if strings.Contains(channel, "spot@public.deals.v3.api@") {
			c.logger.Debug("Processing ticker update",
				zap.String("channel", channel))
			c.handleTickerUpdate(msg)
			return
		}
	}

	// If we get here, we don't know how to handle this message type
	c.logger.Debug("Received unknown message type",
		zap.Any("message", msg))
}

// handlePing responds to ping messages from the server
func (c *Client) handlePing(msg map[string]interface{}) {
	pongMsg := map[string]interface{}{
		"pong": msg["ping"],
	}

	data, err := json.Marshal(pongMsg)
	if err != nil {
		return
	}

	c.connMutex.Lock()
	defer c.connMutex.Unlock()

	if c.conn == nil {
		return
	}

	c.conn.WriteMessage(websocket.TextMessage, data)
}

// handleTickerUpdate processes ticker update messages
func (c *Client) handleTickerUpdate(msg map[string]interface{}) {
	// Extract symbol from channel
	channel, ok := msg["channel"].(string)
	if !ok {
		return
	}

	// Extract symbol from channel (format: spot@public.deals.v3.api@SYMBOL)
	parts := strings.Split(channel, "@")
	if len(parts) < 3 {
		return
	}
	symbol := parts[len(parts)-1]

	// Extract data from the message
	data, ok := msg["data"].(map[string]interface{})
	if !ok {
		return
	}

	// Parse ticker data
	ticker := &models.Ticker{
		Symbol:         symbol,
		Price:          parseFloat(data["c"]),
		PriceChange:    parseFloat(data["p"]),
		PriceChangePct: parseFloat(data["P"]),
		High24h:        parseFloat(data["h"]),
		Low24h:         parseFloat(data["l"]),
		Volume:         parseFloat(data["v"]),
		QuoteVolume:    parseFloat(data["q"]),
		Timestamp:      time.Now(), // Use current time as timestamp
	}

	// Send ticker to channel
	select {
	case c.tickerCh <- ticker:
	default:
		// Channel is full, drop the message
	}
}

// startPingTicker starts a ticker to send periodic ping messages
func (c *Client) startPingTicker() {
	c.mu.Lock()
	// Stop existing ticker if any
	if c.pingTicker != nil {
		c.pingTicker.Stop()
	}

	// Create new ticker
	c.pingTicker = time.NewTicker(c.cfg.WebSocket.PingInterval)
	localPingTicker := c.pingTicker // Create a local reference to avoid nil pointer issues
	
	// Reset last pong time
	c.lastPongTime = time.Now()

	// Create a stop channel for this ticker goroutine
	pingStopCh := make(chan struct{})
	c.pingStopCh = pingStopCh
	c.mu.Unlock()

	// Start goroutine to send pings
	go func() {
		defer func() {
			if r := recover(); r != nil {
				c.logger.Error("Recovered from panic in ping ticker", zap.Any("panic", r))
			}
		}()

		for {
			select {
			case <-pingStopCh:
				return
			case <-localPingTicker.C:
				// Use the local reference instead of accessing c.pingTicker directly
				// to avoid race conditions and nil pointer derefs
				
				c.mu.Lock()
				lastPongTime := c.lastPongTime
				pongTimeout := c.pongTimeout
				c.mu.Unlock()
				
				if time.Since(lastPongTime) > pongTimeout {
					c.logger.Warn("No pong received within timeout period, reconnecting",
						zap.Duration("timeout", pongTimeout),
						zap.Time("last_pong", lastPongTime))

					// Trigger reconnection
					go c.Disconnect()
					return
				}

				// Send ping message
				c.sendPing()
			case <-c.pongCh:
				// Update last pong time
				c.mu.Lock()
				c.lastPongTime = time.Now()
				c.mu.Unlock()
				c.logger.Debug("Pong received, connection healthy")
			case <-c.stopCh:
				// Stop ticker when client is disconnected
				localPingTicker.Stop()
				return
			case <-c.ctx.Done():
				// Stop ticker when context is cancelled
				localPingTicker.Stop()
				return
			}
		}
	}()
}

// stopPingTicker stops the ping ticker and its goroutine
func (c *Client) stopPingTicker() {
	// Stop the ticker
	if c.pingTicker != nil {
		c.pingTicker.Stop()
		c.pingTicker = nil
	}

	// Signal the ping goroutine to stop
	c.mu.Lock()
	if c.pingStopCh != nil {
		close(c.pingStopCh)
		c.pingStopCh = nil
	}
	c.mu.Unlock()
}

// sendPing sends a ping message to the server
func (c *Client) sendPing() {
	// Create ping message
	pingMsg := map[string]interface{}{
		"ping": time.Now().UnixNano() / int64(time.Millisecond),
	}

	// Marshal ping message
	data, err := json.Marshal(pingMsg)
	if err != nil {
		c.logger.Error("Failed to marshal ping message", zap.Error(err))
		return
	}

	// Send ping message
	c.connMutex.Lock()

	if c.conn == nil || !c.connected {
		c.connMutex.Unlock() // Unlock if connection is already gone
		return
	}

	// Store connection locally to use safely
	conn := c.conn

	// Set write deadline
	conn.SetWriteDeadline(time.Now().Add(defaultWriteTimeout))

	// Send ping message
	err = conn.WriteMessage(websocket.TextMessage, data)

	// MUST unlock before potentially calling Disconnect
	c.connMutex.Unlock()

	if err != nil {
		c.logger.Error("Failed to send ping message", zap.Error(err))
		// Trigger reconnection now that mutex is released
		go c.Disconnect() // <-- Now safe to call
		return // Return after signaling disconnect
	}

	// Increment ping sent counter (for testing) - outside connMutex scope
	c.pingMutex.Lock()
	c.pingSentCount++
	c.pingMutex.Unlock()
}

// GetPingSentCount returns the number of ping messages sent (for testing)
func (c *Client) GetPingSentCount() int {
	c.pingMutex.Lock()
	defer c.pingMutex.Unlock()
	return c.pingSentCount
}

// reconnect attempts to reestablish the WebSocket connection
func (c *Client) reconnect() error {
	c.reconnectMutex.Lock()
	defer c.reconnectMutex.Unlock()

	// Check if context was cancelled
	select {
	case <-c.ctx.Done():
		return c.ctx.Err()
	default:
		// Continue with reconnection
	}

	if c.reconnecting {
		return nil
	}
	c.reconnecting = true
	defer func() { c.reconnecting = false }()

	// Increment connection attempts (thread-safe)
	c.connAttemptsMutex.Lock()
	c.connectionAttempts++
	c.connAttemptsMutex.Unlock()

	// Log reconnection attempt
	c.logger.Info("Attempting to reconnect to WebSocket")

	// Stop ping ticker if running
	c.mu.Lock()
	if c.pingTicker != nil {
		c.pingTicker.Stop()
		c.pingTicker = nil
	}
	c.mu.Unlock()

	// Close connection directly without triggering another reconnect
	c.connMutex.Lock()
	if c.conn != nil {
		_ = c.conn.Close()
		c.conn = nil
		c.connected = false
		c.authenticated = false // Reset authentication state
	}
	c.connMutex.Unlock()

	// Create a new stop channel before reconnecting
	c.mu.Lock()
	oldStopCh := c.stopCh
	c.stopCh = make(chan struct{})
	c.mu.Unlock()

	// Close old stop channel if it exists
	if oldStopCh != nil {
		select {
		case <-oldStopCh: // Already closed
		default:
			close(oldStopCh)
		}
	}

	// Create a context with timeout for reconnection
	ctx, cancel := context.WithTimeout(c.ctx, time.Minute)
	defer cancel()

	// Check if we have a custom reconnect handler
	if c.reconnectHandler != nil {
		c.logger.Info("Using custom reconnect handler")
		err := c.reconnectHandler()
		if err != nil {
			c.logger.Error("Custom reconnect handler failed", zap.Error(err))
			return err
		}
		c.logger.Info("Custom reconnect handler succeeded")
		return nil
	}

	// Attempt to reconnect with exponential backoff
	for i := 0; i < c.cfg.WebSocket.MaxReconnectAttempts; i++ {
		// Check for context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Continue reconnection attempt
		}

		// Wait before reconnecting
		delay := c.cfg.WebSocket.ReconnectDelay * time.Duration(1<<uint(i))
		maxDelay := defaultMaxReconnectDelay // Use the default max delay
		if delay > maxDelay {
			delay = maxDelay
		}

		c.logger.Info("Waiting before reconnection attempt",
			zap.Int("attempt", i+1),
			zap.Duration("delay", delay))

		// Use timer instead of Sleep to respect context cancellation
		timer := time.NewTimer(delay)
		select {
		case <-timer.C:
			// Timer completed
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		}

		// Try to connect
		c.logger.Info("Attempting reconnection", zap.Int("attempt", i+1))
		err := c.Connect(ctx)
		if err == nil {
			// Successfully reconnected
			c.logger.Info("Successfully reconnected to WebSocket")

			// For test environments, we might need to give a short delay
			// to ensure the connection is fully established
			time.Sleep(time.Millisecond * 10)

			// Re-authenticate if needed
			if c.listenKey != "" {
				c.logger.Info("Re-authenticating WebSocket connection")
				if authErr := c.Authenticate(ctx); authErr != nil {
					c.logger.Error("Re-authentication failed", zap.Error(authErr))
					// Continue with reconnection even if authentication fails
				}
			}

			// Resubscribe to previous subscriptions
			if len(c.subscribedSymbols) > 0 {
				c.logger.Info("Resubscribing to previous subscriptions",
					zap.Strings("symbols", c.subscribedSymbols))
				subErr := c.SubscribeToTickers(ctx, c.subscribedSymbols)
				if subErr != nil {
					// Log subscription error but don't fail reconnection
					c.logger.Error("Resubscription failed", zap.Error(subErr))
				}
			}
			return nil
		}

		c.logger.Error("Reconnection attempt failed",
			zap.Int("attempt", i+1),
			zap.Error(err))
	}

	// If we get here, all reconnection attempts failed
	c.logger.Error("Failed to reconnect after maximum attempts",
		zap.Int("maxAttempts", c.cfg.WebSocket.MaxReconnectAttempts))

	// If all reconnection attempts fail, mark as not connected
	c.connMutex.Lock()
	c.connected = false
	c.authenticated = false
	c.connMutex.Unlock()

	return fmt.Errorf("failed to reconnect after %d attempts", c.cfg.WebSocket.MaxReconnectAttempts)
}
