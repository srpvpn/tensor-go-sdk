package nfts

import "context"

// NFTsAPI defines the interface for NFTs operations
type NFTsAPI interface {
	// GetNFTsInfo retrieves NFT info based on the mint addresses provided
	GetNFTsInfo(ctx context.Context, req *NFTsInfoRequest) ([]byte, int, error)

	// GetNFTsByCollection retrieves mints based on the collection ID provided
	GetNFTsByCollection(ctx context.Context, req *NFTsByCollectionRequest) ([]byte, int, error)
}
