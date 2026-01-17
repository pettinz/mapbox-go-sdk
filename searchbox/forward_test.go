package searchbox

import (
	"context"
	"net/http"
	"strings"
	"testing"

	internalhttp "github.com/pettinz/mapbox-go-sdk/internal/http"
	"github.com/pettinz/mapbox-go-sdk/internal/testutil"
)

func TestService_Forward(t *testing.T) {
	tests := []struct {
		name           string
		request        *ForwardRequest
		mockStatus     int
		mockResponse   string
		wantErr        bool
		validateResult func(*testing.T, *ForwardResponse)
	}{
		{
			name: "successful forward search",
			request: &ForwardRequest{
				Query:    "Colosseum Rome",
				Language: "en",
				Limit:    intPtr(5),
			},
			mockStatus:   http.StatusOK,
			mockResponse: testutil.SearchBoxForwardResponse,
			wantErr:      false,
			validateResult: func(t *testing.T, resp *ForwardResponse) {
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
			name: "forward search with all options",
			request: &ForwardRequest{
				Query:        "restaurant",
				Autocomplete: boolPtr(true),
				Proximity:    []float64{-122.4194, 37.7749},
				BBox:         []float64{-122.5, 37.7, -122.3, 37.8},
				Country:      []string{"US"},
				Language:     "en",
				Limit:        intPtr(10),
				Types:        []string{"poi"},
				POICategory:  []string{"restaurant"},
				Navigation: &NavigationOptions{
					ETAType: "navigation",
					Origin:  []float64{-122.4, 37.8},
					Profile: "walking",
				},
			},
			mockStatus:   http.StatusOK,
			mockResponse: testutil.SearchBoxForwardResponse,
			wantErr:      false,
		},
		{
			name: "forward search with IP proximity",
			request: &ForwardRequest{
				Query:       "coffee",
				ProximityIP: true,
			},
			mockStatus:   http.StatusOK,
			mockResponse: testutil.SearchBoxForwardResponse,
			wantErr:      false,
		},
		{
			name: "forward search with autocomplete disabled",
			request: &ForwardRequest{
				Query:        "exact address",
				Autocomplete: boolPtr(false),
			},
			mockStatus:   http.StatusOK,
			mockResponse: testutil.SearchBoxForwardResponse,
			wantErr:      false,
		},
		{
			name: "empty query",
			request: &ForwardRequest{
				Query: "",
			},
			mockStatus:   http.StatusOK,
			mockResponse: "{}",
			wantErr:      true,
		},
		{
			name: "query too long",
			request: &ForwardRequest{
				Query: strings.Repeat("a", 257),
			},
			mockStatus:   http.StatusOK,
			mockResponse: "{}",
			wantErr:      true,
		},
		{
			name: "limit out of range",
			request: &ForwardRequest{
				Query: "coffee",
				Limit: intPtr(11),
			},
			mockStatus:   http.StatusOK,
			mockResponse: "{}",
			wantErr:      true,
		},
		{
			name: "API error",
			request: &ForwardRequest{
				Query: "test",
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

			result, err := service.Forward(context.Background(), tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("Forward() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.validateResult != nil {
				tt.validateResult(t, result)
			}
		})
	}
}

func TestBuildForwardQuery(t *testing.T) {
	service := &Service{token: "test-token"}

	req := &ForwardRequest{
		Query:        "restaurant",
		Autocomplete: boolPtr(false),
		Proximity:    []float64{-122.4194, 37.7749},
		BBox:         []float64{-122.5, 37.7, -122.3, 37.8},
		Country:      []string{"US", "CA"},
		Language:     "en",
		Limit:        intPtr(10),
		Types:        []string{"poi"},
		POICategory:  []string{"restaurant", "cafe"},
		Navigation: &NavigationOptions{
			ETAType: "navigation",
			Origin:  []float64{-122.4, 37.8},
			Profile: "driving",
		},
	}

	query := service.buildForwardQuery(req)

	tests := []struct {
		key      string
		expected string
	}{
		{"access_token", "test-token"},
		{"q", "restaurant"},
		{"autocomplete", "false"},
		{"proximity", "-122.4194,37.7749"},
		{"bbox", "-122.5,37.7,-122.3,37.8"},
		{"country", "US,CA"},
		{"language", "en"},
		{"limit", "10"},
		{"types", "poi"},
		{"poi_category", "restaurant,cafe"},
		{"eta_type", "navigation"},
		{"origin", "-122.4,37.8"},
		{"navigation_profile", "driving"},
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

func TestBuildForwardQueryWithIPProximity(t *testing.T) {
	service := &Service{token: "test-token"}

	req := &ForwardRequest{
		Query:       "coffee",
		ProximityIP: true,
	}

	query := service.buildForwardQuery(req)

	if query.Get("proximity") != "ip" {
		t.Errorf("expected proximity=ip, got %q", query.Get("proximity"))
	}
}
