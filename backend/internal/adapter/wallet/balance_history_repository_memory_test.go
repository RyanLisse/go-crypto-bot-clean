package wallet

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
)

func TestInMemoryBalanceHistoryRepository_SaveAndFind(t *testing.T) {
	repo := NewInMemoryBalanceHistoryRepository()
	ctx := context.Background()
	walletID := "wallet123"
	userID := "user123"

	rec1 := &model.BalanceHistory{
		ID:            "1",
		UserID:        userID,
		WalletID:      walletID,
		Balances:      map[model.Asset]*model.Balance{"BTC": {Asset: "BTC", Free: 1.0, Total: 1.0}},
		TotalUSDValue: 30000.0,
		Timestamp:     time.Now().Add(-10 * time.Minute),
	}
	rec2 := &model.BalanceHistory{
		ID:            "2",
		UserID:        userID,
		WalletID:      walletID,
		Balances:      map[model.Asset]*model.Balance{"BTC": {Asset: "BTC", Free: 2.0, Total: 2.0}},
		TotalUSDValue: 60000.0,
		Timestamp:     time.Now().Add(-5 * time.Minute),
	}
	rec3 := &model.BalanceHistory{
		ID:            "3",
		UserID:        userID,
		WalletID:      walletID,
		Balances:      map[model.Asset]*model.Balance{"BTC": {Asset: "BTC", Free: 3.0, Total: 3.0}},
		TotalUSDValue: 90000.0,
		Timestamp:     time.Now(),
	}

	// Save records
	assert.NoError(t, repo.Save(ctx, rec1))
	assert.NoError(t, repo.Save(ctx, rec2))
	assert.NoError(t, repo.Save(ctx, rec3))

	// Find all since rec2.Timestamp
	since := rec2.Timestamp
	found, err := repo.FindByWalletID(ctx, walletID, since)
	assert.NoError(t, err)
	assert.Len(t, found, 2)
	assert.Equal(t, rec2.ID, found[0].ID)
	assert.Equal(t, rec3.ID, found[1].ID)

	// Find latest
	latest, err := repo.FindLatestByWalletID(ctx, walletID)
	assert.NoError(t, err)
	assert.NotNil(t, latest)
	assert.Equal(t, rec3.ID, latest.ID)
}
