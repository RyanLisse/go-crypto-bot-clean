# MEXC Client: Main Implementation

This document explains the implementation of the main MEXC client that combines both REST and WebSocket functionality.

## Main Client Implementation

The main client implements the domain layer's `ExchangeService` interface by composing the REST and WebSocket clients.

```go
// internal/platform/mexc/mexc.go
package mexc

import (
    "context"
    "github.com/ryanlisse/cryptobot/internal/domain/models"
    "github.com/ryanlisse/cryptobot/internal/domain/service"
    "github.com/ryanlisse/cryptobot/internal/platform/mexc/rest"
    "github.com/ryanlisse/cryptobot/internal/platform/mexc/websocket"
    "sync"
    "time"
)

// Client implements the domain.service.ExchangeService interface
type Client struct {
    restClient     *rest.Client
    wsClient       *websocket.Client
    tickerChannels map[string]chan *models.Ticker
    mu             sync.RWMutex
    isConnected    bool
}

// NewClient creates a new MEXC client with both REST and WebSocket capabilities
func NewClient(apiKey, secretKey string, baseURL string) (service.ExchangeService, error) {
    // Create REST client
    restClient, err := rest.NewClient(apiKey, secretKey, baseURL)
    if err != nil {
        return nil, err
    }
    
    // Create WebSocket client
    wsClient, err := websocket.NewClient(apiKey, secretKey)
    if err != nil {
        return nil, err
    }
    
    client := &Client{
        restClient:     restClient,
        wsClient:       wsClient,
        tickerChannels: make(map[string]chan *models.Ticker),
        isConnected:    false,
    }
    
    return client, nil
}

// Connect establishes WebSocket connection and initializes listeners
func (c *Client) Connect(ctx context.Context) error {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    if c.isConnected {
        return nil // Already connected
    }
    
    if err := c.wsClient.Connect(ctx); err != nil {
        return err
    }
    
    c.isConnected = true
    
    // Start listening for WebSocket messages in a separate goroutine
    go c.handleWebSocketMessages(ctx)
    
    return nil
}

// Disconnect closes the WebSocket connection
func (c *Client) Disconnect() error {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    if !c.isConnected {
        return nil // Already disconnected
    }
    
    if err := c.wsClient.Disconnect(); err != nil {
        return err
    }
    
    c.isConnected = false
    return nil
}

// handleWebSocketMessages processes WebSocket messages and routes them to appropriate channels
func (c *Client) handleWebSocketMessages(ctx context.Context) {
    tickerCh := c.wsClient.TickerChannel()
    orderCh := c.wsClient.OrderChannel()
    
    for {
        select {
        case <-ctx.Done():
            return
            
        case ticker := <-tickerCh:
            if ticker == nil {
                continue
            }
            
            c.mu.RLock()
            channels, ok := c.tickerChannels[ticker.Symbol]
            c.mu.RUnlock()
            
            if !ok {
                continue
            }
            
            // Broadcast to all subscribers
            for _, ch := range channels {
                select {
                case ch <- ticker:
                    // Sent successfully
                case default:
                    // Channel is full or closed, skip
                }
            }
            
        case order := <-orderCh:
            // Process order updates if needed
            _ = order
        }
    }
}

// Implementation of domain.service.ExchangeService interface methods

// GetTicker fetches the current ticker data for a specific symbol
func (c *Client) GetTicker(ctx context.Context, symbol string) (*models.Ticker, error) {
    return c.restClient.GetTicker(ctx, symbol)
}

// GetAllTickers fetches tickers for all symbols
func (c *Client) GetAllTickers(ctx context.Context) (map[string]*models.Ticker, error) {
    return c.restClient.GetAllTickers(ctx)
}

// GetKlines fetches candlestick data for a symbol
func (c *Client) GetKlines(ctx context.Context, symbol string, interval string, limit int) ([]*models.Kline, error) {
    return c.restClient.GetKlines(ctx, symbol, interval, limit)
}

// GetWallet fetches account information including balances
func (c *Client) GetWallet(ctx context.Context) (*models.Wallet, error) {
    return c.restClient.GetWallet(ctx)
}

// PlaceOrder sends a new order to the exchange
func (c *Client) PlaceOrder(ctx context.Context, order *models.Order) (*models.Order, error) {
    return c.restClient.PlaceOrder(ctx, order)
}

// CancelOrder cancels an existing order
func (c *Client) CancelOrder(ctx context.Context, orderID string, symbol string) error {
    return c.restClient.CancelOrder(ctx, orderID, symbol)
}

// GetOrder retrieves information for a specific order
func (c *Client) GetOrder(ctx context.Context, orderID string, symbol string) (*models.Order, error) {
    return c.restClient.GetOrder(ctx, orderID, symbol)
}

// GetOpenOrders retrieves all open orders for a symbol
func (c *Client) GetOpenOrders(ctx context.Context, symbol string) ([]*models.Order, error) {
    return c.restClient.GetOpenOrders(ctx, symbol)
}

// GetNewCoins fetches a list of new coins (recent listings)
func (c *Client) GetNewCoins(ctx context.Context) ([]*models.NewCoin, error) {
    return c.restClient.GetNewCoins(ctx)
}

// SubscribeToTickers subscribes to ticker updates for the given symbols
func (c *Client) SubscribeToTickers(ctx context.Context, symbols []string, updates chan<- *models.Ticker) error {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    if !c.isConnected {
        if err := c.Connect(ctx); err != nil {
            return err
        }
    }
    
    // Subscribe to all requested symbols
    if err := c.wsClient.SubscribeToTickers(symbols); err != nil {
        return err
    }
    
    // Register the output channel for each symbol
    for _, symbol := range symbols {
        if _, ok := c.tickerChannels[symbol]; !ok {
            c.tickerChannels[symbol] = make([]chan *models.Ticker, 0, 1)
        }
        c.tickerChannels[symbol] = append(c.tickerChannels[symbol], updates)
    }
    
    return nil
}

// UnsubscribeFromTickers unsubscribes from ticker updates for the given symbols
func (c *Client) UnsubscribeFromTickers(ctx context.Context, symbols []string) error {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    if !c.isConnected {
        return nil // Already disconnected
    }
    
    if err := c.wsClient.UnsubscribeFromTickers(symbols); err != nil {
        return err
    }
    
    // Remove the channels for the unsubscribed symbols
    for _, symbol := range symbols {
        delete(c.tickerChannels, symbol)
    }
    
    return nil
}
```

