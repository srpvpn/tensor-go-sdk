package escrow

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestDepositWithdrawEscrowRequest_Validate(t *testing.T) {
	// Helper function for int32 pointer
	int32Ptr := func(i int32) *int32 {
		return &i
	}

	tests := []struct {
		name    string
		request *DepositWithdrawEscrowRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid deposit request - uppercase",
			request: &DepositWithdrawEscrowRequest{
				Action:    "DEPOSIT",
				Owner:     "11111111111111111111111111111112",
				Lamports:  1000000000.0, // 1 SOL
				Blockhash: "11111111111111111111111111111114",
			},
			wantErr: false,
		},
		{
			name: "valid deposit request - lowercase (auto-normalized)",
			request: &DepositWithdrawEscrowRequest{
				Action:    "deposit",
				Owner:     "11111111111111111111111111111112",
				Lamports:  1000000000.0, // 1 SOL
				Blockhash: "11111111111111111111111111111114",
			},
			wantErr: false,
		},
		{
			name: "valid withdraw request - uppercase",
			request: &DepositWithdrawEscrowRequest{
				Action:    "WITHDRAW",
				Owner:     "11111111111111111111111111111112",
				Lamports:  500000000.0, // 0.5 SOL
				Blockhash: "11111111111111111111111111111114",
			},
			wantErr: false,
		},
		{
			name: "valid withdraw request - mixed case (auto-normalized)",
			request: &DepositWithdrawEscrowRequest{
				Action:    "WithDraw",
				Owner:     "11111111111111111111111111111112",
				Lamports:  500000000.0, // 0.5 SOL
				Blockhash: "11111111111111111111111111111114",
			},
			wantErr: false,
		},
		{
			name: "empty action",
			request: &DepositWithdrawEscrowRequest{
				Action:    "",
				Owner:     "11111111111111111111111111111112",
				Lamports:  1000000000.0,
				Blockhash: "11111111111111111111111111111114",
			},
			wantErr: true,
			errMsg:  "action is required",
		},
		{
			name: "invalid action",
			request: &DepositWithdrawEscrowRequest{
				Action:    "invalid",
				Owner:     "11111111111111111111111111111112",
				Lamports:  1000000000.0,
				Blockhash: "11111111111111111111111111111114",
			},
			wantErr: true,
			errMsg:  "invalid action",
		},
		{
			name: "empty owner",
			request: &DepositWithdrawEscrowRequest{
				Action:    "DEPOSIT",
				Owner:     "",
				Lamports:  1000000000.0,
				Blockhash: "11111111111111111111111111111114",
			},
			wantErr: true,
			errMsg:  "owner is required",
		},
		{
			name: "invalid owner address",
			request: &DepositWithdrawEscrowRequest{
				Action:    "DEPOSIT",
				Owner:     "invalid-address",
				Lamports:  1000000000.0,
				Blockhash: "11111111111111111111111111111114",
			},
			wantErr: true,
			errMsg:  "invalid owner address",
		},
		{
			name: "negative lamports",
			request: &DepositWithdrawEscrowRequest{
				Action:    "DEPOSIT",
				Owner:     "11111111111111111111111111111112",
				Lamports:  -1000.0,
				Blockhash: "11111111111111111111111111111114",
			},
			wantErr: true,
			errMsg:  "lamports must be >= 0",
		},
		{
			name: "empty blockhash",
			request: &DepositWithdrawEscrowRequest{
				Action:    "DEPOSIT",
				Owner:     "11111111111111111111111111111112",
				Lamports:  1000000000.0,
				Blockhash: "",
			},
			wantErr: true,
			errMsg:  "blockhash is required",
		},
		{
			name: "negative compute",
			request: &DepositWithdrawEscrowRequest{
				Action:    "DEPOSIT",
				Owner:     "11111111111111111111111111111112",
				Lamports:  1000000000.0,
				Blockhash: "11111111111111111111111111111114",
				Compute:   int32Ptr(-1),
			},
			wantErr: true,
			errMsg:  "compute must be >= 0",
		},
		{
			name: "valid request with all optional fields",
			request: &DepositWithdrawEscrowRequest{
				Action:                "WITHDRAW",
				Owner:                 "11111111111111111111111111111112",
				Lamports:              2000000000.0, // 2 SOL
				Blockhash:             "11111111111111111111111111111114",
				Compute:               int32Ptr(200000),
				PriorityMicroLamports: int32Ptr(1000),
			},
			wantErr: false,
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
					t.Errorf("DepositWithdrawEscrowRequest.Validate() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("DepositWithdrawEscrowRequest.Validate() error = %v, want error containing %v", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("DepositWithdrawEscrowRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestDepositWithdrawEscrowRequest_JSON(t *testing.T) {
	// Helper function for int32 pointer
	int32Ptr := func(i int32) *int32 {
		return &i
	}

	tests := []struct {
		name     string
		request  *DepositWithdrawEscrowRequest
		expected string
	}{
		{
			name: "minimal request",
			request: &DepositWithdrawEscrowRequest{
				Action:    "DEPOSIT",
				Owner:     "11111111111111111111111111111112",
				Lamports:  1000000000.0,
				Blockhash: "11111111111111111111111111111114",
			},
			expected: `{"action":"DEPOSIT","owner":"11111111111111111111111111111112","lamports":1000000000,"blockhash":"11111111111111111111111111111114"}`,
		},
		{
			name: "full request",
			request: &DepositWithdrawEscrowRequest{
				Action:                "WITHDRAW",
				Owner:                 "11111111111111111111111111111112",
				Lamports:              2000000000.0,
				Blockhash:             "11111111111111111111111111111114",
				Compute:               int32Ptr(200000),
				PriorityMicroLamports: int32Ptr(1000),
			},
			expected: `{"action":"WITHDRAW","owner":"11111111111111111111111111111112","lamports":2000000000,"blockhash":"11111111111111111111111111111114","compute":200000,"priorityMicroLamports":1000}`,
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
			var unmarshaled DepositWithdrawEscrowRequest
			err = json.Unmarshal(jsonBytes, &unmarshaled)
			if err != nil {
				t.Errorf("json.Unmarshal() error = %v", err)
			}
		})
	}
}

func TestDepositWithdrawEscrowRequest_UnmarshalJSON_TrimWhitespace(t *testing.T) {
	jsonStr := `{"action":"  DEPOSIT  ","owner":"  11111111111111111111111111111112  ","lamports":1000000,"blockhash":"  11111111111111111111111111111114  "}`

	var request DepositWithdrawEscrowRequest
	err := json.Unmarshal([]byte(jsonStr), &request)
	if err != nil {
		t.Errorf("json.Unmarshal() error = %v", err)
		return
	}

	expectedAction := "DEPOSIT"
	expectedOwner := "11111111111111111111111111111112"
	expectedBlockhash := "11111111111111111111111111111114"

	if request.Action != expectedAction {
		t.Errorf("Unmarshaled Action = %v, want %v", request.Action, expectedAction)
	}
	if request.Owner != expectedOwner {
		t.Errorf("Unmarshaled Owner = %v, want %v", request.Owner, expectedOwner)
	}
	if request.Blockhash != expectedBlockhash {
		t.Errorf("Unmarshaled Blockhash = %v, want %v", request.Blockhash, expectedBlockhash)
	}
}

func TestDepositWithdrawEscrowRequest_UnmarshalJSON_ActionNormalization(t *testing.T) {
	tests := []struct {
		name           string
		jsonStr        string
		expectedAction string
	}{
		{
			name:           "uppercase action remains uppercase",
			jsonStr:        `{"action":"DEPOSIT","owner":"11111111111111111111111111111112","lamports":1000000,"blockhash":"11111111111111111111111111111114"}`,
			expectedAction: "DEPOSIT",
		},
		{
			name:           "mixed case action normalized to uppercase",
			jsonStr:        `{"action":"WithDraw","owner":"11111111111111111111111111111112","lamports":1000000,"blockhash":"11111111111111111111111111111114"}`,
			expectedAction: "WITHDRAW",
		},
		{
			name:           "lowercase action normalized to uppercase",
			jsonStr:        `{"action":"withdraw","owner":"11111111111111111111111111111112","lamports":1000000,"blockhash":"11111111111111111111111111111114"}`,
			expectedAction: "WITHDRAW",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var request DepositWithdrawEscrowRequest
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

func TestDepositWithdrawEscrowResponse_JSON(t *testing.T) {
	tests := []struct {
		name     string
		response *DepositWithdrawEscrowResponse
		expected string
	}{
		{
			name: "success response",
			response: &DepositWithdrawEscrowResponse{
				Status: "Ok",
			},
			expected: `{"status":"Ok"}`,
		},
		{
			name: "empty status response",
			response: &DepositWithdrawEscrowResponse{
				Status: "",
			},
			expected: `{"status":""}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBytes, err := json.Marshal(tt.response)
			if err != nil {
				t.Errorf("json.Marshal() error = %v", err)
				return
			}

			if string(jsonBytes) != tt.expected {
				t.Errorf("json.Marshal() = %v, want %v", string(jsonBytes), tt.expected)
			}

			// Test unmarshaling
			var unmarshaled DepositWithdrawEscrowResponse
			err = json.Unmarshal(jsonBytes, &unmarshaled)
			if err != nil {
				t.Errorf("json.Unmarshal() error = %v", err)
			}
		})
	}
}
func TestDepositWithdrawEscrowResponse_UnmarshalJSON_FromAPIExample(t *testing.T) {
	// Test with the exact API response format from the specification
	apiResponse := `{"status":"Ok"}`

	var response DepositWithdrawEscrowResponse
	err := json.Unmarshal([]byte(apiResponse), &response)
	if err != nil {
		t.Errorf("json.Unmarshal() error = %v", err)
		return
	}

	expected := DepositWithdrawEscrowResponse{
		Status: "Ok",
	}

	if response != expected {
		t.Errorf("Unmarshaled response = %+v, want %+v", response, expected)
	}
}
