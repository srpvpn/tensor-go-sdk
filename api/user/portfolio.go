package user

import (
	"context"
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

// GetPortfolio retrieves portfolio data for a given wallet address
// Returns: response body, status code, error
func (u *userAPI) GetPortfolio(ctx context.Context, req *PortfolioRequest) ([]byte, int, error) {
	// Validate the request
	if err := req.Validate(); err != nil {
		return nil, 0, fmt.Errorf("request validation failed: %w", err)
	}

	// Build query parameters from the request
	params, err := utils.BuildQueryParams(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to build query parameters: %w", err)
	}

	// Make the HTTP request
	resp, err := u.transport.Get(ctx, "/api/v1/user/portfolio", params)
	if err != nil {
		return nil, 0, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode >= 400 {
		return body, resp.StatusCode, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	return body, resp.StatusCode, nil
}
