package mexc

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/adapter/infrastructure/gateway/mexc/apikeystore"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/adapter/infrastructure/gateway/mexc/rest"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/adapter/infrastructure/gateway/mexc/websocket"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port/gateway"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
)

// Gateway implements the gateway.MEXCGateway interface
// It adapts the MEXC REST and WebSocket clients to the domain gateway interface
type Gateway struct {
	restClient   *rest.Client
	wsClient     *websocket.Client
	keyStore     apikeystore.KeyStore
	logger       *zerolog.Logger
	defaultKeyID string
	config       GatewayConfig
}

// GatewayConfig holds configuration for the MEXC gateway
type GatewayConfig struct {
	DefaultKeyID string
	BaseURL      string
	WSEndpoint   string
}

// DefaultGatewayConfig returns a default configuration for the MEXC gateway
func DefaultGatewayConfig() GatewayConfig {
	return GatewayConfig{
		DefaultKeyID: "default",
		BaseURL:      "https://api.mexc.com",
		WSEndpoint:   "wss://stream.mexc.com/ws",
	}
}

// NewGateway creates a new MEXC gateway
func NewGateway(keyStore apikeystore.KeyStore, config GatewayConfig, logger *zerolog.Logger) (gateway.MEXCGateway, error) {
	// Get API credentials from key store
	creds, err := keyStore.GetAPIKey(config.DefaultKeyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get API key from key store: %w", err)
	}

	// Create REST client
	restClient := rest.NewClient(creds.APIKey, creds.SecretKey, logger)

	// Create WebSocket client
	wsClient := websocket.NewClient(creds.APIKey, creds.SecretKey, logger)

	return &Gateway{
		restClient:   restClient,
		wsClient:     wsClient,
		keyStore:     keyStore,
		logger:       logger,
		defaultKeyID: config.DefaultKeyID,
		config:       config,
	}, nil
}

// GetExchangeInfo retrieves general exchange information
func (g *Gateway) GetExchangeInfo(ctx context.Context) (*model.ExchangeInfo, error) {
	g.logger.Debug().Msg("Getting exchange info from MEXC")

	// Get exchange info from REST client
	exchangeInfo, err := g.restClient.GetExchangeInfo(ctx)
	if err != nil {
		g.logger.Error().Err(err).Msg("Failed to get exchange info from MEXC")
		return nil, fmt.Errorf("failed to get exchange info from MEXC: %w", err)
	}

	return exchangeInfo, nil
}

// GetSymbols retrieves all available trading symbols from MEXC
func (g *Gateway) GetSymbols(ctx context.Context) ([]model.SymbolInfo, error) {
	g.logger.Debug().Msg("Getting symbols from MEXC")

	// Get exchange info from REST client
	exchangeInfo, err := g.restClient.GetExchangeInfo(ctx)
	if err != nil {
		g.logger.Error().Err(err).Msg("Failed to get exchange info from MEXC")
		return nil, fmt.Errorf("failed to get exchange info from MEXC: %w", err)
	}

	return exchangeInfo.Symbols, nil
}

// GetTicker retrieves current ticker data for a symbol
func (g *Gateway) GetTicker(ctx context.Context, symbol string) (model.Ticker, error) {
	g.logger.Debug().Str("symbol", symbol).Msg("Getting ticker from MEXC")

	// Get ticker from REST client
	ticker, err := g.restClient.GetMarketData(ctx, symbol)
	if err != nil {
		g.logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to get ticker from MEXC")
		return model.Ticker{}, fmt.Errorf("failed to get ticker from MEXC: %w", err)
	}

	return *ticker, nil
}

// GetOrderBook retrieves the order book for a symbol
func (g *Gateway) GetOrderBook(ctx context.Context, symbol string, depth int) (model.OrderBook, error) {
	g.logger.Debug().Str("symbol", symbol).Int("depth", depth).Msg("Getting order book from MEXC")

	// Get order book from REST client
	orderBook, err := g.restClient.GetOrderBook(ctx, symbol, depth)
	if err != nil {
		g.logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to get order book from MEXC")
		return model.OrderBook{}, fmt.Errorf("failed to get order book from MEXC: %w", err)
	}

	return *orderBook, nil
}

// GetAccountInfo retrieves account information
func (g *Gateway) GetAccountInfo(ctx context.Context) (model.AccountInfo, error) {
	g.logger.Debug().Msg("Getting account info from MEXC")

	// Get account info from REST client
	accountInfo, err := g.restClient.GetAccountInfo(ctx)
	if err != nil {
		g.logger.Error().Err(err).Msg("Failed to get account info from MEXC")
		return model.AccountInfo{}, fmt.Errorf("failed to get account info from MEXC: %w", err)
	}

	return accountInfo, nil
}

