package user

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/srpvpn/tensor-go-sdk/internal/transport"
	"github.com/srpvpn/tensor-go-sdk/internal/utils"
)

// userAPI implements the UserAPI interface
type userAPI struct {
	transport transport.Transport
}

// New creates a new UserAPI instance with the provided transport
func New(transport transport.Transport) UserAPI {
	return &userAPI{
		transport: transport,
	}
}

// GetPortfolio retrieves the portfolio collections for a given wallet address
func (u *userAPI) GetPortfolio(ctx context.Context, req *PortfolioRequest) (*PortfolioResponse, error) {
	// Validate the request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("request validation failed: %w", err)
	}

	// Build query parameters from the request
	params, err := utils.BuildQueryParams(req)
	if err != nil {
		return nil, fmt.Errorf("failed to build query parameters: %w", err)
	}

	// Make the HTTP request
	resp, err := u.transport.Get(ctx, "/api/v1/user/portfolio", params)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the JSON response
	var portfolioResp PortfolioResponse
	if err := json.Unmarshal(body, &portfolioResp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &portfolioResp, nil
}
