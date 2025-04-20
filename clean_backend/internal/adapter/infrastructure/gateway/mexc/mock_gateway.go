package mexc

import (
	"context"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port/gateway"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
)

// MockGateway is a mock implementation of the gateway.MEXCGateway interface
type MockGateway struct {
	logger *zerolog.Logger
}

// NewMockGateway creates a new mock MEXC gateway
func NewMockGateway(logger *zerolog.Logger) *MockGateway {
	return &MockGateway{
		logger: logger,
	}
}

// GetExchangeInfo retrieves general exchange information
func (g *MockGateway) GetExchangeInfo(ctx context.Context) (*model.ExchangeInfo, error) {
	g.logger.Debug().Msg("Getting exchange info from MEXC (mock)")

	// Return mock exchange info
	return &model.ExchangeInfo{
		Symbols: []model.SymbolInfo{
			{
				Symbol:     "BTCUSDT",
				Status:     "TRADING",
				BaseAsset:  "BTC",
				QuoteAsset: "USDT",
			},
			{
				Symbol:     "ETHUSDT",
				Status:     "TRADING",
				BaseAsset:  "ETH",
				QuoteAsset: "USDT",
			},
		},
	}, nil
}

// GetSymbols retrieves all available trading symbols from MEXC
func (g *MockGateway) GetSymbols(ctx context.Context) ([]model.SymbolInfo, error) {
	g.logger.Debug().Msg("Getting symbols from MEXC (mock)")

	// Return mock symbols
	return []model.SymbolInfo{
		{
			Symbol:     "BTCUSDT",
			Status:     "TRADING",
			BaseAsset:  "BTC",
			QuoteAsset: "USDT",
		},
		{
			Symbol:     "ETHUSDT",
			Status:     "TRADING",
			BaseAsset:  "ETH",
			QuoteAsset: "USDT",
		},
	}, nil
}

// GetTicker retrieves current ticker data for a symbol
func (g *MockGateway) GetTicker(ctx context.Context, symbol string) (model.Ticker, error) {
	g.logger.Debug().Str("symbol", symbol).Msg("Getting ticker from MEXC (mock)")

	// Return mock ticker
	now := time.Now()
	return model.Ticker{
		Symbol:             symbol,
		Exchange:           "MEXC",
		ExchangeName:       "MEXC",
		LastPrice:          50000.0,
		Price:              50000.0,
		PriceChange:        100.0,
		PriceChangePercent: 0.2,
		PercentChange:      0.2,
		HighPrice:          51000.0,
		High24h:            51000.0,
		LowPrice:           49000.0,
		Low24h:             49000.0,
		Volume:             1000.0,
		QuoteVolume:        50000000.0,
		Timestamp:          now,
		LastUpdated:        now,
	}, nil
}

// GetOrderBook retrieves the order book for a symbol
func (g *MockGateway) GetOrderBook(ctx context.Context, symbol string, depth int) (model.OrderBook, error) {
	g.logger.Debug().Str("symbol", symbol).Int("depth", depth).Msg("Getting order book from MEXC (mock)")

	// Return mock order book
	return model.OrderBook{
		Symbol:       symbol,
		LastUpdateID: 12345,
		Bids: []model.OrderBookEntry{
			{Price: 49900.0, Quantity: 1.0},
			{Price: 49800.0, Quantity: 2.0},
		},
		Asks: []model.OrderBookEntry{
			{Price: 50100.0, Quantity: 1.0},
			{Price: 50200.0, Quantity: 2.0},
		},
		Timestamp: time.Now(),
	}, nil
}

// GetAccountInfo retrieves account information
func (g *MockGateway) GetAccountInfo(ctx context.Context) (model.AccountInfo, error) {
	g.logger.Debug().Msg("Getting account info from MEXC (mock)")

	// Return mock account info
	return model.AccountInfo{
		UserID:      "mock-user-123",
		CanTrade:    true,
		CanWithdraw: true,
		CanDeposit:  true,
		Balances: []model.Balance{
			{
				Asset:  "BTC",
				Free:   1.0,
				Locked: 0.0,
				Total:  1.0,
			},
			{
				Asset:  "USDT",
				Free:   50000.0,
				Locked: 0.0,
				Total:  50000.0,
			},
		},
		LastUpdated: time.Now(),
	}, nil
}

