package risk

import (
	"context"
	"time"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/risk/controls"
)

// BalanceHistoryRepository defines the interface for storing and retrieving balance history
type BalanceHistoryRepository interface {
	// AddBalanceRecord adds a new balance record
	AddBalanceRecord(ctx context.Context, balance float64) error

	// GetHistory retrieves balance history for the specified number of days
	GetHistory(ctx context.Context, days int) ([]controls.BalanceHistory, error)

	// GetHighestBalance returns the highest recorded balance
	GetHighestBalance(ctx context.Context) (float64, error)

	// GetBalanceAt returns the balance at a specific time
	GetBalanceAt(ctx context.Context, timestamp time.Time) (float64, error)

	// GetLatestBalance returns the most recent balance record
	GetLatestBalance(ctx context.Context) (*controls.BalanceHistory, error)
}
