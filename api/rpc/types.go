package rpc

import (
	"encoding/json"
)

// PriorityFeesRequest represents the request for getting priority fees
type PriorityFeesRequest struct {
	// No parameters needed for this endpoint
}

// PriorityFeesResponse represents the response from the priority fees endpoint
type PriorityFeesResponse struct {
	Min      int64 `json:"min"`
	Low      int64 `json:"low"`
	Medium   int64 `json:"medium"`
	High     int64 `json:"high"`
	VeryHigh int64 `json:"veryHigh"`
}

// Validator interface for request validation
type Validator interface {
	Validate() error
}

// Validate validates the PriorityFeesRequest fields
func (r *PriorityFeesRequest) Validate() error {
	// No validation needed for this request as it has no parameters
	return nil
}

// MarshalJSON implements custom JSON marshaling for PriorityFeesRequest
func (r *PriorityFeesRequest) MarshalJSON() ([]byte, error) {
	type Alias PriorityFeesRequest
	return json.Marshal((*Alias)(r))
}

// UnmarshalJSON implements custom JSON unmarshaling for PriorityFeesRequest
func (r *PriorityFeesRequest) UnmarshalJSON(data []byte) error {
	type Alias PriorityFeesRequest
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// No normalization needed for this request
	return nil
}

// MarshalJSON implements custom JSON marshaling for PriorityFeesResponse
func (r *PriorityFeesResponse) MarshalJSON() ([]byte, error) {
	type Alias PriorityFeesResponse
	return json.Marshal((*Alias)(r))
}

// UnmarshalJSON implements custom JSON unmarshaling for PriorityFeesResponse
func (r *PriorityFeesResponse) UnmarshalJSON(data []byte) error {
	type Alias PriorityFeesResponse
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	return json.Unmarshal(data, &aux)
}
