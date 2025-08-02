package user

import (
	"context"
)

// GetListings retrieves all active listings for supplied wallets
// Returns: response body, status code, error
func (u *userAPI) GetListings(ctx context.Context, req *ListingsRequest) ([]byte, int, error) {
	return u.executeRequest(ctx, "/api/v1/user/active_listings", req)
}
