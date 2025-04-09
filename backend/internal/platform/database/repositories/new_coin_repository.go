package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/repository"
)

// SQLiteNewCoinRepository implements the NewCoinRepository interface using SQLite
type SQLiteNewCoinRepository struct {
	db *sqlx.DB
}

// NewSQLiteNewCoinRepository creates a new SQLite implementation of NewCoinRepository
func NewSQLiteNewCoinRepository(db *sqlx.DB) repository.NewCoinRepository {
	return &SQLiteNewCoinRepository{
		db: db,
	}
}

// FindAll returns all new coins that haven't been processed
func (r *SQLiteNewCoinRepository) FindAll(ctx context.Context) ([]models.NewCoin, error) {
	var coins []models.NewCoin
	query := `SELECT * FROM new_coins WHERE is_deleted = 0`
	err := r.db.SelectContext(ctx, &coins, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find all new coins: %w", err)
	}
	return coins, nil
}

// FindByID returns a specific new coin by ID
func (r *SQLiteNewCoinRepository) FindByID(ctx context.Context, id int64) (*models.NewCoin, error) {
	var coin models.NewCoin
	query := `SELECT * FROM new_coins WHERE id = ? AND is_deleted = 0`
	err := r.db.GetContext(ctx, &coin, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find new coin by ID: %w", err)
	}
	return &coin, nil
}

// FindBySymbol returns a specific new coin by symbol
func (r *SQLiteNewCoinRepository) FindBySymbol(ctx context.Context, symbol string) (*models.NewCoin, error) {
	var coin models.NewCoin
	query := `SELECT * FROM new_coins WHERE symbol = ? AND is_deleted = 0`
	err := r.db.GetContext(ctx, &coin, query, symbol)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find new coin by symbol: %w", err)
	}
	return &coin, nil
}

// Create adds a new coin listing
func (r *SQLiteNewCoinRepository) Create(ctx context.Context, coin *models.NewCoin) (int64, error) {
	query := `
		INSERT INTO new_coins (
			symbol, found_at, first_open_time, base_volume, quote_volume, status, became_tradable_at, is_processed, is_deleted, is_upcoming
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.ExecContext(
		ctx, query,
		coin.Symbol, coin.FoundAt, coin.FirstOpenTime, coin.BaseVolume, coin.QuoteVolume,
		coin.Status, coin.BecameTradableAt, coin.IsProcessed, coin.IsDeleted, coin.IsUpcoming,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create new coin: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	return id, nil
}

// MarkAsProcessed marks a new coin as processed
func (r *SQLiteNewCoinRepository) MarkAsProcessed(ctx context.Context, id int64) error {
	query := `UPDATE new_coins SET is_processed = 1 WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to mark new coin as processed: %w", err)
	}

	return nil
}

// Delete marks a new coin as deleted
func (r *SQLiteNewCoinRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE new_coins SET is_deleted = 1 WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete new coin: %w", err)
	}

	return nil
}

// FindByDateRange finds coins within a date range
func (r *SQLiteNewCoinRepository) FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]models.NewCoin, error) {
	var coins []models.NewCoin
	query := `SELECT * FROM new_coins WHERE found_at BETWEEN ? AND ? AND is_deleted = 0`
	err := r.db.SelectContext(ctx, &coins, query, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to find coins by date range: %w", err)
	}
	return coins, nil
}

// FindUpcomingCoins finds coins that are scheduled to be listed in the future
func (r *SQLiteNewCoinRepository) FindUpcomingCoins(ctx context.Context) ([]models.NewCoin, error) {
	var coins []models.NewCoin
	currentTime := time.Now()
	query := `SELECT * FROM new_coins WHERE first_open_time > ? AND is_upcoming = 1 AND is_deleted = 0 ORDER BY first_open_time ASC`
	err := r.db.SelectContext(ctx, &coins, query, currentTime)
	if err != nil {
		return nil, fmt.Errorf("failed to find upcoming coins: %w", err)
	}
	return coins, nil
}

// FindUpcomingCoinsByDate finds upcoming coins that will be listed on a specific date
func (r *SQLiteNewCoinRepository) FindUpcomingCoinsByDate(ctx context.Context, date time.Time) ([]models.NewCoin, error) {
	var coins []models.NewCoin
	// Get the end of the day
	endOfDay := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, date.Location())

	query := `SELECT * FROM new_coins WHERE first_open_time BETWEEN ? AND ? AND is_upcoming = 1 AND is_deleted = 0 ORDER BY first_open_time ASC`
	err := r.db.SelectContext(ctx, &coins, query, date, endOfDay)
	if err != nil {
		return nil, fmt.Errorf("failed to find upcoming coins by date: %w", err)
	}
	return coins, nil
}

// Update updates an existing coin
func (r *SQLiteNewCoinRepository) Update(ctx context.Context, coin *models.NewCoin) error {
	query := `
		UPDATE new_coins SET
			symbol = ?,
			found_at = ?,
			first_open_time = ?,
			base_volume = ?,
			quote_volume = ?,
			status = ?,
			became_tradable_at = ?,
			is_processed = ?,
			is_deleted = ?,
			is_upcoming = ?
		WHERE id = ?
	`
	_, err := r.db.ExecContext(
		ctx, query,
		coin.Symbol, coin.FoundAt, coin.FirstOpenTime, coin.BaseVolume, coin.QuoteVolume,
		coin.Status, coin.BecameTradableAt, coin.IsProcessed, coin.IsDeleted, coin.IsUpcoming,
		coin.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update coin: %w", err)
	}

	return nil
}

// FindTradableCoins finds coins that have become tradable (status = "1")
func (r *SQLiteNewCoinRepository) FindTradableCoins(ctx context.Context) ([]models.NewCoin, error) {
	var coins []models.NewCoin
	query := `SELECT * FROM new_coins WHERE status = '1' AND became_tradable_at IS NOT NULL AND is_deleted = 0 ORDER BY became_tradable_at DESC`
	err := r.db.SelectContext(ctx, &coins, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find tradable coins: %w", err)
	}
	return coins, nil
}

// FindTradableCoinsByDate finds coins that became tradable on a specific date
func (r *SQLiteNewCoinRepository) FindTradableCoinsByDate(ctx context.Context, date time.Time) ([]models.NewCoin, error) {
	var coins []models.NewCoin
	// Set time boundaries for the day
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour).Add(-time.Nanosecond)

	query := `SELECT * FROM new_coins WHERE status = '1' AND became_tradable_at BETWEEN ? AND ? AND is_deleted = 0 ORDER BY became_tradable_at DESC`
	err := r.db.SelectContext(ctx, &coins, query, startOfDay, endOfDay)
	if err != nil {
		return nil, fmt.Errorf("failed to find tradable coins by date: %w", err)
	}
	return coins, nil
}
