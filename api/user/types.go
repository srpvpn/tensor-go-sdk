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

// TransactionsRequest represents the request parameters for getting user transactions
type TransactionsRequest struct {
	Wallets []string `json:"wallets"`
	Limit   int32    `json:"limit"`
	TxTypes []string `json:"txTypes"`
	Collid  string   `json:"collId"`
	Cursor  *string  `json:"cursor,omitempty"`
}

// TAmmPoolsRequest  represents the request parameters for getting user TAmm pools
type TAmmPoolsRequest struct {
	Owner         string   `json:"owner"`
	PoolAddresses []string `json:"poolAddresses"`
	Limit         int32    `json:"limit"`
	Cursor        *string  `json:"cursor,omitempty"`
}

// TSwapsPoolsRequest represents the request parameters for getting user TSwap pools
type TSwapsPoolsRequest struct {
	Owner         string   `json:"owner"`
	PoolAddresses []string `json:"poolAddresses"`
	Limit         int32    `json:"limit"`
	Cursor        *string  `json:"cursor,omitempty"`
}

// NFTBidsRequest represents the request parameters for getting user NFT bids
type NFTBidsRequest struct {
	Owner        string   `json:"owner"`
	Limit        int32    `json:"limit"`
	CollId       *string  `json:"collId,omitempty"`
	Cursor       *string  `json:"cursor,omitempty"`
	BidAddresses []string `json:"bidAddresses,omitempty"`
}

// CollectionBidsRequest represents the request parameters for getting user collection bids
type CollectionBidsRequest struct {
	Owner        string   `json:"owner"`
	Limit        int32    `json:"limit"`
	CollId       *string  `json:"collId,omitempty"`
	Cursor       *string  `json:"cursor,omitempty"`
	BidAddresses []string `json:"bidAddresses,omitempty"`
}

