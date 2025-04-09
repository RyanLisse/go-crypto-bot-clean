package status

import (
	"context"
	"sync"
)

// MockWatcherStatus is a mock implementation of WatcherStatus
type MockWatcherStatus struct {
	running bool
	mu      sync.Mutex
}

// IsRunning returns whether the watcher is running
func (m *MockWatcherStatus) IsRunning() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.running
}

// Start starts the watcher
func (m *MockWatcherStatus) Start(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.running = true
	return nil
}

// Stop stops the watcher
func (m *MockWatcherStatus) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.running = false
}

// MockStatusProvider is a mock implementation of StatusProvider
type MockStatusProvider struct {
	newCoinWatcher   *MockWatcherStatus
	positionMonitor  *MockWatcherStatus
}

// NewMockStatusProvider creates a new mock status provider
func NewMockStatusProvider() *MockStatusProvider {
	return &MockStatusProvider{
		newCoinWatcher:   &MockWatcherStatus{},
		positionMonitor:  &MockWatcherStatus{},
	}
}

// GetNewCoinWatcher returns the new coin watcher
func (m *MockStatusProvider) GetNewCoinWatcher() WatcherStatus {
	return m.newCoinWatcher
}

// GetPositionMonitor returns the position monitor
func (m *MockStatusProvider) GetPositionMonitor() WatcherStatus {
	return m.positionMonitor
}
