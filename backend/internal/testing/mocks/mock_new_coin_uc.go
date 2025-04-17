package mocks

import (
	"time"

	"github.com/stretchr/testify/mock"
)

// CoinStatus represents the status of a newly listed coin
type CoinStatus string

const (
	// CoinStatusExpected indicates the coin is expected to be listed
	CoinStatusExpected CoinStatus = "expected"
	// CoinStatusListed indicates the coin is now listed but not yet tradeable
	CoinStatusListed CoinStatus = "listed"
	// CoinStatusTrading indicates active trading has begun
	CoinStatusTrading CoinStatus = "trading"
	// CoinStatusFailed indicates listing process failed or was cancelled
	CoinStatusFailed CoinStatus = "failed"
)

// NewCoin represents a newly listed cryptocurrency
type NewCoin struct {
	ID                    string     `json:"id"`
	Symbol                string     `json:"symbol"`
	Name                  string     `json:"name"`
	Status                CoinStatus `json:"status"`
	ExpectedListingTime   time.Time  `json:"expected_listing_time"`
	BecameTradableAt      *time.Time `json:"became_tradable_at,omitempty"`
	BaseAsset             string     `json:"base_asset"`  // e.g., "BTC"
	QuoteAsset            string     `json:"quote_asset"` // e.g., "USDT"
	MinPrice              float64    `json:"min_price"`   // Minimum allowed price
	MaxPrice              float64    `json:"max_price"`   // Maximum allowed price
	MinQty                float64    `json:"min_qty"`     // Minimum order quantity
	MaxQty                float64    `json:"max_qty"`     // Maximum order quantity
	PriceScale            int        `json:"price_scale"` // Price decimal places
	QtyScale              int        `json:"qty_scale"`   // Quantity decimal places
	IsProcessedForAutobuy bool       `json:"is_processed_for_autobuy"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

// NewCoinEvent represents an event related to a new coin listing
type NewCoinEvent struct {
	ID        string      `json:"id"`
	CoinID    string      `json:"coin_id"`
	EventType string      `json:"event_type"` // e.g., "status_change", "trading_started"
	OldStatus CoinStatus  `json:"old_status,omitempty"`
	NewStatus CoinStatus  `json:"new_status,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
}

// MockNewCoinUseCase is a testable mock for the NewCoinUseCase interface
type MockNewCoinUseCase struct {
	mock.Mock
}

// DetectNewCoins mocks the DetectNewCoins method
func (m *MockNewCoinUseCase) DetectNewCoins() error {
	args := m.Called()
	return args.Error(0)
}

// UpdateCoinStatus mocks the UpdateCoinStatus method
func (m *MockNewCoinUseCase) UpdateCoinStatus(coinID string, newStatus CoinStatus) error {
	args := m.Called(coinID, newStatus)
	return args.Error(0)
}

// GetCoinDetails mocks the GetCoinDetails method
func (m *MockNewCoinUseCase) GetCoinDetails(symbol string) (*NewCoin, error) {
	args := m.Called(symbol)

	var coin *NewCoin
	if args.Get(0) != nil {
		coin = args.Get(0).(*NewCoin)
	}

	return coin, args.Error(1)
}

// ListNewCoins mocks the ListNewCoins method
func (m *MockNewCoinUseCase) ListNewCoins(status CoinStatus, limit, offset int) ([]*NewCoin, error) {
	args := m.Called(status, limit, offset)

	var coins []*NewCoin
	if args.Get(0) != nil {
		coins = args.Get(0).([]*NewCoin)
	}

	return coins, args.Error(1)
}

// GetRecentTradableCoins mocks the GetRecentTradableCoins method
func (m *MockNewCoinUseCase) GetRecentTradableCoins(limit int) ([]*NewCoin, error) {
	args := m.Called(limit)

	var coins []*NewCoin
	if args.Get(0) != nil {
		coins = args.Get(0).([]*NewCoin)
	}

	return coins, args.Error(1)
}

// SubscribeToEvents mocks the SubscribeToEvents method
func (m *MockNewCoinUseCase) SubscribeToEvents(callback func(*NewCoinEvent)) error {
	args := m.Called(callback)
	return args.Error(0)
}

// UnsubscribeFromEvents mocks the UnsubscribeFromEvents method
func (m *MockNewCoinUseCase) UnsubscribeFromEvents(callback func(*NewCoinEvent)) error {
	args := m.Called(callback)
	return args.Error(0)
}
