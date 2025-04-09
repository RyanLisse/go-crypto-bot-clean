package analytics

import (
	"context"
	"fmt"
	"math"
	"time"

	"go.uber.org/zap"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/repository"
)

// TradeAnalyticsService provides methods for analyzing trading performance
type TradeAnalyticsService interface {
	// GetTradeAnalytics generates analytics for a specific time frame
	GetTradeAnalytics(ctx context.Context, timeFrame models.TimeFrame, startTime, endTime time.Time) (*models.TradeAnalytics, error)
	
	// GetTradePerformance retrieves performance metrics for a specific trade
	GetTradePerformance(ctx context.Context, tradeID string) (*models.TradePerformance, error)
	
	// GetAllTradePerformance retrieves performance metrics for all trades in a time range
	GetAllTradePerformance(ctx context.Context, startTime, endTime time.Time) ([]*models.TradePerformance, error)
	
	// GetWinRate calculates the win rate for a time range
	GetWinRate(ctx context.Context, startTime, endTime time.Time) (float64, error)
	
	// GetProfitFactor calculates the profit factor for a time range
	GetProfitFactor(ctx context.Context, startTime, endTime time.Time) (float64, error)
	
	// GetDrawdown calculates the maximum drawdown for a time range
	GetDrawdown(ctx context.Context, startTime, endTime time.Time) (float64, float64, error)
	
	// GetBalanceHistory retrieves the balance history for a time range
	GetBalanceHistory(ctx context.Context, startTime, endTime time.Time, interval time.Duration) ([]models.BalancePoint, error)
	
	// GetPerformanceBySymbol retrieves performance metrics grouped by symbol
	GetPerformanceBySymbol(ctx context.Context, startTime, endTime time.Time) (map[string]models.SymbolPerformance, error)
	
	// GetPerformanceByReason retrieves performance metrics grouped by decision reason
	GetPerformanceByReason(ctx context.Context, startTime, endTime time.Time) (map[string]models.ReasonPerformance, error)
	
	// GetPerformanceByStrategy retrieves performance metrics grouped by strategy
	GetPerformanceByStrategy(ctx context.Context, startTime, endTime time.Time) (map[string]models.StrategyPerformance, error)
}

// tradeAnalyticsService implements TradeAnalyticsService
type tradeAnalyticsService struct {
	closedPositionRepo repository.ClosedPositionRepository
	tradeDecisionRepo  repository.TradeDecisionRepository
	transactionRepo    repository.TransactionRepository
	balanceHistoryRepo repository.BalanceHistoryRepository
	logger             *zap.Logger
}

// NewTradeAnalyticsService creates a new trade analytics service
func NewTradeAnalyticsService(
	closedPositionRepo repository.ClosedPositionRepository,
	tradeDecisionRepo repository.TradeDecisionRepository,
	transactionRepo repository.TransactionRepository,
	balanceHistoryRepo repository.BalanceHistoryRepository,
	logger *zap.Logger,
) TradeAnalyticsService {
	return &tradeAnalyticsService{
		closedPositionRepo: closedPositionRepo,
		tradeDecisionRepo:  tradeDecisionRepo,
		transactionRepo:    transactionRepo,
		balanceHistoryRepo: balanceHistoryRepo,
		logger:             logger,
	}
}

