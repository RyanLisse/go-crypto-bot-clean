package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"
)

// EncryptionService handles encryption and decryption of sensitive data
type EncryptionService interface {
	Encrypt(plaintext string) ([]byte, error)
	Decrypt(ciphertext []byte) (string, error)
}

// AESEncryptionService implements EncryptionService using AES-256-GCM
type AESEncryptionService struct {
	key []byte
}

// NewAESEncryptionService creates a new AESEncryptionService
func NewAESEncryptionService() (*AESEncryptionService, error) {
	keyB64 := os.Getenv("MEXC_CRED_ENCRYPTION_KEY")
	if keyB64 == "" {
		return nil, errors.New("MEXC_CRED_ENCRYPTION_KEY environment variable not set")
	}

	key, err := base64.StdEncoding.DecodeString(keyB64)
	if err != nil {
		return nil, err
	}

	if len(key) != 32 {
		return nil, errors.New("encryption key must be 32 bytes (256 bits)")
	}

	return &AESEncryptionService{
		key: key,
	}, nil
}

// Encrypt encrypts a string using AES-256-GCM
func (s *AESEncryptionService) Encrypt(plaintext string) ([]byte, error) {
	// Create cipher
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return nil, err
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Create nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Encrypt
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	return ciphertext, nil
}

// Decrypt decrypts a string using AES-256-GCM
func (s *AESEncryptionService) Decrypt(ciphertext []byte) (string, error) {
	// Create cipher
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", err
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Check ciphertext length
	if len(ciphertext) < gcm.NonceSize() {
		return "", errors.New("ciphertext too short")
	}

	// Extract nonce and ciphertext
	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// GenerateEncryptionKey generates a new random encryption key
func GenerateEncryptionKey() (string, error) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}
