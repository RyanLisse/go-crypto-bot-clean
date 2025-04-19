package factory

import (
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/adapter/infrastructure/gateway/mexc"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port"
)

// GetMEXCClient returns either a real or mock MEXC client
func (f *AppFactory) GetMEXCClient() port.MEXCClient {
	return f.GetComponent("mexc_client", func() interface{} {
		if f.ShouldUseMock("mexc_client") {
			f.LogMockUsage("MEXCClient")
			return mexc.NewMockMEXCClientV2(f.logger)
		}

		return mexc.NewMEXCClientV2(
			f.config.MEXC.APIKey,
			f.config.MEXC.SecretKey,
			f.logger,
		)
	}).(port.MEXCClient)
}
