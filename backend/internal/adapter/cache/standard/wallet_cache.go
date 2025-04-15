package standard

import (
	"sync"
	"time"

	gocache "github.com/patrickmn/go-cache"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
)

// WalletCache provides in-memory caching for wallet data with TTL
// Thread-safe, uses go-cache under the hood
// Keyed by userID+walletID or another unique key

type WalletCache struct {
	cache *gocache.Cache
	mu    sync.RWMutex
	defaultTTL time.Duration
}

func NewWalletCache(defaultTTL, cleanupInterval time.Duration) *WalletCache {
	return &WalletCache{
		cache: gocache.New(defaultTTL, cleanupInterval),
		defaultTTL: defaultTTL,
	}
}

// Set caches the wallet for the given key (userID+walletID)
func (wc *WalletCache) Set(key string, wallet *model.Wallet, ttl ...time.Duration) {
	wc.mu.Lock()
	defer wc.mu.Unlock()
	var duration time.Duration
	if len(ttl) > 0 {
		duration = ttl[0]
	} else {
		duration = wc.defaultTTL
	}
	wc.cache.Set(key, wallet, duration)
}

// Get retrieves a wallet from the cache
func (wc *WalletCache) Get(key string) (*model.Wallet, bool) {
	wc.mu.RLock()
	defer wc.mu.RUnlock()
	item, found := wc.cache.Get(key)
	if !found {
		return nil, false
	}
	wallet, ok := item.(*model.Wallet)
	return wallet, ok
}

// Delete removes a wallet from the cache
func (wc *WalletCache) Delete(key string) {
	wc.mu.Lock()
	defer wc.mu.Unlock()
	wc.cache.Delete(key)
}

// Clear clears all cached wallets
func (wc *WalletCache) Clear() {
	wc.mu.Lock()
	defer wc.mu.Unlock()
	wc.cache.Flush()
}

// Keys returns all keys in the cache
func (wc *WalletCache) Keys() []string {
	wc.mu.RLock()
	defer wc.mu.RUnlock()
	items := wc.cache.Items()
	keys := make([]string, 0, len(items))
	for k := range items {
		keys = append(keys, k)
	}
	return keys
}

// IsExpired checks if a wallet is expired
func (wc *WalletCache) IsExpired(key string) bool {
	wc.mu.RLock()
	defer wc.mu.RUnlock()
	_, found := wc.cache.Get(key)
	return !found
}
