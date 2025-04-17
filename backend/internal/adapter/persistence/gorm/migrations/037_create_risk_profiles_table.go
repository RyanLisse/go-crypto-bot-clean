package migrations

import (
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type RiskProfileMigration struct{}

func (m RiskProfileMigration) ID() string {
	return "037_create_risk_profiles_table"
}

func (m RiskProfileMigration) Migrate(db *gorm.DB, logger *zerolog.Logger) error {
	logger.Info().Msg("Running migration: " + m.ID())

	// Create risk profiles table
	err := db.Exec(`
		CREATE TABLE IF NOT EXISTS risk_profiles (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL UNIQUE,
			max_position_size REAL NOT NULL,
			max_total_exposure REAL NOT NULL,
			max_drawdown REAL NOT NULL,
			max_leverage REAL NOT NULL,
			max_concentration REAL NOT NULL,
			min_liquidity REAL NOT NULL,
			volatility_threshold REAL NOT NULL,
			daily_loss_limit REAL NOT NULL,
			weekly_loss_limit REAL NOT NULL,
			enable_auto_risk_control BOOLEAN NOT NULL,
			enable_notifications BOOLEAN NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`).Error
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create risk_profiles table")
		return err
	}

	// Create index for user_id
	err = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_risk_profiles_user_id ON risk_profiles(user_id)").Error
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create user_id index on risk_profiles")
		return err
	}

	logger.Info().Msg("Migration completed: " + m.ID())
	return nil
}
