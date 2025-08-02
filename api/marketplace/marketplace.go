package marketplace

import (
	"context"
	"fmt"
	"io"

	"github.com/srpvpn/tensor-go-sdk/internal/transport"
	"github.com/srpvpn/tensor-go-sdk/internal/utils"
)

// marketplaceAPI implements the MarketplaceAPI interface
type marketplaceAPI struct {
	transport transport.Transport
}

// New creates a new MarketplaceAPI instance with the provided transport
func New(transport transport.Transport) MarketplaceAPI {
	return &marketplaceAPI{
		transport: transport,
	}
}

// executeRequest is a helper method that handles the common pattern of:
// 1. Request validation
// 2. Query parameter building
// 3. HTTP request execution
// 4. Response handling
func (m *marketplaceAPI) executeRequest(ctx context.Context, endpoint string, req Validator) ([]byte, int, error) {
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
	resp, err := m.transport.Get(ctx, endpoint, params)
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