## Dependency Injection and Factory Pattern

The client uses dependency injection and the factory pattern to simplify testing and configuration:

```go
// internal/platform/mexc/factory.go
package mexc

import (
    "github.com/ryanlisse/cryptobot/internal/domain/service"
)

// Factory creates MEXC exchange clients
type Factory struct {
    apiKey    string
    secretKey string
    baseURL   string
}

// NewFactory creates a new MEXC client factory
func NewFactory(apiKey, secretKey, baseURL string) *Factory {
    return &Factory{
        apiKey:    apiKey,
        secretKey: secretKey,
        baseURL:   baseURL,
    }
}

// Create instantiates a new MEXC client
func (f *Factory) Create() (service.ExchangeService, error) {
    return NewClient(f.apiKey, f.secretKey, f.baseURL)
}
```

## Usage Example

```go
// Example usage in an application service
func NewTradeService(config *config.Config) (*TradeService, error) {
    // Create the MEXC client factory
    mexcFactory := mexc.NewFactory(
        config.MEXC.APIKey,
        config.MEXC.SecretKey,
        config.MEXC.BaseURL,
    )
    
    // Create the MEXC client
    mexcClient, err := mexcFactory.Create()
    if err != nil {
        return nil, fmt.Errorf("failed to create MEXC client: %w", err)
    }
    
    // Use the client in the service
    return &TradeService{
        exchangeService: mexcClient,
        // ... other dependencies
    }, nil
}
```

This approach makes it easy to:
1. Mock the exchange service in tests
2. Configure the client through dependency injection
3. Switch to a different exchange implementation if needed
