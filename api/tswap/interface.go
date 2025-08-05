package tswap

import (
	"context"
)

// Validator defines the interface for request validation
type Validator interface {
	Validate() error
}

// TSwapAPI defines the interface for TSwap-related operations
type TSwapAPI interface {
	// CloseTSwapPool creates the transaction to close a TSwap pool
	// Returns: response, status code, error
	CloseTSwapPool(ctx context.Context, req *CloseTSwapPoolRequest) (*CloseTSwapPoolResponse, int, error)

	// EditTSwapPool creates the transaction to edit a TSwap pool
	// Returns: response, status code, error
	EditTSwapPool(ctx context.Context, req *EditTSwapPoolRequest) (*EditTSwapPoolResponse, int, error)

	// DepositWithdrawNFT creates the transaction to deposit/withdraw NFT to/from a TSwap pool
	// Returns: response, status code, error
	DepositWithdrawNFT(ctx context.Context, req *DepositWithdrawNFTRequest) (*DepositWithdrawNFTResponse, int, error)

	// DepositWithdrawSOL creates the transaction to deposit/withdraw SOL to/from a TSwap pool
	// Returns: response, status code, error
	DepositWithdrawSOL(ctx context.Context, req *DepositWithdrawSOLRequest) (*DepositWithdrawSOLResponse, int, error)
}
