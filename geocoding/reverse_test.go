package geocoding

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
		validateResult func(*testing.T, *Response)
	}{
		{
			name: "successful reverse geocoding",
			request: &ReverseRequest{
				Longitude: -122.419415,
				Latitude:  37.774929,
				Language:  "en",
				Limit:     intPtr(1),
			},
			mockStatus:   http.StatusOK,
			mockResponse: testutil.ReverseGeocodingResponse,
			wantErr:      false,
			validateResult: func(t *testing.T, resp *Response) {
				if resp.Type != "FeatureCollection" {
					t.Errorf("expected Type FeatureCollection, got %s", resp.Type)
				}
				if len(resp.Features) == 0 {
					t.Error("expected at least one feature")
				}
				if len(resp.Features) > 0 {
					feature := resp.Features[0]
					if feature.Properties.PlaceName == "" {
						t.Error("expected non-empty place name")
					}
				}
			},
		},
		{
			name: "invalid longitude",
			request: &ReverseRequest{
				Longitude: 181,
				Latitude:  0,
			},
			mockStatus:   http.StatusOK,
			mockResponse: "{}",
			wantErr:      true,
		},
		{
			name: "invalid latitude",
			request: &ReverseRequest{
				Longitude: 0,
				Latitude:  91,
			},
			mockStatus:   http.StatusOK,
			mockResponse: "{}",
			wantErr:      true,
		},
		{
			name: "API error",
			request: &ReverseRequest{
				Longitude: 0,
				Latitude:  0,
			},
			mockStatus:   http.StatusNotFound,
			mockResponse: testutil.NotFoundErrorResponse,
			wantErr:      true,
		},
		{
			name: "boundary coordinates - max",
			request: &ReverseRequest{
				Longitude: 180,
				Latitude:  90,
			},
			mockStatus:   http.StatusOK,
			mockResponse: testutil.ReverseGeocodingResponse,
			wantErr:      false,
		},
		{
			name: "boundary coordinates - min",
			request: &ReverseRequest{
				Longitude: -180,
				Latitude:  -90,
			},
			mockStatus:   http.StatusOK,
			mockResponse: testutil.ReverseGeocodingResponse,
			wantErr:      false,
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
		Country:   []string{"US"},
		Language:  "en",
		Limit:     intPtr(5),
		Types:     []string{"place", "region"},
		Worldview: "US",
	}

	query := service.buildReverseQuery(req)

	tests := []struct {
		key      string
		expected string
	}{
		{"access_token", "test-token"},
		{"longitude", "-122.419415"},
		{"latitude", "37.774929"},
		{"country", "US"},
		{"language", "en"},
		{"limit", "5"},
		{"types", "place,region"},
		{"worldview", "US"},
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

func TestBuildReverseQuery_MinimalRequest(t *testing.T) {
	service := &Service{token: "test-token"}

	req := &ReverseRequest{
		Longitude: 0,
		Latitude:  0,
	}

	query := service.buildReverseQuery(req)

	// Check required parameters
	if query.Get("access_token") != "test-token" {
		t.Error("expected access_token to be set")
	}
	if query.Get("longitude") != "0" {
		t.Error("expected longitude to be set")
	}
	if query.Get("latitude") != "0" {
		t.Error("expected latitude to be set")
	}

	// Check optional parameters are not set
	if query.Get("country") != "" {
		t.Error("expected country to be empty")
	}
	if query.Get("language") != "" {
		t.Error("expected language to be empty")
	}
	if query.Get("limit") != "" {
		t.Error("expected limit to be empty")
	}
}
