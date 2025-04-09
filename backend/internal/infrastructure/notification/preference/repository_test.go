package preference

import (
	"context"
	"testing"

	"go-crypto-bot-clean/backend/internal/domain/notification"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	return db, mock
}

func TestPreferenceRepository(t *testing.T) {
	db, mock := setupMockDB(t)
	logger, _ := zap.NewDevelopment()
	repo := NewRepository(db, logger)

	// Mock the AutoMigrate call
	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS "notification_preferences"`).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// Test GetUserPreferences
	ctx := context.Background()
	userID := "user123"

	// Mock the query
	mock.ExpectQuery(`SELECT (.+) FROM "notification_preferences" WHERE user_id = \$1 AND enabled = \$2`).
		WithArgs(userID, true).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "user_id", "channel", "recipient", "enabled"}).
			AddRow(1, "2023-01-01", "2023-01-01", nil, userID, "telegram", "123456789", true).
			AddRow(2, "2023-01-01", "2023-01-01", nil, userID, "slack", "#alerts", true))

	prefs, err := repo.GetUserPreferences(ctx, userID)
	require.NoError(t, err)
	assert.Len(t, prefs, 2)
	assert.Equal(t, userID, prefs[0].UserID)
	assert.Equal(t, "telegram", prefs[0].Channel)
	assert.Equal(t, "123456789", prefs[0].Recipient)
	assert.True(t, prefs[0].Enabled)
	assert.Equal(t, userID, prefs[1].UserID)
	assert.Equal(t, "slack", prefs[1].Channel)
	assert.Equal(t, "#alerts", prefs[1].Recipient)
	assert.True(t, prefs[1].Enabled)

	// Test SaveUserPreference - Update existing
	pref := notification.Preference{
		UserID:    userID,
		Channel:   "telegram",
		Recipient: "987654321",
		Enabled:   true,
	}

	// Mock the query to find existing preference
	mock.ExpectQuery(`SELECT (.+) FROM "notification_preferences" WHERE user_id = \$1 AND channel = \$2`).
		WithArgs(userID, "telegram").
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "user_id", "channel", "recipient", "enabled"}).
			AddRow(1, "2023-01-01", "2023-01-01", nil, userID, "telegram", "123456789", true))

	// Mock the update
	mock.ExpectExec(`UPDATE "notification_preferences" SET (.+) WHERE "id" = \$1`).
		WithArgs(1, sqlmock.AnyArg(), sqlmock.AnyArg(), nil, userID, "telegram", "987654321", true).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.SaveUserPreference(ctx, pref)
	require.NoError(t, err)

	// Test SaveUserPreference - Create new
	pref = notification.Preference{
		UserID:    userID,
		Channel:   "email",
		Recipient: "user@example.com",
		Enabled:   true,
	}

	// Mock the query to find existing preference (not found)
	mock.ExpectQuery(`SELECT (.+) FROM "notification_preferences" WHERE user_id = \$1 AND channel = \$2`).
		WithArgs(userID, "email").
		WillReturnError(gorm.ErrRecordNotFound)

	// Mock the insert
	mock.ExpectExec(`INSERT INTO "notification_preferences"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), nil, userID, "email", "user@example.com", true).
		WillReturnResult(sqlmock.NewResult(3, 1))

	err = repo.SaveUserPreference(ctx, pref)
	require.NoError(t, err)

	// Test DeleteUserPreference
	mock.ExpectExec(`DELETE FROM "notification_preferences" WHERE user_id = \$1 AND channel = \$2`).
		WithArgs(userID, "email").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.DeleteUserPreference(ctx, userID, "email")
	require.NoError(t, err)

	// Test GetAllUserPreferences
	mock.ExpectQuery(`SELECT (.+) FROM "notification_preferences" WHERE user_id = \$1`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "user_id", "channel", "recipient", "enabled"}).
			AddRow(1, "2023-01-01", "2023-01-01", nil, userID, "telegram", "987654321", true).
			AddRow(2, "2023-01-01", "2023-01-01", nil, userID, "slack", "#alerts", true).
			AddRow(4, "2023-01-01", "2023-01-01", nil, userID, "sms", "+1234567890", false))

	allPrefs, err := repo.GetAllUserPreferences(ctx, userID)
	require.NoError(t, err)
	assert.Len(t, allPrefs, 3)
	assert.Equal(t, userID, allPrefs[0].UserID)
	assert.Equal(t, "telegram", allPrefs[0].Channel)
	assert.Equal(t, "987654321", allPrefs[0].Recipient)
	assert.True(t, allPrefs[0].Enabled)
	assert.Equal(t, userID, allPrefs[1].UserID)
	assert.Equal(t, "slack", allPrefs[1].Channel)
	assert.Equal(t, "#alerts", allPrefs[1].Recipient)
	assert.True(t, allPrefs[1].Enabled)
	assert.Equal(t, userID, allPrefs[2].UserID)
	assert.Equal(t, "sms", allPrefs[2].Channel)
	assert.Equal(t, "+1234567890", allPrefs[2].Recipient)
	assert.False(t, allPrefs[2].Enabled)
}
