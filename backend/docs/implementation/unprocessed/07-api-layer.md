# API Layer Implementation

This document provides an overview of the API layer for the Go crypto trading bot. For detailed implementation guidelines, see the following specific documents:

- [07a-api-middleware.md](07a-api-middleware.md) - Middleware components for authentication, logging, etc.
- [07b-api-handlers.md](07b-api-handlers.md) - Request handlers for REST API endpoints
- [07c-api-websocket.md](07c-api-websocket.md) - WebSocket implementation for real-time updates

## 1. Overview

The API layer serves as the interface for external clients to interact with the trading bot. It's designed to provide:

- RESTful endpoints for CRUD operations and trading actions
- Real-time WebSocket data streaming
- Authentication and security controls
- Structured error handling and logging

## 2. API Structure

The API follows a clean, layered architecture:

```
internal/api/
├── router.go          # Main router setup and configuration
├── middleware/        # API middleware components
│   ├── auth.go        # Authentication middleware
│   ├── cors.go        # CORS handling
│   ├── logging.go     # Request logging
│   └── recovery.go    # Panic recovery
├── handlers/          # Request handlers for each endpoint
│   ├── coin_handler.go
│   ├── decision_handler.go
│   ├── trade_handler.go
│   ├── config_handler.go
│   ├── log_handler.go
│   ├── status_handler.go
│   └── websocket_handler.go
├── middleware/        # Middleware components
│   ├── logger.go
│   ├── cors.go
│   ├── recovery.go
│   └── auth.go
├── dto/               # Data Transfer Objects
│   ├── request/
│   └── response/
└── routes.go          # Route definitions
```

## 3. Router Setup

First, let's set up the router with Gin:

```go
// internal/api/routes.go
package api

import (
    "github.com/gin-gonic/gin"
    
    "github.com/ryanlisse/cryptobot-backend/internal/api/handlers"
    "github.com/ryanlisse/cryptobot-backend/internal/api/middleware"
    "github.com/ryanlisse/cryptobot-backend/internal/core"
    "github.com/ryanlisse/cryptobot-backend/internal/domain/repository"
    "github.com/ryanlisse/cryptobot-backend/internal/platform/mexc"
)

// SetupRoutes configures all the routes for the API
func SetupRoutes(
    router *gin.Engine,
    mexcClient mexc.Client,
    repoFactory *repository.Factory,
    coreFactory *core.Factory,
) {
    // Apply global middleware
    router.Use(middleware.Logger())
    router.Use(middleware.CORS())
    router.Use(middleware.Recovery())
    
    // Create handlers
    coinHandler := handlers.NewCoinHandler(repoFactory.BoughtCoin(), repoFactory.NewCoin(), mexcClient)
    decisionHandler := handlers.NewDecisionHandler(repoFactory.PurchaseDecision())
    tradeHandler := handlers.NewTradeHandler(repoFactory.BoughtCoin(), mexcClient, coreFactory.CreateTradeExecutor())
    configHandler := handlers.NewConfigHandler()
    logHandler := handlers.NewLogHandler(repoFactory.LogEvent())
    statusHandler := handlers.NewStatusHandler(coreFactory)
    positionHandler := handlers.NewPositionHandler(repoFactory.Position(), coreFactory.CreatePositionManager())
    riskHandler := handlers.NewRiskHandler(coreFactory.CreateRiskManager())
    
    // API v1 group
    v1 := router.Group("/api/v1")
    {
        // Health check endpoint
        v1.GET("/health", handlers.HealthCheck)
        
        // Coins endpoints
        coins := v1.Group("/coins")
        {
            coins.GET("/bought", coinHandler.GetBoughtCoins)
            coins.GET("/bought/:id", coinHandler.GetBoughtCoin)
            coins.GET("/new", coinHandler.GetNewCoins)
        }
        
        // Trading endpoints
        trades := v1.Group("/trades")
        {
            trades.POST("/buy", tradeHandler.BuyManually)
            trades.POST("/sell/:id", tradeHandler.SellManually)
            trades.GET("/history", tradeHandler.GetTradeHistory)
        }
        
        // Position management endpoints
        positions := v1.Group("/positions")
        {
            positions.GET("", positionHandler.GetPositions)
            positions.GET("/:id", positionHandler.GetPosition)
            positions.POST("", positionHandler.CreatePosition)
            positions.PUT("/:id/stoploss", positionHandler.UpdateStopLoss)
            positions.PUT("/:id/takeprofit", positionHandler.UpdateTakeProfit)
            positions.POST("/:id/close", positionHandler.ClosePosition)
        }
        
        // Risk management endpoints
        risk := v1.Group("/risk")
        {
            risk.GET("/status", riskHandler.GetRiskStatus)
            risk.GET("/position-size", riskHandler.CalculatePositionSize)
            risk.GET("/exposure", riskHandler.GetExposure)
        }
        
        // Configuration endpoints
        v1.GET("/config", configHandler.GetConfig)
        v1.PUT("/config", configHandler.UpdateConfig)
        
        // Decisions endpoints
        v1.GET("/decisions", decisionHandler.GetDecisions)
        v1.GET("/decisions/stats", decisionHandler.GetDecisionStats)
        
        // Logs endpoints
        v1.GET("/logs", logHandler.GetLogs)
        
        // Status endpoints
        v1.GET("/status", statusHandler.GetStatus)
        v1.POST("/status/start", statusHandler.StartProcesses)
        v1.POST("/status/stop", statusHandler.StopProcesses)
        
        // Market data endpoints
        v1.GET("/market/ticker/:symbol", coinHandler.GetTicker)
    }
    
    // WebSocket endpoint
    router.GET("/ws", handlers.NewWebSocketHandler(mexcClient, repoFactory, coreFactory).Handle)
}
```

