package crypto

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// KeyMetadata holds metadata for an encryption key
// Version is a monotonically increasing integer or string (e.g. ISO date)
type KeyMetadata struct {
	Version     string    // e.g. "2025-04-15" or "v1"
	Key         []byte    // actual key bytes
	CreatedAt   time.Time // when key was created/added
	Status      string    // "active", "retired", etc.
}

// KeyRegistry manages encryption keys and their versions
// Thread-safe for concurrent access
type KeyRegistry struct {
	mu        sync.RWMutex
	keys      map[string]*KeyMetadata // version -> metadata
	activeVer string                  // current active version
}

func NewKeyRegistry() *KeyRegistry {
	return &KeyRegistry{
		keys: make(map[string]*KeyMetadata),
	}
}

// AddKey adds a new key and sets it as active if requested
func (r *KeyRegistry) AddKey(version string, key []byte, setActive bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.keys[version]; exists {
		return fmt.Errorf("key version %s already exists", version)
	}
	r.keys[version] = &KeyMetadata{
		Version:   version,
		Key:       key,
		CreatedAt: time.Now(),
		Status:    "active",
	}
	if setActive {
		r.activeVer = version
	}
	return nil
}

// GetActiveKey returns the current active key and its version
func (r *KeyRegistry) GetActiveKey() (*KeyMetadata, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.activeVer == "" {
		return nil, errors.New("no active key")
	}
	key, ok := r.keys[r.activeVer]
	if !ok {
		return nil, errors.New("active key metadata not found")
	}
	return key, nil
}

// GetKey returns the key metadata for a given version
func (r *KeyRegistry) GetKey(version string) (*KeyMetadata, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	key, ok := r.keys[version]
	if !ok {
		return nil, fmt.Errorf("key version %s not found", version)
	}
	return key, nil
}

// RetireKey marks a key as retired (not used for new encryption)
func (r *KeyRegistry) RetireKey(version string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	key, ok := r.keys[version]
	if !ok {
		return fmt.Errorf("key version %s not found", version)
	}
	key.Status = "retired"
	if r.activeVer == version {
		r.activeVer = ""
	}
	return nil
}

// ListKeys returns metadata for all keys
func (r *KeyRegistry) ListKeys() []*KeyMetadata {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []*KeyMetadata
	for _, k := range r.keys {
		out = append(out, k)
	}
	return out
}
