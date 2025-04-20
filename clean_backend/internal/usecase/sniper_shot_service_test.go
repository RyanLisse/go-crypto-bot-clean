package usecase

import (
	"context"
	"testing"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/mocks/domain/port" // Assuming mocks are used
	usecase "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"   // Import the package under test

	"github.com/stretchr/testify/assert"
)

// ... (rest of the file remains the same, but usages need prefixing)

func TestSniperShotService_ExecuteShot_Success(t *testing.T) {
	// ... setup mocks ...
	mockExecutor := new(port.MockTradeExecutor)
	mockRepo := new(port.MockSniperShotRepository)
	mockNotifier := new(port.MockNotifier) // Assuming mock notifier

	// Configure mocks for success scenario
	// ... mock calls ...

	// Create service instance
	service := usecase.NewSniperShotService(mockExecutor, mockRepo, mockNotifier, nil) // Prefix with usecase.

	// Define request
	req := &usecase.SniperShotRequest{ // Prefix with usecase.
		// ... fill request fields ...
	}

	// Execute
	result, err := service.ExecuteShot(context.Background(), req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, usecase.SniperShotResultStatusSuccess, result.Status) // Prefix enum with usecase.
	// ... other assertions ...

	// Verify mock expectations
	mockExecutor.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

// ... (Apply similar prefixing to other test functions: TestSniperShotService_ExecuteShot_Failure_Execution, TestSniperShotService_ExecuteShot_Failure_Save, TestSniperShotService_ExecuteShot_Failure_Notification) ...

func TestSniperShotService_HandleRetry(t *testing.T) {
	// ... setup mocks ...
	mockExecutor := new(port.MockTradeExecutor)
	mockRepo := new(port.MockSniperShotRepository)
	mockNotifier := new(port.MockNotifier)

	// Create service instance
	service := usecase.NewSniperShotService(mockExecutor, mockRepo, mockNotifier, nil) // Prefix with usecase.

	// ... setup for retry scenario ...
	shotToRetry := &model.SniperShot{ /* ... fill fields ... */ }

	// Execute retry
	result, err := service.HandleRetry(context.Background(), shotToRetry)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, usecase.SniperShotResultStatusSuccess, result.Status) // Prefix enum with usecase.
	// ... other assertions ...

	// Verify mock expectations
	mockExecutor.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}