## 4. Middleware Components

### 4.1 Logger Middleware

```go
// internal/api/middleware/logger.go
package middleware

import (
    "time"

    "github.com/gin-gonic/gin"
    "github.com/rs/zerolog/log"
)

// Logger logs information about each request
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        raw := c.Request.URL.RawQuery

        c.Next()

        latency := time.Since(start)
        statusCode := c.Writer.Status()

        if raw != "" {
            path = path + "?" + raw
        }

        log.Info().
            Str("method", c.Request.Method).
            Str("path", path).
            Int("status", statusCode).
            Dur("latency", latency).
            Str("client_ip", c.ClientIP()).
            Str("user_agent", c.Request.UserAgent()).
            Int("size", c.Writer.Size()).
            Msg("API request")
    }
}
```

### 4.2 CORS Middleware

```go
// internal/api/middleware/cors.go
package middleware

import "github.com/gin-gonic/gin"

// CORS handles Cross-Origin Resource Sharing
func CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}
```

### 4.3 Recovery Middleware

```go
// internal/api/middleware/recovery.go
package middleware

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/rs/zerolog/log"
)

// Recovery recovers from panics and logs the error
func Recovery() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                log.Error().
                    Interface("error", err).
                    Str("method", c.Request.Method).
                    Str("path", c.Request.URL.Path).
                    Str("client_ip", c.ClientIP()).
                    Msg("Recovered from panic")

                c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
                    "error": "Internal server error",
                })
            }
        }()

        c.Next()
    }
}
```

### 4.4 Authentication Middleware (Placeholder)

```go
// internal/api/middleware/auth.go
package middleware

import (
    "net/http"
    "os"
    "strings"

    "github.com/gin-gonic/gin"
)

// APIKeyAuth authenticates requests using an API key
func APIKeyAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        apiKey := c.GetHeader("X-API-Key")
        validAPIKey := os.Getenv("API_KEY") // Should come from secure configuration

        if apiKey == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "API key required"})
            return
        }
        if apiKey != validAPIKey {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
            return
        }
        c.Next()
    }
}
```

## 5. API Handlers

### 5.1 Health Check Handler

```go
// internal/api/handlers/health.go
package handlers

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
)

// HealthCheck returns the health status of the API
func HealthCheck(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "status": "ok",
        "time": time.Now().Format(time.RFC3339),
    })
}
```

### 5.2 Coin Handler

```go
// internal/api/handlers/coin_handler.go
package handlers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"

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

### 5.3 Position Handler

```go
// internal/api/handlers/position_handler.go
package handlers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"

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

### 5.4 WebSocket Handler

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
                return true // Allow all origins
            },
        },
        clients: make(map[*websocket.Conn]bool),
    }
}