// GetAssetBalance retrieves the balance for a specific asset
func (g *MockGateway) GetAssetBalance(ctx context.Context, asset string) (decimal.Decimal, error) {
	g.logger.Debug().Str("asset", asset).Msg("Getting asset balance from MEXC (mock)")

	// Return mock balance based on asset
	switch asset {
	case "BTC":
		return decimal.NewFromFloat(1.0), nil
	case "ETH":
		return decimal.NewFromFloat(10.0), nil
	case "USDT":
		return decimal.NewFromFloat(50000.0), nil
	default:
		return decimal.Zero, nil
	}
}

// PlaceOrder places a new order on MEXC
func (g *MockGateway) PlaceOrder(ctx context.Context, orderRequest model.OrderRequest) (model.OrderResponse, error) {
	g.logger.Debug().Str("symbol", orderRequest.Symbol).Str("side", string(orderRequest.Side)).Str("type", string(orderRequest.Type)).Float64("quantity", orderRequest.Quantity).Float64("price", orderRequest.Price).Msg("Placing order on MEXC (mock)")

	// Return mock order response
	now := time.Now()
	return model.OrderResponse{
		Order: model.Order{
			OrderID:       "mock-order-123",
			ClientOrderID: "mock-client-order-123",
			Symbol:        orderRequest.Symbol,
			CreatedAt:     now,
			UpdatedAt:     now,
			Price:         orderRequest.Price,
			Quantity:      orderRequest.Quantity,
			ExecutedQty:   0.0,
			Status:        model.OrderStatusNew,
			Type:          orderRequest.Type,
			Side:          orderRequest.Side,
			Exchange:      "MEXC",
		},
		IsSuccess: true,
	}, nil
}

// CancelOrder cancels an existing order on MEXC
func (g *MockGateway) CancelOrder(ctx context.Context, symbol, orderID string) error {
	g.logger.Debug().Str("symbol", symbol).Str("orderID", orderID).Msg("Cancelling order on MEXC (mock)")

	// Return success
	return nil
}

// SubscribeToTicker subscribes to ticker updates for a symbol
func (g *MockGateway) SubscribeToTicker(ctx context.Context, symbol string, handler func(model.Ticker)) error {
	g.logger.Debug().Str("symbol", symbol).Msg("Subscribing to ticker updates (mock)")

	// Simulate a ticker update in a goroutine
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case t := <-ticker.C:
				// Create a mock ticker
				mockTicker := model.Ticker{
					Symbol:    symbol,
					Exchange:  "MEXC",
					LastPrice: 50000.0 + float64(t.Nanosecond()%1000)/100.0, // Add some variation
					Volume:    1000.0,
					Timestamp: t,
				}

				// Call the handler
				handler(mockTicker)
			}
		}
	}()

	return nil
}

// SubscribeToOrderBook subscribes to order book updates for a symbol
func (g *MockGateway) SubscribeToOrderBook(ctx context.Context, symbol string, depth int, handler func(model.OrderBook)) error {
	g.logger.Debug().Str("symbol", symbol).Int("depth", depth).Msg("Subscribing to order book updates (mock)")

	// Simulate order book updates in a goroutine
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case t := <-ticker.C:
				// Create a mock order book
				mockOrderBook := model.OrderBook{
					Symbol:    symbol,
					Exchange:  "MEXC",
					Timestamp: t,
					Bids:      make([]model.OrderBookEntry, depth),
					Asks:      make([]model.OrderBookEntry, depth),
				}

				// Generate some mock bids and asks
				basePrice := 50000.0
				for i := 0; i < depth; i++ {
					mockOrderBook.Bids[i] = model.OrderBookEntry{
						Price:    basePrice - float64(i)*10.0,
						Quantity: 1.0 + float64(i)*0.1,
					}
					mockOrderBook.Asks[i] = model.OrderBookEntry{
						Price:    basePrice + float64(i)*10.0,
						Quantity: 1.0 + float64(i)*0.1,
					}
				}

				// Call the handler
				handler(mockOrderBook)
			}
		}
	}()

	return nil
}

// Unsubscribe unsubscribes from a channel
func (g *MockGateway) Unsubscribe(ctx context.Context, channel, symbol string) error {
	g.logger.Debug().Str("channel", channel).Str("symbol", symbol).Msg("Unsubscribing from channel (mock)")

	// Return success
	return nil
}

// ChangeAPIKey changes the API key used by the gateway
func (g *MockGateway) ChangeAPIKey(ctx context.Context, keyID string) error {
	g.logger.Debug().Str("keyID", keyID).Msg("Changing API key (mock)")

	// Return success
	return nil
}

// Ensure MockGateway implements gateway.MEXCGateway
var _ gateway.MEXCGateway = (*MockGateway)(nil)
