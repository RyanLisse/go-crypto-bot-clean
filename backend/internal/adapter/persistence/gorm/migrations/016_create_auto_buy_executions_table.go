package migrations

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// AutoBuyExecutionEntity is the GORM model for auto-buy execution data
type AutoBuyExecutionEntity struct {
	ID           string `gorm:"primaryKey"`
	RuleID       string `gorm:"index:idx_auto_buy_execution_rule_id"`
	Symbol       string
	Price        float64
	Quantity     float64
	Amount       float64
	OrderID      string
	Status       string
	ErrorMessage string
	ExecutedAt   string
	CreatedAt    string
	UpdatedAt    string
}

// TableName sets the table name for AutoBuyExecutionEntity
func (AutoBuyExecutionEntity) TableName() string {
	return "auto_buy_executions"
}

// CreateAutoBuyExecutionsTable creates the auto-buy executions table
func CreateAutoBuyExecutionsTable(db *gorm.DB) error {
	logger := log.With().Str("migration", "create_auto_buy_executions_table").Logger()
	logger.Info().Msg("Running migration: Create auto-buy executions table")

	// Create the auto-buy executions table
	if err := db.AutoMigrate(&AutoBuyExecutionEntity{}); err != nil {
		logger.Error().Err(err).Msg("Failed to create auto-buy executions table")
		return err
	}

	logger.Info().Msg("Auto-buy executions table created successfully")
	return nil
}
