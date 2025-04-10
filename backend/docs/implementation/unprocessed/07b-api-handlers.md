# API Handlers Implementation

This document covers the implementation of API handlers for the Go crypto trading bot. These handlers process incoming HTTP requests, interact with core domain services, and return appropriate responses.

## 1. Overview

API handlers act as the interface between external clients and the application's core business logic. They are responsible for:

- Parsing and validating request data
- Calling the appropriate domain services
- Formatting responses
- Handling errors

## 2. Handler Structure

In our hexagonal architecture, handlers are organized within the application layer:

```
internal/api/handlers/
├── health.go            # Health check endpoint
├── coin_handler.go      # Coin-related endpoints
├── trade_handler.go     # Trading endpoints
├── position_handler.go  # Position management endpoints
├── risk_handler.go      # Risk management endpoints
├── config_handler.go    # Configuration endpoints
├── log_handler.go       # Log access endpoints
├── status_handler.go    # System status endpoints
└── websocket_handler.go # WebSocket connections
```

## 3. Core Handler Patterns

All handlers follow a consistent pattern:

```go
// Handler struct with dependencies
type SomeHandler struct {
    // Dependencies injected via constructor
    someRepository repository.SomeRepository
    someService    service.SomeService
}

// Constructor function
func NewSomeHandler(
    someRepository repository.SomeRepository,
    someService service.SomeService,
) *SomeHandler {
    return &SomeHandler{
        someRepository: someRepository,
        someService:    someService,
    }
}

// Handler methods for each endpoint
func (h *SomeHandler) HandleEndpoint(c *gin.Context) {
    // 1. Parse request parameters/body
    // 2. Validate input
    // 3. Call domain service
    // 4. Handle errors
    // 5. Format and return response
}
```

## 4. Basic Health Check Handler

```go
// internal/api/handlers/health.go
package handlers

import (
    "net/http"
    "time"

    " "
)

// HealthCheck returns the health status of the API
func HealthCheck(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "status": "ok",
        "time": time.Now().Format(time.RFC3339),
    })
}
```

## 5. Coin Handler

```go
// internal/api/handlers/coin_handler.go
package handlers

import (
    "net/http"
    "strconv"

    " "

    "github.com/ryanlisse/cryptobot-backend/internal/domain/models"
    "github.com/ryanlisse/cryptobot-backend/internal/domain/repository"
    "github.com/ryanlisse/cryptobot-backend/internal/platform/mexc"
)

// CoinHandler handles API requests related to coins
type CoinHandler struct {
    boughtCoinRepo repository.BoughtCoinRepository
    newCoinRepo    repository.NewCoinRepository
    mexcClient     mexc.Client
}

// NewCoinHandler creates a new coin handler
func NewCoinHandler(
    boughtCoinRepo repository.BoughtCoinRepository,
    newCoinRepo repository.NewCoinRepository,
    mexcClient mexc.Client,
) *CoinHandler {
    return &CoinHandler{
        boughtCoinRepo: boughtCoinRepo,
        newCoinRepo:    newCoinRepo,
        mexcClient:     mexcClient,
    }
}

// GetBoughtCoins returns all bought coins
func (h *CoinHandler) GetBoughtCoins(c *gin.Context) {
    // Parse query parameters
    includeDeletedStr := c.DefaultQuery("include_deleted", "false")
    includeDeleted, _ := strconv.ParseBool(includeDeletedStr)

    // Get coins from repository
    coins, err := h.boughtCoinRepo.FindAll(c.Request.Context(), includeDeleted)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Enrich with current price if requested
    withPriceStr := c.DefaultQuery("with_price", "false")
    withPrice, _ := strconv.ParseBool(withPriceStr)

    if withPrice {
        // Create enriched response
        type coinResponse struct {
            ID               int64     `json:"id"`
            Symbol           string    `json:"symbol"`
            PurchasePrice    float64   `json:"purchase_price"`
            CurrentPrice     float64   `json:"current_price,omitempty"`
            ProfitLoss       float64   `json:"profit_loss,omitempty"`
            ProfitLossPercent float64  `json:"profit_loss_percent,omitempty"`
        }

        responses := make([]coinResponse, 0, len(coins))

        for _, coin := range coins {
            // Skip price enrichment for deleted coins
            if coin.IsDeleted && includeDeleted {
                responses = append(responses, coinResponse{
                    ID:            coin.ID,
                    Symbol:        coin.Symbol,
                    PurchasePrice: coin.PurchasePrice,
                })
                continue
            }

            // Get current price
            ticker, err := h.mexcClient.GetTicker(c.Request.Context(), coin.Symbol)
            if err != nil {
                // Skip price info on error
                responses = append(responses, coinResponse{
                    ID:            coin.ID,
                    Symbol:        coin.Symbol,
                    PurchasePrice: coin.PurchasePrice,
                })
                continue
            }

            // Calculate profit/loss
            currentPrice, _ := strconv.ParseFloat(ticker.LastPrice, 64)
            profitLoss := (currentPrice - coin.PurchasePrice) * coin.Quantity
            profitLossPercent := (currentPrice - coin.PurchasePrice) / coin.PurchasePrice * 100

            responses = append(responses, coinResponse{
                ID:                coin.ID,
                Symbol:            coin.Symbol,
                PurchasePrice:     coin.PurchasePrice,
                CurrentPrice:      currentPrice,
                ProfitLoss:        profitLoss,
                ProfitLossPercent: profitLossPercent,
            })
        }

        c.JSON(http.StatusOK, responses)
        return
    }

    c.JSON(http.StatusOK, coins)
}

// GetTicker returns the current ticker for a symbol
func (h *CoinHandler) GetTicker(c *gin.Context) {
    symbol := c.Param("symbol")
    if symbol == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Symbol parameter is required"})
        return
    }

    ticker, err := h.mexcClient.GetTicker(c.Request.Context(), symbol)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, ticker)
}
```

