package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"go-crypto-bot-clean/backend/api/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// Create a temporary database file
	dbFile, err := os.CreateTemp("", "test-*.db")
	require.NoError(t, err)
	dbFile.Close()

	// Open database connection
	db, err := gorm.Open(sqlite.Open(dbFile.Name()), &gorm.Config{})
	require.NoError(t, err)

	// Run migrations
	err = db.AutoMigrate(
		// User models
		&models.User{},
		&models.UserRole{},
		&models.UserSettings{},
		&models.RefreshToken{},

		// Strategy models
		&models.Strategy{},
		&models.StrategyParameter{},
		&models.StrategyPerformance{},

		// Backtest models
		&models.Backtest{},
		&models.BacktestTrade{},
		&models.BacktestEquity{},
	)
	require.NoError(t, err)

	// Clean up the database file when the test is done
	t.Cleanup(func() {
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.Close()
		}
		os.Remove(dbFile.Name())
	})

	return db
}

func TestUserRepository(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormUserRepository(db)
	ctx := context.Background()

	// Test creating a user
	user := &models.User{
		ID:           "test-user",
		Email:        "test@example.com",
		Username:     "testuser",
		PasswordHash: "password-hash",
		FirstName:    "Test",
		LastName:     "User",
	}

	err := repo.Create(ctx, user)
	assert.NoError(t, err)

	// Test getting a user by ID
	retrievedUser, err := repo.GetByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, retrievedUser.ID)
	assert.Equal(t, user.Email, retrievedUser.Email)
	assert.Equal(t, user.Username, retrievedUser.Username)

	// Test getting a user by email
	retrievedUser, err = repo.GetByEmail(ctx, user.Email)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, retrievedUser.ID)
	assert.Equal(t, user.Email, retrievedUser.Email)

	// Test getting a user by username
	retrievedUser, err = repo.GetByUsername(ctx, user.Username)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, retrievedUser.ID)
	assert.Equal(t, user.Username, retrievedUser.Username)

	// Test updating a user
	user.FirstName = "Updated"
	user.LastName = "Name"
	err = repo.Update(ctx, user)
	assert.NoError(t, err)

	retrievedUser, err = repo.GetByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated", retrievedUser.FirstName)
	assert.Equal(t, "Name", retrievedUser.LastName)

	// Test adding a role
	err = repo.AddRole(ctx, user.ID, "admin")
	assert.NoError(t, err)

	// Test getting roles
	roles, err := repo.GetRoles(ctx, user.ID)
	assert.NoError(t, err)
	assert.Contains(t, roles, "admin")

	// Test removing a role
	err = repo.RemoveRole(ctx, user.ID, "admin")
	assert.NoError(t, err)

	roles, err = repo.GetRoles(ctx, user.ID)
	assert.NoError(t, err)
	assert.NotContains(t, roles, "admin")

	// Test getting and updating settings
	settings, err := repo.GetSettings(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, settings.UserID)
	assert.Equal(t, "light", settings.Theme)

	settings.Theme = "dark"
	err = repo.UpdateSettings(ctx, settings)
	assert.NoError(t, err)

	updatedSettings, err := repo.GetSettings(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "dark", updatedSettings.Theme)

	// Test saving a refresh token
	token := &models.RefreshToken{
		ID:        "token-id",
		UserID:    user.ID,
		Token:     "refresh-token",
		ExpiresAt: time.Now().Add(time.Hour),
	}
	err = repo.SaveRefreshToken(ctx, token)
	assert.NoError(t, err)

	// Test getting a refresh token
	retrievedToken, err := repo.GetRefreshToken(ctx, token.Token)
	assert.NoError(t, err)
	assert.Equal(t, token.ID, retrievedToken.ID)
	assert.Equal(t, token.UserID, retrievedToken.UserID)

	// Test revoking a refresh token
	err = repo.RevokeRefreshToken(ctx, token.Token)
	assert.NoError(t, err)

	_, err = repo.GetRefreshToken(ctx, token.Token)
	assert.Error(t, err)

	// Test updating last login
	err = repo.UpdateLastLogin(ctx, user.ID)
	assert.NoError(t, err)

	retrievedUser, err = repo.GetByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedUser.LastLoginAt)

	// Test deleting a user
	err = repo.Delete(ctx, user.ID)
	assert.NoError(t, err)

	_, err = repo.GetByID(ctx, user.ID)
	assert.Error(t, err)
	assert.Equal(t, ErrUserNotFound, err)
}
