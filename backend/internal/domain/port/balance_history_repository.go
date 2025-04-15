package port

import (
	"context"
	"time"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
)

// BalanceHistoryRepository defines persistence for balance history records
// Used for tracking historical wallet balances (for charting, auditing, etc)
type BalanceHistoryRepository interface {
	Save(ctx context.Context, record *model.BalanceHistory) error
	FindByWalletID(ctx context.Context, walletID string, since time.Time) ([]*model.BalanceHistory, error)
	FindLatestByWalletID(ctx context.Context, walletID string) (*model.BalanceHistory, error)
}
