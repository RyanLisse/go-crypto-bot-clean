package repository

import (
	"context"
	"time"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
)

// BalanceHistory represents a point in the balance history
type BalanceHistory struct {
	ID            int64     `db:"id"`
	Timestamp     time.Time `db:"timestamp"`
	Balance       float64   `db:"balance"`
	Equity        float64   `db:"equity"`
	FreeBalance   float64   `db:"free_balance"`
	LockedBalance float64   `db:"locked_balance"`
	UnrealizedPnL float64   `db:"unrealized_pnl"`
}

// BalanceHistoryRepository defines the interface for balance history persistence
type BalanceHistoryRepository interface {
	// Create adds a new balance history point
	Create(ctx context.Context, history *BalanceHistory) (int64, error)
	
	// GetBalanceHistory retrieves balance history within a time range
	GetBalanceHistory(ctx context.Context, startTime, endTime time.Time) ([]*BalanceHistory, error)
	
	// GetLatestBalance retrieves the latest balance history point
	GetLatestBalance(ctx context.Context) (*BalanceHistory, error)
	
	// GetBalancePoints retrieves balance points for equity curve
	GetBalancePoints(ctx context.Context, startTime, endTime time.Time) ([]models.BalancePoint, error)
}
