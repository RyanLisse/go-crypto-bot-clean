package migrations

import (
	"context"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// CreateEnhancedWalletBalanceHistoryTable creates the enhanced wallet balance history table
type CreateEnhancedWalletBalanceHistoryTable struct {
	logger *zerolog.Logger
}

// NewCreateEnhancedWalletBalanceHistoryTable creates a new migration
func NewCreateEnhancedWalletBalanceHistoryTable(logger *zerolog.Logger) *CreateEnhancedWalletBalanceHistoryTable {
	return &CreateEnhancedWalletBalanceHistoryTable{
		logger: logger,
	}
}

// Name returns the name of the migration
func (m *CreateEnhancedWalletBalanceHistoryTable) Name() string {
	return "create_enhanced_wallet_balance_history_table"
}

// Up runs the migration
func (m *CreateEnhancedWalletBalanceHistoryTable) Up(ctx context.Context, db *gorm.DB) error {
	m.logger.Info().Msg("Running migration: Create enhanced wallet balance history table")

	// Create enhanced wallet balance history table
	err := db.Exec(`
		CREATE TABLE IF NOT EXISTS enhanced_wallet_balance_history (
			id VARCHAR(50) PRIMARY KEY,
			user_id VARCHAR(50) NOT NULL,
			wallet_id VARCHAR(50) NOT NULL,
			balances_json JSON,
			total_usd_value DECIMAL(18,8) NOT NULL DEFAULT 0,
			timestamp DATETIME NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			INDEX idx_enhanced_wallet_balance_history_user_id (user_id),
			INDEX idx_enhanced_wallet_balance_history_wallet_id (wallet_id),
			INDEX idx_enhanced_wallet_balance_history_timestamp (timestamp)
		)
	`).Error

	if err != nil {
		m.logger.Error().Err(err).Msg("Failed to create enhanced wallet balance history table")
		return err
	}

	m.logger.Info().Msg("Enhanced wallet balance history table created successfully")
	return nil
}

// Down rolls back the migration
func (m *CreateEnhancedWalletBalanceHistoryTable) Down(ctx context.Context, db *gorm.DB) error {
	m.logger.Info().Msg("Rolling back migration: Create enhanced wallet balance history table")

	// Drop enhanced wallet balance history table
	err := db.Exec(`DROP TABLE IF EXISTS enhanced_wallet_balance_history`).Error
	if err != nil {
		m.logger.Error().Err(err).Msg("Failed to drop enhanced wallet balance history table")
		return err
	}

	m.logger.Info().Msg("Enhanced wallet balance history table dropped successfully")
	return nil
}
