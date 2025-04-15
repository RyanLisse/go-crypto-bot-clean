package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockTransactionManager is a mock implementation of port.TransactionManager
type MockTransactionManager struct {
	mock.Mock
}

// WithTransaction mocks the WithTransaction method
func (m *MockTransactionManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}
