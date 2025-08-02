package marketplace

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestBuyNFTRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request *BuyNFTRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request",
			request: &BuyNFTRequest{
				Buyer:     "11111111111111111111111111111112",
				Mint:      "11111111111111111111111111111113",
				Owner:     "11111111111111111111111111111114",
				MaxPrice:  1.5,
				Blockhash: "11111111111111111111111111111115",
			},
			wantErr: false,
		},
		{
			name: "empty buyer",
			request: &BuyNFTRequest{
				Buyer:     "",
				Mint:      "11111111111111111111111111111113",
				Owner:     "11111111111111111111111111111114",
				MaxPrice:  1.5,
				Blockhash: "11111111111111111111111111111115",
			},
			wantErr: true,
			errMsg:  "buyer address is required",
		},
		{
			name: "invalid buyer address",
			request: &BuyNFTRequest{
				Buyer:     "invalid",
				Mint:      "11111111111111111111111111111113",
				Owner:     "11111111111111111111111111111114",
				MaxPrice:  1.5,
				Blockhash: "11111111111111111111111111111115",
			},
			wantErr: true,
			errMsg:  "invalid buyer address",
		},
		{
			name: "empty mint",
			request: &BuyNFTRequest{
				Buyer:     "11111111111111111111111111111112",
				Mint:      "",
				Owner:     "11111111111111111111111111111114",
				MaxPrice:  1.5,
				Blockhash: "11111111111111111111111111111115",
			},
			wantErr: true,
			errMsg:  "mint address is required",
		},
		{
			name: "empty owner",
			request: &BuyNFTRequest{
				Buyer:     "11111111111111111111111111111112",
				Mint:      "11111111111111111111111111111113",
				Owner:     "",
				MaxPrice:  1.5,
				Blockhash: "11111111111111111111111111111115",
			},
			wantErr: true,
			errMsg:  "owner address is required",
		},
		{
			name: "negative max price",
			request: &BuyNFTRequest{
				Buyer:     "11111111111111111111111111111112",
				Mint:      "11111111111111111111111111111113",
				Owner:     "11111111111111111111111111111114",
				MaxPrice:  -1.0,
				Blockhash: "11111111111111111111111111111115",
			},
			wantErr: true,
			errMsg:  "maxPrice must be >= 0",
		},
		{
			name: "empty blockhash",
			request: &BuyNFTRequest{
				Buyer:     "11111111111111111111111111111112",
				Mint:      "11111111111111111111111111111113",
				Owner:     "11111111111111111111111111111114",
				MaxPrice:  1.5,
				Blockhash: "",
			},
			wantErr: true,
			errMsg:  "blockhash is required",
		},
		{
			name: "invalid optional royalty percent - too high",
			request: &BuyNFTRequest{
				Buyer:              "11111111111111111111111111111112",
				Mint:               "11111111111111111111111111111113",
				Owner:              "11111111111111111111111111111114",
				MaxPrice:           1.5,
				Blockhash:          "11111111111111111111111111111115",
				OptionalRoyaltyPct: int32Ptr(101),
			},
			wantErr: true,
			errMsg:  "optionalRoyaltyPct must be between 0 and 100",
		},
		{
			name: "invalid optional royalty percent - negative",
			request: &BuyNFTRequest{
				Buyer:              "11111111111111111111111111111112",
				Mint:               "11111111111111111111111111111113",
				Owner:              "11111111111111111111111111111114",
				MaxPrice:           1.5,
				Blockhash:          "11111111111111111111111111111115",
				OptionalRoyaltyPct: int32Ptr(-1),
			},
			wantErr: true,
			errMsg:  "optionalRoyaltyPct must be between 0 and 100",
		},
		{
			name: "negative compute",
			request: &BuyNFTRequest{
				Buyer:     "11111111111111111111111111111112",
				Mint:      "11111111111111111111111111111113",
				Owner:     "11111111111111111111111111111114",
				MaxPrice:  1.5,
				Blockhash: "11111111111111111111111111111115",
				Compute:   int32Ptr(-1),
			},
			wantErr: true,
			errMsg:  "compute must be >= 0",
		},
		{
			name: "negative priority micro lamports",
			request: &BuyNFTRequest{
				Buyer:                 "11111111111111111111111111111112",
				Mint:                  "11111111111111111111111111111113",
				Owner:                 "11111111111111111111111111111114",
				MaxPrice:              1.5,
				Blockhash:             "11111111111111111111111111111115",
				PriorityMicroLamports: int32Ptr(-1),
			},
			wantErr: true,
			errMsg:  "priorityMicroLamports must be >= 0",
		},
		{
			name: "valid request with all optional fields",
			request: &BuyNFTRequest{
				Buyer:                 "11111111111111111111111111111112",
				Mint:                  "11111111111111111111111111111113",
				Owner:                 "11111111111111111111111111111114",
				MaxPrice:              1.5,
				Blockhash:             "11111111111111111111111111111115",
				IncludeTotalCost:      boolPtr(true),
				Payer:                 stringPtr("11111111111111111111111111111116"),
				FeePayer:              stringPtr("11111111111111111111111111111117"),
				OptionalRoyaltyPct:    int32Ptr(5),
				Currency:              stringPtr("11111111111111111111111111111118"),
				TakerBroker:           stringPtr("11111111111111111111111111111119"),
				Compute:               int32Ptr(200000),
				PriorityMicroLamports: int32Ptr(1000),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("BuyNFTRequest.Validate() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg && !contains(err.Error(), tt.errMsg) {
					t.Errorf("BuyNFTRequest.Validate() error = %v, want error containing %v", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("BuyNFTRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestBuyNFTRequest_JSON(t *testing.T) {
	tests := []struct {
		name     string
		request  *BuyNFTRequest
		expected string
	}{
		{
			name: "minimal request",
			request: &BuyNFTRequest{
				Buyer:     "11111111111111111111111111111112",
				Mint:      "11111111111111111111111111111113",
				Owner:     "11111111111111111111111111111114",
				MaxPrice:  1.5,
				Blockhash: "11111111111111111111111111111115",
			},
			expected: `{"buyer":"11111111111111111111111111111112","mint":"11111111111111111111111111111113","owner":"11111111111111111111111111111114","maxPrice":1.5,"blockhash":"11111111111111111111111111111115"}`,
		},
		{
			name: "full request",
			request: &BuyNFTRequest{
				Buyer:                 "11111111111111111111111111111112",
				Mint:                  "11111111111111111111111111111113",
				Owner:                 "11111111111111111111111111111114",
				MaxPrice:              1.5,
				Blockhash:             "11111111111111111111111111111115",
				IncludeTotalCost:      boolPtr(true),
				Payer:                 stringPtr("11111111111111111111111111111116"),
				FeePayer:              stringPtr("11111111111111111111111111111117"),
				OptionalRoyaltyPct:    int32Ptr(5),
				Currency:              stringPtr("11111111111111111111111111111118"),
				TakerBroker:           stringPtr("11111111111111111111111111111119"),
				Compute:               int32Ptr(200000),
				PriorityMicroLamports: int32Ptr(1000),
			},
			expected: `{"buyer":"11111111111111111111111111111112","mint":"11111111111111111111111111111113","owner":"11111111111111111111111111111114","maxPrice":1.5,"blockhash":"11111111111111111111111111111115","includeTotalCost":true,"payer":"11111111111111111111111111111116","feePayer":"11111111111111111111111111111117","optionalRoyaltyPct":5,"currency":"11111111111111111111111111111118","takerBroker":"11111111111111111111111111111119","compute":200000,"priorityMicroLamports":1000}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test marshaling
			data, err := json.Marshal(tt.request)
			if err != nil {
				t.Errorf("json.Marshal() error = %v", err)
				return
			}

			if string(data) != tt.expected {
				t.Errorf("json.Marshal() = %v, want %v", string(data), tt.expected)
			}

			// Test unmarshaling
			var unmarshaled BuyNFTRequest
			err = json.Unmarshal(data, &unmarshaled)
			if err != nil {
				t.Errorf("json.Unmarshal() error = %v", err)
				return
			}

			// Compare fields
			if unmarshaled.Buyer != tt.request.Buyer {
				t.Errorf("Unmarshaled Buyer = %v, want %v", unmarshaled.Buyer, tt.request.Buyer)
			}
			if unmarshaled.Mint != tt.request.Mint {
				t.Errorf("Unmarshaled Mint = %v, want %v", unmarshaled.Mint, tt.request.Mint)
			}
			if unmarshaled.Owner != tt.request.Owner {
				t.Errorf("Unmarshaled Owner = %v, want %v", unmarshaled.Owner, tt.request.Owner)
			}
			if unmarshaled.MaxPrice != tt.request.MaxPrice {
				t.Errorf("Unmarshaled MaxPrice = %v, want %v", unmarshaled.MaxPrice, tt.request.MaxPrice)
			}
		})
	}
}

func TestBuyNFTResponse_JSON(t *testing.T) {
	responseJSON := `{
		"txs": [
			{
				"tx": "transaction_string",
				"txV0": "transaction_v0_string",
				"lastValidBlockHeight": 123456789,
				"metadata": {
					"key1": "value1",
					"key2": 42
				},
				"totalCost": 1.5
			}
		]
	}`

	var response BuyNFTResponse
	err := json.Unmarshal([]byte(responseJSON), &response)
	if err != nil {
		t.Errorf("json.Unmarshal() error = %v", err)
		return
	}

	if len(response.Txs) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(response.Txs))
		return
	}

	tx := response.Txs[0]
	if tx.Tx == nil || *tx.Tx != "transaction_string" {
		t.Errorf("Expected tx = 'transaction_string', got %v", tx.Tx)
	}
	if tx.TxV0 != "transaction_v0_string" {
		t.Errorf("Expected txV0 = 'transaction_v0_string', got %v", tx.TxV0)
	}
	if tx.LastValidBlockHeight == nil || *tx.LastValidBlockHeight != 123456789 {
		t.Errorf("Expected lastValidBlockHeight = 123456789, got %v", tx.LastValidBlockHeight)
	}
	if tx.TotalCost == nil || *tx.TotalCost != 1.5 {
		t.Errorf("Expected totalCost = 1.5, got %v", tx.TotalCost)
	}
}

// Helper functions
func boolPtr(b bool) *bool {
	return &b
}

func stringPtr(s string) *string {
	return &s
}

func int32Ptr(i int32) *int32 {
	return &i
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
func TestSellNFTRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request *SellNFTRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request",
			request: &SellNFTRequest{
				Seller:     "11111111111111111111111111111112",
				Mint:       "11111111111111111111111111111113",
				BidAddress: "11111111111111111111111111111114",
				MinPrice:   1.5,
				Blockhash:  "11111111111111111111111111111115",
			},
			wantErr: false,
		},
		{
			name: "empty seller",
			request: &SellNFTRequest{
				Seller:     "",
				Mint:       "11111111111111111111111111111113",
				BidAddress: "11111111111111111111111111111114",
				MinPrice:   1.5,
				Blockhash:  "11111111111111111111111111111115",
			},
			wantErr: true,
			errMsg:  "seller address is required",
		},
		{
			name: "empty bidAddress",
			request: &SellNFTRequest{
				Seller:     "11111111111111111111111111111112",
				Mint:       "11111111111111111111111111111113",
				BidAddress: "",
				MinPrice:   1.5,
				Blockhash:  "11111111111111111111111111111115",
			},
			wantErr: true,
			errMsg:  "bidAddress is required",
		},
		{
			name: "negative min price",
			request: &SellNFTRequest{
				Seller:     "11111111111111111111111111111112",
				Mint:       "11111111111111111111111111111113",
				BidAddress: "11111111111111111111111111111114",
				MinPrice:   -1.0,
				Blockhash:  "11111111111111111111111111111115",
			},
			wantErr: true,
			errMsg:  "minPrice must be >= 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("SellNFTRequest.Validate() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("SellNFTRequest.Validate() error = %v, want error containing %v", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("SellNFTRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestListNFTRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request *ListNFTRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request",
			request: &ListNFTRequest{
				Mint:      "11111111111111111111111111111112",
				Owner:     "11111111111111111111111111111113",
				Price:     2.5,
				Blockhash: "11111111111111111111111111111114",
			},
			wantErr: false,
		},
		{
			name: "empty mint",
			request: &ListNFTRequest{
				Mint:      "",
				Owner:     "11111111111111111111111111111113",
				Price:     2.5,
				Blockhash: "11111111111111111111111111111114",
			},
			wantErr: true,
			errMsg:  "mint address is required",
		},
		{
			name: "empty owner",
			request: &ListNFTRequest{
				Mint:      "11111111111111111111111111111112",
				Owner:     "",
				Price:     2.5,
				Blockhash: "11111111111111111111111111111114",
			},
			wantErr: true,
			errMsg:  "owner address is required",
		},
		{
			name: "negative price",
			request: &ListNFTRequest{
				Mint:      "11111111111111111111111111111112",
				Owner:     "11111111111111111111111111111113",
				Price:     -1.0,
				Blockhash: "11111111111111111111111111111114",
			},
			wantErr: true,
			errMsg:  "price must be >= 0",
		},
		{
			name: "negative expireIn",
			request: &ListNFTRequest{
				Mint:      "11111111111111111111111111111112",
				Owner:     "11111111111111111111111111111113",
				Price:     2.5,
				Blockhash: "11111111111111111111111111111114",
				ExpireIn:  int32Ptr(-1),
			},
			wantErr: true,
			errMsg:  "expireIn must be >= 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("ListNFTRequest.Validate() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("ListNFTRequest.Validate() error = %v, want error containing %v", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ListNFTRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestDelistNFTRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request *DelistNFTRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request",
			request: &DelistNFTRequest{
				Mint:      "11111111111111111111111111111112",
				Owner:     "11111111111111111111111111111113",
				Blockhash: "11111111111111111111111111111114",
			},
			wantErr: false,
		},
		{
			name: "empty mint",
			request: &DelistNFTRequest{
				Mint:      "",
				Owner:     "11111111111111111111111111111113",
				Blockhash: "11111111111111111111111111111114",
			},
			wantErr: true,
			errMsg:  "mint address is required",
		},
		{
			name: "empty owner",
			request: &DelistNFTRequest{
				Mint:      "11111111111111111111111111111112",
				Owner:     "",
				Blockhash: "11111111111111111111111111111114",
			},
			wantErr: true,
			errMsg:  "owner address is required",
		},
		{
			name: "empty blockhash",
			request: &DelistNFTRequest{
				Mint:      "11111111111111111111111111111112",
				Owner:     "11111111111111111111111111111113",
				Blockhash: "",
			},
			wantErr: true,
			errMsg:  "blockhash is required",
		},
		{
			name: "negative compute",
			request: &DelistNFTRequest{
				Mint:      "11111111111111111111111111111112",
				Owner:     "11111111111111111111111111111113",
				Blockhash: "11111111111111111111111111111114",
				Compute:   int32Ptr(-1),
			},
			wantErr: true,
			errMsg:  "compute must be >= 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("DelistNFTRequest.Validate() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("DelistNFTRequest.Validate() error = %v, want error containing %v", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("DelistNFTRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}
