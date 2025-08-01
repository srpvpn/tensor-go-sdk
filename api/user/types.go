package user

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// ListingsRequest represents the request parameters for getting user listings
type ListingsRequest struct {
	Wallets    []string `json:"wallets"`
	SortBy     string   `json:"sortBy"`
	Limit      int32    `json:"limit"`
	Cursor     *string  `json:"cursor,omitempty"`
	CollId     *string  `json:"collId,omitempty"`
	Currencies []string `json:"currencies,omitempty"`
}

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
var solanaAddressRegex = regexp.MustCompile(`^[1-9A-HJ-NP-Za-km-z]+$`)

func validateSolanaAddress(address string) error {
	if len(address) < 32 || len(address) > 44 {
		return fmt.Errorf("address length must be between 32 and 44 characters")
	}

	if !solanaAddressRegex.MatchString(address) {
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

// Validate validates the ListingsRequest fields
func (r *ListingsRequest) Validate() error {
	if len(r.Wallets) == 0 {
		return fmt.Errorf("at least one wallet address is required")
	}

	for _, wallet := range r.Wallets {
		if err := validateSolanaAddress(wallet); err != nil {
			return fmt.Errorf("invalid wallet address %s: %w", wallet, err)
		}
	}

	if r.Limit <= 0 {
		return fmt.Errorf("limit must be greater than 0")
	}

	validSortOptions := []string{
		"PriceDesc", "NormalizedPriceAsc", "NormalizedPriceDesc", "HybridAmountAsc",
		"HybridAmountDesc", "LastSaleAsc", "LastSaleDesc", "ListedDesc", "RankHrttAsc",
		"RankHrttDesc", "RankStatAsc", "OrdinalAsc", "RankStatDesc", "OrdinalDesc",
		"RankTeamAsc", "RankTeamDesc", "RankTnAsc", "RankTnDesc", "PriceAsc",
	}

	if r.SortBy != "" {
		isValid := false
		for _, validOption := range validSortOptions {
			if r.SortBy == validOption {
				isValid = true
				break
			}
		}
		if !isValid {
			return fmt.Errorf("invalid sortBy value: %s", r.SortBy)
		}
	}

	return nil
}