// Handle handles WebSocket connections
func (h *WebSocketHandler) Handle(c *gin.Context) {
    conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        log.Printf("Error upgrading connection: %v", err)
        return
    }
    
    // Add client to list
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
    
    // Start sending updates
    go h.startUpdates(conn)
    
    // Handle incoming messages
    for {
        messageType, message, err := conn.ReadMessage()
        if err != nil {
            log.Printf("Error reading message: %v", err)
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
        log.Printf("Error parsing message: %v", err)
        return
    }
    
    // Handle different message types
    switch request.Type {
    case "subscribe_ticker":
        var payload struct {
            Symbol string `json:"symbol"`
        }
        json.Unmarshal(request.Payload, &payload)
        
        // Subscribe to ticker updates
        // Implementation depends on your WebSocket architecture
    
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
                log.Printf("Error getting bought coins: %v", err)
                continue
            }
            
            // Send update
            if err := conn.WriteJSON(map[string]interface{}{
                "type": "portfolio_update",
                "data": boughtCoins,
            }); err != nil {
                log.Printf("Error sending update: %v", err)
                return
            }
        }
    }
}
```

## 6. Client Communication

### 6.1 Request/Response DTOs

It's best practice to define clear Data Transfer Objects for API requests and responses:

```go
// internal/api/dto/request/trade_request.go
package request

// BuyRequest represents a request to buy a coin
type BuyRequest struct {
    Symbol   string  `json:"symbol" binding:"required"`
    Quantity float64 `json:"quantity" binding:"required,gt=0"`
}

// SellRequest represents a request to sell a coin
type SellRequest struct {
    Quantity *float64 `json:"quantity"` // Optional, if not provided, sell all
}
```

```go
// internal/api/dto/response/trade_response.go
package response

import "time"

// TradeResponse represents a trade response
type TradeResponse struct {
    ID           int64     `json:"id"`
    Symbol       string    `json:"symbol"`
    Side         string    `json:"side"` // "buy" or "sell"
    Price        float64   `json:"price"`
    Quantity     float64   `json:"quantity"`
    ExecutedAt   time.Time `json:"executed_at"`
    TotalCost    float64   `json:"total_cost,omitempty"`    // Only for buys
    ProfitLoss   float64   `json:"profit_loss,omitempty"`   // Only for sells
    ProfitPercent float64  `json:"profit_percent,omitempty"` // Only for sells
}
```

### 6.2 Error Handling

Standardize API error responses:

```go
// internal/api/errors.go
package api

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

// ErrorResponse represents a standardized API error response
type ErrorResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}

// RespondWithError sends a standardized error response
func RespondWithError(c *gin.Context, statusCode int, message string, details string) {
    c.JSON(statusCode, ErrorResponse{
        Code:    statusCode,
        Message: message,
        Details: details,
    })
}

// Common error response helpers
func BadRequest(c *gin.Context, message string, details string) {
    RespondWithError(c, http.StatusBadRequest, message, details)
}

func Unauthorized(c *gin.Context, message string, details string) {
    RespondWithError(c, http.StatusUnauthorized, message, details)
}

func NotFound(c *gin.Context, message string, details string) {
    RespondWithError(c, http.StatusNotFound, message, details)
}

func InternalError(c *gin.Context, message string, details string) {
    RespondWithError(c, http.StatusInternalServerError, message, details)
}
```

## 7. API Documentation

For API documentation, we recommend using Swagger/OpenAPI. Here's how to set it up with Gin and Swaggo:

```go
// cmd/server/main.go
package main

import (
    "github.com/gin-gonic/gin"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
    
    // Import generated docs
    _ "github.com/ryanlisse/cryptobot-backend/docs"
    "github.com/ryanlisse/cryptobot-backend/internal/api"
)

// @title Crypto Trading Bot API
// @version 1.0
// @description API for cryptocurrency trading bot
// @host localhost:8080
// @BasePath /api/v1
func main() {
    router := gin.Default()
    
    // Setup routes
    // ... (initialization code)
    
    // Setup Swagger
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    
    router.Run(":8080")
}
```

## 8. API Security Best Practices

1. **Use HTTPS**: Always use HTTPS in production
2. **Input Validation**: Validate all input using Gin's binding
3. **Rate Limiting**: Implement rate limiting to prevent abuse
4. **Authentication**: Use API keys or JWT tokens for authentication
5. **Sanitize Outputs**: Never leak sensitive or internal information
6. **Error Handling**: Don't expose implementation details in errors
7. **Logging**: Log all access with appropriate detail
8. **CORS Policy**: Restrict allowed origins in production

## 9. API Testing

For proper API testing:

1. **Unit Tests**: Test individual handlers with mock repositories
2. **Integration Tests**: Test the API with real (test) repositories
3. **End-to-End Tests**: Test the entire API with HTTP clients

Example test:

```go
// internal/api/handlers/coin_handler_test.go
package handlers_test

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
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

This API layer implementation provides a robust, secure interface for interacting with the Go crypto trading bot.
