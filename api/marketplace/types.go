package marketplace

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// EditListingRequest represents the request parameters for editing an NFT listing
type EditListingRequest struct {
	Mint                  string  `json:"mint"`
	Owner                 string  `json:"owner"`
	Price                 float64 `json:"price"`
	Blockhash             string  `json:"blockhash"`
	MakerBroker           *string `json:"makerBroker,omitempty"`
	ExpireIn              *int32  `json:"expireIn,omitempty"`
	FeePayer              *string `json:"feePayer,omitempty"`
	Compute               *int32  `json:"compute,omitempty"`
	PriorityMicroLamports *int32  `json:"priorityMicroLamports,omitempty"`
}

// EditListingResponse represents the response from the edit listing API
type EditListingResponse struct {
	Txs []Transaction `json:"txs"`
}

// EditBidRequest represents the request parameters for editing a bid
type EditBidRequest struct {
	BidStateAddress       string   `json:"bidStateAddress"`
	Blockhash             string   `json:"blockhash"`
	Price                 *float64 `json:"price,omitempty"`
	Quantity              *int32   `json:"quantity,omitempty"`
	ExpireIn              *int32   `json:"expireIn,omitempty"`
	PrivateTaker          *string  `json:"privateTaker,omitempty"`
	UseSharedEscrow       *bool    `json:"useSharedEscrow,omitempty"`
	Compute               *int32   `json:"compute,omitempty"`
	PriorityMicroLamports *int32   `json:"priorityMicroLamports,omitempty"`
}

// EditBidResponse represents the response from the edit bid API
type EditBidResponse struct {
	Txs      []Transaction `json:"txs"`
	BidState string        `json:"bidState"`
}

// CancelBidRequest represents the request parameters for canceling a bid
type CancelBidRequest struct {
	BidStateAddress       string `json:"bidStateAddress"`
	Blockhash             string `json:"blockhash"`
	Compute               *int32 `json:"compute,omitempty"`
	PriorityMicroLamports *int32 `json:"priorityMicroLamports,omitempty"`
}

// CancelBidResponse represents the response from the cancel bid API
type CancelBidResponse struct {
	Txs      []Transaction `json:"txs"`
	BidState string        `json:"bidState"`
}

// PlaceNFTBidRequest represents the request parameters for placing a bid on a single NFT
type PlaceNFTBidRequest struct {
	Owner                 string  `json:"owner"`
	Price                 float64 `json:"price"`
	Mint                  string  `json:"mint"`
	Blockhash             string  `json:"blockhash"`
	MakerBroker           *string `json:"makerBroker,omitempty"`
	UseSharedEscrow       *bool   `json:"useSharedEscrow,omitempty"`
	RentPayer             *string `json:"rentPayer,omitempty"`
	ExpireIn              *int32  `json:"expireIn,omitempty"`
	Compute               *int32  `json:"compute,omitempty"`
	PriorityMicroLamports *int32  `json:"priorityMicroLamports,omitempty"`
}

// PlaceNFTBidResponse represents the response from the place NFT bid API
type PlaceNFTBidResponse struct {
	Message string `json:"message"`
}

// PlaceTraitBidRequest represents the request parameters for placing a trait bid on a collection
type PlaceTraitBidRequest struct {
	Owner                 string   `json:"owner"`
	Price                 float64  `json:"price"`
	Quantity              int32    `json:"quantity"`
	CollId                string   `json:"collId"`
	Blockhash             string   `json:"blockhash"`
	MakerBroker           *string  `json:"makerBroker,omitempty"`
	Traits                []string `json:"traits,omitempty"`
	UseSharedEscrow       *bool    `json:"useSharedEscrow,omitempty"`
	RentPayer             *string  `json:"rentPayer,omitempty"`
	ExpireIn              *int32   `json:"expireIn,omitempty"`
	Compute               *int32   `json:"compute,omitempty"`
	PriorityMicroLamports *int32   `json:"priorityMicroLamports,omitempty"`
}

// PlaceTraitBidResponse represents the response from the place trait bid API
type PlaceTraitBidResponse struct {
	Message string `json:"message"`
}