// GetTradeAnalytics generates analytics for a specific time frame
func (s *tradeAnalyticsService) GetTradeAnalytics(ctx context.Context, timeFrame models.TimeFrame, startTime, endTime time.Time) (*models.TradeAnalytics, error) {
	// Initialize analytics
	analytics := &models.TradeAnalytics{
		TimeFrame:           timeFrame,
		StartTime:           startTime,
		EndTime:             endTime,
		PerformanceByReason: make(map[string]models.ReasonPerformance),
		PerformanceBySymbol: make(map[string]models.SymbolPerformance),
		PerformanceByStrategy: make(map[string]models.StrategyPerformance),
	}

	// Get closed positions in the time range
	positions, err := s.closedPositionRepo.FindByTimeRange(ctx, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get closed positions: %w", err)
	}

	// Calculate overall performance metrics
	var (
		totalProfit     float64
		totalLoss       float64
		winningTrades   int
		losingTrades    int
		largestProfit   float64
		largestLoss     float64
		totalHoldingMs  int64
		winningHoldingMs int64
		losingHoldingMs int64
	)

	// Performance by symbol and reason
	symbolPerf := make(map[string]*models.SymbolPerformance)
	reasonPerf := make(map[string]*models.ReasonPerformance)
	strategyPerf := make(map[string]*models.StrategyPerformance)

	for _, pos := range positions {
		// Count total trades
		analytics.TotalTrades++

		// Calculate profit/loss
		if pos.ProfitLoss > 0 {
			winningTrades++
			totalProfit += pos.ProfitLoss
			if pos.ProfitLoss > largestProfit {
				largestProfit = pos.ProfitLoss
			}
			winningHoldingMs += pos.HoldingTimeMs
		} else {
			losingTrades++
			totalLoss += math.Abs(pos.ProfitLoss)
			if math.Abs(pos.ProfitLoss) > largestLoss {
				largestLoss = math.Abs(pos.ProfitLoss)
			}
			losingHoldingMs += pos.HoldingTimeMs
		}

		// Track total holding time
		totalHoldingMs += pos.HoldingTimeMs

		// Track performance by symbol
		symbol := pos.Symbol
		if _, exists := symbolPerf[symbol]; !exists {
			symbolPerf[symbol] = &models.SymbolPerformance{
				Symbol: symbol,
			}
		}
		symbolPerf[symbol].TotalTrades++
		if pos.ProfitLoss > 0 {
			symbolPerf[symbol].WinningTrades++
			symbolPerf[symbol].TotalProfit += pos.ProfitLoss
		} else {
			symbolPerf[symbol].LosingTrades++
		}

		// Track performance by reason
		reason := pos.ExitReason
		if _, exists := reasonPerf[reason]; !exists {
			reasonPerf[reason] = &models.ReasonPerformance{
				Reason: reason,
			}
		}
		reasonPerf[reason].TotalTrades++
		if pos.ProfitLoss > 0 {
			reasonPerf[reason].WinningTrades++
			reasonPerf[reason].TotalProfit += pos.ProfitLoss
		} else {
			reasonPerf[reason].LosingTrades++
		}

		// Track performance by strategy
		strategy := pos.Strategy
		if strategy == "" {
			strategy = "unknown"
		}
		if _, exists := strategyPerf[strategy]; !exists {
			strategyPerf[strategy] = &models.StrategyPerformance{
				Strategy: strategy,
			}
		}
		strategyPerf[strategy].TotalTrades++
		if pos.ProfitLoss > 0 {
			strategyPerf[strategy].WinningTrades++
			strategyPerf[strategy].TotalProfit += pos.ProfitLoss
		} else {
			strategyPerf[strategy].LosingTrades++
		}
	}

	// Set overall metrics
	analytics.WinningTrades = winningTrades
	analytics.LosingTrades = losingTrades
	analytics.TotalProfit = totalProfit
	analytics.TotalLoss = totalLoss
	analytics.NetProfit = totalProfit - totalLoss
	analytics.LargestProfit = largestProfit
	analytics.LargestLoss = largestLoss

	// Calculate win rate
	if analytics.TotalTrades > 0 {
		analytics.WinRate = float64(winningTrades) / float64(analytics.TotalTrades)
	}

	// Calculate profit factor
	if totalLoss > 0 {
		analytics.ProfitFactor = totalProfit / totalLoss
	} else if totalProfit > 0 {
		analytics.ProfitFactor = math.Inf(1) // Infinite profit factor (no losses)
	}

	// Calculate average profit/loss
	if winningTrades > 0 {
		analytics.AverageProfit = totalProfit / float64(winningTrades)
	}
	if losingTrades > 0 {
		analytics.AverageLoss = totalLoss / float64(losingTrades)
	}

	// Calculate average holding time
	if analytics.TotalTrades > 0 {
		avgHoldingMs := totalHoldingMs / int64(analytics.TotalTrades)
		analytics.AverageHoldingTime = formatDuration(time.Duration(avgHoldingMs) * time.Millisecond)
	}
	if winningTrades > 0 {
		avgWinningHoldingMs := winningHoldingMs / int64(winningTrades)
		analytics.AverageHoldingTimeWinning = formatDuration(time.Duration(avgWinningHoldingMs) * time.Millisecond)
	}
	if losingTrades > 0 {
		avgLosingHoldingMs := losingHoldingMs / int64(losingTrades)
		analytics.AverageHoldingTimeLosing = formatDuration(time.Duration(avgLosingHoldingMs) * time.Millisecond)
	}

	// Calculate trade frequency
	durationDays := endTime.Sub(startTime).Hours() / 24
	if durationDays > 0 {
		analytics.TradesPerDay = float64(analytics.TotalTrades) / durationDays
		analytics.TradesPerWeek = analytics.TradesPerDay * 7
		analytics.TradesPerMonth = analytics.TradesPerDay * 30
	}

	// Calculate performance by symbol
	for symbol, perf := range symbolPerf {
		if perf.TotalTrades > 0 {
			perf.WinRate = float64(perf.WinningTrades) / float64(perf.TotalTrades)
			perf.AverageProfit = perf.TotalProfit / float64(perf.TotalTrades)
			
			// Calculate profit factor for symbol
			totalLoss := float64(perf.LosingTrades) * math.Abs(perf.AverageProfit)
			if totalLoss > 0 {
				perf.ProfitFactor = perf.TotalProfit / totalLoss
			} else if perf.TotalProfit > 0 {
				perf.ProfitFactor = math.Inf(1)
			}
			
			analytics.PerformanceBySymbol[symbol] = *perf
		}
	}

	// Calculate performance by reason
	for reason, perf := range reasonPerf {
		if perf.TotalTrades > 0 {
			perf.WinRate = float64(perf.WinningTrades) / float64(perf.TotalTrades)
			perf.AverageProfit = perf.TotalProfit / float64(perf.TotalTrades)
			
			// Calculate profit factor for reason
			totalLoss := float64(perf.LosingTrades) * math.Abs(perf.AverageProfit)
			if totalLoss > 0 {
				perf.ProfitFactor = perf.TotalProfit / totalLoss
			} else if perf.TotalProfit > 0 {
				perf.ProfitFactor = math.Inf(1)
			}
			
			analytics.PerformanceByReason[reason] = *perf
		}
	}

	// Calculate performance by strategy
	for strategy, perf := range strategyPerf {
		if perf.TotalTrades > 0 {
			perf.WinRate = float64(perf.WinningTrades) / float64(perf.TotalTrades)
			perf.AverageProfit = perf.TotalProfit / float64(perf.TotalTrades)
			
			// Calculate profit factor for strategy
			totalLoss := float64(perf.LosingTrades) * math.Abs(perf.AverageProfit)
			if totalLoss > 0 {
				perf.ProfitFactor = perf.TotalProfit / totalLoss
			} else if perf.TotalProfit > 0 {
				perf.ProfitFactor = math.Inf(1)
			}
			
			analytics.PerformanceByStrategy[strategy] = *perf
		}
	}

	// Get balance history
	balanceHistory, err := s.GetBalanceHistory(ctx, startTime, endTime, getDurationForTimeFrame(timeFrame))
	if err != nil {
		s.logger.Warn("Failed to get balance history", zap.Error(err))
	} else {
		analytics.BalanceHistory = balanceHistory
	}

	// Calculate drawdown
	maxDrawdown, maxDrawdownPercent, err := s.GetDrawdown(ctx, startTime, endTime)
	if err != nil {
		s.logger.Warn("Failed to calculate drawdown", zap.Error(err))
	} else {
		analytics.MaxDrawdown = maxDrawdown
		analytics.MaxDrawdownPercent = maxDrawdownPercent
	}

	return analytics, nil
}

