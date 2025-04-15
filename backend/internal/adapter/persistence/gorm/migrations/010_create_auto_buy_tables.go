package migrations

import (
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm/entity"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// CreateAutoBuyTables creates tables for auto-buy functionality
func CreateAutoBuyTables(db *gorm.DB) error {
	logger := log.With().Str("migration", "create_auto_buy_tables").Logger()
	logger.Info().Msg("Running migration: Create auto-buy tables")

	// Create auto_buy_rules table
	if err := db.AutoMigrate(&entity.AutoBuyRuleEntity{}); err != nil {
		logger.Error().Err(err).Msg("Failed to create auto_buy_rules table")
		return err
	}

	// Create auto_buy_executions table
	if err := db.AutoMigrate(&entity.AutoBuyExecutionEntity{}); err != nil {
		logger.Error().Err(err).Msg("Failed to create auto_buy_executions table")
		return err
	}

	// Add foreign key constraint from auto_buy_executions to auto_buy_rules
	// Skip for SQLite as it has issues with ALTER TABLE ADD CONSTRAINT
	if db.Dialector.Name() != "sqlite" {
		err := db.Exec(
			"ALTER TABLE auto_buy_executions " +
				"ADD CONSTRAINT fk_auto_buy_executions_rule " +
				"FOREIGN KEY (rule_id) REFERENCES auto_buy_rules(id) ON DELETE CASCADE",
		).Error
		if err != nil {
			logger.Error().Err(err).Msg("Failed to add foreign key constraint to auto_buy_executions table")
			return err
		}
		logger.Info().Msg("Added foreign key constraint to auto_buy_executions table")
	} else {
		logger.Info().Msg("Skipping foreign key constraint for SQLite")
	}

	// Add index on timestamp in auto_buy_executions
	indexErr := db.Exec("CREATE INDEX IF NOT EXISTS idx_auto_buy_executions_timestamp ON auto_buy_executions(timestamp)").Error
	if indexErr != nil {
		logger.Error().Err(indexErr).Msg("Failed to create timestamp index on auto_buy_executions table")
		return indexErr
	}

	logger.Info().Msg("Successfully created auto-buy tables")
	return nil
}
