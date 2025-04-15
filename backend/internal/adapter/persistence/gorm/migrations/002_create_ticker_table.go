package migrations

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// TickerEntity is the GORM model for market ticker data
type TickerEntity struct {
	ID            string `gorm:"primaryKey"`
	Symbol        string `gorm:"index:idx_ticker_symbol"`
	Price         float64
	Volume        float64
	High24h       float64
	Low24h        float64
	PriceChange   float64
	PercentChange float64
	LastUpdated   string `gorm:"index:idx_ticker_updated"`
	Exchange      string `gorm:"index:idx_ticker_exchange"`
	CreatedAt     string
	UpdatedAt     string
}

// TableName sets the table name for TickerEntity
func (TickerEntity) TableName() string {
	return "tickers"
}

// CreateTickerTable creates the ticker table
func CreateTickerTable(db *gorm.DB) error {
	logger := log.With().Str("migration", "create_ticker_table").Logger()
	logger.Info().Msg("Running migration: Create ticker table")

	// Create the ticker table
	if err := db.AutoMigrate(&TickerEntity{}); err != nil {
		logger.Error().Err(err).Msg("Failed to create ticker table")
		return err
	}

	logger.Info().Msg("Ticker table created successfully")
	return nil
}
