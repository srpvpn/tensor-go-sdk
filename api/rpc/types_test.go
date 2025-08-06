package rpc

import (
	"encoding/json"
	"testing"
)

func TestPriorityFeesRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request *PriorityFeesRequest
		wantErr bool
	}{
		{
			name:    "valid request",
			request: &PriorityFeesRequest{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("PriorityFeesRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPriorityFeesRequest_JSON(t *testing.T) {
	tests := []struct {
		name     string
		request  *PriorityFeesRequest
		expected string
	}{
		{
			name:     "empty request",
			request:  &PriorityFeesRequest{},
			expected: `{}`,
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
			var unmarshaled PriorityFeesRequest
			err = json.Unmarshal(jsonBytes, &unmarshaled)
			if err != nil {
				t.Errorf("json.Unmarshal() error = %v", err)
			}
		})
	}
}

func TestPriorityFeesResponse_JSON(t *testing.T) {
	tests := []struct {
		name     string
		response *PriorityFeesResponse
		expected string
	}{
		{
			name: "complete response",
			response: &PriorityFeesResponse{
				Min:      0,
				Low:      1000,
				Medium:   5000,
				High:     10000,
				VeryHigh: 50000,
			},
			expected: `{"min":0,"low":1000,"medium":5000,"high":10000,"veryHigh":50000}`,
		},
		{
			name: "zero values response",
			response: &PriorityFeesResponse{
				Min:      0,
				Low:      0,
				Medium:   0,
				High:     0,
				VeryHigh: 0,
			},
			expected: `{"min":0,"low":0,"medium":0,"high":0,"veryHigh":0}`,
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
			var unmarshaled PriorityFeesResponse
			err = json.Unmarshal(jsonBytes, &unmarshaled)
			if err != nil {
				t.Errorf("json.Unmarshal() error = %v", err)
				return
			}

			if unmarshaled != *tt.response {
				t.Errorf("Unmarshaled response = %+v, want %+v", unmarshaled, *tt.response)
			}
		})
	}
}

func TestPriorityFeesResponse_UnmarshalJSON_FromAPIExample(t *testing.T) {
	// Test with the exact API response format from the specification
	apiResponse := `{"min": 0,"low": 0,"medium": 0,"high": 0,"veryHigh": 0}`

	var response PriorityFeesResponse
	err := json.Unmarshal([]byte(apiResponse), &response)
	if err != nil {
		t.Errorf("json.Unmarshal() error = %v", err)
		return
	}

	expected := PriorityFeesResponse{
		Min:      0,
		Low:      0,
		Medium:   0,
		High:     0,
		VeryHigh: 0,
	}

	if response != expected {
		t.Errorf("Unmarshaled response = %+v, want %+v", response, expected)
	}
}
