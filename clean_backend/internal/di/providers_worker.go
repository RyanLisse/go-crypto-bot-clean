package di

import (
	"fmt"

	// Use the correct path for the moved worker:
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/adapter/scheduler"
)

// Worker providers are commented out until the worker package is implemented
/*
// provideNewCoinWorker creates and returns a new NewCoinWorker.
func provideNewCoinWorker(c *Container) (*worker.NewCoinWorker, error) {
	// logger := c.GetLogger()
	// cfg := c.GetConfig()

	// Fetch the MexcSniperService
	sniperService, err := provideMexcSniperService(c) // Fetch from provider
	if err != nil {
		return nil, fmt.Errorf("failed to get MexcSniperService for NewCoinWorker: %w", err)
	}
	if sniperService == nil {
		return nil, fmt.Errorf("MexcSniperService is nil, cannot create NewCoinWorker")
	}

	// Once worker package is implemented, uncomment the following:
	// newCoinWorker := worker.NewNewCoinWorker(sniperService, cfg, *logger)
	// return newCoinWorker, nil

	return nil, fmt.Errorf("worker package not implemented yet")
}
*/

// provideNewCoinWorker creates and returns a new NewCoinWorker.
func provideNewCoinWorker(c *Container) (*scheduler.NewCoinWorker, error) {
	logger := c.GetLogger()
	cfg := c.GetConfig()

	// Fetch the MexcSniperService
	sniperService, err := provideMexcSniperService(c) // Fetch from its provider
	if err != nil {
		return nil, fmt.Errorf("failed to get MexcSniperService for NewCoinWorker: %w", err)
	}
	if sniperService == nil {
		// This check might be redundant if provideMexcSniperService returns an error on nil
		return nil, fmt.Errorf("MexcSniperService is nil, cannot create NewCoinWorker")
	}

	// Create the worker using the concrete type from the scheduler package
	newCoinWorker := scheduler.NewNewCoinWorker(sniperService, cfg, *logger)
	return newCoinWorker, nil
}

// Add providers for other workers here...
