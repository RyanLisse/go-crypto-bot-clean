package port

import (
	"context"
	"time"

	"github.com/neo/crypto-bot/internal/domain/model"
)

// OrderRepository defines the interface for order persistence operations
type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) error
	GetByID(ctx context.Context, id string) (*model.Order, error)
	GetByClientOrderID(ctx context.Context, clientOrderID string) (*model.Order, error)
	Update(ctx context.Context, order *model.Order) error
	GetBySymbol(ctx context.Context, symbol string, limit, offset int) ([]*model.Order, error)
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*model.Order, error)
	GetByStatus(ctx context.Context, status model.OrderStatus, limit, offset int) ([]*model.Order, error)
	Count(ctx context.Context, filters map[string]interface{}) (int64, error)
	Delete(ctx context.Context, id string) error
}

// PositionRepository defines the interface for position persistence operations
type PositionRepository interface {
	Create(ctx context.Context, position *model.Position) error
	GetByID(ctx context.Context, id string) (*model.Position, error)
	Update(ctx context.Context, position *model.Position) error
	GetOpenPositions(ctx context.Context) ([]*model.Position, error)
	GetOpenPositionsBySymbol(ctx context.Context, symbol string) ([]*model.Position, error)
	GetOpenPositionsByType(ctx context.Context, positionType model.PositionType) ([]*model.Position, error)
	GetBySymbol(ctx context.Context, symbol string, limit, offset int) ([]*model.Position, error)
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*model.Position, error)
	GetClosedPositions(ctx context.Context, from, to time.Time, limit, offset int) ([]*model.Position, error)
	Count(ctx context.Context, filters map[string]interface{}) (int64, error)
	Delete(ctx context.Context, id string) error
}

// WalletRepository defines the interface for wallet persistence operations
type WalletRepository interface {
	Save(ctx context.Context, wallet *model.Wallet) error
	GetByUserID(ctx context.Context, userID string) (*model.Wallet, error)
	SaveBalanceHistory(ctx context.Context, history *model.BalanceHistory) error
	GetBalanceHistory(ctx context.Context, userID string, asset model.Asset, from, to time.Time) ([]*model.BalanceHistory, error)
}

// NewCoinRepository defines the interface for new coin persistence operations
type NewCoinRepository interface {
	Save(ctx context.Context, newCoin *model.NewCoin) error
	GetBySymbol(ctx context.Context, symbol string) (*model.NewCoin, error)
	GetRecent(ctx context.Context, limit int) ([]*model.NewCoin, error)
	GetByStatus(ctx context.Context, status model.NewCoinStatus) ([]*model.NewCoin, error)
	Update(ctx context.Context, newCoin *model.NewCoin) error
	// FindRecentlyListed retrieves coins expected to list soon or recently became tradable.
	FindRecentlyListed(ctx context.Context, thresholdTime time.Time) ([]*model.NewCoin, error)
}

// TickerRepository defines the interface for ticker persistence operations
type TickerRepository interface {
	Save(ctx context.Context, ticker *model.Ticker) error
	GetBySymbol(ctx context.Context, symbol string) (*model.Ticker, error)
	GetAll(ctx context.Context) ([]*model.Ticker, error)
	GetRecent(ctx context.Context, limit int) ([]*model.Ticker, error)
	SaveKline(ctx context.Context, kline *model.Kline) error
	GetKlines(ctx context.Context, symbol string, interval model.KlineInterval, from, to time.Time, limit int) ([]*model.Kline, error)
}

// AIConversationRepository defines the interface for AI conversation persistence
type AIConversationRepository interface {
	SaveMessage(ctx context.Context, userID string, message map[string]interface{}) error
	GetConversation(ctx context.Context, userID string, limit int) ([]map[string]interface{}, error)
	ClearConversation(ctx context.Context, userID string) error
}

// StrategyRepository defines the interface for strategy persistence
type StrategyRepository interface {
	SaveConfig(ctx context.Context, strategyID string, config map[string]interface{}) error
	GetConfig(ctx context.Context, strategyID string) (map[string]interface{}, error)
	ListStrategies(ctx context.Context) ([]string, error)
	DeleteStrategy(ctx context.Context, strategyID string) error
}

// NotificationRepository defines the interface for notification persistence
type NotificationRepository interface {
	SavePreferences(ctx context.Context, userID string, preferences map[string]interface{}) error
	GetPreferences(ctx context.Context, userID string) (map[string]interface{}, error)
	SaveNotification(ctx context.Context, notification map[string]interface{}) error
	GetNotifications(ctx context.Context, userID string, limit, offset int) ([]map[string]interface{}, error)
}

// AnalyticsRepository defines the interface for analytics data persistence
type AnalyticsRepository interface {
	SaveMetrics(ctx context.Context, metrics map[string]interface{}) error
	GetMetrics(ctx context.Context, from, to time.Time) ([]map[string]interface{}, error)
	GetPerformanceByStrategy(ctx context.Context, strategyID string, from, to time.Time) (map[string]interface{}, error)
}
