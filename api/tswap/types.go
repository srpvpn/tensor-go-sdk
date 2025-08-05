package tswap

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/srpvpn/tensor-go-sdk/internal/utils"
)

// CloseTSwapPoolRequest represents the request parameters for closing a TSwap pool
type CloseTSwapPoolRequest struct {
	PoolAddress           string `json:"poolAddress"`
	Blockhash             string `json:"blockhash"`
	Compute               *int32 `json:"compute,omitempty"`
	PriorityMicroLamports *int32 `json:"priorityMicroLamports,omitempty"`
}

// CloseTSwapPoolResponse represents the response from the close TSwap pool API
type CloseTSwapPoolResponse struct {
	Txs []Transaction `json:"txs"`
}

// EditTSwapPoolRequest represents the request parameters for editing a TSwap pool
type EditTSwapPoolRequest struct {
	PoolAddress           string   `json:"poolAddress"`
	PoolType              string   `json:"poolType"`
	CurveType             string   `json:"curveType"`
	StartingPrice         float64  `json:"startingPrice"`
	Delta                 float64  `json:"delta"`
	Blockhash             string   `json:"blockhash"`
	MmKeepFeesSeparate    *bool    `json:"mmKeepFeesSeparate,omitempty"`
	MmFeeBps              *float64 `json:"mmFeeBps,omitempty"`
	MaxTakerSellCount     *int32   `json:"maxTakerSellCount,omitempty"`
	UseSharedEscrow       *bool    `json:"useSharedEscrow,omitempty"`
	Compute               *int32   `json:"compute,omitempty"`
	PriorityMicroLamports *int32   `json:"priorityMicroLamports,omitempty"`
}

// EditTSwapPoolResponse represents the response from the edit TSwap pool API
type EditTSwapPoolResponse struct {
	Txs []Transaction `json:"txs"`
}

// DepositWithdrawNFTRequest represents the request parameters for depositing/withdrawing NFT to/from a TSwap pool
type DepositWithdrawNFTRequest struct {
	Action                string  `json:"action"`
	PoolAddress           string  `json:"poolAddress"`
	Mint                  string  `json:"mint"`
	Blockhash             string  `json:"blockhash"`
	Compute               *int32  `json:"compute,omitempty"`
	PriorityMicroLamports *int32  `json:"priorityMicroLamports,omitempty"`
	NftSource             *string `json:"nftSource,omitempty"`
}

// DepositWithdrawNFTResponse represents the response from the deposit/withdraw NFT API
type DepositWithdrawNFTResponse struct {
	Txs []Transaction `json:"txs"`
}

// DepositWithdrawSOLRequest represents the request parameters for depositing/withdrawing SOL to/from a TSwap pool
type DepositWithdrawSOLRequest struct {
	Action                string  `json:"action"`
	PoolAddress           string  `json:"poolAddress"`
	Lamports              float64 `json:"lamports"`
	Blockhash             string  `json:"blockhash"`
	Compute               *int32  `json:"compute,omitempty"`
	PriorityMicroLamports *int32  `json:"priorityMicroLamports,omitempty"`
}

// DepositWithdrawSOLResponse represents the response from the deposit/withdraw SOL API
type DepositWithdrawSOLResponse struct {
	Txs []Transaction `json:"txs"`
}

// Transaction represents a transaction in the response
type Transaction struct {
	Tx                   *string                `json:"tx"`
	TxV0                 string                 `json:"txV0"`
	LastValidBlockHeight *float64               `json:"lastValidBlockHeight"`
	Metadata             map[string]interface{} `json:"metadata"`
}

// Validate validates the CloseTSwapPoolRequest fields
func (r *CloseTSwapPoolRequest) Validate() error {
	if r.PoolAddress == "" {
		return fmt.Errorf("poolAddress is required")
	}

	if err := utils.ValidateWalletAddress(r.PoolAddress); err != nil {
		return fmt.Errorf("invalid poolAddress: %w", err)
	}

	if r.Blockhash == "" {
		return fmt.Errorf("blockhash is required")
	}

	// Validate compute units
	if r.Compute != nil && *r.Compute < 0 {
		return fmt.Errorf("compute must be >= 0")
	}

	// Validate priority micro lamports
	if r.PriorityMicroLamports != nil && *r.PriorityMicroLamports < 0 {
		return fmt.Errorf("priorityMicroLamports must be >= 0")
	}

	return nil
}

