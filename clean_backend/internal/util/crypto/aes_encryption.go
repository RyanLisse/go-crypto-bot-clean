package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// AESEncryptionService implements port.EncryptionService using AES-GCM
type AESEncryptionService struct {
	key    []byte
	logger *zerolog.Logger
}

// NewAESEncryptionService creates a new AESEncryptionService
func NewAESEncryptionService(key []byte, logger *zerolog.Logger) *AESEncryptionService {
	return &AESEncryptionService{
		key:    key,
		logger: logger,
	}
}

// Encrypt encrypts a plaintext string using AES-GCM
func (s *AESEncryptionService) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	// Create cipher
	block, err := aes.NewCipher(s.key)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to create AES cipher")
		return "", fmt.Errorf("failed to create AES cipher: %w", err)
	}

	// Create GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to create GCM")
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Create nonce
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		s.logger.Error().Err(err).Msg("Failed to generate nonce")
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt
	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts a ciphertext string using AES-GCM
func (s *AESEncryptionService) Decrypt(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	// Decode base64
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to decode base64 ciphertext")
		return "", fmt.Errorf("failed to decode base64 ciphertext: %w", err)
	}

	// Create cipher
	block, err := aes.NewCipher(s.key)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to create AES cipher")
		return "", fmt.Errorf("failed to create AES cipher: %w", err)
	}

	// Create GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to create GCM")
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Check if ciphertext is valid
	if len(data) < aesGCM.NonceSize() {
		err := errors.New("ciphertext too short")
		s.logger.Error().Err(err).Msg("Invalid ciphertext")
		return "", err
	}

	// Extract nonce and ciphertext
	nonce, ciphertextBytes := data[:aesGCM.NonceSize()], data[aesGCM.NonceSize():]

	// Decrypt
	plaintext, err := aesGCM.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to decrypt ciphertext")
		return "", fmt.Errorf("failed to decrypt ciphertext: %w", err)
	}

	return string(plaintext), nil
}

// Ensure AESEncryptionService implements port.EncryptionService
var _ port.EncryptionService = (*AESEncryptionService)(nil)
