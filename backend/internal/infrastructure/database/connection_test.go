package database

import (
	"os"
	"testing"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnect(t *testing.T) {
	// Create a test config
	cfg := &config.Config{
		ENV: "test",
		Database: config.DatabaseConfig{
			Driver: "sqlite",
			Path:   ":memory:",
		},
	}

	// Create a logger
	logger := zerolog.New(os.Stdout).Level(zerolog.DebugLevel)

	// Connect to the database
	db, err := Connect(cfg, &logger)
	require.NoError(t, err)
	require.NotNil(t, db)

	// Verify the connection
	sqlDB, err := db.DB()
	require.NoError(t, err)
	require.NotNil(t, sqlDB)

	// Ping the database
	err = sqlDB.Ping()
	require.NoError(t, err)
}

func TestRunMigrations(t *testing.T) {
	// Create a test config
	cfg := &config.Config{
		ENV: "test",
		Database: config.DatabaseConfig{
			Driver: "sqlite",
			Path:   ":memory:",
		},
	}

	// Create a logger
	logger := zerolog.New(os.Stdout).Level(zerolog.DebugLevel)

	// Connect to the database
	db, err := Connect(cfg, &logger)
	require.NoError(t, err)
	require.NotNil(t, db)

	// Run migrations
	err = RunMigrations(db, &logger)
	assert.NoError(t, err)

	// Verify that tables were created
	var count int64
	err = db.Raw("SELECT count(*) FROM sqlite_master WHERE type='table'").Count(&count).Error
	require.NoError(t, err)
	assert.Greater(t, count, int64(0))
}
