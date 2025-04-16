package database

import (
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/rs/zerolog"
	gormdb "gorm.io/gorm"
)

// Connect creates a new database connection
func Connect(cfg *config.Config, logger *zerolog.Logger) (*gormdb.DB, error) {
	// Use the existing GORM connection function
	return gorm.NewDBConnection(cfg, *logger)
}

// RunMigrations runs all database migrations
func RunMigrations(db *gormdb.DB, logger *zerolog.Logger) error {
	// Use the consolidated migrations approach
	return gorm.AutoMigrateModels(db, logger)
}
