package worker

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockNewCoinWorker is a mock implementation of the NewCoinWorker
type MockNewCoinWorker struct {
	mock.Mock
}

// NewMockNewCoinWorker creates a new mock NewCoinWorker
func NewMockNewCoinWorker() *MockNewCoinWorker {
	return &MockNewCoinWorker{}
}

// Start mocks the Start method
func (m *MockNewCoinWorker) Start(ctx context.Context) {
	m.Called(ctx)
}

// Stop mocks the Stop method
func (m *MockNewCoinWorker) Stop() {
	m.Called()
}
