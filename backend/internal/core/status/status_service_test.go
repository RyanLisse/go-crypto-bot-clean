package status

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestWatcherStatus is a mock implementation of WatcherStatus for testing
type TestWatcherStatus struct {
	mock.Mock
}

func (m *TestWatcherStatus) IsRunning() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *TestWatcherStatus) Start(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *TestWatcherStatus) Stop() {
	m.Called()
}

// TestStatusProvider is a mock implementation of StatusProvider for testing
type TestStatusProvider struct {
	mock.Mock
}

func (m *TestStatusProvider) GetNewCoinWatcher() WatcherStatus {
	args := m.Called()
	return args.Get(0).(WatcherStatus)
}

func (m *TestStatusProvider) GetPositionMonitor() WatcherStatus {
	args := m.Called()
	return args.Get(0).(WatcherStatus)
}

func TestStatusService_GetStatus(t *testing.T) {
	// Create mocks
	mockNewCoinWatcher := new(TestWatcherStatus)
	mockPositionMonitor := new(TestWatcherStatus)
	mockProvider := new(TestStatusProvider)

	// Setup expectations
	mockNewCoinWatcher.On("IsRunning").Return(true)
	mockPositionMonitor.On("IsRunning").Return(true)
	mockProvider.On("GetNewCoinWatcher").Return(mockNewCoinWatcher)
	mockProvider.On("GetPositionMonitor").Return(mockPositionMonitor)

	// Create service
	service := NewStatusService(mockProvider, "1.0.0")

	// Get status
	status, err := service.GetStatus()

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, status)
	assert.Equal(t, "1.0.0", status.SystemInfo.Version)
	assert.Equal(t, "healthy", status.OverallStatus)
	assert.Len(t, status.Components, 2)
	assert.Equal(t, "NewCoinWatcher", status.Components[0].Name)
	assert.Equal(t, "PositionMonitor", status.Components[1].Name)
	assert.True(t, status.Components[0].IsRunning)
	assert.True(t, status.Components[1].IsRunning)

	// Verify mocks
	mockNewCoinWatcher.AssertExpectations(t)
	mockPositionMonitor.AssertExpectations(t)
	mockProvider.AssertExpectations(t)
}

func TestStatusService_GetStatus_Degraded(t *testing.T) {
	// Create mocks
	mockNewCoinWatcher := new(TestWatcherStatus)
	mockPositionMonitor := new(TestWatcherStatus)
	mockProvider := new(TestStatusProvider)

	// Setup expectations
	mockNewCoinWatcher.On("IsRunning").Return(true)
	mockPositionMonitor.On("IsRunning").Return(false)
	mockProvider.On("GetNewCoinWatcher").Return(mockNewCoinWatcher)
	mockProvider.On("GetPositionMonitor").Return(mockPositionMonitor)

	// Create service
	service := NewStatusService(mockProvider, "1.0.0")

	// Get status
	status, err := service.GetStatus()

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, status)
	assert.Equal(t, "degraded", status.OverallStatus)
	assert.True(t, status.Components[0].IsRunning)
	assert.False(t, status.Components[1].IsRunning)

	// Verify mocks
	mockNewCoinWatcher.AssertExpectations(t)
	mockPositionMonitor.AssertExpectations(t)
	mockProvider.AssertExpectations(t)
}

func TestStatusService_StartProcesses(t *testing.T) {
	// Create mocks
	mockNewCoinWatcher := new(TestWatcherStatus)
	mockPositionMonitor := new(TestWatcherStatus)
	mockProvider := new(TestStatusProvider)

	// Setup expectations
	mockNewCoinWatcher.On("IsRunning").Return(false)
	mockNewCoinWatcher.On("Start", mock.Anything).Return(nil)
	mockPositionMonitor.On("IsRunning").Return(false)
	mockPositionMonitor.On("Start", mock.Anything).Return(nil)
	mockProvider.On("GetNewCoinWatcher").Return(mockNewCoinWatcher)
	mockProvider.On("GetPositionMonitor").Return(mockPositionMonitor)

	// For GetStatus after starting
	mockNewCoinWatcher.On("IsRunning").Return(true)
	mockPositionMonitor.On("IsRunning").Return(true)

	// Create service
	service := NewStatusService(mockProvider, "1.0.0")

	// Start processes
	status, err := service.StartProcesses(context.Background())

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, status)

	// Verify mocks
	mockNewCoinWatcher.AssertExpectations(t)
	mockPositionMonitor.AssertExpectations(t)
	mockProvider.AssertExpectations(t)
}

func TestStatusService_StopProcesses(t *testing.T) {
	// Create mocks
	mockNewCoinWatcher := new(TestWatcherStatus)
	mockPositionMonitor := new(TestWatcherStatus)
	mockProvider := new(TestStatusProvider)

	// Setup expectations
	mockNewCoinWatcher.On("IsRunning").Return(true)
	mockNewCoinWatcher.On("Stop").Return()
	mockPositionMonitor.On("IsRunning").Return(true)
	mockPositionMonitor.On("Stop").Return()
	mockProvider.On("GetNewCoinWatcher").Return(mockNewCoinWatcher)
	mockProvider.On("GetPositionMonitor").Return(mockPositionMonitor)

	// For GetStatus after stopping
	mockNewCoinWatcher.On("IsRunning").Return(false)
	mockPositionMonitor.On("IsRunning").Return(false)

	// Create service
	service := NewStatusService(mockProvider, "1.0.0")

	// Stop processes
	status, err := service.StopProcesses()

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, status)

	// Verify mocks
	mockNewCoinWatcher.AssertExpectations(t)
	mockPositionMonitor.AssertExpectations(t)
	mockProvider.AssertExpectations(t)
}
