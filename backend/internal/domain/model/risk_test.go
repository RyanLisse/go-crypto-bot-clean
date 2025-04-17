package model_test

import (
	"testing"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRiskAssessment_Resolve(t *testing.T) {
	// Create a new risk assessment
	assessment := &model.RiskAssessment{
		ID:        uuid.New().String(),
		UserID:    "user123",
		Type:      model.RiskTypePosition,
		Level:     model.RiskLevelMedium,
		Status:    model.RiskStatusActive,
		Symbol:    "BTCUSDT",
		Message:   "Test risk assessment",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Initial state check
	assert.Equal(t, model.RiskStatusActive, assessment.Status)
	assert.Nil(t, assessment.ResolvedAt)
	initialUpdatedAt := assessment.UpdatedAt

	// Call the Resolve method
	assessment.Resolve()

	// Check the status has been updated
	assert.Equal(t, model.RiskStatusResolved, assessment.Status)

	// Check that ResolvedAt is set
	assert.NotNil(t, assessment.ResolvedAt)

	// Check that UpdatedAt is updated
	assert.True(t, assessment.UpdatedAt.After(initialUpdatedAt))

	// Check that ResolvedAt is approximately now
	now := time.Now()
	resolvedTime := *assessment.ResolvedAt
	difference := now.Sub(resolvedTime)
	assert.True(t, difference < time.Second, "ResolvedAt should be very close to current time")
}

func TestRiskAssessment_Ignore(t *testing.T) {
	// Create a new risk assessment
	assessment := &model.RiskAssessment{
		ID:        uuid.New().String(),
		UserID:    "user123",
		Type:      model.RiskTypePosition,
		Level:     model.RiskLevelMedium,
		Status:    model.RiskStatusActive,
		Symbol:    "BTCUSDT",
		Message:   "Test risk assessment",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Initial state check
	assert.Equal(t, model.RiskStatusActive, assessment.Status)
	initialUpdatedAt := assessment.UpdatedAt

	// Call the Ignore method
	assessment.Ignore()

	// Check the status has been updated
	assert.Equal(t, model.RiskStatusIgnored, assessment.Status)

	// Check that ResolvedAt remains nil
	assert.Nil(t, assessment.ResolvedAt)

	// Check that UpdatedAt is updated
	assert.True(t, assessment.UpdatedAt.After(initialUpdatedAt))
}

func TestNewRiskAssessment(t *testing.T) {
	// Create a new risk assessment using the constructor
	userID := "user123"
	riskType := model.RiskTypePosition
	level := model.RiskLevelHigh
	message := "Test risk assessment"

	assessment := model.NewRiskAssessment(userID, riskType, level, message)

	// Check that all fields are set correctly
	assert.NotEmpty(t, assessment.ID)
	assert.Equal(t, userID, assessment.UserID)
	assert.Equal(t, riskType, assessment.Type)
	assert.Equal(t, level, assessment.Level)
	assert.Equal(t, model.RiskStatusActive, assessment.Status)
	assert.Equal(t, message, assessment.Message)
	assert.Equal(t, 0.0, assessment.Score)
	assert.WithinDuration(t, time.Now(), assessment.CreatedAt, time.Second)
	assert.WithinDuration(t, time.Now(), assessment.UpdatedAt, time.Second)
	assert.Nil(t, assessment.ResolvedAt)
}
