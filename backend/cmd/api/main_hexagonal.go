package main

import (
	"log"
	"os"

	"go-crypto-bot-clean/backend/internal/api/http"
	"go-crypto-bot-clean/backend/internal/config"
)

func main() {
	// Load configuration
	configPath := getEnv("CONFIG_PATH", "config.yaml")
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create dependency injection container
	container, err := config.NewContainer(cfg)
	if err != nil {
		log.Fatalf("Failed to create container: %v", err)
	}

	// Create and start HTTP server
	server := http.NewServer(container)
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
