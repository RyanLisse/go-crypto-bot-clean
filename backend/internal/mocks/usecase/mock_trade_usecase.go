package mocks

import (
	"github.com/stretchr/testify/mock"
)

// OrderParameters defines parameters for order execution
type OrderParameters struct {
	Symbol      string
	Side        string
	Type        string
	Quantity    float64
	Price       float64
	UserID      string
	Exchange    string
	IsLeveraged bool
	Leverage    int
}

// MockTradeUsecase is a mock implementation of the TradeUsecase interface
type MockTradeUsecase struct {
	mock.Mock
}

// ExecuteMarketBuy provides a mock function with given fields: order
func (m *MockTradeUsecase) ExecuteMarketBuy(order OrderParameters) error {
	args := m.Called(order)
	return args.Error(0)
}
