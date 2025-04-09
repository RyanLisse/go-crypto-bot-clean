# WebSocket Implementation

This document covers the implementation of WebSocket functionality for the Go crypto trading bot, enabling real-time data streaming and interactive features.

## 1. Overview

WebSockets provide bidirectional communication channels between the server and clients, allowing for real-time updates without the overhead of repeated HTTP requests. In our trading bot, WebSockets are used for:

- Real-time price updates
- Trade notifications
- Portfolio status changes
- Live trading signals

## 2. WebSocket Handler Structure

The WebSocket handler is implemented within the API layer:

```
internal/api/handlers/
└── websocket_handler.go  # WebSocket implementation
```

## 3. WebSocket Handler Implementation

```go
// internal/api/handlers/websocket_handler.go
package handlers

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "sync"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"

    "github.com/ryanlisse/cryptobot-backend/internal/domain/models"
    "github.com/ryanlisse/cryptobot-backend/internal/domain/service"
)

// WebSocketHandler handles WebSocket connections
type WebSocketHandler struct {
    tradeService     service.TradeService
    portfolioService service.PortfolioService
    marketService    service.MarketService
    
    // Track active connections
    clients    map[*websocket.Conn]bool
    clientsMux sync.Mutex
    
    // Configure WebSocket
    upgrader websocket.Upgrader
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(
    tradeService service.TradeService,
    portfolioService service.PortfolioService,
    marketService service.MarketService,
) *WebSocketHandler {
    return &WebSocketHandler{
        tradeService:     tradeService,
        portfolioService: portfolioService,
        marketService:    marketService,
        clients:          make(map[*websocket.Conn]bool),
        upgrader: websocket.Upgrader{
            ReadBufferSize:  1024,
            WriteBufferSize: 1024,
            // Allow all origins for development, restrict in production
            CheckOrigin: func(r *http.Request) bool {
                return true
            },
        },
    }
}

// HandleConnection upgrades HTTP connection to WebSocket
func (h *WebSocketHandler) HandleConnection(c *gin.Context) {
    // Upgrade the HTTP connection to a WebSocket connection
    conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        log.Printf("Failed to upgrade connection: %v", err)
        return
    }
    
    // Register new client
    h.clientsMux.Lock()
    h.clients[conn] = true
    h.clientsMux.Unlock()
    
    // Clean up connection when done
    defer func() {
        conn.Close()
        h.clientsMux.Lock()
        delete(h.clients, conn)
        h.clientsMux.Unlock()
    }()
    
    // Start goroutines for handling this connection
    go h.handleMessages(conn)
    go h.startUpdates(conn)
}

// handleMessages processes incoming messages from clients
func (h *WebSocketHandler) handleMessages(conn *websocket.Conn) {
    for {
        // Read message
        _, message, err := conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, 
                websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("WebSocket error: %v", err)
            }
            break
        }
        
        // Process message
        h.handleMessage(conn, message)
    }
}

// handleMessage processes a single message
func (h *WebSocketHandler) handleMessage(conn *websocket.Conn, message []byte) {
    var request struct {
        Type    string          `json:"type"`
        Payload json.RawMessage `json:"payload"`
    }
    
    if err := json.Unmarshal(message, &request); err != nil {
        h.sendError(conn, "Invalid message format")
        return
    }
    
    // Handle different message types
    switch request.Type {
    case "subscribe_ticker":
        h.handleSubscribeTicker(conn, request.Payload)
    case "manual_trade":
        h.handleManualTrade(conn, request.Payload)
    case "update_position":
        h.handleUpdatePosition(conn, request.Payload)
    default:
        h.sendError(conn, "Unknown message type")
    }
}

// handleSubscribeTicker subscribes to ticker updates
func (h *WebSocketHandler) handleSubscribeTicker(conn *websocket.Conn, payload json.RawMessage) {
    var request struct {
        Symbol string `json:"symbol"`
    }
    
    if err := json.Unmarshal(payload, &request); err != nil {
        h.sendError(conn, "Invalid payload for ticker subscription")
        return
    }
    
    // Store subscription info and respond with confirmation
    h.sendJSON(conn, map[string]interface{}{
        "type":    "subscription_success",
        "message": "Subscribed to " + request.Symbol,
    })
}

// startUpdates sends periodic updates to the client
func (h *WebSocketHandler) startUpdates(conn *websocket.Conn) {
    // Create a ticker for regular updates
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            // Get portfolio summary
            portfolio, err := h.portfolioService.GetSummary(context.Background())
            if err != nil {
                log.Printf("Failed to get portfolio: %v", err)
                continue
            }
            
            // Send update
            if err := h.sendJSON(conn, map[string]interface{}{
                "type":     "portfolio_update",
                "payload":  portfolio,
                "timestamp": time.Now().Unix(),
            }); err != nil {
                return // Connection is likely closed
            }
        }
    }
}

// sendJSON sends JSON data to the client
func (h *WebSocketHandler) sendJSON(conn *websocket.Conn, data interface{}) error {
    return conn.WriteJSON(data)
}

// sendError sends an error message to the client
func (h *WebSocketHandler) sendError(conn *websocket.Conn, message string) {
    h.sendJSON(conn, map[string]interface{}{
        "type":  "error",
        "error": message,
    })
}

// BroadcastTrade broadcasts a trade notification to all connected clients
func (h *WebSocketHandler) BroadcastTrade(trade *models.BoughtCoin) {
    h.clientsMux.Lock()
    defer h.clientsMux.Unlock()
    
    message := map[string]interface{}{
        "type": "trade_notification",
        "payload": map[string]interface{}{
            "symbol":        trade.Symbol,
            "price":         trade.PurchasePrice,
            "quantity":      trade.Quantity,
            "purchased_at":  trade.PurchasedAt,
            "type":          "buy",
        },
        "timestamp": time.Now().Unix(),
    }
    
    // Send to all clients
    for client := range h.clients {
        if err := client.WriteJSON(message); err != nil {
            log.Printf("Error broadcasting to client: %v", err)
            client.Close()
            delete(h.clients, client)
        }
    }
}

// BroadcastAlert broadcasts an alert to all connected clients
func (h *WebSocketHandler) BroadcastAlert(alertType, message string) {
    h.clientsMux.Lock()
    defer h.clientsMux.Unlock()
    
    alert := map[string]interface{}{
        "type": "alert",
        "payload": map[string]interface{}{
            "alert_type": alertType,
            "message":    message,
        },
        "timestamp": time.Now().Unix(),
    }
    
    // Send to all clients
    for client := range h.clients {
        if err := client.WriteJSON(alert); err != nil {
            log.Printf("Error broadcasting to client: %v", err)
            client.Close()
            delete(h.clients, client)
        }
    }
}
```

