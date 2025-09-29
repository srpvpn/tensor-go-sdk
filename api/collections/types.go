package collections

import (
	"fmt"
	"strings"
)

// GetVerifiedCollectionsRequest represents the request parameters for getting verified collections
type GetVerifiedCollectionsRequest struct {
	SortBy       string   `json:"sortBy"`                 // required - The order in which collections are returned
	Limit        int32    `json:"limit"`                  // required - The number of collections returned (1 to 100)
	SlugDisplays []string `json:"slugDisplays,omitempty"` // Slugs used in tensor.trade/trade/ urls
	CollIds      []string `json:"collIds,omitempty"`      // Collection IDs to filter by
	Vocs         []string `json:"vocs,omitempty"`         // Verified on-chain collection mints (max 10)
	Fvcs         []string `json:"fvcs,omitempty"`         // First verified creators (max 10)
	Page         *int32   `json:"page,omitempty"`         // The page number of the response (≥ 1)
}

// GetVerifiedCollectionsResponse represents the response from the verified collections API
type GetVerifiedCollectionsResponse struct {
	Page        int32                `json:"page"`
	Total       int32                `json:"total"`
	Collections []CollectionDetailed `json:"collections"`
}

// CollectionDetailed represents a detailed collection with all its metadata and stats
type CollectionDetailed struct {
	Name                string          `json:"name"`
	CollId              string          `json:"collId"`
	SlugDisplay         string          `json:"slugDisplay"`
	SlugMe              string          `json:"slugMe"`
	Symbol              string          `json:"symbol"`
	Description         string          `json:"description"`
	TeamId              string          `json:"teamId"`
	Website             string          `json:"website,omitempty"`
	Discord             string          `json:"discord,omitempty"`
	Twitter             string          `json:"twitter,omitempty"`
	ImageUri            string          `json:"imageUri"`
	TensorVerified      bool            `json:"tensorVerified"`
	TensorWhitelisted   bool            `json:"tensorWhitelisted"`
	WhitelistV2Pda      []interface{}   `json:"whitelistV2Pda,omitempty"`
	TokenStandard       string          `json:"tokenStandard"`
	Compressed          bool            `json:"compressed"`
	Inscription         bool            `json:"inscription"`
	InscriptionMetaplex bool            `json:"inscriptionMetaplex"`
	Spl20               bool            `json:"spl20"`
	SellRoyaltyFeeBPS   int32           `json:"sellRoyaltyFeeBPS"`
	Stats               CollectionStats `json:"stats"`
	CreatedAt           string          `json:"createdAt"`
	UpdatedAt           string          `json:"updatedAt"`
	FirstListDate       string          `json:"firstListDate,omitempty"`
	Hidden              bool            `json:"hidden"`
	TokenProgram        string          `json:"tokenProgram"`
}

// CollectionStats represents statistics for a collection
type CollectionStats struct {
	BuyNowPrice         string  `json:"buyNowPrice"`
	BuyNowPriceNetFees  string  `json:"buyNowPriceNetFees"`
	Floor1h             float64 `json:"floor1h"`
	Floor24h            float64 `json:"floor24h"`
	Floor7d             float64 `json:"floor7d"`
	MarketCap           string  `json:"marketCap"`
	NumBids             int32   `json:"numBids"`
	NumListed           int32   `json:"numListed"`
	NumListed1h         float64 `json:"numListed1h"`
	NumListed24h        float64 `json:"numListed24h"`
	NumListed7d         float64 `json:"numListed7d"`
	NumMints            int32   `json:"numMints"`
	PctListed           float64 `json:"pctListed"`
	Sales1h             int32   `json:"sales1h"`
	Sales24h            int32   `json:"sales24h"`
	Sales7d             int32   `json:"sales7d"`
	SalesAll            int32   `json:"salesAll"`
	SellNowPrice        string  `json:"sellNowPrice"`
	SellNowPriceNetFees string  `json:"sellNowPriceNetFees"`
	Volume1h            string  `json:"volume1h"`
	Volume24h           string  `json:"volume24h"`
	Volume7d            string  `json:"volume7d"`
	VolumeAll           string  `json:"volumeAll"`
}

// Validate validates the GetVerifiedCollectionsRequest fields
func (r *GetVerifiedCollectionsRequest) Validate() error {
	if r.SortBy == "" {
		return fmt.Errorf("sortBy is required")
	}

	// Validate sortBy format - should contain a colon for direction (e.g., "statsV2.volume1h:desc")
	if !strings.Contains(r.SortBy, ":") {
		return fmt.Errorf("sortBy must include direction (e.g., 'statsV2.volume1h:desc')")
	}

	if r.Limit <= 0 {
		return fmt.Errorf("limit must be greater than 0")
	}

	if r.Limit > 100 {
		return fmt.Errorf("limit must be 100 or less")
	}

	// Validate max vocs/fvcs (max 10)
	if len(r.Vocs) > 10 {
		return fmt.Errorf("maximum 10 vocs allowed")
	}

	if len(r.Fvcs) > 10 {
		return fmt.Errorf("maximum 10 fvcs allowed")
	}

	// Validate page number if provided
	if r.Page != nil && *r.Page < 1 {
		return fmt.Errorf("page must be 1 or greater")
	}

	return nil
}
