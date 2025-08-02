package user

import (
	"context"
)

// GetInventoryForCollection retrieves all active listings for supplied wallets
// Returns: response body, status code, error
func (u *userAPI) GetInventoryForCollection(ctx context.Context, req *InventoryForCollectionRequest) ([]byte, int, error) {
	return u.executeRequest(ctx, "/api/v1/user/inventory_by_collection", req)
}
