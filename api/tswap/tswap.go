package tswap

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/srpvpn/tensor-go-sdk/internal/transport"
	"github.com/srpvpn/tensor-go-sdk/internal/utils"
)

// tswapAPI implements the TSwapAPI interface
type tswapAPI struct {
	transport transport.Transport
}

// New creates a new TSwap service
func New(t transport.Transport) TSwapAPI {
	return &tswapAPI{
		transport: t,
	}
}

// CloseTSwapPool creates the transaction to close a TSwap pool
// Returns: response, status code, error
func (s *tswapAPI) CloseTSwapPool(ctx context.Context, req *CloseTSwapPoolRequest) (*CloseTSwapPoolResponse, int, error) {
	body, statusCode, err := s.executeRequest(ctx, "/api/v1/tx/tswap/close_order", req)
	if err != nil {
		return nil, statusCode, err
	}

	var response CloseTSwapPoolResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, statusCode, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, statusCode, nil
}

// EditTSwapPool creates the transaction to edit a TSwap pool
// Returns: response, status code, error
func (s *tswapAPI) EditTSwapPool(ctx context.Context, req *EditTSwapPoolRequest) (*EditTSwapPoolResponse, int, error) {
	body, statusCode, err := s.executeRequest(ctx, "/api/v1/tx/tswap/edit_order", req)
	if err != nil {
		return nil, statusCode, err
	}

	var response EditTSwapPoolResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, statusCode, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, statusCode, nil
}

// DepositWithdrawNFT creates the transaction to deposit/withdraw NFT to/from a TSwap pool
// Returns: response, status code, error
func (s *tswapAPI) DepositWithdrawNFT(ctx context.Context, req *DepositWithdrawNFTRequest) (*DepositWithdrawNFTResponse, int, error) {
	body, statusCode, err := s.executeRequest(ctx, "/api/v1/tx/tswap/deposit_withdraw", req)
	if err != nil {
		return nil, statusCode, err
	}

	var response DepositWithdrawNFTResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, statusCode, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, statusCode, nil
}

// DepositWithdrawSOL creates the transaction to deposit/withdraw SOL to/from a TSwap pool
// Returns: response, status code, error
func (s *tswapAPI) DepositWithdrawSOL(ctx context.Context, req *DepositWithdrawSOLRequest) (*DepositWithdrawSOLResponse, int, error) {
	body, statusCode, err := s.executeRequest(ctx, "/api/v1/tx/tswap/deposit_withdraw_sol", req)
	if err != nil {
		return nil, statusCode, err
	}

	var response DepositWithdrawSOLResponse
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
func (s *tswapAPI) executeRequest(ctx context.Context, endpoint string, req Validator) ([]byte, int, error) {
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
