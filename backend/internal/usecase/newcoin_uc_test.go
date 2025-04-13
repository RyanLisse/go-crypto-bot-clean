package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/neo/crypto-bot/internal/domain/event"
	"github.com/neo/crypto-bot/internal/domain/model"
	"github.com/neo/crypto-bot/internal/usecase"
)

// Note: Temporary SymbolInfo is defined in newcoin_uc.go (usecase package)

// MockNewCoinRepository is a mock type for the NewCoinRepository type
type MockNewCoinRepository struct {
	mock.Mock
}

func (m *MockNewCoinRepository) Save(ctx context.Context, coin *model.NewCoin) error {
	args := m.Called(ctx, coin)
	return args.Error(0)
}

func (m *MockNewCoinRepository) FindBySymbol(ctx context.Context, symbol string) (*model.NewCoin, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.NewCoin), args.Error(1)
}

func (m *MockNewCoinRepository) Update(ctx context.Context, coin *model.NewCoin) error {
	args := m.Called(ctx, coin)
	return args.Error(0)
}

func (m *MockNewCoinRepository) GetByStatus(ctx context.Context, status model.NewCoinStatus) ([]*model.NewCoin, error) {
	args := m.Called(ctx, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.NewCoin), args.Error(1)
}

func (m *MockNewCoinRepository) FindRecentlyListed(ctx context.Context, thresholdTime time.Time) ([]*model.NewCoin, error) {
	args := m.Called(ctx, thresholdTime)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.NewCoin), args.Error(1)
}

// MockEventBus is a mock type for the EventBus type
type MockEventBus struct {
	mock.Mock
}

func (m *MockEventBus) Publish(ctx context.Context, e event.DomainEvent) error {
	args := m.Called(ctx, e)
	return args.Error(0)
}

// MockMarketDataService is a mock type for the MarketDataService type
// Assuming MarketDataService has a method like GetSymbolInfo
type MockMarketDataService struct {
	mock.Mock
}

// Assuming a method like GetSymbolInfo exists in MarketDataService port
func (m *MockMarketDataService) GetSymbolInfo(ctx context.Context, symbol string) (*usecase.SymbolInfo, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.SymbolInfo), args.Error(1)
}

// --- Test Cases ---

func TestNewCoinUsecase_CheckNewListings_CoinBecomesTradable(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockNewCoinRepository)
	mockBus := new(MockEventBus)
	mockMarketSvc := new(MockMarketDataService)

	// --- Arrange ---
	symbol := "NEWCOIN/USDT"
	expectedListingTime := time.Now().Add(-1 * time.Hour)
	now := time.Now()

	existingCoin := &model.NewCoin{
		ID:                  "uuid-1",
		Symbol:              symbol,
		ExpectedListingTime: expectedListingTime,
		Status:              model.StatusExpected,
		CreatedAt:           now.Add(-2 * time.Hour),
		UpdatedAt:           now.Add(-2 * time.Hour),
	}
	mockRepo.On("FindRecentlyListed", mock.AnythingOfType("time.Time")).Return([]*model.NewCoin{existingCoin}, nil).Once()

	symbolInfo := &usecase.SymbolInfo{
		Symbol: symbol,
		Status: "TRADING",
	}
	mockMarketSvc.On("GetSymbolInfo", ctx, symbol).Return(symbolInfo, nil).Once()

	mockRepo.On("Update", ctx, mock.MatchedBy(func(coin *model.NewCoin) bool {
		return coin.Symbol == symbol &&
			coin.Status == model.StatusTrading &&
			coin.BecameTradableAt != nil &&
			!coin.IsProcessedForAutobuy
	})).Return(nil).Once()

	var publishedEvent event.DomainEvent
	mockBus.On("Publish", ctx, mock.AnythingOfType("*event.NewCoinTradable")).Run(func(args mock.Arguments) {
		publishedEvent = args.Get(1).(event.DomainEvent)
	}).Return(nil).Once()

	// --- Act ---
	uc := usecase.NewNewCoinUsecase(mockRepo, mockBus, mockMarketSvc)
	err := uc.CheckNewListings(ctx)

	// --- Assert ---
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockMarketSvc.AssertExpectations(t)
	mockBus.AssertExpectations(t)

	assert.NotNil(t, publishedEvent)
	if publishedEvent != nil {
		assert.Equal(t, event.NewCoinTradableEvent, publishedEvent.Type())
		assert.Equal(t, symbol, publishedEvent.AggregateID())
		concreteEvent, ok := publishedEvent.(*event.NewCoinTradable)
		assert.True(t, ok)
		if ok {
			assert.Equal(t, symbol, concreteEvent.Symbol)
			assert.WithinDuration(t, now, concreteEvent.TradableAt, 5*time.Second)
		}
	}
}

// TODO: Define the actual `model.SymbolInfo` struct returned by the market service or adjust the mock MarketDataService to return the correct type and status field.
