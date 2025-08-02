package marketplace

import "context"

// Validator defines the interface for request validation
type Validator interface {
	Validate() error
}

// MarketplaceAPI defines the interface for marketplace-related API operations
type MarketplaceAPI interface {
	// BuyNFT creates the transaction to purchase an NFT
	// Returns: response body, status code, error
	BuyNFT(ctx context.Context, req *BuyNFTRequest) (*BuyNFTResponse, int, error)
}