## 4. Client-Side Implementation

Below is an example of how to connect to the WebSocket from a JavaScript client:

```javascript
// client-side implementation
const socket = new WebSocket('ws://localhost:8080/api/v1/ws');

// Connection opened
socket.addEventListener('open', (event) => {
    console.log('Connected to WebSocket server');
    
    // Subscribe to ticker updates
    socket.send(JSON.stringify({
        type: 'subscribe_ticker',
        payload: {
            symbol: 'BTCUSDT'
        }
    }));
});

// Listen for messages
socket.addEventListener('message', (event) => {
    const data = JSON.parse(event.data);
    
    switch(data.type) {
        case 'portfolio_update':
            updatePortfolioUI(data.payload);
            break;
        case 'trade_notification':
            showTradeNotification(data.payload);
            break;
        case 'alert':
            showAlert(data.payload.alert_type, data.payload.message);
            break;
        case 'error':
            console.error('WebSocket error:', data.error);
            break;
    }
});

// Connection closed
socket.addEventListener('close', (event) => {
    console.log('Disconnected from WebSocket server');
});

// Handle errors
socket.addEventListener('error', (event) => {
    console.error('WebSocket error:', event);
});
```

## 5. WebSocket Integration with Core Services

The WebSocket handler integrates with core services to provide real-time updates:

