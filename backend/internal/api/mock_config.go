package api

import "time"

// MockAccountConfig is a mock implementation of the account.Config interface
type MockAccountConfig struct{}

// GetRiskThreshold returns a mock risk threshold
func (m *MockAccountConfig) GetRiskThreshold() float64 {
	return 0.1 // 10% risk threshold
}

// GetCacheTTL returns a mock cache TTL
func (m *MockAccountConfig) GetCacheTTL() time.Duration {
	return 5 * time.Minute // 5 minute cache TTL
}
