package escrow

import "context"

// EscrowAPI defines the interface for Shared Escrow operations
type EscrowAPI interface {
	// DepositWithdrawEscrow creates the transaction to deposit or withdraw from an escrow account
	DepositWithdrawEscrow(ctx context.Context, req *DepositWithdrawEscrowRequest) (*DepositWithdrawEscrowResponse, int, error)
}
