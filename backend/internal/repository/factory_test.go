package repository

import (
	"testing"
	"time"

	"go-crypto-bot-clean/backend/internal/repository/balance"
	"go-crypto-bot-clean/backend/internal/repository/database"
	"go-crypto-bot-clean/backend/internal/repository/report"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestFactory(t *testing.T) {
	// Create a test logger
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	// Create a test database
	dbConfig := database.Config{
		DatabasePath:    ":memory:", // Use in-memory SQLite for testing
		MaxOpenConns:    5,
		MaxIdleConns:    2,
		ConnMaxLifetime: 5 * time.Minute,
	}

	// Create SQLite repository
	db, err := database.NewSQLiteRepository(dbConfig)
	require.NoError(t, err)
	defer db.Close()

	// Create factory
	factory := NewFactory(db, logger)

	// Test creating balance history repository
	t.Run("NewBalanceHistoryRepository", func(t *testing.T) {
		repo := factory.NewBalanceHistoryRepository()
		assert.NotNil(t, repo)
		assert.IsType(t, &balance.BalanceHistoryRepository{}, repo)
	})

	// Test creating report repository
	t.Run("NewReportRepository", func(t *testing.T) {
		repo := factory.NewReportRepository()
		assert.NotNil(t, repo)
		assert.IsType(t, &report.ReportRepository{}, repo)
	})
}
