package user

import (
	"context"
)

// GetTSwapPools retrieves TSwap pools owned by an address.
// Returns: response body, status code, error
func (u *userAPI) GetTSwapPools(ctx context.Context, req *TSwapsPoolsRequest) ([]byte, int, error) {
	return u.executeRequest(ctx, "/api/v1/user/amm_pools", req)
}
