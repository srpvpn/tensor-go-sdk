package user

import (
	"context"
)

// GetNFTBids retrieves all single NFT bids made by a supplied wallet
// Returns: response body, status code, error
func (u *userAPI) GetNFTBids(ctx context.Context, req *NFTBidsRequest) ([]byte, int, error) {
	return u.executeRequest(ctx, "/api/v1/user/nft_bids", req)
}

// GetCollectionBids retrieves all collection bids made by a supplied wallet
// Returns: response body, status code, error
func (u *userAPI) GetCollectionBids(ctx context.Context, req *CollectionBidsRequest) ([]byte, int, error) {
	return u.executeRequest(ctx, "/api/v1/user/coll_bids", req)
}

// GetTraitBids retrieves all trait bids made by a supplied wallet
// Returns: response body, status code, error
func (u *userAPI) GetTraitBids(ctx context.Context, req *TraitBidsRequest) ([]byte, int, error) {
	return u.executeRequest(ctx, "/api/v1/user/trait_bids", req)
}
