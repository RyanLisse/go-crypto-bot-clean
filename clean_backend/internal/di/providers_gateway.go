package di

import (
	mexcGatewayAdapter "github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/adapter/infrastructure/gateway"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/adapter/infrastructure/gateway/mexc" // Import for MEXCClientV2
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port"
	portgateway "github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port/gateway"
)

// ProvideMEXCGateway creates an instance of the MEXC gateway adapter.
func ProvideMEXCGateway(c *Container) (portgateway.MEXCGateway, error) {
	cfg := c.GetConfig().MEXC
	logger := c.GetLogger()
	// Returns *mexcAdapter which implements portgateway.MEXCGateway
	return mexcGatewayAdapter.NewMEXCGateway(cfg, logger)
}

// ProvideMEXCClient creates an instance that satisfies the port.MEXCClient interface.
func ProvideMEXCClient(c *Container) (port.MEXCClient, error) {
	cfg := c.GetConfig().MEXC
	logger := c.GetLogger()
	// Returns *mexc.MEXCClientV2 which implements port.MEXCClient
	// Ensure MEXCClientV2 constructor exists and returns the correct interface type implicitly or explicitly.
	// Assuming NewMEXCClientV2 exists in the mexc package:
	return mexc.NewMEXCClientV2(cfg.APIKey, cfg.SecretKey, logger), nil // Return as port.MEXCClient
}

// ProvideListingDetector creates an instance that satisfies the port.ListingDetector interface.
func ProvideListingDetector(c *Container) (port.ListingDetector, error) {
	cfg := c.GetConfig().MEXC
	logger := c.GetLogger()
	// Instantiate the concrete type *mexc.MEXCClientV2
	// Go will implicitly check if it satisfies port.ListingDetector upon return.
	concreteClient := mexc.NewMEXCClientV2(cfg.APIKey, cfg.SecretKey, logger)
	return concreteClient, nil // Return the concrete type
}

// ProvideClerkGateway creates an instance of the Clerk gateway.
func ProvideClerkGateway(c *Container) (portgateway.ClerkGateway, error) {
	cfg := c.GetConfig()
	logger := c.GetLogger()
	return mexcGatewayAdapter.NewClerkGateway(cfg, logger), nil
}

// ProvideAIGateway creates an instance of the AI gateway.
func ProvideAIGateway(c *Container) (portgateway.AIGateway, error) {
	cfg := c.GetConfig().AI
	logger := c.GetLogger()
	return mexcGatewayAdapter.NewAIGateway(cfg, logger)
}

// TODO: Add providers for other gateways (Notification, Wallet, System, Trade) if needed.
