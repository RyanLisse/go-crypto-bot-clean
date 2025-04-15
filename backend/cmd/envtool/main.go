package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/util/crypto"
)

func main() {
	// Define command-line flags
	encryptCmd := flag.Bool("encrypt", false, "Encrypt an environment file")
	decryptCmd := flag.Bool("decrypt", false, "Decrypt an environment file")
	inputFlag := flag.String("input", "", "Input file path")
	outputFlag := flag.String("output", "", "Output file path")
	keyFlag := flag.String("key", "", "Encryption key (base64-encoded)")

	// Parse flags
	flag.Parse()

	// Validate flags
	if !*encryptCmd && !*decryptCmd {
		fmt.Println("Error: Either -encrypt or -decrypt must be specified")
		flag.Usage()
		os.Exit(1)
	}

	if *inputFlag == "" {
		fmt.Println("Error: Input file path must be specified")
		flag.Usage()
		os.Exit(1)
	}

	if *outputFlag == "" {
		fmt.Println("Error: Output file path must be specified")
		flag.Usage()
		os.Exit(1)
	}

	if *keyFlag == "" {
		// Use environment variable if key flag not provided
		*keyFlag = os.Getenv("ENCRYPTION_KEY")
		if *keyFlag == "" {
			fmt.Println("Error: Encryption key must be provided via -key flag or ENCRYPTION_KEY environment variable")
			os.Exit(1)
		}
	}

	// Set encryption key environment variable
	os.Setenv("ENCRYPTION_KEY", *keyFlag)

	// Create encryption service factory
	factory, err := crypto.NewEncryptionServiceFactory()
	if err != nil {
		fmt.Printf("Error creating encryption service factory: %v\n", err)
		os.Exit(1)
	}

	// Get encryption service
	encryptionSvc, err := factory.GetEncryptionService(crypto.BasicEncryptionService)
	if err != nil {
		fmt.Printf("Error getting encryption service: %v\n", err)
		os.Exit(1)
	}

	// Create environment manager
	envManager := crypto.NewEnvManager(encryptionSvc, "")

	// Perform operation
	if *encryptCmd {
		fmt.Printf("Encrypting %s to %s...\n", *inputFlag, *outputFlag)
		if err := envManager.EncryptEnvFile(*inputFlag, *outputFlag); err != nil {
			fmt.Printf("Error encrypting environment file: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Encryption successful")
	} else {
		fmt.Printf("Decrypting %s to %s...\n", *inputFlag, *outputFlag)
		if err := envManager.DecryptEnvFile(*inputFlag, *outputFlag); err != nil {
			fmt.Printf("Error decrypting environment file: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Decryption successful")
	}
}
