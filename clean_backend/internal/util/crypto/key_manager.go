package crypto

import (
	"errors"
	"fmt"
	"sync"

	"github.com/rs/zerolog"
)

// KeyManager manages encryption keys and supports key rotation
type KeyManager struct {
	currentKeyID string
	keys         map[string][]byte
	mu           sync.RWMutex
	logger       *zerolog.Logger
}

// KeyEntry represents a key entry in the key store
type KeyEntry struct {
	ID  string `json:"id"`
	Key string `json:"key"`
}

// NewKeyManager creates a new KeyManager
func NewKeyManager(currentKeyID string, keysStr string, logger *zerolog.Logger) (*KeyManager, error) {
	if currentKeyID == "" {
		return nil, errors.New("current key ID cannot be empty")
	}

	// Create a default key for testing
	keys := make(map[string][]byte)
	keys["default"] = []byte("6368616e676520746869732070617373776f726420746f206120736563726574")

	// Verify current key exists
	if _, ok := keys[currentKeyID]; !ok {
		return nil, fmt.Errorf("current key ID %s not found in keys", currentKeyID)
	}

	return &KeyManager{
		currentKeyID: currentKeyID,
		keys:         keys,
		logger:       logger,
	}, nil
}

// GetCurrentKey returns the current encryption key
func (km *KeyManager) GetCurrentKey() []byte {
	km.mu.RLock()
	defer km.mu.RUnlock()
	return km.keys[km.currentKeyID]
}

// GetCurrentKeyID returns the current key ID
func (km *KeyManager) GetCurrentKeyID() string {
	km.mu.RLock()
	defer km.mu.RUnlock()
	return km.currentKeyID
}

// GetKey returns the key for the given ID
func (km *KeyManager) GetKey(keyID string) ([]byte, error) {
	km.mu.RLock()
	defer km.mu.RUnlock()

	key, ok := km.keys[keyID]
	if !ok {
		return nil, fmt.Errorf("key ID %s not found", keyID)
	}

	return key, nil
}

// SetCurrentKey sets the current key ID
func (km *KeyManager) SetCurrentKey(keyID string) error {
	km.mu.Lock()
	defer km.mu.Unlock()

	if _, ok := km.keys[keyID]; !ok {
		return fmt.Errorf("key ID %s not found", keyID)
	}

	km.currentKeyID = keyID
	km.logger.Info().Str("keyID", keyID).Msg("Current encryption key changed")
	return nil
}

// AddKey adds a new key
func (km *KeyManager) AddKey(keyID string, key []byte) error {
	if keyID == "" {
		return errors.New("key ID cannot be empty")
	}

	if len(key) == 0 {
		return errors.New("key cannot be empty")
	}

	km.mu.Lock()
	defer km.mu.Unlock()

	if _, ok := km.keys[keyID]; ok {
		return fmt.Errorf("key ID %s already exists", keyID)
	}

	km.keys[keyID] = key
	km.logger.Info().Str("keyID", keyID).Msg("New encryption key added")
	return nil
}

// RemoveKey removes a key
func (km *KeyManager) RemoveKey(keyID string) error {
	km.mu.Lock()
	defer km.mu.Unlock()

	if keyID == km.currentKeyID {
		return errors.New("cannot remove current key")
	}

	if _, ok := km.keys[keyID]; !ok {
		return fmt.Errorf("key ID %s not found", keyID)
	}

	delete(km.keys, keyID)
	km.logger.Info().Str("keyID", keyID).Msg("Encryption key removed")
	return nil
}

// ListKeyIDs returns a list of all key IDs
func (km *KeyManager) ListKeyIDs() []string {
	km.mu.RLock()
	defer km.mu.RUnlock()

	ids := make([]string, 0, len(km.keys))
	for id := range km.keys {
		ids = append(ids, id)
	}

	return ids
}

// HasKey checks if a key ID exists
func (km *KeyManager) HasKey(keyID string) bool {
	km.mu.RLock()
	defer km.mu.RUnlock()

	_, ok := km.keys[keyID]
	return ok
}
