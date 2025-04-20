package standard

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
)

func TestWalletCache_SetGetDelete(t *testing.T) {
	cache := NewWalletCache(2*time.Second, 1*time.Minute)
	key := "user1_wallet1"
	wallet := &model.Wallet{ID: "wallet1", UserID: "user1"}

	// Set and Get
	cache.Set(key, wallet)
	got, found := cache.Get(key)
	assert.True(t, found)
	assert.Equal(t, wallet, got)

	// Delete
	cache.Delete(key)
	_, found = cache.Get(key)
	assert.False(t, found)
}

func TestWalletCache_TTLExpiry(t *testing.T) {
	cache := NewWalletCache(1*time.Second, 1*time.Minute)
	key := "user2_wallet2"
	wallet := &model.Wallet{ID: "wallet2", UserID: "user2"}

	cache.Set(key, wallet)
	_, found := cache.Get(key)
	assert.True(t, found)

	time.Sleep(1100 * time.Millisecond)
	_, found = cache.Get(key)
	assert.False(t, found)
}

func TestWalletCache_ClearAndKeys(t *testing.T) {
	cache := NewWalletCache(5*time.Second, 1*time.Minute)
	cache.Set("k1", &model.Wallet{ID: "w1"})
	cache.Set("k2", &model.Wallet{ID: "w2"})

	keys := cache.Keys()
	assert.ElementsMatch(t, []string{"k1", "k2"}, keys)

	cache.Clear()
	keys = cache.Keys()
	assert.Empty(t, keys)
}

func TestWalletCache_IsExpired(t *testing.T) {
	cache := NewWalletCache(1*time.Second, 1*time.Minute)
	key := "user3_wallet3"
	wallet := &model.Wallet{ID: "wallet3", UserID: "user3"}
	cache.Set(key, wallet)
	assert.False(t, cache.IsExpired(key))
	time.Sleep(1100 * time.Millisecond)
	assert.True(t, cache.IsExpired(key))
}
