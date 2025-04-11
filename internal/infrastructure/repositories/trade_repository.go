package infrastructure

import (
	"context"
	"database/sql"
	"time"
)

type TradeRepository struct {
	db *sql.DB
}

// DeleteOlderThan implements ports.TradeRepository
func (r *TradeRepository) DeleteOlderThan(ctx context.Context, timestamp time.Time) error {
	query := `DELETE FROM trades WHERE created_at < $1`
	_, err := r.db.ExecContext(ctx, query, timestamp)
	return err
}
