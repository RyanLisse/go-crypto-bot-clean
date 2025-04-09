package persistence_test

import (
	"context"
	"testing"

	"go-crypto-bot-clean/backend/internal/domain/notification"
	"go-crypto-bot-clean/backend/internal/domain/notification/ports"
	gorm_persistence "go-crypto-bot-clean/backend/internal/infrastructure/persistence/gorm" // Alias for GORM implementation package
	"go-crypto-bot-clean/backend/internal/infrastructure/persistence/memory"                // Use memory implementation for initial test structure
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Test Suite Runner for NotificationPreferenceRepository implementations
func testNotificationPreferenceRepository(t *testing.T, setupRepo func() (ports.NotificationPreferenceRepository, func())) {
	t.Run("GetUserPreferences", func(t *testing.T) {
		t.Run("should return empty slice when no preferences exist", func(t *testing.T) {
			repo, teardown := setupRepo()
			defer teardown()

			ctx := context.Background()
			userID := "user-no-prefs"

			prefs, err := repo.GetUserPreferences(ctx, userID)

			require.NoError(t, err)
			assert.Empty(t, prefs)
		})

		t.Run("should return only enabled preferences for the user", func(t *testing.T) {
			repo, teardown := setupRepo()
			defer teardown()

			ctx := context.Background()
			userID := "user-with-prefs"

			// Seed data (specific to implementation, need setup/teardown or interface)
			// This part highlights the need for a seeding mechanism within setupRepo
			// For now, assuming the interface might have a Save method for testing setup
			saver, ok := repo.(interface { // Type assert for a hypothetical Save method for testing
				SaveUserPreference(context.Context, notification.Preference) error
			})
			if !ok {
				t.Skip("Skipping test: Repository implementation does not support seeding via SaveUserPreference")
			}

			pref1 := notification.Preference{UserID: userID, Channel: "telegram", Recipient: "12345", Enabled: true}
			pref2 := notification.Preference{UserID: userID, Channel: "slack", Recipient: "U67890", Enabled: false}
			pref3 := notification.Preference{UserID: userID, Channel: "email", Recipient: "test@example.com", Enabled: true}
			prefOtherUser := notification.Preference{UserID: "other-user", Channel: "telegram", Recipient: "99999", Enabled: true}

			require.NoError(t, saver.SaveUserPreference(ctx, pref1))
			require.NoError(t, saver.SaveUserPreference(ctx, pref2))
			require.NoError(t, saver.SaveUserPreference(ctx, pref3))
			require.NoError(t, saver.SaveUserPreference(ctx, prefOtherUser))

			prefs, err := repo.GetUserPreferences(ctx, userID)

			require.NoError(t, err)
			assert.Len(t, prefs, 2) // Only pref1 and pref3 should be returned
			assert.Contains(t, prefs, pref1)
			assert.NotContains(t, prefs, pref2)
			assert.Contains(t, prefs, pref3)
		})

		t.Run("should return preferences correctly across multiple users", func(t *testing.T) {
			repo, teardown := setupRepo()
			defer teardown()

			ctx := context.Background()
			userID1 := "user-multi-1"
			userID2 := "user-multi-2"

			saver, ok := repo.(interface {
				SaveUserPreference(context.Context, notification.Preference) error
			})
			if !ok {
				t.Skip("Skipping test: Repository implementation does not support seeding via SaveUserPreference")
			}

			pref1_1 := notification.Preference{UserID: userID1, Channel: "telegram", Recipient: "111", Enabled: true}
			pref1_2 := notification.Preference{UserID: userID1, Channel: "slack", Recipient: "S111", Enabled: true}
			pref2_1 := notification.Preference{UserID: userID2, Channel: "telegram", Recipient: "222", Enabled: true}
			pref2_2 := notification.Preference{UserID: userID2, Channel: "email", Recipient: "user2@example.com", Enabled: false} // Disabled

			require.NoError(t, saver.SaveUserPreference(ctx, pref1_1))
			require.NoError(t, saver.SaveUserPreference(ctx, pref1_2))
			require.NoError(t, saver.SaveUserPreference(ctx, pref2_1))
			require.NoError(t, saver.SaveUserPreference(ctx, pref2_2))

			prefs1, err1 := repo.GetUserPreferences(ctx, userID1)
			prefs2, err2 := repo.GetUserPreferences(ctx, userID2)

			require.NoError(t, err1)
			require.NoError(t, err2)
			assert.Len(t, prefs1, 2)
			assert.Contains(t, prefs1, pref1_1)
			assert.Contains(t, prefs1, pref1_2)

			assert.Len(t, prefs2, 1)
			assert.Contains(t, prefs2, pref2_1)
		})

		// TODO: Add test case for context cancellation if applicable to the implementation
	})

	// TODO: Add test suites for SaveUserPreference, DeleteUserPreference when they are added to the interface
}

// TestInMemoryNotificationPreferenceRepository runs the test suite for the memory implementation.
func TestInMemoryNotificationPreferenceRepository(t *testing.T) {
	// Setup function for the in-memory repository
	setupFunc := func() (ports.NotificationPreferenceRepository, func()) {
		repo := memory.NewInMemoryNotificationPreferenceRepository() // Assumes this constructor exists
		teardown := func() {
			// No specific teardown needed for this simple in-memory repo
		}
		return repo, teardown
	}

	testNotificationPreferenceRepository(t, setupFunc)
}

// TestGormNotificationPreferenceRepository runs the test suite for the GORM implementation.
func TestGormNotificationPreferenceRepository(t *testing.T) {
	// Setup function for the GORM repository using in-memory SQLite
	setupFunc := func() (ports.NotificationPreferenceRepository, func()) {
		// Use in-memory SQLite for testing: "file::memory:?cache=shared"
		// Alternatively, use a temporary file: "test_gorm_prefs.db"
		db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
		if err != nil {
			t.Fatalf("Failed to connect to in-memory sqlite: %v", err)
		}

		// Run migrations (assuming the Preference model is defined in gorm_persistence)
		// Drop table first to ensure clean state for each test setup with shared cache
		db.Migrator().DropTable(&gorm_persistence.PreferenceGORM{})
		err = db.AutoMigrate(&gorm_persistence.PreferenceGORM{}) // Assumes PreferenceGORM model exists
		if err != nil {
			t.Fatalf("Failed to auto-migrate PreferenceGORM model: %v", err)
		}

		repo := gorm_persistence.NewGormNotificationPreferenceRepository(db) // Assumes constructor exists

		teardown := func() {
			// Drop the table after tests are done to clean up
			db.Migrator().DropTable(&gorm_persistence.PreferenceGORM{})
			// No explicit close needed for in-memory, but good practice if using file DB
			// sqlDB, _ := db.DB()
			// if sqlDB != nil { sqlDB.Close() }
			// If using a file DB, os.Remove("test_gorm_prefs.db")
		}
		return repo, teardown
	}

	testNotificationPreferenceRepository(t, setupFunc)
}
