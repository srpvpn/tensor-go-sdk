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

	// EditListing creates the transaction to edit an NFT listing
	// Returns: response body, status code, error
	EditListing(ctx context.Context, req *EditListingRequest) (*EditListingResponse, int, error)

	// PlaceNFTBid creates the transaction to place a bid on a single NFT
	// Returns: response body, status code, error
	PlaceNFTBid(ctx context.Context, req *PlaceNFTBidRequest) (*PlaceNFTBidResponse, int, error)

	// PlaceTraitBid creates the transaction to place a trait bid on a collection
	// Returns: response body, status code, error
	PlaceTraitBid(ctx context.Context, req *PlaceTraitBidRequest) (*PlaceTraitBidResponse, int, error)

	// PlaceCollectionBid creates the transaction to place a collection wide bid
	// Returns: response body, status code, error
	PlaceCollectionBid(ctx context.Context, req *PlaceCollectionBidRequest) (*PlaceCollectionBidResponse, int, error)

	// EditBid creates the transaction to edit a bid
	// Returns: response body, status code, error
	EditBid(ctx context.Context, req *EditBidRequest) (*EditBidResponse, int, error)

	// CancelBid creates the transaction to cancel a bid
	// Returns: response body, status code, error
	CancelBid(ctx context.Context, req *CancelBidRequest) (*CancelBidResponse, int, error)
}
