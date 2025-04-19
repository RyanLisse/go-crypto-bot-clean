package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

const (
	// WebSocket endpoint for MEXC
	mexcWSEndpoint = "wss://stream.mexc.com/ws"

	// Ping interval
	pingInterval = 30 * time.Second

	// Reconnect delay
	reconnectDelay = 5 * time.Second

	// Message types
	messageTypePing = "ping"
	messageTypeSub  = "sub"
	messageTypeUnsub = "unsub"
)

// Client represents a WebSocket client for MEXC
type Client struct {
	conn      *websocket.Conn
	url       string
	apiKey    string
	apiSecret string
	logger    *zerolog.Logger
	
	// Subscriptions
	subscriptions map[string]chan []byte
	subMutex      sync.RWMutex
	
	// Connection management
	isConnected bool
	connMutex   sync.RWMutex
	ctx         context.Context
	cancel      context.CancelFunc
}

// Message represents a WebSocket message
type Message struct {
	Type    string      `json:"op"`
	Symbol  string      `json:"symbol,omitempty"`
	Channel string      `json:"channel,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// NewClient creates a new WebSocket client
func NewClient(apiKey, apiSecret string, logger *zerolog.Logger) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &Client{
		url:           mexcWSEndpoint,
		apiKey:        apiKey,
		apiSecret:     apiSecret,
		logger:        logger,
		subscriptions: make(map[string]chan []byte),
		ctx:           ctx,
		cancel:        cancel,
	}
}

// Connect establishes a WebSocket connection
func (c *Client) Connect() error {
	c.connMutex.Lock()
	defer c.connMutex.Unlock()
	
	if c.isConnected {
		return nil
	}
	
	u, err := url.Parse(c.url)
	if err != nil {
		return fmt.Errorf("failed to parse URL: %w", err)
	}
	
	c.logger.Info().Str("url", u.String()).Msg("Connecting to MEXC WebSocket")
	
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to connect to WebSocket: %w", err)
	}
	
	c.conn = conn
	c.isConnected = true
	
	// Start message handler
	go c.handleMessages()
	
	// Start ping handler
	go c.pingHandler()
	
	return nil
}

// Disconnect closes the WebSocket connection
func (c *Client) Disconnect() error {
	c.connMutex.Lock()
	defer c.connMutex.Unlock()
	
	if !c.isConnected {
		return nil
	}
	
	c.cancel()
	
	err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		c.logger.Warn().Err(err).Msg("Error sending close message")
	}
	
	err = c.conn.Close()
	if err != nil {
		return fmt.Errorf("failed to close WebSocket connection: %w", err)
	}
	
	c.isConnected = false
	c.conn = nil
	
	return nil
}

// Subscribe subscribes to a channel
func (c *Client) Subscribe(channel, symbol string) (chan []byte, error) {
	if err := c.ensureConnected(); err != nil {
		return nil, err
	}
	
	key := fmt.Sprintf("%s:%s", channel, symbol)
	
	c.subMutex.Lock()
	defer c.subMutex.Unlock()
	
	// Check if already subscribed
	if ch, ok := c.subscriptions[key]; ok {
		return ch, nil
	}
	
	// Create subscription message
	msg := Message{
		Type:    messageTypeSub,
		Channel: channel,
		Symbol:  symbol,
	}
	
	// Send subscription message
	if err := c.sendJSON(msg); err != nil {
		return nil, fmt.Errorf("failed to send subscription message: %w", err)
	}
	
	// Create channel for subscription
	ch := make(chan []byte, 100)
	c.subscriptions[key] = ch
	
	c.logger.Info().Str("channel", channel).Str("symbol", symbol).Msg("Subscribed to channel")
	
	return ch, nil
}

// Unsubscribe unsubscribes from a channel
func (c *Client) Unsubscribe(channel, symbol string) error {
	if err := c.ensureConnected(); err != nil {
		return err
	}
	
	key := fmt.Sprintf("%s:%s", channel, symbol)
	
	c.subMutex.Lock()
	defer c.subMutex.Unlock()
	
	// Check if subscribed
	ch, ok := c.subscriptions[key]
	if !ok {
		return nil
	}
	
	// Create unsubscription message
	msg := Message{
		Type:    messageTypeUnsub,
		Channel: channel,
		Symbol:  symbol,
	}
	
	// Send unsubscription message
	if err := c.sendJSON(msg); err != nil {
		return fmt.Errorf("failed to send unsubscription message: %w", err)
	}
	
	// Close channel and remove subscription
	close(ch)
	delete(c.subscriptions, key)
	
	c.logger.Info().Str("channel", channel).Str("symbol", symbol).Msg("Unsubscribed from channel")
	
	return nil
}

// handleMessages handles incoming WebSocket messages
func (c *Client) handleMessages() {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			if !c.isConnected {
				time.Sleep(reconnectDelay)
				continue
			}
			
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				c.logger.Error().Err(err).Msg("Error reading message")
				c.reconnect()
				continue
			}
			
			// Handle message
			go c.processMessage(message)
		}
	}
}

// processMessage processes an incoming message
func (c *Client) processMessage(message []byte) {
	// Parse message to determine channel and symbol
	var msg struct {
		Channel string `json:"c"`
		Symbol  string `json:"s"`
	}
	
	if err := json.Unmarshal(message, &msg); err != nil {
		c.logger.Error().Err(err).Str("message", string(message)).Msg("Failed to parse message")
		return
	}
	
	// Route message to appropriate subscription
	key := fmt.Sprintf("%s:%s", msg.Channel, msg.Symbol)
	
	c.subMutex.RLock()
	ch, ok := c.subscriptions[key]
	c.subMutex.RUnlock()
	
	if ok {
		select {
		case ch <- message:
			// Message sent to channel
		default:
			c.logger.Warn().Str("channel", msg.Channel).Str("symbol", msg.Symbol).Msg("Channel buffer full, dropping message")
		}
	}
}

// pingHandler sends periodic ping messages to keep the connection alive
func (c *Client) pingHandler() {
	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			if !c.isConnected {
				continue
			}
			
			msg := Message{
				Type: messageTypePing,
				Data: time.Now().UnixNano() / int64(time.Millisecond),
			}
			
			if err := c.sendJSON(msg); err != nil {
				c.logger.Error().Err(err).Msg("Failed to send ping")
				c.reconnect()
			}
		}
	}
}

// reconnect attempts to reconnect to the WebSocket
func (c *Client) reconnect() {
	c.connMutex.Lock()
	defer c.connMutex.Unlock()
	
	if !c.isConnected {
		return
	}
	
	c.logger.Info().Msg("Reconnecting to MEXC WebSocket")
	
	// Close existing connection
	if c.conn != nil {
		_ = c.conn.Close()
	}
	
	c.isConnected = false
	c.conn = nil
	
	// Attempt to reconnect
	go func() {
		for {
			select {
			case <-c.ctx.Done():
				return
			default:
				time.Sleep(reconnectDelay)
				
				if err := c.Connect(); err != nil {
					c.logger.Error().Err(err).Msg("Failed to reconnect")
					continue
				}
				
				// Resubscribe to channels
				c.resubscribe()
				return
			}
		}
	}()
}

// resubscribe resubscribes to all channels
func (c *Client) resubscribe() {
	c.subMutex.RLock()
	defer c.subMutex.RUnlock()
	
	for key, _ := range c.subscriptions {
		var channel, symbol string
		fmt.Sscanf(key, "%s:%s", &channel, &symbol)
		
		msg := Message{
			Type:    messageTypeSub,
			Channel: channel,
			Symbol:  symbol,
		}
		
		if err := c.sendJSON(msg); err != nil {
			c.logger.Error().Err(err).Str("channel", channel).Str("symbol", symbol).Msg("Failed to resubscribe")
		} else {
			c.logger.Info().Str("channel", channel).Str("symbol", symbol).Msg("Resubscribed to channel")
		}
	}
}

// sendJSON sends a JSON message over the WebSocket
func (c *Client) sendJSON(v interface{}) error {
	c.connMutex.RLock()
	defer c.connMutex.RUnlock()
	
	if !c.isConnected {
		return fmt.Errorf("not connected")
	}
	
	return c.conn.WriteJSON(v)
}

// ensureConnected ensures that the client is connected
func (c *Client) ensureConnected() error {
	c.connMutex.RLock()
	isConnected := c.isConnected
	c.connMutex.RUnlock()
	
	if !isConnected {
		return c.Connect()
	}
	
	return nil
}
