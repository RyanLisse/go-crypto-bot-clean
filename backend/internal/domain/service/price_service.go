package service

import (
	"context"
	"time"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
)

// PriceService defines the interface for price data
type PriceService interface {
	// GetPrice returns the current price for a symbol
	GetPrice(ctx context.Context, symbol string) (float64, error)

	// GetTicker returns the current ticker for a symbol
	GetTicker(ctx context.Context, symbol string) (*models.Ticker, error)

	// GetKlines returns historical kline data
	GetKlines(ctx context.Context, symbol, interval string, limit int) ([]*models.Kline, error)

	// GetPriceHistory returns historical price data
	GetPriceHistory(ctx context.Context, symbol string, startTime, endTime time.Time) ([]float64, error)
}
