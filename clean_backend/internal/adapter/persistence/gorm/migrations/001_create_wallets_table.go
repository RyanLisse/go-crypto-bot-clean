package migrations

import (
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/adapter/persistence/gorm/entity"
)

// CreateWalletsTable creates the wallets table
func CreateWalletsTable(db *gorm.DB, logger *zerolog.Logger) error {
	migrationName := "create_wallets_table"
	log := logger.With().Str("migration", migrationName).Logger()
	log.Info().Msg("Running migration")

	// AutoMigrate will create the table, including columns, primary key, and GORM-defined indexes
	if err := db.AutoMigrate(&entity.WalletEntity{}); err != nil {
		log.Error().Err(err).Msg("Failed to migrate wallets table")
		return err
	}

	// Optional: Add specific indexes not automatically created by GORM tags if needed
	// Example: Creating the specific index for primary wallets if GORM didn't handle the combined index tag correctly
	// if db.Dialector.Name() != "sqlite" { // SQLite might have issues with complex WHERE clauses in index creation
	// 	err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_wallets_user_primary ON wallets(user_id, is_primary) WHERE is_primary = TRUE`).Error
	// 	if err != nil {
	// 		log.Warn().Err(err).Msg("Could not create specific primary wallet index (may be harmless)")
	// 	}
	// }

	log.Info().Msg("Migration completed successfully")
	return nil
}
