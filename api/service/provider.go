package service

import (
	"go-crypto-bot-clean/backend/internal/auth"
	"go-crypto-bot-clean/backend/internal/backtest"
	"go-crypto-bot-clean/backend/internal/domain/strategy"
)

// Provider provides access to all services
type Provider struct {
	BacktestService *BacktestService
	StrategyService *StrategyService
	AuthService     *AuthService
	UserService     *UserService
}

// NewProvider creates a new service provider
func NewProvider(
	backtestService *backtest.Service,
	strategyFactory *strategy.Factory,
	authService *auth.Service,
) *Provider {
	return &Provider{
		BacktestService: NewBacktestService(backtestService),
		StrategyService: NewStrategyService(strategyFactory),
		AuthService:     NewAuthService(authService),
		UserService:     NewUserService(),
	}
}
