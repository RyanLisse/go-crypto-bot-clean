package crypto

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEnvEncryptionService is a mock implementation of EncryptionService
type MockEnvEncryptionService struct {
	mock.Mock
}

func (m *MockEnvEncryptionService) Encrypt(plaintext string) ([]byte, error) {
	args := m.Called(plaintext)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockEnvEncryptionService) Decrypt(ciphertext []byte) (string, error) {
	args := m.Called(ciphertext)
	return args.String(0), args.Error(1)
}

func TestEnvManager_SaveAndLoadEnv(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "env-manager-test-*")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create env file path
	envPath := filepath.Join(tempDir, ".env")

	// Create mock encryption service
	mockEncryptionSvc := new(MockEnvEncryptionService)

	// Test data
	vars := map[string]string{
		"TEST_KEY1": "test-value-1",
		"TEST_KEY2": "test-value-2",
	}

	// Mock encryption
	mockEncryptionSvc.On("Encrypt", "test-value-1").Return([]byte("encrypted-value-1"), nil)
	mockEncryptionSvc.On("Encrypt", "test-value-2").Return([]byte("encrypted-value-2"), nil)

	// Mock decryption
	mockEncryptionSvc.On("Decrypt", []byte("encrypted-value-1")).Return("test-value-1", nil)
	mockEncryptionSvc.On("Decrypt", []byte("encrypted-value-2")).Return("test-value-2", nil)

	// Create env manager
	manager := NewEnvManager(mockEncryptionSvc, envPath)

	// Save env variables
	err = manager.SaveEnv(vars, true)
	assert.NoError(t, err)

	// Clear environment variables
	os.Unsetenv("TEST_KEY1")
	os.Unsetenv("TEST_KEY2")

	// Load env variables
	err = manager.LoadEnv()
	assert.NoError(t, err)

	// Check environment variables
	assert.Equal(t, "test-value-1", os.Getenv("TEST_KEY1"))
	assert.Equal(t, "test-value-2", os.Getenv("TEST_KEY2"))

	// Verify mocks
	mockEncryptionSvc.AssertExpectations(t)
}

func TestEnvManager_EncryptAndDecryptEnvFile(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "env-manager-encrypt-test-*")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create env file paths
	plainEnvPath := filepath.Join(tempDir, ".env.plain")
	encryptedEnvPath := filepath.Join(tempDir, ".env.encrypted")
	decryptedEnvPath := filepath.Join(tempDir, ".env.decrypted")

	// Create mock encryption service
	mockEncryptionSvc := new(MockEnvEncryptionService)

	// Create plain env file
	plainEnvContent := `# Environment variables
# Test file

TEST_KEY1=test-value-1
TEST_KEY2=test-value-2
`
	err = os.WriteFile(plainEnvPath, []byte(plainEnvContent), 0600)
	assert.NoError(t, err)

	// Mock encryption
	mockEncryptionSvc.On("Encrypt", "test-value-1").Return([]byte("encrypted-value-1"), nil)
	mockEncryptionSvc.On("Encrypt", "test-value-2").Return([]byte("encrypted-value-2"), nil)

	// Mock decryption
	mockEncryptionSvc.On("Decrypt", []byte("encrypted-value-1")).Return("test-value-1", nil)
	mockEncryptionSvc.On("Decrypt", []byte("encrypted-value-2")).Return("test-value-2", nil)

	// Create env manager
	manager := NewEnvManager(mockEncryptionSvc, "")

	// Encrypt env file
	err = manager.EncryptEnvFile(plainEnvPath, encryptedEnvPath)
	assert.NoError(t, err)

	// Decrypt env file
	err = manager.DecryptEnvFile(encryptedEnvPath, decryptedEnvPath)
	assert.NoError(t, err)

	// Read decrypted env file
	decryptedContent, err := os.ReadFile(decryptedEnvPath)
	assert.NoError(t, err)

	// Check decrypted content
	assert.Contains(t, string(decryptedContent), "TEST_KEY1=test-value-1")
	assert.Contains(t, string(decryptedContent), "TEST_KEY2=test-value-2")

	// Verify mocks
	mockEncryptionSvc.AssertExpectations(t)
}
