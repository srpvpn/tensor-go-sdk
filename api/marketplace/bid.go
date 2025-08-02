package marketplace

import (
	"context"
	"encoding/json"
	"fmt"
)

// PlaceNFTBid creates the transaction to place a bid on a single NFT
// Returns: response body, status code, error
func (m *marketplaceAPI) PlaceNFTBid(ctx context.Context, req *PlaceNFTBidRequest) (*PlaceNFTBidResponse, int, error) {
	// Execute the request using the helper method
	body, statusCode, err := m.executeRequest(ctx, "/api/v1/tx/bid", req)
	if err != nil {
		return nil, statusCode, err
	}

	// Parse the JSON response into the structured response
	var response PlaceNFTBidResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, statusCode, fmt.Errorf("failed to parse response JSON: %w", err)
	}

	return &response, statusCode, nil
}

// PlaceTraitBid creates the transaction to place a trait bid on a collection
// Returns: response body, status code, error
func (m *marketplaceAPI) PlaceTraitBid(ctx context.Context, req *PlaceTraitBidRequest) (*PlaceTraitBidResponse, int, error) {
	// Execute the request using the helper method
	body, statusCode, err := m.executeRequest(ctx, "/api/v1/tx/trait_bid", req)
	if err != nil {
		return nil, statusCode, err
	}

	// Parse the JSON response into the structured response
	var response PlaceTraitBidResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, statusCode, fmt.Errorf("failed to parse response JSON: %w", err)
	}

	return &response, statusCode, nil
}

// PlaceCollectionBid creates the transaction to place a collection wide bid
// Returns: response body, status code, error
func (m *marketplaceAPI) PlaceCollectionBid(ctx context.Context, req *PlaceCollectionBidRequest) (*PlaceCollectionBidResponse, int, error) {
	// Execute the request using the helper method
	body, statusCode, err := m.executeRequest(ctx, "/api/v1/tx/collection_bid", req)
	if err != nil {
		return nil, statusCode, err
	}

	// Parse the JSON response into the structured response
	var response PlaceCollectionBidResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, statusCode, fmt.Errorf("failed to parse response JSON: %w", err)
	}

	return &response, statusCode, nil
}
