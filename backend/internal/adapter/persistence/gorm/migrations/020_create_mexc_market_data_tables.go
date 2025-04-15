package migrations

import (
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm/entity"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// CreateMexcMarketDataTables creates tables for MEXC market data
func CreateMexcMarketDataTables(db *gorm.DB) error {
	logger := log.With().Str("migration", "create_mexc_market_data_tables").Logger()
	logger.Info().Msg("Running migration: Create MEXC market data tables")

	// Create tables using AutoMigrate
	entities := []interface{}{
		&entity.MexcTickerEntity{},
		&entity.MexcCandleEntity{},
		&entity.MexcOrderBookEntity{},
		&entity.MexcOrderBookEntryEntity{},
		&entity.MexcSymbolEntity{},
		&entity.MexcSyncStateEntity{},
	}

	for _, e := range entities {
		if err := db.AutoMigrate(e); err != nil {
			logger.Error().Err(err).Interface("entity", e).Msg("Failed to create table")
			return err
		}
	}

	// Add foreign key from order book entries to order books
	// Skip for SQLite as it has issues with ALTER TABLE ADD CONSTRAINT
	if db.Dialector.Name() != "sqlite" {
		err := db.Exec(
			"ALTER TABLE mexc_orderbook_entries " +
				"ADD CONSTRAINT fk_mexc_orderbook_entries_orderbook " +
				"FOREIGN KEY (order_book_id) REFERENCES mexc_orderbooks(id) ON DELETE CASCADE",
		).Error
		if err != nil {
			logger.Error().Err(err).Msg("Failed to add foreign key constraint to mexc_orderbook_entries")
			return err
		}
		logger.Info().Msg("Added foreign key constraint to mexc_orderbook_entries table")
	} else {
		logger.Info().Msg("Skipping foreign key constraint for SQLite")
	}

	// Create indexes for performance
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_mexc_tickers_timestamp ON mexc_tickers(timestamp)",
		"CREATE INDEX IF NOT EXISTS idx_mexc_candles_open_time ON mexc_candles(open_time)",
		"CREATE INDEX IF NOT EXISTS idx_mexc_orderbooks_timestamp ON mexc_orderbooks(timestamp)",
		"CREATE INDEX IF NOT EXISTS idx_mexc_sync_states_data_type ON mexc_sync_states(data_type)",
	}

	for _, idx := range indexes {
		if err := db.Exec(idx).Error; err != nil {
			logger.Error().Err(err).Str("index", idx).Msg("Failed to create index")
			return err
		}
	}

	// Initialize default sync states
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
		}
	}

	logger.Info().Msg("Successfully created MEXC market data tables")
	return nil
}
