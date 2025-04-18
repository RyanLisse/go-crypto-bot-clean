package mocks

package mocks

import mock "github.com/stretchr/testify/mock"

// MarketDataService is an autogenerated mock type for the MarketDataService type
type MarketDataService struct {
	mock.Mock
}

// GetMarketData provides a mock function with given fields: symbol
func (_m *MarketDataService) GetMarketData(symbol string) (float64, float64, error) {
	ret := _m.Called(symbol)

	if len(ret) == 0 {
		panic("no return value specified for GetMarketData")
	}

	var r0 float64
	var r1 float64
	var r2 error
	if rf, ok := ret.Get(0).(func(string) (float64, float64, error)); ok {
		return rf(symbol)
	}
	if rf, ok := ret.Get(0).(func(string) float64); ok {
		r0 = rf(symbol)
	} else {
		r0 = ret.Get(0).(float64)
	}

	if rf, ok := ret.Get(1).(func(string) float64); ok {
		r1 = rf(symbol)
	} else {
		r1 = ret.Get(1).(float64)
	}

	if rf, ok := ret.Get(2).(func(string) error); ok {
		r2 = rf(symbol)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// NewMarketDataService creates a new instance of MarketDataService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMarketDataService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MarketDataService {
	mock := &MarketDataService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