// GetTradePerformance retrieves performance metrics for a specific trade
func (s *tradeAnalyticsService) GetTradePerformance(ctx context.Context, tradeID string) (*models.TradePerformance, error) {
	// Get the closed position
	position, err := s.closedPositionRepo.FindByID(ctx, tradeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get closed position: %w", err)
	}

	// Create trade performance
	performance := &models.TradePerformance{
		TradeID:         position.ID,
		Symbol:          position.Symbol,
		EntryTime:       position.OpenTime,
		ExitTime:        position.CloseTime,
		EntryPrice:      position.EntryPrice,
		ExitPrice:       position.ExitPrice,
		Quantity:        position.Quantity,
		ProfitLoss:      position.ProfitLoss,
		ProfitLossPercent: position.ProfitLossPercentage,
		HoldingTime:     formatDuration(time.Duration(position.HoldingTimeMs) * time.Millisecond),
		HoldingTimeMs:   position.HoldingTimeMs,
		EntryReason:     position.EntryReason,
		ExitReason:      position.ExitReason,
		Strategy:        position.Strategy,
		StopLoss:        position.InitialStopLoss,
		TakeProfit:      position.InitialTakeProfit,
		RiskRewardRatio: position.RiskRewardRatio,
		ActualRR:        position.ActualRR,
		ExpectedValue:   position.ExpectedValue,
	}

	return performance, nil
}

