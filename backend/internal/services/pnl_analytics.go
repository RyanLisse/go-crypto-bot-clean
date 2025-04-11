package services

import (
	"errors"
	"time"

	"go-crypto-bot-clean/backend/api/models"
)

// PriceProvider abstracts current price retrieval
type PriceProvider interface {
	GetPrice(symbol string) (float64, error)
}

// AnalyticsResult holds aggregated analytics metrics
type AnalyticsResult struct {
	TotalRealizedPnL   float64
	TotalUnrealizedPnL float64
	WinCount           int
	LossCount          int
	AverageDuration    time.Duration
	TotalTrades        int
}

// PnLService provides P&L calculation methods
type PnLService interface {
	CalculateUnrealizedPnL(positions []models.Position, priceProvider PriceProvider) (float64, error)
	CalculateRealizedPnL(closedPositions []models.Position) (float64, error)
}

// AnalyticsService provides position analytics
type AnalyticsService interface {
	GetPositionAnalytics(positions []models.Position, priceProvider PriceProvider) (AnalyticsResult, error)
}

// pnlServiceImpl implements PnLService
type pnlServiceImpl struct{}

func NewPnLService() PnLService {
	return &pnlServiceImpl{}
}

func (s *pnlServiceImpl) CalculateUnrealizedPnL(positions []models.Position, priceProvider PriceProvider) (float64, error) {
	var total float64
	for _, pos := range positions {
		if pos.Status != "open" {
			continue
		}
		price, err := priceProvider.GetPrice(pos.Symbol)
		if err != nil {
			return 0, err
		}
		unrealized := (price - pos.EntryPrice) * pos.Quantity
		total += unrealized
	}
	return total, nil
}

func (s *pnlServiceImpl) CalculateRealizedPnL(closedPositions []models.Position) (float64, error) {
	var total float64
	for _, pos := range closedPositions {
		if pos.Status != "closed" {
			continue
		}
		if pos.CloseTime == nil {
			return 0, errors.New("closed position missing CloseTime")
		}
		realized := (pos.CurrentPrice - pos.EntryPrice) * pos.Quantity
		total += realized
	}
	return total, nil
}

// analyticsServiceImpl implements AnalyticsService
type analyticsServiceImpl struct {
	pnlSvc PnLService
}

func NewAnalyticsService(pnlSvc PnLService) AnalyticsService {
	return &analyticsServiceImpl{pnlSvc: pnlSvc}
}

func (s *analyticsServiceImpl) GetPositionAnalytics(positions []models.Position, priceProvider PriceProvider) (AnalyticsResult, error) {
	var result AnalyticsResult
	var totalDuration time.Duration

	// Separate open and closed positions
	var openPositions, closedPositions []models.Position
	for _, pos := range positions {
		if pos.Status == "closed" {
			closedPositions = append(closedPositions, pos)
		} else if pos.Status == "open" {
			openPositions = append(openPositions, pos)
		}
	}

	// Calculate PnL
	unrealizedPnL, err := s.pnlSvc.CalculateUnrealizedPnL(openPositions, priceProvider)
	if err != nil {
		return result, err
	}
	realizedPnL, err := s.pnlSvc.CalculateRealizedPnL(closedPositions)
	if err != nil {
		return result, err
	}
	result.TotalUnrealizedPnL = unrealizedPnL
	result.TotalRealizedPnL = realizedPnL

	// Count wins/losses and compute average duration
	for _, pos := range closedPositions {
		pnl := (pos.CurrentPrice - pos.EntryPrice) * pos.Quantity
		if pnl >= 0 {
			result.WinCount++
		} else {
			result.LossCount++
		}
		if pos.OpenTime != nil && pos.CloseTime != nil {
			duration := pos.CloseTime.Sub(*pos.OpenTime)
			totalDuration += duration
		}
	}

	result.TotalTrades = len(closedPositions)
	if result.TotalTrades > 0 {
		result.AverageDuration = totalDuration / time.Duration(result.TotalTrades)
	}

	return result, nil
}
