package main

import (
	"context"
	"fmt"
	"time"

	"github.com/srpvpn/tensor-go-sdk/api/user"
	"github.com/srpvpn/tensor-go-sdk/client"
)

func main() {
	fmt.Println("=== Tensor SDK Offline Demo ===")
	fmt.Println("This demo shows SDK usage patterns without making real API calls")

	// Example 1: Client creation and configuration
	fmt.Println("1. Creating SDK client with custom configuration...")

	config := &client.Config{
		BaseURL: "https://api.mainnet.tensordev.io",
		Timeout: 30 * time.Second,
		// APIKey: "your-api-key-here", // Uncomment when you have an API key
	}

	tensorClient := client.New(config)
	defer tensorClient.Close()

	fmt.Println("✓ Client created successfully")

	// Example 2: Request structure demonstration
	fmt.Println("\n2. Building portfolio requests...")

	// Basic request
	basicRequest := &user.PortfolioRequest{
		Wallet: "DRpbCBMxVnDK7maPM5tGv6MvB3v1sRMC86PZ8okm21hy",
	}

	fmt.Printf("Basic request: Wallet=%s\n", basicRequest.Wallet)

	// Advanced request with all options
	includeBids := true
	includeFavs := true
	includeUnverified := false
	includeCompressed := true

	advancedRequest := &user.PortfolioRequest{
		Wallet:                "DRpbCBMxVnDK7maPM5tGv6MvB3v1sRMC86PZ8okm21hy",
		IncludeBidCount:       &includeBids,
		IncludeFavouriteCount: &includeFavs,
		IncludeUnverified:     &includeUnverified,
		IncludeCompressed:     &includeCompressed,
		Currencies:            []string{"SOL", "USDC"},
	}

	fmt.Printf("Advanced request: Wallet=%s, IncludeBids=%t, IncludeFavs=%t\n",
		advancedRequest.Wallet, *advancedRequest.IncludeBidCount, *advancedRequest.IncludeFavouriteCount)

	// Example 3: Validation demonstration
	fmt.Println("\n3. Request validation...")

	// Valid request
	if err := basicRequest.Validate(); err != nil {
		fmt.Printf("❌ Basic request validation failed: %v\n", err)
	} else {
		fmt.Println("✓ Basic request validation passed")
	}

	// Invalid request (empty wallet)
	invalidRequest := &user.PortfolioRequest{
		Wallet: "",
	}

	if err := invalidRequest.Validate(); err != nil {
		fmt.Printf("✓ Invalid request properly rejected: %v\n", err)
	}

	// Invalid request (bad wallet format)
	badWalletRequest := &user.PortfolioRequest{
		Wallet: "invalid-wallet-123",
	}

	if err := badWalletRequest.Validate(); err != nil {
		fmt.Printf("✓ Bad wallet format properly rejected: %v\n", err)
	}

	// Example 4: Context usage patterns
	fmt.Println("\n4. Context usage patterns...")

	// Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Printf("✓ Context with 5-second timeout created (has deadline: %t)\n", ctx.Err() == nil)

	// Context with cancellation
	cancelCtx, cancelFunc := context.WithCancel(context.Background())

	// Simulate cancellation after 1 second
	go func() {
		time.Sleep(1 * time.Second)
		cancelFunc()
	}()

	fmt.Printf("✓ Cancellable context created (initially active: %t)\n", cancelCtx.Err() == nil)

	// Example 5: Response structure demonstration
	fmt.Println("\n5. Response structure examples...")

	// Create a sample response to show the structure
	sampleResponse := &user.PortfolioResponse{
		Message: "Portfolio retrieved successfully",
		Collections: []user.Collection{
			{
				ID:         "collection-1",
				Name:       "Sample NFT Collection",
				Symbol:     "SAMPLE",
				Image:      "https://example.com/image.png",
				FloorPrice: 1.5,
				Volume24h:  125.75,
				BidCount:   &[]int{42}[0],
				FavCount:   &[]int{156}[0],
				Verified:   true,
				Compressed: false,
			},
			{
				ID:         "collection-2",
				Name:       "Another Collection",
				Symbol:     "ANOTHER",
				Image:      "https://example.com/image2.png",
				FloorPrice: 0.8,
				Volume24h:  89.25,
				Verified:   false,
				Compressed: true,
			},
		},
	}

	fmt.Printf("Sample response message: %s\n", sampleResponse.Message)
	fmt.Printf("Number of collections: %d\n", len(sampleResponse.Collections))

	for i, collection := range sampleResponse.Collections {
		fmt.Printf("  Collection %d: %s (%s)\n", i+1, collection.Name, collection.Symbol)
		fmt.Printf("    Floor Price: %.2f SOL\n", collection.FloorPrice)
		fmt.Printf("    24h Volume: %.2f SOL\n", collection.Volume24h)
		fmt.Printf("    Verified: %t\n", collection.Verified)
		if collection.BidCount != nil {
			fmt.Printf("    Bid Count: %d\n", *collection.BidCount)
		}
		if collection.FavCount != nil {
			fmt.Printf("    Favorite Count: %d\n", *collection.FavCount)
		}
	}

	// Example 6: Error handling patterns
	fmt.Println("\n6. Error handling patterns...")

	// This would be how you handle different types of errors in real usage:
	fmt.Println("In real usage, you would handle errors like this:")
	codeExample := `
	response, err := tensorClient.User.GetPortfolio(ctx, request)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "validation"):
			log.Printf("Validation error: %%v", err)
		case strings.Contains(err.Error(), "timeout"):
			log.Printf("Request timeout: %%v", err)
		case strings.Contains(err.Error(), "403"):
			log.Printf("Authentication error: %%v", err)
		default:
			log.Printf("Unknown error: %%v", err)
		}
		return
	}

	for _, collection := range response.Collections {
		fmt.Printf("Collection: %%s, Floor: %%.2f SOL\n", 
			collection.Name, collection.FloorPrice)
	}
`
fmt.Println(codeExample)


	fmt.Println("\n=== Demo completed successfully! ===")
	fmt.Println("\nTo run with real API calls, use the basic_usage example:")
	fmt.Println("go run examples/basic_usage/main.go")
}
