package factory

import (
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port"
	portservice "github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port/service"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/usecase"
	"github.com/rs/zerolog"
)

// UseCaseFactory is responsible for creating use case instances with their dependencies.
// It acts as an abstraction over the direct instantiation within the DI container,
// potentially allowing for different implementations (e.g., mocks for testing)
// based on configuration or environment in the future.
type UseCaseFactory struct {
	// Dependencies required by the factory to create use cases
	Logger         *zerolog.Logger
	TradeExecutor  port.TradeExecutor
	OrderRepo      port.OrderRepository
	MarketService  port.MarketDataService    // Needed for PriceChecker for SniperShotService
	TriggerService port.TriggerService       // Needed for SniperShotService
	WalletRepo     port.WalletRepository     // Needed for WalletUseCase
	AuthService    port.AuthServiceInterface // Needed for ProtectedUserUseCase
	// Add other dependencies as needed (e.g., SymbolRepo, etc.)
}

// NewUseCaseFactory creates a new UseCaseFactory.
func NewUseCaseFactory(
	logger *zerolog.Logger,
	tradeExecutor port.TradeExecutor,
	orderRepo port.OrderRepository,
	marketService port.MarketDataService,
	triggerService port.TriggerService,
	walletRepo port.WalletRepository,
	authService port.AuthServiceInterface,
	// Add other dependencies...
) *UseCaseFactory {
	return &UseCaseFactory{
		Logger:         logger,
		TradeExecutor:  tradeExecutor,
		OrderRepo:      orderRepo,
		MarketService:  marketService,
		TriggerService: triggerService,
		WalletRepo:     walletRepo,
		AuthService:    authService,
		// Assign other dependencies...
	}
}

// BuildSniperShotService creates a new instance of the SniperShotService use case.
func (f *UseCaseFactory) BuildSniperShotService() *usecase.SniperShotService {
	// The PriceChecker dependency is created using the MarketDataService
	priceChecker := usecase.NewMarketDataPriceChecker(f.MarketService)

	// Create the real SniperShotService use case with its dependencies
	// Note: TriggerService might be nil if not always required/provided
	return usecase.NewSniperShotService(
		f.Logger,
		f.TradeExecutor,
		f.OrderRepo,
		priceChecker,
		f.TriggerService,
	)
}

// BuildWalletUseCase creates a new instance of the WalletUseCase.
// Note: It returns the port.WalletService interface type, as defined by the use case constructor.
func (f *UseCaseFactory) BuildWalletUseCase() port.WalletService {
	return usecase.NewWalletUseCase(
		f.WalletRepo,
		f.Logger,
	)
}

// BuildProtectedUserUseCase creates a new instance of the ProtectedUserUseCase.
// Note: It returns the ports.ProtectedUserService interface type, as defined by the use case constructor.
func (f *UseCaseFactory) BuildProtectedUserUseCase() portservice.ProtectedUserService {
	return usecase.NewProtectedUserUseCase(
		f.AuthService,
	)
}

// TODO: Add factory methods for other potential use cases
