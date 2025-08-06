package escrow

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/srpvpn/tensor-go-sdk/internal/transport"
	"github.com/srpvpn/tensor-go-sdk/internal/utils"
)

// escrowAPI implements the EscrowAPI interface
type escrowAPI struct {
	transport transport.Transport
}

// New creates a new Escrow service
func New(t transport.Transport) EscrowAPI {
	return &escrowAPI{
		transport: t,
	}
}

// DepositWithdrawEscrow creates the transaction to deposit or withdraw from an escrow account
// Returns: response, status code, error
func (s *escrowAPI) DepositWithdrawEscrow(ctx context.Context, req *DepositWithdrawEscrowRequest) (*DepositWithdrawEscrowResponse, int, error) {
	body, statusCode, err := s.executeRequest(ctx, "/api/v1/tx/deposit_withdraw_escrow", req)
	if err != nil {
		return nil, statusCode, err
	}

	var response DepositWithdrawEscrowResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, statusCode, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, statusCode, nil
}

// executeRequest is a helper method that handles the common pattern of:
// 1. Request validation
// 2. Query parameter building
// 3. HTTP request execution
// 4. Response handling
func (s *escrowAPI) executeRequest(ctx context.Context, endpoint string, req Validator) ([]byte, int, error) {
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
	resp, err := s.transport.Get(ctx, endpoint, params)
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
