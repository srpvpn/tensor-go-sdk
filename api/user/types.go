package user

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// PortfolioRequest represents the request parameters for getting user portfolio
type PortfolioRequest struct {
	Wallet                string   `json:"wallet"`
	IncludeBidCount       *bool    `json:"includeBidCount,omitempty"`
	IncludeFavouriteCount *bool    `json:"includeFavouriteCount,omitempty"`
	IncludeUnverified     *bool    `json:"includeUnverified,omitempty"`
	IncludeCompressed     *bool    `json:"includeCompressed,omitempty"`
	Currencies            []string `json:"currencies,omitempty"`
}

// PortfolioResponse represents the response from the portfolio API
type PortfolioResponse struct {
	Message     string       `json:"message"`
	Collections []Collection `json:"collections,omitempty"`
}

// Collection represents an NFT collection in the user's portfolio
type Collection struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Symbol     string  `json:"symbol"`
	Image      string  `json:"image"`
	FloorPrice float64 `json:"floorPrice"`
	Volume24h  float64 `json:"volume24h"`
	BidCount   *int    `json:"bidCount,omitempty"`
	FavCount   *int    `json:"favCount,omitempty"`
	Verified   bool    `json:"verified"`
	Compressed bool    `json:"compressed"`
}

// Validate validates the PortfolioRequest fields
func (r *PortfolioRequest) Validate() error {
	if r.Wallet == "" {
		return fmt.Errorf("wallet address is required")
	}

	if err := validateSolanaAddress(r.Wallet); err != nil {
		return fmt.Errorf("invalid wallet address: %w", err)
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

// MarshalJSON implements custom JSON marshaling for PortfolioRequest
func (r *PortfolioRequest) MarshalJSON() ([]byte, error) {
	type Alias PortfolioRequest
	return json.Marshal((*Alias)(r))
}

// UnmarshalJSON implements custom JSON unmarshaling for PortfolioRequest
func (r *PortfolioRequest) UnmarshalJSON(data []byte) error {
	type Alias PortfolioRequest
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Normalize wallet address
	r.Wallet = strings.TrimSpace(r.Wallet)

	return nil
}
