package factory

import (
	"path/filepath"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/csv"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
	gormdb "gorm.io/gorm"
)

// TradeHistoryFactory creates trade history components
type TradeHistoryFactory struct {
	config *config.Config
	logger *zerolog.Logger
	db     *gormdb.DB
}

// NewTradeHistoryFactory creates a new trade history factory
func NewTradeHistoryFactory(config *config.Config, logger *zerolog.Logger, db *gormdb.DB) *TradeHistoryFactory {
	return &TradeHistoryFactory{
		config: config,
		logger: logger,
		db:     db,
	}
}

// CreateTradeHistoryRepository creates a new trade history repository
func (f *TradeHistoryFactory) CreateTradeHistoryRepository() port.TradeHistoryRepository {
	// Create logger for the repository
	repoLogger := f.logger.With().Str("component", "trade_history_repository").Logger()
	
	// Create and return the repository
	return gorm.NewTradeHistoryRepository(f.db, &repoLogger)
}

// CreateTradeHistoryWriter creates a new CSV trade history writer
func (f *TradeHistoryFactory) CreateTradeHistoryWriter() (*csv.TradeHistoryWriter, error) {
	// Create logger for the writer
	writerLogger := f.logger.With().Str("component", "trade_history_writer").Logger()
	
	// Determine base directory
	baseDir := filepath.Join(filepath.Dir(f.config.Database.Path), "csv")
	
	// Create writer config
	writerConfig := csv.TradeHistoryWriterConfig{
		Enabled:           true,
		BaseDirectory:     baseDir,
		TradeFilename:     "trade_records.csv",
		DetectionFilename: "detection_logs.csv",
		FlushInterval:     30 * time.Second,
	}
	
	// Create and return the writer
	return csv.NewTradeHistoryWriter(writerConfig, &writerLogger)
}