// MarshalJSON implements custom JSON marshaling for CloseTSwapPoolRequest
func (r *CloseTSwapPoolRequest) MarshalJSON() ([]byte, error) {
	type Alias CloseTSwapPoolRequest
	return json.Marshal((*Alias)(r))
}

// UnmarshalJSON implements custom JSON unmarshaling for CloseTSwapPoolRequest
func (r *CloseTSwapPoolRequest) UnmarshalJSON(data []byte) error {
	type Alias CloseTSwapPoolRequest
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Normalize addresses
	r.PoolAddress = strings.TrimSpace(r.PoolAddress)
	r.Blockhash = strings.TrimSpace(r.Blockhash)

	return nil
}

// Validate validates the EditTSwapPoolRequest fields
func (r *EditTSwapPoolRequest) Validate() error {
	if r.PoolAddress == "" {
		return fmt.Errorf("poolAddress is required")
	}

	if err := utils.ValidateWalletAddress(r.PoolAddress); err != nil {
		return fmt.Errorf("invalid poolAddress: %w", err)
	}

	if r.PoolType == "" {
		return fmt.Errorf("poolType is required")
	}

	// Validate pool type
	validPoolTypes := []string{"TOKEN", "NFT", "TRADE"}
	isValidPoolType := false
	for _, validType := range validPoolTypes {
		if r.PoolType == validType {
			isValidPoolType = true
			break
		}
	}
	if !isValidPoolType {
		return fmt.Errorf("invalid poolType: %s, must be one of: %v", r.PoolType, validPoolTypes)
	}

	if r.CurveType == "" {
		return fmt.Errorf("curveType is required")
	}

	// Validate curve type
	validCurveTypes := []string{"linear", "exponential"}
	isValidCurveType := false
	for _, validType := range validCurveTypes {
		if r.CurveType == validType {
			isValidCurveType = true
			break
		}
	}
	if !isValidCurveType {
		return fmt.Errorf("invalid curveType: %s, must be one of: %v", r.CurveType, validCurveTypes)
	}

	if r.StartingPrice < 0 {
		return fmt.Errorf("startingPrice must be >= 0")
	}

	if r.Delta < 0 {
		return fmt.Errorf("delta must be >= 0")
	}

	if r.Blockhash == "" {
		return fmt.Errorf("blockhash is required")
	}

	// Validate optional fields
	if r.MmFeeBps != nil && (*r.MmFeeBps < 0 || *r.MmFeeBps > 10000) {
		return fmt.Errorf("mmFeeBps must be between 0 and 10000 basis points")
	}

	if r.MaxTakerSellCount != nil && *r.MaxTakerSellCount < 0 {
		return fmt.Errorf("maxTakerSellCount must be >= 0")
	}

	// Validate compute units
	if r.Compute != nil && *r.Compute < 0 {
		return fmt.Errorf("compute must be >= 0")
	}

	// Validate priority micro lamports
	if r.PriorityMicroLamports != nil && *r.PriorityMicroLamports < 0 {
		return fmt.Errorf("priorityMicroLamports must be >= 0")
	}

	return nil
}

// MarshalJSON implements custom JSON marshaling for EditTSwapPoolRequest
func (r *EditTSwapPoolRequest) MarshalJSON() ([]byte, error) {
	type Alias EditTSwapPoolRequest
	return json.Marshal((*Alias)(r))
}

// UnmarshalJSON implements custom JSON unmarshaling for EditTSwapPoolRequest
func (r *EditTSwapPoolRequest) UnmarshalJSON(data []byte) error {
	type Alias EditTSwapPoolRequest
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Normalize addresses and strings
	r.PoolAddress = strings.TrimSpace(r.PoolAddress)
	r.PoolType = strings.TrimSpace(r.PoolType)
	r.CurveType = strings.TrimSpace(r.CurveType)
	r.Blockhash = strings.TrimSpace(r.Blockhash)

	return nil
}

