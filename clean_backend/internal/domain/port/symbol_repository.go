package port

import (
	"context"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
)

// SymbolRepository defines methods for symbol data persistence operations
type SymbolRepository interface {
	// Query operations
	GetBySymbol(ctx context.Context, symbol string) (*model.Symbol, error)
	GetBySymbolAndExchange(ctx context.Context, symbol, exchange string) (*model.Symbol, error)
	ListAll(ctx context.Context, limit, offset int) ([]*model.Symbol, error)
	ListByExchange(ctx context.Context, exchange string, limit, offset int) ([]*model.Symbol, error)
	ListByBaseAsset(ctx context.Context, baseAsset string, limit, offset int) ([]*model.Symbol, error)
	ListByQuoteAsset(ctx context.Context, quoteAsset string, limit, offset int) ([]*model.Symbol, error)
	Search(ctx context.Context, query string, limit, offset int) ([]*model.Symbol, error)
	
	// Count operations
	Count(ctx context.Context) (int64, error)
	CountByExchange(ctx context.Context, exchange string) (int64, error)
	CountByBaseAsset(ctx context.Context, baseAsset string) (int64, error)
	CountByQuoteAsset(ctx context.Context, quoteAsset string) (int64, error)
	
	// Write operations
	Save(ctx context.Context, symbol *model.Symbol) error
	SaveBatch(ctx context.Context, symbols []*model.Symbol) error
	Update(ctx context.Context, symbol *model.Symbol) error
	UpdateStatus(ctx context.Context, symbol, exchange, status string) error
	Delete(ctx context.Context, symbol, exchange string) error
	
	// Sync state operations
	GetLastUpdated(ctx context.Context, exchange string) (time.Time, error)
	SetLastUpdated(ctx context.Context, exchange string, timestamp time.Time) error
}
