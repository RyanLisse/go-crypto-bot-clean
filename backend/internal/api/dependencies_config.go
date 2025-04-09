package api

import (
	"github.com/ryanlisse/go-crypto-bot/internal/api/handlers"
)

// InitializeConfigDependencies initializes the Config dependencies
func (d *Dependencies) InitializeConfigDependencies() {
	// Create config handler with the loaded config
	d.ConfigHandler = handlers.NewConfigHandler(d.Config)
}
