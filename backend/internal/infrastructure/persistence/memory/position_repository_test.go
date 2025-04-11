package memory

import (
	"context"
	"testing"

	"go-crypto-bot-clean/backend/internal/domain/models"

	"github.com/stretchr/testify/assert"
)

func TestPositionRepository_Create(t *testing.T) {
	repo := NewPositionRepository()
	ctx := context.Background()

	position := &models.Position{
		ID:     "position1",
		Symbol: "BTC-USD",
		Status: models.PositionStatusOpen,
	}

	err := repo.Create(ctx, position)
	assert.NoError(t, err)

	// Test creating position with empty ID
	emptyPosition := &models.Position{
		Symbol: "BTC-USD",
	}
	err = repo.Create(ctx, emptyPosition)
	assert.Error(t, err)
}

func TestPositionRepository_GetByID(t *testing.T) {
	repo := NewPositionRepository()
	ctx := context.Background()

	position := &models.Position{
		ID:     "position1",
		Symbol: "BTC-USD",
		Status: models.PositionStatusOpen,
	}

	// Create the position first
	err := repo.Create(ctx, position)
	assert.NoError(t, err)

	// Test retrieving the position
	retrieved, err := repo.GetByID(ctx, "position1")
	assert.NoError(t, err)
	assert.Equal(t, position.ID, retrieved.ID)
	assert.Equal(t, position.Symbol, retrieved.Symbol)
	assert.Equal(t, position.Status, retrieved.Status)

	// Test retrieving non-existent position
	_, err = repo.GetByID(ctx, "nonexistent")
	assert.Error(t, err)
}

func TestPositionRepository_List(t *testing.T) {
	repo := NewPositionRepository()
	ctx := context.Background()

	// Create test positions
	positions := []*models.Position{
		{ID: "position1", Symbol: "BTC-USD", Status: models.PositionStatusOpen},
		{ID: "position2", Symbol: "BTC-USD", Status: models.PositionStatusClosed},
		{ID: "position3", Symbol: "ETH-USD", Status: models.PositionStatusOpen},
	}

	for _, position := range positions {
		err := repo.Create(ctx, position)
		assert.NoError(t, err)
	}

	// Test listing all positions
	allPositions, err := repo.List(ctx, "")
	assert.NoError(t, err)
	assert.Len(t, allPositions, 3)

	// Test filtering by status
	openPositions, err := repo.List(ctx, string(models.PositionStatusOpen))
	assert.NoError(t, err)
	assert.Len(t, openPositions, 2)
}

func TestPositionRepository_Update(t *testing.T) {
	repo := NewPositionRepository()
	ctx := context.Background()

	// Create a position
	position := &models.Position{
		ID:     "position1",
		Symbol: "BTC-USD",
		Status: models.PositionStatusOpen,
	}
	err := repo.Create(ctx, position)
	assert.NoError(t, err)

	// Update the position
	position.Status = models.PositionStatusClosed
	err = repo.Update(ctx, position)
	assert.NoError(t, err)

	// Verify the update
	updated, err := repo.GetByID(ctx, "position1")
	assert.NoError(t, err)
	assert.Equal(t, models.PositionStatusClosed, updated.Status)

	// Test updating non-existent position
	nonExistentPosition := &models.Position{ID: "nonexistent"}
	err = repo.Update(ctx, nonExistentPosition)
	assert.Error(t, err)
}

func TestPositionRepository_Delete(t *testing.T) {
	repo := NewPositionRepository()
	ctx := context.Background()

	// Create a position
	position := &models.Position{
		ID:     "position1",
		Symbol: "BTC-USD",
		Status: models.PositionStatusOpen,
	}
	err := repo.Create(ctx, position)
	assert.NoError(t, err)

	// Delete the position
	err = repo.Delete(ctx, "position1")
	assert.NoError(t, err)

	// Verify the position is deleted
	_, err = repo.GetByID(ctx, "position1")
	assert.Error(t, err)

	// Test deleting non-existent position
	err = repo.Delete(ctx, "nonexistent")
	assert.Error(t, err)
}

func TestPositionRepository_GetOpenPositionBySymbol(t *testing.T) {
	repo := NewPositionRepository()
	ctx := context.Background()

	// Create test positions
	positions := []*models.Position{
		{ID: "position1", Symbol: "BTC-USD", Status: models.PositionStatusOpen},
		{ID: "position2", Symbol: "BTC-USD", Status: models.PositionStatusClosed},
		{ID: "position3", Symbol: "ETH-USD", Status: models.PositionStatusOpen},
	}

	for _, position := range positions {
		err := repo.Create(ctx, position)
		assert.NoError(t, err)
	}

	// Test getting open position for BTC-USD
	btcPosition, err := repo.GetOpenPositionBySymbol(ctx, "BTC-USD")
	assert.NoError(t, err)
	assert.Equal(t, "position1", btcPosition.ID)

	// Test getting open position for non-existent symbol
	_, err = repo.GetOpenPositionBySymbol(ctx, "LTC-USD")
	assert.Error(t, err)
}
