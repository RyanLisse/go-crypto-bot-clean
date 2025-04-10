package repository

import (
	"context"
	"testing"
	"time"

	"go-crypto-bot-clean/backend/api/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBacktestRepository(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormBacktestRepository(db)
	ctx := context.Background()

	// Create a test user and strategy first
	userRepo := NewGormUserRepository(db)
	strategyRepo := NewGormStrategyRepository(db)

	user := &models.User{
		ID:           "test-user",
		Email:        "test@example.com",
		Username:     "testuser",
		PasswordHash: "password-hash",
	}
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	strategy := &models.Strategy{
		ID:          "test-strategy",
		Name:        "Test Strategy",
		Description: "A test strategy",
		Parameters: models.Parameters{
			"param1": "value1",
			"param2": 42,
		},
		IsEnabled: true,
		UserID:    user.ID,
	}
	err = strategyRepo.Create(ctx, strategy)
	require.NoError(t, err)

	// Test creating a backtest
	backtest := &models.Backtest{
		ID:             "test-backtest",
		UserID:         user.ID,
		StrategyID:     strategy.ID,
		Name:           "Test Backtest",
		Description:    "A test backtest",
		StartDate:      time.Now().AddDate(0, -1, 0),
		EndDate:        time.Now(),
		InitialBalance: 10000,
		FinalBalance:   12000,
		TotalTrades:    50,
		WinningTrades:  30,
		LosingTrades:   20,
		WinRate:        0.6,
		ProfitFactor:   1.5,
		SharpeRatio:    1.2,
		MaxDrawdown:    0.1,
		Parameters: models.Parameters{
			"param1": "value1",
			"param2": 42,
		},
		Status: "completed",
	}

	err = repo.Create(ctx, backtest)
	assert.NoError(t, err)

	// Test getting a backtest by ID
	retrievedBacktest, err := repo.GetByID(ctx, backtest.ID)
	assert.NoError(t, err)
	assert.Equal(t, backtest.ID, retrievedBacktest.ID)
	assert.Equal(t, backtest.UserID, retrievedBacktest.UserID)
	assert.Equal(t, backtest.StrategyID, retrievedBacktest.StrategyID)
	assert.Equal(t, backtest.Name, retrievedBacktest.Name)
	assert.Equal(t, backtest.Description, retrievedBacktest.Description)
	assert.Equal(t, backtest.InitialBalance, retrievedBacktest.InitialBalance)
	assert.Equal(t, backtest.FinalBalance, retrievedBacktest.FinalBalance)
	assert.Equal(t, backtest.TotalTrades, retrievedBacktest.TotalTrades)
	assert.Equal(t, backtest.WinningTrades, retrievedBacktest.WinningTrades)
	assert.Equal(t, backtest.LosingTrades, retrievedBacktest.LosingTrades)
	assert.Equal(t, backtest.WinRate, retrievedBacktest.WinRate)
	assert.Equal(t, backtest.ProfitFactor, retrievedBacktest.ProfitFactor)
	assert.Equal(t, backtest.SharpeRatio, retrievedBacktest.SharpeRatio)
	assert.Equal(t, backtest.MaxDrawdown, retrievedBacktest.MaxDrawdown)
	assert.Equal(t, backtest.Status, retrievedBacktest.Status)
	assert.Equal(t, "value1", retrievedBacktest.Parameters["param1"])
	assert.Equal(t, float64(42), retrievedBacktest.Parameters["param2"])

	// Test getting backtests by user ID
	backtests, err := repo.GetByUserID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Len(t, backtests, 1)
	assert.Equal(t, backtest.ID, backtests[0].ID)

	// Test getting backtests by strategy ID
	backtests, err = repo.GetByStrategyID(ctx, strategy.ID)
	assert.NoError(t, err)
	assert.Len(t, backtests, 1)
	assert.Equal(t, backtest.ID, backtests[0].ID)

	// Test updating a backtest
	backtest.Name = "Updated Backtest"
	backtest.Description = "An updated test backtest"
	backtest.FinalBalance = 13000
	backtest.Parameters["param3"] = "value3"
	err = repo.Update(ctx, backtest)
	assert.NoError(t, err)

	retrievedBacktest, err = repo.GetByID(ctx, backtest.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Backtest", retrievedBacktest.Name)
	assert.Equal(t, "An updated test backtest", retrievedBacktest.Description)
	assert.Equal(t, 13000.0, retrievedBacktest.FinalBalance)
	assert.Equal(t, "value3", retrievedBacktest.Parameters["param3"])

	// Test adding a trade
	trade := &models.BacktestTrade{
		BacktestID: backtest.ID,
		Symbol:     "BTC/USD",
		EntryTime:  time.Now().Add(-time.Hour),
		EntryPrice: 50000,
		Quantity:   1,
		Direction:  "long",
	}
	err = repo.AddTrade(ctx, trade)
	assert.NoError(t, err)

	// Test getting trades
	trades, err := repo.GetTrades(ctx, backtest.ID)
	assert.NoError(t, err)
	assert.Len(t, trades, 1)
	assert.Equal(t, "BTC/USD", trades[0].Symbol)
	assert.Equal(t, 50000.0, trades[0].EntryPrice)
	assert.Equal(t, 1.0, trades[0].Quantity)
	assert.Equal(t, "long", trades[0].Direction)

	// Test adding an equity point
	equity := &models.BacktestEquity{
		BacktestID: backtest.ID,
		Timestamp:  time.Now(),
		Equity:     11000,
		Balance:    11000,
		Drawdown:   0.05,
	}
	err = repo.AddEquityPoint(ctx, equity)
	assert.NoError(t, err)

	// Test getting equity curve
	equityCurve, err := repo.GetEquityCurve(ctx, backtest.ID)
	assert.NoError(t, err)
	assert.Len(t, equityCurve, 1)
	assert.Equal(t, 11000.0, equityCurve[0].Equity)
	assert.Equal(t, 11000.0, equityCurve[0].Balance)
	assert.Equal(t, 0.05, equityCurve[0].Drawdown)

	// Test deleting a backtest
	err = repo.Delete(ctx, backtest.ID)
	assert.NoError(t, err)

	_, err = repo.GetByID(ctx, backtest.ID)
	assert.Error(t, err)
	assert.Equal(t, ErrBacktestNotFound, err)

	// Verify that trades and equity points were also deleted
	trades, err = repo.GetTrades(ctx, backtest.ID)
	assert.Error(t, err)
	assert.Equal(t, ErrBacktestNotFound, err)

	equityCurve, err = repo.GetEquityCurve(ctx, backtest.ID)
	assert.Error(t, err)
	assert.Equal(t, ErrBacktestNotFound, err)
}
