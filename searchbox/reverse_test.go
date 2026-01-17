package searchbox

import (
	"context"
	"net/http"
	"testing"

	internalhttp "github.com/pettinz/mapbox-go-sdk/internal/http"
	"github.com/pettinz/mapbox-go-sdk/internal/testutil"
)

func TestService_Reverse(t *testing.T) {
	tests := []struct {
		name           string
		request        *ReverseRequest
		mockStatus     int
		mockResponse   string
		wantErr        bool
		validateResult func(*testing.T, *ReverseResponse)
	}{
		{
			name: "successful reverse geocoding",
			request: &ReverseRequest{
				Longitude: -122.419415,
				Latitude:  37.774929,
			},
			mockStatus:   http.StatusOK,
			mockResponse: testutil.SearchBoxReverseResponse,
			wantErr:      false,
			validateResult: func(t *testing.T, resp *ReverseResponse) {
				if resp.Type != "FeatureCollection" {
					t.Errorf("expected Type FeatureCollection, got %s", resp.Type)
				}
				if len(resp.Features) == 0 {
					t.Error("expected at least one feature")
				}
				if len(resp.Features) > 0 {
					feature := resp.Features[0]
					if feature.Properties.Name == "" {
						t.Error("expected non-empty name")
					}
					if feature.Geometry.Coordinates[0] == 0 && feature.Geometry.Coordinates[1] == 0 {
						t.Error("expected non-zero coordinates")
					}
				}
			},
		},
		{
			name: "reverse with all options",
			request: &ReverseRequest{
				Longitude: -122.419415,
				Latitude:  37.774929,
				Country:   []string{"US"},
				Language:  "en",
				Limit:     intPtr(5),
				Types:     []string{"address", "street"},
			},
			mockStatus:   http.StatusOK,
			mockResponse: testutil.SearchBoxReverseResponse,
			wantErr:      false,
		},
		{
			name: "longitude out of range - too low",
			request: &ReverseRequest{
				Longitude: -181.0,
				Latitude:  37.774929,
			},
			mockStatus:   http.StatusOK,
			mockResponse: "{}",
			wantErr:      true,
		},
		{
			name: "longitude out of range - too high",
			request: &ReverseRequest{
				Longitude: 181.0,
				Latitude:  37.774929,
			},
			mockStatus:   http.StatusOK,
			mockResponse: "{}",
			wantErr:      true,
		},
		{
			name: "latitude out of range - too low",
			request: &ReverseRequest{
				Longitude: -122.419415,
				Latitude:  -91.0,
			},
			mockStatus:   http.StatusOK,
			mockResponse: "{}",
			wantErr:      true,
		},
		{
			name: "latitude out of range - too high",
			request: &ReverseRequest{
				Longitude: -122.419415,
				Latitude:  91.0,
			},
			mockStatus:   http.StatusOK,
			mockResponse: "{}",
			wantErr:      true,
		},
		{
			name: "limit out of range",
			request: &ReverseRequest{
				Longitude: -122.419415,
				Latitude:  37.774929,
				Limit:     intPtr(11),
			},
			mockStatus:   http.StatusOK,
			mockResponse: "{}",
			wantErr:      true,
		},
		{
			name: "API error",
			request: &ReverseRequest{
				Longitude: -122.419415,
				Latitude:  37.774929,
			},
			mockStatus:   http.StatusUnauthorized,
			mockResponse: testutil.ErrorResponse,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := testutil.MockServer(t, testutil.MockResponse(tt.mockStatus, tt.mockResponse))
			defer server.Close()

			httpClient := internalhttp.New(server.URL, nil)
			service := New("test-token", httpClient)

			result, err := service.Reverse(context.Background(), tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("Reverse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.validateResult != nil {
				tt.validateResult(t, result)
			}
		})
	}
}

func TestBuildReverseQuery(t *testing.T) {
	service := &Service{token: "test-token"}

	req := &ReverseRequest{
		Longitude: -122.419415,
		Latitude:  37.774929,
		Country:   []string{"US", "CA"},
		Language:  "en",
		Limit:     intPtr(10),
		Types:     []string{"address", "street", "place"},
	}

	query := service.buildReverseQuery(req)

	tests := []struct {
		key      string
		expected string
	}{
		{"access_token", "test-token"},
		{"longitude", "-122.419415"},
		{"latitude", "37.774929"},
		{"country", "US,CA"},
		{"language", "en"},
		{"limit", "10"},
		{"types", "address,street,place"},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			actual := query.Get(tt.key)
			if actual != tt.expected {
				t.Errorf("expected %s=%q, got %q", tt.key, tt.expected, actual)
			}
		})
	}
}
