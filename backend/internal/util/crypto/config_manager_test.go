package crypto

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEncryptionService is a mock implementation of EncryptionService
type MockConfigEncryptionService struct {
	mock.Mock
}

func (m *MockConfigEncryptionService) Encrypt(plaintext string) ([]byte, error) {
	args := m.Called(plaintext)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockConfigEncryptionService) Decrypt(ciphertext []byte) (string, error) {
	args := m.Called(ciphertext)
	return args.String(0), args.Error(1)
}

func TestConfigManager(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "config-manager-test-*")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create config file path
	configPath := filepath.Join(tempDir, "config.enc")

	// Create mock encryption service
	mockEncryptionSvc := new(MockConfigEncryptionService)

	// Test data
	configJSON := `{"key1":"value1","key2":"value2"}`
	encryptedConfig := []byte("encrypted-config")

	// Mock encryption and decryption
	mockEncryptionSvc.On("Encrypt", configJSON).Return(encryptedConfig, nil)
	mockEncryptionSvc.On("Decrypt", encryptedConfig).Return(configJSON, nil)

	// Create config manager
	manager, err := NewConfigManager(mockEncryptionSvc, configPath)
	assert.NoError(t, err)
	assert.NotNil(t, manager)

	// Test setting values
	err = manager.SetValue("key1", "value1")
	assert.NoError(t, err)

	err = manager.SetValue("key2", "value2")
	assert.NoError(t, err)

	// Test getting values
	value1, err := manager.GetValue("key1")
	assert.NoError(t, err)
	assert.Equal(t, "value1", value1)

	value2, err := manager.GetValue("key2")
	assert.NoError(t, err)
	assert.Equal(t, "value2", value2)

	// Test getting non-existent value
	_, err = manager.GetValue("key3")
	assert.Error(t, err)

	// Test getting all values
	values := manager.GetAllValues()
	assert.Equal(t, 2, len(values))
	assert.Equal(t, "value1", values["key1"])
	assert.Equal(t, "value2", values["key2"])

	// Test setting multiple values
	err = manager.SetMultipleValues(map[string]string{
		"key3": "value3",
		"key4": "value4",
	})
	assert.NoError(t, err)

	values = manager.GetAllValues()
	assert.Equal(t, 4, len(values))
	assert.Equal(t, "value3", values["key3"])
	assert.Equal(t, "value4", values["key4"])

	// Test deleting a value
	err = manager.DeleteValue("key1")
	assert.NoError(t, err)

	_, err = manager.GetValue("key1")
	assert.Error(t, err)

	values = manager.GetAllValues()
	assert.Equal(t, 3, len(values))

	// Test clearing all values
	err = manager.Clear()
	assert.NoError(t, err)

	values = manager.GetAllValues()
	assert.Equal(t, 0, len(values))

	// Verify mocks
	mockEncryptionSvc.AssertExpectations(t)
}

func TestConfigManager_LoadConfig(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "config-manager-load-test-*")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create config file path
	configPath := filepath.Join(tempDir, "config.enc")

	// Create mock encryption service
	mockEncryptionSvc := new(MockConfigEncryptionService)

	// Test data
	configJSON := `{"key1":"value1","key2":"value2"}`
	encryptedConfig := []byte("encrypted-config")

	// Create encrypted config file
	err = os.WriteFile(configPath, encryptedConfig, 0600)
	assert.NoError(t, err)

	// Mock decryption
	mockEncryptionSvc.On("Decrypt", encryptedConfig).Return(configJSON, nil)

	// Create config manager
	manager, err := NewConfigManager(mockEncryptionSvc, configPath)
	assert.NoError(t, err)
	assert.NotNil(t, manager)

	// Test getting values
	value1, err := manager.GetValue("key1")
	assert.NoError(t, err)
	assert.Equal(t, "value1", value1)

	value2, err := manager.GetValue("key2")
	assert.NoError(t, err)
	assert.Equal(t, "value2", value2)

	// Verify mocks
	mockEncryptionSvc.AssertExpectations(t)
}