## 6. Trade Handler

```go
// internal/api/handlers/trade_handler.go
package handlers

import (
    "net/http"
    "strconv"

    " "

    "github.com/ryanlisse/cryptobot-backend/internal/domain/models"
    "github.com/ryanlisse/cryptobot-backend/internal/domain/repository"
    "github.com/ryanlisse/cryptobot-backend/internal/domain/service"
    "github.com/ryanlisse/cryptobot-backend/internal/platform/mexc"
)

// TradeHandler handles trading-related API requests
type TradeHandler struct {
    boughtCoinRepo repository.BoughtCoinRepository
    mexcClient     mexc.Client
    tradeExecutor  service.TradeExecutor
}

// NewTradeHandler creates a new trade handler
func NewTradeHandler(
    boughtCoinRepo repository.BoughtCoinRepository,
    mexcClient mexc.Client,
    tradeExecutor service.TradeExecutor,
) *TradeHandler {
    return &TradeHandler{
        boughtCoinRepo: boughtCoinRepo,
        mexcClient:     mexcClient,
        tradeExecutor:  tradeExecutor,
    }
}

// BuyManually handles manual buy requests
func (h *TradeHandler) BuyManually(c *gin.Context) {
    var req struct {
        Symbol   string  `json:"symbol" binding:"required"`
        Quantity float64 `json:"quantity" binding:"required,gt=0"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // Execute purchase
    boughtCoin, err := h.tradeExecutor.ExecutePurchase(c.Request.Context(), req.Symbol, req.Quantity)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, boughtCoin)
}

