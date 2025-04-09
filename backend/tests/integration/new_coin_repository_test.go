package integration

import (
	"context"
	"testing"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
	"go-crypto-bot-clean/backend/internal/platform/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCoinRepository(t *testing.T) {
	repo := database.NewSQLiteNewCoinRepository(testDB)
	ctx := context.Background()

	// Test Create
	now := time.Now().Round(time.Second) // SQLite truncates timestamps to seconds
	// Create a coin with all required fields
	firstOpenTime := now.Add(-1 * time.Hour)
	becameTradableAt := now.Add(-30 * time.Minute)
	coin := &models.NewCoin{
		Symbol:           "ETHUSDT",
		FoundAt:          now,
		FirstOpenTime:    &firstOpenTime,
		BecameTradableAt: &becameTradableAt,
		BaseVolume:       1000.0,
		QuoteVolume:      3000000.0,
		Status:           "1",
		IsProcessed:      false,
		IsDeleted:        false,
		IsUpcoming:       false,
	}

	// Create a new coin
	id, err := repo.Create(ctx, coin)
	require.NoError(t, err)
	assert.Greater(t, id, int64(0))

	// Test FindByID
	foundCoin, err := repo.FindByID(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, coin.Symbol, foundCoin.Symbol)
	assert.Equal(t, coin.FoundAt.Unix(), foundCoin.FoundAt.Unix())
	assert.Equal(t, coin.BaseVolume, foundCoin.BaseVolume)
	assert.Equal(t, coin.QuoteVolume, foundCoin.QuoteVolume)
	assert.Equal(t, coin.IsProcessed, foundCoin.IsProcessed)
	assert.Equal(t, coin.IsDeleted, foundCoin.IsDeleted)

	// Test FindBySymbol
	foundCoin, err = repo.FindBySymbol(ctx, "ETHUSDT")
	require.NoError(t, err)
	assert.Equal(t, id, foundCoin.ID)
	assert.Equal(t, coin.Symbol, foundCoin.Symbol)

	// Test FindAll
	coins, err := repo.FindAll(ctx)
	require.NoError(t, err)
	assert.Len(t, coins, 1)

	// Test MarkAsProcessed
	err = repo.MarkAsProcessed(ctx, id)
	require.NoError(t, err)

	updatedCoin, err := repo.FindByID(ctx, id)
	require.NoError(t, err)
	assert.True(t, updatedCoin.IsProcessed)

	// Test Delete
	err = repo.Delete(ctx, id)
	require.NoError(t, err)

	// Verify deletion works as expected
	coins, err = repo.FindAll(ctx)
	require.NoError(t, err)
	assert.Len(t, coins, 0)
}
