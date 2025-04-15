package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/util/crypto"
)

func main() {
	// Define command-line flags
	generateCmd := flag.Bool("generate", false, "Generate a new encryption key")
	rotateCmd := flag.Bool("rotate", false, "Rotate encryption keys")
	bitsFlag := flag.Int("bits", 256, "Key size in bits (must be a multiple of 8)")
	envFlag := flag.Bool("env", false, "Output in environment variable format")

	// Parse flags
	flag.Parse()

	// Create key generator
	keyGen := crypto.NewKeyGenerator()

	// Generate key
	if *generateCmd {
		if *envFlag {
			// Generate key configuration
			config, err := keyGen.GenerateKeyConfig()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error generating key configuration: %v\n", err)
				os.Exit(1)
			}

			// Print configuration
			for k, v := range config {
				fmt.Printf("export %s=\"%s\"\n", k, v)
			}
		} else {
			// Generate single key
			key, err := keyGen.GenerateKey(*bitsFlag)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error generating key: %v\n", err)
				os.Exit(1)
			}

			fmt.Println(key)
		}
		return
	}

	// Rotate keys
	if *rotateCmd {
		// Get current configuration from environment
		currentConfig := make(map[string]string)
		currentConfig["ENCRYPTION_CURRENT_KEY_ID"] = os.Getenv("ENCRYPTION_CURRENT_KEY_ID")
		currentConfig["ENCRYPTION_KEYS"] = os.Getenv("ENCRYPTION_KEYS")

		// Rotate keys
		newConfig, err := keyGen.RotateKeyConfig(currentConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error rotating keys: %v\n", err)
			os.Exit(1)
		}

		// Print new configuration
		for k, v := range newConfig {
			fmt.Printf("export %s=\"%s\"\n", k, v)
		}
		return
	}

	// If no command specified, print usage
	flag.Usage()
}