// GetAssetBalance retrieves the balance for a specific asset
func (g *Gateway) GetAssetBalance(ctx context.Context, asset string) (decimal.Decimal, error) {
	g.logger.Debug().Str("asset", asset).Msg("Getting asset balance from MEXC")

	// Get account info from REST client
	accountInfo, err := g.restClient.GetAccountInfo(ctx)
	if err != nil {
		g.logger.Error().Err(err).Str("asset", asset).Msg("Failed to get account info from MEXC")
		return decimal.Zero, fmt.Errorf("failed to get account info from MEXC: %w", err)
	}

	// Find the balance for the specified asset
	for _, balance := range accountInfo.Balances {
		if string(balance.Asset) == asset {
			return decimal.NewFromFloat(balance.Free), nil
		}
	}

	// Asset not found, return zero balance
	return decimal.Zero, nil
}

// PlaceOrder places a new order on MEXC
func (g *Gateway) PlaceOrder(ctx context.Context, orderRequest model.OrderRequest) (model.OrderResponse, error) {
	g.logger.Debug().Str("symbol", orderRequest.Symbol).Str("side", string(orderRequest.Side)).Str("type", string(orderRequest.Type)).Float64("quantity", orderRequest.Quantity).Float64("price", orderRequest.Price).Msg("Placing order on MEXC")

	// Place order using REST client
	result, err := g.restClient.PlaceOrder(ctx, orderRequest.Symbol, orderRequest.Side, orderRequest.Type, orderRequest.Quantity, orderRequest.Price, orderRequest.TimeInForce)
	if err != nil {
		g.logger.Error().Err(err).Str("symbol", orderRequest.Symbol).Msg("Failed to place order on MEXC")
		return model.OrderResponse{IsSuccess: false}, fmt.Errorf("failed to place order on MEXC: %w", err)
	}

	// Convert to OrderResponse
	orderResponse := model.OrderResponse{
		Order: model.Order{
			OrderID:       result.OrderID,
			ClientOrderID: result.ClientOrderID,
			Symbol:        result.Symbol,
			CreatedAt:     result.CreatedAt,
			Price:         result.Price,
			Quantity:      result.Quantity,
			ExecutedQty:   result.ExecutedQty,
			Status:        result.Status,
			Type:          result.Type,
			Side:          result.Side,
			Exchange:      "MEXC",
		},
		IsSuccess: true,
	}

	return orderResponse, nil
}

// CancelOrder cancels an existing order on MEXC
func (g *Gateway) CancelOrder(ctx context.Context, symbol, orderID string) error {
	g.logger.Debug().Str("symbol", symbol).Str("orderID", orderID).Msg("Cancelling order on MEXC")

	// Cancel order using REST client
	err := g.restClient.CancelOrder(ctx, symbol, orderID)
	if err != nil {
		g.logger.Error().Err(err).Str("symbol", symbol).Str("orderID", orderID).Msg("Failed to cancel order on MEXC")
		return fmt.Errorf("failed to cancel order on MEXC: %w", err)
	}

	return nil
}

// SubscribeToTicker subscribes to ticker updates for a symbol
func (g *Gateway) SubscribeToTicker(ctx context.Context, symbol string, handler func(model.Ticker)) error {
	g.logger.Debug().Str("symbol", symbol).Msg("Subscribing to ticker updates")

	// Connect to WebSocket if not already connected
	if err := g.wsClient.Connect(); err != nil {
		g.logger.Error().Err(err).Msg("Failed to connect to WebSocket")
		return fmt.Errorf("failed to connect to WebSocket: %w", err)
	}

	// Subscribe to ticker channel
	ch, err := g.wsClient.Subscribe("ticker", symbol)
	if err != nil {
		g.logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to subscribe to ticker channel")
		return fmt.Errorf("failed to subscribe to ticker channel: %w", err)
	}

	// Start goroutine to handle messages
	go func() {
		defer func() {
			if r := recover(); r != nil {
				g.logger.Error().Interface("recover", r).Msg("Recovered from panic in ticker handler")
			}
		}()

		for {
			select {
			case <-ctx.Done():
				g.logger.Debug().Str("symbol", symbol).Msg("Context cancelled, stopping ticker subscription")
				_ = g.wsClient.Unsubscribe("ticker", symbol)
				return
			case msg, ok := <-ch:
				if !ok {
					g.logger.Debug().Str("symbol", symbol).Msg("Channel closed, stopping ticker subscription")
					return
				}

				// Parse message
				var tickerMsg struct {
					Symbol    string `json:"s"`
					LastPrice string `json:"c"`
					Volume    string `json:"v"`
					Timestamp int64  `json:"E"`
				}

				if err := json.Unmarshal(msg, &tickerMsg); err != nil {
					g.logger.Error().Err(err).Str("message", string(msg)).Msg("Failed to parse ticker message")
					continue
				}

				// Convert to Ticker model
				lastPrice, _ := strconv.ParseFloat(tickerMsg.LastPrice, 64)
				volume, _ := strconv.ParseFloat(tickerMsg.Volume, 64)

				ticker := model.Ticker{
					Symbol:    tickerMsg.Symbol,
					Exchange:  "MEXC",
					LastPrice: lastPrice,
					Volume:    volume,
					Timestamp: time.Unix(0, tickerMsg.Timestamp*int64(time.Millisecond)),
				}

				// Call handler
				handler(ticker)
			}
		}
	}()

	return nil
}

