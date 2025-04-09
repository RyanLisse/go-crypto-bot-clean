package repositories

import (
	"context"
	"time"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
)

// BoughtCoinRepository defines the interface for bought coin repository operations
type BoughtCoinRepository interface {
	FindAll(ctx context.Context) ([]*models.BoughtCoin, error)
	FindBySymbol(ctx context.Context, symbol string) (*models.BoughtCoin, error)
	FindByID(ctx context.Context, id int64) (*models.BoughtCoin, error)
	FindAllActive(ctx context.Context) ([]*models.BoughtCoin, error)
	Save(ctx context.Context, coin *models.BoughtCoin) error
	Delete(ctx context.Context, symbol string) error
	DeleteByID(ctx context.Context, id int64) error
	UpdatePrice(ctx context.Context, symbol string, price float64) error
	HardDelete(ctx context.Context, symbol string) error
	Count(ctx context.Context) (int64, error)
}

// NewCoinRepository defines the interface for new coin repository operations
type NewCoinRepository interface {
	FindAll(ctx context.Context) ([]*models.NewCoin, error)
	FindBySymbol(ctx context.Context, symbol string) (*models.NewCoin, error)
	FindByID(ctx context.Context, id int64) (*models.NewCoin, error)
	FindByDate(ctx context.Context, date time.Time) ([]*models.NewCoin, error)
	FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*models.NewCoin, error)
	FindUpcoming(ctx context.Context) ([]*models.NewCoin, error)
	FindUpcomingByDate(ctx context.Context, date time.Time) ([]*models.NewCoin, error)
	FindTradable(ctx context.Context) ([]*models.NewCoin, error)
	FindTradableToday(ctx context.Context) ([]*models.NewCoin, error)
	Save(ctx context.Context, coin *models.NewCoin) error
	SaveAll(ctx context.Context, coins []*models.NewCoin) error
	Delete(ctx context.Context, symbol string) error
	MarkAsProcessed(ctx context.Context, id int64) error
	Count(ctx context.Context) (int64, error)
}
