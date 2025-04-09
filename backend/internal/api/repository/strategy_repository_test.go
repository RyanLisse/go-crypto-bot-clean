package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go-crypto-bot-clean/backend/internal/api/models"
)

func TestStrategyRepository(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormStrategyRepository(db)
	ctx := context.Background()

	// Create a test user first
	userRepo := NewGormUserRepository(db)
	user := &models.User{
		ID:           "test-user",
		Email:        "test@example.com",
		Username:     "testuser",
		PasswordHash: "password-hash",
	}
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	// Test creating a strategy
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

	err = repo.Create(ctx, strategy)
	assert.NoError(t, err)

	// Test getting a strategy by ID
	retrievedStrategy, err := repo.GetByID(ctx, strategy.ID)
	assert.NoError(t, err)
	assert.Equal(t, strategy.ID, retrievedStrategy.ID)
	assert.Equal(t, strategy.Name, retrievedStrategy.Name)
	assert.Equal(t, strategy.Description, retrievedStrategy.Description)
	assert.Equal(t, strategy.IsEnabled, retrievedStrategy.IsEnabled)
	assert.Equal(t, strategy.UserID, retrievedStrategy.UserID)
	assert.Equal(t, "value1", retrievedStrategy.Parameters["param1"])
	assert.Equal(t, float64(42), retrievedStrategy.Parameters["param2"])

	// Test getting strategies by user ID
	strategies, err := repo.GetByUserID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Len(t, strategies, 1)
	assert.Equal(t, strategy.ID, strategies[0].ID)

	// Test updating a strategy
	strategy.Name = "Updated Strategy"
	strategy.Description = "An updated test strategy"
	strategy.Parameters["param3"] = "value3"
	err = repo.Update(ctx, strategy)
	assert.NoError(t, err)

	retrievedStrategy, err = repo.GetByID(ctx, strategy.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Strategy", retrievedStrategy.Name)
	assert.Equal(t, "An updated test strategy", retrievedStrategy.Description)
	assert.Equal(t, "value3", retrievedStrategy.Parameters["param3"])

	// Test adding a parameter
	param := &models.StrategyParameter{
		StrategyID:  strategy.ID,
		Name:        "testParam",
		Type:        "string",
		Description: "A test parameter",
		Default:     "default",
		Required:    true,
	}
	err = repo.AddParameter(ctx, param)
	assert.NoError(t, err)

	// Test getting parameters
	params, err := repo.GetParameters(ctx, strategy.ID)
	assert.NoError(t, err)
	assert.Len(t, params, 1)
	assert.Equal(t, "testParam", params[0].Name)
	assert.Equal(t, "string", params[0].Type)
	assert.Equal(t, "A test parameter", params[0].Description)
	assert.Equal(t, "default", params[0].Default)
	assert.Equal(t, true, params[0].Required)

	// Test updating a parameter
	param.Description = "An updated parameter"
	err = repo.UpdateParameter(ctx, param)
	assert.NoError(t, err)

	params, err = repo.GetParameters(ctx, strategy.ID)
	assert.NoError(t, err)
	assert.Equal(t, "An updated parameter", params[0].Description)

	// Test adding performance metrics
	performance := &models.StrategyPerformance{
		StrategyID:   strategy.ID,
		WinRate:      0.65,
		ProfitFactor: 1.5,
		SharpeRatio:  1.2,
		MaxDrawdown:  0.1,
		TotalTrades:  100,
		PeriodStart:  time.Now().AddDate(0, -1, 0),
		PeriodEnd:    time.Now(),
	}
	err = repo.AddPerformance(ctx, performance)
	assert.NoError(t, err)

	// Test getting performance metrics
	performances, err := repo.GetPerformance(ctx, strategy.ID)
	assert.NoError(t, err)
	assert.Len(t, performances, 1)
	assert.Equal(t, 0.65, performances[0].WinRate)
	assert.Equal(t, 1.5, performances[0].ProfitFactor)
	assert.Equal(t, 1.2, performances[0].SharpeRatio)
	assert.Equal(t, 0.1, performances[0].MaxDrawdown)
	assert.Equal(t, 100, performances[0].TotalTrades)

	// Test deleting a parameter
	err = repo.DeleteParameter(ctx, param.ID)
	assert.NoError(t, err)

	params, err = repo.GetParameters(ctx, strategy.ID)
	assert.NoError(t, err)
	assert.Len(t, params, 0)

	// Test deleting a strategy
	err = repo.Delete(ctx, strategy.ID)
	assert.NoError(t, err)

	_, err = repo.GetByID(ctx, strategy.ID)
	assert.Error(t, err)
	assert.Equal(t, ErrStrategyNotFound, err)
}
