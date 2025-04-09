// newcoin_service_test.go

package newcoin

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
)

// MockExchangeService simulates an exchange service for testing
type MockExchangeService struct {
	mock.Mock
}

func (m *MockExchangeService) Connect(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockExchangeService) Disconnect() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockExchangeService) GetNewCoins(ctx context.Context) ([]*models.NewCoin, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.NewCoin), args.Error(1)
}

func (m *MockExchangeService) GetTicker(ctx context.Context, symbol string) (*models.Ticker, error) {
	args := m.Called(ctx, symbol)
	return args.Get(0).(*models.Ticker), args.Error(1)
}

func (m *MockExchangeService) GetAllTickers(ctx context.Context) (map[string]*models.Ticker, error) {
	args := m.Called(ctx)
	return args.Get(0).(map[string]*models.Ticker), args.Error(1)
}

func (m *MockExchangeService) GetKlines(ctx context.Context, symbol, interval string, limit int) ([]*models.Kline, error) {
	args := m.Called(ctx, symbol, interval, limit)
	return args.Get(0).([]*models.Kline), args.Error(1)
}

func (m *MockExchangeService) GetWallet(ctx context.Context) (*models.Wallet, error) {
	args := m.Called(ctx)
	return args.Get(0).(*models.Wallet), args.Error(1)
}

func (m *MockExchangeService) PlaceOrder(ctx context.Context, order *models.Order) (*models.Order, error) {
	args := m.Called(ctx, order)
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockExchangeService) CancelOrder(ctx context.Context, orderID, symbol string) error {
	args := m.Called(ctx, orderID, symbol)
	return args.Error(0)
}

func (m *MockExchangeService) GetOrder(ctx context.Context, orderID, symbol string) (*models.Order, error) {
	args := m.Called(ctx, orderID, symbol)
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockExchangeService) GetOpenOrders(ctx context.Context, symbol string) ([]*models.Order, error) {
	args := m.Called(ctx, symbol)
	return args.Get(0).([]*models.Order), args.Error(1)
}

func (m *MockExchangeService) SubscribeToTickers(ctx context.Context, symbols []string, updates chan<- *models.Ticker) error {
	args := m.Called(ctx, symbols, updates)
	return args.Error(0)
}

func (m *MockExchangeService) UnsubscribeFromTickers(ctx context.Context, symbols []string) error {
	args := m.Called(ctx, symbols)
	return args.Error(0)
}

// MockNewCoinRepository simulates a new coin repository for testing
type MockNewCoinRepository struct {
	mock.Mock
}

func (m *MockNewCoinRepository) FindAll(ctx context.Context) ([]models.NewCoin, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.NewCoin), args.Error(1)
}

func (m *MockNewCoinRepository) FindByID(ctx context.Context, id int64) (*models.NewCoin, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.NewCoin), args.Error(1)
}

func (m *MockNewCoinRepository) FindBySymbol(ctx context.Context, symbol string) (*models.NewCoin, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.NewCoin), args.Error(1)
}

func (m *MockNewCoinRepository) Create(ctx context.Context, coin *models.NewCoin) (int64, error) {
	args := m.Called(ctx, coin)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockNewCoinRepository) MarkAsProcessed(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockNewCoinRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockNewCoinRepository) FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]models.NewCoin, error) {
	args := m.Called(ctx, startDate, endDate)
	return args.Get(0).([]models.NewCoin), args.Error(1)
}

func (m *MockNewCoinRepository) FindUpcomingCoins(ctx context.Context) ([]models.NewCoin, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.NewCoin), args.Error(1)
}

func (m *MockNewCoinRepository) FindUpcomingCoinsByDate(ctx context.Context, date time.Time) ([]models.NewCoin, error) {
	args := m.Called(ctx, date)
	return args.Get(0).([]models.NewCoin), args.Error(1)
}

func (m *MockNewCoinRepository) Update(ctx context.Context, coin *models.NewCoin) error {
	args := m.Called(ctx, coin)
	return args.Error(0)
}

func (m *MockNewCoinRepository) FindTradableCoins(ctx context.Context) ([]models.NewCoin, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.NewCoin), args.Error(1)
}

func (m *MockNewCoinRepository) FindTradableCoinsByDate(ctx context.Context, date time.Time) ([]models.NewCoin, error) {
	args := m.Called(ctx, date)
	return args.Get(0).([]models.NewCoin), args.Error(1)
}

