package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upCreateWalletEntitiesTable, downCreateWalletEntitiesTable)
}

func upCreateWalletEntitiesTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS wallet_entities (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id VARCHAR(50) NOT NULL,
			exchange VARCHAR(20) NOT NULL,
			last_updated TIMESTAMP NOT NULL,
			total_usd_value DECIMAL(18,8) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);

		CREATE INDEX IF NOT EXISTS idx_wallet_entities_user_id ON wallet_entities(user_id);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_wallet_entities_user_exchange ON wallet_entities(user_id, exchange);
	`)
	return err
}

func downCreateWalletEntitiesTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
		DROP TABLE IF EXISTS wallet_entities;
	`)
	return err
}
