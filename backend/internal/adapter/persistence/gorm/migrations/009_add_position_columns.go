package migrations

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// AddPositionColumns adds additional columns to the positions table
func AddPositionColumns(db *gorm.DB) error {
	logger := log.With().Str("migration", "add_position_columns").Logger()
	logger.Info().Msg("Running migration: Add position columns")

	// SQLite doesn't support ADD COLUMN IF NOT EXISTS, so we need to check if the columns exist
	// and add them one by one if they don't

	// Check if max_drawdown column exists
	var hasMaxDrawdown int
	db.Raw("SELECT COUNT(*) FROM pragma_table_info('positions') WHERE name = 'max_drawdown'").Scan(&hasMaxDrawdown)
	if hasMaxDrawdown == 0 {
		if err := db.Exec("ALTER TABLE positions ADD COLUMN max_drawdown DECIMAL(24,8)").Error; err != nil {
			logger.Error().Err(err).Msg("Failed to add max_drawdown column")
			return err
		}
	}

	// Check if max_profit column exists
	var hasMaxProfit int
	db.Raw("SELECT COUNT(*) FROM pragma_table_info('positions') WHERE name = 'max_profit'").Scan(&hasMaxProfit)
	if hasMaxProfit == 0 {
		if err := db.Exec("ALTER TABLE positions ADD COLUMN max_profit DECIMAL(24,8)").Error; err != nil {
			logger.Error().Err(err).Msg("Failed to add max_profit column")
			return err
		}
	}

	// Check if risk_reward_ratio column exists
	var hasRiskRewardRatio int
	db.Raw("SELECT COUNT(*) FROM pragma_table_info('positions') WHERE name = 'risk_reward_ratio'").Scan(&hasRiskRewardRatio)
	if hasRiskRewardRatio == 0 {
		if err := db.Exec("ALTER TABLE positions ADD COLUMN risk_reward_ratio DECIMAL(24,8)").Error; err != nil {
			logger.Error().Err(err).Msg("Failed to add risk_reward_ratio column")
			return err
		}
	}

	// Check if entry_order_ids column exists
	var hasEntryOrderIds int
	db.Raw("SELECT COUNT(*) FROM pragma_table_info('positions') WHERE name = 'entry_order_ids'").Scan(&hasEntryOrderIds)
	if hasEntryOrderIds == 0 {
		if err := db.Exec("ALTER TABLE positions ADD COLUMN entry_order_ids TEXT").Error; err != nil {
			logger.Error().Err(err).Msg("Failed to add entry_order_ids column")
			return err
		}
	}

	// Check if exit_order_ids column exists
	var hasExitOrderIds int
	db.Raw("SELECT COUNT(*) FROM pragma_table_info('positions') WHERE name = 'exit_order_ids'").Scan(&hasExitOrderIds)
	if hasExitOrderIds == 0 {
		if err := db.Exec("ALTER TABLE positions ADD COLUMN exit_order_ids TEXT").Error; err != nil {
			logger.Error().Err(err).Msg("Failed to add exit_order_ids column")
			return err
		}
	}

	logger.Info().Msg("Position columns added successfully")
	return nil
}
