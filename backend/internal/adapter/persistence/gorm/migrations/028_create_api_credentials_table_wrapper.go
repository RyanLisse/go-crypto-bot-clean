package migrations

import (
	"gorm.io/gorm"
)

// CreateAPICredentialsTable creates the API credentials table
func CreateAPICredentialsTable(db *gorm.DB) error {
	return db.Exec(`
		CREATE TABLE IF NOT EXISTS api_credentials (
			id VARCHAR(50) PRIMARY KEY,
			user_id VARCHAR(50) NOT NULL,
			exchange VARCHAR(20) NOT NULL,
			api_key VARCHAR(100) NOT NULL,
			api_secret BLOB NOT NULL,
			label VARCHAR(50),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);

		CREATE INDEX IF NOT EXISTS idx_api_credentials_user_id ON api_credentials(user_id);
		CREATE INDEX IF NOT EXISTS idx_api_credentials_exchange ON api_credentials(exchange);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_api_credentials_user_exchange_label ON api_credentials(user_id, exchange, label);
	`).Error
}
