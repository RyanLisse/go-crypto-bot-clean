// Package mexc provides a unified client for interacting with the MEXC exchange
package mexc

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/ryanlisse/go-crypto-bot/internal/config"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/service"
	"github.com/ryanlisse/go-crypto-bot/internal/platform/mexc/rest"
	"github.com/ryanlisse/go-crypto-bot/internal/platform/mexc/websocket"
	"github.com/ryanlisse/go-crypto-bot/pkg/ratelimiter"
)

// Client implements the service.ExchangeService interface by combining
// REST and WebSocket clients for MEXC exchange.
type Client struct {
	restClient      *rest.Client
	wsClient        *websocket.Client
	connRateLimiter *ratelimiter.TokenBucketRateLimiter
	subRateLimiter  *ratelimiter.TokenBucketRateLimiter
	tickerChannels  map[string][]chan<- *models.Ticker
	mu              sync.RWMutex
	isConnected     bool
	logger          *log.Logger

	// For managing goroutines
	distributeTickersDone chan struct{}
	distributeTickersStop chan struct{}
}

// ClientOption configures a MEXC client
type ClientOption func(*Client)

// WithLogger sets a custom logger for the client
func WithLogger(logger *log.Logger) ClientOption {
	return func(c *Client) {
		c.logger = logger
	}
}

// NewClient creates a new MEXC client with both REST and WebSocket capabilities.
// It initializes shared rate limiters and manages both client types.
func NewClient(cfg *config.Config, options ...ClientOption) (service.ExchangeService, error) {
	// Initialize default logger if none provided
	logger := log.New(log.Writer(), "[MEXC] ", log.LstdFlags)

	// Create shared rate limiters
	connRateLimiter := ratelimiter.NewTokenBucketRateLimiter(
		cfg.ConnectionRateLimiter.RequestsPerSecond,
		float64(cfg.ConnectionRateLimiter.BurstCapacity),
	)

	subRateLimiter := ratelimiter.NewTokenBucketRateLimiter(
		cfg.SubscriptionRateLimiter.RequestsPerSecond,
		float64(cfg.SubscriptionRateLimiter.BurstCapacity),
	)

	// Create REST client
	restClient, err := rest.NewClient(
		cfg.Mexc.APIKey,
		cfg.Mexc.SecretKey,
		rest.WithBaseURL(cfg.Mexc.BaseURL),
		rest.WithPublicRateLimiter(connRateLimiter),
		rest.WithPrivateRateLimiter(connRateLimiter),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create REST client: %w", err)
	}

	// Create WebSocket client
	wsClient, err := websocket.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create WebSocket client: %w", err)
	}

	// Create unified client
	client := &Client{
		restClient:      restClient,
		wsClient:        wsClient,
		connRateLimiter: connRateLimiter,
		subRateLimiter:  subRateLimiter,
		tickerChannels:  make(map[string][]chan<- *models.Ticker),
		logger:          logger,

		// Initialize channels for goroutine management
		distributeTickersDone: make(chan struct{}),
		distributeTickersStop: make(chan struct{}),
	}

	// Apply options
	for _, option := range options {
		option(client)
	}

	// Start background goroutine to distribute tickers to subscribers
	go client.distributeTickers()

	return client, nil
}

// Connect establishes a connection to the WebSocket API
func (c *Client) Connect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isConnected {
		return nil
	}

	if err := c.wsClient.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect to WebSocket: %w", err)
	}

	c.isConnected = true
	return nil
}

// Disconnect closes the WebSocket connection and cleans up resources
func (c *Client) Disconnect() error {
	c.mu.Lock()

	if !c.isConnected {
		c.mu.Unlock()
		return nil
	}

	// Signal the ticker distribution goroutine to stop
	close(c.distributeTickersStop)

	// Disconnect the websocket
	var wsErr error
	if err := c.wsClient.Disconnect(); err != nil {
		wsErr = fmt.Errorf("failed to disconnect from WebSocket: %w", err)
	}

	c.isConnected = false
	c.mu.Unlock()

	// Wait for distributor goroutine to finish
	<-c.distributeTickersDone

	// Create new channels for potential reconnection
	c.mu.Lock()
	c.distributeTickersStop = make(chan struct{})
	c.distributeTickersDone = make(chan struct{})
	c.mu.Unlock()

	// Restart goroutine if we reconnect
	if wsErr == nil {
		go c.distributeTickers()
	}

	return wsErr
}

