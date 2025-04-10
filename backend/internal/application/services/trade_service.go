package services

import (
	"context"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
	"go-crypto-bot-clean/backend/internal/domain/ports"

	"github.com/google/uuid"
)

// TradeService handles business logic related to trades
type TradeService struct {
	tradeRepository ports.TradeRepository
}

// NewTradeService creates a new TradeService
func NewTradeService(tradeRepository ports.TradeRepository) *TradeService {
	return &TradeService{
		tradeRepository: tradeRepository,
	}
}

// GetTradeByID retrieves a trade by its ID
func (s *TradeService) GetTradeByID(ctx context.Context, id string) (*models.Trade, error) {
	return s.tradeRepository.GetByID(ctx, id)
}

// GetTradesBySymbol retrieves trades by symbol
func (s *TradeService) GetTradesBySymbol(ctx context.Context, symbol string, limit int) ([]*models.Trade, error) {
	return s.tradeRepository.GetBySymbol(ctx, symbol, limit)
}

// GetTradesByOrderID retrieves trades by order ID
func (s *TradeService) GetTradesByOrderID(ctx context.Context, orderID string) ([]*models.Trade, error) {
	return s.tradeRepository.GetByOrderID(ctx, orderID)
}

// RecordTrade records a new trade
func (s *TradeService) RecordTrade(ctx context.Context, orderID, symbol string, isBuyer bool, quantity, price, fee float64) (*models.Trade, error) {
	trade := &models.Trade{
		ID:        uuid.New().String(),
		OrderID:   orderID,
		Symbol:    symbol,
		IsBuyer:   isBuyer,
		Quantity:  quantity,
		Price:     price,
		Fee:       fee,
		TradeTime: time.Now(),
	}

	if err := s.tradeRepository.Store(ctx, trade); err != nil {
		return nil, err
	}

	return trade, nil
}
