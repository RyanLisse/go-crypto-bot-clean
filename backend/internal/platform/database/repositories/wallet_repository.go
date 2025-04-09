package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"

	"go-crypto-bot-clean/backend/internal/domain/models"
	"go-crypto-bot-clean/backend/internal/domain/repository"
)

// SQLiteWalletRepository implements the WalletRepository interface for SQLite
type SQLiteWalletRepository struct {
	db *sqlx.DB
}

// NewSQLiteWalletRepository creates a new SQLite wallet repository
func NewSQLiteWalletRepository(db *sqlx.DB) repository.WalletRepository {
	return &SQLiteWalletRepository{
		db: db,
	}
}

// GetWallet retrieves the wallet from the database
func (r *SQLiteWalletRepository) GetWallet(ctx context.Context) (*models.Wallet, error) {
	// In a real implementation, we would retrieve the wallet data from the database
	// For now, we'll return a mock wallet

	// Query for wallet in the database
	var wallet models.Wallet
	var updatedAt time.Time

	err := r.db.QueryRowContext(ctx, "SELECT updated_at FROM wallets LIMIT 1").Scan(&updatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	// Query for balances
	rows, err := r.db.QueryContext(ctx, "SELECT asset, free, locked FROM asset_balances")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	balances := make(map[string]*models.AssetBalance)
	for rows.Next() {
		var asset string
		var free, locked float64

		if err := rows.Scan(&asset, &free, &locked); err != nil {
			return nil, err
		}

		balances[asset] = &models.AssetBalance{
			Asset:  asset,
			Free:   free,
			Locked: locked,
			Total:  free + locked,
		}
	}

	wallet.Balances = balances
	wallet.UpdatedAt = updatedAt

	return &wallet, nil
}

// SaveWallet saves the wallet to the database
func (r *SQLiteWalletRepository) SaveWallet(ctx context.Context, wallet *models.Wallet) (*models.Wallet, error) {
	// In a real implementation, we would save the wallet data to the database
	// For now, we'll just return the wallet

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Update wallet timestamp
	now := time.Now()
	_, err = tx.ExecContext(ctx, "INSERT OR REPLACE INTO wallets (id, updated_at) VALUES (1, ?)", now)
	if err != nil {
		return nil, err
	}

	// Update asset balances
	for asset, balance := range wallet.Balances {
		_, err = tx.ExecContext(ctx,
			"INSERT OR REPLACE INTO asset_balances (asset, free, locked, updated_at) VALUES (?, ?, ?, ?)",
			asset, balance.Free, balance.Locked, now)
		if err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	// Update the wallet's timestamp
	wallet.UpdatedAt = now

	return wallet, nil
}
