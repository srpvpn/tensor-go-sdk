package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/srpvpn/tensor-go-sdk/api/user"
	"github.com/srpvpn/tensor-go-sdk/client"
)

func main() {
	fmt.Println("=== Testing NFT Bids Endpoint ===")

	// Create client
	cfg := client.Config{}
	tensorClient := client.New(&cfg)
	defer tensorClient.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test 1: Valid request
	fmt.Println("Test 1: Valid NFT Bids request")
	validRequest := &user.NFTBidsRequest{
		Owner: "11111111111111111111111111111112", // Valid test wallet
		Limit: 10,
	}

	body, status, err := tensorClient.User.GetNFTBids(ctx, validRequest)
	if err != nil {
		fmt.Printf("Expected error (no API key): %v\n", err)
	} else {
		fmt.Printf("Success! Status: %d, Body: %s\n", status, string(body))
	}

	// Test 2: Validation errors
	fmt.Println("\nTest 2: Validation errors")

	// Empty owner
	invalidRequest1 := &user.NFTBidsRequest{
		Owner: "",
		Limit: 10,
	}
	_, _, err = tensorClient.User.GetNFTBids(ctx, invalidRequest1)
	if err != nil {
		fmt.Printf("✓ Empty owner validation: %v\n", err)
	}

	// Invalid limit
	invalidRequest2 := &user.NFTBidsRequest{
		Owner: "11111111111111111111111111111112",
		Limit: 600, // Too high
	}
	_, _, err = tensorClient.User.GetNFTBids(ctx, invalidRequest2)
	if err != nil {
		fmt.Printf("✓ Invalid limit validation: %v\n", err)
	}

	// Invalid owner address
	invalidRequest3 := &user.NFTBidsRequest{
		Owner: "invalid-address",
		Limit: 10,
	}
	_, _, err = tensorClient.User.GetNFTBids(ctx, invalidRequest3)
	if err != nil {
		fmt.Printf("✓ Invalid owner address validation: %v\n", err)
	}

	// Test 3: Full request with all optional parameters
	fmt.Println("\nTest 3: Full request with optional parameters")
	fullRequest := &user.NFTBidsRequest{
		Owner:        "11111111111111111111111111111112",
		Limit:        50,
		CollId:       stringPtr("some-collection-id"),
		Cursor:       stringPtr("some-cursor"),
		BidAddresses: []string{"11111111111111111111111111111113"},
	}

	_, _, err = tensorClient.User.GetNFTBids(ctx, fullRequest)
	if err != nil {
		fmt.Printf("Expected error (no API key): %v\n", err)
	} else {
		fmt.Println("✓ Full request processed successfully")
	}

	fmt.Println("\n=== All tests completed! ===")
}

func stringPtr(s string) *string {
	return &s
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
