package errors

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name     string
		apiError *APIError
		expected string
	}{
		{
			name: "error without details",
			apiError: &APIError{
				Code:    400,
				Message: "bad request",
			},
			expected: "API error 400: bad request",
		},
		{
			name: "error with details",
			apiError: &APIError{
				Code:    422,
				Message: "validation error",
				Details: "wallet field is required",
			},
			expected: "API error 422: validation error (wallet field is required)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.apiError.Error(); got != tt.expected {
				t.Errorf("APIError.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNetworkError_Error(t *testing.T) {
	err := &NetworkError{
		Op:  "GET",
		Err: fmt.Errorf("connection refused"),
	}

	expected := "network error during GET: connection refused"
	if got := err.Error(); got != expected {
		t.Errorf("NetworkError.Error() = %v, want %v", got, expected)
	}
}

func TestValidationError_Error(t *testing.T) {
	err := &ValidationError{
		Field:   "wallet",
		Message: "invalid format",
	}

	expected := "validation error for field 'wallet': invalid format"
	if got := err.Error(); got != expected {
		t.Errorf("ValidationError.Error() = %v, want %v", got, expected)
	}
}

func TestParseAPIError(t *testing.T) {
	tests := []struct {
		name          string
		statusCode    int
		responseBody  string
		expectedError string
		expectedCode  int
		shouldBeNil   bool
	}{
		{
			name:         "success response",
			statusCode:   200,
			responseBody: `{"message": "success"}`,
			shouldBeNil:  true,
		},
		{
			name:          "JSON error response",
			statusCode:    422,
			responseBody:  `{"message": "validation failed", "details": "wallet is required"}`,
			expectedError: "API error 422: validation failed (wallet is required)",
			expectedCode:  422,
		},
		{
			name:          "JSON error with error field",
			statusCode:    400,
			responseBody:  `{"error": "bad request"}`,
			expectedError: "API error 400: bad request",
			expectedCode:  400,
		},
		{
			name:          "non-JSON error response",
			statusCode:    500,
			responseBody:  "Internal Server Error",
			expectedError: "API error 500: internal server error",
			expectedCode:  500,
		},
		{
			name:          "empty response body",
			statusCode:    404,
			responseBody:  "",
			expectedError: "API error 404: Not Found",
			expectedCode:  404,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{
				StatusCode: tt.statusCode,
				Body:       io.NopCloser(bytes.NewBufferString(tt.responseBody)),
			}

			err := ParseAPIError(resp)

			if tt.shouldBeNil {
				if err != nil {
					t.Errorf("ParseAPIError() expected nil, got %v", err)
				}
				return
			}

			if err == nil {
				t.Errorf("ParseAPIError() expected error, got nil")
				return
			}

			apiErr, ok := err.(*APIError)
			if !ok {
				t.Errorf("ParseAPIError() expected *APIError, got %T", err)
				return
			}

			if apiErr.Code != tt.expectedCode {
				t.Errorf("ParseAPIError() code = %v, want %v", apiErr.Code, tt.expectedCode)
			}

			if apiErr.Error() != tt.expectedError {
				t.Errorf("ParseAPIError() error = %v, want %v", apiErr.Error(), tt.expectedError)
			}
		})
	}
}

func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      *APIError
		wantCode int
	}{
		{"ErrInvalidWallet", ErrInvalidWallet, 400},
		{"ErrValidation", ErrValidation, 422},
		{"ErrUnauthorized", ErrUnauthorized, 401},
		{"ErrRateLimit", ErrRateLimit, 429},
		{"ErrInternalServer", ErrInternalServer, 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Code != tt.wantCode {
				t.Errorf("%s.Code = %v, want %v", tt.name, tt.err.Code, tt.wantCode)
			}
			if tt.err.Message == "" {
				t.Errorf("%s.Message is empty", tt.name)
			}
		})
	}
}
