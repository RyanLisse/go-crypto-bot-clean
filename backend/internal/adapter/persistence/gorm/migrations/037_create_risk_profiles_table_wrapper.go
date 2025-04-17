package migrations

import (
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// CreateRiskProfilesTable runs the migration to create risk profiles table
func CreateRiskProfilesTable(db *gorm.DB, logger *zerolog.Logger) error {
	migration := RiskProfileMigration{}
	return migration.Migrate(db, logger)
}
