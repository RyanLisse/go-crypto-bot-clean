package factory

import (
	"os"
	"testing"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/config"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnectTurso_LocalOnly(t *testing.T) {
	// Create a test logger
	logger := zerolog.New(zerolog.NewTestWriter(t))

	// Create a temporary directory for the test database
	tempDir, err := os.MkdirTemp("", "turso-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a config with local database only
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Type:                   "turso",
			DSN:                    "file:" + tempDir + "/test.db",
			MaxIdleConns:           5,
			MaxOpenConns:           10,
			ConnMaxLifetimeMinutes: 60,
			TursoURL:               "", // Empty URL to force local-only mode
			AuthToken:              "",
		},
	}

	// Create a GORM config
	gormConfig := createGormConfig(cfg, &logger)

	// Connect to the database
	db, err := connectTurso(cfg, gormConfig, &logger)
	require.NoError(t, err)
	require.NotNil(t, db)

	// Verify the connection works
	sqlDB, err := db.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Ping())

	// Close the connection
	require.NoError(t, sqlDB.Close())
}

func TestOpenLocalDatabase(t *testing.T) {
	// Create a test logger
	logger := zerolog.New(zerolog.NewTestWriter(t))

	// Create a temporary directory for the test database
	tempDir, err := os.MkdirTemp("", "turso-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a config
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			MaxIdleConns:           5,
			MaxOpenConns:           10,
			ConnMaxLifetimeMinutes: 60,
		},
	}

	// Create a GORM config
	gormConfig := createGormConfig(cfg, &logger)

	// Open a local database
	dbPath := tempDir + "/test.db"
	db, err := openLocalDatabase(dbPath, gormConfig, cfg, &logger)
	require.NoError(t, err)
	require.NotNil(t, db)

	// Verify the connection works
	sqlDB, err := db.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Ping())

	// Close the connection
	require.NoError(t, sqlDB.Close())
}

func TestSetupPeriodicSync(t *testing.T) {
	// Create a test logger
	logger := zerolog.New(zerolog.NewTestWriter(t))

	// Create a config
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			MaxIdleConns:           5,
			MaxOpenConns:           10,
			ConnMaxLifetimeMinutes: 60,
		},
	}

	// Set environment variables for testing
	os.Setenv("TURSO_SYNC_ENABLED", "false")
	defer os.Unsetenv("TURSO_SYNC_ENABLED")

	// Create a mock database
	db, err := openLocalDatabase(":memory:", createGormConfig(cfg, &logger), cfg, &logger)
	require.NoError(t, err)
	require.NotNil(t, db)

	// Test with sync disabled
	setupPeriodicSync(db, cfg, &logger)

	// Set environment variables for testing
	os.Setenv("TURSO_SYNC_ENABLED", "true")
	os.Setenv("TURSO_SYNC_INTERVAL_SECONDS", "1")
	defer os.Unsetenv("TURSO_SYNC_INTERVAL_SECONDS")

	// Test with sync enabled but short interval
	// This is just to verify it doesn't crash, not to test the actual sync
	setupPeriodicSync(db, cfg, &logger)

	// Wait a moment to ensure the goroutine starts
	time.Sleep(10 * time.Millisecond)

	// Close the connection
	sqlDB, err := db.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())
}

func TestMaskDSN(t *testing.T) {
	testCases := []struct {
		name     string
		dsn      string
		expected string
	}{
		{
			name:     "SQLite file path",
			dsn:      "file:/path/to/database.db",
			expected: "file:/path/to/database.db",
		},
		{
			name:     "Memory database",
			dsn:      ":memory:",
			expected: ":memory:",
		},
		{
			name:     "Turso URL with auth token",
			dsn:      "libsql://my-db.turso.io?authToken=secret123",
			expected: "libsql://my-db.turso.io?authToken=***",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			masked := maskDSN(tc.dsn)
			assert.Equal(t, tc.expected, masked)
		})
	}
}