// Validate validates the DepositWithdrawNFTRequest fields
func (r *DepositWithdrawNFTRequest) Validate() error {
	if r.Action == "" {
		return fmt.Errorf("action is required")
	}

	// Validate action type (accepts any case, normalized to uppercase for NFT operations)
	actionUpper := strings.ToUpper(r.Action)
	validActions := []string{"DEPOSIT", "WITHDRAW"}
	isValidAction := false
	for _, validAction := range validActions {
		if actionUpper == validAction {
			isValidAction = true
			break
		}
	}
	if !isValidAction {
		return fmt.Errorf("invalid action: %s, must be 'deposit' or 'withdraw' (case insensitive)", r.Action)
	}

	if r.PoolAddress == "" {
		return fmt.Errorf("poolAddress is required")
	}

	if err := utils.ValidateWalletAddress(r.PoolAddress); err != nil {
		return fmt.Errorf("invalid poolAddress: %w", err)
	}

	if r.Mint == "" {
		return fmt.Errorf("mint is required")
	}

	if err := utils.ValidateWalletAddress(r.Mint); err != nil {
		return fmt.Errorf("invalid mint address: %w", err)
	}

	if r.Blockhash == "" {
		return fmt.Errorf("blockhash is required")
	}

	// Validate optional NFT source if provided
	if r.NftSource != nil && *r.NftSource != "" {
		if err := utils.ValidateWalletAddress(*r.NftSource); err != nil {
			return fmt.Errorf("invalid nftSource address: %w", err)
		}
	}

	// Validate compute units
	if r.Compute != nil && *r.Compute < 0 {
		return fmt.Errorf("compute must be >= 0")
	}

	// Validate priority micro lamports
	if r.PriorityMicroLamports != nil && *r.PriorityMicroLamports < 0 {
		return fmt.Errorf("priorityMicroLamports must be >= 0")
	}

	return nil
}

// MarshalJSON implements custom JSON marshaling for DepositWithdrawNFTRequest
func (r *DepositWithdrawNFTRequest) MarshalJSON() ([]byte, error) {
	type Alias DepositWithdrawNFTRequest
	return json.Marshal((*Alias)(r))
}

// UnmarshalJSON implements custom JSON unmarshaling for DepositWithdrawNFTRequest
func (r *DepositWithdrawNFTRequest) UnmarshalJSON(data []byte) error {
	type Alias DepositWithdrawNFTRequest
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Normalize addresses and strings
	r.Action = strings.ToUpper(strings.TrimSpace(r.Action)) // NFT operations require uppercase
	r.PoolAddress = strings.TrimSpace(r.PoolAddress)
	r.Mint = strings.TrimSpace(r.Mint)
	r.Blockhash = strings.TrimSpace(r.Blockhash)

	if r.NftSource != nil {
		trimmed := strings.TrimSpace(*r.NftSource)
		r.NftSource = &trimmed
	}

	return nil
}

// Validate validates the DepositWithdrawSOLRequest fields
func (r *DepositWithdrawSOLRequest) Validate() error {
	if r.Action == "" {
		return fmt.Errorf("action is required")
	}

	// Validate action type (accepts any case, normalized to uppercase for SOL operations)
	actionUpper := strings.ToUpper(r.Action)
	validActions := []string{"DEPOSIT", "WITHDRAW"}
	isValidAction := false
	for _, validAction := range validActions {
		if actionUpper == validAction {
			isValidAction = true
			break
		}
	}
	if !isValidAction {
		return fmt.Errorf("invalid action: %s, must be 'deposit' or 'withdraw' (case insensitive)", r.Action)
	}

	if r.PoolAddress == "" {
		return fmt.Errorf("poolAddress is required")
	}

	if err := utils.ValidateWalletAddress(r.PoolAddress); err != nil {
		return fmt.Errorf("invalid poolAddress: %w", err)
	}

	if r.Lamports < 0 {
		return fmt.Errorf("lamports must be >= 0")
	}

	if r.Blockhash == "" {
		return fmt.Errorf("blockhash is required")
	}

	// Validate compute units
	if r.Compute != nil && *r.Compute < 0 {
		return fmt.Errorf("compute must be >= 0")
	}

	// Validate priority micro lamports
	if r.PriorityMicroLamports != nil && *r.PriorityMicroLamports < 0 {
		return fmt.Errorf("priorityMicroLamports must be >= 0")
	}

	return nil
}

// MarshalJSON implements custom JSON marshaling for DepositWithdrawSOLRequest
func (r *DepositWithdrawSOLRequest) MarshalJSON() ([]byte, error) {
	type Alias DepositWithdrawSOLRequest
	return json.Marshal((*Alias)(r))
}

// UnmarshalJSON implements custom JSON unmarshaling for DepositWithdrawSOLRequest
func (r *DepositWithdrawSOLRequest) UnmarshalJSON(data []byte) error {
	type Alias DepositWithdrawSOLRequest
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Normalize addresses and strings
	r.Action = strings.ToUpper(strings.TrimSpace(r.Action)) // SOL operations require uppercase
	r.PoolAddress = strings.TrimSpace(r.PoolAddress)
	r.Blockhash = strings.TrimSpace(r.Blockhash)

	return nil
}
