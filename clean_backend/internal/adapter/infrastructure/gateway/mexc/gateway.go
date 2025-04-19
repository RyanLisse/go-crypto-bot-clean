package mexc

import (
	"context"
	"fmt"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port/gateway"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
)

// Gateway implements the gateway.MEXCGateway interface
// It adapts the MEXCClient to the domain gateway interface
type Gateway struct {
	client port.MEXCClient
	logger *zerolog.Logger
}

// NewGateway creates a new MEXC gateway
func NewGateway(client port.MEXCClient, logger *zerolog.Logger) gateway.MEXCGateway {
	return &Gateway{
		client: client,
		logger: logger,
	}
}

// GetSymbols retrieves all available trading symbols from MEXC
func (g *Gateway) GetSymbols(ctx context.Context) ([]model.SymbolInfo, error) {
	g.logger.Debug().Msg("Getting symbols from MEXC")

	// Get exchange info from client
	exchangeInfo, err := g.client.GetExchangeInfo(ctx)
	if err != nil {
		g.logger.Error().Err(err).Msg("Failed to get exchange info from MEXC")
		return nil, fmt.Errorf("failed to get exchange info from MEXC: %w", err)
	}

	return exchangeInfo.Symbols, nil
}

// GetTicker retrieves current ticker data for a symbol
func (g *Gateway) GetTicker(ctx context.Context, symbol string) (model.Ticker, error) {
	g.logger.Debug().Str("symbol", symbol).Msg("Getting ticker from MEXC")

	// Get ticker from client
	ticker, err := g.client.GetMarketData(ctx, symbol)
	if err != nil {
		g.logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to get ticker from MEXC")
		return model.Ticker{}, fmt.Errorf("failed to get ticker from MEXC: %w", err)
	}

	return *ticker, nil
}

// GetOrderBook retrieves the order book for a symbol
func (g *Gateway) GetOrderBook(ctx context.Context, symbol string, depth int) (model.OrderBook, error) {
	g.logger.Debug().Str("symbol", symbol).Int("depth", depth).Msg("Getting order book from MEXC")

	// Get order book from client
	orderBook, err := g.client.GetOrderBook(ctx, symbol, depth)
	if err != nil {
		g.logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to get order book from MEXC")
		return model.OrderBook{}, fmt.Errorf("failed to get order book from MEXC: %w", err)
	}

	return *orderBook, nil
}

// GetAccountInfo retrieves account information
func (g *Gateway) GetAccountInfo(ctx context.Context) (model.AccountInfo, error) {
	g.logger.Debug().Msg("Getting account info from MEXC")

	// This is a placeholder implementation
	// In a real implementation, you would call the client to get account info
	// For now, we'll return a mock account info
	return model.AccountInfo{
		UserID:      "user123",
		CanTrade:    true,
		CanWithdraw: true,
		CanDeposit:  true,
		Balances:    []model.Balance{},
		LastUpdated: time.Now(),
	}, nil
}

// GetAssetBalance retrieves the balance for a specific asset
func (g *Gateway) GetAssetBalance(ctx context.Context, asset string) (decimal.Decimal, error) {
	g.logger.Debug().Str("asset", asset).Msg("Getting asset balance from MEXC")

	// This is a placeholder implementation
	// In a real implementation, you would call the client to get the asset balance
	// For now, we'll return a mock balance
	return decimal.NewFromFloat(1000.0), nil
}

// PlaceOrder places a new order on MEXC
func (g *Gateway) PlaceOrder(ctx context.Context, order *model.Order) (model.OrderResult, error) {
	g.logger.Debug().Str("symbol", order.Symbol).Str("side", string(order.Side)).Str("type", string(order.Type)).Float64("quantity", order.Quantity).Float64("price", order.Price).Msg("Placing order on MEXC")

	// Place order using client
	result, err := g.client.PlaceOrder(ctx, order.Symbol, order.Side, order.Type, order.Quantity, order.Price, order.TimeInForce)
	if err != nil {
		g.logger.Error().Err(err).Str("symbol", order.Symbol).Msg("Failed to place order on MEXC")
		return model.OrderResult{}, fmt.Errorf("failed to place order on MEXC: %w", err)
	}

	// Convert to OrderResult
	orderResult := model.OrderResult{
		OrderID:       result.OrderID,
		ClientOrderID: result.ClientOrderID,
		Symbol:        result.Symbol,
		TransactTime:  result.CreatedAt,
		Price:         result.Price,
		OrigQty:       result.Quantity,
		ExecutedQty:   result.ExecutedQty,
		Status:        result.Status,
		Type:          result.Type,
		Side:          result.Side,
	}

	return orderResult, nil
}

// CancelOrder cancels an existing order on MEXC
func (g *Gateway) CancelOrder(ctx context.Context, symbol, orderID string) error {
	g.logger.Debug().Str("symbol", symbol).Str("orderID", orderID).Msg("Cancelling order on MEXC")

	// This is a placeholder implementation
	// In a real implementation, you would call the client to cancel the order
	// For now, we'll just return nil (success)
	return nil
}
