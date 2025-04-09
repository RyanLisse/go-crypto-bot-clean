package gorm

import (
	"context"
	"testing"

	notification_domain "go-crypto-bot-clean/backend/internal/domain/notification"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestNotificationPreferenceRepo_Upsert(t *testing.T) {
	// Setup in-memory SQLite database
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)

	// Run migrations
	err = db.AutoMigrate(&PreferenceGORM{})
	require.NoError(t, err)

	// Create repository
	repo := NewGormNotificationPreferenceRepository(db)

	// Test data
	testUserID := "test-user-1"
	testChannel := "telegram"
	enabled := true

	// Create a domain preference
	pref := notification_domain.Preference{
		UserID:    testUserID,
		Channel:   testChannel,
		Recipient: "12345",
		Enabled:   enabled,
	}

	// Test saving preference
	ctx := context.Background()
	err = repo.SaveUserPreference(ctx, pref)
	require.NoError(t, err)

	// Verify saved preference
	var result PreferenceGORM
	err = db.Where("user_id = ? AND channel = ?", testUserID, testChannel).First(&result).Error
	require.NoError(t, err)

	// Verify fields
	assert.Equal(t, testUserID, result.UserID)
	assert.Equal(t, testChannel, result.Channel)
	assert.Equal(t, "12345", result.Recipient)
	assert.NotNil(t, result.Enabled)
	assert.True(t, *result.Enabled)

	// Test updating preference
	pref.Enabled = false
	err = repo.SaveUserPreference(ctx, pref)
	require.NoError(t, err)

	// Verify updated preference
	err = db.Where("user_id = ? AND channel = ?", testUserID, testChannel).First(&result).Error
	require.NoError(t, err)
	assert.NotNil(t, result.Enabled)
	assert.False(t, *result.Enabled)
}

func setupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&PreferenceGORM{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestNotificationPreferenceRepository(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&PreferenceGORM{})
	require.NoError(t, err)

	repo := NewGormNotificationPreferenceRepository(db)
	ctx := context.Background()

	testUserID := "test-user"

	// Save some preferences
	pref1 := notification_domain.Preference{UserID: testUserID, Channel: "email", Recipient: "test@example.com", Enabled: true}
	pref2 := notification_domain.Preference{UserID: testUserID, Channel: "sms", Recipient: "+123456", Enabled: false}
	pref3 := notification_domain.Preference{UserID: "other-user", Channel: "email", Recipient: "other@example.com", Enabled: true}

	require.NoError(t, repo.SaveUserPreference(ctx, pref1))
	require.NoError(t, repo.SaveUserPreference(ctx, pref2))
	require.NoError(t, repo.SaveUserPreference(ctx, pref3))

	t.Run("GetUserPreferences - Existing User with Enabled Prefs", func(t *testing.T) {
		prefs, err := repo.GetUserPreferences(ctx, testUserID)
		require.NoError(t, err)
		assert.Len(t, prefs, 1) // Only enabled prefs should be returned
		assert.Equal(t, pref1.Channel, prefs[0].Channel)
		assert.Equal(t, pref1.Recipient, prefs[0].Recipient)
		assert.True(t, prefs[0].Enabled)
	})

	t.Run("GetUserPreferences - User with Only Disabled Prefs", func(t *testing.T) {
		// Temporarily disable pref1 for this user
		pref1Disabled := pref1
		pref1Disabled.Enabled = false
		require.NoError(t, repo.SaveUserPreference(ctx, pref1Disabled))

		prefs, err := repo.GetUserPreferences(ctx, testUserID)
		require.NoError(t, err)
		assert.Empty(t, prefs) // No enabled prefs

		// Re-enable pref1 for other tests
		require.NoError(t, repo.SaveUserPreference(ctx, pref1))
	})

	t.Run("GetUserPreferences - Non-Existent User", func(t *testing.T) {
		prefs, err := repo.GetUserPreferences(ctx, "non-existent-user")
		require.NoError(t, err)
		assert.Empty(t, prefs)
	})

	// --- Test SaveUserPreference (Upsert Logic) ---
	t.Run("SaveUserPreference - Update Existing", func(t *testing.T) {
		updatedPref1 := notification_domain.Preference{UserID: testUserID, Channel: "email", Recipient: "updated@example.com", Enabled: false}
		err := repo.SaveUserPreference(ctx, updatedPref1)
		require.NoError(t, err)

		// Verify update using GetUserPreferences (should be empty as it's disabled)
		prefs, err := repo.GetUserPreferences(ctx, testUserID)
		require.NoError(t, err)
		assert.Empty(t, prefs)

		// Restore original pref1 state
		require.NoError(t, repo.SaveUserPreference(ctx, pref1))
	})

	t.Run("SaveUserPreference - Insert New", func(t *testing.T) {
		newPref := notification_domain.Preference{UserID: testUserID, Channel: "slack", Recipient: "slack-id", Enabled: true}
		err := repo.SaveUserPreference(ctx, newPref)
		require.NoError(t, err)

		prefs, err := repo.GetUserPreferences(ctx, testUserID)
		require.NoError(t, err)
		assert.Len(t, prefs, 2) // Should now have email and slack enabled

		foundSlack := false
		for _, p := range prefs {
			if p.Channel == "slack" {
				assert.Equal(t, newPref.Recipient, p.Recipient)
				assert.True(t, p.Enabled)
				foundSlack = true
				break
			}
		}
		assert.True(t, foundSlack, "Newly inserted slack preference not found or not enabled")
	})
}
