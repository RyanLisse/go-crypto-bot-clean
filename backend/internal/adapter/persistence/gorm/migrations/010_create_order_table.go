package migrations

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// OrderEntity is the GORM model for order data
type OrderEntity struct {
	ID            string `gorm:"primaryKey"`
	UserID        string `gorm:"index:idx_order_user_id"`
	Symbol        string
	Side          string
	Type          string
	Status        string `gorm:"index:idx_order_status"`
	Price         float64
	Quantity      float64
	FilledQty     float64
	RemainingQty  float64
	ClientOrderID string `gorm:"index:idx_order_client_id"`
	Exchange      string
	CreatedAt     string
	UpdatedAt     string
}

// TableName sets the table name for OrderEntity
func (OrderEntity) TableName() string {
	return "orders"
}

// CreateOrderTable creates the order table
func CreateOrderTable(db *gorm.DB) error {
	logger := log.With().Str("migration", "create_order_table").Logger()
	logger.Info().Msg("Running migration: Create order table")

	// Create the order table
	if err := db.AutoMigrate(&OrderEntity{}); err != nil {
		logger.Error().Err(err).Msg("Failed to create order table")
		return err
	}

	logger.Info().Msg("Order table created successfully")
	return nil
}
