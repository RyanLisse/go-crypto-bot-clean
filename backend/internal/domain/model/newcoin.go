package model

import (
	"time"
)

// Status represents the status of a newly listed coin
type Status string

const (
	// StatusExpected indicates the coin is expected to be listed
	StatusExpected Status = "expected"
	// StatusListed indicates the coin is now listed but not yet tradeable
	StatusListed Status = "listed"
	// StatusTrading indicates active trading has begun
	StatusTrading Status = "trading"
	// StatusFailed indicates listing process failed or was cancelled
	StatusFailed Status = "failed"
)

// NewCoin represents a newly listed cryptocurrency
type NewCoin struct {
	ID                    string     `json:"id"`
	Symbol                string     `json:"symbol"`
	Name                  string     `json:"name"`
	Status                Status     `json:"status"`
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

// MarkAsTradable updates the coin's status to trading and sets the tradable time
func (c *NewCoin) MarkAsTradable(tradableTime time.Time) {
	c.Status = StatusTrading
	c.BecameTradableAt = &tradableTime
	c.UpdatedAt = time.Now()
}

// NewCoinEvent represents an event related to a new coin listing
type NewCoinEvent struct {
	ID        string      `json:"id"`
	CoinID    string      `json:"coin_id"`
	EventType string      `json:"event_type"` // e.g., "status_change", "trading_started"
	OldStatus Status      `json:"old_status,omitempty"`
	NewStatus Status      `json:"new_status,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
}

// NewCoinRepository defines the interface for new coin data persistence
type NewCoinRepository interface {
	// Create stores a new coin in the repository
	Create(coin *NewCoin) error
	// Update updates an existing coin's information
	Update(coin *NewCoin) error
	// GetByID retrieves a coin by its ID
	GetByID(id string) (*NewCoin, error)
	// GetBySymbol retrieves a coin by its trading symbol
	GetBySymbol(symbol string) (*NewCoin, error)
	// List retrieves all coins with optional filtering
	List(status Status, limit, offset int) ([]*NewCoin, error)
	// GetRecent retrieves recently listed coins that are now tradable
	GetRecent(limit int) ([]*NewCoin, error)
	// CreateEvent stores a new coin event
	CreateEvent(event *NewCoinEvent) error
	// GetEvents retrieves events for a specific coin
	GetEvents(coinID string, limit, offset int) ([]*NewCoinEvent, error)
}

// NewCoinService defines the interface for new coin business logic
type NewCoinService interface {
	// DetectNewCoins checks for newly listed coins on MEXC
	DetectNewCoins() error
	// UpdateCoinStatus updates a coin's status and creates an event
	UpdateCoinStatus(coinID string, newStatus Status) error
	// GetCoinDetails retrieves detailed information about a coin
	GetCoinDetails(symbol string) (*NewCoin, error)
	// ListNewCoins retrieves a list of new coins with optional filtering
	ListNewCoins(status Status, limit, offset int) ([]*NewCoin, error)
	// GetRecentTradableCoins retrieves recently listed coins that are now tradable
	GetRecentTradableCoins(limit int) ([]*NewCoin, error)
	// SubscribeToEvents allows subscribing to new coin events
	SubscribeToEvents(callback func(*NewCoinEvent)) error
	// UnsubscribeFromEvents removes an event subscription
	UnsubscribeFromEvents(callback func(*NewCoinEvent)) error
}
