package migrations

import (
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm/entity"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm/repo"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// MigrationRecord tracks applied migrations
type MigrationRecord struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"uniqueIndex"`
	AppliedAt time.Time `gorm:"autoCreateTime"`
}

// TableName sets the table name for MigrationRecord
func (MigrationRecord) TableName() string {
	return "migrations"
}

// RunConsolidatedMigrations is the main entry point for database migrations
// It uses GORM's AutoMigrate as the standardized migration strategy
func RunConsolidatedMigrations(db *gorm.DB, logger *zerolog.Logger) error {
	migrationStart := time.Now()
	logger.Info().Msg("Starting consolidated database migrations")

	// Create migrations table if it doesn't exist
	if err := db.AutoMigrate(&MigrationRecord{}); err != nil {
		logger.Error().Err(err).Msg("Failed to create migrations table")
		return err
	}

	// List of all entity models to migrate
	entities := []interface{}{
		// User and authentication entities
		&entity.UserEntity{},
		&entity.APICredentialEntity{},

		// Wallet and balance entities
		&repo.EnhancedWalletEntity{},
		&repo.EnhancedWalletBalanceEntity{},
		&repo.EnhancedWalletBalanceHistoryEntity{},

		// Market data entities
		&entity.Symbol{},
		&entity.Ticker{},
		&entity.OrderBook{},
		&entity.Candle{},
		&entity.MexcTickerEntity{},
		&entity.MexcCandleEntity{},
		&entity.MexcOrderBookEntity{},
		&entity.MexcOrderBookEntryEntity{},
		&entity.MexcSymbolEntity{},
		&entity.MexcSyncStateEntity{},

		// Trading entities
		&entity.Position{},
		&entity.OrderEntity{},
		&entity.TransactionEntity{},
		&entity.AutoBuyRuleEntity{},
		&entity.AutoBuyExecutionEntity{},

		// System entities
		&entity.StatusEntity{},
		&repo.StatusRecord{},
	}

	// Run AutoMigrate on all models
	for _, model := range entities {
		entityName := getEntityName(model)
		logger.Debug().Str("entity", entityName).Msg("Migrating entity")
		
		if err := db.AutoMigrate(model); err != nil {
			logger.Error().Err(err).Str("entity", entityName).Msg("Failed to migrate entity")
			return err
		}
		
		logger.Debug().Str("entity", entityName).Msg("Successfully migrated entity")
	}

	// Add foreign key constraints (skip for SQLite as it has issues with ALTER TABLE ADD CONSTRAINT)
	if db.Dialector.Name() != "sqlite" {
		constraints := []struct {
			query string
			desc  string
		}{
			{
				query: "ALTER TABLE auto_buy_executions " +
					"ADD CONSTRAINT IF NOT EXISTS fk_auto_buy_executions_rule " +
					"FOREIGN KEY (rule_id) REFERENCES auto_buy_rules(id) ON DELETE CASCADE",
				desc: "foreign key from auto_buy_executions to auto_buy_rules",
			},
			{
				query: "ALTER TABLE mexc_orderbook_entries " +
					"ADD CONSTRAINT IF NOT EXISTS fk_mexc_orderbook_entries_orderbook " +
					"FOREIGN KEY (order_book_id) REFERENCES mexc_orderbooks(id) ON DELETE CASCADE",
				desc: "foreign key from mexc_orderbook_entries to mexc_orderbooks",
			},
		}

		for _, constraint := range constraints {
			err := db.Exec(constraint.query).Error
			if err != nil {
				logger.Error().Err(err).Str("constraint", constraint.desc).Msg("Failed to add constraint")
				return err
			}
			logger.Info().Str("constraint", constraint.desc).Msg("Successfully added constraint")
		}
	} else {
		logger.Info().Msg("Skipping foreign key constraints for SQLite")
	}

	// Initialize default sync states for MEXC data
	syncStates := []entity.MexcSyncStateEntity{
		{
			DataType:     "tickers",
			SyncInterval: 60, // 1 minute
			Status:       "idle",
		},
		{
			DataType:     "candles",
			SyncInterval: 300, // 5 minutes
			Status:       "idle",
		},
		{
			DataType:     "orderbooks",
			SyncInterval: 30, // 30 seconds
			Status:       "idle",
		},
		{
			DataType:     "symbols",
			SyncInterval: 3600, // 1 hour
			Status:       "idle",
		},
	}

	for _, state := range syncStates {
		var count int64
		db.Model(&entity.MexcSyncStateEntity{}).Where("data_type = ?", state.DataType).Count(&count)
		if count == 0 {
			if err := db.Create(&state).Error; err != nil {
				logger.Error().Err(err).Str("dataType", state.DataType).Msg("Failed to create sync state")
				return err
			}
			logger.Info().Str("dataType", state.DataType).Msg("Created default sync state")
		}
	}

	// Record this migration
	migration := MigrationRecord{
		Name: "consolidated_migration_" + time.Now().Format("20060102150405"),
	}
	if err := db.Create(&migration).Error; err != nil {
		logger.Warn().Err(err).Msg("Failed to record migration, but schema changes were applied")
	}

	logger.Info().
		Str("total_duration", time.Since(migrationStart).String()).
		Msg("All database migrations completed successfully")
	return nil
}

// getEntityName returns the name of an entity for logging purposes
func getEntityName(model interface{}) string {
	if t, ok := model.(interface{ TableName() string }); ok {
		return t.TableName()
	}
	return "unknown"
}
