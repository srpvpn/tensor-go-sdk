package user

import "context"

// UserAPI defines the interface for user-related API operations
type UserAPI interface {
	// GetPortfolio retrieves the portfolio collections for a given wallet address
	GetPortfolio(ctx context.Context, req *PortfolioRequest) (*PortfolioResponse, error)
}
