package nfts

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestNFTsInfoRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request *NFTsInfoRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request with single mint",
			request: &NFTsInfoRequest{
				Mints: []string{"11111111111111111111111111111112"},
			},
			wantErr: false,
		},
		{
			name: "valid request with multiple mints",
			request: &NFTsInfoRequest{
				Mints: []string{
					"11111111111111111111111111111112",
					"11111111111111111111111111111113",
					"11111111111111111111111111111114",
				},
			},
			wantErr: false,
		},
		{
			name: "empty mints array",
			request: &NFTsInfoRequest{
				Mints: []string{},
			},
			wantErr: true,
			errMsg:  "mints is required and cannot be empty",
		},
		{
			name: "nil mints array",
			request: &NFTsInfoRequest{
				Mints: nil,
			},
			wantErr: true,
			errMsg:  "mints is required and cannot be empty",
		},
		{
			name: "empty mint address",
			request: &NFTsInfoRequest{
				Mints: []string{"11111111111111111111111111111112", ""},
			},
			wantErr: true,
			errMsg:  "mint address at index 1 cannot be empty",
		},
		{
			name: "invalid mint address",
			request: &NFTsInfoRequest{
				Mints: []string{"invalid-address"},
			},
			wantErr: true,
			errMsg:  "invalid mint address at index 0",
		},
	}

	// Helper function to check if error contains expected message
	contains := func(s, substr string) bool {
		return strings.Contains(s, substr)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("NFTsInfoRequest.Validate() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("NFTsInfoRequest.Validate() error = %v, want error containing %v", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("NFTsInfoRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestNFTsByCollectionRequest_Validate(t *testing.T) {
	// Helper functions for pointers
	boolPtr := func(b bool) *bool { return &b }
	stringPtr := func(s string) *string { return &s }
	int32Ptr := func(i int32) *int32 { return &i }
	float64Ptr := func(f float64) *float64 { return &f }

	tests := []struct {
		name    string
		request *NFTsByCollectionRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid minimal request",
			request: &NFTsByCollectionRequest{
				CollId: "collection-id",
				SortBy: "PriceAsc",
				Limit:  50,
			},
			wantErr: false,
		},
		{
			name: "valid request with all optional fields",
			request: &NFTsByCollectionRequest{
				CollId:            "collection-id",
				SortBy:            "PriceDesc",
				Limit:             100,
				OnlyListings:      boolPtr(true),
				Mints:             []string{"11111111111111111111111111111112"},
				Cursor:            stringPtr("cursor-string"),
				ListingSources:    []string{"tensor", "magiceden"},
				MinPrice:          float64Ptr(0.5),
				MaxPrice:          float64Ptr(10.0),
				TraitCountMin:     int32Ptr(1),
				TraitCountMax:     int32Ptr(10),
				Name:              stringPtr("Cool NFT"),
				ExcludeOwners:     []string{"11111111111111111111111111111113"},
				IncludeOwners:     []string{"11111111111111111111111111111114"},
				IncludeCurrencies: []string{"SOL", "USDC"},
				Traits:            []string{`{"trait_type": ["value1", "value2"]}`},
				RaritySystem:      stringPtr("tensor"),
				RarityMin:         float64Ptr(1.0),
				RarityMax:         float64Ptr(100.0),
				OnlyInscriptions:  boolPtr(false),
				ImmutableStatus:   stringPtr("mutable"),
			},
			wantErr: false,
		},
		{
			name: "empty collId",
			request: &NFTsByCollectionRequest{
				CollId: "",
				SortBy: "PriceAsc",
				Limit:  50,
			},
			wantErr: true,
			errMsg:  "collId is required",
		},
		{
			name: "empty sortBy",
			request: &NFTsByCollectionRequest{
				CollId: "collection-id",
				SortBy: "",
				Limit:  50,
			},
			wantErr: true,
			errMsg:  "sortBy is required",
		},
		{
			name: "limit too low",
			request: &NFTsByCollectionRequest{
				CollId: "collection-id",
				SortBy: "PriceAsc",
				Limit:  0,
			},
			wantErr: true,
			errMsg:  "limit must be between 1 and 250",
		},
		{
			name: "limit too high",
			request: &NFTsByCollectionRequest{
				CollId: "collection-id",
				SortBy: "PriceAsc",
				Limit:  251,
			},
			wantErr: true,
			errMsg:  "limit must be between 1 and 250",
		},
		{
			name: "invalid mint address",
			request: &NFTsByCollectionRequest{
				CollId: "collection-id",
				SortBy: "PriceAsc",
				Limit:  50,
				Mints:  []string{"invalid-address"},
			},
			wantErr: true,
			errMsg:  "invalid mint address at index 0",
		},
		{
			name: "invalid exclude owner address",
			request: &NFTsByCollectionRequest{
				CollId:        "collection-id",
				SortBy:        "PriceAsc",
				Limit:         50,
				ExcludeOwners: []string{"invalid-address"},
			},
			wantErr: true,
			errMsg:  "invalid exclude owner address at index 0",
		},
		{
			name: "invalid include owner address",
			request: &NFTsByCollectionRequest{
				CollId:        "collection-id",
				SortBy:        "PriceAsc",
				Limit:         50,
				IncludeOwners: []string{"invalid-address"},
			},
			wantErr: true,
			errMsg:  "invalid include owner address at index 0",
		},
		{
			name: "negative minPrice",
			request: &NFTsByCollectionRequest{
				CollId:   "collection-id",
				SortBy:   "PriceAsc",
				Limit:    50,
				MinPrice: float64Ptr(-1.0),
			},
			wantErr: true,
			errMsg:  "minPrice must be >= 0",
		},
		{
			name: "negative maxPrice",
			request: &NFTsByCollectionRequest{
				CollId:   "collection-id",
				SortBy:   "PriceAsc",
				Limit:    50,
				MaxPrice: float64Ptr(-1.0),
			},
			wantErr: true,
			errMsg:  "maxPrice must be >= 0",
		},
		{
			name: "negative traitCountMin",
			request: &NFTsByCollectionRequest{
				CollId:        "collection-id",
				SortBy:        "PriceAsc",
				Limit:         50,
				TraitCountMin: int32Ptr(-1),
			},
			wantErr: true,
			errMsg:  "traitCountMin must be >= 0",
		},
		{
			name: "invalid traitCountMax",
			request: &NFTsByCollectionRequest{
				CollId:        "collection-id",
				SortBy:        "PriceAsc",
				Limit:         50,
				TraitCountMax: int32Ptr(0),
			},
			wantErr: true,
			errMsg:  "traitCountMax must be >= 1",
		},
		{
			name: "negative rarityMin",
			request: &NFTsByCollectionRequest{
				CollId:    "collection-id",
				SortBy:    "PriceAsc",
				Limit:     50,
				RarityMin: float64Ptr(-1.0),
			},
			wantErr: true,
			errMsg:  "rarityMin must be >= 0",
		},
		{
			name: "negative rarityMax",
			request: &NFTsByCollectionRequest{
				CollId:    "collection-id",
				SortBy:    "PriceAsc",
				Limit:     50,
				RarityMax: float64Ptr(-1.0),
			},
			wantErr: true,
			errMsg:  "rarityMax must be >= 0",
		},
	}

	// Helper function to check if error contains expected message
	contains := func(s, substr string) bool {
		return strings.Contains(s, substr)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("NFTsByCollectionRequest.Validate() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("NFTsByCollectionRequest.Validate() error = %v, want error containing %v", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("NFTsByCollectionRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestNFTsInfoRequest_JSON(t *testing.T) {
	tests := []struct {
		name     string
		request  *NFTsInfoRequest
		expected string
	}{
		{
			name: "single mint",
			request: &NFTsInfoRequest{
				Mints: []string{"11111111111111111111111111111112"},
			},
			expected: `{"mints":["11111111111111111111111111111112"]}`,
		},
		{
			name: "multiple mints",
			request: &NFTsInfoRequest{
				Mints: []string{
					"11111111111111111111111111111112",
					"11111111111111111111111111111113",
				},
			},
			expected: `{"mints":["11111111111111111111111111111112","11111111111111111111111111111113"]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBytes, err := json.Marshal(tt.request)
			if err != nil {
				t.Errorf("json.Marshal() error = %v", err)
				return
			}

			if string(jsonBytes) != tt.expected {
				t.Errorf("json.Marshal() = %v, want %v", string(jsonBytes), tt.expected)
			}

			// Test unmarshaling
			var unmarshaled NFTsInfoRequest
			err = json.Unmarshal(jsonBytes, &unmarshaled)
			if err != nil {
				t.Errorf("json.Unmarshal() error = %v", err)
			}
		})
	}
}

func TestNFTsByCollectionRequest_JSON(t *testing.T) {
	// Helper functions for pointers
	boolPtr := func(b bool) *bool { return &b }
	stringPtr := func(s string) *string { return &s }
	float64Ptr := func(f float64) *float64 { return &f }

	tests := []struct {
		name     string
		request  *NFTsByCollectionRequest
		expected string
	}{
		{
			name: "minimal request",
			request: &NFTsByCollectionRequest{
				CollId: "collection-id",
				SortBy: "PriceAsc",
				Limit:  50,
			},
			expected: `{"collId":"collection-id","sortBy":"PriceAsc","limit":50}`,
		},
		{
			name: "request with optional fields",
			request: &NFTsByCollectionRequest{
				CollId:       "collection-id",
				SortBy:       "PriceDesc",
				Limit:        100,
				OnlyListings: boolPtr(true),
				MinPrice:     float64Ptr(1.5),
				MaxPrice:     float64Ptr(10.0),
				Name:         stringPtr("Cool NFT"),
			},
			expected: `{"collId":"collection-id","sortBy":"PriceDesc","limit":100,"onlyListings":true,"minPrice":1.5,"maxPrice":10,"name":"Cool NFT"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBytes, err := json.Marshal(tt.request)
			if err != nil {
				t.Errorf("json.Marshal() error = %v", err)
				return
			}

			if string(jsonBytes) != tt.expected {
				t.Errorf("json.Marshal() = %v, want %v", string(jsonBytes), tt.expected)
			}

			// Test unmarshaling
			var unmarshaled NFTsByCollectionRequest
			err = json.Unmarshal(jsonBytes, &unmarshaled)
			if err != nil {
				t.Errorf("json.Unmarshal() error = %v", err)
			}
		})
	}
}

func TestNFTsInfoRequest_UnmarshalJSON_TrimWhitespace(t *testing.T) {
	jsonStr := `{"mints":["  11111111111111111111111111111112  ","  11111111111111111111111111111113  "]}`

	var request NFTsInfoRequest
	err := json.Unmarshal([]byte(jsonStr), &request)
	if err != nil {
		t.Errorf("json.Unmarshal() error = %v", err)
		return
	}

	expectedMints := []string{
		"11111111111111111111111111111112",
		"11111111111111111111111111111113",
	}

	if len(request.Mints) != len(expectedMints) {
		t.Errorf("Unmarshaled Mints length = %v, want %v", len(request.Mints), len(expectedMints))
		return
	}

	for i, mint := range request.Mints {
		if mint != expectedMints[i] {
			t.Errorf("Unmarshaled Mint[%d] = %v, want %v", i, mint, expectedMints[i])
		}
	}
}

func TestNFTsByCollectionRequest_UnmarshalJSON_TrimWhitespace(t *testing.T) {
	jsonStr := `{"collId":"  collection-id  ","sortBy":"  PriceAsc  ","limit":50,"name":"  Cool NFT  "}`

	var request NFTsByCollectionRequest
	err := json.Unmarshal([]byte(jsonStr), &request)
	if err != nil {
		t.Errorf("json.Unmarshal() error = %v", err)
		return
	}

	expectedCollId := "collection-id"
	expectedSortBy := "PriceAsc"
	expectedName := "Cool NFT"

	if request.CollId != expectedCollId {
		t.Errorf("Unmarshaled CollId = %v, want %v", request.CollId, expectedCollId)
	}
	if request.SortBy != expectedSortBy {
		t.Errorf("Unmarshaled SortBy = %v, want %v", request.SortBy, expectedSortBy)
	}
	if request.Name == nil || *request.Name != expectedName {
		t.Errorf("Unmarshaled Name = %v, want %v", request.Name, expectedName)
	}
}
