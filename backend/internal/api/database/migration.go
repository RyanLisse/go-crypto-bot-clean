package database

import (
	"fmt"

	"go-crypto-bot-clean/backend/internal/api/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// MigrationManager handles database migrations
type MigrationManager struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager(db *gorm.DB, logger *zap.Logger) *MigrationManager {
	if logger == nil {
		var err error
		logger, err = zap.NewProduction()
		if err != nil {
			// Fallback to a simple logger if zap fails
			logger = zap.NewExample()
		}
	}

	return &MigrationManager{
		db:     db,
		logger: logger,
	}
}

// RunMigrations runs all migrations
func (m *MigrationManager) RunMigrations() error {
	m.logger.Info("Running database migrations")

	// Define all models to migrate
	models := []interface{}{
		// User models
		&models.User{},
		&models.UserRole{},
		&models.UserSettings{},
		&models.RefreshToken{},

		// Strategy models
		&models.Strategy{},
		&models.StrategyParameter{},
		&models.StrategyPerformance{},

		// Backtest models
		&models.Backtest{},
		&models.BacktestTrade{},
		&models.BacktestEquity{},
	}

	// Run migrations
	err := m.db.AutoMigrate(models...)
	if err != nil {
		m.logger.Error("Failed to run migrations", zap.Error(err))
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	m.logger.Info("Database migrations completed successfully")
	return nil
}

// SeedDatabase seeds the database with initial data
func (m *MigrationManager) SeedDatabase() error {
	m.logger.Info("Seeding database with initial data")

	// Check if admin user exists
	var count int64
	if err := m.db.Model(&models.User{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to count users: %w", err)
	}

	// If no users exist, create an admin user
	if count == 0 {
		m.logger.Info("Creating admin user")

		// Create admin user
		adminUser := models.User{
			ID:           "admin",
			Email:        "admin@example.com",
			Username:     "admin",
			PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // password: admin
			FirstName:    "Admin",
			LastName:     "User",
		}

		// Create admin role
		adminRole := models.UserRole{
			UserID: adminUser.ID,
			Role:   "admin",
		}

		// Create user settings
		adminSettings := models.UserSettings{
			UserID:              adminUser.ID,
			Theme:               "dark",
			Language:            "en",
			TimeZone:            "UTC",
			NotificationsEnabled: true,
			EmailNotifications:   true,
			PushNotifications:    false,
			DefaultCurrency:      "USD",
		}

		// Use a transaction to ensure all or nothing
		err := m.db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&adminUser).Error; err != nil {
				return fmt.Errorf("failed to create admin user: %w", err)
			}

			if err := tx.Create(&adminRole).Error; err != nil {
				return fmt.Errorf("failed to create admin role: %w", err)
			}

			if err := tx.Create(&adminSettings).Error; err != nil {
				return fmt.Errorf("failed to create admin settings: %w", err)
			}

			return nil
		})

		if err != nil {
			m.logger.Error("Failed to seed database", zap.Error(err))
			return fmt.Errorf("failed to seed database: %w", err)
		}

		m.logger.Info("Admin user created successfully")
	} else {
		m.logger.Info("Database already contains users, skipping seeding")
	}

	m.logger.Info("Database seeding completed successfully")
	return nil
}
