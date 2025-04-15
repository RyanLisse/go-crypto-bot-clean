package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upConsolidateWalletTables, downConsolidateWalletTables)
}

func upConsolidateWalletTables(tx *sql.Tx) error {
	// This migration ensures that the enhanced wallet tables are properly created
	// It's a safety measure in case the previous migration (031_create_enhanced_wallet_table.go)
	// was not executed or had issues

	// Create enhanced wallet table if it doesn't exist
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS enhanced_wallets (
			id VARCHAR(50) PRIMARY KEY,
			user_id VARCHAR(50) NOT NULL,
			exchange VARCHAR(50),
			type VARCHAR(20) NOT NULL,
			status VARCHAR(20) NOT NULL,
			total_usd_value DECIMAL(18,8) NOT NULL DEFAULT 0,
			metadata JSON,
			last_updated TIMESTAMP NOT NULL,
			last_sync_at TIMESTAMP,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);

		CREATE INDEX IF NOT EXISTS idx_enhanced_wallets_user_id ON enhanced_wallets(user_id);
		CREATE INDEX IF NOT EXISTS idx_enhanced_wallets_type ON enhanced_wallets(type);
		CREATE INDEX IF NOT EXISTS idx_enhanced_wallets_status ON enhanced_wallets(status);
	`)
	if err != nil {
		return err
	}

	// Create enhanced wallet balances table if it doesn't exist
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS enhanced_wallet_balances (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			wallet_id VARCHAR(50) NOT NULL,
			asset VARCHAR(20) NOT NULL,
			free DECIMAL(18,8) NOT NULL DEFAULT 0,
			locked DECIMAL(18,8) NOT NULL DEFAULT 0,
			total DECIMAL(18,8) NOT NULL DEFAULT 0,
			usd_value DECIMAL(18,8) NOT NULL DEFAULT 0,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (wallet_id) REFERENCES enhanced_wallets(id) ON DELETE CASCADE
		);

		CREATE INDEX IF NOT EXISTS idx_enhanced_wallet_balances_wallet_id ON enhanced_wallet_balances(wallet_id);
		CREATE INDEX IF NOT EXISTS idx_enhanced_wallet_balances_asset ON enhanced_wallet_balances(asset);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_enhanced_wallet_balances_wallet_id_asset ON enhanced_wallet_balances(wallet_id, asset);
	`)
	if err != nil {
		return err
	}

	// Create enhanced wallet balance history table if it doesn't exist
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS enhanced_wallet_balance_history (
			id VARCHAR(50) PRIMARY KEY,
			user_id VARCHAR(50) NOT NULL,
			wallet_id VARCHAR(50) NOT NULL,
			asset VARCHAR(20) NOT NULL,
			free DECIMAL(18,8) NOT NULL DEFAULT 0,
			locked DECIMAL(18,8) NOT NULL DEFAULT 0,
			total DECIMAL(18,8) NOT NULL DEFAULT 0,
			usd_value DECIMAL(18,8) NOT NULL DEFAULT 0,
			timestamp TIMESTAMP NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (wallet_id) REFERENCES enhanced_wallets(id) ON DELETE CASCADE
		);

		CREATE INDEX IF NOT EXISTS idx_enhanced_wallet_balance_history_user_id ON enhanced_wallet_balance_history(user_id);
		CREATE INDEX IF NOT EXISTS idx_enhanced_wallet_balance_history_wallet_id ON enhanced_wallet_balance_history(wallet_id);
		CREATE INDEX IF NOT EXISTS idx_enhanced_wallet_balance_history_asset ON enhanced_wallet_balance_history(asset);
		CREATE INDEX IF NOT EXISTS idx_enhanced_wallet_balance_history_timestamp ON enhanced_wallet_balance_history(timestamp);
	`)

	return err
}

func downConsolidateWalletTables(tx *sql.Tx) error {
	// This is a no-op since we don't want to drop the tables
	// The tables will be dropped by the down migration of 031_create_enhanced_wallet_table.go
	return nil
}
