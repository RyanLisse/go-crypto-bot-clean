package migrations

import (
	"context"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// CreateEnhancedWalletsTable creates the enhanced wallets table
type CreateEnhancedWalletsTable struct {
	logger *zerolog.Logger
}

// NewCreateEnhancedWalletsTable creates a new migration
func NewCreateEnhancedWalletsTable(logger *zerolog.Logger) *CreateEnhancedWalletsTable {
	return &CreateEnhancedWalletsTable{
		logger: logger,
	}
}

// Name returns the name of the migration
func (m *CreateEnhancedWalletsTable) Name() string {
	return "create_enhanced_wallets_table"
}

// Up runs the migration
func (m *CreateEnhancedWalletsTable) Up(ctx context.Context, db *gorm.DB) error {
	m.logger.Info().Msg("Running migration: Create enhanced wallets table")

	// Create enhanced wallets table
	err := db.Exec(`
		CREATE TABLE IF NOT EXISTS enhanced_wallets (
			id VARCHAR(50) PRIMARY KEY,
			user_id VARCHAR(50) NOT NULL,
			exchange VARCHAR(50),
			type VARCHAR(20) NOT NULL,
			status VARCHAR(20) NOT NULL,
			total_usd_value DECIMAL(18,8) NOT NULL DEFAULT 0,
			metadata JSON,
			last_updated DATETIME NOT NULL,
			last_sync_at DATETIME,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			INDEX idx_enhanced_wallets_user_id (user_id),
			INDEX idx_enhanced_wallets_exchange (exchange),
			INDEX idx_enhanced_wallets_type (type),
			INDEX idx_enhanced_wallets_status (status)
		)
	`).Error

	if err != nil {
		m.logger.Error().Err(err).Msg("Failed to create enhanced wallets table")
		return err
	}

	// Create enhanced wallet balances table
	err = db.Exec(`
		CREATE TABLE IF NOT EXISTS enhanced_wallet_balances (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			wallet_id VARCHAR(50) NOT NULL,
			asset VARCHAR(20) NOT NULL,
			free DECIMAL(18,8) NOT NULL DEFAULT 0,
			locked DECIMAL(18,8) NOT NULL DEFAULT 0,
			total DECIMAL(18,8) NOT NULL DEFAULT 0,
			usd_value DECIMAL(18,8) NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			INDEX idx_enhanced_wallet_balances_wallet_id (wallet_id),
			INDEX idx_enhanced_wallet_balances_asset (asset)
		)
	`).Error

	if err != nil {
		m.logger.Error().Err(err).Msg("Failed to create enhanced wallet balances table")
		return err
	}

	m.logger.Info().Msg("Enhanced wallets tables created successfully")
	return nil
}

// Down rolls back the migration
func (m *CreateEnhancedWalletsTable) Down(ctx context.Context, db *gorm.DB) error {
	m.logger.Info().Msg("Rolling back migration: Create enhanced wallets table")

	// Drop enhanced wallet balances table
	err := db.Exec(`DROP TABLE IF EXISTS enhanced_wallet_balances`).Error
	if err != nil {
		m.logger.Error().Err(err).Msg("Failed to drop enhanced wallet balances table")
		return err
	}

	// Drop enhanced wallets table
	err = db.Exec(`DROP TABLE IF EXISTS enhanced_wallets`).Error
	if err != nil {
		m.logger.Error().Err(err).Msg("Failed to drop enhanced wallets table")
		return err
	}

	m.logger.Info().Msg("Enhanced wallets tables dropped successfully")
	return nil
}
