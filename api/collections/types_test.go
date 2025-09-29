package collections

import (
	"strings"
	"testing"
)

func TestGetVerifiedCollectionsRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     *GetVerifiedCollectionsRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request with all required fields",
			req: &GetVerifiedCollectionsRequest{
				SortBy: "statsV2.volume1h:desc",
				Limit:  50,
			},
			wantErr: false,
		},
		{
			name: "valid request with all fields",
			req: &GetVerifiedCollectionsRequest{
				SortBy:       "statsV2.volume24h:desc",
				Limit:        100,
				SlugDisplays: []string{"portalsuniverse", "degods"},
				CollIds:      []string{"coll-id-1", "coll-id-2"},
				Vocs:         []string{"voc1", "voc2"},
				Fvcs:         []string{"fvc1", "fvc2"},
				Page:         func() *int32 { v := int32(1); return &v }(),
			},
			wantErr: false,
		},
		{
			name: "missing sortBy",
			req: &GetVerifiedCollectionsRequest{
				Limit: 50,
			},
			wantErr: true,
			errMsg:  "sortBy is required",
		},
		{
			name: "invalid sortBy format - no direction",
			req: &GetVerifiedCollectionsRequest{
				SortBy: "statsV2.volume1h",
				Limit:  50,
			},
			wantErr: true,
			errMsg:  "sortBy must include direction",
		},
		{
			name: "limit is zero",
			req: &GetVerifiedCollectionsRequest{
				SortBy: "statsV2.volume1h:desc",
				Limit:  0,
			},
			wantErr: true,
			errMsg:  "limit must be greater than 0",
		},
		{
			name: "limit is negative",
			req: &GetVerifiedCollectionsRequest{
				SortBy: "statsV2.volume1h:desc",
				Limit:  -10,
			},
			wantErr: true,
			errMsg:  "limit must be greater than 0",
		},
		{
			name: "limit exceeds maximum",
			req: &GetVerifiedCollectionsRequest{
				SortBy: "statsV2.volume1h:desc",
				Limit:  101,
			},
			wantErr: true,
			errMsg:  "limit must be 100 or less",
		},
		{
			name: "too many vocs",
			req: &GetVerifiedCollectionsRequest{
				SortBy: "statsV2.volume1h:desc",
				Limit:  50,
				Vocs:   []string{"v1", "v2", "v3", "v4", "v5", "v6", "v7", "v8", "v9", "v10", "v11"},
			},
			wantErr: true,
			errMsg:  "maximum 10 vocs allowed",
		},
		{
			name: "too many fvcs",
			req: &GetVerifiedCollectionsRequest{
				SortBy: "statsV2.volume1h:desc",
				Limit:  50,
				Fvcs:   []string{"f1", "f2", "f3", "f4", "f5", "f6", "f7", "f8", "f9", "f10", "f11"},
			},
			wantErr: true,
			errMsg:  "maximum 10 fvcs allowed",
		},
		{
			name: "page is zero",
			req: &GetVerifiedCollectionsRequest{
				SortBy: "statsV2.volume1h:desc",
				Limit:  50,
				Page:   func() *int32 { v := int32(0); return &v }(),
			},
			wantErr: true,
			errMsg:  "page must be 1 or greater",
		},
		{
			name: "valid request with exactly 10 vocs",
			req: &GetVerifiedCollectionsRequest{
				SortBy: "statsV2.volume1h:desc",
				Limit:  50,
				Vocs:   []string{"v1", "v2", "v3", "v4", "v5", "v6", "v7", "v8", "v9", "v10"},
			},
			wantErr: false,
		},
		{
			name: "valid request with exactly 10 fvcs",
			req: &GetVerifiedCollectionsRequest{
				SortBy: "statsV2.volume1h:desc",
				Limit:  50,
				Fvcs:   []string{"f1", "f2", "f3", "f4", "f5", "f6", "f7", "f8", "f9", "f10"},
			},
			wantErr: false,
		},
		{
			name: "valid request with page 1",
			req: &GetVerifiedCollectionsRequest{
				SortBy: "statsV2.volume1h:desc",
				Limit:  50,
				Page:   func() *int32 { v := int32(1); return &v }(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
					return
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("expected error to contain %q, got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}
