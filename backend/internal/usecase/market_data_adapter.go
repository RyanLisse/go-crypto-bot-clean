package usecase

import (
	"context"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// MarketDataUseCase implements the port.MarketDataUseCaseInterface
type MarketDataUseCase struct {
	marketDataService port.MarketDataService
	logger            *zerolog.Logger
}

// Ensure MarketDataUseCase implements the port.MarketDataUseCaseInterface
var _ port.MarketDataUseCaseInterface = (*MarketDataUseCase)(nil)

// NewMarketDataUseCase creates a new market data use case
func NewMarketDataUseCase(marketDataService port.MarketDataService, logger *zerolog.Logger) port.MarketDataUseCaseInterface {
	return &MarketDataUseCase{
		marketDataService: marketDataService,
		logger:            logger,
	}
}

// GetCandles implements the MarketDataUseCaseInterface
func (uc *MarketDataUseCase) GetCandles(ctx context.Context, exchange, symbol string, interval model.KlineInterval, start, end time.Time, limit int) ([]*model.Kline, error) {
	// If start and end time are provided, use historical prices
	if !start.IsZero() && !end.IsZero() {
		return uc.marketDataService.GetHistoricalPrices(ctx, symbol, start, end, interval)
	}
	// Otherwise use regular candles
	return uc.marketDataService.GetCandles(ctx, symbol, interval, limit)
}

// GetTicker retrieves the current ticker for a symbol
func (uc *MarketDataUseCase) GetTicker(ctx context.Context, exchange, symbol string) (*model.Ticker, error) {
	return uc.marketDataService.GetTicker(ctx, symbol)
}

// GetOrderBook retrieves the order book for a symbol with the specified depth
func (uc *MarketDataUseCase) GetOrderBook(ctx context.Context, exchange, symbol string) (*model.OrderBook, error) {
	// Default depth if not specified
	const defaultDepth = 50
	return uc.marketDataService.GetOrderBook(ctx, symbol, defaultDepth)
}

// GetAllSymbols retrieves information about all available symbols
func (uc *MarketDataUseCase) GetSymbols(ctx context.Context) ([]*model.Symbol, error) {
	return uc.marketDataService.GetAllSymbols(ctx)
}

// GetSymbolInfo retrieves detailed information about a specific symbol
func (uc *MarketDataUseCase) GetSymbol(ctx context.Context, symbol string) (*model.Symbol, error) {
	return uc.marketDataService.GetSymbolInfo(ctx, symbol)
}

// GetAllTickers implements the MarketDataUseCaseInterface
func (uc *MarketDataUseCase) GetAllTickers(ctx context.Context, exchange string) ([]*model.Ticker, error) {
	// Since the marketDataService doesn't have a direct GetAllTickers method,
	// we'll need to get all symbols and then get tickers for each one
	symbols, err := uc.marketDataService.GetAllSymbols(ctx)
	if err != nil {
		return nil, err
	}

	tickers := make([]*model.Ticker, 0, len(symbols))
	for _, symbol := range symbols {
		ticker, err := uc.marketDataService.GetTicker(ctx, symbol.Symbol)
		if err != nil {
			uc.logger.Warn().Err(err).Str("symbol", symbol.Symbol).Msg("Failed to get ticker for symbol")
			continue
		}
		tickers = append(tickers, ticker)
	}

	return tickers, nil
}

// GetLatestTickers implements the MarketDataUseCaseInterface
func (uc *MarketDataUseCase) GetLatestTickers(ctx context.Context) ([]*model.Ticker, error) {
	const defaultExchange = "mexc"
	return uc.GetAllTickers(ctx, defaultExchange)
}
