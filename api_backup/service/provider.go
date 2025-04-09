package service

import (
	"go-crypto-bot-clean/api/repository"
	"go-crypto-bot-clean/backend/pkg/auth"
	"go-crypto-bot-clean/backend/pkg/backtest"
	"go-crypto-bot-clean/backend/pkg/strategy"
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
	userRepo repository.UserRepository,
	strategyRepo repository.StrategyRepository,
	backtestRepo repository.BacktestRepository,
) *Provider {
	return &Provider{
		BacktestService: NewBacktestService(backtestService),
		StrategyService: NewStrategyService(strategyFactory),
		AuthService:     NewAuthService(authService, userRepo),
		UserService:     NewUserService(userRepo),
	}
}
