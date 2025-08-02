package marketplace

import (
	"context"
	"encoding/json"
	"fmt"
)

// BuyNFT creates the transaction to purchase an NFT
// Returns: response body, status code, error
func (m *marketplaceAPI) BuyNFT(ctx context.Context, req *BuyNFTRequest) (*BuyNFTResponse, int, error) {
	// Execute the request using the helper method
	body, statusCode, err := m.executeRequest(ctx, "/api/v1/tx/buy", req)
	if err != nil {
		return nil, statusCode, err
	}

	// Parse the JSON response into the structured response
	var response BuyNFTResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, statusCode, fmt.Errorf("failed to parse response JSON: %w", err)
	}

	return &response, statusCode, nil
}

// SellNFT creates the transaction to accept a bid on an NFT
// Returns: response body, status code, error
func (m *marketplaceAPI) SellNFT(ctx context.Context, req *SellNFTRequest) (*SellNFTResponse, int, error) {
	// Execute the request using the helper method
	body, statusCode, err := m.executeRequest(ctx, "/api/v1/tx/sell", req)
	if err != nil {
		return nil, statusCode, err
	}

	// Parse the JSON response into the structured response
	var response SellNFTResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, statusCode, fmt.Errorf("failed to parse response JSON: %w", err)
	}

	return &response, statusCode, nil
}

// ListNFT creates the transaction to list an NFT
// Returns: response body, status code, error
func (m *marketplaceAPI) ListNFT(ctx context.Context, req *ListNFTRequest) (*ListNFTResponse, int, error) {
	// Execute the request using the helper method
	body, statusCode, err := m.executeRequest(ctx, "/api/v1/tx/list", req)
	if err != nil {
		return nil, statusCode, err
	}

	// Parse the JSON response into the structured response
	var response ListNFTResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, statusCode, fmt.Errorf("failed to parse response JSON: %w", err)
	}

	return &response, statusCode, nil
}

// DelistNFT creates the transaction to delist an NFT
// Returns: response body, status code, error
func (m *marketplaceAPI) DelistNFT(ctx context.Context, req *DelistNFTRequest) (*DelistNFTResponse, int, error) {
	// Execute the request using the helper method
	body, statusCode, err := m.executeRequest(ctx, "/api/v1/tx/delist", req)
	if err != nil {
		return nil, statusCode, err
	}

	// Parse the JSON response into the structured response
	var response DelistNFTResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, statusCode, fmt.Errorf("failed to parse response JSON: %w", err)
	}

	return &response, statusCode, nil
}

// EditListing creates the transaction to edit an NFT listing
// Returns: response body, status code, error
func (m *marketplaceAPI) EditListing(ctx context.Context, req *EditListingRequest) (*EditListingResponse, int, error) {
	// Execute the request using the helper method
	body, statusCode, err := m.executeRequest(ctx, "/api/v1/tx/edit", req)
	if err != nil {
		return nil, statusCode, err
	}

	// Parse the JSON response into the structured response
	var response EditListingResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, statusCode, fmt.Errorf("failed to parse response JSON: %w", err)
	}

	return &response, statusCode, nil
}
