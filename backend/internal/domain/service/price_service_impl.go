package service

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/interfaces"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
)

// PriceServiceImpl implements the interfaces.PriceService interface
type PriceServiceImpl struct {
	exchangeService interfaces.ExchangeService
	logger          *zap.Logger
}

// NewPriceService creates a new interfaces.PriceService
func NewPriceService(exchangeService interfaces.ExchangeService, logger *zap.Logger) interfaces.PriceService {
	return &PriceServiceImpl{
		exchangeService: exchangeService,
		logger:          logger,
	}
}

// GetPrice returns the current price for a symbol
func (s *PriceServiceImpl) GetPrice(ctx context.Context, symbol string) (float64, error) {
	ticker, err := s.exchangeService.GetTicker(ctx, symbol)
	if err != nil {
		return 0, fmt.Errorf("failed to get ticker: %w", err)
	}
	return ticker.Price, nil
}

// GetTicker returns the current ticker for a symbol
func (s *PriceServiceImpl) GetTicker(ctx context.Context, symbol string) (*models.Ticker, error) {
	return s.exchangeService.GetTicker(ctx, symbol)
}

// GetKlines returns historical kline data
func (s *PriceServiceImpl) GetKlines(ctx context.Context, symbol, interval string, limit int) ([]*models.Kline, error) {
	return s.exchangeService.GetKlines(ctx, symbol, interval, limit)
}

// GetPriceHistory returns historical price data
func (s *PriceServiceImpl) GetPriceHistory(ctx context.Context, symbol string, startTime, endTime time.Time) ([]float64, error) {
	// Get klines for the specified time range
	interval := "1h" // Default interval
	klines, err := s.exchangeService.GetKlines(ctx, symbol, interval, 1000)
	if err != nil {
		return nil, fmt.Errorf("failed to get klines: %w", err)
	}

	// Filter klines by time range and extract close prices
	prices := make([]float64, 0, len(klines))
	for _, kline := range klines {
		if (startTime.IsZero() || kline.OpenTime.After(startTime) || kline.OpenTime.Equal(startTime)) &&
			(endTime.IsZero() || kline.OpenTime.Before(endTime) || kline.OpenTime.Equal(endTime)) {
			prices = append(prices, kline.Close)
		}
	}

	return prices, nil
}
