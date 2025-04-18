package mocks

package mocks

import (
	context "context"

	market "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/market"
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// MarketDataServiceInterface is an autogenerated mock type for the MarketDataServiceInterface type
type MarketDataServiceInterface struct {
	mock.Mock
}

// GetHistoricalPrices provides a mock function with given fields: ctx, symbol, startTime, endTime
func (_m *MarketDataServiceInterface) GetHistoricalPrices(ctx context.Context, symbol string, startTime time.Time, endTime time.Time) ([]market.Ticker, error) {
	ret := _m.Called(ctx, symbol, startTime, endTime)

	if len(ret) == 0 {
		panic("no return value specified for GetHistoricalPrices")
	}

	var r0 []market.Ticker
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, time.Time, time.Time) ([]market.Ticker, error)); ok {
		return rf(ctx, symbol, startTime, endTime)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, time.Time, time.Time) []market.Ticker); ok {
		r0 = rf(ctx, symbol, startTime, endTime)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]market.Ticker)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, time.Time, time.Time) error); ok {
		r1 = rf(ctx, symbol, startTime, endTime)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RefreshTicker provides a mock function with given fields: ctx, symbol
func (_m *MarketDataServiceInterface) RefreshTicker(ctx context.Context, symbol string) (*market.Ticker, error) {
	ret := _m.Called(ctx, symbol)

	if len(ret) == 0 {
		panic("no return value specified for RefreshTicker")
	}

	var r0 *market.Ticker
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*market.Ticker, error)); ok {
		return rf(ctx, symbol)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *market.Ticker); ok {
		r0 = rf(ctx, symbol)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*market.Ticker)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, symbol)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMarketDataServiceInterface creates a new instance of MarketDataServiceInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMarketDataServiceInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *MarketDataServiceInterface {
	mock := &MarketDataServiceInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
