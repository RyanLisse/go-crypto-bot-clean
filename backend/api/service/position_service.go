package service

import (
	"context"
	"time"

	"go-crypto-bot-clean/backend/api/repository"
)

// OrderEvent represents an order execution event
type OrderEvent struct {
	PositionID string
	Status     string // e.g., "filled", "closed"
	ClosePrice float64
	CloseTime  time.Time
}

// MarketTick represents a market data tick
type MarketTick struct {
	Symbol string
	Price  float64
}

// PositionService handles position updates and tracking
type PositionService struct {
	repo repository.PositionRepository
}

// NewPositionService creates a new PositionService
func NewPositionService(repo repository.PositionRepository) *PositionService {
	return &PositionService{repo: repo}
}

// UpdatePositionPrice updates the current price of a position
func (s *PositionService) UpdatePositionPrice(ctx context.Context, positionID string, newPrice float64) error {
	return s.repo.UpdatePrice(ctx, positionID, newPrice)
}

// ClosePosition marks a position as closed with close price and time
func (s *PositionService) ClosePosition(ctx context.Context, positionID string, closePrice float64, closeTime time.Time) error {
	return s.repo.MarkClosed(ctx, positionID, closePrice, closeTime)
}

// HandleOrderEvent processes an order execution event
func (s *PositionService) HandleOrderEvent(ctx context.Context, event OrderEvent) error {
	if event.Status == "closed" {
		return s.ClosePosition(ctx, event.PositionID, event.ClosePrice, event.CloseTime)
	}
	return nil
}

// HandleMarketTick updates prices of all open positions for the symbol
func (s *PositionService) HandleMarketTick(ctx context.Context, tick MarketTick) error {
	positions, err := s.repo.GetByStatus(ctx, "open")
	if err != nil {
		return err
	}

	for _, pos := range positions {
		if pos.Symbol == tick.Symbol {
			err := s.UpdatePositionPrice(ctx, pos.ID, tick.Price)
			if err != nil {
				// Log error and continue updating other positions
				continue
			}
		}
	}
	return nil
}
