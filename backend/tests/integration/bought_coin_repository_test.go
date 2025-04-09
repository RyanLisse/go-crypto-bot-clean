package integration

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
	"github.com/ryanlisse/go-crypto-bot/internal/platform/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testDB *sqlx.DB

func TestMain(m *testing.M) {
	// Setup test database
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}

	testDB = db

	// Run migrations
	err = database.RunMigrations(testDB)
	if err != nil {
		panic(err)
	}

	// Run tests
	exitCode := m.Run()

	// Cleanup
	testDB.Close()

	os.Exit(exitCode)
}

func TestBoughtCoinRepository(t *testing.T) {
	repo := database.NewSQLiteBoughtCoinRepository(testDB)
	ctx := context.Background()

	// Test Create
	now := time.Now().Round(time.Second) // SQLite truncates timestamps to seconds
	coin := &models.BoughtCoin{
		Symbol:        "BTCUSDT",
		PurchasePrice: 50000.0,
		Quantity:      0.1,
		BoughtAt:      now,
		StopLoss:      47500.0,
		TakeProfit:    55000.0,
		CurrentPrice:  51000.0,
		IsDeleted:     false,
		UpdatedAt:     now,
	}

	// Create a new coin
	id, err := repo.Create(ctx, coin)
	require.NoError(t, err)
	assert.Greater(t, id, int64(0))

	// Test FindByID
	foundCoin, err := repo.FindByID(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, coin.Symbol, foundCoin.Symbol)
	assert.Equal(t, coin.PurchasePrice, foundCoin.PurchasePrice)
	assert.Equal(t, coin.Quantity, foundCoin.Quantity)
	assert.Equal(t, coin.BoughtAt.Unix(), foundCoin.BoughtAt.Unix())
	assert.Equal(t, coin.StopLoss, foundCoin.StopLoss)
	assert.Equal(t, coin.TakeProfit, foundCoin.TakeProfit)
	assert.Equal(t, coin.CurrentPrice, foundCoin.CurrentPrice)

	// Test FindBySymbol
	foundCoin, err = repo.FindBySymbol(ctx, "BTCUSDT")
	require.NoError(t, err)
	assert.Equal(t, id, foundCoin.ID)
	assert.Equal(t, coin.Symbol, foundCoin.Symbol)

	// Test Update
	foundCoin.CurrentPrice = 52000.0
	foundCoin.UpdatedAt = time.Now().Round(time.Second)
	err = repo.Update(ctx, foundCoin)
	require.NoError(t, err)

	updatedCoin, err := repo.FindByID(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, 52000.0, updatedCoin.CurrentPrice)

	// Test FindAll
	coins, err := repo.FindAll(ctx)
	require.NoError(t, err)
	assert.Len(t, coins, 1)

	// Test Delete
	err = repo.Delete(ctx, id)
	require.NoError(t, err)

	// Verify deletion works as expected
	coins, err = repo.FindAll(ctx)
	require.NoError(t, err)
	assert.Len(t, coins, 0)
}
