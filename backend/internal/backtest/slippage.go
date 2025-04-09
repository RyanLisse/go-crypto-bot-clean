package backtest

import (
	"math/rand"
	"time"
)

// SlippageModel defines the interface for simulating slippage
type SlippageModel interface {
	// CalculateSlippage calculates the slippage for a trade
	CalculateSlippage(symbol string, side string, quantity float64, price float64, timestamp time.Time) float64
}

// NoSlippage implements the SlippageModel interface with no slippage
type NoSlippage struct{}

// CalculateSlippage calculates the slippage for a trade (always 0)
func (s *NoSlippage) CalculateSlippage(symbol string, side string, quantity float64, price float64, timestamp time.Time) float64 {
	return 0.0
}

// FixedSlippage implements the SlippageModel interface with a fixed slippage
type FixedSlippage struct {
	SlippagePercent float64
}

// NewFixedSlippage creates a new FixedSlippage model
func NewFixedSlippage(slippagePercent float64) *FixedSlippage {
	return &FixedSlippage{
		SlippagePercent: slippagePercent,
	}
}

// CalculateSlippage calculates the slippage for a trade
func (s *FixedSlippage) CalculateSlippage(symbol string, side string, quantity float64, price float64, timestamp time.Time) float64 {
	return price * s.SlippagePercent / 100.0
}

// VariableSlippage implements the SlippageModel interface with variable slippage
type VariableSlippage struct {
	MinSlippagePercent float64
	MaxSlippagePercent float64
	random             *rand.Rand
}

// NewVariableSlippage creates a new VariableSlippage model
func NewVariableSlippage(minSlippagePercent, maxSlippagePercent float64) *VariableSlippage {
	return &VariableSlippage{
		MinSlippagePercent: minSlippagePercent,
		MaxSlippagePercent: maxSlippagePercent,
		random:             rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// CalculateSlippage calculates the slippage for a trade
func (s *VariableSlippage) CalculateSlippage(symbol string, side string, quantity float64, price float64, timestamp time.Time) float64 {
	// Generate a random slippage percentage between min and max
	slippagePercent := s.MinSlippagePercent + s.random.Float64()*(s.MaxSlippagePercent-s.MinSlippagePercent)
	return price * slippagePercent / 100.0
}

// OrderBookSlippage implements the SlippageModel interface with order book-based slippage
type OrderBookSlippage struct {
	DataProvider DataProvider
}

// NewOrderBookSlippage creates a new OrderBookSlippage model
func NewOrderBookSlippage(dataProvider DataProvider) *OrderBookSlippage {
	return &OrderBookSlippage{
		DataProvider: dataProvider,
	}
}

// CalculateSlippage calculates the slippage for a trade based on order book depth
func (s *OrderBookSlippage) CalculateSlippage(symbol string, side string, quantity float64, price float64, timestamp time.Time) float64 {
	// Get order book snapshot
	orderBook, err := s.DataProvider.GetOrderBook(nil, symbol, timestamp)
	if err != nil || orderBook == nil {
		// Fallback to a default slippage if order book is not available
		return price * 0.001 // 0.1% default slippage
	}

	// Calculate slippage based on order book depth
	var slippage float64
	remainingQuantity := quantity

	if side == "BUY" {
		// For buy orders, we need to look at the asks
		for _, ask := range orderBook.Asks {
			if remainingQuantity <= 0 {
				break
			}

			// Calculate how much we can buy at this price level
			quantityAtLevel := ask.Quantity
			if quantityAtLevel > remainingQuantity {
				quantityAtLevel = remainingQuantity
			}

			// Calculate slippage for this level
			slippage += quantityAtLevel * (ask.Price - price)
			remainingQuantity -= quantityAtLevel
		}
	} else {
		// For sell orders, we need to look at the bids
		for _, bid := range orderBook.Bids {
			if remainingQuantity <= 0 {
				break
			}

			// Calculate how much we can sell at this price level
			quantityAtLevel := bid.Quantity
			if quantityAtLevel > remainingQuantity {
				quantityAtLevel = remainingQuantity
			}

			// Calculate slippage for this level
			slippage += quantityAtLevel * (price - bid.Price)
			remainingQuantity -= quantityAtLevel
		}
	}

	// If we couldn't fill the entire order from the order book, use a default slippage for the remainder
	if remainingQuantity > 0 {
		slippage += remainingQuantity * price * 0.002 // 0.2% default slippage for the remainder
	}

	// Convert total slippage to price slippage
	return slippage / quantity
}
