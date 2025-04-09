package unit

import (
	"testing"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
	"github.com/stretchr/testify/assert"
)

func TestNewCoin(t *testing.T) {
	now := time.Now()
	coin := models.NewCoin{
		ID:          1,
		Symbol:      "NEWUSDT",
		FoundAt:     now,
		BaseVolume:  1000.0,
		QuoteVolume: 50000.0,
		IsProcessed: false,
		IsDeleted:   false,
	}

	assert.Equal(t, int64(1), coin.ID)
	assert.Equal(t, "NEWUSDT", coin.Symbol)
	assert.Equal(t, now, coin.FoundAt)
	assert.Equal(t, 1000.0, coin.BaseVolume)
	assert.Equal(t, 50000.0, coin.QuoteVolume)
	assert.False(t, coin.IsProcessed)
	assert.False(t, coin.IsDeleted)
}
