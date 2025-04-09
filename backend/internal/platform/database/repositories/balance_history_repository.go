package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/repository"
)

// SQLiteBalanceHistoryRepository implements the BalanceHistoryRepository interface using SQLite
type SQLiteBalanceHistoryRepository struct {
	db *sqlx.DB
}

// NewSQLiteBalanceHistoryRepository creates a new SQLite-based balance history repository
func NewSQLiteBalanceHistoryRepository(db *sqlx.DB) repository.BalanceHistoryRepository {
	return &SQLiteBalanceHistoryRepository{
		db: db,
	}
}

// Create adds a new balance history point
func (r *SQLiteBalanceHistoryRepository) Create(ctx context.Context, history *repository.BalanceHistory) (int64, error) {
	query := `
		INSERT INTO balance_history (
			timestamp, balance, equity, free_balance, locked_balance, unrealized_pnl
		) VALUES (?, ?, ?, ?, ?, ?)
		RETURNING id
	`
	
	var id int64
	err := r.db.QueryRowContext(ctx, query,
		history.Timestamp,
		history.Balance,
		history.Equity,
		history.FreeBalance,
		history.LockedBalance,
		history.UnrealizedPnL,
	).Scan(&id)
	
	if err != nil {
		return 0, fmt.Errorf("failed to insert balance history: %w", err)
	}
	
	history.ID = id
	return id, nil
}

// GetBalanceHistory retrieves balance history within a time range
func (r *SQLiteBalanceHistoryRepository) GetBalanceHistory(ctx context.Context, startTime, endTime time.Time) ([]*repository.BalanceHistory, error) {
	query := `
		SELECT * FROM balance_history
		WHERE timestamp BETWEEN ? AND ?
		ORDER BY timestamp ASC
	`
	
	var history []*repository.BalanceHistory
	err := r.db.SelectContext(ctx, &history, query, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance history: %w", err)
	}
	
	return history, nil
}

// GetLatestBalance retrieves the latest balance history point
func (r *SQLiteBalanceHistoryRepository) GetLatestBalance(ctx context.Context) (*repository.BalanceHistory, error) {
	query := `
		SELECT * FROM balance_history
		ORDER BY timestamp DESC
		LIMIT 1
	`
	
	var history repository.BalanceHistory
	err := r.db.GetContext(ctx, &history, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest balance: %w", err)
	}
	
	return &history, nil
}

// GetBalancePoints retrieves balance points for equity curve
func (r *SQLiteBalanceHistoryRepository) GetBalancePoints(ctx context.Context, startTime, endTime time.Time) ([]models.BalancePoint, error) {
	query := `
		SELECT timestamp, balance FROM balance_history
		WHERE timestamp BETWEEN ? AND ?
		ORDER BY timestamp ASC
	`
	
	type balancePoint struct {
		Timestamp time.Time `db:"timestamp"`
		Balance   float64   `db:"balance"`
	}
	
	var points []balancePoint
	err := r.db.SelectContext(ctx, &points, query, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance points: %w", err)
	}
	
	// Convert to models.BalancePoint
	result := make([]models.BalancePoint, len(points))
	for i, p := range points {
		result[i] = models.BalancePoint{
			Timestamp: p.Timestamp,
			Balance:   p.Balance,
		}
	}
	
	return result, nil
}
