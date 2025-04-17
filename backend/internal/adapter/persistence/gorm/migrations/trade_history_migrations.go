package migrations

import (
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// MigrateTradeHistory creates the trade history tables
func MigrateTradeHistory(db *gorm.DB, logger *zerolog.Logger) error {
	logger.Info().Msg("Running trade history migrations")

	// Create trade_records table
	err := db.Exec(`
		CREATE TABLE IF NOT EXISTS trade_records (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			symbol TEXT NOT NULL,
			side TEXT NOT NULL,
			type TEXT NOT NULL,
			quantity REAL NOT NULL,
			price REAL NOT NULL,
			amount REAL NOT NULL,
			fee REAL NOT NULL,
			fee_currency TEXT NOT NULL,
			order_id TEXT NOT NULL,
			trade_id TEXT,
			execution_time TIMESTAMP NOT NULL,
			strategy TEXT,
			notes TEXT,
			tags TEXT,
			metadata TEXT,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`).Error
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create trade_records table")
		return err
	}

	// Create indexes for trade_records
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_trade_records_user_id ON trade_records(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_trade_records_symbol ON trade_records(symbol)",
		"CREATE INDEX IF NOT EXISTS idx_trade_records_side ON trade_records(side)",
		"CREATE INDEX IF NOT EXISTS idx_trade_records_order_id ON trade_records(order_id)",
		"CREATE INDEX IF NOT EXISTS idx_trade_records_trade_id ON trade_records(trade_id)",
		"CREATE INDEX IF NOT EXISTS idx_trade_records_execution_time ON trade_records(execution_time)",
		"CREATE INDEX IF NOT EXISTS idx_trade_records_strategy ON trade_records(strategy)",
	}

	for _, idx := range indexes {
		if err := db.Exec(idx).Error; err != nil {
			logger.Error().Err(err).Str("index", idx).Msg("Failed to create index")
			return err
		}
	}

	// Create detection_logs table
	err = db.Exec(`
		CREATE TABLE IF NOT EXISTS detection_logs (
			id TEXT PRIMARY KEY,
			type TEXT NOT NULL,
			symbol TEXT NOT NULL,
			value REAL NOT NULL,
			threshold REAL NOT NULL,
			description TEXT,
			metadata TEXT,
			detected_at TIMESTAMP NOT NULL,
			processed_at TIMESTAMP,
			processed BOOLEAN NOT NULL DEFAULT FALSE,
			result TEXT,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`).Error
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create detection_logs table")
		return err
	}

	// Create indexes for detection_logs
	indexes = []string{
		"CREATE INDEX IF NOT EXISTS idx_detection_logs_type ON detection_logs(type)",
		"CREATE INDEX IF NOT EXISTS idx_detection_logs_symbol ON detection_logs(symbol)",
		"CREATE INDEX IF NOT EXISTS idx_detection_logs_detected_at ON detection_logs(detected_at)",
		"CREATE INDEX IF NOT EXISTS idx_detection_logs_processed ON detection_logs(processed)",
	}

	for _, idx := range indexes {
		if err := db.Exec(idx).Error; err != nil {
			logger.Error().Err(err).Str("index", idx).Msg("Failed to create index")
			return err
		}
	}

	logger.Info().Msg("Trade history migrations completed successfully")
	return nil
}
