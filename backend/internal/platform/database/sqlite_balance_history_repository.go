package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/risk/controls"
)

// SQLiteBalanceHistoryRepository implements the BalanceHistoryRepository interface
type SQLiteBalanceHistoryRepository struct {
	db *sql.DB
}

// NewSQLiteBalanceHistoryRepository creates a new SQLiteBalanceHistoryRepository
func NewSQLiteBalanceHistoryRepository(db *sql.DB) *SQLiteBalanceHistoryRepository {
	return &SQLiteBalanceHistoryRepository{
		db: db,
	}
}

// AddBalanceRecord adds a new balance record
func (r *SQLiteBalanceHistoryRepository) AddBalanceRecord(ctx context.Context, balance float64) error {
	query := `
		INSERT INTO balance_history (balance, timestamp)
		VALUES (?, ?)
	`

	_, err := r.db.ExecContext(ctx, query, balance, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("failed to insert balance record: %w", err)
	}

	return nil
}

// GetHistory retrieves balance history for the specified number of days
func (r *SQLiteBalanceHistoryRepository) GetHistory(ctx context.Context, days int) ([]controls.BalanceHistory, error) {
	query := `
		SELECT id, balance, timestamp
		FROM balance_history
		WHERE timestamp >= datetime('now', '-' || ? || ' days')
		ORDER BY timestamp ASC
	`

	rows, err := r.db.QueryContext(ctx, query, days)
	if err != nil {
		return nil, fmt.Errorf("failed to query balance history: %w", err)
	}
	defer rows.Close()

	var history []controls.BalanceHistory
	for rows.Next() {
		var record controls.BalanceHistory
		var timestamp string

		err := rows.Scan(&record.ID, &record.Balance, &timestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to scan balance record: %w", err)
		}

		// Parse timestamp
		record.Timestamp, err = time.Parse(time.RFC3339, timestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to parse timestamp: %w", err)
		}

		history = append(history, record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating balance records: %w", err)
	}

	return history, nil
}

// GetHighestBalance returns the highest recorded balance
func (r *SQLiteBalanceHistoryRepository) GetHighestBalance(ctx context.Context) (float64, error) {
	query := `
		SELECT MAX(balance)
		FROM balance_history
	`

	var highestBalance float64
	err := r.db.QueryRowContext(ctx, query).Scan(&highestBalance)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to get highest balance: %w", err)
	}

	return highestBalance, nil
}

// GetBalanceAt returns the balance at a specific time
func (r *SQLiteBalanceHistoryRepository) GetBalanceAt(ctx context.Context, timestamp time.Time) (float64, error) {
	query := `
		SELECT balance
		FROM balance_history
		WHERE timestamp <= ?
		ORDER BY timestamp DESC
		LIMIT 1
	`

	var balance float64
	err := r.db.QueryRowContext(ctx, query, timestamp.UTC().Format(time.RFC3339)).Scan(&balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to get balance at %s: %w", timestamp, err)
	}

	return balance, nil
}

// GetLatestBalance returns the most recent balance record
func (r *SQLiteBalanceHistoryRepository) GetLatestBalance(ctx context.Context) (*controls.BalanceHistory, error) {
	query := `
		SELECT id, balance, timestamp
		FROM balance_history
		ORDER BY timestamp DESC
		LIMIT 1
	`

	var record controls.BalanceHistory
	var timestamp string

	err := r.db.QueryRowContext(ctx, query).Scan(&record.ID, &record.Balance, &timestamp)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get latest balance: %w", err)
	}

	// Parse timestamp
	record.Timestamp, err = time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse timestamp: %w", err)
	}

	return &record, nil
}
