package backtest

import (
	"math/rand"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
)

// // SlippageModel defines the interface for simulating slippage.
// // It uses the price, quantity, and the order side to calculate the slippage.
// type SlippageModel interface {
// 	CalculateSlippage(price float64, quantity float64, side models.OrderSide) float64
// }

// FixedSlippageModel implements a fixed percentage slippage model.
type FixedSlippageModel struct {
	Percentage float64
}

// NewFixedSlippageModel creates a new instance of FixedSlippageModel.
func NewFixedSlippageModel(percentage float64) *FixedSlippageModel {
	return &FixedSlippageModel{Percentage: percentage}
}

// CalculateSlippage calculates the slippage as a fixed percentage of the price.
func (s *FixedSlippageModel) CalculateSlippage(price float64, quantity float64, side models.OrderSide) float64 {
	return price * (s.Percentage / 100.0)
}

// NoSlippage is a slippage model that always returns zero slippage.
type NoSlippage struct{}

// CalculateSlippage for NoSlippage always returns 0.0.
func (s *NoSlippage) CalculateSlippage(price float64, quantity float64, side models.OrderSide) float64 {
	return 0.0
}

// VariableSlippage implements a slippage model with variable, random slippage percentage.
type VariableSlippage struct {
	MinSlippagePercent float64
	MaxSlippagePercent float64
	random             *rand.Rand
}

// NewVariableSlippage creates a new instance of VariableSlippage.
func NewVariableSlippage(minSlippagePercent, maxSlippagePercent float64) *VariableSlippage {
	return &VariableSlippage{
		MinSlippagePercent: minSlippagePercent,
		MaxSlippagePercent: maxSlippagePercent,
		random:             rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// CalculateSlippage calculates slippage using a random percentage between the defined min and max.
func (s *VariableSlippage) CalculateSlippage(price float64, quantity float64, side models.OrderSide) float64 {
	slippagePercent := s.MinSlippagePercent + s.random.Float64()*(s.MaxSlippagePercent-s.MinSlippagePercent)
	return price * slippagePercent / 100.0
}

// // OrderBookSlippage implements the SlippageModel interface with order book-based slippage
// // NOTE: This requires DataProvider to have a GetOrderBook method, which it currently doesn't.
// type OrderBookSlippage struct {
// 	DataProvider DataProvider
// }

// // NewOrderBookSlippage creates a new OrderBookSlippage model
// func NewOrderBookSlippage(dataProvider DataProvider) *OrderBookSlippage {
// 	return &OrderBookSlippage{
// 		DataProvider: dataProvider,
// 	}
// }

// // CalculateSlippage calculates the slippage for a trade based on order book depth
// // NOTE: Requires symbol and timestamp, which are not part of the current interface signature.
// func (s *OrderBookSlippage) CalculateSlippage(price float64, quantity float64, side models.OrderSide) float64 {
// 	 symbol := "" // Placeholder - need symbol
// 	 timestamp := time.Now() // Placeholder - need timestamp

// 	// Get order book snapshot
// 	orderBook, err := s.DataProvider.GetOrderBook(context.Background(), symbol, timestamp)
// 	if err != nil || orderBook == nil {
// 		// Fallback to a default slippage if order book is not available
// 		fmt.Printf("Error getting order book for slippage calc: %v\n", err)
// 		return price * 0.001 // 0.1% default slippage
// 	}

// 	// Calculate slippage based on order book depth
// 	var slippage float64
// 	remainingQuantity := quantity

// 	if side == models.OrderSideBuy {
// 		// For buy orders, we need to look at the asks
// 		for _, ask := range orderBook.Asks {
// 			if remainingQuantity <= 0 {
// 				break
// 			}
// 			quantityAtLevel := ask.Quantity
// 			if quantityAtLevel > remainingQuantity {
// 				quantityAtLevel = remainingQuantity
// 			}
// 			slippage += quantityAtLevel * (ask.Price - price)
// 			remainingQuantity -= quantityAtLevel
// 		}
// 	} else {
// 		// For sell orders, we need to look at the bids
// 		for _, bid := range orderBook.Bids {
// 			if remainingQuantity <= 0 {
// 				break
// 			}
// 			quantityAtLevel := bid.Quantity
// 			if quantityAtLevel > remainingQuantity {
// 				quantityAtLevel = remainingQuantity
// 			}
// 			slippage += quantityAtLevel * (price - bid.Price)
// 			remainingQuantity -= quantityAtLevel
// 		}
// 	}

// 	if remainingQuantity > 0 {
// 		slippage += remainingQuantity * price * 0.002 // 0.2% default slippage for the remainder
// 	}

// 	// Convert total slippage to price slippage per unit
// 	if quantity == 0 { return 0.0 }
// 	return slippage / quantity
// }

// // calculateDepth calculates the cumulative quantity within a certain price range
// // Note: Needs models.OrderBookEntry
// func calculateDepth(entries []models.OrderBookEntry, price float64) float64 {
// 	var depth float64
// 	priceRange := price * 0.01 // Consider depth within 1% of the price

// 	for _, entry := range entries {
// 		if math.Abs(entry.Price-price) <= priceRange {
// 			depth += entry.Quantity
// 		}
// 	}
// 	return depth
// }

// Interface compliance checks.
var (
	_ SlippageModel = (*FixedSlippageModel)(nil)
	_ SlippageModel = (*NoSlippage)(nil)
	_ SlippageModel = (*VariableSlippage)(nil)
)

// var _ SlippageModel = (*OrderBookSlippage)(nil) // Keep commented
