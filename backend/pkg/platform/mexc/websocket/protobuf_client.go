package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	mexcproto "github.com/RyanLisse/go-crypto-bot-clean/backend/pkg/platform/mexc/websocket/proto"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"golang.org/x/time/rate"
	"google.golang.org/protobuf/proto"
)

const (
	// WebSocket endpoints
	defaultMexcWSProtobufURL = "wss://wbs.mexc.com/ws" // Updated URL for Protobuf-based WebSocket

	// Reconnection parameters
	// Removed duplicate constants to avoid redeclaration conflicts with client.go
	// reconnectDelay    = 5 * time.Second
	// pingInterval      = 20 * time.Second
	// maxReconnectTries = 5

	// Message types
	// msgTypePing        = "ping"
	// msgTypePong        = "pong"
	// msgTypeSubscribe   = "sub"
	// msgTypeUnsubscribe = "unsub"

	// Channel types
	channelNewListings  = "spot@public.newlistings.v3.api"
	channelSymbolStatus = "spot@public.symbolstatus.v3.api"
)

// MessageHandler is a function type for handling WebSocket messages
type MessageHandler func(message *mexcproto.MexcMessage) error

// ProtobufClient represents a WebSocket client for the MEXC exchange using Protocol Buffers
type ProtobufClient struct {
	conn             *websocket.Conn
	url              string
	subscriptions    map[string]bool
	handlers         map[string][]MessageHandler
	mu               sync.RWMutex
	isConnected      bool
	reconnectTries   int
	ctx              context.Context
	cancel           context.CancelFunc
	reconnectHandler func() error
	rateLimiter      *rate.Limiter
	logger           *zerolog.Logger
}

// IsConnected returns whether the ProtobufClient is connected
func (c *ProtobufClient) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.isConnected
}

// NewProtobufClient creates a new WebSocket client with Protocol Buffer support
func NewProtobufClient(ctx context.Context, logger *zerolog.Logger) *ProtobufClient {
	// Create a rate limiter with MEXC's WebSocket API limits (10 requests per second)
	limiter := rate.NewLimiter(rate.Limit(10), 20) // 10 requests/sec, burst 20

	// Allow override of WebSocket URL via environment variable
	wsURL := os.Getenv("MEXC_WS_PROTOBUF_URL")
	if wsURL == "" {
		wsURL = defaultMexcWSProtobufURL
	}

	ctx, cancel := context.WithCancel(ctx)
	return &ProtobufClient{
		url:           wsURL,
		subscriptions: make(map[string]bool),
		handlers:      make(map[string][]MessageHandler),
		ctx:           ctx,
		cancel:        cancel,
		rateLimiter:   limiter,
		logger:        logger,
	}
}

// Connect establishes a WebSocket connection to MEXC
func (c *ProtobufClient) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isConnected {
		return nil
	}

	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(c.url, nil)
	if err != nil {
		c.reconnectTries++
		return fmt.Errorf("failed to connect to WebSocket: %w", err)
	}

	c.conn = conn
	c.isConnected = true
	c.reconnectTries = 0

	// Start handling messages in a separate goroutine
	go c.handleMessages()

	// Start ping/pong heartbeat to keep connection alive
	go c.keepAlive()

	// Resubscribe to all previous subscriptions after reconnect
	if len(c.subscriptions) > 0 {
		for channel := range c.subscriptions {
			if err := c.sendSubscribeRequest(channel); err != nil {
				c.logger.Error().Err(err).Str("channel", channel).Msg("Failed to resubscribe to channel")
			}
		}
	}

	c.logger.Info().Str("url", c.url).Msg("Connected to MEXC WebSocket")
	return nil
}

// Disconnect closes the WebSocket connection
func (c *ProtobufClient) Disconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.isConnected {
		return nil
	}

	// Cancel context to stop all goroutines
	c.cancel()

	// Close WebSocket connection
	if c.conn != nil {
		err := c.conn.Close()
		c.conn = nil
		c.isConnected = false
		if err != nil {
			return fmt.Errorf("failed to close WebSocket connection: %w", err)
		}
	}

	c.logger.Info().Msg("Disconnected from MEXC WebSocket")
	return nil
}

// SubscribeToNewListings subscribes to new coin listing announcements
func (c *ProtobufClient) SubscribeToNewListings() error {
	return c.subscribe(channelNewListings)
}

// SubscribeToSymbolStatus subscribes to symbol status updates
func (c *ProtobufClient) SubscribeToSymbolStatus() error {
	return c.subscribe(channelSymbolStatus)
}

