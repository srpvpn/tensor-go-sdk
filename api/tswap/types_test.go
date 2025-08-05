package tswap

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestCloseTSwapPoolRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request *CloseTSwapPoolRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request",
			request: &CloseTSwapPoolRequest{
				PoolAddress: "11111111111111111111111111111112",
				Blockhash:   "11111111111111111111111111111113",
			},
			wantErr: false,
		},
		{
			name: "empty pool address",
			request: &CloseTSwapPoolRequest{
				PoolAddress: "",
				Blockhash:   "11111111111111111111111111111113",
			},
			wantErr: true,
			errMsg:  "poolAddress is required",
		},
		{
			name: "invalid pool address",
			request: &CloseTSwapPoolRequest{
				PoolAddress: "invalid",
				Blockhash:   "11111111111111111111111111111113",
			},
			wantErr: true,
			errMsg:  "invalid poolAddress",
		},
		{
			name: "empty blockhash",
			request: &CloseTSwapPoolRequest{
				PoolAddress: "11111111111111111111111111111112",
				Blockhash:   "",
			},
			wantErr: true,
			errMsg:  "blockhash is required",
		},
		{
			name: "negative compute",
			request: &CloseTSwapPoolRequest{
				PoolAddress: "11111111111111111111111111111112",
				Blockhash:   "11111111111111111111111111111113",
				Compute:     int32Ptr(-1),
			},
			wantErr: true,
			errMsg:  "compute must be >= 0",
		},
		{
			name: "negative priority micro lamports",
			request: &CloseTSwapPoolRequest{
				PoolAddress:           "11111111111111111111111111111112",
				Blockhash:             "11111111111111111111111111111113",
				PriorityMicroLamports: int32Ptr(-1),
			},
			wantErr: true,
			errMsg:  "priorityMicroLamports must be >= 0",
		},
		{
			name: "valid request with optional fields",
			request: &CloseTSwapPoolRequest{
				PoolAddress:           "11111111111111111111111111111112",
				Blockhash:             "11111111111111111111111111111113",
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
					t.Errorf("CloseTSwapPoolRequest.Validate() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("CloseTSwapPoolRequest.Validate() error = %v, want error containing %v", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("CloseTSwapPoolRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestEditTSwapPoolRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request *EditTSwapPoolRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request",
			request: &EditTSwapPoolRequest{
				PoolAddress:   "11111111111111111111111111111112",
				PoolType:      "TOKEN",
				CurveType:     "linear",
				StartingPrice: 1.5,
				Delta:         0.1,
				Blockhash:     "11111111111111111111111111111113",
			},
			wantErr: false,
		},
		{
			name: "empty pool address",
			request: &EditTSwapPoolRequest{
				PoolAddress:   "",
				PoolType:      "TOKEN",
				CurveType:     "linear",
				StartingPrice: 1.5,
				Delta:         0.1,
				Blockhash:     "11111111111111111111111111111113",
			},
			wantErr: true,
			errMsg:  "poolAddress is required",
		},
		{
			name: "invalid pool type",
			request: &EditTSwapPoolRequest{
				PoolAddress:   "11111111111111111111111111111112",
				PoolType:      "INVALID",
				CurveType:     "linear",
				StartingPrice: 1.5,
				Delta:         0.1,
				Blockhash:     "11111111111111111111111111111113",
			},
			wantErr: true,
			errMsg:  "invalid poolType",
		},
		{
			name: "invalid curve type",
			request: &EditTSwapPoolRequest{
				PoolAddress:   "11111111111111111111111111111112",
				PoolType:      "TOKEN",
				CurveType:     "invalid",
				StartingPrice: 1.5,
				Delta:         0.1,
				Blockhash:     "11111111111111111111111111111113",
			},
			wantErr: true,
			errMsg:  "invalid curveType",
		},
		{
			name: "negative starting price",
			request: &EditTSwapPoolRequest{
				PoolAddress:   "11111111111111111111111111111112",
				PoolType:      "TOKEN",
				CurveType:     "linear",
				StartingPrice: -1.0,
				Delta:         0.1,
				Blockhash:     "11111111111111111111111111111113",
			},
			wantErr: true,
			errMsg:  "startingPrice must be >= 0",
		},
		{
			name: "negative delta",
			request: &EditTSwapPoolRequest{
				PoolAddress:   "11111111111111111111111111111112",
				PoolType:      "TOKEN",
				CurveType:     "linear",
				StartingPrice: 1.5,
				Delta:         -0.1,
				Blockhash:     "11111111111111111111111111111113",
			},
			wantErr: true,
			errMsg:  "delta must be >= 0",
		},
		{
			name: "invalid mm fee bps - too high",
			request: &EditTSwapPoolRequest{
				PoolAddress:   "11111111111111111111111111111112",
				PoolType:      "TOKEN",
				CurveType:     "linear",
				StartingPrice: 1.5,
				Delta:         0.1,
				Blockhash:     "11111111111111111111111111111113",
				MmFeeBps:      float64Ptr(10001),
			},
			wantErr: true,
			errMsg:  "mmFeeBps must be between 0 and 10000 basis points",
		},
		{
			name: "negative max taker sell count",
			request: &EditTSwapPoolRequest{
				PoolAddress:       "11111111111111111111111111111112",
				PoolType:          "TOKEN",
				CurveType:         "linear",
				StartingPrice:     1.5,
				Delta:             0.1,
				Blockhash:         "11111111111111111111111111111113",
				MaxTakerSellCount: int32Ptr(-1),
			},
			wantErr: true,
			errMsg:  "maxTakerSellCount must be >= 0",
		},
		{
			name: "valid request with all optional fields",
			request: &EditTSwapPoolRequest{
				PoolAddress:           "11111111111111111111111111111112",
				PoolType:              "NFT",
				CurveType:             "exponential",
				StartingPrice:         2.5,
				Delta:                 0.2,
				Blockhash:             "11111111111111111111111111111113",
				MmKeepFeesSeparate:    boolPtr(true),
				MmFeeBps:              float64Ptr(250.5),
				MaxTakerSellCount:     int32Ptr(10),
				UseSharedEscrow:       boolPtr(false),
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
					t.Errorf("EditTSwapPoolRequest.Validate() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("EditTSwapPoolRequest.Validate() error = %v, want error containing %v", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("EditTSwapPoolRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestCloseTSwapPoolRequest_JSON(t *testing.T) {
	tests := []struct {
		name     string
		request  *CloseTSwapPoolRequest
		expected string
	}{
		{
			name: "minimal request",
			request: &CloseTSwapPoolRequest{
				PoolAddress: "11111111111111111111111111111112",
				Blockhash:   "11111111111111111111111111111113",
			},
			expected: `{"poolAddress":"11111111111111111111111111111112","blockhash":"11111111111111111111111111111113"}`,
		},
		{
			name: "full request",
			request: &CloseTSwapPoolRequest{
				PoolAddress:           "11111111111111111111111111111112",
				Blockhash:             "11111111111111111111111111111113",
				Compute:               int32Ptr(200000),
				PriorityMicroLamports: int32Ptr(1000),
			},
			expected: `{"poolAddress":"11111111111111111111111111111112","blockhash":"11111111111111111111111111111113","compute":200000,"priorityMicroLamports":1000}`,
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
			var unmarshaled CloseTSwapPoolRequest
			err = json.Unmarshal(data, &unmarshaled)
			if err != nil {
				t.Errorf("json.Unmarshal() error = %v", err)
				return
			}

			// Compare fields
			if unmarshaled.PoolAddress != tt.request.PoolAddress {
				t.Errorf("Unmarshaled PoolAddress = %v, want %v", unmarshaled.PoolAddress, tt.request.PoolAddress)
			}
			if unmarshaled.Blockhash != tt.request.Blockhash {
				t.Errorf("Unmarshaled Blockhash = %v, want %v", unmarshaled.Blockhash, tt.request.Blockhash)
			}
		})
	}
}

func TestEditTSwapPoolRequest_JSON(t *testing.T) {
	tests := []struct {
		name     string
		request  *EditTSwapPoolRequest
		expected string
	}{
		{
			name: "minimal request",
			request: &EditTSwapPoolRequest{
				PoolAddress:   "11111111111111111111111111111112",
				PoolType:      "TOKEN",
				CurveType:     "linear",
				StartingPrice: 1.5,
				Delta:         0.1,
				Blockhash:     "11111111111111111111111111111113",
			},
			expected: `{"poolAddress":"11111111111111111111111111111112","poolType":"TOKEN","curveType":"linear","startingPrice":1.5,"delta":0.1,"blockhash":"11111111111111111111111111111113"}`,
		},
		{
			name: "full request",
			request: &EditTSwapPoolRequest{
				PoolAddress:           "11111111111111111111111111111112",
				PoolType:              "NFT",
				CurveType:             "exponential",
				StartingPrice:         2.5,
				Delta:                 0.2,
				Blockhash:             "11111111111111111111111111111113",
				MmKeepFeesSeparate:    boolPtr(true),
				MmFeeBps:              float64Ptr(250.5),
				MaxTakerSellCount:     int32Ptr(10),
				UseSharedEscrow:       boolPtr(false),
				Compute:               int32Ptr(200000),
				PriorityMicroLamports: int32Ptr(1000),
			},
			expected: `{"poolAddress":"11111111111111111111111111111112","poolType":"NFT","curveType":"exponential","startingPrice":2.5,"delta":0.2,"blockhash":"11111111111111111111111111111113","mmKeepFeesSeparate":true,"mmFeeBps":250.5,"maxTakerSellCount":10,"useSharedEscrow":false,"compute":200000,"priorityMicroLamports":1000}`,
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
			var unmarshaled EditTSwapPoolRequest
			err = json.Unmarshal(data, &unmarshaled)
			if err != nil {
				t.Errorf("json.Unmarshal() error = %v", err)
				return
			}

			// Compare fields
			if unmarshaled.PoolAddress != tt.request.PoolAddress {
				t.Errorf("Unmarshaled PoolAddress = %v, want %v", unmarshaled.PoolAddress, tt.request.PoolAddress)
			}
			if unmarshaled.PoolType != tt.request.PoolType {
				t.Errorf("Unmarshaled PoolType = %v, want %v", unmarshaled.PoolType, tt.request.PoolType)
			}
		})
	}
}

func TestCloseTSwapPoolRequest_UnmarshalJSON_TrimWhitespace(t *testing.T) {
	jsonStr := `{"poolAddress":"  11111111111111111111111111111112  ","blockhash":"  11111111111111111111111111111113  "}`

	var request CloseTSwapPoolRequest
	err := json.Unmarshal([]byte(jsonStr), &request)
	if err != nil {
		t.Errorf("json.Unmarshal() error = %v", err)
		return
	}

	expectedPoolAddress := "11111111111111111111111111111112"
	expectedBlockhash := "11111111111111111111111111111113"

	if request.PoolAddress != expectedPoolAddress {
		t.Errorf("Unmarshaled PoolAddress = %v, want %v", request.PoolAddress, expectedPoolAddress)
	}
	if request.Blockhash != expectedBlockhash {
		t.Errorf("Unmarshaled Blockhash = %v, want %v", request.Blockhash, expectedBlockhash)
	}
}

func TestEditTSwapPoolRequest_UnmarshalJSON_TrimWhitespace(t *testing.T) {
	jsonStr := `{"poolAddress":"  11111111111111111111111111111112  ","poolType":"  TOKEN  ","curveType":"  linear  ","startingPrice":1.5,"delta":0.1,"blockhash":"  11111111111111111111111111111113  "}`

	var request EditTSwapPoolRequest
	err := json.Unmarshal([]byte(jsonStr), &request)
	if err != nil {
		t.Errorf("json.Unmarshal() error = %v", err)
		return
	}

	expectedPoolAddress := "11111111111111111111111111111112"
	expectedPoolType := "TOKEN"
	expectedCurveType := "linear"
	expectedBlockhash := "11111111111111111111111111111113"

	if request.PoolAddress != expectedPoolAddress {
		t.Errorf("Unmarshaled PoolAddress = %v, want %v", request.PoolAddress, expectedPoolAddress)
	}
	if request.PoolType != expectedPoolType {
		t.Errorf("Unmarshaled PoolType = %v, want %v", request.PoolType, expectedPoolType)
	}
	if request.CurveType != expectedCurveType {
		t.Errorf("Unmarshaled CurveType = %v, want %v", request.CurveType, expectedCurveType)
	}
	if request.Blockhash != expectedBlockhash {
		t.Errorf("Unmarshaled Blockhash = %v, want %v", request.Blockhash, expectedBlockhash)
	}
}

// Helper functions
func boolPtr(b bool) *bool {
	return &b
}

func int32Ptr(i int32) *int32 {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
func TestDepositWithdrawNFTRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request *DepositWithdrawNFTRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid deposit request - uppercase",
			request: &DepositWithdrawNFTRequest{
				Action:      "DEPOSIT",
				PoolAddress: "11111111111111111111111111111112",
				Mint:        "11111111111111111111111111111113",
				Blockhash:   "11111111111111111111111111111114",
			},
			wantErr: false,
		},
		{
			name: "valid deposit request - lowercase (auto-normalized)",
			request: &DepositWithdrawNFTRequest{
				Action:      "deposit",
				PoolAddress: "11111111111111111111111111111112",
				Mint:        "11111111111111111111111111111113",
				Blockhash:   "11111111111111111111111111111114",
			},
			wantErr: false,
		},
		{
			name: "valid withdraw request - uppercase",
			request: &DepositWithdrawNFTRequest{
				Action:      "WITHDRAW",
				PoolAddress: "11111111111111111111111111111112",
				Mint:        "11111111111111111111111111111113",
				Blockhash:   "11111111111111111111111111111114",
			},
			wantErr: false,
		},
		{
			name: "valid withdraw request - mixed case (auto-normalized)",
			request: &DepositWithdrawNFTRequest{
				Action:      "WithDraw",
				PoolAddress: "11111111111111111111111111111112",
				Mint:        "11111111111111111111111111111113",
				Blockhash:   "11111111111111111111111111111114",
			},
			wantErr: false,
		},
		{
			name: "empty action",
			request: &DepositWithdrawNFTRequest{
				Action:      "",
				PoolAddress: "11111111111111111111111111111112",
				Mint:        "11111111111111111111111111111113",
				Blockhash:   "11111111111111111111111111111114",
			},
			wantErr: true,
			errMsg:  "action is required",
		},
		{
			name: "invalid action",
			request: &DepositWithdrawNFTRequest{
				Action:      "INVALID",
				PoolAddress: "11111111111111111111111111111112",
				Mint:        "11111111111111111111111111111113",
				Blockhash:   "11111111111111111111111111111114",
			},
			wantErr: true,
			errMsg:  "invalid action",
		},
		{
			name: "empty pool address",
			request: &DepositWithdrawNFTRequest{
				Action:      "DEPOSIT",
				PoolAddress: "",
				Mint:        "11111111111111111111111111111113",
				Blockhash:   "11111111111111111111111111111114",
			},
			wantErr: true,
			errMsg:  "poolAddress is required",
		},
		{
			name: "invalid pool address",
			request: &DepositWithdrawNFTRequest{
				Action:      "DEPOSIT",
				PoolAddress: "invalid",
				Mint:        "11111111111111111111111111111113",
				Blockhash:   "11111111111111111111111111111114",
			},
			wantErr: true,
			errMsg:  "invalid poolAddress",
		},
		{
			name: "empty mint",
			request: &DepositWithdrawNFTRequest{
				Action:      "DEPOSIT",
				PoolAddress: "11111111111111111111111111111112",
				Mint:        "",
				Blockhash:   "11111111111111111111111111111114",
			},
			wantErr: true,
			errMsg:  "mint is required",
		},
		{
			name: "invalid mint address",
			request: &DepositWithdrawNFTRequest{
				Action:      "DEPOSIT",
				PoolAddress: "11111111111111111111111111111112",
				Mint:        "invalid",
				Blockhash:   "11111111111111111111111111111114",
			},
			wantErr: true,
			errMsg:  "invalid mint address",
		},
		{
			name: "empty blockhash",
			request: &DepositWithdrawNFTRequest{
				Action:      "DEPOSIT",
				PoolAddress: "11111111111111111111111111111112",
				Mint:        "11111111111111111111111111111113",
				Blockhash:   "",
			},
			wantErr: true,
			errMsg:  "blockhash is required",
		},
		{
			name: "invalid nft source",
			request: &DepositWithdrawNFTRequest{
				Action:      "DEPOSIT",
				PoolAddress: "11111111111111111111111111111112",
				Mint:        "11111111111111111111111111111113",
				Blockhash:   "11111111111111111111111111111114",
				NftSource:   stringPtr("invalid"),
			},
			wantErr: true,
			errMsg:  "invalid nftSource address",
		},
		{
			name: "negative compute",
			request: &DepositWithdrawNFTRequest{
				Action:      "DEPOSIT",
				PoolAddress: "11111111111111111111111111111112",
				Mint:        "11111111111111111111111111111113",
				Blockhash:   "11111111111111111111111111111114",
				Compute:     int32Ptr(-1),
			},
			wantErr: true,
			errMsg:  "compute must be >= 0",
		},
		{
			name: "valid request with all optional fields",
			request: &DepositWithdrawNFTRequest{
				Action:                "WITHDRAW",
				PoolAddress:           "11111111111111111111111111111112",
				Mint:                  "11111111111111111111111111111113",
				Blockhash:             "11111111111111111111111111111114",
				Compute:               int32Ptr(200000),
				PriorityMicroLamports: int32Ptr(1000),
				NftSource:             stringPtr("11111111111111111111111111111115"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("DepositWithdrawNFTRequest.Validate() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("DepositWithdrawNFTRequest.Validate() error = %v, want error containing %v", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("DepositWithdrawNFTRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestDepositWithdrawSOLRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request *DepositWithdrawSOLRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid deposit request - uppercase",
			request: &DepositWithdrawSOLRequest{
				Action:      "DEPOSIT",
				PoolAddress: "11111111111111111111111111111112",
				Lamports:    1000000.0, // 1 SOL
				Blockhash:   "11111111111111111111111111111114",
			},
			wantErr: false,
		},
		{
			name: "valid deposit request - lowercase (auto-normalized)",
			request: &DepositWithdrawSOLRequest{
				Action:      "deposit",
				PoolAddress: "11111111111111111111111111111112",
				Lamports:    1000000.0, // 1 SOL
				Blockhash:   "11111111111111111111111111111114",
			},
			wantErr: false,
		},
		{
			name: "valid withdraw request - uppercase",
			request: &DepositWithdrawSOLRequest{
				Action:      "WITHDRAW",
				PoolAddress: "11111111111111111111111111111112",
				Lamports:    500000.0, // 0.5 SOL
				Blockhash:   "11111111111111111111111111111114",
			},
			wantErr: false,
		},
		{
			name: "valid withdraw request - mixed case (auto-normalized)",
			request: &DepositWithdrawSOLRequest{
				Action:      "WithDraw",
				PoolAddress: "11111111111111111111111111111112",
				Lamports:    500000.0, // 0.5 SOL
				Blockhash:   "11111111111111111111111111111114",
			},
			wantErr: false,
		},
		{
			name: "empty action",
			request: &DepositWithdrawSOLRequest{
				Action:      "",
				PoolAddress: "11111111111111111111111111111112",
				Lamports:    1000000.0,
				Blockhash:   "11111111111111111111111111111114",
			},
			wantErr: true,
			errMsg:  "action is required",
		},
		{
			name: "invalid action",
			request: &DepositWithdrawSOLRequest{
				Action:      "INVALID",
				PoolAddress: "11111111111111111111111111111112",
				Lamports:    1000000.0,
				Blockhash:   "11111111111111111111111111111114",
			},
			wantErr: true,
			errMsg:  "invalid action",
		},

		{
			name: "empty pool address",
			request: &DepositWithdrawSOLRequest{
				Action:      "deposit",
				PoolAddress: "",
				Lamports:    1000000.0,
				Blockhash:   "11111111111111111111111111111114",
			},
			wantErr: true,
			errMsg:  "poolAddress is required",
		},
		{
			name: "invalid pool address",
			request: &DepositWithdrawSOLRequest{
				Action:      "deposit",
				PoolAddress: "invalid",
				Lamports:    1000000.0,
				Blockhash:   "11111111111111111111111111111114",
			},
			wantErr: true,
			errMsg:  "invalid poolAddress",
		},
		{
			name: "negative lamports",
			request: &DepositWithdrawSOLRequest{
				Action:      "deposit",
				PoolAddress: "11111111111111111111111111111112",
				Lamports:    -1000.0,
				Blockhash:   "11111111111111111111111111111114",
			},
			wantErr: true,
			errMsg:  "lamports must be >= 0",
		},
		{
			name: "empty blockhash",
			request: &DepositWithdrawSOLRequest{
				Action:      "deposit",
				PoolAddress: "11111111111111111111111111111112",
				Lamports:    1000000.0,
				Blockhash:   "",
			},
			wantErr: true,
			errMsg:  "blockhash is required",
		},
		{
			name: "negative compute",
			request: &DepositWithdrawSOLRequest{
				Action:      "deposit",
				PoolAddress: "11111111111111111111111111111112",
				Lamports:    1000000.0,
				Blockhash:   "11111111111111111111111111111114",
				Compute:     int32Ptr(-1),
			},
			wantErr: true,
			errMsg:  "compute must be >= 0",
		},
		{
			name: "valid request with all optional fields",
			request: &DepositWithdrawSOLRequest{
				Action:                "WITHDRAW",
				PoolAddress:           "11111111111111111111111111111112",
				Lamports:              2000000.0, // 2 SOL
				Blockhash:             "11111111111111111111111111111114",
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
					t.Errorf("DepositWithdrawSOLRequest.Validate() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("DepositWithdrawSOLRequest.Validate() error = %v, want error containing %v", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("DepositWithdrawSOLRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestDepositWithdrawNFTRequest_JSON(t *testing.T) {
	tests := []struct {
		name     string
		request  *DepositWithdrawNFTRequest
		expected string
	}{
		{
			name: "minimal request",
			request: &DepositWithdrawNFTRequest{
				Action:      "DEPOSIT",
				PoolAddress: "11111111111111111111111111111112",
				Mint:        "11111111111111111111111111111113",
				Blockhash:   "11111111111111111111111111111114",
			},
			expected: `{"action":"DEPOSIT","poolAddress":"11111111111111111111111111111112","mint":"11111111111111111111111111111113","blockhash":"11111111111111111111111111111114"}`,
		},
		{
			name: "full request",
			request: &DepositWithdrawNFTRequest{
				Action:                "WITHDRAW",
				PoolAddress:           "11111111111111111111111111111112",
				Mint:                  "11111111111111111111111111111113",
				Blockhash:             "11111111111111111111111111111114",
				Compute:               int32Ptr(200000),
				PriorityMicroLamports: int32Ptr(1000),
				NftSource:             stringPtr("11111111111111111111111111111115"),
			},
			expected: `{"action":"WITHDRAW","poolAddress":"11111111111111111111111111111112","mint":"11111111111111111111111111111113","blockhash":"11111111111111111111111111111114","compute":200000,"priorityMicroLamports":1000,"nftSource":"11111111111111111111111111111115"}`,
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
			var unmarshaled DepositWithdrawNFTRequest
			err = json.Unmarshal(data, &unmarshaled)
			if err != nil {
				t.Errorf("json.Unmarshal() error = %v", err)
				return
			}

			// Compare fields
			if unmarshaled.Action != tt.request.Action {
				t.Errorf("Unmarshaled Action = %v, want %v", unmarshaled.Action, tt.request.Action)
			}
			if unmarshaled.PoolAddress != tt.request.PoolAddress {
				t.Errorf("Unmarshaled PoolAddress = %v, want %v", unmarshaled.PoolAddress, tt.request.PoolAddress)
			}
			if unmarshaled.Mint != tt.request.Mint {
				t.Errorf("Unmarshaled Mint = %v, want %v", unmarshaled.Mint, tt.request.Mint)
			}
		})
	}
}

func TestDepositWithdrawSOLRequest_JSON(t *testing.T) {
	tests := []struct {
		name     string
		request  *DepositWithdrawSOLRequest
		expected string
	}{
		{
			name: "minimal request",
			request: &DepositWithdrawSOLRequest{
				Action:      "DEPOSIT",
				PoolAddress: "11111111111111111111111111111112",
				Lamports:    1000000.0,
				Blockhash:   "11111111111111111111111111111114",
			},
			expected: `{"action":"DEPOSIT","poolAddress":"11111111111111111111111111111112","lamports":1000000,"blockhash":"11111111111111111111111111111114"}`,
		},
		{
			name: "full request",
			request: &DepositWithdrawSOLRequest{
				Action:                "WITHDRAW",
				PoolAddress:           "11111111111111111111111111111112",
				Lamports:              2000000.0,
				Blockhash:             "11111111111111111111111111111114",
				Compute:               int32Ptr(200000),
				PriorityMicroLamports: int32Ptr(1000),
			},
			expected: `{"action":"WITHDRAW","poolAddress":"11111111111111111111111111111112","lamports":2000000,"blockhash":"11111111111111111111111111111114","compute":200000,"priorityMicroLamports":1000}`,
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
			var unmarshaled DepositWithdrawSOLRequest
			err = json.Unmarshal(data, &unmarshaled)
			if err != nil {
				t.Errorf("json.Unmarshal() error = %v", err)
				return
			}

			// Compare fields
			if unmarshaled.Action != tt.request.Action {
				t.Errorf("Unmarshaled Action = %v, want %v", unmarshaled.Action, tt.request.Action)
			}
			if unmarshaled.PoolAddress != tt.request.PoolAddress {
				t.Errorf("Unmarshaled PoolAddress = %v, want %v", unmarshaled.PoolAddress, tt.request.PoolAddress)
			}
			if unmarshaled.Lamports != tt.request.Lamports {
				t.Errorf("Unmarshaled Lamports = %v, want %v", unmarshaled.Lamports, tt.request.Lamports)
			}
		})
	}
}

func TestDepositWithdrawNFTRequest_UnmarshalJSON_TrimWhitespace(t *testing.T) {
	jsonStr := `{"action":"  DEPOSIT  ","poolAddress":"  11111111111111111111111111111112  ","mint":"  11111111111111111111111111111113  ","blockhash":"  11111111111111111111111111111114  ","nftSource":"  11111111111111111111111111111115  "}`

	var request DepositWithdrawNFTRequest
	err := json.Unmarshal([]byte(jsonStr), &request)
	if err != nil {
		t.Errorf("json.Unmarshal() error = %v", err)
		return
	}

	expectedAction := "DEPOSIT"
	expectedPoolAddress := "11111111111111111111111111111112"
	expectedMint := "11111111111111111111111111111113"
	expectedBlockhash := "11111111111111111111111111111114"
	expectedNftSource := "11111111111111111111111111111115"

	if request.Action != expectedAction {
		t.Errorf("Unmarshaled Action = %v, want %v", request.Action, expectedAction)
	}
	if request.PoolAddress != expectedPoolAddress {
		t.Errorf("Unmarshaled PoolAddress = %v, want %v", request.PoolAddress, expectedPoolAddress)
	}
	if request.Mint != expectedMint {
		t.Errorf("Unmarshaled Mint = %v, want %v", request.Mint, expectedMint)
	}
	if request.Blockhash != expectedBlockhash {
		t.Errorf("Unmarshaled Blockhash = %v, want %v", request.Blockhash, expectedBlockhash)
	}
	if request.NftSource == nil || *request.NftSource != expectedNftSource {
		t.Errorf("Unmarshaled NftSource = %v, want %v", request.NftSource, expectedNftSource)
	}
}

func TestDepositWithdrawNFTRequest_UnmarshalJSON_ActionNormalization(t *testing.T) {
	tests := []struct {
		name           string
		jsonStr        string
		expectedAction string
	}{
		{
			name:           "lowercase action normalized to uppercase",
			jsonStr:        `{"action":"deposit","poolAddress":"11111111111111111111111111111112","mint":"11111111111111111111111111111113","blockhash":"11111111111111111111111111111114"}`,
			expectedAction: "DEPOSIT",
		},
		{
			name:           "mixed case action normalized to uppercase",
			jsonStr:        `{"action":"WithDraw","poolAddress":"11111111111111111111111111111112","mint":"11111111111111111111111111111113","blockhash":"11111111111111111111111111111114"}`,
			expectedAction: "WITHDRAW",
		},
		{
			name:           "uppercase action remains uppercase",
			jsonStr:        `{"action":"WITHDRAW","poolAddress":"11111111111111111111111111111112","mint":"11111111111111111111111111111113","blockhash":"11111111111111111111111111111114"}`,
			expectedAction: "WITHDRAW",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var request DepositWithdrawNFTRequest
			err := json.Unmarshal([]byte(tt.jsonStr), &request)
			if err != nil {
				t.Errorf("json.Unmarshal() error = %v", err)
				return
			}

			if request.Action != tt.expectedAction {
				t.Errorf("Unmarshaled Action = %v, want %v", request.Action, tt.expectedAction)
			}
		})
	}
}

func TestDepositWithdrawSOLRequest_UnmarshalJSON_TrimWhitespace(t *testing.T) {
	jsonStr := `{"action":"  deposit  ","poolAddress":"  11111111111111111111111111111112  ","lamports":1000000,"blockhash":"  11111111111111111111111111111114  "}`

	var request DepositWithdrawSOLRequest
	err := json.Unmarshal([]byte(jsonStr), &request)
	if err != nil {
		t.Errorf("json.Unmarshal() error = %v", err)
		return
	}

	expectedAction := "DEPOSIT"
	expectedPoolAddress := "11111111111111111111111111111112"
	expectedBlockhash := "11111111111111111111111111111114"

	if request.Action != expectedAction {
		t.Errorf("Unmarshaled Action = %v, want %v", request.Action, expectedAction)
	}
	if request.PoolAddress != expectedPoolAddress {
		t.Errorf("Unmarshaled PoolAddress = %v, want %v", request.PoolAddress, expectedPoolAddress)
	}
	if request.Blockhash != expectedBlockhash {
		t.Errorf("Unmarshaled Blockhash = %v, want %v", request.Blockhash, expectedBlockhash)
	}
}

func TestDepositWithdrawSOLRequest_UnmarshalJSON_ActionNormalization(t *testing.T) {
	tests := []struct {
		name           string
		jsonStr        string
		expectedAction string
	}{
		{
			name:           "uppercase action remains uppercase",
			jsonStr:        `{"action":"DEPOSIT","poolAddress":"11111111111111111111111111111112","lamports":1000000,"blockhash":"11111111111111111111111111111114"}`,
			expectedAction: "DEPOSIT",
		},
		{
			name:           "mixed case action normalized to uppercase",
			jsonStr:        `{"action":"WithDraw","poolAddress":"11111111111111111111111111111112","lamports":1000000,"blockhash":"11111111111111111111111111111114"}`,
			expectedAction: "WITHDRAW",
		},
		{
			name:           "lowercase action normalized to uppercase",
			jsonStr:        `{"action":"withdraw","poolAddress":"11111111111111111111111111111112","lamports":1000000,"blockhash":"11111111111111111111111111111114"}`,
			expectedAction: "WITHDRAW",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var request DepositWithdrawSOLRequest
			err := json.Unmarshal([]byte(tt.jsonStr), &request)
			if err != nil {
				t.Errorf("json.Unmarshal() error = %v", err)
				return
			}

			if request.Action != tt.expectedAction {
				t.Errorf("Unmarshaled Action = %v, want %v", request.Action, tt.expectedAction)
			}
		})
	}
}

// Helper function for string pointer
func stringPtr(s string) *string {
	return &s
}
