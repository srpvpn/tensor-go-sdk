package collections

import "context"

// Validator defines the interface for request validation
type Validator interface {
	Validate() error
}

// CollectionsAPI defines the interface for collections-related API operations
type CollectionsAPI interface {
	// GetVerifiedCollections retrieves all verified collections based on parameters provided
	// Returns: response body, status code, error
	GetVerifiedCollections(ctx context.Context, req *GetVerifiedCollectionsRequest) ([]byte, int, error)
}
