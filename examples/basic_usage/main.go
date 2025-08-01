package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/srpvpn/tensor-go-sdk/api/user"
	"github.com/srpvpn/tensor-go-sdk/client"
)

func main() {
	// Example 1: Basic usage with default configuration
	fmt.Println("=== Tensor SDK Basic Usage Example ===\n")

	// Create a client with default configuration
	// This will use the default base URL and timeout
	tensorClient := client.New(nil)
	defer tensorClient.Close()

	// Example wallet address (this is a well-known Solana address for demonstration)
	walletAddress := "DRpbCBMxVnDK7maPM5tGv6MvB3v1sRMC86PZ8okm21hy"

	// Example 2: Basic portfolio request
	fmt.Printf("Fetching portfolio for wallet: %s\n", walletAddress)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a basic portfolio request
	request := &user.PortfolioRequest{
		Wallet: walletAddress,
	}

	// Execute the request
	response, err := tensorClient.User.GetPortfolio(ctx, request)
	if err != nil {
		handleError("Failed to get portfolio", err)
		return
	}

	// Display basic results
	fmt.Printf("✓ Request successful: %s\n", response.Message)
	fmt.Printf("Found %d collections in portfolio\n\n", len(response.Collections))

	// Display collection details
	if len(response.Collections) > 0 {
		fmt.Println("=== Portfolio Collections ===")
		for i, collection := range response.Collections {
			fmt.Printf("%d. %s (%s)\n", i+1, collection.Name, collection.Symbol)
			fmt.Printf("   Floor Price: %.4f SOL\n", collection.FloorPrice)
			fmt.Printf("   24h Volume: %.4f SOL\n", collection.Volume24h)
			fmt.Printf("   Verified: %t\n", collection.Verified)
			if collection.BidCount != nil {
				fmt.Printf("   Bid Count: %d\n", *collection.BidCount)
			}
			if collection.FavCount != nil {
				fmt.Printf("   Favorite Count: %d\n", *collection.FavCount)
			}
			fmt.Println()
		}
	}

	// Example 3: Advanced usage with custom configuration and options
	fmt.Println("=== Advanced Usage Example ===\n")

	// Create client with custom configuration
	config := &client.Config{
		BaseURL: "https://api.mainnet.tensordev.io",
		Timeout: 45 * time.Second,
		// APIKey: "your-api-key-here", // Uncomment if you have an API key
	}

	advancedClient := client.New(config)
	defer advancedClient.Close()

	// Create an advanced request with all optional parameters
	includeBidCount := true
	includeFavCount := true
	includeUnverified := false
	includeCompressed := true

	advancedRequest := &user.PortfolioRequest{
		Wallet:                walletAddress,
		IncludeBidCount:       &includeBidCount,
		IncludeFavouriteCount: &includeFavCount,
		IncludeUnverified:     &includeUnverified,
		IncludeCompressed:     &includeCompressed,
		Currencies:            []string{"SOL"}, // You can specify multiple currencies
	}

	fmt.Printf("Fetching detailed portfolio for wallet: %s\n", walletAddress)
	fmt.Println("Options: Include bid count, favorite count, compressed collections")

	// Execute the advanced request
	advancedResponse, err := advancedClient.User.GetPortfolio(ctx, advancedRequest)
	if err != nil {
		handleError("Failed to get advanced portfolio", err)
		return
	}

	fmt.Printf("✓ Advanced request successful: %s\n", advancedResponse.Message)
	fmt.Printf("Found %d collections with detailed information\n\n", len(advancedResponse.Collections))

	// Example 4: Error handling scenarios
	fmt.Println("=== Error Handling Examples ===\n")

	// Example with invalid wallet address
	fmt.Println("Testing with invalid wallet address...")
	invalidRequest := &user.PortfolioRequest{
		Wallet: "invalid-wallet-address",
	}

	_, err = tensorClient.User.GetPortfolio(ctx, invalidRequest)
	if err != nil {
		fmt.Printf("✓ Expected error caught: %v\n", err)
	}

	// Example with empty wallet address
	fmt.Println("Testing with empty wallet address...")
	emptyRequest := &user.PortfolioRequest{
		Wallet: "",
	}

	_, err = tensorClient.User.GetPortfolio(ctx, emptyRequest)
	if err != nil {
		fmt.Printf("✓ Expected validation error caught: %v\n", err)
	}

	// Example 5: Context cancellation
	fmt.Println("\nTesting context cancellation...")

	// Create a context that will be cancelled quickly
	shortCtx, shortCancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer shortCancel()

	_, err = tensorClient.User.GetPortfolio(shortCtx, request)
	if err != nil {
		fmt.Printf("✓ Context timeout handled: %v\n", err)
	}

	fmt.Println("\n=== Example completed successfully! ===")
}

// handleError provides centralized error handling with different error types
func handleError(operation string, err error) {
	fmt.Printf("❌ %s: %v\n", operation, err)

	// You can add more sophisticated error handling here
	// For example, checking for specific error types:

	// if apiErr, ok := err.(*errors.APIError); ok {
	//     fmt.Printf("API Error Code: %d\n", apiErr.Code)
	//     fmt.Printf("API Error Message: %s\n", apiErr.Message)
	// }

	// For demonstration purposes, we'll continue execution
	// In a real application, you might want to exit or retry
}

// init function to set up logging and environment
func init() {
	// Set up basic logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Check if we're in a development environment
	if os.Getenv("TENSOR_SDK_DEBUG") == "true" {
		fmt.Println("Debug mode enabled")
	}
}
