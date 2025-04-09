package controls

import (
	"context"
	"fmt"
	"time"
)

// BalanceHistory tracks account balance over time for drawdown calculation
type BalanceHistory struct {
	ID        int64     `json:"id"`
	Balance   float64   `json:"balance"`
	Timestamp time.Time `json:"timestamp"`
}

// BalanceHistoryRepository defines the interface for storing and retrieving balance history
type BalanceHistoryRepository interface {
	// AddBalanceRecord adds a new balance record
	AddBalanceRecord(ctx context.Context, balance float64) error

	// GetHistory retrieves balance history for the specified number of days
	GetHistory(ctx context.Context, days int) ([]BalanceHistory, error)

	// GetHighestBalance returns the highest recorded balance
	GetHighestBalance(ctx context.Context) (float64, error)

	// GetBalanceAt returns the balance at a specific time
	GetBalanceAt(ctx context.Context, timestamp time.Time) (float64, error)

	// GetLatestBalance returns the most recent balance record
	GetLatestBalance(ctx context.Context) (*BalanceHistory, error)
}

// DrawdownMonitor calculates and monitors drawdown
type DrawdownMonitor struct {
	balanceRepo BalanceHistoryRepository
	logger      Logger
}

// NewDrawdownMonitor creates a new DrawdownMonitor
func NewDrawdownMonitor(balanceRepo BalanceHistoryRepository, logger Logger) *DrawdownMonitor {
	return &DrawdownMonitor{
		balanceRepo: balanceRepo,
		logger:      logger,
	}
}

// CalculateDrawdown computes the maximum peak-to-trough drawdown
func (dm *DrawdownMonitor) CalculateDrawdown(ctx context.Context, days int) (float64, error) {
	// Get historical balance data
	history, err := dm.balanceRepo.GetHistory(ctx, days)
	if err != nil {
		return 0, fmt.Errorf("failed to get balance history: %w", err)
	}

	if len(history) < 2 {
		return 0, nil // Not enough data to calculate drawdown
	}

	// Find peak and calculate drawdown
	var maxDrawdown float64
	var peak float64

	for _, entry := range history {
		if entry.Balance > peak {
			peak = entry.Balance
		}

		if peak > 0 {
			drawdown := (peak - entry.Balance) / peak
			if drawdown > maxDrawdown {
				maxDrawdown = drawdown
			}
		}
	}

	dm.logger.Info("Calculated drawdown",
		"days", days,
		"max_drawdown", maxDrawdown,
		"peak_balance", peak)

	return maxDrawdown, nil
}

// CheckDrawdownLimit verifies if trading should be allowed based on drawdown
func (dm *DrawdownMonitor) CheckDrawdownLimit(ctx context.Context, maxDrawdownPercent float64) (bool, error) {
	// Calculate drawdown over the last 90 days
	drawdown, err := dm.CalculateDrawdown(ctx, 90)
	if err != nil {
		return false, err
	}

	maxAllowed := maxDrawdownPercent / 100
	allowed := drawdown < maxAllowed

	if !allowed {
		dm.logger.Warn("Trading disabled due to drawdown limit",
			"current_drawdown", drawdown,
			"max_allowed", maxAllowed)
	}

	return allowed, nil
}

// RecordBalance adds a new balance record to the repository
func (dm *DrawdownMonitor) RecordBalance(ctx context.Context, balance float64) error {
	err := dm.balanceRepo.AddBalanceRecord(ctx, balance)
	if err != nil {
		return fmt.Errorf("failed to record balance: %w", err)
	}

	dm.logger.Info("Recorded balance", "balance", balance)
	return nil
}
