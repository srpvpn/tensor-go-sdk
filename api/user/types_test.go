package user

import (
	"encoding/json"
	"testing"
)

func TestPortfolioRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request *PortfolioRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request",
			request: &PortfolioRequest{
				Wallet: "11111111111111111111111111111112",
			},
			wantErr: false,
		},
		{
			name: "empty wallet",
			request: &PortfolioRequest{
				Wallet: "",
			},
			wantErr: true,
			errMsg:  "wallet address is required",
		},
		{
			name: "invalid wallet - too short",
			request: &PortfolioRequest{
				Wallet: "short",
			},
			wantErr: true,
			errMsg:  "invalid wallet address",
		},
		{
			name: "invalid wallet - too long",
			request: &PortfolioRequest{
				Wallet: "11111111111111111111111111111111111111111111111",
			},
			wantErr: true,
			errMsg:  "invalid wallet address",
		},
		{
			name: "invalid wallet - invalid characters",
			request: &PortfolioRequest{
				Wallet: "1111111111111111111111111111111O", // Contains 'O'
			},
			wantErr: true,
			errMsg:  "invalid wallet address",
		},
		{
			name: "valid wallet with optional fields",
			request: &PortfolioRequest{
				Wallet:                "11111111111111111111111111111112",
				IncludeBidCount:       boolPtr(true),
				IncludeFavouriteCount: boolPtr(false),
				IncludeUnverified:     boolPtr(true),
				IncludeCompressed:     boolPtr(false),
				Currencies:            []string{"SOL", "USDC"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("PortfolioRequest.Validate() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg && !contains(err.Error(), tt.errMsg) {
					t.Errorf("PortfolioRequest.Validate() error = %v, want error containing %v", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("PortfolioRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestPortfolioRequest_JSON(t *testing.T) {
	tests := []struct {
		name     string
		request  *PortfolioRequest
		expected string
	}{
		{
			name: "minimal request",
			request: &PortfolioRequest{
				Wallet: "11111111111111111111111111111112",
			},
			expected: `{"wallet":"11111111111111111111111111111112"}`,
		},
		{
			name: "full request",
			request: &PortfolioRequest{
				Wallet:                "11111111111111111111111111111112",
				IncludeBidCount:       boolPtr(true),
				IncludeFavouriteCount: boolPtr(false),
				IncludeUnverified:     boolPtr(true),
				IncludeCompressed:     boolPtr(false),
				Currencies:            []string{"SOL", "USDC"},
			},
			expected: `{"wallet":"11111111111111111111111111111112","includeBidCount":true,"includeFavouriteCount":false,"includeUnverified":true,"includeCompressed":false,"currencies":["SOL","USDC"]}`,
		},
		{
			name: "request with nil optional fields",
			request: &PortfolioRequest{
				Wallet:                "11111111111111111111111111111112",
				IncludeBidCount:       nil,
				IncludeFavouriteCount: nil,
				IncludeUnverified:     nil,
				IncludeCompressed:     nil,
				Currencies:            nil,
			},
			expected: `{"wallet":"11111111111111111111111111111112"}`,
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
			var unmarshaled PortfolioRequest
			err = json.Unmarshal(data, &unmarshaled)
			if err != nil {
				t.Errorf("json.Unmarshal() error = %v", err)
				return
			}

			// Compare fields
			if unmarshaled.Wallet != tt.request.Wallet {
				t.Errorf("Unmarshaled wallet = %v, want %v", unmarshaled.Wallet, tt.request.Wallet)
			}

			if !equalBoolPtr(unmarshaled.IncludeBidCount, tt.request.IncludeBidCount) {
				t.Errorf("Unmarshaled IncludeBidCount = %v, want %v", unmarshaled.IncludeBidCount, tt.request.IncludeBidCount)
			}

			if !equalBoolPtr(unmarshaled.IncludeFavouriteCount, tt.request.IncludeFavouriteCount) {
				t.Errorf("Unmarshaled IncludeFavouriteCount = %v, want %v", unmarshaled.IncludeFavouriteCount, tt.request.IncludeFavouriteCount)
			}

			if !equalStringSlice(unmarshaled.Currencies, tt.request.Currencies) {
				t.Errorf("Unmarshaled Currencies = %v, want %v", unmarshaled.Currencies, tt.request.Currencies)
			}
		})
	}
}

func TestPortfolioResponse_JSON(t *testing.T) {
	tests := []struct {
		name     string
		response *PortfolioResponse
		jsonStr  string
	}{
		{
			name: "empty response",
			response: &PortfolioResponse{
				Message:     "success",
				Collections: []Collection{},
			},
			jsonStr: `{"message":"success"}`,
		},
		{
			name: "response with collections",
			response: &PortfolioResponse{
				Message: "success",
				Collections: []Collection{
					{
						ID:         "collection1",
						Name:       "Test Collection",
						Symbol:     "TEST",
						Image:      "https://example.com/image.png",
						FloorPrice: 1.5,
						Volume24h:  100.0,
						BidCount:   intPtr(5),
						FavCount:   intPtr(10),
						Verified:   true,
						Compressed: false,
					},
				},
			},
			jsonStr: `{"message":"success","collections":[{"id":"collection1","name":"Test Collection","symbol":"TEST","image":"https://example.com/image.png","floorPrice":1.5,"volume24h":100,"bidCount":5,"favCount":10,"verified":true,"compressed":false}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test marshaling
			data, err := json.Marshal(tt.response)
			if err != nil {
				t.Errorf("json.Marshal() error = %v", err)
				return
			}

			if string(data) != tt.jsonStr {
				t.Errorf("json.Marshal() = %v, want %v", string(data), tt.jsonStr)
			}

			// Test unmarshaling
			var unmarshaled PortfolioResponse
			err = json.Unmarshal([]byte(tt.jsonStr), &unmarshaled)
			if err != nil {
				t.Errorf("json.Unmarshal() error = %v", err)
				return
			}

			if unmarshaled.Message != tt.response.Message {
				t.Errorf("Unmarshaled Message = %v, want %v", unmarshaled.Message, tt.response.Message)
			}

			if len(unmarshaled.Collections) != len(tt.response.Collections) {
				t.Errorf("Unmarshaled Collections length = %v, want %v", len(unmarshaled.Collections), len(tt.response.Collections))
			}
		})
	}
}

func TestCollection_JSON(t *testing.T) {
	collection := Collection{
		ID:         "test-id",
		Name:       "Test Collection",
		Symbol:     "TEST",
		Image:      "https://example.com/image.png",
		FloorPrice: 2.5,
		Volume24h:  250.0,
		BidCount:   intPtr(15),
		FavCount:   nil, // Test nil pointer
		Verified:   true,
		Compressed: false,
	}

	expectedJSON := `{"id":"test-id","name":"Test Collection","symbol":"TEST","image":"https://example.com/image.png","floorPrice":2.5,"volume24h":250,"bidCount":15,"verified":true,"compressed":false}`

	// Test marshaling
	data, err := json.Marshal(collection)
	if err != nil {
		t.Errorf("json.Marshal() error = %v", err)
		return
	}

	if string(data) != expectedJSON {
		t.Errorf("json.Marshal() = %v, want %v", string(data), expectedJSON)
	}

	// Test unmarshaling
	var unmarshaled Collection
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Errorf("json.Unmarshal() error = %v", err)
		return
	}

	if unmarshaled.ID != collection.ID {
		t.Errorf("Unmarshaled ID = %v, want %v", unmarshaled.ID, collection.ID)
	}
	if unmarshaled.Name != collection.Name {
		t.Errorf("Unmarshaled Name = %v, want %v", unmarshaled.Name, collection.Name)
	}
	if !equalIntPtr(unmarshaled.BidCount, collection.BidCount) {
		t.Errorf("Unmarshaled BidCount = %v, want %v", unmarshaled.BidCount, collection.BidCount)
	}
	if !equalIntPtr(unmarshaled.FavCount, collection.FavCount) {
		t.Errorf("Unmarshaled FavCount = %v, want %v", unmarshaled.FavCount, collection.FavCount)
	}
}

func TestPortfolioRequest_UnmarshalJSON_TrimWhitespace(t *testing.T) {
	jsonStr := `{"wallet":"  11111111111111111111111111111112  "}`

	var request PortfolioRequest
	err := json.Unmarshal([]byte(jsonStr), &request)
	if err != nil {
		t.Errorf("json.Unmarshal() error = %v", err)
		return
	}

	expected := "11111111111111111111111111111112"
	if request.Wallet != expected {
		t.Errorf("Unmarshaled wallet = %v, want %v", request.Wallet, expected)
	}
}

// Helper functions
func boolPtr(b bool) *bool {
	return &b
}

func intPtr(i int) *int {
	return &i
}

func equalBoolPtr(a, b *bool) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func equalIntPtr(a, b *int) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func equalStringSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || (len(s) > len(substr) && s[:len(substr)] == substr) || (len(s) > len(substr) && s[len(s)-len(substr):] == substr) || (len(s) > len(substr) && findInString(s, substr)))
}

func findInString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestNFTBidsRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request *NFTBidsRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request",
			request: &NFTBidsRequest{
				Owner: "11111111111111111111111111111112",
				Limit: 100,
			},
			wantErr: false,
		},
		{
			name: "empty owner",
			request: &NFTBidsRequest{
				Owner: "",
				Limit: 100,
			},
			wantErr: true,
			errMsg:  "owner wallet address is required",
		},
		{
			name: "invalid owner - too short",
			request: &NFTBidsRequest{
				Owner: "short",
				Limit: 100,
			},
			wantErr: true,
			errMsg:  "invalid owner wallet address",
		},
		{
			name: "invalid owner - invalid characters",
			request: &NFTBidsRequest{
				Owner: "1111111111111111111111111111111O", // Contains 'O'
				Limit: 100,
			},
			wantErr: true,
			errMsg:  "invalid owner wallet address",
		},
		{
			name: "limit too low",
			request: &NFTBidsRequest{
				Owner: "11111111111111111111111111111112",
				Limit: 0,
			},
			wantErr: true,
			errMsg:  "limit must be between 1 and 500",
		},
		{
			name: "limit too high",
			request: &NFTBidsRequest{
				Owner: "11111111111111111111111111111112",
				Limit: 501,
			},
			wantErr: true,
			errMsg:  "limit must be between 1 and 500",
		},
		{
			name: "valid request with optional fields",
			request: &NFTBidsRequest{
				Owner:        "11111111111111111111111111111112",
				Limit:        50,
				CollId:       stringPtr("collection123"),
				Cursor:       stringPtr("cursor123"),
				BidAddresses: []string{"11111111111111111111111111111113", "11111111111111111111111111111114"},
			},
			wantErr: false,
		},
		{
			name: "invalid bid address",
			request: &NFTBidsRequest{
				Owner:        "11111111111111111111111111111112",
				Limit:        50,
				BidAddresses: []string{"invalid_address"},
			},
			wantErr: true,
			errMsg:  "invalid bid address",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("NFTBidsRequest.Validate() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg && !contains(err.Error(), tt.errMsg) {
					t.Errorf("NFTBidsRequest.Validate() error = %v, want error containing %v", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("NFTBidsRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestNFTBidsRequest_JSON(t *testing.T) {
	tests := []struct {
		name     string
		request  *NFTBidsRequest
		expected string
	}{
		{
			name: "minimal request",
			request: &NFTBidsRequest{
				Owner: "11111111111111111111111111111112",
				Limit: 100,
			},
			expected: `{"owner":"11111111111111111111111111111112","limit":100}`,
		},
		{
			name: "full request",
			request: &NFTBidsRequest{
				Owner:        "11111111111111111111111111111112",
				Limit:        50,
				CollId:       stringPtr("collection123"),
				Cursor:       stringPtr("cursor123"),
				BidAddresses: []string{"11111111111111111111111111111113", "11111111111111111111111111111114"},
			},
			expected: `{"owner":"11111111111111111111111111111112","limit":50,"collId":"collection123","cursor":"cursor123","bidAddresses":["11111111111111111111111111111113","11111111111111111111111111111114"]}`,
		},
		{
			name: "request with nil optional fields",
			request: &NFTBidsRequest{
				Owner:        "11111111111111111111111111111112",
				Limit:        100,
				CollId:       nil,
				Cursor:       nil,
				BidAddresses: nil,
			},
			expected: `{"owner":"11111111111111111111111111111112","limit":100}`,
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
			var unmarshaled NFTBidsRequest
			err = json.Unmarshal(data, &unmarshaled)
			if err != nil {
				t.Errorf("json.Unmarshal() error = %v", err)
				return
			}

			// Compare fields
			if unmarshaled.Owner != tt.request.Owner {
				t.Errorf("Unmarshaled Owner = %v, want %v", unmarshaled.Owner, tt.request.Owner)
			}

			if unmarshaled.Limit != tt.request.Limit {
				t.Errorf("Unmarshaled Limit = %v, want %v", unmarshaled.Limit, tt.request.Limit)
			}

			if !equalStringPtr(unmarshaled.CollId, tt.request.CollId) {
				t.Errorf("Unmarshaled CollId = %v, want %v", unmarshaled.CollId, tt.request.CollId)
			}

			if !equalStringPtr(unmarshaled.Cursor, tt.request.Cursor) {
				t.Errorf("Unmarshaled Cursor = %v, want %v", unmarshaled.Cursor, tt.request.Cursor)
			}

			if !equalStringSlice(unmarshaled.BidAddresses, tt.request.BidAddresses) {
				t.Errorf("Unmarshaled BidAddresses = %v, want %v", unmarshaled.BidAddresses, tt.request.BidAddresses)
			}
		})
	}
}

// Helper function for string pointer comparison
func stringPtr(s string) *string {
	return &s
}

func equalStringPtr(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func TestCollectionBidsRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request *CollectionBidsRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request",
			request: &CollectionBidsRequest{
				Owner: "11111111111111111111111111111112",
				Limit: 100,
			},
			wantErr: false,
		},
		{
			name: "empty owner",
			request: &CollectionBidsRequest{
				Owner: "",
				Limit: 100,
			},
			wantErr: true,
			errMsg:  "owner wallet address is required",
		},
		{
			name: "invalid owner - too short",
			request: &CollectionBidsRequest{
				Owner: "short",
				Limit: 100,
			},
			wantErr: true,
			errMsg:  "invalid owner wallet address",
		},
		{
			name: "invalid owner - invalid characters",
			request: &CollectionBidsRequest{
				Owner: "1111111111111111111111111111111O", // Contains 'O'
				Limit: 100,
			},
			wantErr: true,
			errMsg:  "invalid owner wallet address",
		},
		{
			name: "limit too low",
			request: &CollectionBidsRequest{
				Owner: "11111111111111111111111111111112",
				Limit: 0,
			},
			wantErr: true,
			errMsg:  "limit must be between 1 and 500",
		},
		{
			name: "limit too high",
			request: &CollectionBidsRequest{
				Owner: "11111111111111111111111111111112",
				Limit: 501,
			},
			wantErr: true,
			errMsg:  "limit must be between 1 and 500",
		},
		{
			name: "valid request with optional fields",
			request: &CollectionBidsRequest{
				Owner:        "11111111111111111111111111111112",
				Limit:        50,
				CollId:       stringPtr("collection123"),
				Cursor:       stringPtr("cursor123"),
				BidAddresses: []string{"11111111111111111111111111111113", "11111111111111111111111111111114"},
			},
			wantErr: false,
		},
		{
			name: "invalid bid address",
			request: &CollectionBidsRequest{
				Owner:        "11111111111111111111111111111112",
				Limit:        50,
				BidAddresses: []string{"invalid_address"},
			},
			wantErr: true,
			errMsg:  "invalid bid address",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("CollectionBidsRequest.Validate() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg && !contains(err.Error(), tt.errMsg) {
					t.Errorf("CollectionBidsRequest.Validate() error = %v, want error containing %v", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("CollectionBidsRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestCollectionBidsRequest_JSON(t *testing.T) {
	tests := []struct {
		name     string
		request  *CollectionBidsRequest
		expected string
	}{
		{
			name: "minimal request",
			request: &CollectionBidsRequest{
				Owner: "11111111111111111111111111111112",
				Limit: 100,
			},
			expected: `{"owner":"11111111111111111111111111111112","limit":100}`,
		},
		{
			name: "full request",
			request: &CollectionBidsRequest{
				Owner:        "11111111111111111111111111111112",
				Limit:        50,
				CollId:       stringPtr("collection123"),
				Cursor:       stringPtr("cursor123"),
				BidAddresses: []string{"11111111111111111111111111111113", "11111111111111111111111111111114"},
			},
			expected: `{"owner":"11111111111111111111111111111112","limit":50,"collId":"collection123","cursor":"cursor123","bidAddresses":["11111111111111111111111111111113","11111111111111111111111111111114"]}`,
		},
		{
			name: "request with nil optional fields",
			request: &CollectionBidsRequest{
				Owner:        "11111111111111111111111111111112",
				Limit:        100,
				CollId:       nil,
				Cursor:       nil,
				BidAddresses: nil,
			},
			expected: `{"owner":"11111111111111111111111111111112","limit":100}`,
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
			var unmarshaled CollectionBidsRequest
			err = json.Unmarshal(data, &unmarshaled)
			if err != nil {
				t.Errorf("json.Unmarshal() error = %v", err)
				return
			}

			// Compare fields
			if unmarshaled.Owner != tt.request.Owner {
				t.Errorf("Unmarshaled Owner = %v, want %v", unmarshaled.Owner, tt.request.Owner)
			}

			if unmarshaled.Limit != tt.request.Limit {
				t.Errorf("Unmarshaled Limit = %v, want %v", unmarshaled.Limit, tt.request.Limit)
			}

			if !equalStringPtr(unmarshaled.CollId, tt.request.CollId) {
				t.Errorf("Unmarshaled CollId = %v, want %v", unmarshaled.CollId, tt.request.CollId)
			}

			if !equalStringPtr(unmarshaled.Cursor, tt.request.Cursor) {
				t.Errorf("Unmarshaled Cursor = %v, want %v", unmarshaled.Cursor, tt.request.Cursor)
			}

			if !equalStringSlice(unmarshaled.BidAddresses, tt.request.BidAddresses) {
				t.Errorf("Unmarshaled BidAddresses = %v, want %v", unmarshaled.BidAddresses, tt.request.BidAddresses)
			}
		})
	}
}
func TestTraitBidsRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request *TraitBidsRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request",
			request: &TraitBidsRequest{
				Owner: "11111111111111111111111111111112",
				Limit: 100,
			},
			wantErr: false,
		},
		{
			name: "empty owner",
			request: &TraitBidsRequest{
				Owner: "",
				Limit: 100,
			},
			wantErr: true,
			errMsg:  "owner wallet address is required",
		},
		{
			name: "invalid owner - too short",
			request: &TraitBidsRequest{
				Owner: "short",
				Limit: 100,
			},
			wantErr: true,
			errMsg:  "invalid owner wallet address",
		},
		{
			name: "invalid owner - invalid characters",
			request: &TraitBidsRequest{
				Owner: "1111111111111111111111111111111O", // Contains 'O'
				Limit: 100,
			},
			wantErr: true,
			errMsg:  "invalid owner wallet address",
		},
		{
			name: "limit too low",
			request: &TraitBidsRequest{
				Owner: "11111111111111111111111111111112",
				Limit: 0,
			},
			wantErr: true,
			errMsg:  "limit must be between 1 and 500",
		},
		{
			name: "limit too high",
			request: &TraitBidsRequest{
				Owner: "11111111111111111111111111111112",
				Limit: 501,
			},
			wantErr: true,
			errMsg:  "limit must be between 1 and 500",
		},
		{
			name: "valid request with optional fields",
			request: &TraitBidsRequest{
				Owner:        "11111111111111111111111111111112",
				Limit:        50,
				CollId:       stringPtr("collection123"),
				Cursor:       stringPtr("cursor123"),
				BidAddresses: []string{"11111111111111111111111111111113", "11111111111111111111111111111114"},
			},
			wantErr: false,
		},
		{
			name: "invalid bid address",
			request: &TraitBidsRequest{
				Owner:        "11111111111111111111111111111112",
				Limit:        50,
				BidAddresses: []string{"invalid_address"},
			},
			wantErr: true,
			errMsg:  "invalid bid address",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("TraitBidsRequest.Validate() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg && !contains(err.Error(), tt.errMsg) {
					t.Errorf("TraitBidsRequest.Validate() error = %v, want error containing %v", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("TraitBidsRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestTraitBidsRequest_JSON(t *testing.T) {
	tests := []struct {
		name     string
		request  *TraitBidsRequest
		expected string
	}{
		{
			name: "minimal request",
			request: &TraitBidsRequest{
				Owner: "11111111111111111111111111111112",
				Limit: 100,
			},
			expected: `{"owner":"11111111111111111111111111111112","limit":100}`,
		},
		{
			name: "full request",
			request: &TraitBidsRequest{
				Owner:        "11111111111111111111111111111112",
				Limit:        50,
				CollId:       stringPtr("collection123"),
				Cursor:       stringPtr("cursor123"),
				BidAddresses: []string{"11111111111111111111111111111113", "11111111111111111111111111111114"},
			},
			expected: `{"owner":"11111111111111111111111111111112","limit":50,"collId":"collection123","cursor":"cursor123","bidAddresses":["11111111111111111111111111111113","11111111111111111111111111111114"]}`,
		},
		{
			name: "request with nil optional fields",
			request: &TraitBidsRequest{
				Owner:        "11111111111111111111111111111112",
				Limit:        100,
				CollId:       nil,
				Cursor:       nil,
				BidAddresses: nil,
			},
			expected: `{"owner":"11111111111111111111111111111112","limit":100}`,
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
			var unmarshaled TraitBidsRequest
			err = json.Unmarshal(data, &unmarshaled)
			if err != nil {
				t.Errorf("json.Unmarshal() error = %v", err)
				return
			}

			// Compare fields
			if unmarshaled.Owner != tt.request.Owner {
				t.Errorf("Unmarshaled Owner = %v, want %v", unmarshaled.Owner, tt.request.Owner)
			}

			if unmarshaled.Limit != tt.request.Limit {
				t.Errorf("Unmarshaled Limit = %v, want %v", unmarshaled.Limit, tt.request.Limit)
			}

			if !equalStringPtr(unmarshaled.CollId, tt.request.CollId) {
				t.Errorf("Unmarshaled CollId = %v, want %v", unmarshaled.CollId, tt.request.CollId)
			}

			if !equalStringPtr(unmarshaled.Cursor, tt.request.Cursor) {
				t.Errorf("Unmarshaled Cursor = %v, want %v", unmarshaled.Cursor, tt.request.Cursor)
			}

			if !equalStringSlice(unmarshaled.BidAddresses, tt.request.BidAddresses) {
				t.Errorf("Unmarshaled BidAddresses = %v, want %v", unmarshaled.BidAddresses, tt.request.BidAddresses)
			}
		})
	}
}