// RegisterNewListingHandler registers a handler for new listing events
func (c *ProtobufClient) RegisterNewListingHandler(handler MessageHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handlers[channelNewListings] = append(c.handlers[channelNewListings], handler)
}

// RegisterSymbolStatusHandler registers a handler for symbol status updates
func (c *ProtobufClient) RegisterSymbolStatusHandler(handler MessageHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handlers[channelSymbolStatus] = append(c.handlers[channelSymbolStatus], handler)
}

// internal subscribe method to handle subscriptions
func (c *ProtobufClient) subscribe(channel string) error {
	if err := c.ensureConnected(); err != nil {
		return err
	}

	// Add to subscriptions map
	c.mu.Lock()
	c.subscriptions[channel] = true
	c.mu.Unlock()

	return c.sendSubscribeRequest(channel)
}

// sendSubscribeRequest sends a subscription request to the WebSocket
func (c *ProtobufClient) sendSubscribeRequest(channel string) error {
	// Create subscription message
	subMsg := map[string]interface{}{
		"method": msgTypeSubscribe,
		"params": []string{channel},
		"id":     time.Now().UnixNano(),
	}

	return c.sendMessage(subMsg)
}

// handleMessages processes incoming WebSocket messages
func (c *ProtobufClient) handleMessages() {
	defer func() {
		if r := recover(); r != nil {
			c.logger.Error().Interface("panic", r).Msg("Recovered from panic in handleMessages")
			go c.reconnect()
		}
	}()

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			c.mu.RLock()
			conn := c.conn
			isConnected := c.isConnected
			c.mu.RUnlock()

			if !isConnected || conn == nil {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			// Read message
			_, message, err := conn.ReadMessage()
			if err != nil {
				c.logger.Error().Err(err).Msg("Error reading WebSocket message")
				go c.reconnect()
				return
			}

			// Process the message
			if err := c.processMessage(message); err != nil {
				c.logger.Error().Err(err).Msg("Error processing WebSocket message")
			}
		}
	}
}

// processMessage processes a raw WebSocket message
func (c *ProtobufClient) processMessage(data []byte) error {
	// First try to parse as Protocol Buffer
	var pbMsg mexcproto.MexcMessage
	if err := proto.Unmarshal(data, &pbMsg); err == nil {
		return c.handleProtobufMessage(&pbMsg)
	}

	// Fallback to JSON parsing for control messages
	var jsonMsg map[string]interface{}
	if err := json.Unmarshal(data, &jsonMsg); err != nil {
		return fmt.Errorf("failed to parse message as JSON: %w", err)
	}

	// Handle JSON message (usually control messages like ping/pong)
	return c.handleJSONMessage(jsonMsg)
}

// handleProtobufMessage processes a Protocol Buffer message
func (c *ProtobufClient) handleProtobufMessage(msg *mexcproto.MexcMessage) error {
	c.logger.Debug().
		Str("channel", msg.Channel).
		Str("symbol", msg.Symbol).
		Str("data_type", msg.DataType).
		Int64("timestamp", msg.Timestamp).
		Msg("Received Protocol Buffer message")

	// Handle different message types
	switch msg.DataType {
	case "ping":
		return c.handlePing(msg.Timestamp)
	case "pong":
		// Nothing to do for pong responses
		return nil
	case "error":
		if msg.GetErrorResponse() != nil {
			c.logger.Error().
				Int32("code", msg.GetErrorResponse().Code).
				Str("message", msg.GetErrorResponse().Message).
				Msg("Received error message")
		}
		return nil
	default:
		// Dispatch message to registered handlers
		return c.dispatchMessage(msg)
	}
}

// handleJSONMessage processes a JSON message (usually control messages)
func (c *ProtobufClient) handleJSONMessage(msg map[string]interface{}) error {
	// Check for ping message
	if _, ok := msg["ping"]; ok {
		// Extract timestamp
		var timestamp int64
		if ts, ok := msg["ping"].(float64); ok {
			timestamp = int64(ts)
		} else {
			timestamp = time.Now().UnixMilli()
		}
		return c.handlePing(timestamp)
	}

	// Check for subscription response
	if method, ok := msg["method"].(string); ok && method == "sub" {
		result, _ := msg["result"].(bool)
		id, _ := msg["id"].(float64)

		if !result {
			c.logger.Error().
				Float64("id", id).
				Interface("msg", msg).
				Msg("Subscription failed")
		} else {
			c.logger.Info().
				Float64("id", id).
				Msg("Subscription successful")
		}
		return nil
	}

	// Log other messages for debugging
	c.logger.Debug().Interface("message", msg).Msg("Received JSON message")
	return nil
}

