package user

import (
	"context"
)

// GetTransactions retrieves all NFT transactions for a supplied wallet.
// Returns: response body, status code, error
func (u *userAPI) GetTransactions(ctx context.Context, req *TransactionsRequest) ([]byte, int, error) {
	return u.executeRequest(ctx, "/api/v1/user/transactions", req)
}
