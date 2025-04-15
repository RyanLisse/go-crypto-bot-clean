package config

import (
	"testing"
	"github.com/rs/zerolog"
)

func TestLoadConfig_ReturnsConfig(t *testing.T) {
	logger := zerolog.Nop()
	cfg := LoadConfig(&logger)
	if cfg == nil {
		t.Error("LoadConfig() returned nil")
	}
}
