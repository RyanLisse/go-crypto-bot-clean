package migrations

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// PositionEntity is the GORM model for position data
type PositionEntity struct {
	ID            string `gorm:"primaryKey"`
	UserID        string `gorm:"index:idx_position_user_id"`
	Symbol        string `gorm:"index:idx_position_symbol"`
	Side          string
	Status        string `gorm:"index:idx_position_status"`
	Type          string
	EntryPrice    float64
	Quantity      float64
	CurrentPrice  float64
	PnL           float64
	PnLPercent    float64
	StopLoss      *float64
	TakeProfit    *float64
	StrategyID    *string
	Notes         string
	OpenedAt      string
	ClosedAt      *string
	LastUpdatedAt string
	CreatedAt     string
	UpdatedAt     string
}

// TableName sets the table name for PositionEntity
func (PositionEntity) TableName() string {
	return "positions"
}

// CreatePositionTable creates the position table
func CreatePositionTable(db *gorm.DB) error {
	logger := log.With().Str("migration", "create_position_table").Logger()
	logger.Info().Msg("Running migration: Create position table")

	// Create the position table
	if err := db.AutoMigrate(&PositionEntity{}); err != nil {
		logger.Error().Err(err).Msg("Failed to create position table")
		return err
	}

	logger.Info().Msg("Position table created successfully")
	return nil
}
