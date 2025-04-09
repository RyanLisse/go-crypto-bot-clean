package unit

import (
	"testing"
	"time"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
	"github.com/stretchr/testify/assert"
)

func TestBoughtCoin(t *testing.T) {
	now := time.Now()
	coin := models.BoughtCoin{
		ID:            1,
		Symbol:        "BTCUSDT",
		PurchasePrice: 50000.0,
		Quantity:      0.1,
		BoughtAt:      now,
		StopLoss:      47500.0,
		TakeProfit:    55000.0,
		CurrentPrice:  51000.0,
		IsDeleted:     false,
		UpdatedAt:     now,
	}

	assert.Equal(t, int64(1), coin.ID)
	assert.Equal(t, "BTCUSDT", coin.Symbol)
	assert.Equal(t, 50000.0, coin.PurchasePrice)
	assert.Equal(t, 0.1, coin.Quantity)
	assert.Equal(t, now, coin.BoughtAt)
	assert.Equal(t, 47500.0, coin.StopLoss)
	assert.Equal(t, 55000.0, coin.TakeProfit)
	assert.Equal(t, 51000.0, coin.CurrentPrice)
	assert.False(t, coin.IsDeleted)
	assert.Equal(t, now, coin.UpdatedAt)
}
