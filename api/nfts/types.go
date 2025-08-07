package nfts

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/srpvpn/tensor-go-sdk/internal/utils"
)

// NFTsInfoRequest represents the request for getting NFT info
type NFTsInfoRequest struct {
	Mints []string `json:"mints"` // The mint addresses to fetch
}

// NFTsByCollectionRequest represents the request for getting NFTs by collection
type NFTsByCollectionRequest struct {
	CollId            string   `json:"collId"`                      // The collection ID of the mint to filter for
	SortBy            string   `json:"sortBy"`                      // The order in with the NFTs should be returned
	Limit             int32    `json:"limit"`                       // 1 to 250 Number of mint addresses to return
	OnlyListings      *bool    `json:"onlyListings,omitempty"`      // Hide unlisted NFTs
	Mints             []string `json:"mints,omitempty"`             // The list of mints for filter for
	Cursor            *string  `json:"cursor,omitempty"`            // The cursor string received in the previous response
	ListingSources    []string `json:"listingSources,omitempty"`    // Sources to agregate listings from
	MinPrice          *float64 `json:"minPrice,omitempty"`          // The minimum price of to filter for
	MaxPrice          *float64 `json:"maxPrice,omitempty"`          // The maximum price to filter for
	TraitCountMin     *int32   `json:"traitCountMin,omitempty"`     // Minimum number of traits to filter for
	TraitCountMax     *int32   `json:"traitCountMax,omitempty"`     // Maximum number of traits to filter for
	Name              *string  `json:"name,omitempty"`              // Name of the NFT to filter for
	ExcludeOwners     []string `json:"excludeOwners,omitempty"`     // Owners to exclude in results
	IncludeOwners     []string `json:"includeOwners,omitempty"`     // Owners to include in results
	IncludeCurrencies []string `json:"includeCurrencies,omitempty"` // Currencies to include in results
	Traits            []string `json:"traits,omitempty"`            // Traits and values to filter for
	RaritySystem      *string  `json:"raritySystem,omitempty"`      // Rarity System to use when filtering for rarity
	RarityMin         *float64 `json:"rarityMin,omitempty"`         // Minimum rarity points to return in results
	RarityMax         *float64 `json:"rarityMax,omitempty"`         // Maximum rarity points to return in results
	OnlyInscriptions  *bool    `json:"onlyInscriptions,omitempty"`  // Filter to include only Solana Inscriptions
	ImmutableStatus   *string  `json:"immutableStatus,omitempty"`   // Filter the immutability of the Inscriptions
}

// Validator interface for request validation
type Validator interface {
	Validate() error
}

// Validate validates the NFTsInfoRequest fields
func (r *NFTsInfoRequest) Validate() error {
	if len(r.Mints) == 0 {
		return fmt.Errorf("mints is required and cannot be empty")
	}

	// Validate each mint address
	for i, mint := range r.Mints {
		if mint == "" {
			return fmt.Errorf("mint address at index %d cannot be empty", i)
		}
		if err := utils.ValidateWalletAddress(mint); err != nil {
			return fmt.Errorf("invalid mint address at index %d: %w", i, err)
		}
	}

	return nil
}

