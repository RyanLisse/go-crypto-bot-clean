package repo_test

import (
	"context"
	"testing"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm/repo"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupRiskProfileTestDB creates a new in-memory SQLite database for testing
func setupRiskProfileTestDB(t *testing.T) (*gorm.DB, func()) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err, "Failed to connect to in-memory database")

	// Create risk_profiles table
	err = db.Exec(`
		CREATE TABLE IF NOT EXISTS risk_profiles (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL UNIQUE,
			max_position_size REAL NOT NULL,
			max_total_exposure REAL NOT NULL,
			max_drawdown REAL NOT NULL,
			max_leverage REAL NOT NULL,
			max_concentration REAL NOT NULL,
			min_liquidity REAL NOT NULL,
			volatility_threshold REAL NOT NULL,
			daily_loss_limit REAL NOT NULL,
			weekly_loss_limit REAL NOT NULL,
			enable_auto_risk_control BOOLEAN NOT NULL,
			enable_notifications BOOLEAN NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`).Error
	require.NoError(t, err, "Failed to create risk_profiles table")

	// Create index for user_id
	err = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_risk_profiles_user_id ON risk_profiles(user_id)").Error
	require.NoError(t, err, "Failed to create index")

	// Return DB and cleanup function
	sqlDB, err := db.DB()
	require.NoError(t, err)

	return db, func() {
		sqlDB.Close()
	}
}

func TestGormRiskProfileRepository(t *testing.T) {
	db, cleanup := setupRiskProfileTestDB(t)
	defer cleanup()

	repository := repo.NewGormRiskProfileRepository(db)
	ctx := context.Background()

	t.Run("Save and GetByUserID", func(t *testing.T) {
		// Create a risk profile
		userID := "user123"
		profile := model.RiskProfile{
			ID:                    uuid.New().String(),
			UserID:                userID,
			MaxPositionSize:       2000.0,
			MaxTotalExposure:      10000.0,
			MaxDrawdown:           0.15,
			MaxLeverage:           5.0,
			MaxConcentration:      0.25,
			MinLiquidity:          20000.0,
			VolatilityThreshold:   0.08,
			DailyLossLimit:        200.0,
			WeeklyLossLimit:       1000.0,
			EnableAutoRiskControl: true,
			EnableNotifications:   true,
			CreatedAt:             time.Now(),
			UpdatedAt:             time.Now(),
		}

		// Save to database
		err := repository.Save(ctx, &profile)
		require.NoError(t, err)

		// Get by user ID
		retrieved, err := repository.GetByUserID(ctx, userID)
		require.NoError(t, err)
		assert.Equal(t, profile.ID, retrieved.ID)
		assert.Equal(t, profile.UserID, retrieved.UserID)
		assert.Equal(t, profile.MaxPositionSize, retrieved.MaxPositionSize)
		assert.Equal(t, profile.MaxTotalExposure, retrieved.MaxTotalExposure)
		assert.Equal(t, profile.MaxDrawdown, retrieved.MaxDrawdown)
		assert.Equal(t, profile.MaxLeverage, retrieved.MaxLeverage)
		assert.Equal(t, profile.MaxConcentration, retrieved.MaxConcentration)
		assert.Equal(t, profile.MinLiquidity, retrieved.MinLiquidity)
		assert.Equal(t, profile.VolatilityThreshold, retrieved.VolatilityThreshold)
		assert.Equal(t, profile.DailyLossLimit, retrieved.DailyLossLimit)
		assert.Equal(t, profile.WeeklyLossLimit, retrieved.WeeklyLossLimit)
		assert.Equal(t, profile.EnableAutoRiskControl, retrieved.EnableAutoRiskControl)
		assert.Equal(t, profile.EnableNotifications, retrieved.EnableNotifications)
	})

	t.Run("Update", func(t *testing.T) {
		// Create a risk profile
		userID := "user456"
		profile := model.RiskProfile{
			ID:                    uuid.New().String(),
			UserID:                userID,
			MaxPositionSize:       1000.0,
			MaxTotalExposure:      5000.0,
			MaxDrawdown:           0.1,
			MaxLeverage:           3.0,
			MaxConcentration:      0.2,
			MinLiquidity:          10000.0,
			VolatilityThreshold:   0.05,
			DailyLossLimit:        100.0,
			WeeklyLossLimit:       500.0,
			EnableAutoRiskControl: true,
			EnableNotifications:   true,
			CreatedAt:             time.Now(),
			UpdatedAt:             time.Now(),
		}

		// Save to database
		err := repository.Save(ctx, &profile)
		require.NoError(t, err)

		// Update profile
		profile.MaxPositionSize = 1500.0
		profile.EnableAutoRiskControl = false
		err = repository.Save(ctx, &profile)
		require.NoError(t, err)

		// Get updated profile
		retrieved, err := repository.GetByUserID(ctx, userID)
		require.NoError(t, err)
		assert.Equal(t, 1500.0, retrieved.MaxPositionSize)
		assert.False(t, retrieved.EnableAutoRiskControl)
	})

	t.Run("GetDefault", func(t *testing.T) {
		// Get profile for non-existent user
		userID := "non_existent_user"
		profile, err := repository.GetByUserID(ctx, userID)
		require.NoError(t, err)
		assert.NotEmpty(t, profile.ID)
		assert.Equal(t, userID, profile.UserID)
		assert.Equal(t, 1000.0, profile.MaxPositionSize)
		assert.Equal(t, 5000.0, profile.MaxTotalExposure)
		assert.Equal(t, 0.1, profile.MaxDrawdown)
		assert.Equal(t, 3.0, profile.MaxLeverage)
		assert.Equal(t, 0.2, profile.MaxConcentration)
		assert.Equal(t, 10000.0, profile.MinLiquidity)
		assert.Equal(t, 0.05, profile.VolatilityThreshold)
		assert.Equal(t, 100.0, profile.DailyLossLimit)
		assert.Equal(t, 500.0, profile.WeeklyLossLimit)
		assert.True(t, profile.EnableAutoRiskControl)
		assert.True(t, profile.EnableNotifications)
	})

	t.Run("Delete", func(t *testing.T) {
		// Create a risk profile
		userID := "user789"
		profile := &model.RiskProfile{
			ID:                    uuid.New().String(),
			UserID:                userID,
			MaxPositionSize:       1000.0,
			MaxTotalExposure:      5000.0,
			MaxDrawdown:           0.1,
			MaxLeverage:           3.0,
			MaxConcentration:      0.2,
			MinLiquidity:          10000.0,
			VolatilityThreshold:   0.05,
			DailyLossLimit:        100.0,
			WeeklyLossLimit:       500.0,
			EnableAutoRiskControl: true,
			EnableNotifications:   true,
			CreatedAt:             time.Now(),
			UpdatedAt:             time.Now(),
		}

		// Save to database
		err := repository.Save(ctx, profile)
		require.NoError(t, err)

		// Delete profile
		err = repository.Delete(ctx, profile.ID)
		require.NoError(t, err)

		// Verify profile is deleted (should get default profile)
		retrieved, err := repository.GetByUserID(ctx, userID)
		require.NoError(t, err)
		assert.NotEqual(t, profile.ID, retrieved.ID) // Should be a new ID
		assert.Equal(t, userID, retrieved.UserID)
	})
}
