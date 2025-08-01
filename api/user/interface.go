package user

import "context"

// UserAPI defines the interface for user-related API operations
type UserAPI interface {
	// GetPortfolio retrieves portfolio data for a given wallet address
	// Returns: response body, status code, error
	GetPortfolio(ctx context.Context, req *PortfolioRequest) ([]byte, int, error)
}
