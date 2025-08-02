package user

import (
	"context"
)

// GetTAmmPools retrieves TSwap pools owned by an address.
// Returns: response body, status code, error
func (u *userAPI) GetTAmmPools(ctx context.Context, req *TAmmPoolsRequest) ([]byte, int, error) {
	return u.executeRequest(ctx, "/api/v1/user/tamm_pools", req)
}
