package factory

import (
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/delivery/http/handler"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/wallet"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// AddressValidatorFactory creates address validator components
type AddressValidatorFactory struct {
	cfg    *config.Config
	logger *zerolog.Logger
	db     *gorm.DB
}

// NewAddressValidatorFactory creates a new AddressValidatorFactory
func NewAddressValidatorFactory(cfg *config.Config, logger *zerolog.Logger, db *gorm.DB) *AddressValidatorFactory {
	return &AddressValidatorFactory{
		cfg:    cfg,
		logger: logger,
		db:     db,
	}
}

// CreateAddressValidatorService creates an address validator service
func (f *AddressValidatorFactory) CreateAddressValidatorService(
	providerRegistry *wallet.ProviderRegistry,
) usecase.AddressValidatorService {
	return usecase.NewAddressValidatorService(
		providerRegistry,
		f.logger,
	)
}

// CreateAddressValidatorHandler creates an address validator handler
func (f *AddressValidatorFactory) CreateAddressValidatorHandler(
	addressValidatorService usecase.AddressValidatorService,
) *handler.AddressValidatorHandler {
	return handler.NewAddressValidatorHandler(
		addressValidatorService,
		f.logger,
	)
}
