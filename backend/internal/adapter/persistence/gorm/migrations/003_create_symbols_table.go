package migrations

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// SymbolEntity is the GORM model for trading pair information
type SymbolEntity struct {
	Symbol            string `gorm:"primaryKey"`
	BaseAsset         string
	QuoteAsset        string
	Exchange          string `gorm:"index:idx_symbol_exchange"`
	Status            string
	MinPrice          float64
	MaxPrice          float64
	PricePrecision    int
	MinQty            float64
	MaxQty            float64
	QtyPrecision      int
	AllowedOrderTypes string
	CreatedAt         string
	UpdatedAt         string
}

// TableName sets the table name for SymbolEntity
func (SymbolEntity) TableName() string {
	return "symbols"
}

// CreateSymbolsTable creates the symbols table
func CreateSymbolsTable(db *gorm.DB) error {
	logger := log.With().Str("migration", "create_symbols_table").Logger()
	logger.Info().Msg("Running migration: Create symbols table")

	// Create the symbols table
	if err := db.AutoMigrate(&SymbolEntity{}); err != nil {
		logger.Error().Err(err).Msg("Failed to create symbols table")
		return err
	}

	logger.Info().Msg("Symbols table created successfully")
	return nil
}