// Test cases for NewCoinDetectionService
func TestDetectNewCoins_NewCoinsDetected(t *testing.T) {
	ctx := context.Background()
	mockExchangeService := new(MockExchangeService)
	mockRepo := new(MockNewCoinRepository)

	// Prepare test data
	now := time.Now()
	newCoins := []*models.NewCoin{
		{
			Symbol:      "NEWCOIN1",
			FoundAt:     now,
			BaseVolume:  1000.0,
			QuoteVolume: 2000.0,
		},
		{
			Symbol:      "NEWCOIN2",
			FoundAt:     now,
			BaseVolume:  1500.0,
			QuoteVolume: 3000.0,
		},
	}

	// Set up expectations
	mockExchangeService.On("GetNewCoins", ctx).Return(newCoins, nil)
	mockRepo.On("FindBySymbol", ctx, "NEWCOIN1").Return(nil, errors.New("not found"))
	mockRepo.On("FindBySymbol", ctx, "NEWCOIN2").Return(nil, errors.New("not found"))
	mockRepo.On("Create", ctx, newCoins[0]).Return(int64(1), nil)
	mockRepo.On("Create", ctx, newCoins[1]).Return(int64(2), nil)

	// Create service
	service := NewNewCoinDetectionService(mockExchangeService, mockRepo)

	// Execute
	detectedCoins, err := service.DetectNewCoins(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, detectedCoins, 2)
	assert.Equal(t, "NEWCOIN1", detectedCoins[0].Symbol)
	assert.Equal(t, "NEWCOIN2", detectedCoins[1].Symbol)

	// Verify mock expectations
	mockExchangeService.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestDetectNewCoins_ExistingCoinsNotSaved(t *testing.T) {
	ctx := context.Background()
	mockExchangeService := new(MockExchangeService)
	mockRepo := new(MockNewCoinRepository)

	// Prepare test data
	now := time.Now()
	newCoins := []*models.NewCoin{
		{
			Symbol:      "EXISTINGCOIN",
			FoundAt:     now,
			BaseVolume:  1000.0,
			QuoteVolume: 2000.0,
		},
	}

	// Set up expectations
	mockExchangeService.On("GetNewCoins", ctx).Return(newCoins, nil)
	mockRepo.On("FindBySymbol", ctx, "EXISTINGCOIN").Return(&models.NewCoin{Symbol: "EXISTINGCOIN"}, nil)

	// Create service
	service := NewNewCoinDetectionService(mockExchangeService, mockRepo)

	// Execute
	detectedCoins, err := service.DetectNewCoins(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, detectedCoins, 0)

	// Verify mock expectations
	mockExchangeService.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestDetectNewCoins_ExchangeServiceError(t *testing.T) {
	ctx := context.Background()
	mockExchangeService := new(MockExchangeService)
	mockRepo := new(MockNewCoinRepository)

	// Set up expectations
	mockExchangeService.On("GetNewCoins", ctx).Return(nil, errors.New("exchange service error"))

	// Create service
	service := NewNewCoinDetectionService(mockExchangeService, mockRepo)

	// Execute
	detectedCoins, err := service.DetectNewCoins(ctx)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, detectedCoins)
	assert.Contains(t, err.Error(), "failed to retrieve new coins")

	// Verify mock expectations
	mockExchangeService.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestStartWatching_Cancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	mockExchangeService := new(MockExchangeService)
	mockRepo := new(MockNewCoinRepository)

	// Create service
	service := NewNewCoinDetectionService(mockExchangeService, mockRepo)

	// Set up expectations
	mockExchangeService.On("GetNewCoins", mock.Anything).Return([]*models.NewCoin{}, nil)

	// Start watching in a goroutine
	go func() {
		err := service.StartWatching(ctx, 100*time.Millisecond)
		assert.Equal(t, context.Canceled, err)
	}()

	// Cancel after a short delay
	time.Sleep(250 * time.Millisecond)
	cancel()

	// Verify mock expectations
	mockExchangeService.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestStopWatching(t *testing.T) {
	mockExchangeService := new(MockExchangeService)
	mockRepo := new(MockNewCoinRepository)

	// Create service
	service := NewNewCoinDetectionService(mockExchangeService, mockRepo)

	// Stop watching
	service.StopWatching()

	// No assertions needed, just ensuring no panic occurs
}
