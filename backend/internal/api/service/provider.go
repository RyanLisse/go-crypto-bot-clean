package service

import (
	"go-crypto-bot-clean/backend/internal/api/repository"
	internalAuth "go-crypto-bot-clean/backend/internal/auth" // Use internal/auth
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
	authProvider internalAuth.AuthProvider, // Expect interface from internal/auth
	userRepo repository.UserRepository,
	strategyRepo repository.StrategyRepository,
	backtestRepo repository.BacktestRepository,
) *Provider {
	return &Provider{
		BacktestService: NewBacktestService(backtestService),
		StrategyService: NewStrategyService(strategyFactory),
		AuthService:     NewAuthService(authProvider, userRepo), // Pass interface
		UserService:     NewUserService(userRepo),
	}
}

// HasBacktestService checks if the backtest service is available
func (p *Provider) HasBacktestService() bool {
	return p != nil && p.BacktestService != nil
}

// HasStrategyService checks if the strategy service is available
func (p *Provider) HasStrategyService() bool {
	return p != nil && p.StrategyService != nil
}

// HasAuthService checks if the authentication service is available
func (p *Provider) HasAuthService() bool {
	return p != nil && p.AuthService != nil
}

// HasUserService checks if the user service is available
func (p *Provider) HasUserService() bool {
	return p != nil && p.UserService != nil
}
