package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/gorilla/websocket"
	"golang.org/x/time/rate"
)

const (
	// WebSocket endpoints
	mexcWSBaseURL = "wss://wss.mexc.com/ws"

	// Reconnection parameters
	reconnectDelay    = 5 * time.Second
	pingInterval      = 20 * time.Second
	maxReconnectTries = 5

	// Message types
	msgTypePing        = "ping"
	msgTypePong        = "pong"
	msgTypeSubscribe   = "sub"
	msgTypeUnsubscribe = "unsub"
)

// Client represents a WebSocket client for the MEXC exchange
type Client struct {
	conn             *websocket.Conn
	url              string
	subscriptions    map[string]bool
	mu               sync.RWMutex
	isConnected      bool
	reconnectTries   int
	ctx              context.Context
	cancel           context.CancelFunc
	messageHandler   func([]byte) error
	reconnectHandler func() error
	rateLimiter      *rate.Limiter
}

// NewClient creates a new WebSocket client
func NewClient(ctx context.Context) *Client {
	// Create a rate limiter with MEXC's WebSocket API limits (10 requests per second)
	limiter := rate.NewLimiter(rate.Limit(10), 20) // 10 requests/sec, burst 20

	ctx, cancel := context.WithCancel(ctx)
	return &Client{
		url:           mexcWSBaseURL,
		subscriptions: make(map[string]bool),

		ctx:         ctx,
		cancel:      cancel,
		rateLimiter: limiter,
	}
}

// Connect establishes a WebSocket connection to MEXC
func (c *Client) Connect() error {
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
				log.Printf("Failed to resubscribe to %s: %v", channel, err)
			}
		}
	}

	return nil
}

// Disconnect closes the WebSocket connection
func (c *Client) Disconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.isConnected || c.conn == nil {
		return nil
	}

	c.cancel()
	c.isConnected = false

	// Close the connection
	err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		return fmt.Errorf("failed to send close message: %w", err)
	}

	return c.conn.Close()
}

// IsConnected returns the connection status
func (c *Client) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.isConnected
}

// SubscribeToTicker subscribes to ticker updates for a symbol
func (c *Client) SubscribeToTicker(symbol string) error {
	channel := fmt.Sprintf("spot@public.ticker.v3.api@%s", symbol)
	return c.subscribe(channel)
}

// SubscribeToKlines subscribes to kline updates for a symbol and interval
func (c *Client) SubscribeToKlines(symbol string, interval model.KlineInterval) error {
	intervalStr := ""
	switch interval {
	case model.KlineInterval1m:
		intervalStr = "Min1"
	case model.KlineInterval5m:
		intervalStr = "Min5"
	case model.KlineInterval15m:
		intervalStr = "Min15"
	case model.KlineInterval30m:
		intervalStr = "Min30"
	case model.KlineInterval1h:
		intervalStr = "Min60"
	case model.KlineInterval4h:
		intervalStr = "Hour4"
	case model.KlineInterval1d:
		intervalStr = "Day1"
	case model.KlineInterval1w:
		intervalStr = "Week1"
	default:
		return fmt.Errorf("unsupported kline interval: %s", interval)
	}

	channel := fmt.Sprintf("spot@public.kline.v3.api@%s@%s", symbol, intervalStr)
	return c.subscribe(channel)
}

// SubscribeToOrderBook subscribes to order book updates for a symbol
func (c *Client) SubscribeToOrderBook(symbol string) error {
	channel := fmt.Sprintf("spot@public.bookTicker.v3.api@%s", symbol)
	return c.subscribe(channel)
}

// internal subscribe method to handle subscriptions
func (c *Client) subscribe(channel string) error {
	if err := c.ensureConnected(); err != nil {
		return err
	}

	// Add to subscriptions map
	c.mu.Lock()
	c.subscriptions[channel] = true
	c.mu.Unlock()

	return c.sendSubscribeRequest(channel)
}

// Unsubscribe removes a subscription
func (c *Client) Unsubscribe(channel string) error {
	if err := c.ensureConnected(); err != nil {
		return err
	}

	// Remove from subscriptions map
	c.mu.Lock()
	delete(c.subscriptions, channel)
	c.mu.Unlock()

	return c.sendUnsubscribeRequest(channel)
}

// ensureConnected checks and establishes connection if needed
func (c *Client) ensureConnected() error {
	c.mu.RLock()
	isConnected := c.isConnected
	c.mu.RUnlock()

	if !isConnected {
		return c.Connect()
	}
	return nil
}

// handleMessages processes incoming WebSocket messages
func (c *Client) handleMessages() {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			if !c.IsConnected() {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			_, message, err := c.conn.ReadMessage()
			if err != nil {
				log.Printf("Error reading message: %v", err)
				c.handleDisconnect()
				return
			}

			// Handle the received message
			if err := c.processMessage(message); err != nil {
				log.Printf("Error processing message: %v", err)
			}
		}
	}
}

// keepAlive sends ping messages periodically to keep the connection alive
func (c *Client) keepAlive() {
	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			if !c.IsConnected() {
				continue
			}

			pingMsg := map[string]any{"method": msgTypePing}
			if err := c.sendMessage(pingMsg); err != nil {
				log.Printf("Failed to send ping: %v", err)
				c.handleDisconnect()
			}
		}
	}
}

