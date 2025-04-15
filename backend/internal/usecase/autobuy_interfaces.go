package usecase

// ConfigLoader loads autobuy configuration
type ConfigLoader interface {
	LoadAutoBuyConfig() (*AutoBuyConfig, error)
}

// NewCoinRepository handles coin processing state
type NewCoinRepository interface {
	IsProcessedForAutobuy(symbol string) bool
	MarkAsProcessed(symbol string) error
}

// MarketDataService provides market data for coins
type MarketDataService interface {
	GetMarketData(symbol string) (price float64, volume float64, err error)
}

// RiskUsecase handles risk assessments for orders
type RiskUsecase interface {
	CheckRisk(order OrderParameters) error
}

// TradeUsecase executes trading operations
type TradeUsecase interface {
	ExecuteMarketBuy(order OrderParameters) error
}

// NotificationService sends notifications about events
type NotificationService interface {
	Notify(message string)
}
