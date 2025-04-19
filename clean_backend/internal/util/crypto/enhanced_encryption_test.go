package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnhancedEncryptionService(t *testing.T) {
	// Setup logger
	logger := zerolog.New(zerolog.NewTestWriter(t))

	// Setup key manager with multiple keys
	keysJSON := `[
		{"id":"key1","key":"0123456789abcdef0123456789abcdef"},
		{"id":"key2","key":"fedcba9876543210fedcba9876543210"}
	]`
	keyManager, err := NewKeyManager("key1", keysJSON, &logger)
	require.NoError(t, err)
	require.NotNil(t, keyManager)

	// Create encryption service
	encryptionService := NewEnhancedEncryptionService(keyManager, &logger)
	require.NotNil(t, encryptionService)

	t.Run("Encrypt and decrypt with current key", func(t *testing.T) {
		plaintext := "This is a secret message"
		ciphertext, err := encryptionService.Encrypt(plaintext)
		require.NoError(t, err)
		require.NotEmpty(t, ciphertext)

		decrypted, err := encryptionService.Decrypt(ciphertext)
		require.NoError(t, err)
		assert.Equal(t, plaintext, decrypted)
	})

	t.Run("Encrypt with key1, rotate to key2, and decrypt", func(t *testing.T) {
		// Encrypt with key1
		plaintext := "This is a secret message"
		ciphertext, err := encryptionService.Encrypt(plaintext)
		require.NoError(t, err)
		require.NotEmpty(t, ciphertext)

		// Rotate to key2
		err = keyManager.SetCurrentKey("key2")
		require.NoError(t, err)

		// Decrypt with key2 (should still work because key1 is still available)
		decrypted, err := encryptionService.Decrypt(ciphertext)
		require.NoError(t, err)
		assert.Equal(t, plaintext, decrypted)
	})

	t.Run("Re-encrypt with current key", func(t *testing.T) {
		// Encrypt with key1
		plaintext := "This is a secret message"
		ciphertext, err := encryptionService.Encrypt(plaintext)
		require.NoError(t, err)
		require.NotEmpty(t, ciphertext)

		// Rotate to key2
		err = keyManager.SetCurrentKey("key2")
		require.NoError(t, err)

		// Re-encrypt with key2
		newCiphertext, err := encryptionService.ReEncryptWithCurrentKey(ciphertext)
		require.NoError(t, err)
		require.NotEqual(t, ciphertext, newCiphertext)

		// Decrypt with key2
		decrypted, err := encryptionService.Decrypt(newCiphertext)
		require.NoError(t, err)
		assert.Equal(t, plaintext, decrypted)
	})

	t.Run("Handle empty plaintext", func(t *testing.T) {
		ciphertext, err := encryptionService.Encrypt("")
		require.NoError(t, err)
		assert.Empty(t, ciphertext)

		decrypted, err := encryptionService.Decrypt("")
		require.NoError(t, err)
		assert.Empty(t, decrypted)
	})

	t.Run("Decrypt legacy format simulation", func(t *testing.T) {
		// Reset current key to key1 for predictable encryption
		err := keyManager.SetCurrentKey("key1")
		require.NoError(t, err)

		// Encrypt with EnhancedService using key1
		plaintext := "This is a legacy-style message"
		// Manually create legacy format: base64(nonce + Seal(nonce, nonce, text, nil))
		key1Bytes, _ := keyManager.GetKey("key1")
		block1, _ := aes.NewCipher(key1Bytes)
		aesGCM1, _ := cipher.NewGCM(block1)
		nonce := make([]byte, aesGCM1.NonceSize())
		_, err = io.ReadFull(rand.Reader, nonce)
		require.NoError(t, err)
		t.Logf("Encrypting with Nonce: %x", nonce) // DEBUG LOG
		plaintextBytes := []byte(plaintext)
		aadBytes := nonce // AAD is the nonce itself in legacy format
		t.Logf("Seal Plaintext bytes: %x", plaintextBytes)
		t.Logf("Seal AAD bytes: %x", aadBytes)
		sealedData := aesGCM1.Seal(nonce, nonce, plaintextBytes, nil) // Keep original Seal call
		legacyCiphertext := base64.StdEncoding.EncodeToString(sealedData)
		require.NotEmpty(t, legacyCiphertext)

		// Ensure ciphertext is not JSON-like
		assert.False(t, strings.HasPrefix(legacyCiphertext, "eyJ"))

		// Rotate key manager to key2 (optional, but simulates state change)
		err = keyManager.SetCurrentKey("key2")
		require.NoError(t, err)

		// Decrypt with EnhancedService (should trigger decryptLegacy path)
		decrypted, err := encryptionService.Decrypt(legacyCiphertext)
		require.NoError(t, err, "Decryption failed") // This is the critical check
		assert.Equal(t, plaintext, decrypted)
	})
}
