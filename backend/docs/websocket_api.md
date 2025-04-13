# WebSocket API Documentation

This document describes the WebSocket API for real-time market data.

## Connection

Connect to the WebSocket API at:

```
ws://<host>/api/v1/ws/market
```

Upon successful connection, you will receive a welcome message:

```json
{
  "type": "info",
  "channel": "system",
  "data": "Connected to market data WebSocket"
}
```

## Rate Limiting

The WebSocket API is rate-limited to 10 messages per second with a burst capacity of 30 messages. If you exceed this limit, you will receive an error message:

```json
{
  "type": "error",
  "channel": "system",
  "data": "Rate limit exceeded. Please slow down your requests."
}
```

## Subscription

To receive real-time data, you need to subscribe to specific channels. The available channels are:

- `tickers`: Real-time ticker updates
- `candles`: Real-time candle (k-line) updates

### Subscribe to Tickers

To subscribe to ticker updates for specific symbols:

```json
{
  "action": "subscribe",
  "channel": "tickers",
  "symbols": ["BTCUSDT", "ETHUSDT"]
}
```

If you want to subscribe to all tickers, use an empty symbols array:

```json
{
  "action": "subscribe",
  "channel": "tickers",
  "symbols": []
}
```

### Subscribe to Candles

To subscribe to candle updates for specific symbols and interval:

```json
{
  "action": "subscribe",
  "channel": "candles",
  "symbols": ["BTCUSDT", "ETHUSDT"],
  "interval": "1h"
}
```

The `interval` parameter is required for candle subscriptions. Valid intervals are:

- `1m`: 1 minute
- `5m`: 5 minutes
- `15m`: 15 minutes
- `30m`: 30 minutes
- `1h`: 1 hour
- `4h`: 4 hours
- `1d`: 1 day
- `1w`: 1 week
- `1M`: 1 month

### Unsubscribe

To unsubscribe from a channel:

```json
{
  "action": "unsubscribe",
  "channel": "tickers",
  "symbols": ["BTCUSDT", "ETHUSDT"]
}
```

For candles, include the interval:

```json
{
  "action": "unsubscribe",
  "channel": "candles",
  "symbols": ["BTCUSDT", "ETHUSDT"],
  "interval": "1h"
}
```

## Responses

### Subscription Response

After subscribing, you will receive a confirmation message:

```json
{
  "type": "subscription",
  "channel": "tickers",
  "data": {
    "status": "success",
    "action": "subscribe",
    "symbols": ["BTCUSDT", "ETHUSDT"]
  }
}
```

### Ticker Data

Ticker updates are sent in the following format:

```json
{
  "type": "data",
  "channel": "tickers",
  "symbol": "BTCUSDT",
  "data": {
    "symbol": "BTCUSDT",
    "exchange": "mexc",
    "price": 50000.0,
    "volume": 100.0,
    "high24h": 51000.0,
    "low24h": 49000.0,
    "priceChange": 1000.0,
    "percentChange": 2.0,
    "lastUpdated": "2023-06-01T12:00:00Z"
  }
}
```

### Candle Data

Candle updates are sent in the following format:

```json
{
  "type": "data",
  "channel": "candles",
  "symbol": "BTCUSDT",
  "data": [
    {
      "symbol": "BTCUSDT",
      "exchange": "mexc",
      "interval": "1h",
      "openTime": "2023-06-01T12:00:00Z",
      "closeTime": "2023-06-01T13:00:00Z",
      "open": 50000.0,
      "high": 51000.0,
      "low": 49000.0,
      "close": 50500.0,
      "volume": 100.0,
      "quoteVolume": 5000000.0,
      "tradeCount": 1000,
      "complete": true
    }
  ]
}
```

## Error Handling

If an error occurs, you will receive an error message:

```json
{
  "type": "error",
  "channel": "system",
  "data": "Error message"
}
```

Common errors:

- `Invalid subscription request format`: The subscription request is not properly formatted
- `Invalid action. Must be 'subscribe' or 'unsubscribe'`: The action is not valid
- `Unsupported channel. Supported channels: 'tickers', 'candles'`: The channel is not supported
- `Interval is required for candle subscriptions`: The interval is missing for candle subscriptions
- `Rate limit exceeded. Please slow down your requests.`: You have exceeded the rate limit

## Connection Maintenance

The WebSocket connection has a 60-second inactivity timeout. To keep the connection alive, you should:

1. Respond to ping frames with pong frames (handled automatically by most WebSocket clients)
2. Send a message at least once every 60 seconds if you're not receiving any data

## Example Usage

### JavaScript Example

```javascript
const ws = new WebSocket('ws://localhost:8080/api/v1/ws/market');

ws.onopen = () => {
  console.log('Connected to WebSocket');
  
  // Subscribe to tickers
  ws.send(JSON.stringify({
    action: 'subscribe',
    channel: 'tickers',
    symbols: ['BTCUSDT', 'ETHUSDT']
  }));
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log('Received message:', message);
  
  if (message.type === 'data' && message.channel === 'tickers') {
    // Handle ticker data
    console.log(`${message.symbol} price: ${message.data.price}`);
  }
};

ws.onerror = (error) => {
  console.error('WebSocket error:', error);
};

ws.onclose = () => {
  console.log('WebSocket connection closed');
};
```

### Go Example

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

type SubscriptionRequest struct {
	Action   string   `json:"action"`
	Channel  string   `json:"channel"`
	Symbols  []string `json:"symbols,omitempty"`
	Interval string   `json:"interval,omitempty"`
}

type WebSocketMessage struct {
	Type    string      `json:"type"`
	Channel string      `json:"channel"`
	Symbol  string      `json:"symbol,omitempty"`
	Data    interface{} `json:"data"`
}

func main() {
	// Connect to WebSocket
	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/api/v1/ws/market", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	// Handle interrupt signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Subscribe to tickers
	subscriptionReq := SubscriptionRequest{
		Action:  "subscribe",
		Channel: "tickers",
		Symbols: []string{"BTCUSDT", "ETHUSDT"},
	}
	err = c.WriteJSON(subscriptionReq)
	if err != nil {
		log.Println("write:", err)
		return
	}

	// Process messages
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			var msg WebSocketMessage
			if err := json.Unmarshal(message, &msg); err != nil {
				log.Println("unmarshal:", err)
				continue
			}

			if msg.Type == "data" && msg.Channel == "tickers" {
				fmt.Printf("%s price: %v\n", msg.Symbol, msg.Data.(map[string]interface{})["price"])
			} else {
				fmt.Printf("Received: %s\n", message)
			}
		}
	}()

	// Keep connection alive
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			// Send ping
			if err := c.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			// Cleanly close the connection
			log.Println("interrupt")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
```
