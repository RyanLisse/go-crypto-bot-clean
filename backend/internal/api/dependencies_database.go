package api

import (
	"log"

	"github.com/ryanlisse/go-crypto-bot/internal/platform/database/gorm"
	gormrepo "github.com/ryanlisse/go-crypto-bot/internal/platform/database/gorm/repositories"
	"go.uber.org/zap"
)

// InitializeDatabaseDependencies initializes the database and repositories
func (d *Dependencies) InitializeDatabaseDependencies() {
	// Create logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Printf("Failed to create logger: %v", err)
		return
	}

	// Initialize GORM database
	db, err := gorm.NewDatabase(gorm.Config{
		Path:   d.Config.Database.Path,
		Debug:  d.Config.App.Debug,
		Logger: logger,
	})
	if err != nil {
		logger.Error("Failed to initialize database", zap.Error(err))
		return
	}

	// Run migrations
	if err := db.Migrate(); err != nil {
		logger.Error("Failed to run database migrations", zap.Error(err))
		return
	}

	// Create repositories using GORM
	boughtCoinRepo := gormrepo.NewGORMBoughtCoinRepository(db.DB, logger)
	newCoinRepo := gormrepo.NewGORMNewCoinRepository(db.DB, logger)

	// Store repositories in dependencies
	d.BoughtCoinRepository = boughtCoinRepo
	d.NewCoinRepository = newCoinRepo
}
