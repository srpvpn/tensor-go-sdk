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