// SellManually handles manual sell requests
func (h *TradeHandler) SellManually(c *gin.Context) {
    // Get coin ID from path
    idStr := c.Param("id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
        return
    }
    
    // Parse request body for optional quantity
    var req struct {
        Quantity *float64 `json:"quantity"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // Get the coin to sell
    coin, err := h.boughtCoinRepo.FindByID(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    if coin == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Coin not found"})
        return
    }
    
    // Determine quantity to sell
    quantity := coin.Quantity
    if req.Quantity != nil {
        quantity = *req.Quantity
        if quantity <= 0 || quantity > coin.Quantity {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quantity"})
            return
        }
    }
    
    // Execute sale
    result, err := h.tradeExecutor.ExecuteSale(c.Request.Context(), coin, quantity)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, result)
}
```

## 7. Position Handler

```go
// internal/api/handlers/position_handler.go
package handlers

import (
    "net/http"
    "strconv"

    " "

    "github.com/ryanlisse/cryptobot-backend/internal/domain/models"
    "github.com/ryanlisse/cryptobot-backend/internal/domain/repository"
    "github.com/ryanlisse/cryptobot-backend/internal/domain/service"
)

// PositionHandler handles API requests related to trading positions
type PositionHandler struct {
    positionRepo    repository.PositionRepository
    positionManager service.PositionManager
}

// NewPositionHandler creates a new position handler
func NewPositionHandler(
    positionRepo repository.PositionRepository,
    positionManager service.PositionManager,
) *PositionHandler {
    return &PositionHandler{
        positionRepo:    positionRepo,
        positionManager: positionManager,
    }
}

// GetPositions returns all positions with optional status filter
func (h *PositionHandler) GetPositions(c *gin.Context) {
    status := c.DefaultQuery("status", "")
    
    positions, err := h.positionRepo.FindAll(c.Request.Context(), status)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, positions)
}

// CreatePosition creates a new trading position
func (h *PositionHandler) CreatePosition(c *gin.Context) {
    var req struct {
        Symbol    string  `json:"symbol" binding:"required"`
        Quantity  float64 `json:"quantity" binding:"required,gt=0"`
        EntryPrice float64 `json:"entry_price" binding:"required,gt=0"`
        StopLoss  *float64 `json:"stop_loss"`
        TakeProfit *float64 `json:"take_profit"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // Create order model
    order := models.Order{
        Symbol:     req.Symbol,
        Quantity:   req.Quantity,
        Price:      req.EntryPrice,
        StopLoss:   req.StopLoss,
        TakeProfit: req.TakeProfit,
    }
    
    // Enter position
    position, err := h.positionManager.EnterPosition(c.Request.Context(), order)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusCreated, position)
}

// UpdateStopLoss updates a position's stop-loss
func (h *PositionHandler) UpdateStopLoss(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
        return
    }
    
    var req struct {
        StopLoss float64 `json:"stop_loss" binding:"required,gt=0"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // Update stop-loss
    err = h.positionManager.UpdateStopLoss(c.Request.Context(), id, req.StopLoss)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.Status(http.StatusNoContent)
}
```

## 8. Risk Handler

```go
// internal/api/handlers/risk_handler.go
package handlers

import (
    "net/http"
    "strconv"

    " "

    "github.com/ryanlisse/cryptobot-backend/internal/domain/service"
)

// RiskHandler handles risk management API requests
type RiskHandler struct {
    riskManager service.RiskManager
}

// NewRiskHandler creates a new risk handler
func NewRiskHandler(riskManager service.RiskManager) *RiskHandler {
    return &RiskHandler{
        riskManager: riskManager,
    }
}

// GetRiskStatus returns the current risk status
func (h *RiskHandler) GetRiskStatus(c *gin.Context) {
    status, err := h.riskManager.GetRiskStatus(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, status)
}

// CalculatePositionSize calculates the recommended position size
func (h *RiskHandler) CalculatePositionSize(c *gin.Context) {
    symbol := c.Query("symbol")
    if symbol == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Symbol parameter is required"})
        return
    }
    
    accountBalanceStr := c.DefaultQuery("account_balance", "0")
    accountBalance, err := strconv.ParseFloat(accountBalanceStr, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account balance"})
        return
    }
    
    // Calculate position size
    size, err := h.riskManager.CalculatePositionSize(c.Request.Context(), symbol, accountBalance)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "symbol": symbol,
        "position_size": size,
    })
}
```

## 9. WebSocket Handler

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

    " "
    "github.com/gorilla/websocket"

    "github.com/ryanlisse/cryptobot-backend/internal/core"
    "github.com/ryanlisse/cryptobot-backend/internal/domain/repository"
    "github.com/ryanlisse/cryptobot-backend/internal/platform/mexc"
)

