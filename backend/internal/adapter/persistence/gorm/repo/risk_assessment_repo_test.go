package repo_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm/repo"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates a new in-memory SQLite database for testing
func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err, "Failed to connect to in-memory database")

	// Create risk_assessments table
	err = db.Exec(`
		CREATE TABLE IF NOT EXISTS risk_assessments (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			type TEXT NOT NULL, 
			level TEXT NOT NULL,
			status TEXT NOT NULL,
			symbol TEXT,
			position_id TEXT,
			order_id TEXT,
			score REAL,
			message TEXT NOT NULL,
			recommendation TEXT,
			metadata_json TEXT,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			resolved_at TIMESTAMP
		)
	`).Error
	require.NoError(t, err, "Failed to create risk_assessments table")

	// Create indices for faster searches
	indices := []string{
		"CREATE INDEX IF NOT EXISTS idx_risk_assessments_user_id ON risk_assessments(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_risk_assessments_type ON risk_assessments(type)",
		"CREATE INDEX IF NOT EXISTS idx_risk_assessments_level ON risk_assessments(level)",
		"CREATE INDEX IF NOT EXISTS idx_risk_assessments_status ON risk_assessments(status)",
		"CREATE INDEX IF NOT EXISTS idx_risk_assessments_symbol ON risk_assessments(symbol)",
		"CREATE INDEX IF NOT EXISTS idx_risk_assessments_created_at ON risk_assessments(created_at)",
	}

	for _, index := range indices {
		err = db.Exec(index).Error
		require.NoError(t, err, "Failed to create index")
	}

	// Return DB and cleanup function
	sqlDB, err := db.DB()
	require.NoError(t, err)

	return db, func() {
		sqlDB.Close()
	}
}

// getTestLogger creates a logger for testing
func getTestLogger() zerolog.Logger {
	return zerolog.New(os.Stdout).With().Timestamp().Logger()
}