// SubscribeToTickers subscribes to ticker updates for the given symbols
func (c *Client) SubscribeToTickers(ctx context.Context, symbols []string, updates chan<- *models.Ticker) error {
	if err := c.ensureConnected(ctx); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Subscribe to tickers via WebSocket - use provided context
	if err := c.wsClient.SubscribeToTickers(ctx, symbols); err != nil {
		return fmt.Errorf("failed to subscribe to tickers: %w", err)
	}

	// Register the updates channel for each symbol
	for _, symbol := range symbols {
		c.tickerChannels[symbol] = append(c.tickerChannels[symbol], updates)
	}

	return nil
}

// UnsubscribeFromTickers removes the subscription for the given symbols
func (c *Client) UnsubscribeFromTickers(ctx context.Context, symbols []string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.isConnected {
		return nil
	}

	// Unsubscribe from tickers via WebSocket
	if err := c.wsClient.UnsubscribeFromTickers(symbols); err != nil {
		return fmt.Errorf("failed to unsubscribe from tickers: %w", err)
	}

	// Remove the tickers
	for _, symbol := range symbols {
		delete(c.tickerChannels, symbol)
	}

	return nil
}

// GetTicker retrieves the current ticker for a symbol using the REST API
func (c *Client) GetTicker(ctx context.Context, symbol string) (*models.Ticker, error) {
	return c.restClient.GetTicker(ctx, symbol)
}

// GetAllTickers retrieves all current tickers using the REST API
func (c *Client) GetAllTickers(ctx context.Context) (map[string]*models.Ticker, error) {
	return c.restClient.GetAllTickers(ctx)
}

// GetOrderBook retrieves the current order book for a symbol using the REST API
func (c *Client) GetOrderBook(ctx context.Context, symbol string, limit int) (*models.OrderBookUpdate, error) {
	return c.restClient.GetOrderBook(ctx, symbol, limit)
}

// GetWallet retrieves the current wallet balance using the REST API
func (c *Client) GetWallet(ctx context.Context) (*models.Wallet, error) {
	return c.restClient.GetWallet(ctx)
}

// PlaceOrder places a new order using the REST API
func (c *Client) PlaceOrder(ctx context.Context, order *models.Order) (*models.Order, error) {
	return c.restClient.PlaceOrder(ctx, order)
}

// CancelOrder cancels an existing order using the REST API
func (c *Client) CancelOrder(ctx context.Context, orderID, symbol string) error {
	return c.restClient.CancelOrder(ctx, orderID, symbol)
}

// GetOrder retrieves an order's details using the REST API
func (c *Client) GetOrder(ctx context.Context, orderID, symbol string) (*models.Order, error) {
	return c.restClient.GetOrder(ctx, orderID, symbol)
}

// GetOpenOrders retrieves all open orders using the REST API
func (c *Client) GetOpenOrders(ctx context.Context, symbol string) ([]*models.Order, error) {
	return c.restClient.GetOpenOrders(ctx, symbol)
}

// GetKlines retrieves klines/candlesticks for a symbol using the REST API
func (c *Client) GetKlines(ctx context.Context, symbol, interval string, limit int) ([]*models.Kline, error) {
	return c.restClient.GetKlines(ctx, symbol, interval, limit)
}

// GetNewCoins retrieves new coins added to the exchange using the REST API
func (c *Client) GetNewCoins(ctx context.Context) ([]*models.NewCoin, error) {
	return c.restClient.GetNewCoins(ctx)
}

// Helper methods

// distributeTickers forwards ticker updates from the WebSocket to the registered channels
func (c *Client) distributeTickers() {
	defer func() {
		// Signal that we're done when exiting
		close(c.distributeTickersDone)
	}()

	ticker := c.wsClient.TickerChannel()
	for {
		select {
		case <-c.distributeTickersStop:
			// Clean exit requested
			return
		case t, ok := <-ticker:
			if !ok {
				// Channel was closed
				return
			}

			c.mu.RLock()
			channels, exists := c.tickerChannels[t.Symbol]
			c.mu.RUnlock()

			if !exists {
				continue
			}

			for _, ch := range channels {
				select {
				case ch <- t:
					// Ticker sent successfully
				default:
					c.logger.Printf("Warning: Ticker channel for %s is full or closed, update dropped", t.Symbol)
				}
			}
		}
	}
}

// ensureConnected makes sure the client is connected to the WebSocket
func (c *Client) ensureConnected(ctx context.Context) error {
	c.mu.RLock()
	connected := c.isConnected
	c.mu.RUnlock()

	if !connected {
		return c.Connect(ctx)
	}
	return nil
}
