package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ryanlisse/go-crypto-bot/internal/api"
	"github.com/ryanlisse/go-crypto-bot/internal/config"
)

func main() {
	// Define command-line flags
	defaultPort := os.Getenv("PORT")
	if defaultPort == "" {
		defaultPort = "8080"
	}
	port := flag.String("port", defaultPort, "Port to run the API server on")
	helpFlag := flag.Bool("help", false, "Display help information")

	// Parse command-line flags
	flag.Parse()

	// Display help information if requested
	if *helpFlag || (len(os.Args) > 1 && os.Args[1] == "--help") {
		fmt.Println("Crypto Bot API Server")
		fmt.Println("An API server for the cryptocurrency trading bot.")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  crypto-bot [flags] [command]")
		fmt.Println()
		fmt.Println("Commands:")
		fmt.Println("  serve\t\tStart the API server (default)")
		fmt.Println()
		fmt.Println("Flags:")
		flag.PrintDefaults()
		return
	}

	// Check if a command was provided
	command := "serve"
	if flag.NArg() > 0 {
		command = flag.Arg(0)
	}

	// Load configuration
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		fmt.Printf("Warning: Error loading config: %v\n", err)
		fmt.Println("Continuing with environment variables...")
	}

	// Handle different commands
	switch command {
	case "serve":
		serveAPI(cfg, *port)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Run with --help for usage information.")
	}
}

func serveAPI(cfg *config.Config, port string) {
	// Initialize dependencies
	deps := api.NewDependencies(cfg)

	// Create server
	server := api.NewServer(deps, ":"+port)

	// Start server in a goroutine
	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan

	log.Printf("Received signal: %v", sig)
	log.Println("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Stop(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server stopped")
}
