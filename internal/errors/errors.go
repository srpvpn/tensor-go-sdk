package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// APIError represents an error returned by the Tensor API
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("API error %d: %s (%s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("API error %d: %s", e.Code, e.Message)
}

// Predefined API errors
var (
	ErrInvalidWallet  = &APIError{Code: 400, Message: "invalid wallet address"}
	ErrValidation     = &APIError{Code: 422, Message: "validation error"}
	ErrUnauthorized   = &APIError{Code: 401, Message: "unauthorized"}
	ErrRateLimit      = &APIError{Code: 429, Message: "rate limit exceeded"}
	ErrInternalServer = &APIError{Code: 500, Message: "internal server error"}
)

// NetworkError represents a network-related error
type NetworkError struct {
	Op  string
	Err error
}

// Error implements the error interface
func (e *NetworkError) Error() string {
	return fmt.Sprintf("network error during %s: %v", e.Op, e.Err)
}

// Unwrap returns the underlying error
func (e *NetworkError) Unwrap() error {
	return e.Err
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error for field '%s': %s", e.Field, e.Message)
}

// ParseAPIError parses an HTTP response and returns an appropriate error
func ParseAPIError(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	apiErr := &APIError{
		Code: resp.StatusCode,
	}

	// Try to parse JSON error response
	var errorResponse struct {
		Message string `json:"message"`
		Details string `json:"details,omitempty"`
		Error   string `json:"error,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err == nil {
		if errorResponse.Message != "" {
			apiErr.Message = errorResponse.Message
		} else if errorResponse.Error != "" {
			apiErr.Message = errorResponse.Error
		}
		apiErr.Details = errorResponse.Details
	}

	// Set default messages if not provided
	if apiErr.Message == "" {
		switch resp.StatusCode {
		case 400:
			apiErr.Message = "bad request"
		case 401:
			apiErr.Message = "unauthorized"
		case 422:
			apiErr.Message = "validation error"
		case 429:
			apiErr.Message = "rate limit exceeded"
		case 500:
			apiErr.Message = "internal server error"
		default:
			apiErr.Message = http.StatusText(resp.StatusCode)
		}
	}

	return apiErr
}