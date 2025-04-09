package balance

import (
	"context"
	"fmt"
	"time"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/repository"
	"github.com/ryanlisse/go-crypto-bot/internal/repository/database"
)

// BalanceHistoryRepository implements the repository.BalanceHistoryRepository interface
// using our database abstraction layer
type BalanceHistoryRepository struct {
	db database.Repository
}

// NewBalanceHistoryRepository creates a new balance history repository
func NewBalanceHistoryRepository(db database.Repository) repository.BalanceHistoryRepository {
	return &BalanceHistoryRepository{
		db: db,
	}
}

// Create adds a new balance history point
func (r *BalanceHistoryRepository) Create(ctx context.Context, history *repository.BalanceHistory) (int64, error) {
	query := `
		INSERT INTO balance_history (
			timestamp, balance, equity, free_balance, locked_balance, unrealized_pnl
		) VALUES (?, ?, ?, ?, ?, ?)
		RETURNING id
	`
	
	var id int64
	err := r.db.QueryRow(ctx, query,
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
func (r *BalanceHistoryRepository) GetBalanceHistory(ctx context.Context, startTime, endTime time.Time) ([]*repository.BalanceHistory, error) {
	query := `
		SELECT id, timestamp, balance, equity, free_balance, locked_balance, unrealized_pnl 
		FROM balance_history
		WHERE timestamp BETWEEN ? AND ?
		ORDER BY timestamp ASC
	`
	
	rows, err := r.db.Query(ctx, query, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance history: %w", err)
	}
	defer rows.Close()
	
	var history []*repository.BalanceHistory
	for rows.Next() {
		var h repository.BalanceHistory
		if err := rows.Scan(
			&h.ID,
			&h.Timestamp,
			&h.Balance,
			&h.Equity,
			&h.FreeBalance,
			&h.LockedBalance,
			&h.UnrealizedPnL,
		); err != nil {
			return nil, fmt.Errorf("failed to scan balance history: %w", err)
		}
		history = append(history, &h)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating balance history: %w", err)
	}
	
	return history, nil
}

// GetLatestBalance retrieves the latest balance history point
func (r *BalanceHistoryRepository) GetLatestBalance(ctx context.Context) (*repository.BalanceHistory, error) {
	query := `
		SELECT id, timestamp, balance, equity, free_balance, locked_balance, unrealized_pnl
		FROM balance_history
		ORDER BY timestamp DESC
		LIMIT 1
	`
	
	var h repository.BalanceHistory
	err := r.db.QueryRow(ctx, query).Scan(
		&h.ID,
		&h.Timestamp,
		&h.Balance,
		&h.Equity,
		&h.FreeBalance,
		&h.LockedBalance,
		&h.UnrealizedPnL,
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to get latest balance: %w", err)
	}
	
	return &h, nil
}

// GetBalancePoints retrieves balance points for equity curve
func (r *BalanceHistoryRepository) GetBalancePoints(ctx context.Context, startTime, endTime time.Time) ([]models.BalancePoint, error) {
	query := `
		SELECT timestamp, balance FROM balance_history
		WHERE timestamp BETWEEN ? AND ?
		ORDER BY timestamp ASC
	`
	
	rows, err := r.db.Query(ctx, query, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance points: %w", err)
	}
	defer rows.Close()
	
	var points []models.BalancePoint
	for rows.Next() {
		var p models.BalancePoint
		if err := rows.Scan(&p.Timestamp, &p.Balance); err != nil {
			return nil, fmt.Errorf("failed to scan balance point: %w", err)
		}
		points = append(points, p)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating balance points: %w", err)
	}
	
	return points, nil
}
