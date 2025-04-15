package usecase

import (
	"context"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/market"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
)

// Ensure MarketDataUseCase implements the MarketDataUseCaseInterface
var _ port.MarketDataUseCaseInterface = (*MarketDataUseCase)(nil)

// GetSymbols implements the MarketDataUseCaseInterface
func (uc *MarketDataUseCase) GetSymbols(ctx context.Context) ([]*market.Symbol, error) {
	return uc.symbolRepo.GetAll(ctx)
}

// GetSymbol implements the MarketDataUseCaseInterface
func (uc *MarketDataUseCase) GetSymbol(ctx context.Context, symbol string) (*market.Symbol, error) {
	return uc.GetSymbolInfo(ctx, symbol)
}

// GetOrderBook implements the MarketDataUseCaseInterface
func (uc *MarketDataUseCase) GetOrderBook(ctx context.Context, exchange, symbol string) (*market.OrderBook, error) {
	// This is a mock implementation for now
	return &market.OrderBook{
		Symbol:      symbol,
		Exchange:    exchange,
		LastUpdated: time.Now(),
		Bids: []market.OrderBookEntry{
			{Price: 40000.0, Quantity: 1.5},
			{Price: 39900.0, Quantity: 2.0},
			{Price: 39800.0, Quantity: 3.0},
		},
		Asks: []market.OrderBookEntry{
			{Price: 40100.0, Quantity: 1.0},
			{Price: 40200.0, Quantity: 2.5},
			{Price: 40300.0, Quantity: 1.8},
		},
	}, nil
}

// GetAllTickers implements the MarketDataUseCaseInterface
func (uc *MarketDataUseCase) GetAllTickers(ctx context.Context, exchange string) ([]*market.Ticker, error) {
	return uc.marketRepo.GetAllTickers(ctx, exchange)
}