```go
// Example integration with the TradeService
// internal/domain/service/trade_service.go

// ExecutePurchase executes a purchase and broadcasts via WebSocket
func (s *TradeService) ExecutePurchase(ctx context.Context, symbol string, amount float64) (*models.BoughtCoin, error) {
    // Perform purchase logic...
    
    // After successful purchase
    if s.wsHandler != nil {
        s.wsHandler.BroadcastTrade(boughtCoin)
    }
    
    return boughtCoin, nil
}
```

## 6. Integration with Router

Register the WebSocket handler with the router:

```go
// internal/api/router.go
package api

import (
    "github.com/gin-gonic/gin"
    "github.com/ryanlisse/cryptobot-backend/internal/api/handlers"
)

// SetupRouter configures the API routes
func SetupRouter(
    coinHandler *handlers.CoinHandler,
    tradeHandler *handlers.TradeHandler,
    wsHandler *handlers.WebSocketHandler,
    // Other handlers...
) *gin.Engine {
    router := gin.Default()
    
    // API v1 routes
    v1 := router.Group("/api/v1")
    {
        // WebSocket endpoint
        v1.GET("/ws", wsHandler.HandleConnection)
        
        // Other routes...
    }
    
    return router
}
```

## 7. Security Considerations

When implementing WebSockets in production:

1. **Authentication**: Implement token-based authentication for WebSocket connections
2. **Rate Limiting**: Prevent abuse by limiting connection frequency
3. **Origin Checking**: Restrict connections to known origins
4. **Data Validation**: Validate all incoming WebSocket messages
5. **Timeouts**: Implement connection and read timeouts

Example of adding authentication to WebSockets:

```go
// Middleware to authenticate WebSocket connections
func WebSocketAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.Query("token")
        if token == "" {
            token = c.GetHeader("Authorization")
        }
        
        // Validate token
        if !isValidToken(token) {
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }
        
        // Set user identity in context
        userID, _ := getUserIDFromToken(token)
        c.Set("userID", userID)
        
        c.Next()
    }
}

// Apply middleware to WebSocket route
v1.GET("/ws", WebSocketAuthMiddleware(), wsHandler.HandleConnection)
```

## 8. Testing WebSocket Functionality

Example of testing the WebSocket handler:

```go
// internal/api/handlers/websocket_handler_test.go
package handlers_test

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    
    "github.com/ryanlisse/cryptobot-backend/internal/api/handlers"
    "github.com/ryanlisse/cryptobot-backend/internal/domain/service/mocks"
)

func TestWebSocketHandler_HandleConnection(t *testing.T) {
    // Setup mocks
    mockTradeService := new(mocks.TradeService)
    mockPortfolioService := new(mocks.PortfolioService)
    mockMarketService := new(mocks.MarketService)
    
    // Create handler
    wsHandler := handlers.NewWebSocketHandler(
        mockTradeService,
        mockPortfolioService,
        mockMarketService,
    )
    
    // Setup router
    gin.SetMode(gin.TestMode)
    router := gin.New()
    router.GET("/ws", wsHandler.HandleConnection)
    
    // Start test server
    server := httptest.NewServer(router)
    defer server.Close()
    
    // Convert http URL to ws URL
    wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
    
    // Connect to WebSocket
    ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
    assert.NoError(t, err)
    defer ws.Close()
    
    // Test sending a message
    err = ws.WriteJSON(map[string]interface{}{
        "type": "subscribe_ticker",
        "payload": map[string]string{
            "symbol": "BTCUSDT",
        },
    })
    assert.NoError(t, err)
    
    // Read response
    var response map[string]interface{}
    err = ws.ReadJSON(&response)
    assert.NoError(t, err)
    
    // Verify response
    assert.Equal(t, "subscription_success", response["type"])
}
```

For more details on API layer implementation, see [07-api-layer.md](07-api-layer.md) and [07b-api-handlers.md](07b-api-handlers.md).