// PlaceCollectionBidRequest represents the request parameters for placing a collection wide bid
type PlaceCollectionBidRequest struct {
	Owner                 string   `json:"owner"`
	Price                 float64  `json:"price"`
	Quantity              int32    `json:"quantity"`
	CollId                string   `json:"collId"`
	Blockhash             string   `json:"blockhash"`
	MakerBroker           *string  `json:"makerBroker,omitempty"`
	UseSharedEscrow       *bool    `json:"useSharedEscrow,omitempty"`
	RentPayer             *string  `json:"rentPayer,omitempty"`
	ExpireIn              *int32   `json:"expireIn,omitempty"`
	TopUp                 *float64 `json:"topUp,omitempty"`
	Compute               *int32   `json:"compute,omitempty"`
	PriorityMicroLamports *int32   `json:"priorityMicroLamports,omitempty"`
}

// PlaceCollectionBidResponse represents the response from the place collection bid API
type PlaceCollectionBidResponse struct {
	Message string `json:"message"`
}

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

// SellNFTRequest represents the request parameters for selling an NFT (accepting a bid)
type SellNFTRequest struct {
	Seller                string  `json:"seller"`
	Mint                  string  `json:"mint"`
	BidAddress            string  `json:"bidAddress"`
	MinPrice              float64 `json:"minPrice"`
	Blockhash             string  `json:"blockhash"`
	TakerBroker           *string `json:"takerBroker,omitempty"`
	FeePayer              *string `json:"feePayer,omitempty"`
	OptionalRoyaltyPct    *int32  `json:"optionalRoyaltyPct,omitempty"`
	Currency              *string `json:"currency,omitempty"`
	DelegateSigner        *bool   `json:"delegateSigner,omitempty"`
	IncludeProof          *bool   `json:"includeProof,omitempty"`
	Compute               *int32  `json:"compute,omitempty"`
	PriorityMicroLamports *int32  `json:"priorityMicroLamports,omitempty"`
}

// SellNFTResponse represents the response from the sell NFT API
type SellNFTResponse struct {
	Txs []Transaction `json:"txs"`
}

// ListNFTRequest represents the request parameters for listing an NFT
type ListNFTRequest struct {
	Mint                  string  `json:"mint"`
	Owner                 string  `json:"owner"`
	Price                 float64 `json:"price"`
	Blockhash             string  `json:"blockhash"`
	MakerBroker           *string `json:"makerBroker,omitempty"`
	Payer                 *string `json:"payer,omitempty"`
	FeePayer              *string `json:"feePayer,omitempty"`
	RentPayer             *string `json:"rentPayer,omitempty"`
	Currency              *string `json:"currency,omitempty"`
	ExpireIn              *int32  `json:"expireIn,omitempty"`
	PrivateTaker          *string `json:"privateTaker,omitempty"`
	DelegateSigner        *bool   `json:"delegateSigner,omitempty"`
	Compute               *int32  `json:"compute,omitempty"`
	PriorityMicroLamports *int32  `json:"priorityMicroLamports,omitempty"`
}

// ListNFTResponse represents the response from the list NFT API
type ListNFTResponse struct {
	Txs []Transaction `json:"txs"`
}

// DelistNFTRequest represents the request parameters for delisting an NFT
type DelistNFTRequest struct {
	Mint                  string  `json:"mint"`
	Owner                 string  `json:"owner"`
	Blockhash             string  `json:"blockhash"`
	FeePayer              *string `json:"feePayer,omitempty"`
	Compute               *int32  `json:"compute,omitempty"`
	PriorityMicroLamports *int32  `json:"priorityMicroLamports,omitempty"`
}

