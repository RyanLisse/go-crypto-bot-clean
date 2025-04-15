package crypto

import (
	"encoding/base64"
	"errors"
	"os"
	"sync"
	"time"
)

// EncryptionServiceType represents the type of encryption service
type EncryptionServiceType string

const (
	// BasicEncryptionService is the basic encryption service
	BasicEncryptionService EncryptionServiceType = "basic"

	// EnhancedEncryptionService is the enhanced encryption service with key rotation
	EnhancedEncryptionServiceType EncryptionServiceType = "enhanced"
)

// EncryptionServiceFactory creates encryption services
type EncryptionServiceFactory struct {
	keyManager KeyManager
	services   map[EncryptionServiceType]EncryptionService
	mu         sync.RWMutex
}

// NewEncryptionServiceFactory creates a new EncryptionServiceFactory
func NewEncryptionServiceFactory() (*EncryptionServiceFactory, error) {
	// Create key manager
	keyManager, err := createKeyManager()
	if err != nil {
		return nil, err
	}

	factory := &EncryptionServiceFactory{
		keyManager: keyManager,
		services:   make(map[EncryptionServiceType]EncryptionService),
	}

	return factory, nil
}

// GetEncryptionService returns an encryption service of the specified type
func (f *EncryptionServiceFactory) GetEncryptionService(serviceType EncryptionServiceType) (EncryptionService, error) {
	f.mu.RLock()
	service, ok := f.services[serviceType]
	f.mu.RUnlock()

	if ok {
		return service, nil
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	// Check again in case another goroutine created the service
	service, ok = f.services[serviceType]
	if ok {
		return service, nil
	}

	// Create new service
	var err error
	switch serviceType {
	case BasicEncryptionService:
		service, err = f.createBasicEncryptionService()
	case EnhancedEncryptionServiceType:
		service, err = f.createEnhancedEncryptionService()
	default:
		return nil, errors.New("unknown encryption service type")
	}

	if err != nil {
		return nil, err
	}

	f.services[serviceType] = service
	return service, nil
}

// createBasicEncryptionService creates a basic encryption service
func (f *EncryptionServiceFactory) createBasicEncryptionService() (EncryptionService, error) {
	key, err := f.keyManager.GetCurrentKey()
	if err != nil {
		return nil, err
	}

	return &AESEncryptionService{
		key: key,
	}, nil
}

// createEnhancedEncryptionService creates an enhanced encryption service
func (f *EncryptionServiceFactory) createEnhancedEncryptionService() (EncryptionService, error) {
	return NewEnhancedEncryptionService(f.keyManager), nil
}

// createKeyManager creates a key manager based on environment variables
func createKeyManager() (KeyManager, error) {
	// Check if we should use the environment key manager
	if os.Getenv("ENCRYPTION_KEYS") != "" {
		return NewEnvKeyManager()
	}

	// Fallback to a simple key manager with a single key
	keyB64 := os.Getenv("ENCRYPTION_KEY")
	if keyB64 == "" {
		// For development, use a default key
		keyB64 = "Wn3PvhLOYk0QpFdod9qUDRRik9cI8jD3noi0TgrTJ1M="
	}

	key, err := base64.StdEncoding.DecodeString(keyB64)
	if err != nil || len(key) != 32 {
		return nil, errors.New("invalid encryption key")
	}

	manager := &EnvKeyManager{
		keys:       make(map[string]EncryptionKey),
		currentKey: "default",
	}

	manager.keys["default"] = EncryptionKey{
		ID:        "default",
		Key:       key,
		CreatedAt: time.Now(),
	}

	return manager, nil
}
