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
