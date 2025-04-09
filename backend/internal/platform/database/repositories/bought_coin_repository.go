package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
	"go-crypto-bot-clean/backend/internal/domain/repository"
	"github.com/jmoiron/sqlx"
)

// SQLiteBoughtCoinRepository implements the BoughtCoinRepository interface using SQLite
type SQLiteBoughtCoinRepository struct {
	db *sqlx.DB
}

// NewSQLiteBoughtCoinRepository creates a new SQLite implementation of BoughtCoinRepository
func NewSQLiteBoughtCoinRepository(db *sqlx.DB) repository.BoughtCoinRepository {
	return &SQLiteBoughtCoinRepository{
		db: db,
	}
}

// FindAll returns all bought coins that haven't been deleted
func (r *SQLiteBoughtCoinRepository) FindAll(ctx context.Context) ([]models.BoughtCoin, error) {
	var coins []models.BoughtCoin
	query := `SELECT * FROM bought_coins WHERE is_deleted = 0`
	err := r.db.SelectContext(ctx, &coins, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find all bought coins: %w", err)
	}
	return coins, nil
}

// FindByID returns a specific bought coin by ID
func (r *SQLiteBoughtCoinRepository) FindByID(ctx context.Context, id int64) (*models.BoughtCoin, error) {
	var coin models.BoughtCoin
	query := `SELECT * FROM bought_coins WHERE id = ? AND is_deleted = 0`
	err := r.db.GetContext(ctx, &coin, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find bought coin by ID: %w", err)
	}
	return &coin, nil
}

// FindBySymbol returns a specific bought coin by symbol
func (r *SQLiteBoughtCoinRepository) FindBySymbol(ctx context.Context, symbol string) (*models.BoughtCoin, error) {
	var coin models.BoughtCoin
	query := `SELECT * FROM bought_coins WHERE symbol = ? AND is_deleted = 0`
	err := r.db.GetContext(ctx, &coin, query, symbol)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find bought coin by symbol: %w", err)
	}
	return &coin, nil
}

// Create adds a new bought coin
func (r *SQLiteBoughtCoinRepository) Create(ctx context.Context, coin *models.BoughtCoin) (int64, error) {
	query := `
		INSERT INTO bought_coins (
			symbol, purchase_price, quantity, bought_at,
			stop_loss, take_profit, current_price, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.ExecContext(
		ctx, query,
		coin.Symbol, coin.PurchasePrice, coin.Quantity, coin.BoughtAt,
		coin.StopLoss, coin.TakeProfit, coin.CurrentPrice, coin.UpdatedAt,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create bought coin: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	return id, nil
}

// Update modifies an existing bought coin
func (r *SQLiteBoughtCoinRepository) Update(ctx context.Context, coin *models.BoughtCoin) error {
	query := `
		UPDATE bought_coins
		SET symbol = ?, purchase_price = ?, quantity = ?, bought_at = ?,
			stop_loss = ?, take_profit = ?, current_price = ?, updated_at = ?
		WHERE id = ?
	`
	_, err := r.db.ExecContext(
		ctx, query,
		coin.Symbol, coin.PurchasePrice, coin.Quantity, coin.BoughtAt,
		coin.StopLoss, coin.TakeProfit, coin.CurrentPrice, time.Now(),
		coin.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update bought coin: %w", err)
	}

	return nil
}

// Delete marks a bought coin as deleted
func (r *SQLiteBoughtCoinRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE bought_coins SET is_deleted = 1, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete bought coin: %w", err)
	}

	return nil
}
