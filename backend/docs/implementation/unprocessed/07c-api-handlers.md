# API Handlers Implementation

This document details the implementation of API handlers in the Go cryptocurrency trading bot. Each handler is responsible for processing specific types of client requests and returning appropriate responses.

## Table of Contents

1. [Overview](#overview)
2. [Handler Structure](#handler-structure)
3. [Health Check Handler](#health-check-handler)
4. [Coin Handler](#coin-handler)
5. [Trade Handler](#trade-handler)
6. [Status Handler](#status-handler)
7. [Config Handler](#config-handler)
8. [Decision Handler](#decision-handler)
9. [Log Handler](#log-handler)
10. [WebSocket Handler](#websocket-handler)

## Overview

The API handlers layer serves as the interface between the client applications and the core business logic. It handles:

- Request validation and parsing
- Error handling and appropriate HTTP responses
- Data serialization/deserialization
- Authentication and authorization (via middleware)
- Real-time updates via WebSockets

All handlers follow a consistent pattern of dependency injection through constructors, clear separation of concerns, proper error handling, and consistent response formats.

## Handler Structure

Each handler is implemented as a struct with the dependencies it needs, following the dependency injection pattern:

```go
// Generic handler structure pattern
type SomeHandler struct {
    // Dependencies injected via constructor
    dependency1 SomeDependency
    dependency2 AnotherDependency
}

// Constructor for dependency injection
func NewSomeHandler(dep1 SomeDependency, dep2 AnotherDependency) *SomeHandler {
    return &SomeHandler{
        dependency1: dep1,
        dependency2: dep2,
    }
}

// Handler methods
func (h *SomeHandler) HandleSomeRequest(c *gin.Context) {
    // Request handling logic
}
```

## Health Check Handler

A simple handler that returns the API's health status, useful for monitoring and load balancers.

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

## Coin Handler

Manages coin-related operations, including listing bought coins and new coins.

```go
// internal/api/handlers/coin_handler.go
package handlers

import (
    "net/http"
    "strconv"
    "context"

    " "

    "github.com/ryanlisse/cryptobot/internal/domain/models"
    "github.com/ryanlisse/cryptobot/internal/domain/repository"
    "github.com/ryanlisse/cryptobot/internal/mexc"
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

// Key methods:
// - GetBoughtCoins: Lists all bought coins with optional price enrichment
// - GetBoughtCoin: Retrieves a specific bought coin by ID
// - GetNewCoins: Lists new coins with optional status filtering
// - GetTicker: Returns current market data for a symbol
```

### Key Implementation Details

The `GetBoughtCoins` method supports several query parameters:
- `include_deleted`: Boolean to include or exclude sold coins
- `with_price`: Boolean to enrich the response with current prices and profit/loss calculations

When `with_price=true`, the handler:
1. Fetches each coin's current price from the exchange
2. Calculates profit/loss amounts and percentages
3. Returns an enriched response with this additional information

## Trade Handler

Manages trading operations, including manual buying and selling.

```go
// internal/api/handlers/trade_handler.go
package handlers

import (
    "context"
    "net/http"
    "strconv"
    "time"

    " "

    "github.com/ryanlisse/cryptobot/internal/domain/models"
    "github.com/ryanlisse/cryptobot/internal/domain/repository"
    "github.com/ryanlisse/cryptobot/internal/mexc"
)

// TradeExecutor is an interface for executing trades
type TradeExecutor interface {
    ExecuteBuy(ctx context.Context, symbol string, price, quantity float64) error
    ExecuteSell(ctx context.Context, coin *models.BoughtCoin, price, quantity float64, reason string) error
}

// TradeHandler handles API requests related to trading
type TradeHandler struct {
    boughtCoinRepo repository.BoughtCoinRepository
    mexcClient     mexc.Client
    executor       TradeExecutor
}

// Key methods:
// - BuyManually: Handles manual buy requests 
// - SellManually: Handles manual sell requests for a specific coin
```

### Key Implementation Details

1. The handler delegates actual trading operations to the `TradeExecutor` interface
2. For buying, it accepts a JSON payload with symbol, quantity, and price
3. For selling, it:
   - Retrieves the coin from the repository
   - Gets the current market price
   - Executes the sell through the trade executor
   - Calculates and returns profit/loss information

## Status Handler

Manages the status of bot components and allows starting/stopping processes.

```go
// internal/api/handlers/status_handler.go
package handlers

// WatcherStatus represents components that can report their running state
type WatcherStatus interface {
    IsRunning() bool
    Start(ctx context.Context) error
    Stop()
}

// StatusProvider is an interface for retrieving component status
type StatusProvider interface {
    GetNewCoinWatcher() WatcherStatus
    GetPositionMonitor() WatcherStatus
}

// StatusHandler handles API requests related to bot status
type StatusHandler struct {
    provider StatusProvider
}

// Key methods:
// - GetStatus: Returns the current status of the bot components
// - StartProcesses: Starts the bot processes
// - StopProcesses: Stops the bot processes
```

### Key Implementation Details

This handler uses the `StatusProvider` interface to:
1. Check the current status of each component
2. Start or stop components as requested
3. Return the updated status after operations

## Config Handler

Manages application configuration through the API.

```go
// internal/api/handlers/config_handler.go
package handlers

// ConfigHandler handles API requests related to configuration
type ConfigHandler struct {
    config *config.Config
}

// Key methods:
// - GetConfig: Returns the current configuration (excluding sensitive info)
// - UpdateConfig: Updates the configuration based on request
```

### Key Implementation Details

The configuration handler:
1. Exposes only non-sensitive configuration values
2. Allows partial updates (only specified fields are changed)
3. Persists configuration changes to storage
4. Implements validation to ensure configuration remains valid

## Decision Handler

Provides access to purchase decisions made by the bot.

```go
// internal/api/handlers/decision_handler.go
package handlers

// DecisionHandler handles API requests related to purchase decisions
type DecisionHandler struct {
    decisionRepo repository.PurchaseDecisionRepository
}

// Key methods:
// - GetDecisions: Returns purchase decisions with pagination
// - GetDecisionStats: Returns statistics about purchase decisions
```

## Log Handler

Provides access to system log events with filtering capabilities.

```go
// internal/api/handlers/log_handler.go
package handlers

// LogHandler handles API requests related to log events
type LogHandler struct {
    logRepo repository.LogEventRepository
}

// Key methods:
// - GetLogs: Returns log events with filtering by level, component, etc.
```

## WebSocket Handler

Implements real-time updates via WebSockets.

```go
// internal/api/handlers/websocket_handler.go
package handlers

// WebSocketHandler handles WebSocket connections
type WebSocketHandler struct {
    mexcClient  mexc.Client
    repoFactory *repository.Factory
    clients     map[*websocket.Conn]bool
    mutex       sync.Mutex
    upgrader    websocket.Upgrader
}

// Key methods:
// - Handle: Handles WebSocket connections and lifecycle
// - Broadcast: Sends a message to all connected clients
```

### Key Implementation Details

The WebSocket handler provides:
1. Connection upgrading from HTTP to WebSocket protocol
2. Client tracking and cleanup for closed connections
3. Initial data payload when clients connect
4. Regular updates of portfolio data
5. Broadcasting capability for system events

The WebSocket functionality enables real-time updates such as:
- Portfolio value changes
- New trade executions
- System status changes
- New coin detections

## Handler Registration

All handlers are registered with the router in the `internal/api/router.go` file:

```go
// RegisterHandlers registers all API handlers with the router
func RegisterHandlers(r *gin.Engine, deps *Dependencies) {
    // Health check
    r.GET("/health", handlers.HealthCheck)
    
    // Coin handlers
    coinHandler := handlers.NewCoinHandler(deps.BoughtCoinRepo, deps.NewCoinRepo, deps.MexcClient)
    r.GET("/coins/bought", coinHandler.GetBoughtCoins)
    r.GET("/coins/bought/:id", coinHandler.GetBoughtCoin)
    r.GET("/coins/new", coinHandler.GetNewCoins)
    r.GET("/ticker/:symbol", coinHandler.GetTicker)
    
    // Other handlers...
    
    // WebSocket
    wsHandler := handlers.NewWebSocketHandler(deps.MexcClient, deps.RepoFactory)
    r.GET("/ws", wsHandler.Handle)
}
```

## Integration with Other Components

The API handlers integrate with other components as follows:

1. **Domain Services**: Handlers delegate business logic to domain services
2. **Repositories**: Data access is done through repository interfaces
3. **External APIs**: Exchange communication uses the MEXC client
4. **Event System**: WebSocket handler broadcasts system events

This separation of concerns ensures that handlers focus on HTTP-specific concerns while delegating business logic to the appropriate domain services.
