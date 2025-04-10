package memory

import (
	"context"
	"testing"

	"go-crypto-bot-clean/backend/internal/domain/models"

	"github.com/stretchr/testify/assert"
)

func TestOrderRepository_Create(t *testing.T) {
	repo := NewOrderRepository()
	ctx := context.Background()

	order := &models.Order{
		ID:     "order1",
		Symbol: "BTC-USD",
		Status: models.OrderStatusNew,
	}

	err := repo.Create(ctx, order)
	assert.NoError(t, err)

	// Test creating order with empty ID
	emptyOrder := &models.Order{
		Symbol: "BTC-USD",
	}
	err = repo.Create(ctx, emptyOrder)
	assert.Error(t, err)
}

func TestOrderRepository_GetByID(t *testing.T) {
	repo := NewOrderRepository()
	ctx := context.Background()

	order := &models.Order{
		ID:     "order1",
		Symbol: "BTC-USD",
		Status: models.OrderStatusNew,
	}

	// Create the order first
	err := repo.Create(ctx, order)
	assert.NoError(t, err)

	// Test retrieving the order
	retrieved, err := repo.GetByID(ctx, "order1")
	assert.NoError(t, err)
	assert.Equal(t, order.ID, retrieved.ID)
	assert.Equal(t, order.Symbol, retrieved.Symbol)
	assert.Equal(t, order.Status, retrieved.Status)

	// Test retrieving non-existent order
	_, err = repo.GetByID(ctx, "nonexistent")
	assert.Error(t, err)
}

func TestOrderRepository_List(t *testing.T) {
	repo := NewOrderRepository()
	ctx := context.Background()

	// Create test orders
	orders := []*models.Order{
		{ID: "order1", Symbol: "BTC-USD", Status: models.OrderStatusNew},
		{ID: "order2", Symbol: "BTC-USD", Status: models.OrderStatusFilled},
		{ID: "order3", Symbol: "ETH-USD", Status: models.OrderStatusNew},
	}

	for _, order := range orders {
		err := repo.Create(ctx, order)
		assert.NoError(t, err)
	}

	// Test listing all orders
	allOrders, err := repo.List(ctx, "", "")
	assert.NoError(t, err)
	assert.Len(t, allOrders, 3)

	// Test filtering by symbol
	btcOrders, err := repo.List(ctx, "BTC-USD", "")
	assert.NoError(t, err)
	assert.Len(t, btcOrders, 2)

	// Test filtering by status
	newOrders, err := repo.List(ctx, "", models.OrderStatusNew)
	assert.NoError(t, err)
	assert.Len(t, newOrders, 2)

	// Test filtering by both symbol and status
	btcNewOrders, err := repo.List(ctx, "BTC-USD", models.OrderStatusNew)
	assert.NoError(t, err)
	assert.Len(t, btcNewOrders, 1)
}

func TestOrderRepository_Update(t *testing.T) {
	repo := NewOrderRepository()
	ctx := context.Background()

	// Create an order
	order := &models.Order{
		ID:     "order1",
		Symbol: "BTC-USD",
		Status: models.OrderStatusNew,
	}
	err := repo.Create(ctx, order)
	assert.NoError(t, err)

	// Update the order
	order.Status = models.OrderStatusFilled
	err = repo.Update(ctx, order)
	assert.NoError(t, err)

	// Verify the update
	updated, err := repo.GetByID(ctx, "order1")
	assert.NoError(t, err)
	assert.Equal(t, models.OrderStatusFilled, updated.Status)

	// Test updating non-existent order
	nonExistentOrder := &models.Order{ID: "nonexistent"}
	err = repo.Update(ctx, nonExistentOrder)
	assert.Error(t, err)
}

func TestOrderRepository_Delete(t *testing.T) {
	repo := NewOrderRepository()
	ctx := context.Background()

	// Create an order
	order := &models.Order{
		ID:     "order1",
		Symbol: "BTC-USD",
		Status: models.OrderStatusNew,
	}
	err := repo.Create(ctx, order)
	assert.NoError(t, err)

	// Delete the order
	err = repo.Delete(ctx, "order1")
	assert.NoError(t, err)

	// Verify the order is deleted
	_, err = repo.GetByID(ctx, "order1")
	assert.Error(t, err)

	// Test deleting non-existent order
	err = repo.Delete(ctx, "nonexistent")
	assert.Error(t, err)
}
