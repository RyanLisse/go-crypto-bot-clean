package crypto

import (
	"encoding/hex"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeyManager(t *testing.T) {
	// Setup logger
	logger := zerolog.New(zerolog.NewTestWriter(t))

	t.Run("Create key manager with valid keys", func(t *testing.T) {
		keysJSON := `[
			{"id":"key1","key":"0123456789abcdef0123456789abcdef"},
			{"id":"key2","key":"fedcba9876543210fedcba9876543210"}
		]`
		keyManager, err := NewKeyManager("key1", keysJSON, &logger)
		require.NoError(t, err)
		require.NotNil(t, keyManager)

		// Check current key (expect hex decoded bytes)
		expectedKey1, _ := hex.DecodeString("0123456789abcdef0123456789abcdef")
		assert.Equal(t, "key1", keyManager.GetCurrentKeyID())
		assert.Equal(t, expectedKey1, keyManager.GetCurrentKey())

		// Check key retrieval (expect hex decoded bytes)
		expectedKey2, _ := hex.DecodeString("fedcba9876543210fedcba9876543210")
		key, err := keyManager.GetKey("key2")
		require.NoError(t, err)
		assert.Equal(t, expectedKey2, key)

		// Check key existence
		assert.True(t, keyManager.HasKey("key1"))
		assert.True(t, keyManager.HasKey("key2"))
		assert.False(t, keyManager.HasKey("key3"))

		// Check key IDs
		ids := keyManager.ListKeyIDs()
		assert.Len(t, ids, 2)
		assert.Contains(t, ids, "key1")
		assert.Contains(t, ids, "key2")
	})

	t.Run("Create key manager with invalid keys", func(t *testing.T) {
		// Empty key ID
		_, err := NewKeyManager("", `[{"id":"key1","key":"0123456789abcdef0123456789abcdef"}]`, &logger)
		assert.Error(t, err)

		// Empty keys JSON
		_, err = NewKeyManager("key1", "", &logger)
		assert.Error(t, err)

		// Invalid JSON
		_, err = NewKeyManager("key1", "invalid json", &logger)
		assert.Error(t, err)

		// Current key not found
		_, err = NewKeyManager("key3", `[{"id":"key1","key":"0123456789abcdef0123456789abcdef"}]`, &logger)
		assert.Error(t, err)
	})

	t.Run("Set current key", func(t *testing.T) {
		keysJSON := `[
			{"id":"key1","key":"0123456789abcdef0123456789abcdef"},
			{"id":"key2","key":"fedcba9876543210fedcba9876543210"}
		]`
		keyManager, err := NewKeyManager("key1", keysJSON, &logger)
		require.NoError(t, err)

		// Set current key to key2
		err = keyManager.SetCurrentKey("key2")
		require.NoError(t, err)
		assert.Equal(t, "key2", keyManager.GetCurrentKeyID())
		// Check current key (expect hex decoded bytes)
		expectedKey2, _ := hex.DecodeString("fedcba9876543210fedcba9876543210")
		assert.Equal(t, expectedKey2, keyManager.GetCurrentKey())

		// Try to set current key to non-existent key
		err = keyManager.SetCurrentKey("key3")
		assert.Error(t, err)
	})

	t.Run("Add and remove keys", func(t *testing.T) {
		keysJSON := `[
			{"id":"key1","key":"0123456789abcdef0123456789abcdef"}
		]`
		keyManager, err := NewKeyManager("key1", keysJSON, &logger)
		require.NoError(t, err)

		// Add a new key
		err = keyManager.AddKey("key2", []byte("fedcba9876543210fedcba9876543210"))
		require.NoError(t, err)
		assert.True(t, keyManager.HasKey("key2"))

		// Try to add a key with existing ID
		err = keyManager.AddKey("key2", []byte("newkey"))
		assert.Error(t, err)

		// Try to add a key with empty ID
		err = keyManager.AddKey("", []byte("newkey"))
		assert.Error(t, err)

		// Try to add a key with empty key
		err = keyManager.AddKey("key3", []byte{})
		assert.Error(t, err)

		// Remove a key
		err = keyManager.RemoveKey("key2")
		require.NoError(t, err)
		assert.False(t, keyManager.HasKey("key2"))

		// Try to remove current key
		err = keyManager.RemoveKey("key1")
		assert.Error(t, err)

		// Try to remove non-existent key
		err = keyManager.RemoveKey("key3")
		assert.Error(t, err)
	})
}
