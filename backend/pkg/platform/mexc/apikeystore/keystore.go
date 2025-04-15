package apikeystore

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

var (
	// ErrKeyNotFound is returned when an API key is not found
	ErrKeyNotFound = errors.New("API key not found")

	// ErrEncryptionKeyRequired is returned when an encryption key is required but not provided
	ErrEncryptionKeyRequired = errors.New("encryption key required")

	// ErrInvalidEncryptionKey is returned when an invalid encryption key is used
	ErrInvalidEncryptionKey = errors.New("invalid encryption key")
)

// APIKeyCredentials represents API key credentials
type APIKeyCredentials struct {
	APIKey    string
	SecretKey string
}

// KeyStore defines the interface for API key storage
type KeyStore interface {
	// GetAPIKey retrieves API key credentials for a given key ID
	GetAPIKey(keyID string) (*APIKeyCredentials, error)

	// SetAPIKey stores API key credentials for a given key ID
	SetAPIKey(keyID string, creds *APIKeyCredentials) error

	// DeleteAPIKey removes API key credentials for a given key ID
	DeleteAPIKey(keyID string) error
}

// MemoryKeyStore implements an in-memory key store
type MemoryKeyStore struct {
	keys  map[string]*APIKeyCredentials
	mutex sync.RWMutex
}

// NewMemoryKeyStore creates a new in-memory key store
func NewMemoryKeyStore() *MemoryKeyStore {
	return &MemoryKeyStore{
		keys: make(map[string]*APIKeyCredentials),
	}
}

// GetAPIKey retrieves API key credentials for a given key ID
func (s *MemoryKeyStore) GetAPIKey(keyID string) (*APIKeyCredentials, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	creds, exists := s.keys[keyID]
	if !exists {
		return nil, ErrKeyNotFound
	}

	// Return a copy to prevent modification of stored credentials
	return &APIKeyCredentials{
		APIKey:    creds.APIKey,
		SecretKey: creds.SecretKey,
	}, nil
}

// SetAPIKey stores API key credentials for a given key ID
func (s *MemoryKeyStore) SetAPIKey(keyID string, creds *APIKeyCredentials) error {
	if creds == nil {
		return errors.New("credentials cannot be nil")
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Store a copy to prevent modification of stored credentials
	s.keys[keyID] = &APIKeyCredentials{
		APIKey:    creds.APIKey,
		SecretKey: creds.SecretKey,
	}

	return nil
}

// DeleteAPIKey removes API key credentials for a given key ID
func (s *MemoryKeyStore) DeleteAPIKey(keyID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.keys, keyID)
	return nil
}

// FileKeyStore implements a file-based key store with encryption
type FileKeyStore struct {
	filePath      string
	encryptionKey []byte
	mutex         sync.RWMutex
	keys          map[string]*APIKeyCredentials
}

// NewFileKeyStore creates a new file-based key store
func NewFileKeyStore(filePath string, encryptionKey []byte) (*FileKeyStore, error) {
	if len(encryptionKey) != 32 {
		return nil, errors.New("encryption key must be 32 bytes")
	}

	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	store := &FileKeyStore{
		filePath:      filePath,
		encryptionKey: encryptionKey,
		keys:          make(map[string]*APIKeyCredentials),
	}

	// Load existing keys if file exists
	if _, err := os.Stat(filePath); err == nil {
		if err := store.load(); err != nil {
			return nil, fmt.Errorf("failed to load API keys from file: %w", err)
		}
	}

	return store, nil
}

// GetAPIKey retrieves API key credentials for a given key ID
func (s *FileKeyStore) GetAPIKey(keyID string) (*APIKeyCredentials, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	creds, exists := s.keys[keyID]
	if !exists {
		return nil, ErrKeyNotFound
	}

	// Return a copy to prevent modification of stored credentials
	return &APIKeyCredentials{
		APIKey:    creds.APIKey,
		SecretKey: creds.SecretKey,
	}, nil
}

// SetAPIKey stores API key credentials for a given key ID
func (s *FileKeyStore) SetAPIKey(keyID string, creds *APIKeyCredentials) error {
	if creds == nil {
		return errors.New("credentials cannot be nil")
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Store a copy to prevent modification of stored credentials
	s.keys[keyID] = &APIKeyCredentials{
		APIKey:    creds.APIKey,
		SecretKey: creds.SecretKey,
	}

	// Save to file
	return s.save()
}

// DeleteAPIKey removes API key credentials for a given key ID
func (s *FileKeyStore) DeleteAPIKey(keyID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.keys, keyID)
	return s.save()
}

// load loads API keys from file
func (s *FileKeyStore) load() error {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Decrypt data
	decryptedData, err := decrypt(data, s.encryptionKey)
	if err != nil {
		return fmt.Errorf("failed to decrypt data: %w", err)
	}

	// Unmarshal data
	var keys map[string]*APIKeyCredentials
	if err := json.Unmarshal(decryptedData, &keys); err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}

	s.keys = keys
	return nil
}

