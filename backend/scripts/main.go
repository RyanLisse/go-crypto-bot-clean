package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/scripts/scripts"
)

func main() {
	// Get command line arguments
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Please specify a script to run:")
		fmt.Println("  check-mexc-balance - Check MEXC account balance")
		fmt.Println("  test-mexc-endpoints - Test MEXC API endpoints")
		return
	}

	// Run the specified script
	switch strings.ToLower(args[0]) {
	case "check-mexc-balance":
		scripts.CheckMexcBalance()
	case "test-mexc-endpoints":
		scripts.TestMexcEndpoints()
	default:
		fmt.Printf("Unknown script: %s\n", args[0])
		fmt.Println("Available scripts:")
		fmt.Println("  check-mexc-balance - Check MEXC account balance")
		fmt.Println("  test-mexc-endpoints - Test MEXC API endpoints")
	}
}
