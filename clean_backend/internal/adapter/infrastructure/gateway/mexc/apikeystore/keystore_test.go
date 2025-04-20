package apikeystore

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryKeyStore(t *testing.T) {
	store := NewMemoryKeyStore()

	// Test setting and getting a key
	creds := &APIKeyCredentials{
		APIKey:    "test-api-key",
		SecretKey: "test-secret-key",
	}
	err := store.SetAPIKey("test", creds)
	require.NoError(t, err)

	// Try to get the key
	got, err := store.GetAPIKey("test")
	require.NoError(t, err)
	assert.Equal(t, creds.APIKey, got.APIKey)
	assert.Equal(t, creds.SecretKey, got.SecretKey)

	// Test getting a non-existent key
	_, err = store.GetAPIKey("non-existent")
	assert.ErrorIs(t, err, ErrKeyNotFound)

	// Test deleting a key
	err = store.DeleteAPIKey("test")
	require.NoError(t, err)

	// Verify key was deleted
	_, err = store.GetAPIKey("test")
	assert.ErrorIs(t, err, ErrKeyNotFound)
}

func TestFileKeyStore(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "apikey-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create a file store
	filePath := filepath.Join(tmpDir, "keys.dat")
	encryptionKey := make([]byte, 32) // All zeros for testing
	store, err := NewFileKeyStore(filePath, encryptionKey)
	require.NoError(t, err)

	// Test setting and getting a key
	creds := &APIKeyCredentials{
		APIKey:    "file-api-key",
		SecretKey: "file-secret-key",
	}
	err = store.SetAPIKey("file-test", creds)
	require.NoError(t, err)

	// Try to get the key
	got, err := store.GetAPIKey("file-test")
	require.NoError(t, err)
	assert.Equal(t, creds.APIKey, got.APIKey)
	assert.Equal(t, creds.SecretKey, got.SecretKey)

	// Create a new store instance that should load keys from file
	store2, err := NewFileKeyStore(filePath, encryptionKey)
	require.NoError(t, err)

	// Try to get the key from the new store
	got, err = store2.GetAPIKey("file-test")
	require.NoError(t, err)
	assert.Equal(t, creds.APIKey, got.APIKey)
	assert.Equal(t, creds.SecretKey, got.SecretKey)

	// Test deleting a key
	err = store2.DeleteAPIKey("file-test")
	require.NoError(t, err)

	// Verify key was deleted
	_, err = store2.GetAPIKey("file-test")
	assert.ErrorIs(t, err, ErrKeyNotFound)
}

func TestEnvironmentKeyStore(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("TEST_API_KEY_env-test", "env-api-key")
	os.Setenv("TEST_SECRET_KEY_env-test", "env-secret-key")
	defer func() {
		os.Unsetenv("TEST_API_KEY_env-test")
		os.Unsetenv("TEST_SECRET_KEY_env-test")
	}()

	store := NewEnvironmentKeyStore("TEST_API_KEY_", "TEST_SECRET_KEY_")

	// Try to get the key
	got, err := store.GetAPIKey("env-test")
	require.NoError(t, err)
	assert.Equal(t, "env-api-key", got.APIKey)
	assert.Equal(t, "env-secret-key", got.SecretKey)

	// Test getting a non-existent key
	_, err = store.GetAPIKey("non-existent")
	assert.ErrorIs(t, err, ErrKeyNotFound)

	// Test setting a key (should fail)
	err = store.SetAPIKey("env-test", &APIKeyCredentials{})
	assert.Error(t, err)

	// Test deleting a key (should fail)
	err = store.DeleteAPIKey("env-test")
	assert.Error(t, err)
}

func TestCompositeKeyStore(t *testing.T) {
	// Create stores
	memoryStore := NewMemoryKeyStore()
	envStore := NewEnvironmentKeyStore("TEST_API_KEY_", "TEST_SECRET_KEY_")

	// Set environment variables for testing
	os.Setenv("TEST_API_KEY_env-comp-test", "env-comp-api-key")
	os.Setenv("TEST_SECRET_KEY_env-comp-test", "env-comp-secret-key")
	defer func() {
		os.Unsetenv("TEST_API_KEY_env-comp-test")
		os.Unsetenv("TEST_SECRET_KEY_env-comp-test")
	}()

	// Set a key in memory store
	memCreds := &APIKeyCredentials{
		APIKey:    "mem-comp-api-key",
		SecretKey: "mem-comp-secret-key",
	}
	err := memoryStore.SetAPIKey("mem-comp-test", memCreds)
	require.NoError(t, err)

	// Create composite store with memory first, then env
	compositeStore := NewCompositeKeyStore(memoryStore, envStore)

	// Try to get keys from both sources
	gotMem, err := compositeStore.GetAPIKey("mem-comp-test")
	require.NoError(t, err)
	assert.Equal(t, memCreds.APIKey, gotMem.APIKey)
	assert.Equal(t, memCreds.SecretKey, gotMem.SecretKey)

	gotEnv, err := compositeStore.GetAPIKey("env-comp-test")
	require.NoError(t, err)
	assert.Equal(t, "env-comp-api-key", gotEnv.APIKey)
	assert.Equal(t, "env-comp-secret-key", gotEnv.SecretKey)

	// Test setting a key in composite store (should go to memory store)
	newCreds := &APIKeyCredentials{
		APIKey:    "new-comp-api-key",
		SecretKey: "new-comp-secret-key",
	}
	err = compositeStore.SetAPIKey("new-comp-test", newCreds)
	require.NoError(t, err)

	// Check key exists in memory store
	gotNew, err := memoryStore.GetAPIKey("new-comp-test")
	require.NoError(t, err)
	assert.Equal(t, newCreds.APIKey, gotNew.APIKey)
	assert.Equal(t, newCreds.SecretKey, gotNew.SecretKey)
}

func TestBase64KeyStore(t *testing.T) {
	memoryStore := NewMemoryKeyStore()
	base64Store := NewBase64KeyStore(memoryStore)

	// Test setting and getting a key
	creds := &APIKeyCredentials{
		APIKey:    "base64-api-key",
		SecretKey: "base64-secret-key",
	}
	err := base64Store.SetAPIKey("base64-test", creds)
	require.NoError(t, err)

	// Try to get the key
	got, err := base64Store.GetAPIKey("base64-test")
	require.NoError(t, err)
	assert.Equal(t, creds.APIKey, got.APIKey)
	assert.Equal(t, creds.SecretKey, got.SecretKey)

	// Verify underlying store has encoded values
	raw, err := memoryStore.GetAPIKey("base64-test")
	require.NoError(t, err)
	fmt.Println("Encoded API key:", raw.APIKey)
	fmt.Println("Encoded Secret key:", raw.SecretKey)
	assert.NotEqual(t, creds.APIKey, raw.APIKey)
	assert.NotEqual(t, creds.SecretKey, raw.SecretKey)
}
