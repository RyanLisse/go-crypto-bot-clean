package migrations

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// AutoBuyRuleEntity is the GORM model for auto-buy rule data
type AutoBuyRuleEntity struct {
	ID              string `gorm:"primaryKey"`
	UserID          string `gorm:"index:idx_auto_buy_rule_user_id"`
	Name            string
	Description     string
	TriggerType     string
	TriggerValue    float64
	Amount          float64
	MaxPrice        float64
	MinVolume       float64
	Enabled         bool
	QuoteAsset      string
	ExcludeSymbols  string
	IncludeSymbols  string
	CooldownMinutes int
	CreatedAt       string
	UpdatedAt       string
}

// TableName sets the table name for AutoBuyRuleEntity
func (AutoBuyRuleEntity) TableName() string {
	return "auto_buy_rules"
}

// CreateAutoBuyRulesTable creates the auto-buy rules table
func CreateAutoBuyRulesTable(db *gorm.DB) error {
	logger := log.With().Str("migration", "create_auto_buy_rules_table").Logger()
	logger.Info().Msg("Running migration: Create auto-buy rules table")

	// Create the auto-buy rules table
	if err := db.AutoMigrate(&AutoBuyRuleEntity{}); err != nil {
		logger.Error().Err(err).Msg("Failed to create auto-buy rules table")
		return err
	}

	logger.Info().Msg("Auto-buy rules table created successfully")
	return nil
}
