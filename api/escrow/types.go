package escrow

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/srpvpn/tensor-go-sdk/internal/utils"
)

// DepositWithdrawEscrowRequest represents the request for depositing/withdrawing from escrow
type DepositWithdrawEscrowRequest struct {
	Action                string  `json:"action"`                          // The action to perform. Either "deposit" or "withdraw"
	Owner                 string  `json:"owner"`                           // The owner of the Margin account
	Lamports              float64 `json:"lamports"`                        // The amount of SOL to deposit/withdraw
	Blockhash             string  `json:"blockhash"`                       // The blockhash to be passed into the transaction
	Compute               *int32  `json:"compute,omitempty"`               // Compute units for the transaction
	PriorityMicroLamports *int32  `json:"priorityMicroLamports,omitempty"` // The priority in micro-lamports to be used for the transaction
}

// DepositWithdrawEscrowResponse represents the response from the deposit/withdraw escrow endpoint
type DepositWithdrawEscrowResponse struct {
	Status string `json:"status"` // Success status, typically "Ok"
}

// Validator interface for request validation
type Validator interface {
	Validate() error
}

// Validate validates the DepositWithdrawEscrowRequest fields
func (r *DepositWithdrawEscrowRequest) Validate() error {
	if r.Action == "" {
		return fmt.Errorf("action is required")
	}

	// Validate action type (accepts any case, normalized to uppercase for escrow operations)
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

	if r.Owner == "" {
		return fmt.Errorf("owner is required")
	}

	if err := utils.ValidateWalletAddress(r.Owner); err != nil {
		return fmt.Errorf("invalid owner address: %w", err)
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

// MarshalJSON implements custom JSON marshaling for DepositWithdrawEscrowRequest
func (r *DepositWithdrawEscrowRequest) MarshalJSON() ([]byte, error) {
	type Alias DepositWithdrawEscrowRequest
	return json.Marshal((*Alias)(r))
}

// UnmarshalJSON implements custom JSON unmarshaling for DepositWithdrawEscrowRequest
func (r *DepositWithdrawEscrowRequest) UnmarshalJSON(data []byte) error {
	type Alias DepositWithdrawEscrowRequest
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Normalize addresses and strings
	r.Action = strings.ToUpper(strings.TrimSpace(r.Action)) // Escrow operations require uppercase
	r.Owner = strings.TrimSpace(r.Owner)
	r.Blockhash = strings.TrimSpace(r.Blockhash)

	return nil
}

// MarshalJSON implements custom JSON marshaling for DepositWithdrawEscrowResponse
func (r *DepositWithdrawEscrowResponse) MarshalJSON() ([]byte, error) {
	type Alias DepositWithdrawEscrowResponse
	return json.Marshal((*Alias)(r))
}

// UnmarshalJSON implements custom JSON unmarshaling for DepositWithdrawEscrowResponse
func (r *DepositWithdrawEscrowResponse) UnmarshalJSON(data []byte) error {
	type Alias DepositWithdrawEscrowResponse
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	return json.Unmarshal(data, &aux)
}
