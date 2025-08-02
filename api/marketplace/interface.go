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

	// SellNFT creates the transaction to accept a bid on an NFT
	// Returns: response body, status code, error
	SellNFT(ctx context.Context, req *SellNFTRequest) (*SellNFTResponse, int, error)

	// ListNFT creates the transaction to list an NFT
	// Returns: response body, status code, error
	ListNFT(ctx context.Context, req *ListNFTRequest) (*ListNFTResponse, int, error)

	// DelistNFT creates the transaction to delist an NFT
	// Returns: response body, status code, error
	DelistNFT(ctx context.Context, req *DelistNFTRequest) (*DelistNFTResponse, int, error)
}
