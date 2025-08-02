package user

import (
	"context"
)

// GetEscrowAccounts retrieves details for all escrow accounts for a supplied wallet
// Returns: response body, status code, error
func (u *userAPI) GetEscrowAccounts(ctx context.Context, req *EscrowAccountsRequest) ([]byte, int, error) {
	return u.executeRequest(ctx, "/api/v1/user/escrow_accounts", req)
}