// GetAllTradePerformance retrieves performance metrics for all trades in a time range
func (s *tradeAnalyticsService) GetAllTradePerformance(ctx context.Context, startTime, endTime time.Time) ([]*models.TradePerformance, error) {
	// Get closed positions in the time range
	positions, err := s.closedPositionRepo.FindByTimeRange(ctx, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get closed positions: %w", err)
	}

	// Create trade performances
	performances := make([]*models.TradePerformance, 0, len(positions))
	for _, pos := range positions {
		performance := &models.TradePerformance{
			TradeID:         pos.ID,
			Symbol:          pos.Symbol,
			EntryTime:       pos.OpenTime,
			ExitTime:        pos.CloseTime,
			EntryPrice:      pos.EntryPrice,
			ExitPrice:       pos.ExitPrice,
			Quantity:        pos.Quantity,
			ProfitLoss:      pos.ProfitLoss,
			ProfitLossPercent: pos.ProfitLossPercentage,
			HoldingTime:     formatDuration(time.Duration(pos.HoldingTimeMs) * time.Millisecond),
			HoldingTimeMs:   pos.HoldingTimeMs,
			EntryReason:     pos.EntryReason,
			ExitReason:      pos.ExitReason,
			Strategy:        pos.Strategy,
			StopLoss:        pos.InitialStopLoss,
			TakeProfit:      pos.InitialTakeProfit,
			RiskRewardRatio: pos.RiskRewardRatio,
			ActualRR:        pos.ActualRR,
			ExpectedValue:   pos.ExpectedValue,
		}
		performances = append(performances, performance)
	}

	return performances, nil
}

// GetWinRate calculates the win rate for a time range
func (s *tradeAnalyticsService) GetWinRate(ctx context.Context, startTime, endTime time.Time) (float64, error) {
	// Get closed positions in the time range
	positions, err := s.closedPositionRepo.FindByTimeRange(ctx, startTime, endTime)
	if err != nil {
		return 0, fmt.Errorf("failed to get closed positions: %w", err)
	}

	// Calculate win rate
	totalTrades := len(positions)
	winningTrades := 0
	for _, pos := range positions {
		if pos.ProfitLoss > 0 {
			winningTrades++
		}
	}

	if totalTrades > 0 {
		return float64(winningTrades) / float64(totalTrades), nil
	}
	return 0, nil
}

// GetProfitFactor calculates the profit factor for a time range
func (s *tradeAnalyticsService) GetProfitFactor(ctx context.Context, startTime, endTime time.Time) (float64, error) {
	// Get closed positions in the time range
	positions, err := s.closedPositionRepo.FindByTimeRange(ctx, startTime, endTime)
	if err != nil {
		return 0, fmt.Errorf("failed to get closed positions: %w", err)
	}

	// Calculate profit factor
	totalProfit := 0.0
	totalLoss := 0.0
	for _, pos := range positions {
		if pos.ProfitLoss > 0 {
			totalProfit += pos.ProfitLoss
		} else {
			totalLoss += math.Abs(pos.ProfitLoss)
		}
	}

	if totalLoss > 0 {
		return totalProfit / totalLoss, nil
	} else if totalProfit > 0 {
		return math.Inf(1), nil // Infinite profit factor (no losses)
	}
	return 0, nil
}

// GetDrawdown calculates the maximum drawdown for a time range
func (s *tradeAnalyticsService) GetDrawdown(ctx context.Context, startTime, endTime time.Time) (float64, float64, error) {
	// Get balance history
	balanceHistory, err := s.GetBalanceHistory(ctx, startTime, endTime, time.Hour)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get balance history: %w", err)
	}

	if len(balanceHistory) < 2 {
		return 0, 0, nil
	}

	// Calculate maximum drawdown
	maxBalance := balanceHistory[0].Balance
	maxDrawdown := 0.0
	maxDrawdownPercent := 0.0

	for _, point := range balanceHistory {
		if point.Balance > maxBalance {
			maxBalance = point.Balance
		}

		drawdown := maxBalance - point.Balance
		drawdownPercent := drawdown / maxBalance

		if drawdown > maxDrawdown {
			maxDrawdown = drawdown
			maxDrawdownPercent = drawdownPercent
		}
	}

	return maxDrawdown, maxDrawdownPercent, nil
}

