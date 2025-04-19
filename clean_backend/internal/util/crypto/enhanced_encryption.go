package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// EnhancedEncryptionService implements EncryptionService with key rotation support
type EnhancedEncryptionService struct {
	keyManager *KeyManager
	logger     *zerolog.Logger
}

// EncryptedData represents encrypted data with metadata
type EncryptedData struct {
	KeyID      string `json:"kid"`
	Ciphertext string `json:"ct"`
	Version    int    `json:"v"`
}

// NewEnhancedEncryptionService creates a new EnhancedEncryptionService
func NewEnhancedEncryptionService(keyManager *KeyManager, logger *zerolog.Logger) *EnhancedEncryptionService {
	return &EnhancedEncryptionService{
		keyManager: keyManager,
		logger:     logger,
	}
}

// Encrypt encrypts a plaintext string using AES-GCM with the current key
func (s *EnhancedEncryptionService) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	// Get current key
	keyID := s.keyManager.GetCurrentKeyID()
	key := s.keyManager.GetCurrentKey()

	// Create cipher
	block, err := aes.NewCipher(key)
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
	encodedCiphertext := base64.StdEncoding.EncodeToString(ciphertext)

	// Create encrypted data
	data := EncryptedData{
		KeyID:      keyID,
		Ciphertext: encodedCiphertext,
		Version:    1,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to marshal encrypted data")
		return "", fmt.Errorf("failed to marshal encrypted data: %w", err)
	}

	return base64.StdEncoding.EncodeToString(jsonData), nil
}

// Decrypt decrypts a ciphertext string using AES-GCM
func (s *EnhancedEncryptionService) Decrypt(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	// Check if it's a legacy format (not base64 encoded JSON)
	if !strings.HasPrefix(ciphertext, "eyJ") {
		// Try to decrypt with all known keys
		return s.decryptLegacy(ciphertext) // Pass only ciphertext
	}

	// Decode base64
	jsonData, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to decode base64 ciphertext")
		return "", fmt.Errorf("failed to decode base64 ciphertext: %w", err)
	}

	// Unmarshal JSON
	var data EncryptedData
	if err := json.Unmarshal(jsonData, &data); err != nil {
		s.logger.Error().Err(err).Msg("Failed to unmarshal encrypted data")
		return "", fmt.Errorf("failed to unmarshal encrypted data: %w", err)
	}

	// Get key
	key, err := s.keyManager.GetKey(data.KeyID)
	if err != nil {
		s.logger.Error().Err(err).Str("keyID", data.KeyID).Msg("Failed to get key")
		return "", fmt.Errorf("failed to get key: %w", err)
	}

	// Decode ciphertext
	decodedCiphertext, err := base64.StdEncoding.DecodeString(data.Ciphertext)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to decode base64 ciphertext")
		return "", fmt.Errorf("failed to decode base64 ciphertext: %w", err)
	}

	// Create cipher
	block, err := aes.NewCipher(key)
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
	if len(decodedCiphertext) < aesGCM.NonceSize() {
		err := errors.New("ciphertext too short")
		s.logger.Error().Err(err).Msg("Invalid ciphertext")
		return "", err
	}

	// Extract nonce and ciphertext
	nonce, ciphertextBytes := decodedCiphertext[:aesGCM.NonceSize()], decodedCiphertext[aesGCM.NonceSize():]

	// Decrypt
	plaintext, err := aesGCM.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to decrypt ciphertext")
		return "", fmt.Errorf("failed to decrypt ciphertext: %w", err)
	}

	return string(plaintext), nil
}

// decryptLegacy attempts to decrypt legacy ciphertext by trying all known keys.
func (s *EnhancedEncryptionService) decryptLegacy(ciphertext string) (string, error) {
	// Decode base64
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		s.logger.Warn().Err(err).Msg("Failed to decode potential legacy base64 ciphertext")
		return "", fmt.Errorf("failed to decode legacy base64 ciphertext: %w", err)
	}

	// Iterate through all known keys
	keyIDs := s.keyManager.ListKeyIDs()
	var lastErr error
	for _, keyID := range keyIDs {
		key, keyErr := s.keyManager.GetKey(keyID)
		if keyErr != nil {
			s.logger.Warn().Err(keyErr).Str("keyID", keyID).Msg("Failed to retrieve key for legacy decryption attempt")
			continue // Try next key
		}

		// Attempt decryption with this key
		plaintext, decryptErr := s.tryLegacyDecryptWithKey(data, key)
		if decryptErr == nil {
			// Success!
			return plaintext, nil
		}
		lastErr = decryptErr // Keep track of the last error (likely auth failed)
	}

	// If no key worked, return the last error encountered
	s.logger.Error().Err(lastErr).Msg("Failed to decrypt legacy ciphertext with any known key")
	if lastErr != nil {
		return "", fmt.Errorf("failed to decrypt legacy ciphertext with any known key: %w", lastErr)
	}
	// Should not happen if keyIDs is not empty, but return a generic error just in case
	return "", errors.New("failed to decrypt legacy ciphertext: no keys available or decryption failed")
}

// tryLegacyDecryptWithKey attempts decryption with a specific key (helper for decryptLegacy)
func (s *EnhancedEncryptionService) tryLegacyDecryptWithKey(data []byte, key []byte) (string, error) {
	// Create cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher for key: %w", err)
	}

	// Create GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM for key: %w", err)
	}

	// Determine nonce size and check data length
	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short for nonce size")
	}

	// Extract nonce and ciphertext using dynamic nonce size
	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]

	// Decrypt, AAD must be nil to match legacy encryption
	// aadBytes := nonce // AAD is NOT the nonce
	// Correctly log the hex representation of the key bytes being used
	s.logger.Debug().Str("key_id_attempt_hex", hex.EncodeToString(key)).Hex("nonce", nonce).Hex("ciphertext", ciphertextBytes).Msg("Attempting aesGCM.Open with nil AAD") // DEBUG LOG
	plaintext, err := aesGCM.Open(nil, nonce, ciphertextBytes, nil)                                                                                                       // Pass nil for AAD
	if err != nil {
		// Don't log error here extensively, as failure is expected when trying wrong keys
		s.logger.Debug().Err(err).Str("key_id_attempt_hex", hex.EncodeToString(key)).Msg("aesGCM.Open failed for key")
		return "", err // Return the raw error (e.g., message authentication failed)
	}

	return string(plaintext), nil
}

// ReEncryptWithCurrentKey re-encrypts data with the current key
func (s *EnhancedEncryptionService) ReEncryptWithCurrentKey(ciphertext string) (string, error) {
	// Decrypt with any key
	plaintext, err := s.Decrypt(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt for re-encryption: %w", err)
	}

	// Encrypt with current key
	return s.Encrypt(plaintext)
}

// Ensure EnhancedEncryptionService implements EncryptionService
var _ port.EncryptionService = (*EnhancedEncryptionService)(nil)
var _ port.EnhancedEncryptionService = (*EnhancedEncryptionService)(nil)
