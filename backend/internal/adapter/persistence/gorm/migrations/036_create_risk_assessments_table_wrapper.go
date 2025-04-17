package migrations

import (
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// CreateRiskAssessmentsTable runs the migration to create risk assessments table
func CreateRiskAssessmentsTable(db *gorm.DB, logger *zerolog.Logger) error {
	migration := RiskAssessmentMigration{}
	return migration.Migrate(db, logger)
}
