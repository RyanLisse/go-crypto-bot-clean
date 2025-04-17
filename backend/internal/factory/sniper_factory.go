package factory

import (
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	domainservice "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/service"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/service"
	"github.com/rs/zerolog"
)

// SniperFactory creates sniper service components
type SniperFactory struct {
	mexcClient              port.MEXCClient
	symbolRepo              port.SymbolRepository
	orderRepo               port.OrderRepository
	marketService           port.MarketDataService
	listingDetectionService *service.NewListingDetectionService
	logger                  *zerolog.Logger
}

// NewSniperFactory creates a new sniper factory
func NewSniperFactory(
	mexcClient port.MEXCClient,
	symbolRepo port.SymbolRepository,
	orderRepo port.OrderRepository,
	marketService port.MarketDataService,
	listingDetectionService *service.NewListingDetectionService,
	logger *zerolog.Logger,
) *SniperFactory {
	return &SniperFactory{
		mexcClient:              mexcClient,
		symbolRepo:              symbolRepo,
		orderRepo:               orderRepo,
		marketService:           marketService,
		listingDetectionService: listingDetectionService,
		logger:                  logger,
	}
}

// CreateSniperService creates a new sniper service
func (f *SniperFactory) CreateSniperService() port.SniperService {
	// Create logger for the sniper service
	sniperLogger := f.logger.With().Str("component", "sniper_service").Logger()

	// Create and return the sniper service
	return domainservice.NewMexcSniperService(
		f.mexcClient,
		f.symbolRepo,
		f.orderRepo,
		f.marketService,
		f.listingDetectionService,
		&sniperLogger,
	)
}
