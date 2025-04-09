package repository

import (
	"context"
	"time"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
)

// NewCoinRepository defines operations for managing new coins
type NewCoinRepository interface {
	// FindAll returns all new coins that haven't been processed
	FindAll(ctx context.Context) ([]models.NewCoin, error)

	// FindByID returns a specific new coin by ID
	FindByID(ctx context.Context, id int64) (*models.NewCoin, error)

	// FindBySymbol returns a specific new coin by symbol
	FindBySymbol(ctx context.Context, symbol string) (*models.NewCoin, error)

	// Create adds a new coin listing
	Create(ctx context.Context, coin *models.NewCoin) (int64, error)

	// Update updates an existing coin
	Update(ctx context.Context, coin *models.NewCoin) error

	// MarkAsProcessed marks a new coin as processed
	MarkAsProcessed(ctx context.Context, id int64) error

	// Delete marks a new coin as deleted
	Delete(ctx context.Context, id int64) error

	// FindByDateRange finds coins within a date range
	FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]models.NewCoin, error)

	// FindUpcomingCoins finds coins that are scheduled to be listed in the future
	FindUpcomingCoins(ctx context.Context) ([]models.NewCoin, error)

	// FindUpcomingCoinsByDate finds upcoming coins that will be listed on a specific date
	FindUpcomingCoinsByDate(ctx context.Context, date time.Time) ([]models.NewCoin, error)

	// FindTradableCoins finds coins that have become tradable (status = "1")
	FindTradableCoins(ctx context.Context) ([]models.NewCoin, error)

	// FindTradableCoinsByDate finds coins that became tradable on a specific date
	FindTradableCoinsByDate(ctx context.Context, date time.Time) ([]models.NewCoin, error)
}
