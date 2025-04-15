package crypto

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCredentialEncryptionDecryption_Versioned(t *testing.T) {
	r := NewKeyRegistry()
	key1 := []byte("0123456789abcdef0123456789abcdef") // 32 bytes for AES-256
	key2 := []byte("abcdef0123456789abcdef0123456789")

	// Add two keys, v1 active, then v2 active
	assert.NoError(t, r.AddKey("v1", key1, true))
	assert.NoError(t, r.AddKey("v2", key2, true))

	plain := "supersecret"

	// Encrypt with v2 (active)
	ciphertext, keyVersion, err := EncryptCredentialSecret(plain, r)
	assert.NoError(t, err)
	assert.Equal(t, "v2", keyVersion)
	assert.NotEmpty(t, ciphertext)

	// Decrypt with v2
	decrypted, err := DecryptCredentialSecret(ciphertext, keyVersion, r)
	assert.NoError(t, err)
	assert.Equal(t, plain, decrypted)

	// Retire v2 (no active key now)
	assert.NoError(t, r.RetireKey("v2"))
	// Set v1 as active
	r.mu.Lock()
	r.activeVer = "v1"
	r.mu.Unlock()

	// Decrypt ciphertext encrypted with v2 (should still work)
	decrypted, err = DecryptCredentialSecret(ciphertext, "v2", r)
	assert.NoError(t, err)
	assert.Equal(t, plain, decrypted)

	// Encrypt with v1 (now active)
	ciphertext2, keyVersion2, err := EncryptCredentialSecret(plain, r)
	assert.NoError(t, err)
	assert.Equal(t, "v1", keyVersion2)
	decrypted, err = DecryptCredentialSecret(ciphertext2, keyVersion2, r)
	assert.NoError(t, err)
	assert.Equal(t, plain, decrypted)
}
