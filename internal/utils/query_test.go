package utils

import (
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/srpvpn/tensor-go-sdk/internal/errors"
)

func TestBuildQueryParams(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected url.Values
		wantErr  bool
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: url.Values{},
			wantErr:  false,
		},
		{
			name: "simple struct",
			input: struct {
				Wallet string `json:"wallet"`
				Count  int    `json:"count"`
			}{
				Wallet: "test-wallet",
				Count:  5,
			},
			expected: url.Values{
				"wallet": []string{"test-wallet"},
				"count":  []string{"5"},
			},
			wantErr: false,
		},
		{
			name: "struct with omitempty",
			input: struct {
				Wallet   string  `json:"wallet"`
				Optional *string `json:"optional,omitempty"`
				Count    int     `json:"count,omitempty"`
			}{
				Wallet: "test-wallet",
				Count:  0, // should be omitted
			},
			expected: url.Values{
				"wallet": []string{"test-wallet"},
			},
			wantErr: false,
		},
		{
			name: "struct with pointer fields",
			input: struct {
				Wallet   string `json:"wallet"`
				Optional *bool  `json:"optional,omitempty"`
			}{
				Wallet:   "test-wallet",
				Optional: boolPtr(true),
			},
			expected: url.Values{
				"wallet":   []string{"test-wallet"},
				"optional": []string{"true"},
			},
			wantErr: false,
		},
		{
			name: "struct with slice",
			input: struct {
				Wallet     string   `json:"wallet"`
				Currencies []string `json:"currencies,omitempty"`
			}{
				Wallet:     "test-wallet",
				Currencies: []string{"SOL", "USDC"},
			},
			expected: url.Values{
				"wallet":     []string{"test-wallet"},
				"currencies": []string{"SOL,USDC"},
			},
			wantErr: false,
		},
		{
			name: "struct with ignored fields",
			input: struct {
				Wallet   string `json:"wallet"`
				Ignored  string `json:"-"`
				NoTag    string
				Internal string `json:""`
			}{
				Wallet:   "test-wallet",
				Ignored:  "should-be-ignored",
				NoTag:    "no-tag",
				Internal: "internal",
			},
			expected: url.Values{
				"wallet": []string{"test-wallet"},
			},
			wantErr: false,
		},
		{
			name:     "non-struct input",
			input:    "not-a-struct",
			expected: nil,
			wantErr:  true,
		},
		{
			name: "pointer to struct",
			input: &struct {
				Wallet string `json:"wallet"`
			}{
				Wallet: "test-wallet",
			},
			expected: url.Values{
				"wallet": []string{"test-wallet"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BuildQueryParams(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildQueryParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("BuildQueryParams() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestValidateWalletAddress(t *testing.T) {
	tests := []struct {
		name    string
		wallet  string
		wantErr bool
		errType string
	}{
		{
			name:    "valid wallet address",
			wallet:  "11111111111111111111111111111112", // Valid base58 address
			wantErr: false,
		},
		{
			name:    "empty wallet",
			wallet:  "",
			wantErr: true,
			errType: "required",
		},
		{
			name:    "too short wallet",
			wallet:  "short",
			wantErr: true,
			errType: "length",
		},
		{
			name:    "too long wallet",
			wallet:  "this-wallet-address-is-way-too-long-to-be-valid-solana-address",
			wantErr: true,
			errType: "length",
		},
		{
			name:    "invalid characters",
			wallet:  "11111111111111111111111111111110", // Contains '0' which is not valid base58
			wantErr: true,
			errType: "characters",
		},
		{
			name:    "another valid address",
			wallet:  "So11111111111111111111111111111111111111112", // Another valid format
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateWalletAddress(tt.wallet)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateWalletAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				validationErr, ok := err.(*errors.ValidationError)
				if !ok {
					t.Errorf("ValidateWalletAddress() expected ValidationError, got %T", err)
					return
				}

				if validationErr.Field != "wallet" {
					t.Errorf("ValidateWalletAddress() field = %v, want wallet", validationErr.Field)
				}

				// Check error message contains expected type
				if tt.errType != "" && !strings.Contains(validationErr.Message, tt.errType) {
					t.Errorf("ValidateWalletAddress() message = %v, want to contain %v", validationErr.Message, tt.errType)
				}
			}
		})
	}
}

func TestIsZeroValue(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{"zero bool", false, true},
		{"non-zero bool", true, false},
		{"zero int", 0, true},
		{"non-zero int", 42, false},
		{"zero string", "", true},
		{"non-zero string", "hello", false},
		{"nil pointer", (*string)(nil), true},
		{"non-nil pointer", stringPtr("hello"), false},
		{"empty slice", []string{}, true},
		{"non-empty slice", []string{"hello"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := reflect.ValueOf(tt.value)
			if got := isZeroValue(v); got != tt.expected {
				t.Errorf("isZeroValue() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFieldToString(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected string
		wantErr  bool
	}{
		{"string", "hello", "hello", false},
		{"bool true", true, "true", false},
		{"bool false", false, "false", false},
		{"int", 42, "42", false},
		{"float", 3.14, "3.14", false},
		{"string slice", []string{"a", "b", "c"}, "a,b,c", false},
		{"empty string slice", []string{}, "", false},
		{"nil pointer", (*string)(nil), "", false},
		{"non-nil pointer", stringPtr("hello"), "hello", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := reflect.ValueOf(tt.value)
			got, err := fieldToString(v)
			if (err != nil) != tt.wantErr {
				t.Errorf("fieldToString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("fieldToString() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// Helper functions for tests
func boolPtr(b bool) *bool {
	return &b
}

func stringPtr(s string) *string {
	return &s
}
