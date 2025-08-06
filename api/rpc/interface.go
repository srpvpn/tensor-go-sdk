package rpc

import "context"

// RPCAPI defines the interface for RPC operations
type RPCAPI interface {
	// GetPriorityFees retrieves market-based priority fees for transaction creation
	GetPriorityFees(ctx context.Context, req *PriorityFeesRequest) (*PriorityFeesResponse, int, error)
}
