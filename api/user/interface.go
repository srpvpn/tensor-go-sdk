package user

import "context"

// Validator defines the interface for request validation
type Validator interface {
	Validate() error
}

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

	// GetTraitBids retrieves all trait bids made by a supplied wallet
	// Returns: response body, status code, error
	GetTraitBids(ctx context.Context, req *TraitBidsRequest) ([]byte, int, error)

	//GetTSwapPools retrieves TSwap pools owned by an address.
	//Returns: response body, status code, error
	GetTSwapPools(ctx context.Context, req *TSwapsPoolsRequest) ([]byte, int, error)

	//GetTAmmPools retrieves TAmm pools owned by an address.
	//Returns: response body, status code, error
	GetTAmmPools(ctx context.Context, req *TAmmPoolsRequest) ([]byte, int, error)

	// GetTransactions retrieves all NFT transactions for a supplied wallet.
	// Returns: response body, status code, error
	GetTransactions(ctx context.Context, req *TransactionsRequest) ([]byte, int, error)

	// GetEscrowAccounts retrieves details for all escrow accounts for a supplied wallet
	// Returns: response body, status code, error
	GetEscrowAccounts(ctx context.Context, req *EscrowAccountsRequest) ([]byte, int, error)

	// GetInventoryForCollection Retrieves details for all NFTs owned by a wallet for a collection
	// Returns: response body, status code, error
	GetInventoryForCollection(ctx context.Context, req *InventoryForCollectionRequest) ([]byte, int, error)

}