// Validate validates the NFTsByCollectionRequest fields
func (r *NFTsByCollectionRequest) Validate() error {
	if r.CollId == "" {
		return fmt.Errorf("collId is required")
	}

	if r.SortBy == "" {
		return fmt.Errorf("sortBy is required")
	}

	if r.Limit < 1 || r.Limit > 250 {
		return fmt.Errorf("limit must be between 1 and 250")
	}

	// Validate optional mint addresses if provided
	for i, mint := range r.Mints {
		if mint == "" {
			return fmt.Errorf("mint address at index %d cannot be empty", i)
		}
		if err := utils.ValidateWalletAddress(mint); err != nil {
			return fmt.Errorf("invalid mint address at index %d: %w", i, err)
		}
	}

	// Validate optional owner addresses if provided
	for i, owner := range r.ExcludeOwners {
		if owner == "" {
			return fmt.Errorf("exclude owner address at index %d cannot be empty", i)
		}
		if err := utils.ValidateWalletAddress(owner); err != nil {
			return fmt.Errorf("invalid exclude owner address at index %d: %w", i, err)
		}
	}

	for i, owner := range r.IncludeOwners {
		if owner == "" {
			return fmt.Errorf("include owner address at index %d cannot be empty", i)
		}
		if err := utils.ValidateWalletAddress(owner); err != nil {
			return fmt.Errorf("invalid include owner address at index %d: %w", i, err)
		}
	}

	// Validate price ranges
	if r.MinPrice != nil && *r.MinPrice < 0 {
		return fmt.Errorf("minPrice must be >= 0")
	}

	if r.MaxPrice != nil && *r.MaxPrice < 0 {
		return fmt.Errorf("maxPrice must be >= 0")
	}

	// Validate trait count ranges
	if r.TraitCountMin != nil && *r.TraitCountMin < 0 {
		return fmt.Errorf("traitCountMin must be >= 0")
	}

	if r.TraitCountMax != nil && *r.TraitCountMax < 1 {
		return fmt.Errorf("traitCountMax must be >= 1")
	}

	// Validate rarity ranges
	if r.RarityMin != nil && *r.RarityMin < 0 {
		return fmt.Errorf("rarityMin must be >= 0")
	}

	if r.RarityMax != nil && *r.RarityMax < 0 {
		return fmt.Errorf("rarityMax must be >= 0")
	}

	return nil
}

// MarshalJSON implements custom JSON marshaling for NFTsInfoRequest
func (r *NFTsInfoRequest) MarshalJSON() ([]byte, error) {
	type Alias NFTsInfoRequest
	return json.Marshal((*Alias)(r))
}

// UnmarshalJSON implements custom JSON unmarshaling for NFTsInfoRequest
func (r *NFTsInfoRequest) UnmarshalJSON(data []byte) error {
	type Alias NFTsInfoRequest
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Normalize mint addresses
	for i, mint := range r.Mints {
		r.Mints[i] = strings.TrimSpace(mint)
	}

	return nil
}

// MarshalJSON implements custom JSON marshaling for NFTsByCollectionRequest
func (r *NFTsByCollectionRequest) MarshalJSON() ([]byte, error) {
	type Alias NFTsByCollectionRequest
	return json.Marshal((*Alias)(r))
}

// UnmarshalJSON implements custom JSON unmarshaling for NFTsByCollectionRequest
func (r *NFTsByCollectionRequest) UnmarshalJSON(data []byte) error {
	type Alias NFTsByCollectionRequest
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Normalize addresses and strings
	r.CollId = strings.TrimSpace(r.CollId)
	r.SortBy = strings.TrimSpace(r.SortBy)

	// Normalize mint addresses
	for i, mint := range r.Mints {
		r.Mints[i] = strings.TrimSpace(mint)
	}

	// Normalize owner addresses
	for i, owner := range r.ExcludeOwners {
		r.ExcludeOwners[i] = strings.TrimSpace(owner)
	}

	for i, owner := range r.IncludeOwners {
		r.IncludeOwners[i] = strings.TrimSpace(owner)
	}

	// Normalize optional string fields
	if r.Cursor != nil {
		trimmed := strings.TrimSpace(*r.Cursor)
		r.Cursor = &trimmed
	}

	if r.Name != nil {
		trimmed := strings.TrimSpace(*r.Name)
		r.Name = &trimmed
	}

	if r.RaritySystem != nil {
		trimmed := strings.TrimSpace(*r.RaritySystem)
		r.RaritySystem = &trimmed
	}

	if r.ImmutableStatus != nil {
		trimmed := strings.TrimSpace(*r.ImmutableStatus)
		r.ImmutableStatus = &trimmed
	}

	return nil
}
