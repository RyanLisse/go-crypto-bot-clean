package api

import (
	"log"
	"time"

	"go-crypto-bot-clean/backend/internal/platform/database/gorm"
	gormrepo "go-crypto-bot-clean/backend/internal/platform/database/gorm/repositories"
	"go.uber.org/zap"
)

// InitializeDatabaseDependencies initializes the database and repositories
func (d *Dependencies) InitializeDatabaseDependencies() {
	// Create logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Printf("Failed to create logger: %v", err)
		// Consider returning an error or panicking if logger is essential
		return
	}

	// Initialize GORM database using the refactored NewDatabase
	gormDB, err := gorm.NewDatabase(gorm.Config{
		Path:            d.Config.Database.Path,
		Debug:           d.Config.App.Debug,
		Logger:          logger,
		MaxIdleConns:    d.Config.Database.MaxIdleConns,
		MaxOpenConns:    d.Config.Database.MaxOpenConns,
		ConnMaxLifetime: time.Duration(d.Config.Database.ConnMaxLifetimeSeconds) * time.Second, // Assuming config is in seconds
	})
	if err != nil {
		logger.Fatal("Failed to initialize GORM database", zap.Error(err)) // Fatal error if DB fails
		return
	}

	// Run migrations using the standalone function
	if err := gorm.RunMigrations(gormDB, logger); err != nil {
		logger.Fatal("Failed to run database migrations", zap.Error(err)) // Fatal error if migrations fail
		return
	}

	// Create repositories using GORM DB instance directly
	boughtCoinRepo := gormrepo.NewGORMBoughtCoinRepository(gormDB, logger)
	newCoinRepo := gormrepo.NewGORMNewCoinRepository(gormDB, logger)

	// Store repositories in dependencies
	d.BoughtCoinRepository = boughtCoinRepo
	d.NewCoinRepository = newCoinRepo

	logger.Info("Database dependencies initialized successfully")
}
