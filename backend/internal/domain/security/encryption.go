package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
)

// EncryptionService provides methods for encrypting and decrypting data
type EncryptionService interface {
	// Encrypt encrypts data
	Encrypt(data []byte) (string, error)

	// Decrypt decrypts data
	Decrypt(encryptedData string) ([]byte, error)

	// EncryptString encrypts a string
	EncryptString(data string) (string, error)

	// DecryptString decrypts a string
	DecryptString(encryptedData string) (string, error)
}

// encryptionServiceImpl implements EncryptionService
type encryptionServiceImpl struct {
	key []byte
}

// NewEncryptionService creates a new encryption service
func NewEncryptionService() (EncryptionService, error) {
	// Get encryption key from environment variable
	keyStr := os.Getenv("ENCRYPTION_KEY")
	if keyStr == "" {
		// Generate a random key if not provided
		key := make([]byte, 32) // AES-256
		if _, err := io.ReadFull(rand.Reader, key); err != nil {
			return nil, fmt.Errorf("failed to generate encryption key: %w", err)
		}
		return &encryptionServiceImpl{key: key}, nil
	}

	// Decode base64 key
	key, err := base64.StdEncoding.DecodeString(keyStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode encryption key: %w", err)
	}

	// Ensure key is the correct length
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, errors.New("encryption key must be 16, 24, or 32 bytes (AES-128, AES-192, or AES-256)")
	}

	return &encryptionServiceImpl{key: key}, nil
}

// Encrypt encrypts data using AES-GCM
func (s *encryptionServiceImpl) Encrypt(data []byte) (string, error) {
	// Create cipher block
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher block: %w", err)
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Create nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to create nonce: %w", err)
	}

	// Encrypt data
	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	// Encode as base64
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts data using AES-GCM
func (s *encryptionServiceImpl) Decrypt(encryptedData string) ([]byte, error) {
	// Decode base64
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	// Create cipher block
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher block: %w", err)
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Check ciphertext length
	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}

	// Extract nonce and ciphertext
	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]

	// Decrypt data
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

// EncryptString encrypts a string
func (s *encryptionServiceImpl) EncryptString(data string) (string, error) {
	return s.Encrypt([]byte(data))
}

// DecryptString decrypts a string
func (s *encryptionServiceImpl) DecryptString(encryptedData string) (string, error) {
	plaintext, err := s.Decrypt(encryptedData)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
