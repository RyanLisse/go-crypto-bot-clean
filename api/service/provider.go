package service

import (
	"go-crypto-bot-clean/api/repository"
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
	userRepo repository.UserRepository,
	strategyRepo repository.StrategyRepository,
	backtestRepo repository.BacktestRepository,
) *Provider {
	return &Provider{
		BacktestService: NewBacktestService(backtestService, backtestRepo),
		StrategyService: NewStrategyService(strategyFactory, strategyRepo),
		AuthService:     NewAuthService(authService, userRepo),
		UserService:     NewUserService(userRepo),
	}
}
