package mocks

import (
	"github.com/stretchr/testify/mock"
)

// MockEncryptionService is a mock implementation of the crypto.EncryptionService interface
type MockEncryptionService struct {
	mock.Mock
}

func (m *MockEncryptionService) Encrypt(plaintext string) ([]byte, error) {
	args := m.Called(plaintext)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockEncryptionService) Decrypt(ciphertext []byte) (string, error) {
	args := m.Called(ciphertext)
	return args.String(0), args.Error(1)
}
