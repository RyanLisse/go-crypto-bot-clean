package port

import (
	"context"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
)

// MarketRepository defines methods for market data persistence operations
type MarketRepository interface {
	// Ticker operations
	SaveTicker(ctx context.Context, ticker *model.Ticker) error
	GetTicker(ctx context.Context, symbol, exchange string) (*model.Ticker, error)
	GetAllTickers(ctx context.Context, exchange string) ([]*model.Ticker, error)
	GetLatestTickers(ctx context.Context, limit int) ([]*model.Ticker, error)

	// Kline/Candle operations
	SaveKline(ctx context.Context, kline *model.Kline) error
	SaveKlines(ctx context.Context, klines []*model.Kline) error
	GetKline(ctx context.Context, symbol, exchange string, interval model.KlineInterval, openTime time.Time) (*model.Kline, error)
	GetKlines(ctx context.Context, symbol, exchange string, interval model.KlineInterval, start, end time.Time, limit int) ([]*model.Kline, error)
	GetLatestKline(ctx context.Context, symbol, exchange string, interval model.KlineInterval) (*model.Kline, error)

	// Data management
	PurgeOldData(ctx context.Context, olderThan time.Time) error
}