// save saves API keys to file
func (s *FileKeyStore) save() error {
	// Marshal data
	data, err := json.Marshal(s.keys)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	// Encrypt data
	encryptedData, err := encrypt(data, s.encryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt data: %w", err)
	}

	// Write to file (use temporary file and rename for atomicity)
	tempFile := s.filePath + ".tmp"
	if err := os.WriteFile(tempFile, encryptedData, 0600); err != nil {
		return fmt.Errorf("failed to write temporary file: %w", err)
	}

	if err := os.Rename(tempFile, s.filePath); err != nil {
		return fmt.Errorf("failed to rename temporary file: %w", err)
	}

	return nil
}

// encrypt encrypts data using AES-GCM
func encrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// decrypt decrypts data using AES-GCM
func decrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(data) < gcm.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// EnvironmentKeyStore implements a key store that loads keys from environment variables
type EnvironmentKeyStore struct {
	apiKeyPrefix    string
	secretKeyPrefix string
}

// NewEnvironmentKeyStore creates a new environment variable key store
func NewEnvironmentKeyStore(apiKeyPrefix, secretKeyPrefix string) *EnvironmentKeyStore {
	return &EnvironmentKeyStore{
		apiKeyPrefix:    apiKeyPrefix,
		secretKeyPrefix: secretKeyPrefix,
	}
}

// GetAPIKey retrieves API key credentials for a given key ID
func (s *EnvironmentKeyStore) GetAPIKey(keyID string) (*APIKeyCredentials, error) {
	apiKeyEnv := s.apiKeyPrefix + keyID
	secretKeyEnv := s.secretKeyPrefix + keyID

	apiKey := os.Getenv(apiKeyEnv)
	secretKey := os.Getenv(secretKeyEnv)

	if apiKey == "" || secretKey == "" {
		return nil, ErrKeyNotFound
	}

	return &APIKeyCredentials{
		APIKey:    apiKey,
		SecretKey: secretKey,
	}, nil
}

// SetAPIKey is not supported for environment variables
func (s *EnvironmentKeyStore) SetAPIKey(keyID string, creds *APIKeyCredentials) error {
	return errors.New("setting API keys in environment variables is not supported")
}

// DeleteAPIKey is not supported for environment variables
func (s *EnvironmentKeyStore) DeleteAPIKey(keyID string) error {
	return errors.New("deleting API keys from environment variables is not supported")
}

// CompositeKeyStore implements a key store that tries multiple key stores in sequence
type CompositeKeyStore struct {
	stores []KeyStore
}

// NewCompositeKeyStore creates a new composite key store
func NewCompositeKeyStore(stores ...KeyStore) *CompositeKeyStore {
	return &CompositeKeyStore{
		stores: stores,
	}
}

// GetAPIKey retrieves API key credentials for a given key ID
func (s *CompositeKeyStore) GetAPIKey(keyID string) (*APIKeyCredentials, error) {
	for _, store := range s.stores {
		creds, err := store.GetAPIKey(keyID)
		if err == nil {
			return creds, nil
		}
		if !errors.Is(err, ErrKeyNotFound) {
			return nil, err
		}
	}
	return nil, ErrKeyNotFound
}

// SetAPIKey stores API key credentials for a given key ID
func (s *CompositeKeyStore) SetAPIKey(keyID string, creds *APIKeyCredentials) error {
	if len(s.stores) == 0 {
		return errors.New("no key stores available")
	}
	// Set in the first store (primary)
	return s.stores[0].SetAPIKey(keyID, creds)
}

// DeleteAPIKey removes API key credentials for a given key ID
func (s *CompositeKeyStore) DeleteAPIKey(keyID string) error {
	var lastErr error
	for _, store := range s.stores {
		if err := store.DeleteAPIKey(keyID); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// Base64KeyStore implements a key store that encodes/decodes keys in base64
type Base64KeyStore struct {
	store KeyStore
}

// NewBase64KeyStore creates a new base64 key store
func NewBase64KeyStore(store KeyStore) *Base64KeyStore {
	return &Base64KeyStore{
		store: store,
	}
}

// GetAPIKey retrieves API key credentials for a given key ID
func (s *Base64KeyStore) GetAPIKey(keyID string) (*APIKeyCredentials, error) {
	creds, err := s.store.GetAPIKey(keyID)
	if err != nil {
		return nil, err
	}

	apiKey, err := base64.StdEncoding.DecodeString(creds.APIKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode API key: %w", err)
	}

	secretKey, err := base64.StdEncoding.DecodeString(creds.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode secret key: %w", err)
	}

	return &APIKeyCredentials{
		APIKey:    string(apiKey),
		SecretKey: string(secretKey),
	}, nil
}

// SetAPIKey stores API key credentials for a given key ID
func (s *Base64KeyStore) SetAPIKey(keyID string, creds *APIKeyCredentials) error {
	encodedCreds := &APIKeyCredentials{
		APIKey:    base64.StdEncoding.EncodeToString([]byte(creds.APIKey)),
		SecretKey: base64.StdEncoding.EncodeToString([]byte(creds.SecretKey)),
	}
	return s.store.SetAPIKey(keyID, encodedCreds)
}

// DeleteAPIKey removes API key credentials for a given key ID
func (s *Base64KeyStore) DeleteAPIKey(keyID string) error {
	return s.store.DeleteAPIKey(keyID)
}
