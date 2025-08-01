package user

import "context"

// UserAPI defines the interface for user-related API operations
type UserAPI interface {
	// GetPortfolio retrieves portfolio data for a given wallet address
	// Returns: response body, status code, error
	GetPortfolio(ctx context.Context, req *PortfolioRequest) ([]byte, int, error)

	// GetListings retrieves all active listings for a supplied wallet
	// Returns: response body, status code, error
	GetListings(ctx context.Context, req *ListingsRequest) ([]byte, int, error)

	// GetNFTBids retrieves all single NFT bids made by a supplied wallet
	// Returns: response body, status code, error
	GetNFTBids(ctx context.Context, req *NFTBidsRequest) ([]byte, int, error)

	// GetCollectionBids retrieves all collection bids made by a supplied wallet
	// Returns: response body, status code, error
	GetCollectionBids(ctx context.Context, req *CollectionBidsRequest) ([]byte, int, error)
}
