package gorm

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm/entity"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm/repo"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/rs/zerolog"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// NewDBConnection creates a new GORM database connection
func NewDBConnection(cfg *config.Config, logger zerolog.Logger) (*gorm.DB, error) {
	// Ensure the database directory exists
	dbDir := filepath.Dir(cfg.Database.Path)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Configure GORM logger based on environment
	var logLevel gormlogger.LogLevel
	if cfg.ENV == "development" {
		logLevel = gormlogger.Info
	} else {
		logLevel = gormlogger.Error
	}

	// Create a log writer that uses our zerolog instance
	logWriter := log.New(zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false}, "[GORM] ", log.LstdFlags)

	gormLogger := gormlogger.New(
		logWriter,
		gormlogger.Config{
			SlowThreshold:             2 * time.Second, // Threshold for slow SQL queries
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// Connect to the database based on driver type
	var db *gorm.DB
	var err error

	switch cfg.Database.Driver {
	case "sqlite":
		// Connect to SQLite database
		db, err = gorm.Open(sqlite.Open(cfg.Database.Path), &gorm.Config{
			Logger: gormLogger,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to connect to SQLite database: %w", err)
		}
		logger.Info().Str("path", cfg.Database.Path).Msg("Connected to SQLite database")

	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Database.Driver)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to configure connection pool: %w", err)
	}

	// Set connection pool parameters - use sensible defaults for SQLite
	// SQLite only supports one writer at a time
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetMaxOpenConns(10) // Lower for SQLite
	sqlDB.SetConnMaxLifetime(time.Minute * 30)

	logger.Info().Msg("Database connection established successfully")
	return db, nil
}

// AutoMigrateModels performs automatic migrations for all required models
func AutoMigrateModels(db *gorm.DB, logger *zerolog.Logger) error {
	migrationStart := time.Now()
	logger.Info().Msg("Starting database migrations")

	// Slice of entities to migrate
	entities := []interface{}{
		&TickerEntity{},
		&SymbolEntity{},
		&PositionEntity{},
		&WalletEntity{},
		&OrderEntity{},
		&repo.StatusRecord{},
		&entity.AutoBuyRuleEntity{},
		&entity.AutoBuyExecutionEntity{},
		&entity.MexcTickerEntity{},
		&entity.MexcCandleEntity{},
		&entity.MexcOrderBookEntity{},
		&entity.MexcOrderBookEntryEntity{},
		&entity.MexcSymbolEntity{},
		&entity.MexcSyncStateEntity{},
		&repo.EnhancedWalletEntity{},
		&repo.EnhancedWalletBalanceEntity{},
		&repo.EnhancedWalletBalanceHistoryEntity{},
		// Add other entities as they are implemented
	}

	// Migrate each entity
	for _, entity := range entities {
		typeName := reflect.TypeOf(entity).Elem().Name()
		start := time.Now()

		if err := db.AutoMigrate(entity); err != nil {
			logger.Error().
				Err(err).
				Str("entity", typeName).
				Str("duration", time.Since(start).String()).
				Msg("Failed to migrate entity")
			return fmt.Errorf("failed to migrate %s: %w", typeName, err)
		}

		logger.Info().
			Str("entity", typeName).
			Str("duration", time.Since(start).String()).
			Msg("Successfully migrated entity")
	}

	// Add foreign key constraints (skip for SQLite as it has issues with ALTER TABLE ADD CONSTRAINT)
	// SQLite supports foreign keys but they must be defined when the table is created
	// GORM should handle this automatically with the references in the struct tags
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

	logger.Info().
		Str("total_duration", time.Since(migrationStart).String()).
		Msg("All database migrations completed successfully")
	return nil
}