// GetBalanceHistory retrieves the balance history for a time range
func (s *tradeAnalyticsService) GetBalanceHistory(ctx context.Context, startTime, endTime time.Time, interval time.Duration) ([]models.BalancePoint, error) {
	// Get balance history from repository
	balanceHistory, err := s.balanceHistoryRepo.GetBalanceHistory(ctx, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance history: %w", err)
	}

	// If no balance history, try to reconstruct from transactions
	if len(balanceHistory) == 0 {
		return s.reconstructBalanceHistory(ctx, startTime, endTime, interval)
	}

	// Convert to BalancePoint
	points := make([]models.BalancePoint, 0, len(balanceHistory))
	for _, bh := range balanceHistory {
		points = append(points, models.BalancePoint{
			Timestamp: bh.Timestamp,
			Balance:   bh.Balance,
		})
	}

	return points, nil
}

// reconstructBalanceHistory reconstructs balance history from transactions
func (s *tradeAnalyticsService) reconstructBalanceHistory(ctx context.Context, startTime, endTime time.Time, interval time.Duration) ([]models.BalancePoint, error) {
	// Get transactions
	transactions, err := s.transactionRepo.FindByTimeRange(ctx, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	if len(transactions) == 0 {
		return nil, nil
	}

	// Create balance points at regular intervals
	points := make([]models.BalancePoint, 0)
	
	// Start with the first transaction's balance
	currentTime := startTime
	currentBalance := transactions[0].Balance
	
	// Add initial point
	points = append(points, models.BalancePoint{
		Timestamp: currentTime,
		Balance:   currentBalance,
	})
	
	// Add points at regular intervals
	for currentTime.Before(endTime) {
		currentTime = currentTime.Add(interval)
		
		// Find the latest transaction before currentTime
		for _, tx := range transactions {
			if tx.Timestamp.Before(currentTime) || tx.Timestamp.Equal(currentTime) {
				currentBalance = tx.Balance
			} else {
				break
			}
		}
		
		// Add point
		points = append(points, models.BalancePoint{
			Timestamp: currentTime,
			Balance:   currentBalance,
		})
	}
	
	return points, nil
}

// GetPerformanceBySymbol retrieves performance metrics grouped by symbol
func (s *tradeAnalyticsService) GetPerformanceBySymbol(ctx context.Context, startTime, endTime time.Time) (map[string]models.SymbolPerformance, error) {
	// Get analytics
	analytics, err := s.GetTradeAnalytics(ctx, models.TimeFrameAll, startTime, endTime)
	if err != nil {
		return nil, err
	}
	
	return analytics.PerformanceBySymbol, nil
}

// GetPerformanceByReason retrieves performance metrics grouped by decision reason
func (s *tradeAnalyticsService) GetPerformanceByReason(ctx context.Context, startTime, endTime time.Time) (map[string]models.ReasonPerformance, error) {
	// Get analytics
	analytics, err := s.GetTradeAnalytics(ctx, models.TimeFrameAll, startTime, endTime)
	if err != nil {
		return nil, err
	}
	
	return analytics.PerformanceByReason, nil
}

// GetPerformanceByStrategy retrieves performance metrics grouped by strategy
func (s *tradeAnalyticsService) GetPerformanceByStrategy(ctx context.Context, startTime, endTime time.Time) (map[string]models.StrategyPerformance, error) {
	// Get analytics
	analytics, err := s.GetTradeAnalytics(ctx, models.TimeFrameAll, startTime, endTime)
	if err != nil {
		return nil, err
	}
	
	return analytics.PerformanceByStrategy, nil
}

// Helper functions

// formatDuration formats a duration as a human-readable string
func formatDuration(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	
	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	} else if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}

// getDurationForTimeFrame returns an appropriate interval duration for a time frame
func getDurationForTimeFrame(timeFrame models.TimeFrame) time.Duration {
	switch timeFrame {
	case models.TimeFrameDay:
		return time.Hour
	case models.TimeFrameWeek:
		return 6 * time.Hour
	case models.TimeFrameMonth:
		return 24 * time.Hour
	case models.TimeFrameQuarter:
		return 24 * time.Hour
	case models.TimeFrameYear:
		return 7 * 24 * time.Hour
	default:
		return 24 * time.Hour
	}
}
