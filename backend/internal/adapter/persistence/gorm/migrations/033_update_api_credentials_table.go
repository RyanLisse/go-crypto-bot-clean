package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upUpdateAPICredentialsTable033, downUpdateAPICredentialsTable033)
}

func upUpdateAPICredentialsTable033(tx *sql.Tx) error {
	_, err := tx.Exec(`
		-- Check if the table exists
		CREATE TABLE IF NOT EXISTS api_credentials (
			id VARCHAR(50) PRIMARY KEY,
			user_id VARCHAR(50) NOT NULL,
			exchange VARCHAR(20) NOT NULL,
			api_key VARCHAR(100) NOT NULL,
			api_secret BLOB NOT NULL,
			label VARCHAR(50),
			status VARCHAR(20) NOT NULL DEFAULT 'active',
			last_used TIMESTAMP NULL,
			last_verified TIMESTAMP NULL,
			expires_at TIMESTAMP NULL,
			rotation_due TIMESTAMP NULL,
			failure_count INTEGER NOT NULL DEFAULT 0,
			metadata JSON NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);

		-- Create indexes if they don't exist
		CREATE INDEX IF NOT EXISTS idx_api_credentials_user_id ON api_credentials(user_id);
		CREATE INDEX IF NOT EXISTS idx_api_credentials_exchange ON api_credentials(exchange);
		CREATE INDEX IF NOT EXISTS idx_api_credentials_status ON api_credentials(status);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_api_credentials_user_exchange_label ON api_credentials(user_id, exchange, label);

		-- Add new columns to existing table if they don't exist
		-- SQLite doesn't support ADD COLUMN IF NOT EXISTS, so we need to check if the columns exist
		
		-- Check if status column exists
		SELECT CASE 
			WHEN COUNT(*) = 0 THEN
				-- Add status column if it doesn't exist
				ALTER TABLE api_credentials ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT 'active'
			ELSE
				-- Do nothing
				SELECT 1
		END
		FROM pragma_table_info('api_credentials') WHERE name = 'status';

		-- Check if last_used column exists
		SELECT CASE 
			WHEN COUNT(*) = 0 THEN
				-- Add last_used column if it doesn't exist
				ALTER TABLE api_credentials ADD COLUMN last_used TIMESTAMP NULL
			ELSE
				-- Do nothing
				SELECT 1
		END
		FROM pragma_table_info('api_credentials') WHERE name = 'last_used';

		-- Check if last_verified column exists
		SELECT CASE 
			WHEN COUNT(*) = 0 THEN
				-- Add last_verified column if it doesn't exist
				ALTER TABLE api_credentials ADD COLUMN last_verified TIMESTAMP NULL
			ELSE
				-- Do nothing
				SELECT 1
		END
		FROM pragma_table_info('api_credentials') WHERE name = 'last_verified';

		-- Check if expires_at column exists
		SELECT CASE 
			WHEN COUNT(*) = 0 THEN
				-- Add expires_at column if it doesn't exist
				ALTER TABLE api_credentials ADD COLUMN expires_at TIMESTAMP NULL
			ELSE
				-- Do nothing
				SELECT 1
		END
		FROM pragma_table_info('api_credentials') WHERE name = 'expires_at';

		-- Check if rotation_due column exists
		SELECT CASE 
			WHEN COUNT(*) = 0 THEN
				-- Add rotation_due column if it doesn't exist
				ALTER TABLE api_credentials ADD COLUMN rotation_due TIMESTAMP NULL
			ELSE
				-- Do nothing
				SELECT 1
		END
		FROM pragma_table_info('api_credentials') WHERE name = 'rotation_due';

		-- Check if failure_count column exists
		SELECT CASE 
			WHEN COUNT(*) = 0 THEN
				-- Add failure_count column if it doesn't exist
				ALTER TABLE api_credentials ADD COLUMN failure_count INTEGER NOT NULL DEFAULT 0
			ELSE
				-- Do nothing
				SELECT 1
		END
		FROM pragma_table_info('api_credentials') WHERE name = 'failure_count';

		-- Check if metadata column exists
		SELECT CASE 
			WHEN COUNT(*) = 0 THEN
				-- Add metadata column if it doesn't exist
				ALTER TABLE api_credentials ADD COLUMN metadata JSON NULL
			ELSE
				-- Do nothing
				SELECT 1
		END
		FROM pragma_table_info('api_credentials') WHERE name = 'metadata';
	`)
	return err
}

func downUpdateAPICredentialsTable033(tx *sql.Tx) error {
	// We don't want to drop the table or remove columns in the down migration
	// as it could lead to data loss. Instead, we'll just do nothing.
	return nil
}
