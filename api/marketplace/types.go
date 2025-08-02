package marketplace

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// BuyNFTRequest represents the request parameters for buying an NFT
type BuyNFTRequest struct {
	Buyer                 string  `json:"buyer"`
	Mint                  string  `json:"mint"`
	Owner                 string  `json:"owner"`
	MaxPrice              float64 `json:"maxPrice"`
	Blockhash             string  `json:"blockhash"`
	IncludeTotalCost      *bool   `json:"includeTotalCost,omitempty"`
	Payer                 *string `json:"payer,omitempty"`
	FeePayer              *string `json:"feePayer,omitempty"`
	OptionalRoyaltyPct    *int32  `json:"optionalRoyaltyPct,omitempty"`
	Currency              *string `json:"currency,omitempty"`
	TakerBroker           *string `json:"takerBroker,omitempty"`
	Compute               *int32  `json:"compute,omitempty"`
	PriorityMicroLamports *int32  `json:"priorityMicroLamports,omitempty"`
}

// BuyNFTResponse represents the response from the buy NFT API
type BuyNFTResponse struct {
	Txs []Transaction `json:"txs"`
}

// Transaction represents a transaction in the response
type Transaction struct {
	Tx                   *string                `json:"tx"`
	TxV0                 string                 `json:"txV0"`
	LastValidBlockHeight *float64               `json:"lastValidBlockHeight"`
	Metadata             map[string]interface{} `json:"metadata"`
	TotalCost            *float64               `json:"totalCost,omitempty"`
}

// Validate validates the BuyNFTRequest fields
func (r *BuyNFTRequest) Validate() error {
	if r.Buyer == "" {
		return fmt.Errorf("buyer address is required")
	}

	if err := validateSolanaAddress(r.Buyer); err != nil {
		return fmt.Errorf("invalid buyer address: %w", err)
	}

	if r.Mint == "" {
		return fmt.Errorf("mint address is required")
	}

	if err := validateSolanaAddress(r.Mint); err != nil {
		return fmt.Errorf("invalid mint address: %w", err)
	}

	if r.Owner == "" {
		return fmt.Errorf("owner address is required")
	}

	if err := validateSolanaAddress(r.Owner); err != nil {
		return fmt.Errorf("invalid owner address: %w", err)
	}

	if r.MaxPrice < 0 {
		return fmt.Errorf("maxPrice must be >= 0")
	}

	if r.Blockhash == "" {
		return fmt.Errorf("blockhash is required")
	}

	// Validate optional addresses if provided
	if r.Payer != nil {
		if err := validateSolanaAddress(*r.Payer); err != nil {
			return fmt.Errorf("invalid payer address: %w", err)
		}
	}

	if r.FeePayer != nil {
		if err := validateSolanaAddress(*r.FeePayer); err != nil {
			return fmt.Errorf("invalid feePayer address: %w", err)
		}
	}

	if r.Currency != nil {
		if err := validateSolanaAddress(*r.Currency); err != nil {
			return fmt.Errorf("invalid currency address: %w", err)
		}
	}

	if r.TakerBroker != nil {
		if err := validateSolanaAddress(*r.TakerBroker); err != nil {
			return fmt.Errorf("invalid takerBroker address: %w", err)
		}
	}

	// Validate optional royalty percent
	if r.OptionalRoyaltyPct != nil {
		if *r.OptionalRoyaltyPct < 0 || *r.OptionalRoyaltyPct > 100 {
			return fmt.Errorf("optionalRoyaltyPct must be between 0 and 100")
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

// validateSolanaAddress validates that the provided string is a valid Solana address
func validateSolanaAddress(address string) error {
	// Solana addresses are base58 encoded and typically 32-44 characters long
	if len(address) < 32 || len(address) > 44 {
		return fmt.Errorf("address length must be between 32 and 44 characters")
	}

	// Check for valid base58 characters (no 0, O, I, l)
	validBase58 := regexp.MustCompile(`^[1-9A-HJ-NP-Za-km-z]+$`)
	if !validBase58.MatchString(address) {
		return fmt.Errorf("address contains invalid characters")
	}

	return nil
}

// MarshalJSON implements custom JSON marshaling for BuyNFTRequest
func (r *BuyNFTRequest) MarshalJSON() ([]byte, error) {
	type Alias BuyNFTRequest
	return json.Marshal((*Alias)(r))
}

// UnmarshalJSON implements custom JSON unmarshaling for BuyNFTRequest
func (r *BuyNFTRequest) UnmarshalJSON(data []byte) error {
	type Alias BuyNFTRequest
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Normalize addresses
	r.Buyer = strings.TrimSpace(r.Buyer)
	r.Mint = strings.TrimSpace(r.Mint)
	r.Owner = strings.TrimSpace(r.Owner)
	r.Blockhash = strings.TrimSpace(r.Blockhash)

	return nil
}