// SubscribeToOrderBook subscribes to order book updates for a symbol
func (g *Gateway) SubscribeToOrderBook(ctx context.Context, symbol string, depth int, handler func(model.OrderBook)) error {
	g.logger.Debug().Str("symbol", symbol).Int("depth", depth).Msg("Subscribing to order book updates")

	// Connect to WebSocket if not already connected
	if err := g.wsClient.Connect(); err != nil {
		g.logger.Error().Err(err).Msg("Failed to connect to WebSocket")
		return fmt.Errorf("failed to connect to WebSocket: %w", err)
	}

	// Subscribe to order book channel
	ch, err := g.wsClient.Subscribe("depth", symbol)
	if err != nil {
		g.logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to subscribe to order book channel")
		return fmt.Errorf("failed to subscribe to order book channel: %w", err)
	}

	// Start goroutine to handle messages
	go func() {
		defer func() {
			if r := recover(); r != nil {
				g.logger.Error().Interface("recover", r).Msg("Recovered from panic in order book handler")
			}
		}()

		for {
			select {
			case <-ctx.Done():
				g.logger.Debug().Str("symbol", symbol).Msg("Context cancelled, stopping order book subscription")
				_ = g.wsClient.Unsubscribe("depth", symbol)
				return
			case msg, ok := <-ch:
				if !ok {
					g.logger.Debug().Str("symbol", symbol).Msg("Channel closed, stopping order book subscription")
					return
				}

				// Parse message
				var orderBookMsg struct {
					Symbol    string     `json:"s"`
					Timestamp int64      `json:"E"`
					Bids      [][]string `json:"b"`
					Asks      [][]string `json:"a"`
				}

				if err := json.Unmarshal(msg, &orderBookMsg); err != nil {
					g.logger.Error().Err(err).Str("message", string(msg)).Msg("Failed to parse order book message")
					continue
				}

				// Convert to OrderBook model
				orderBook := model.OrderBook{
					Symbol:    orderBookMsg.Symbol,
					Exchange:  "MEXC",
					Timestamp: time.Unix(0, orderBookMsg.Timestamp*int64(time.Millisecond)),
					Bids:      make([]model.OrderBookEntry, 0, len(orderBookMsg.Bids)),
					Asks:      make([]model.OrderBookEntry, 0, len(orderBookMsg.Asks)),
				}

				// Parse bids
				for _, bid := range orderBookMsg.Bids {
					if len(bid) < 2 {
						continue
					}
					price, _ := strconv.ParseFloat(bid[0], 64)
					quantity, _ := strconv.ParseFloat(bid[1], 64)
					orderBook.Bids = append(orderBook.Bids, model.OrderBookEntry{
						Price:    price,
						Quantity: quantity,
					})
				}

				// Parse asks
				for _, ask := range orderBookMsg.Asks {
					if len(ask) < 2 {
						continue
					}
					price, _ := strconv.ParseFloat(ask[0], 64)
					quantity, _ := strconv.ParseFloat(ask[1], 64)
					orderBook.Asks = append(orderBook.Asks, model.OrderBookEntry{
						Price:    price,
						Quantity: quantity,
					})
				}

				// Limit depth if specified
				if depth > 0 {
					if len(orderBook.Bids) > depth {
						orderBook.Bids = orderBook.Bids[:depth]
					}
					if len(orderBook.Asks) > depth {
						orderBook.Asks = orderBook.Asks[:depth]
					}
				}

				// Call handler
				handler(orderBook)
			}
		}
	}()

	return nil
}

// Unsubscribe unsubscribes from a channel
func (g *Gateway) Unsubscribe(ctx context.Context, channel, symbol string) error {
	g.logger.Debug().Str("channel", channel).Str("symbol", symbol).Msg("Unsubscribing from channel")

	// Unsubscribe from channel
	err := g.wsClient.Unsubscribe(channel, symbol)
	if err != nil {
		g.logger.Error().Err(err).Str("channel", channel).Str("symbol", symbol).Msg("Failed to unsubscribe from channel")
		return fmt.Errorf("failed to unsubscribe from channel: %w", err)
	}

	return nil
}

// ChangeAPIKey changes the API key used by the gateway
func (g *Gateway) ChangeAPIKey(ctx context.Context, keyID string) error {
	g.logger.Debug().Str("keyID", keyID).Msg("Changing API key")

	// Get API credentials from key store
	creds, err := g.keyStore.GetAPIKey(keyID)
	if err != nil {
		g.logger.Error().Err(err).Str("keyID", keyID).Msg("Failed to get API key from key store")
		return fmt.Errorf("failed to get API key from key store: %w", err)
	}

	// Create new REST client with the new credentials
	g.restClient = rest.NewClient(creds.APIKey, creds.SecretKey, g.logger)

	// Disconnect existing WebSocket client if connected
	if g.wsClient != nil {
		_ = g.wsClient.Disconnect()
	}

	// Create new WebSocket client with the new credentials
	g.wsClient = websocket.NewClient(creds.APIKey, creds.SecretKey, g.logger)

	// Update default key ID
	g.defaultKeyID = keyID

	return nil
}