// handleDisconnect manages reconnection logic
func (c *Client) handleDisconnect() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Already handling reconnect or intentionally disconnected
	if !c.isConnected || c.ctx.Err() != nil {
		return
	}

	c.isConnected = false
	if c.conn != nil {
		_ = c.conn.Close()
		c.conn = nil
	}

	// Try to reconnect
	go func() {
		for c.reconnectTries < maxReconnectTries {
			time.Sleep(reconnectDelay)

			if err := c.Connect(); err != nil {
				log.Printf("Failed to reconnect (attempt %d/%d): %v",
					c.reconnectTries, maxReconnectTries, err)
			} else {
				log.Printf("Successfully reconnected after %d attempts", c.reconnectTries)
				return
			}
		}
		log.Printf("Failed to reconnect after %d attempts, giving up", maxReconnectTries)
	}()
}

// sendMessage sends a message to the WebSocket connection
func (c *Client) sendMessage(msg any) error {
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

// sendSubscribeRequest sends a subscription request
func (c *Client) sendSubscribeRequest(channel string) error {
	msg := map[string]any{
		"method": msgTypeSubscribe,
		"params": []string{channel},
		"id":     time.Now().UnixNano(),
	}
	return c.sendMessage(msg)
}

// sendUnsubscribeRequest sends an unsubscription request
func (c *Client) sendUnsubscribeRequest(channel string) error {
	msg := map[string]any{
		"method": msgTypeUnsubscribe,
		"params": []string{channel},
		"id":     time.Now().UnixNano(),
	}
	return c.sendMessage(msg)
}

// processMessage handles incoming messages based on their type
func (c *Client) processMessage(data []byte) error {
	// Skip empty messages
	if len(data) == 0 {
		return nil
	}

	// Parse the raw message to determine its type
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("failed to parse message: %w", err)
	}

	// Handle pong response
	if _, ok := rawMsg["pong"]; ok {
		return nil
	}

	// Handle data messages
	if chdata, ok := rawMsg["data"]; ok {
		if symbol, ok := rawMsg["symbol"]; ok {
			var symbolStr string
			if err := json.Unmarshal(symbol, &symbolStr); err != nil {
				return fmt.Errorf("failed to parse symbol: %w", err)
			}

			// Determine the channel type
			if _, ok := rawMsg["ticker"]; ok {
				return c.handleTickerUpdate(symbolStr, chdata)
			} else if _, ok := rawMsg["kline"]; ok {
				return c.handleKlineUpdate(symbolStr, chdata)
			} else if _, ok := rawMsg["bookTicker"]; ok {
				return c.handleOrderBookUpdate(symbolStr, chdata)
			}
		}
	}

	// No specific handler matched
	return nil
}

// handleTickerUpdate processes ticker updates from the WebSocket
func (c *Client) handleTickerUpdate(_ string, data json.RawMessage) error {
	var tickerData struct {
		Symbol             string `json:"symbol"`
		LastPrice          string `json:"lastPrice"`
		PriceChange        string `json:"priceChange"`
		PriceChangePercent string `json:"priceChangePercent"`
		HighPrice          string `json:"highPrice"`
		LowPrice           string `json:"lowPrice"`
		Volume             string `json:"volume"`
		QuoteVolume        string `json:"quoteVolume"`
		OpenPrice          string `json:"openPrice"`
		PrevClosePrice     string `json:"prevClosePrice"`
		BidPrice           string `json:"bidPrice"`
		BidQuantity        string `json:"bidQty"`
		AskPrice           string `json:"askPrice"`
		AskQuantity        string `json:"askQty"`
		TradeCount         int64  `json:"count"`
		Timestamp          int64  `json:"timestamp"`
	}

	if err := json.Unmarshal(data, &tickerData); err != nil {
		return fmt.Errorf("failed to unmarshal ticker data: %w", err)
	}

	// No further processing needed
	return nil
}

// handleKlineUpdate processes kline updates from the WebSocket
func (c *Client) handleKlineUpdate(_ string, data json.RawMessage) error {
	var klineData struct {
		Symbol    string `json:"symbol"`
		Interval  string `json:"interval"`
		OpenTime  int64  `json:"startTime"`
		CloseTime int64  `json:"endTime"`
		Open      string `json:"open"`
		High      string `json:"high"`
		Low       string `json:"low"`
		Close     string `json:"close"`
		Volume    string `json:"volume"`
		Amount    string `json:"amount"`
		TradeNum  int64  `json:"tradeNum"`
		IsClosed  bool   `json:"isClosed"`
	}

	if err := json.Unmarshal(data, &klineData); err != nil {
		return fmt.Errorf("failed to unmarshal kline data: %w", err)
	}

	// No further processing needed
	return nil
}

// handleOrderBookUpdate processes order book updates from the WebSocket
func (c *Client) handleOrderBookUpdate(_ string, data json.RawMessage) error {
	var orderBookData struct {
		Symbol       string     `json:"symbol"`
		LastUpdateID int64      `json:"lastUpdateId"`
		Bids         [][]string `json:"bids"`
		Asks         [][]string `json:"asks"`
	}

	if err := json.Unmarshal(data, &orderBookData); err != nil {
		return fmt.Errorf("failed to unmarshal order book data: %w", err)
	}

	// No further processing needed
	return nil
}
