package migrations

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// WalletEntity is the GORM model for wallet data
type WalletEntity struct {
	ID          string `gorm:"primaryKey"`
	UserID      string `gorm:"index:idx_wallet_user_id"`
	Exchange    string
	LastUpdated string
	CreatedAt   string
	UpdatedAt   string
}

// TableName sets the table name for WalletEntity
func (WalletEntity) TableName() string {
	return "wallets"
}

// BalanceEntity is the GORM model for balance data
type BalanceEntity struct {
	ID        string `gorm:"primaryKey"`
	WalletID  string `gorm:"index:idx_balance_wallet_id"`
	Asset     string
	Free      float64
	Locked    float64
	CreatedAt string
	UpdatedAt string
}

// TableName sets the table name for BalanceEntity
func (BalanceEntity) TableName() string {
	return "balances"
}

// CreateWalletTable creates the wallet and balance tables
func CreateWalletTable(db *gorm.DB) error {
	logger := log.With().Str("migration", "create_wallet_table").Logger()
	logger.Info().Msg("Running migration: Create wallet table")

	// Create the wallet table
	if err := db.AutoMigrate(&WalletEntity{}); err != nil {
		logger.Error().Err(err).Msg("Failed to create wallet table")
		return err
	}

	// Create the balance table
	if err := db.AutoMigrate(&BalanceEntity{}); err != nil {
		logger.Error().Err(err).Msg("Failed to create balance table")
		return err
	}

	logger.Info().Msg("Wallet and balance tables created successfully")
	return nil
}