// handlePing responds to ping messages with pong
func (c *ProtobufClient) handlePing(timestamp int64) error {
	pongMsg := map[string]interface{}{
		"pong": timestamp,
	}
	return c.sendMessage(pongMsg)
}

// dispatchMessage dispatches a message to registered handlers
func (c *ProtobufClient) dispatchMessage(msg *mexcproto.MexcMessage) error {
	c.mu.RLock()
	handlers, ok := c.handlers[msg.Channel]
	c.mu.RUnlock()

	if !ok || len(handlers) == 0 {
		// No handlers registered for this channel
		return nil
	}

	// Call all registered handlers
	for _, handler := range handlers {
		if err := handler(msg); err != nil {
			c.logger.Error().
				Err(err).
				Str("channel", msg.Channel).
				Msg("Handler error")
		}
	}

	return nil
}

// keepAlive sends periodic ping messages to keep the connection alive
func (c *ProtobufClient) keepAlive() {
	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			if err := c.sendPing(); err != nil {
				c.logger.Error().Err(err).Msg("Failed to send ping")
				go c.reconnect()
				return
			}
		}
	}
}

// sendPing sends a ping message
func (c *ProtobufClient) sendPing() error {
	pingMsg := map[string]interface{}{
		"ping": time.Now().UnixMilli(),
	}
	return c.sendMessage(pingMsg)
}

// reconnect attempts to reconnect to the WebSocket
func (c *ProtobufClient) reconnect() {
	c.mu.Lock()
	if !c.isConnected {
		c.mu.Unlock()
		return
	}

	// Close existing connection
	if c.conn != nil {
		_ = c.conn.Close()
		c.conn = nil
	}
	c.isConnected = false
	reconnectTries := c.reconnectTries
	c.mu.Unlock()

	// Check if we've exceeded max reconnect tries
	if reconnectTries >= maxReconnectTries {
		c.logger.Error().Int("tries", reconnectTries).Msg("Exceeded maximum reconnection attempts")
		return
	}

	// Exponential backoff
	delay := reconnectDelay * time.Duration(1<<uint(reconnectTries))
	if delay > 1*time.Minute {
		delay = 1 * time.Minute
	}

	c.logger.Info().Dur("delay", delay).Msg("Reconnecting to WebSocket")
	time.Sleep(delay)

	// Attempt to reconnect
	if err := c.Connect(); err != nil {
		c.logger.Error().Err(err).Msg("Failed to reconnect")
	}

	// Call reconnect handler if set
	if c.reconnectHandler != nil {
		if err := c.reconnectHandler(); err != nil {
			c.logger.Error().Err(err).Msg("Error in reconnect handler")
		}
	}
}

// ensureConnected ensures the client is connected
func (c *ProtobufClient) ensureConnected() error {
	c.mu.RLock()
	isConnected := c.isConnected
	c.mu.RUnlock()

	if !isConnected {
		return c.Connect()
	}
	return nil
}

// sendMessage sends a message to the WebSocket connection
func (c *ProtobufClient) sendMessage(msg interface{}) error {
	// Check if we have enough tokens in the rate limiter
	_ = c.rateLimiter.Wait(context.Background())

	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.isConnected || c.conn == nil {
		return errors.New("not connected to WebSocket")
	}

	// Send the message
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

// ConvertToNewCoin converts a Protocol Buffer NewListing to a domain model NewCoin
func ConvertToNewCoin(listing *mexcproto.NewListing) *model.NewCoin {
	// Commented out status mapping due to undefined model status constants
	// var status model.Status
	// switch listing.Status {
	// case "PENDING":
	// 	status = model.StatusPending
	// case "PRE_TRADING":
	// 	status = model.StatusPreTrading
	// case "TRADING":
	// 	status = model.StatusTrading
	// case "POST_TRADING":
	// 	status = model.StatusPostTrading
	// case "END_OF_DAY":
	// 	status = model.StatusEndOfDay
	// case "HALT":
	// 	status = model.StatusHalt
	// case "AUCTION_MATCH":
	// 	status = model.StatusAuctionMatch
	// case "BREAK":
	// 	status = model.StatusBreak
	// default:
	// 	status = model.StatusUnknown
	// }

	// Convert listing time and trading time
	listingTime := time.Unix(0, listing.ListingTime*int64(time.Millisecond))
	tradingTime := time.Unix(0, listing.TradingTime*int64(time.Millisecond))

	// Create NewCoin object
	coin := &model.NewCoin{
		Symbol:     listing.Symbol,
		BaseAsset:  listing.BaseAsset,
		QuoteAsset: listing.QuoteAsset,
		// Status:          status,
		ExpectedListingTime: listingTime,
		BecameTradableAt:    &tradingTime,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	return coin
}
