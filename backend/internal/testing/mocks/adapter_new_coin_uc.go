package mocks

import (
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
)

// NewCoinUseCaseAdapter adapts our mock to be used with the actual domain models
type NewCoinUseCaseAdapter struct {
	Mock *MockNewCoinUseCase
}

// DetectNewCoins delegates to the mock
func (a *NewCoinUseCaseAdapter) DetectNewCoins() error {
	return a.Mock.DetectNewCoins()
}

// UpdateCoinStatus converts Status to CoinStatus and delegates to the mock
func (a *NewCoinUseCaseAdapter) UpdateCoinStatus(coinID string, newStatus model.CoinStatus) error {
	// Convert model.CoinStatus to our internal CoinStatus
	return a.Mock.UpdateCoinStatus(coinID, CoinStatus(newStatus))
}

// GetCoinDetails delegates and converts the result
func (a *NewCoinUseCaseAdapter) GetCoinDetails(symbol string) (*model.NewCoin, error) {
	mockCoin, err := a.Mock.GetCoinDetails(symbol)
	if err != nil || mockCoin == nil {
		return nil, err
	}

	// Convert our mock NewCoin to model.NewCoin
	return &model.NewCoin{
		ID:                    mockCoin.ID,
		Symbol:                mockCoin.Symbol,
		Name:                  mockCoin.Name,
		Status:                model.CoinStatus(mockCoin.Status),
		ExpectedListingTime:   mockCoin.ExpectedListingTime,
		BecameTradableAt:      mockCoin.BecameTradableAt,
		BaseAsset:             mockCoin.BaseAsset,
		QuoteAsset:            mockCoin.QuoteAsset,
		MinPrice:              mockCoin.MinPrice,
		MaxPrice:              mockCoin.MaxPrice,
		MinQty:                mockCoin.MinQty,
		MaxQty:                mockCoin.MaxQty,
		PriceScale:            mockCoin.PriceScale,
		QtyScale:              mockCoin.QtyScale,
		IsProcessedForAutobuy: mockCoin.IsProcessedForAutobuy,
		CreatedAt:             mockCoin.CreatedAt,
		UpdatedAt:             mockCoin.UpdatedAt,
	}, nil
}

// ListNewCoins delegates and converts results
func (a *NewCoinUseCaseAdapter) ListNewCoins(status model.CoinStatus, limit, offset int) ([]*model.NewCoin, error) {
	mockCoins, err := a.Mock.ListNewCoins(CoinStatus(status), limit, offset)
	if err != nil || mockCoins == nil {
		return nil, err
	}

	// Convert our mock NewCoin slice to model.NewCoin slice
	result := make([]*model.NewCoin, len(mockCoins))
	for i, mockCoin := range mockCoins {
		result[i] = &model.NewCoin{
			ID:                    mockCoin.ID,
			Symbol:                mockCoin.Symbol,
			Name:                  mockCoin.Name,
			Status:                model.CoinStatus(mockCoin.Status),
			ExpectedListingTime:   mockCoin.ExpectedListingTime,
			BecameTradableAt:      mockCoin.BecameTradableAt,
			BaseAsset:             mockCoin.BaseAsset,
			QuoteAsset:            mockCoin.QuoteAsset,
			MinPrice:              mockCoin.MinPrice,
			MaxPrice:              mockCoin.MaxPrice,
			MinQty:                mockCoin.MinQty,
			MaxQty:                mockCoin.MaxQty,
			PriceScale:            mockCoin.PriceScale,
			QtyScale:              mockCoin.QtyScale,
			IsProcessedForAutobuy: mockCoin.IsProcessedForAutobuy,
			CreatedAt:             mockCoin.CreatedAt,
			UpdatedAt:             mockCoin.UpdatedAt,
		}
	}
	return result, nil
}

// GetRecentTradableCoins delegates and converts results
func (a *NewCoinUseCaseAdapter) GetRecentTradableCoins(limit int) ([]*model.NewCoin, error) {
	mockCoins, err := a.Mock.GetRecentTradableCoins(limit)
	if err != nil || mockCoins == nil {
		return nil, err
	}

	// Convert our mock NewCoin slice to model.NewCoin slice
	result := make([]*model.NewCoin, len(mockCoins))
	for i, mockCoin := range mockCoins {
		result[i] = &model.NewCoin{
			ID:                    mockCoin.ID,
			Symbol:                mockCoin.Symbol,
			Name:                  mockCoin.Name,
			Status:                model.CoinStatus(mockCoin.Status),
			ExpectedListingTime:   mockCoin.ExpectedListingTime,
			BecameTradableAt:      mockCoin.BecameTradableAt,
			BaseAsset:             mockCoin.BaseAsset,
			QuoteAsset:            mockCoin.QuoteAsset,
			MinPrice:              mockCoin.MinPrice,
			MaxPrice:              mockCoin.MaxPrice,
			MinQty:                mockCoin.MinQty,
			MaxQty:                mockCoin.MaxQty,
			PriceScale:            mockCoin.PriceScale,
			QtyScale:              mockCoin.QtyScale,
			IsProcessedForAutobuy: mockCoin.IsProcessedForAutobuy,
			CreatedAt:             mockCoin.CreatedAt,
			UpdatedAt:             mockCoin.UpdatedAt,
		}
	}
	return result, nil
}

// SubscribeToEvents delegates and converts callback
func (a *NewCoinUseCaseAdapter) SubscribeToEvents(callback func(*model.NewCoinEvent)) error {
	// Wrap the callback to convert model.NewCoinEvent to our NewCoinEvent
	wrappedCallback := func(mockEvent *NewCoinEvent) {
		if mockEvent == nil {
			callback(nil)
			return
		}

		// Convert our mock NewCoinEvent to model.NewCoinEvent
		event := &model.NewCoinEvent{
			ID:        mockEvent.ID,
			CoinID:    mockEvent.CoinID,
			EventType: mockEvent.EventType,
			OldStatus: model.CoinStatus(mockEvent.OldStatus),
			NewStatus: model.CoinStatus(mockEvent.NewStatus),
			Data:      mockEvent.Data,
			CreatedAt: mockEvent.CreatedAt,
		}
		callback(event)
	}

	return a.Mock.SubscribeToEvents(wrappedCallback)
}

// UnsubscribeFromEvents delegates and converts callback
func (a *NewCoinUseCaseAdapter) UnsubscribeFromEvents(callback func(*model.NewCoinEvent)) error {
	// We can't directly compare functions, so we'll have to delegate this to the mock
	// In real usage, you'd need to maintain a mapping of callbacks
	// For testing purposes, we'll just accept any callback
	return a.Mock.UnsubscribeFromEvents(func(*NewCoinEvent) {})
}

// NewMockNewCoinUseCase creates a new mock with adapter
func NewMockNewCoinUseCase() *NewCoinUseCaseAdapter {
	return &NewCoinUseCaseAdapter{
		Mock: &MockNewCoinUseCase{},
	}
}