// WebSocketHandler handles WebSocket connections
type WebSocketHandler struct {
    mexcClient  mexc.Client
    repoFactory *repository.Factory
    coreFactory *core.Factory
    upgrader    websocket.Upgrader
    clients     map[*websocket.Conn]bool
    mutex       sync.Mutex
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(
    mexcClient mexc.Client,
    repoFactory *repository.Factory,
    coreFactory *core.Factory,
) *WebSocketHandler {
    return &WebSocketHandler{
        mexcClient:  mexcClient,
        repoFactory: repoFactory,
        coreFactory: coreFactory,
        upgrader: websocket.Upgrader{
            ReadBufferSize:  1024,
            WriteBufferSize: 1024,
            CheckOrigin: func(r *http.Request) bool {
                return true // Allow all origins in development
            },
        },
        clients: make(map[*websocket.Conn]bool),
    }
}

// Handle handles WebSocket connections
func (h *WebSocketHandler) Handle(c *gin.Context) {
    // Upgrade HTTP connection to WebSocket
    conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        log.Printf("Error upgrading connection: %v", err)
        return
    }
    
    // Register client
    h.mutex.Lock()
    h.clients[conn] = true
    h.mutex.Unlock()
    
    // Clean up when connection closes
    defer func() {
        h.mutex.Lock()
        delete(h.clients, conn)
        h.mutex.Unlock()
        conn.Close()
    }()
    
    // Start sending regular updates
    go h.startUpdates(conn)
    
    // Process incoming messages
    for {
        messageType, message, err := conn.ReadMessage()
        if err != nil {
            break
        }
        
        if messageType == websocket.TextMessage {
            h.handleMessage(conn, message)
        }
    }
}

// handleMessage processes incoming WebSocket messages
func (h *WebSocketHandler) handleMessage(conn *websocket.Conn, message []byte) {
    var request struct {
        Type    string          `json:"type"`
        Payload json.RawMessage `json:"payload"`
    }
    
    if err := json.Unmarshal(message, &request); err != nil {
        return
    }
    
    switch request.Type {
    case "subscribe_ticker":
        var payload struct {
            Symbol string `json:"symbol"`
        }
        json.Unmarshal(request.Payload, &payload)
        
        // Handle ticker subscription
    
    case "ping":
        conn.WriteJSON(map[string]string{"type": "pong"})
    }
}

// startUpdates sends periodic updates to the client
func (h *WebSocketHandler) startUpdates(conn *websocket.Conn) {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            // Get portfolio status
            boughtCoins, err := h.repoFactory.BoughtCoin().FindAll(context.Background(), false)
            if err != nil {
                continue
            }
            
            // Send update
            if err := conn.WriteJSON(map[string]interface{}{
                "type": "portfolio_update",
                "data": boughtCoins,
            }); err != nil {
                return
            }
        }
    }
}
```

## 10. Testing API Handlers

```go
// internal/api/handlers/coin_handler_test.go
package handlers_test

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    " "
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"

    "github.com/ryanlisse/cryptobot-backend/internal/api/handlers"
    "github.com/ryanlisse/cryptobot-backend/internal/domain/models"
    "github.com/ryanlisse/cryptobot-backend/internal/mocks"
)

func TestGetBoughtCoins(t *testing.T) {
    // Setup
    gin.SetMode(gin.TestMode)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    
    // Create mocks
    mockBoughtCoinRepo := new(mocks.BoughtCoinRepository)
    mockNewCoinRepo := new(mocks.NewCoinRepository)
    mockMexcClient := new(mocks.MexcClient)
    
    // Set expectations
    mockBoughtCoins := []models.BoughtCoin{
        {ID: 1, Symbol: "BTCUSDT", PurchasePrice: 50000},
        {ID: 2, Symbol: "ETHUSDT", PurchasePrice: 3000},
    }
    mockBoughtCoinRepo.On("FindAll", mock.Anything, false).Return(mockBoughtCoins, nil)
    
    // Create handler
    handler := handlers.NewCoinHandler(mockBoughtCoinRepo, mockNewCoinRepo, mockMexcClient)
    
    // Setup request
    c.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/coins/bought", nil)
    c.Request.URL.RawQuery = "include_deleted=false"
    
    // Execute
    handler.GetBoughtCoins(c)
    
    // Assert
    assert.Equal(t, http.StatusOK, w.Code)
    
    var response []models.BoughtCoin
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Len(t, response, 2)
    assert.Equal(t, "BTCUSDT", response[0].Symbol)
    
    // Verify expectations
    mockBoughtCoinRepo.AssertExpectations(t)
}
```

These API handler implementations provide a robust interface for interacting with the Go crypto trading bot's core business logic while following clean architecture principles.
