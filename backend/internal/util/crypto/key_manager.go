package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

// KeyManager manages encryption keys, including rotation and secure storage
type KeyManager interface {
	// GetCurrentKey returns the current encryption key
	GetCurrentKey() ([]byte, error)

	// GetKeyByID returns a specific encryption key by ID
	GetKeyByID(keyID string) ([]byte, error)

	// RotateKey generates a new encryption key and makes it the current key
	RotateKey() (string, error)

	// AddKey adds a new encryption key with the given ID
	AddKey(keyID string, key []byte) error
}

// EncryptionKey represents an encryption key with metadata
type EncryptionKey struct {
	ID        string
	Key       []byte
	CreatedAt time.Time
}

// EnvKeyManager implements KeyManager using environment variables
type EnvKeyManager struct {
	keys       map[string]EncryptionKey
	currentKey string
	mu         sync.RWMutex
}

// NewEnvKeyManager creates a new EnvKeyManager
func NewEnvKeyManager() (*EnvKeyManager, error) {
	manager := &EnvKeyManager{
		keys: make(map[string]EncryptionKey),
	}

	// Load keys from environment variables
	if err := manager.loadKeysFromEnv(); err != nil {
		return nil, err
	}

	return manager, nil
}

// loadKeysFromEnv loads encryption keys from environment variables
func (m *EnvKeyManager) loadKeysFromEnv() error {
	// Get current key ID
	currentKeyID := os.Getenv("ENCRYPTION_CURRENT_KEY_ID")
	if currentKeyID == "" {
		return errors.New("ENCRYPTION_CURRENT_KEY_ID environment variable not set")
	}

	// Get keys
	keysEnv := os.Getenv("ENCRYPTION_KEYS")
	if keysEnv == "" {
		return errors.New("ENCRYPTION_KEYS environment variable not set")
	}

	// Parse keys
	keyPairs := strings.Split(keysEnv, ",")
	for _, pair := range keyPairs {
		parts := strings.Split(pair, ":")
		if len(parts) != 2 {
			return errors.New("invalid key format in ENCRYPTION_KEYS")
		}

		keyID := parts[0]
		keyB64 := parts[1]

		key, err := base64.StdEncoding.DecodeString(keyB64)
		if err != nil {
			return err
		}

		if len(key) != 32 {
			return errors.New("encryption key must be 32 bytes (256 bits)")
		}

		m.keys[keyID] = EncryptionKey{
			ID:        keyID,
			Key:       key,
			CreatedAt: time.Now(), // We don't have the actual creation time
		}
	}

	// Verify current key exists
	if _, ok := m.keys[currentKeyID]; !ok {
		return errors.New("current key ID not found in ENCRYPTION_KEYS")
	}

	m.currentKey = currentKeyID
	return nil
}

// GetCurrentKey returns the current encryption key
func (m *EnvKeyManager) GetCurrentKey() ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	key, ok := m.keys[m.currentKey]
	if !ok {
		return nil, errors.New("current key not found")
	}

	return key.Key, nil
}

// GetKeyByID returns a specific encryption key by ID
func (m *EnvKeyManager) GetKeyByID(keyID string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	key, ok := m.keys[keyID]
	if !ok {
		return nil, errors.New("key not found")
	}

	return key.Key, nil
}

// RotateKey generates a new encryption key and makes it the current key
func (m *EnvKeyManager) RotateKey() (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Generate new key
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return "", err
	}

	// Generate new key ID
	keyID := generateKeyID()

	// Add new key
	m.keys[keyID] = EncryptionKey{
		ID:        keyID,
		Key:       key,
		CreatedAt: time.Now(),
	}

	// Update current key
	m.currentKey = keyID

	// Return the new key ID and base64-encoded key
	return keyID + ":" + base64.StdEncoding.EncodeToString(key), nil
}

// AddKey adds a new encryption key with the given ID
func (m *EnvKeyManager) AddKey(keyID string, key []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(key) != 32 {
		return errors.New("encryption key must be 32 bytes (256 bits)")
	}

	m.keys[keyID] = EncryptionKey{
		ID:        keyID,
		Key:       key,
		CreatedAt: time.Now(),
	}

	return nil
}

// generateKeyID generates a new random key ID
func generateKeyID() string {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		// Fallback to timestamp if random fails
		return "key-" + time.Now().Format("20060102150405")
	}
	return "key-" + base64.URLEncoding.EncodeToString(b)[:10]
}
