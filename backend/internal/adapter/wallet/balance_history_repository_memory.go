package wallet

import (
	"context"
	"sync"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
)

// InMemoryBalanceHistoryRepository implements BalanceHistoryRepository with in-memory storage (for prototyping/tests)
type InMemoryBalanceHistoryRepository struct {
	mu      sync.RWMutex
	records map[string][]*model.BalanceHistory // walletID -> list (sorted by timestamp asc)
}

func NewInMemoryBalanceHistoryRepository() port.BalanceHistoryRepository {
	return &InMemoryBalanceHistoryRepository{
		records: make(map[string][]*model.BalanceHistory),
	}
}

func (r *InMemoryBalanceHistoryRepository) Save(ctx context.Context, record *model.BalanceHistory) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	arr := r.records[record.WalletID]
	arr = append(arr, record)
	// Sort by timestamp ascending
	for i := len(arr) - 1; i > 0; i-- {
		if arr[i].Timestamp.Before(arr[i-1].Timestamp) {
			arr[i], arr[i-1] = arr[i-1], arr[i]
		}
	}
	r.records[record.WalletID] = arr
	return nil
}

func (r *InMemoryBalanceHistoryRepository) FindByWalletID(ctx context.Context, walletID string, since time.Time) ([]*model.BalanceHistory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	arr := r.records[walletID]
	var result []*model.BalanceHistory
	for _, rec := range arr {
		if rec.Timestamp.After(since) || rec.Timestamp.Equal(since) {
			result = append(result, rec)
		}
	}
	return result, nil
}

func (r *InMemoryBalanceHistoryRepository) FindLatestByWalletID(ctx context.Context, walletID string) (*model.BalanceHistory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	arr := r.records[walletID]
	if len(arr) == 0 {
		return nil, nil
	}
	return arr[len(arr)-1], nil
}
