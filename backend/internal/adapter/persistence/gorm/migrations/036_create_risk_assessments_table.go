package migrations

import (
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type RiskAssessmentMigration struct{}

func (m RiskAssessmentMigration) ID() string {
	return "036_create_risk_assessments_table"
}

func (m RiskAssessmentMigration) Migrate(db *gorm.DB, logger *zerolog.Logger) error {
	logger.Info().Msg("Running migration: " + m.ID())

	// Create risk assessments table
	err := db.Exec(`
		CREATE TABLE IF NOT EXISTS risk_assessments (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			type TEXT NOT NULL, 
			level TEXT NOT NULL,
			status TEXT NOT NULL,
			symbol TEXT,
			position_id TEXT,
			order_id TEXT,
			score REAL,
			message TEXT NOT NULL,
			recommendation TEXT,
			metadata_json TEXT,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			resolved_at TIMESTAMP
		)
	`).Error
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create risk_assessments table")
		return err
	}

	// Create indices for faster searches
	err = db.Exec("CREATE INDEX IF NOT EXISTS idx_risk_assessments_user_id ON risk_assessments(user_id)").Error
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create user_id index on risk_assessments")
		return err
	}

	err = db.Exec("CREATE INDEX IF NOT EXISTS idx_risk_assessments_type ON risk_assessments(type)").Error
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create type index on risk_assessments")
		return err
	}

	err = db.Exec("CREATE INDEX IF NOT EXISTS idx_risk_assessments_level ON risk_assessments(level)").Error
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create level index on risk_assessments")
		return err
	}

	err = db.Exec("CREATE INDEX IF NOT EXISTS idx_risk_assessments_status ON risk_assessments(status)").Error
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create status index on risk_assessments")
		return err
	}

	err = db.Exec("CREATE INDEX IF NOT EXISTS idx_risk_assessments_symbol ON risk_assessments(symbol)").Error
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create symbol index on risk_assessments")
		return err
	}

	err = db.Exec("CREATE INDEX IF NOT EXISTS idx_risk_assessments_created_at ON risk_assessments(created_at)").Error
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create created_at index on risk_assessments")
		return err
	}

	logger.Info().Msg("Migration completed: " + m.ID())
	return nil
}
