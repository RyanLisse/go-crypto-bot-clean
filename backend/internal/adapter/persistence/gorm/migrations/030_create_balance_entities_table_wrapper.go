package migrations

import (
	"gorm.io/gorm"
)

// CreateBalanceEntitiesTable creates the balance entities table
func CreateBalanceEntitiesTable(db *gorm.DB) error {
	return db.Exec(`
		CREATE TABLE IF NOT EXISTS balance_entities (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			wallet_id INTEGER NOT NULL,
			asset VARCHAR(20) NOT NULL,
			free DECIMAL(18,8) NOT NULL,
			locked DECIMAL(18,8) NOT NULL,
			total DECIMAL(18,8) NOT NULL,
			usd_value DECIMAL(18,8) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (wallet_id) REFERENCES wallet_entities(id) ON DELETE CASCADE
		);

		CREATE INDEX IF NOT EXISTS idx_balance_entities_wallet_id ON balance_entities(wallet_id);
		CREATE INDEX IF NOT EXISTS idx_balance_entities_asset ON balance_entities(asset);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_balance_entities_wallet_asset ON balance_entities(wallet_id, asset);
	`).Error
}
