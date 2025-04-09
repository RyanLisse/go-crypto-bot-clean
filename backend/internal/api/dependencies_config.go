package api

import (
	"go-crypto-bot-clean/backend/internal/api/handlers"
)

// InitializeConfigDependencies initializes the Config dependencies
func (d *Dependencies) InitializeConfigDependencies() {
	// Create config handler with the loaded config
	d.ConfigHandler = handlers.NewConfigHandler(d.Config)
}
