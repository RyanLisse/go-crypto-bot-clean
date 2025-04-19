package model

// RateLimit represents a rate limit for an exchange
type RateLimit struct {
	// Type is the type of rate limit (e.g., REQUEST_WEIGHT, ORDERS)
	Type string `json:"rateLimitType"`

	// Interval is the interval for the rate limit (e.g., MINUTE, SECOND)
	Interval string `json:"interval"`

	// IntervalNum is the number of intervals (e.g., 1, 5, 15)
	IntervalNum int `json:"intervalNum"`

	// Limit is the maximum number of requests/orders allowed in the interval
	Limit int `json:"limit"`
}

// ExchangeInfo definition removed from here. See exchange.go for the canonical definition
// related to symbol details and filters.
// type ExchangeInfo struct { ... }
