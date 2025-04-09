package balance

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go-crypto-bot-clean/backend/internal/domain/repository"
	"go-crypto-bot-clean/backend/internal/repository/database"
)

func setupTestDB(t *testing.T) database.Repository {
	// Create a temporary database file
	dbPath := "./test_balance_history.db"

	// Remove any existing test database
	os.Remove(dbPath)

	// Create configuration
	config := database.Config{
		DatabasePath:    dbPath,
		MaxOpenConns:    5,
		MaxIdleConns:    2,
		ConnMaxLifetime: 1 * time.Minute,
	}

	// Create repository
	repo, err := database.NewSQLiteRepository(config)
	require.NoError(t, err)

	// Create the balance history table
	ctx := context.Background()
	_, err = repo.Execute(ctx, `
		CREATE TABLE IF NOT EXISTS balance_history (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp DATETIME NOT NULL,
			balance REAL NOT NULL,
			equity REAL NOT NULL,
			free_balance REAL NOT NULL,
			locked_balance REAL NOT NULL,
			unrealized_pnl REAL NOT NULL
		)
	`)
	require.NoError(t, err)

	return repo
}

func TestBalanceHistoryRepository(t *testing.T) {
	// Setup
	ctx := context.Background()
	db := setupTestDB(t)
	defer db.Close()
	defer os.Remove("./test_balance_history.db")

	repo := NewBalanceHistoryRepository(db)

	// Test Create
	now := time.Now().UTC().Truncate(time.Second)
	history := &repository.BalanceHistory{
		Timestamp:     now,
		Balance:       1000.0,
		Equity:        1100.0,
		FreeBalance:   900.0,
		LockedBalance: 100.0,
		UnrealizedPnL: 100.0,
	}

	id, err := repo.Create(ctx, history)
	require.NoError(t, err)
	assert.Greater(t, id, int64(0))
	assert.Equal(t, id, history.ID)

	// Test GetLatestBalance
	latest, err := repo.GetLatestBalance(ctx)
	require.NoError(t, err)
	assert.Equal(t, id, latest.ID)
	assert.Equal(t, now, latest.Timestamp)
	assert.Equal(t, 1000.0, latest.Balance)
	assert.Equal(t, 1100.0, latest.Equity)
	assert.Equal(t, 900.0, latest.FreeBalance)
	assert.Equal(t, 100.0, latest.LockedBalance)
	assert.Equal(t, 100.0, latest.UnrealizedPnL)

	// Add more data for time range tests
	yesterday := now.Add(-24 * time.Hour)
	twoDaysAgo := now.Add(-48 * time.Hour)

	_, err = repo.Create(ctx, &repository.BalanceHistory{
		Timestamp:     yesterday,
		Balance:       950.0,
		Equity:        1000.0,
		FreeBalance:   850.0,
		LockedBalance: 100.0,
		UnrealizedPnL: 50.0,
	})
	require.NoError(t, err)

	_, err = repo.Create(ctx, &repository.BalanceHistory{
		Timestamp:     twoDaysAgo,
		Balance:       900.0,
		Equity:        950.0,
		FreeBalance:   800.0,
		LockedBalance: 100.0,
		UnrealizedPnL: 50.0,
	})
	require.NoError(t, err)

	// Test GetBalanceHistory
	historyList, err := repo.GetBalanceHistory(ctx, twoDaysAgo, now)
	require.NoError(t, err)
	assert.Len(t, historyList, 3)
	assert.Equal(t, twoDaysAgo, historyList[0].Timestamp)
	assert.Equal(t, yesterday, historyList[1].Timestamp)
	assert.Equal(t, now, historyList[2].Timestamp)

	// Test GetBalancePoints
	points, err := repo.GetBalancePoints(ctx, twoDaysAgo, now)
	require.NoError(t, err)
	assert.Len(t, points, 3)
	assert.Equal(t, twoDaysAgo, points[0].Timestamp)
	assert.Equal(t, 900.0, points[0].Balance)
	assert.Equal(t, yesterday, points[1].Timestamp)
	assert.Equal(t, 950.0, points[1].Balance)
	assert.Equal(t, now, points[2].Timestamp)
	assert.Equal(t, 1000.0, points[2].Balance)
}
