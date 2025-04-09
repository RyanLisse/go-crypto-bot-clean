package websocket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// Client represents a single WebSocket connection.
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan WSMessage
}

// NewClient creates a new client instance.
func NewClient(hub *Hub, conn *websocket.Conn) *Client {
	return &Client{
		hub:  hub,
		conn: conn,
		send: make(chan WSMessage, 256),
	}
}

// ReadPump pumps messages from the WebSocket connection to the hub (if needed).
func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
		// Optionally handle incoming messages here
		var req struct {
			Type    string          `json:"type"`
			Payload json.RawMessage `json:"payload"`
		}
		if err := json.Unmarshal(message, &req); err != nil {
			c.send <- WSMessage{
				Type:      ErrorType,
				Timestamp: time.Now().Unix(),
				Payload:   ErrorPayload{Message: "Invalid message format"},
			}
			continue
		}
		// Handle different message types
		switch req.Type {
		case "subscribe_ticker", "subscribe":
			// Parse subscription payload
			var subPayload struct {
				Channel string   `json:"channel"`
				Symbols []string `json:"symbols,omitempty"`
			}
			if err := json.Unmarshal(req.Payload, &subPayload); err != nil {
				c.send <- WSMessage{
					Type:      ErrorType,
					Timestamp: time.Now().Unix(),
					Payload:   ErrorPayload{Message: "Invalid subscription format"},
				}
				continue
			}

			// Add symbols to market data service if provided
			if subPayload.Channel == "market_data" && len(subPayload.Symbols) > 0 {
				// This would require access to the market data service
				// For now, just acknowledge the subscription
			}

			c.send <- WSMessage{
				Type:      SubscriptionSuccessType,
				Timestamp: time.Now().Unix(),
				Payload: SubscriptionSuccessPayload{
					Message: "Subscribed to " + subPayload.Channel,
					Channel: subPayload.Channel,
				},
			}

		case "ping":
			// Respond with pong
			c.send <- WSMessage{
				Type:      PongType,
				Timestamp: time.Now().Unix(),
				Payload:   PongPayload{Timestamp: time.Now().Unix()},
			}

		default:
			// Unknown message type
			c.send <- WSMessage{
				Type:      ErrorType,
				Timestamp: time.Now().Unix(),
				Payload:   ErrorPayload{Message: "Unknown message type: " + req.Type},
			}
		}
	}
}

// WritePump pumps messages from the hub to the WebSocket connection.
func (c *Client) WritePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// Hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteJSON(message); err != nil {
				return
			}
		case <-ticker.C:
			// Send ping
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
