// Package cache provides caching mechanisms for MEXC API responses
package cache

import (
	"fmt"
	"sync"
	"time"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
)

// Cache is a generic cache interface
type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, ttl time.Duration)
	Delete(key string)
	Clear()
}

// TickerCache is a cache for ticker data
type TickerCache struct {
	data  map[string]*cacheItem
	mutex sync.RWMutex
}

// cacheItem represents a cached item with expiration
type cacheItem struct {
	value      interface{}
	expiration time.Time
}

// NewTickerCache creates a new ticker cache
func NewTickerCache() *TickerCache {
	return &TickerCache{
		data: make(map[string]*cacheItem),
	}
}

// Get retrieves a value from the cache
func (c *TickerCache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, exists := c.data[key]
	if !exists {
		return nil, false
	}

	// Check if the item has expired
	if time.Now().After(item.expiration) {
		return nil, false
	}

	return item.value, true
}

// Set adds a value to the cache with a TTL
func (c *TickerCache) Set(key string, value interface{}, ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data[key] = &cacheItem{
		value:      value,
		expiration: time.Now().Add(ttl),
	}
}

// Delete removes a value from the cache
func (c *TickerCache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.data, key)
}

// Clear removes all values from the cache
func (c *TickerCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data = make(map[string]*cacheItem)
}

// GetTicker retrieves a ticker from the cache
func (c *TickerCache) GetTicker(symbol string) (*models.Ticker, bool) {
	value, exists := c.Get(symbol)
	if !exists {
		return nil, false
	}

	ticker, ok := value.(*models.Ticker)
	return ticker, ok
}

// SetTicker adds a ticker to the cache
func (c *TickerCache) SetTicker(symbol string, ticker *models.Ticker, ttl time.Duration) {
	c.Set(symbol, ticker, ttl)
}

// GetAllTickers retrieves all tickers from the cache
func (c *TickerCache) GetAllTickers() (map[string]*models.Ticker, bool) {
	value, exists := c.Get("all_tickers")
	if !exists {
		return nil, false
	}

	tickers, ok := value.(map[string]*models.Ticker)
	return tickers, ok
}

// SetAllTickers adds all tickers to the cache
func (c *TickerCache) SetAllTickers(tickers map[string]*models.Ticker, ttl time.Duration) {
	c.Set("all_tickers", tickers, ttl)
}

// OrderBookCache is a cache for order book data
type OrderBookCache struct {
	data  map[string]*cacheItem
	mutex sync.RWMutex
}

// NewOrderBookCache creates a new order book cache
func NewOrderBookCache() *OrderBookCache {
	return &OrderBookCache{
		data: make(map[string]*cacheItem),
	}
}

// Get retrieves a value from the cache
func (c *OrderBookCache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, exists := c.data[key]
	if !exists {
		return nil, false
	}

	// Check if the item has expired
	if time.Now().After(item.expiration) {
		return nil, false
	}

	return item.value, true
}

// Set adds a value to the cache with a TTL
func (c *OrderBookCache) Set(key string, value interface{}, ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data[key] = &cacheItem{
		value:      value,
		expiration: time.Now().Add(ttl),
	}
}

// Delete removes a value from the cache
func (c *OrderBookCache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.data, key)
}

// Clear removes all values from the cache
func (c *OrderBookCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data = make(map[string]*cacheItem)
}

// GetOrderBook retrieves an order book from the cache
func (c *OrderBookCache) GetOrderBook(symbol string) (*models.OrderBookUpdate, bool) {
	value, exists := c.Get(symbol)
	if !exists {
		return nil, false
	}

	orderBook, ok := value.(*models.OrderBookUpdate)
	return orderBook, ok
}

// SetOrderBook adds an order book to the cache
func (c *OrderBookCache) SetOrderBook(symbol string, orderBook *models.OrderBookUpdate, ttl time.Duration) {
	c.Set(symbol, orderBook, ttl)
}

// KlineCache is a cache for kline data
type KlineCache struct {
	data  map[string]*cacheItem
	mutex sync.RWMutex
}

// NewKlineCache creates a new kline cache
func NewKlineCache() *KlineCache {
	return &KlineCache{
		data: make(map[string]*cacheItem),
	}
}

// Get retrieves a value from the cache
func (c *KlineCache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, exists := c.data[key]
	if !exists {
		return nil, false
	}

	// Check if the item has expired
	if time.Now().After(item.expiration) {
		return nil, false
	}

	return item.value, true
}

// Set adds a value to the cache with a TTL
func (c *KlineCache) Set(key string, value interface{}, ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data[key] = &cacheItem{
		value:      value,
		expiration: time.Now().Add(ttl),
	}
}

// Delete removes a value from the cache
func (c *KlineCache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.data, key)
}

// Clear removes all values from the cache
func (c *KlineCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data = make(map[string]*cacheItem)
}

// GetKlines retrieves klines from the cache
func (c *KlineCache) GetKlines(symbol, interval string, limit int) ([]*models.Kline, bool) {
	key := formatKlineKey(symbol, interval, limit)
	value, exists := c.Get(key)
	if !exists {
		return nil, false
	}

	klines, ok := value.([]*models.Kline)
	return klines, ok
}

// SetKlines adds klines to the cache
func (c *KlineCache) SetKlines(symbol, interval string, limit int, klines []*models.Kline, ttl time.Duration) {
	key := formatKlineKey(symbol, interval, limit)
	c.Set(key, klines, ttl)
}

// formatKlineKey formats a key for kline cache
func formatKlineKey(symbol, interval string, limit int) string {
	return fmt.Sprintf("%s_%s_%d", symbol, interval, limit)
}

// NewCoinCache is a cache for new coin data
type NewCoinCache struct {
	data  map[string]*cacheItem
	mutex sync.RWMutex
}

// NewNewCoinCache creates a new coin cache
func NewNewCoinCache() *NewCoinCache {
	return &NewCoinCache{
		data: make(map[string]*cacheItem),
	}
}

// Get retrieves a value from the cache
func (c *NewCoinCache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, exists := c.data[key]
	if !exists {
		return nil, false
	}

	// Check if the item has expired
	if time.Now().After(item.expiration) {
		return nil, false
	}

	return item.value, true
}

// Set adds a value to the cache with a TTL
func (c *NewCoinCache) Set(key string, value interface{}, ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data[key] = &cacheItem{
		value:      value,
		expiration: time.Now().Add(ttl),
	}
}

// Delete removes a value from the cache
func (c *NewCoinCache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.data, key)
}

// Clear removes all values from the cache
func (c *NewCoinCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data = make(map[string]*cacheItem)
}

// GetNewCoins retrieves new coins from the cache
func (c *NewCoinCache) GetNewCoins() ([]*models.NewCoin, bool) {
	value, exists := c.Get("new_coins")
	if !exists {
		return nil, false
	}

	newCoins, ok := value.([]*models.NewCoin)
	return newCoins, ok
}

// SetNewCoins adds new coins to the cache
func (c *NewCoinCache) SetNewCoins(newCoins []*models.NewCoin, ttl time.Duration) {
	c.Set("new_coins", newCoins, ttl)
}
