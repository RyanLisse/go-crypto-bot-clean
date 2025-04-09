package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"go-crypto-bot-clean/backend/internal/domain/models"
	"go.uber.org/zap"
)

// SubscribeToAccountUpdates subscribes to account updates
func (c *Client) SubscribeToAccountUpdates(ctx context.Context, callback func(*models.Wallet)) error {
	c.logger.Info("Subscribing to account updates")

	// Check if client is connected
	if !c.IsConnected() {
		c.logger.Warn("WebSocket client is not connected, cannot subscribe to account updates")
		return fmt.Errorf("websocket client is not connected")
	}

	// Check if we have a listen key (required for private streams)
	if c.listenKey == "" {
		c.logger.Warn("Listen key is required for account updates")
		return fmt.Errorf("listen key is required for account updates")
	}

	// Register the callback
	if c.accountHandler != nil {
		c.logger.Debug("Registering account update callback with account handler")
		c.accountHandler.SubscribeToAccountUpdates(ctx, callback)
	} else {
		c.logger.Warn("Account handler is nil, cannot register callback")
		return fmt.Errorf("account handler is nil")
	}

	// Authenticate the WebSocket connection if not already authenticated
	c.connMutex.RLock()
	isAuthenticated := c.authenticated
	c.connMutex.RUnlock()

	if !isAuthenticated {
		c.logger.Debug("Authenticating WebSocket connection")
		if err := c.Authenticate(ctx); err != nil {
			c.logger.Error("Failed to authenticate WebSocket connection", zap.Error(err))
			return fmt.Errorf("failed to authenticate WebSocket connection: %w", err)
		}
	} else {
		c.logger.Debug("WebSocket connection already authenticated")
	}

	// Subscribe to account updates
	subMsg := map[string]interface{}{
		"method": "SUBSCRIPTION",
		"params": []string{"spot@private.account.v3.api"},
		"id":     time.Now().UnixNano(),
	}

	// Apply rate limiting
	if err := c.subRateLimiter.Wait(ctx); err != nil {
		c.logger.Error("Rate limit error when subscribing to account updates", zap.Error(err))
		return fmt.Errorf("rate limit error: %w", err)
	}

	// Send subscription message
	data, err := json.Marshal(subMsg)
	if err != nil {
		c.logger.Error("Error marshaling subscription message", zap.Error(err))
		return fmt.Errorf("error marshaling subscription message: %w", err)
	}

	c.connMutex.Lock()
	defer c.connMutex.Unlock()

	if c.conn == nil {
		c.logger.Error("WebSocket connection is nil")
		return fmt.Errorf("websocket connection is nil")
	}

	c.logger.Debug("Sending account subscription message",
		zap.String("message", string(data)))

	if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		c.logger.Error("Error sending subscription message", zap.Error(err))
		return fmt.Errorf("error sending subscription message: %w", err)
	}

	c.logger.Info("Successfully subscribed to account updates")
	return nil
}

// IsConnected returns whether the client is connected
func (c *Client) IsConnected() bool {
	c.connMutex.RLock()
	defer c.connMutex.RUnlock()
	return c.connected && c.conn != nil
}

// SetListenKey sets the listen key for private streams
func (c *Client) SetListenKey(listenKey string) {
	c.connMutex.Lock()
	defer c.connMutex.Unlock()
	c.logger.Info("Setting listen key for private streams")
	c.listenKey = listenKey
}

// UnsubscribeFromAccountUpdates unsubscribes from account updates
func (c *Client) UnsubscribeFromAccountUpdates(ctx context.Context) error {
	c.logger.Info("Unsubscribing from account updates")

	// Check if client is connected
	if !c.IsConnected() {
		c.logger.Warn("WebSocket client is not connected, cannot unsubscribe from account updates")
		return fmt.Errorf("websocket client is not connected")
	}

	// Unsubscribe from account updates
	unsubMsg := map[string]interface{}{
		"method": "UNSUBSCRIPTION",
		"params": []string{"spot@private.account.v3.api"},
		"id":     time.Now().UnixNano(),
	}

	// Apply rate limiting
	if err := c.subRateLimiter.Wait(ctx); err != nil {
		c.logger.Error("Rate limit error when unsubscribing from account updates", zap.Error(err))
		return fmt.Errorf("rate limit error: %w", err)
	}

	// Send unsubscription message
	data, err := json.Marshal(unsubMsg)
	if err != nil {
		c.logger.Error("Error marshaling unsubscription message", zap.Error(err))
		return fmt.Errorf("error marshaling unsubscription message: %w", err)
	}

	c.connMutex.Lock()
	defer c.connMutex.Unlock()

	if c.conn == nil {
		c.logger.Error("WebSocket connection is nil")
		return fmt.Errorf("websocket connection is nil")
	}

	c.logger.Debug("Sending account unsubscription message",
		zap.String("message", string(data)))

	if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		c.logger.Error("Error sending unsubscription message", zap.Error(err))
		return fmt.Errorf("error sending unsubscription message: %w", err)
	}

	// Unregister callbacks
	if c.accountHandler != nil {
		c.accountHandler.UnsubscribeFromAccountUpdates()
	}

	c.logger.Info("Successfully unsubscribed from account updates")
	return nil
}
