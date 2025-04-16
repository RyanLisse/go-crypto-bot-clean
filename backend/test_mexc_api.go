package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		// Try parent directory
		err = godotenv.Load("../.env")
		if err != nil {
			fmt.Println("Error loading .env file:", err)
			return
		}
	}

	// Get API credentials from environment variables
	apiKey := os.Getenv("MEXC_API_KEY")
	apiSecret := os.Getenv("MEXC_SECRET_KEY")

	if apiKey == "" || apiSecret == "" {
		fmt.Println("MEXC_API_KEY and MEXC_SECRET_KEY environment variables must be set")
		return
	}

	fmt.Printf("Using API Key: %s...\n", apiKey[:5]+"...")
	fmt.Printf("Using API Secret: %s...\n", apiSecret[:5]+"...")

	// Create timestamp for the request
	timestamp := time.Now().UnixMilli()
	
	// Create query parameters
	params := fmt.Sprintf("timestamp=%d", timestamp)
	
	// Generate signature
	h := hmac.New(sha256.New, []byte(apiSecret))
	h.Write([]byte(params))
	signature := hex.EncodeToString(h.Sum(nil))
	
	// Add signature to parameters
	url := fmt.Sprintf("https://api.mexc.com/api/v3/account?%s&signature=%s", params, signature)
	
	// Create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	
	// Add API key header
	req.Header.Set("X-MBX-APIKEY", apiKey)
	fmt.Println("Set header X-MBX-APIKEY:", apiKey)
	
	// Send request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()
	
	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	
	// Print response
	fmt.Println("Response status:", resp.Status)
	fmt.Println("Response body:", string(body))
}