// DelistNFTResponse represents the response from the delist NFT API
type DelistNFTResponse struct {
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

// Validate validates the SellNFTRequest fields
func (r *SellNFTRequest) Validate() error {
	if r.Seller == "" {
		return fmt.Errorf("seller address is required")
	}

	if err := validateSolanaAddress(r.Seller); err != nil {
		return fmt.Errorf("invalid seller address: %w", err)
	}

	if r.Mint == "" {
		return fmt.Errorf("mint address is required")
	}

	if err := validateSolanaAddress(r.Mint); err != nil {
		return fmt.Errorf("invalid mint address: %w", err)
	}

	if r.BidAddress == "" {
		return fmt.Errorf("bidAddress is required")
	}

	if err := validateSolanaAddress(r.BidAddress); err != nil {
		return fmt.Errorf("invalid bidAddress: %w", err)
	}

	if r.MinPrice < 0 {
		return fmt.Errorf("minPrice must be >= 0")
	}

	if r.Blockhash == "" {
		return fmt.Errorf("blockhash is required")
	}

	// Validate optional addresses if provided
	if r.TakerBroker != nil {
		if err := validateSolanaAddress(*r.TakerBroker); err != nil {
			return fmt.Errorf("invalid takerBroker address: %w", err)
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

// Validate validates the ListNFTRequest fields
func (r *ListNFTRequest) Validate() error {
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

	if r.Price < 0 {
		return fmt.Errorf("price must be >= 0")
	}

	if r.Blockhash == "" {
		return fmt.Errorf("blockhash is required")
	}

	// Validate optional addresses if provided
	if r.MakerBroker != nil {
		if err := validateSolanaAddress(*r.MakerBroker); err != nil {
			return fmt.Errorf("invalid makerBroker address: %w", err)
		}
	}

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

	if r.RentPayer != nil {
		if err := validateSolanaAddress(*r.RentPayer); err != nil {
			return fmt.Errorf("invalid rentPayer address: %w", err)
		}
	}

	if r.Currency != nil {
		if err := validateSolanaAddress(*r.Currency); err != nil {
			return fmt.Errorf("invalid currency address: %w", err)
		}
	}

	if r.PrivateTaker != nil {
		if err := validateSolanaAddress(*r.PrivateTaker); err != nil {
			return fmt.Errorf("invalid privateTaker address: %w", err)
		}
	}

	// Validate expireIn
	if r.ExpireIn != nil && *r.ExpireIn < 0 {
		return fmt.Errorf("expireIn must be >= 0")
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

// Validate validates the DelistNFTRequest fields
func (r *DelistNFTRequest) Validate() error {
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

	if r.Blockhash == "" {
		return fmt.Errorf("blockhash is required")
	}

	// Validate optional addresses if provided
	if r.FeePayer != nil {
		if err := validateSolanaAddress(*r.FeePayer); err != nil {
			return fmt.Errorf("invalid feePayer address: %w", err)
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

// Validate validates the EditListingRequest fields
func (r *EditListingRequest) Validate() error {
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

	if r.Price < 0 {
		return fmt.Errorf("price must be >= 0")
	}

	if r.Blockhash == "" {
		return fmt.Errorf("blockhash is required")
	}

	// Validate optional addresses if provided
	if r.MakerBroker != nil {
		if err := validateSolanaAddress(*r.MakerBroker); err != nil {
			return fmt.Errorf("invalid makerBroker address: %w", err)
		}
	}

	if r.FeePayer != nil {
		if err := validateSolanaAddress(*r.FeePayer); err != nil {
			return fmt.Errorf("invalid feePayer address: %w", err)
		}
	}

	// Validate expireIn
	if r.ExpireIn != nil && *r.ExpireIn < 0 {
		return fmt.Errorf("expireIn must be >= 0")
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

// Validate validates the PlaceNFTBidRequest fields
func (r *PlaceNFTBidRequest) Validate() error {
	if r.Owner == "" {
		return fmt.Errorf("owner address is required")
	}

	if err := validateSolanaAddress(r.Owner); err != nil {
		return fmt.Errorf("invalid owner address: %w", err)
	}

	if r.Price < 0 {
		return fmt.Errorf("price must be >= 0")
	}

	if r.Mint == "" {
		return fmt.Errorf("mint address is required")
	}

	if err := validateSolanaAddress(r.Mint); err != nil {
		return fmt.Errorf("invalid mint address: %w", err)
	}

	if r.Blockhash == "" {
		return fmt.Errorf("blockhash is required")
	}

	// Validate optional addresses if provided
	if r.MakerBroker != nil {
		if err := validateSolanaAddress(*r.MakerBroker); err != nil {
			return fmt.Errorf("invalid makerBroker address: %w", err)
		}
	}

	if r.RentPayer != nil {
		if err := validateSolanaAddress(*r.RentPayer); err != nil {
			return fmt.Errorf("invalid rentPayer address: %w", err)
		}
	}

	// Validate expireIn
	if r.ExpireIn != nil && *r.ExpireIn < 0 {
		return fmt.Errorf("expireIn must be >= 0")
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

// Validate validates the PlaceTraitBidRequest fields
func (r *PlaceTraitBidRequest) Validate() error {
	if r.Owner == "" {
		return fmt.Errorf("owner address is required")
	}

	if err := validateSolanaAddress(r.Owner); err != nil {
		return fmt.Errorf("invalid owner address: %w", err)
	}

	if r.Price < 0 {
		return fmt.Errorf("price must be >= 0")
	}

	if r.Quantity < 1 {
		return fmt.Errorf("quantity must be >= 1")
	}

	if r.CollId == "" {
		return fmt.Errorf("collId is required")
	}

	if r.Blockhash == "" {
		return fmt.Errorf("blockhash is required")
	}

	// Validate optional addresses if provided
	if r.MakerBroker != nil {
		if err := validateSolanaAddress(*r.MakerBroker); err != nil {
			return fmt.Errorf("invalid makerBroker address: %w", err)
		}
	}

	if r.RentPayer != nil {
		if err := validateSolanaAddress(*r.RentPayer); err != nil {
			return fmt.Errorf("invalid rentPayer address: %w", err)
		}
	}

	// Validate expireIn
	if r.ExpireIn != nil && *r.ExpireIn < 0 {
		return fmt.Errorf("expireIn must be >= 0")
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

// Validate validates the PlaceCollectionBidRequest fields
func (r *PlaceCollectionBidRequest) Validate() error {
	if r.Owner == "" {
		return fmt.Errorf("owner address is required")
	}

	if err := validateSolanaAddress(r.Owner); err != nil {
		return fmt.Errorf("invalid owner address: %w", err)
	}

	if r.Price < 0 {
		return fmt.Errorf("price must be >= 0")
	}

	if r.Quantity < 1 {
		return fmt.Errorf("quantity must be >= 1")
	}

	if r.CollId == "" {
		return fmt.Errorf("collId is required")
	}

	if r.Blockhash == "" {
		return fmt.Errorf("blockhash is required")
	}

	// Validate optional addresses if provided
	if r.MakerBroker != nil {
		if err := validateSolanaAddress(*r.MakerBroker); err != nil {
			return fmt.Errorf("invalid makerBroker address: %w", err)
		}
	}

	if r.RentPayer != nil {
		if err := validateSolanaAddress(*r.RentPayer); err != nil {
			return fmt.Errorf("invalid rentPayer address: %w", err)
		}
	}

	// Validate expireIn
	if r.ExpireIn != nil && *r.ExpireIn < 0 {
		return fmt.Errorf("expireIn must be >= 0")
	}

	// Validate topUp
	if r.TopUp != nil && *r.TopUp < 0 {
		return fmt.Errorf("topUp must be >= 0")
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



// Validate validates the EditBidRequest fields
func (r *EditBidRequest) Validate() error {
	if r.BidStateAddress == "" {
		return fmt.Errorf("bidStateAddress is required")
	}

	if err := validateSolanaAddress(r.BidStateAddress); err != nil {
		return fmt.Errorf("invalid bidStateAddress: %w", err)
	}

	if r.Blockhash == "" {
		return fmt.Errorf("blockhash is required")
	}

	// Validate price if provided
	if r.Price != nil && *r.Price < 0 {
		return fmt.Errorf("price must be >= 0")
	}

	// Validate quantity if provided
	if r.Quantity != nil && *r.Quantity < 1 {
		return fmt.Errorf("quantity must be >= 1")
	}

	// Validate expireIn if provided
	if r.ExpireIn != nil && *r.ExpireIn < 0 {
		return fmt.Errorf("expireIn must be >= 0")
	}

	// Validate privateTaker if provided
	if r.PrivateTaker != nil {
		if err := validateSolanaAddress(*r.PrivateTaker); err != nil {
			return fmt.Errorf("invalid privateTaker address: %w", err)
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

// Validate validates the CancelBidRequest fields
func (r *CancelBidRequest) Validate() error {
	if r.BidStateAddress == "" {
		return fmt.Errorf("bidStateAddress is required")
	}

	if err := validateSolanaAddress(r.BidStateAddress); err != nil {
		return fmt.Errorf("invalid bidStateAddress: %w", err)
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
