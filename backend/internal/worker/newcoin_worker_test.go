package worker

import (
	"context"
	"testing"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	mocksuc "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/mocks/usecase"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestNewCoinWorker_Start(t *testing.T) {
	// Create mock
	mockNewCoinUC := new(mocksuc.MockNewCoinUseCase)

	// Set up expectations
	mockNewCoinUC.On("DetectNewCoins").Return(nil)

	// Create config and logger
	cfg := &config.Config{}
	logger := zerolog.New(zerolog.NewTestWriter(t))

	// Create worker
	worker := NewNewCoinWorker(mockNewCoinUC, cfg, logger)

	// Start worker in separate goroutine with a context that will be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	go worker.Start(ctx)

	// Sleep to allow the worker to process at least one detection
	time.Sleep(100 * time.Millisecond)

	// Stop the worker
	cancel()
	worker.Stop()

	// Verify expectations
	time.Sleep(100 * time.Millisecond) // Give time for worker to stop
	mockNewCoinUC.AssertExpectations(t)
}

func TestNewCoinWorker_Stop(t *testing.T) {
	// Create mock
	mockNewCoinUC := new(mocksuc.MockNewCoinUseCase)

	// Set up expectations - we may not call this if worker stops quickly
	mockNewCoinUC.On("DetectNewCoins").Return(nil).Maybe()

	// Create config and logger
	cfg := &config.Config{}
	logger := zerolog.New(zerolog.NewTestWriter(t))

	// Create worker
	worker := NewNewCoinWorker(mockNewCoinUC, cfg, logger)

	// Start worker in separate goroutine
	ctx := context.Background()
	go func() {
		time.Sleep(50 * time.Millisecond) // Let it start
		worker.Stop()                     // Stop it
	}()

	// Run it
	worker.Start(ctx)

	// Assert that stopped flag was set
	assert.True(t, worker.stopped)
}

func TestNewCoinWorker_runDetection(t *testing.T) {
	// Create mock
	mockNewCoinUC := new(mocksuc.MockNewCoinUseCase)

	// Set up expectations
	mockNewCoinUC.On("DetectNewCoins").Return(nil)

	// Create config and logger
	cfg := &config.Config{}
	logger := zerolog.New(zerolog.NewTestWriter(t))

	// Create worker
	worker := NewNewCoinWorker(mockNewCoinUC, cfg, logger)

	// Run detection
	ctx := context.Background()
	worker.runDetection(ctx)

	// Verify expectations
	mockNewCoinUC.AssertExpectations(t)
}

func TestNewCoinWorker_runDetection_Error(t *testing.T) {
	// Create mock
	mockNewCoinUC := new(mocksuc.MockNewCoinUseCase)

	// Set up expectations - simulate an error
	expectedErr := error(assert.AnError)
	mockNewCoinUC.On("DetectNewCoins").Return(expectedErr)

	// Create config and logger
	cfg := &config.Config{}
	logger := zerolog.New(zerolog.NewTestWriter(t))

	// Create worker
	worker := NewNewCoinWorker(mockNewCoinUC, cfg, logger)

	// Run detection
	ctx := context.Background()
	worker.runDetection(ctx)

	// Verify expectations - should handle error gracefully
	mockNewCoinUC.AssertExpectations(t)
}