// TraitBidsRequest represents the request parameters for getting user trait bids
type TraitBidsRequest struct {
	Owner        string   `json:"owner"`
	Limit        int32    `json:"limit"`
	CollId       *string  `json:"collId,omitempty"`
	Cursor       *string  `json:"cursor,omitempty"`
	BidAddresses []string `json:"bidAddresses,omitempty"`
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

// Validate validates the NFTBidsRequest fields
func (r *NFTBidsRequest) Validate() error {
	if r.Owner == "" {
		return fmt.Errorf("owner wallet address is required")
	}

	if err := validateSolanaAddress(r.Owner); err != nil {
		return fmt.Errorf("invalid owner wallet address: %w", err)
	}

	if r.Limit <= 0 || r.Limit > 500 {
		return fmt.Errorf("limit must be between 1 and 500")
	}

	// Validate bid addresses if provided
	for _, bidAddr := range r.BidAddresses {
		if err := validateSolanaAddress(bidAddr); err != nil {
			return fmt.Errorf("invalid bid address %s: %w", bidAddr, err)
		}
	}

	return nil
}

// Validate validates the CollectionBidsRequest fields
func (r *CollectionBidsRequest) Validate() error {
	if r.Owner == "" {
		return fmt.Errorf("owner wallet address is required")
	}

	if err := validateSolanaAddress(r.Owner); err != nil {
		return fmt.Errorf("invalid owner wallet address: %w", err)
	}

	if r.Limit <= 0 || r.Limit > 500 {
		return fmt.Errorf("limit must be between 1 and 500")
	}

	// Validate bid addresses if provided
	for _, bidAddr := range r.BidAddresses {
		if err := validateSolanaAddress(bidAddr); err != nil {
			return fmt.Errorf("invalid bid address %s: %w", bidAddr, err)
		}
	}

	return nil
}

// Validate validates the TraitBidsRequest fields
func (r *TraitBidsRequest) Validate() error {
	if r.Owner == "" {
		return fmt.Errorf("owner wallet address is required")
	}

	if err := validateSolanaAddress(r.Owner); err != nil {
		return fmt.Errorf("invalid owner wallet address: %w", err)
	}

	if r.Limit <= 0 || r.Limit > 500 {
		return fmt.Errorf("limit must be between 1 and 500")
	}

	// Validate bid addresses if provided
	for _, bidAddr := range r.BidAddresses {
		if err := validateSolanaAddress(bidAddr); err != nil {
			return fmt.Errorf("invalid bid address %s: %w", bidAddr, err)
		}
	}

	return nil
}

func (r *TSwapsPoolsRequest) Validate() error {
	if r.Owner == "" {
		return fmt.Errorf("owner wallet address is required")
	}

	if err := validateSolanaAddress(r.Owner); err != nil {
		return fmt.Errorf("invalid owner wallet address: %w", err)
	}

	if r.Limit <= 0 || r.Limit > 500 {
		return fmt.Errorf("limit must be between 1 and 500")
	}

	// Validate pool addresses if provided
	for _, poolAddr := range r.PoolAddresses {
		if err := validateSolanaAddress(poolAddr); err != nil {
			return fmt.Errorf("invalid pool address %s: %w", poolAddr, err)
		}
	}

	return nil
}

func (r *TAmmPoolsRequest) Validate() error {
	if r.Owner == "" {
		return fmt.Errorf("owner wallet address is required")
	}

	if err := validateSolanaAddress(r.Owner); err != nil {
		return fmt.Errorf("invalid owner wallet address: %w", err)
	}

	if r.Limit <= 0 || r.Limit > 500 {
		return fmt.Errorf("limit must be between 1 and 500")
	}

	// Validate pool addresses if provided
	for _, poolAddr := range r.PoolAddresses {
		if err := validateSolanaAddress(poolAddr); err != nil {
			return fmt.Errorf("invalid pool address %s: %w", poolAddr, err)
		}
	}

	return nil
}


func (r *TransactionsRequest) Validate() error {
	if len(r.Wallets) == 0 {
		return fmt.Errorf("at least one wallet address is required")
	}

	for _, wallet := range r.Wallets {
		if err := validateSolanaAddress(wallet); err != nil {
			return fmt.Errorf("invalid wallet address %s: %w", wallet, err)
		}
	}

	if r.Limit <= 0 || r.Limit > 500 {
		return fmt.Errorf("limit must be between 1 and 500")
	}

	validTxTypes := []string{
		"LIST", "DELIST", "ADJUST_PRICE", "PLACE_BID", "CANCEL_BID", "SALE_BUY_NOW", "SALE_ACCEPT_BID", "TRANSFER",
		"FAILED", "OTHER", "AUCTION_CREATE", "AUCTION_PLACE_BID", "AUCTION_SETTLE", "AUCTION_CANCEL",
		"CREATE_MINT", "UPDATE_MINT",
		"SWAP_INIT_POOL", "SWAP_CLOSE_POOL", "SWAP_EDIT_POOL", "SWAP_DEPOSIT_NFT", "SWAP_DEPOSIT_SOL",
		"SWAP_BUY_NFT", "SWAP_SELL_NFT", "SWAP_WITHDRAW_NFT", "SWAP_WITHDRAW_SOL", "SWAP_WITHDRAW_MM_FEE",
		"SWAP_EDIT_SINGLE_LISTING", "SWAP_DEPOSIT_LIQ", "SWAP_WITHDRAW_LIQ", "SWAP_LIST", "SWAP_DELIST", "SWAP_BUY_SINGLE_LISTING",
		"ELIXIR_APPRAISE", "ELIXIR_FRACTIONALIZE", "ELIXIR_FUSE", "ELIXIR_POOL_DEPOSIT_FNFT", "ELIXIR_POOL_WITHDRAW_FNFT",
		"ELIXIR_POOL_EXCHANGE_FNFT", "ELIXIR_SELL_PNFT", "ELIXIR_BUY_PNFT", "ELIXIR_COMPOSED_BUY_NFT", "ELIXIR_COMPOSED_SELL_NFT",
		"MARGIN_INIT", "MARGIN_DEPOSIT", "MARGIN_WITHDRAW", "MARGIN_CLOSE", "MARGIN_ATTACH", "MARGIN_DETACH",
		"OTC_BUNDLED_MAKE_OFFER", "OTC_BUNDLED_TAKE_OFFER", "OTC_BUNDLED_TAKER_WITHDRAW", "OTC_BUNDLED_MAKER_WITHDRAW",
		"STAKE", "UNSTAKE",
		"ROLL_COMMIT", "ROLL_FULFILL_NONE", "ROLL_FULFILL_REWARD", "ROLL_FULFILL_SOL",
		"LOCK_UPSERT_ORDER", "LOCK_LOCK_ORDER", "LOCK_CLOSE_ORDER", "LOCK_WITHDRAW_NFT", "LOCK_DEPOSIT_NFT",
		"LOCK_WITHDRAW_COLLATERAL", "LOCK_CLAIM_TOKENS", "LOCK_CLAIM_NFT", "LOCK_ORDER_SELL_NFT", "LOCK_ORDER_BUY_NFT",
		"LOCK_MARKET_SELL_NFT", "LOCK_MARKET_BUY_NFT",
	}
	

	for _, txType := range r.TxTypes {
		isValid := false
		for _, validType := range validTxTypes {
			if txType == validType {
				isValid = true
				break
			}
		}
		if !isValid {
			return fmt.Errorf("invalid txType value: %s", txType)
		}
	}

	return nil
}