func TestGormRiskAssessmentRepository(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repository := repo.NewGormRiskAssessmentRepository(db)
	ctx := context.Background()

	t.Run("Create and GetByID", func(t *testing.T) {
		// Create a risk assessment
		assessment := &model.RiskAssessment{
			ID:             uuid.New().String(),
			UserID:         "user123",
			Type:           model.RiskTypePosition,
			Level:          model.RiskLevelHigh,
			Status:         model.RiskStatusActive,
			Symbol:         "BTCUSDT",
			PositionID:     "pos123",
			Score:          75.0,
			Message:        "Position size exceeds recommended limit",
			Recommendation: "Consider reducing position size",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		// Save to database
		err := repository.Create(ctx, assessment)
		require.NoError(t, err)

		// Get by ID
		retrieved, err := repository.GetByID(ctx, assessment.ID)
		require.NoError(t, err)
		assert.Equal(t, assessment.ID, retrieved.ID)
		assert.Equal(t, assessment.UserID, retrieved.UserID)
		assert.Equal(t, assessment.Type, retrieved.Type)
		assert.Equal(t, assessment.Level, retrieved.Level)
		assert.Equal(t, assessment.Status, retrieved.Status)
		assert.Equal(t, assessment.Symbol, retrieved.Symbol)
		assert.Equal(t, assessment.PositionID, retrieved.PositionID)
		assert.Equal(t, assessment.Message, retrieved.Message)
		assert.Equal(t, assessment.Recommendation, retrieved.Recommendation)
	})

	t.Run("Update", func(t *testing.T) {
		// Create a risk assessment
		assessment := &model.RiskAssessment{
			ID:        uuid.New().String(),
			UserID:    "user123",
			Type:      model.RiskTypePosition,
			Level:     model.RiskLevelMedium,
			Status:    model.RiskStatusActive,
			Symbol:    "ETHUSDT",
			Message:   "Medium risk identified",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Save to database
		err := repository.Create(ctx, assessment)
		require.NoError(t, err)

		// Update assessment
		assessment.Level = model.RiskLevelHigh
		assessment.Message = "Risk level increased"
		assessment.Status = model.RiskStatusResolved
		resolvedTime := time.Now()
		assessment.ResolvedAt = &resolvedTime

		err = repository.Update(ctx, assessment)
		require.NoError(t, err)

		// Get updated assessment
		retrieved, err := repository.GetByID(ctx, assessment.ID)
		require.NoError(t, err)
		assert.Equal(t, model.RiskLevelHigh, retrieved.Level)
		assert.Equal(t, "Risk level increased", retrieved.Message)
		assert.Equal(t, model.RiskStatusResolved, retrieved.Status)
		assert.NotNil(t, retrieved.ResolvedAt)
	})

	t.Run("GetByUserID", func(t *testing.T) {
		// Create multiple risk assessments for the same user
		userID := "user456"
		for i := 0; i < 3; i++ {
			assessment := &model.RiskAssessment{
				ID:        uuid.New().String(),
				UserID:    userID,
				Type:      model.RiskTypeExposure,
				Level:     model.RiskLevelMedium,
				Status:    model.RiskStatusActive,
				Message:   "Test assessment",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			err := repository.Create(ctx, assessment)
			require.NoError(t, err)
		}

		// Retrieve with pagination
		assessments, err := repository.GetByUserID(ctx, userID, 2, 0)
		require.NoError(t, err)
		assert.Len(t, assessments, 2)

		// Get active assessments
		activeAssessments, err := repository.GetActiveByUserID(ctx, userID)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(activeAssessments), 3)
	})

	t.Run("GetByType", func(t *testing.T) {
		// Create assessments with different types
		types := []model.RiskType{
			model.RiskTypePosition,
			model.RiskTypeVolatility,
			model.RiskTypeLiquidity,
		}

		for _, riskType := range types {
			assessment := &model.RiskAssessment{
				ID:        uuid.New().String(),
				UserID:    "user789",
				Type:      riskType,
				Level:     model.RiskLevelMedium,
				Status:    model.RiskStatusActive,
				Message:   "Test assessment for " + string(riskType),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			err := repository.Create(ctx, assessment)
			require.NoError(t, err)
		}

		// Get by type
		liquidityAssessments, err := repository.GetByType(ctx, model.RiskTypeLiquidity, 10, 0)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(liquidityAssessments), 1)
		for _, assessment := range liquidityAssessments {
			assert.Equal(t, model.RiskTypeLiquidity, assessment.Type)
		}
	})

	t.Run("GetByLevel", func(t *testing.T) {
		// Create assessments with different levels
		userId := "user101"
		levels := []model.RiskLevel{
			model.RiskLevelLow,
			model.RiskLevelMedium,
			model.RiskLevelHigh,
			model.RiskLevelCritical,
		}

		for _, level := range levels {
			assessment := &model.RiskAssessment{
				ID:        uuid.New().String(),
				UserID:    userId,
				Type:      model.RiskTypePosition,
				Level:     level,
				Status:    model.RiskStatusActive,
				Message:   "Test assessment for " + string(level),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			err := repository.Create(ctx, assessment)
			require.NoError(t, err)
		}

		// Get by level
		highRiskAssessments, err := repository.GetByLevel(ctx, model.RiskLevelHigh, 10, 0)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(highRiskAssessments), 1)
		for _, assessment := range highRiskAssessments {
			assert.Equal(t, model.RiskLevelHigh, assessment.Level)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		// Create assessment
		assessment := &model.RiskAssessment{
			ID:        uuid.New().String(),
			UserID:    "user202",
			Type:      model.RiskTypePosition,
			Level:     model.RiskLevelMedium,
			Status:    model.RiskStatusActive,
			Message:   "Assessment to delete",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := repository.Create(ctx, assessment)
		require.NoError(t, err)

		// Delete assessment
		err = repository.Delete(ctx, assessment.ID)
		require.NoError(t, err)

		// Verify deletion
		_, err = repository.GetByID(ctx, assessment.ID)
		assert.Error(t, err) // Should get an error (not found)
	})

	t.Run("Resolve Method", func(t *testing.T) {
		// Create a risk assessment
		assessment := &model.RiskAssessment{
			ID:        uuid.New().String(),
			UserID:    "user123",
			Type:      model.RiskTypePosition,
			Level:     model.RiskLevelMedium,
			Status:    model.RiskStatusActive,
			Symbol:    "ETHUSDT",
			Message:   "Active risk assessment",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Save to database
		err := repository.Create(ctx, assessment)
		require.NoError(t, err)

		// Retrieve to verify initial state
		initialAssessment, err := repository.GetByID(ctx, assessment.ID)
		require.NoError(t, err)
		assert.Equal(t, model.RiskStatusActive, initialAssessment.Status)
		assert.Nil(t, initialAssessment.ResolvedAt)

		// Call the Resolve method
		initialAssessment.Resolve()

		// Update in the database
		err = repository.Update(ctx, initialAssessment)
		require.NoError(t, err)

		// Retrieve again to verify changes
		resolvedAssessment, err := repository.GetByID(ctx, assessment.ID)
		require.NoError(t, err)

		// Verify the status changed and ResolvedAt was set
		assert.Equal(t, model.RiskStatusResolved, resolvedAssessment.Status)
		assert.NotNil(t, resolvedAssessment.ResolvedAt)
		assert.True(t, resolvedAssessment.ResolvedAt.After(assessment.CreatedAt),
			"ResolvedAt should be after CreatedAt")
	})

	t.Run("Ignore Method", func(t *testing.T) {
		// Create a risk assessment
		assessment := &model.RiskAssessment{
			ID:        uuid.New().String(),
			UserID:    "user123",
			Type:      model.RiskTypePosition,
			Level:     model.RiskLevelMedium,
			Status:    model.RiskStatusActive,
			Symbol:    "BTCUSDT",
			Message:   "Active risk assessment",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Save to database
		err := repository.Create(ctx, assessment)
		require.NoError(t, err)

		// Retrieve to verify initial state
		initialAssessment, err := repository.GetByID(ctx, assessment.ID)
		require.NoError(t, err)
		assert.Equal(t, model.RiskStatusActive, initialAssessment.Status)

		// Record initial updated time
		initialUpdatedAt := initialAssessment.UpdatedAt

		// Call the Ignore method
		initialAssessment.Ignore()

		// Update in the database
		err = repository.Update(ctx, initialAssessment)
		require.NoError(t, err)

		// Retrieve again to verify changes
		ignoredAssessment, err := repository.GetByID(ctx, assessment.ID)
		require.NoError(t, err)

		// Verify the status changed and UpdatedAt was updated
		assert.Equal(t, model.RiskStatusIgnored, ignoredAssessment.Status)
		assert.True(t, ignoredAssessment.UpdatedAt.After(initialUpdatedAt),
			"UpdatedAt should be after the initial UpdatedAt")
		// ResolvedAt should still be nil
		assert.Nil(t, ignoredAssessment.ResolvedAt)
	})
}
