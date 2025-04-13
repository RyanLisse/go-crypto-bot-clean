package memory

import (
	"sync"
	"time"

	"github.com/neo/crypto-bot/internal/domain/port"
)

// GenericCache provides a generic cache implementation with TTL support
type GenericCache[T any] struct {
	value      *T
	expiration time.Time
	ttl        time.Duration
	mutex      sync.RWMutex
}

// NewGenericCache creates a new cache instance with the specified TTL
func NewGenericCache[T any](ttl time.Duration) port.Cache[T] {
	if ttl <= 0 {
		ttl = 5 * time.Minute // Default TTL
	}

	return &GenericCache[T]{
		value:      nil,
		expiration: time.Time{},
		ttl:        ttl,
		mutex:      sync.RWMutex{},
	}
}

// Get retrieves the cached value if it exists and is not expired
// Returns the value and a boolean indicating if the value was found
func (c *GenericCache[T]) Get() (*T, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.value == nil || time.Now().After(c.expiration) {
		return nil, false
	}

	return c.value, true
}

// Set stores a value in the cache with the configured TTL
func (c *GenericCache[T]) Set(value *T) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.value = value
	c.expiration = time.Now().Add(c.ttl)
}

// GetOrSet retrieves the cached value if valid, or sets it using the provided function
func (c *GenericCache[T]) GetOrSet(fetchFn func() (*T, error)) (*T, error) {
	// Check cache first
	if value, found := c.Get(); found {
		return value, nil
	}

	// Fetch new value
	value, err := fetchFn()
	if err != nil {
		return nil, err
	}

	// Store in cache
	c.Set(value)
	return value, nil
}

// Invalidate clears the cached value
func (c *GenericCache[T]) Invalidate() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.value = nil
	c.expiration = time.Time{}
}

// UpdateTTL changes the TTL for the cache
func (c *GenericCache[T]) UpdateTTL(ttl time.Duration) {
	if ttl <= 0 {
		return
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.ttl = ttl

	// Update expiration for existing value if present
	if c.value != nil && !c.expiration.IsZero() {
		c.expiration = time.Now().Add(ttl)
	}
}
