package searchbox

import (
	"context"
	"net/http"
	"strings"
	"testing"

	internalhttp "github.com/pettinz/mapbox-go-sdk/internal/http"
	"github.com/pettinz/mapbox-go-sdk/internal/testutil"
)

func TestService_Suggest(t *testing.T) {
	tests := []struct {
		name           string
		request        *SuggestRequest
		mockStatus     int
		mockResponse   string
		wantErr        bool
		validateResult func(*testing.T, *SuggestResponse)
	}{
		{
			name: "successful suggest",
			request: &SuggestRequest{
				Query:        "coffee",
				SessionToken: "550e8400-e29b-41d4-a716-446655440000",
				Proximity:    []float64{-122.4194, 37.7749},
				Limit:        intPtr(5),
			},
			mockStatus:   http.StatusOK,
			mockResponse: testutil.SearchBoxSuggestResponse,
			wantErr:      false,
			validateResult: func(t *testing.T, resp *SuggestResponse) {
				if len(resp.Suggestions) == 0 {
					t.Error("expected at least one suggestion")
				}
				if len(resp.Suggestions) > 0 {
					suggestion := resp.Suggestions[0]
					if suggestion.MapboxID == "" {
						t.Error("expected non-empty mapbox_id")
					}
					if suggestion.Name == "" {
						t.Error("expected non-empty name")
					}
				}
			},
		},
		{
			name: "suggest with all options",
			request: &SuggestRequest{
				Query:        "restaurant",
				SessionToken: "550e8400-e29b-41d4-a716-446655440000",
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
					Profile: "driving",
				},
			},
			mockStatus:   http.StatusOK,
			mockResponse: testutil.SearchBoxSuggestResponse,
			wantErr:      false,
		},
		{
			name: "suggest with IP proximity",
			request: &SuggestRequest{
				Query:        "coffee",
				SessionToken: "550e8400-e29b-41d4-a716-446655440000",
				ProximityIP:  true,
			},
			mockStatus:   http.StatusOK,
			mockResponse: testutil.SearchBoxSuggestResponse,
			wantErr:      false,
		},
		{
			name: "empty query",
			request: &SuggestRequest{
				Query:        "",
				SessionToken: "550e8400-e29b-41d4-a716-446655440000",
			},
			mockStatus:   http.StatusOK,
			mockResponse: "{}",
			wantErr:      true,
		},
		{
			name: "query too long",
			request: &SuggestRequest{
				Query:        strings.Repeat("a", 257),
				SessionToken: "550e8400-e29b-41d4-a716-446655440000",
			},
			mockStatus:   http.StatusOK,
			mockResponse: "{}",
			wantErr:      true,
		},
		{
			name: "missing session token",
			request: &SuggestRequest{
				Query: "coffee",
			},
			mockStatus:   http.StatusOK,
			mockResponse: "{}",
			wantErr:      true,
		},
		{
			name: "limit out of range",
			request: &SuggestRequest{
				Query:        "coffee",
				SessionToken: "550e8400-e29b-41d4-a716-446655440000",
				Limit:        intPtr(11),
			},
			mockStatus:   http.StatusOK,
			mockResponse: "{}",
			wantErr:      true,
		},
		{
			name: "API error",
			request: &SuggestRequest{
				Query:        "test",
				SessionToken: "550e8400-e29b-41d4-a716-446655440000",
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

			result, err := service.Suggest(context.Background(), tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("Suggest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.validateResult != nil {
				tt.validateResult(t, result)
			}
		})
	}
}

func TestService_Retrieve(t *testing.T) {
	tests := []struct {
		name           string
		request        *RetrieveRequest
		mockStatus     int
		mockResponse   string
		wantErr        bool
		validateResult func(*testing.T, *RetrieveResponse)
	}{
		{
			name: "successful retrieve",
			request: &RetrieveRequest{
				MapboxID:     "dXJuOm1ieHBvaTphYmNkZWY",
				SessionToken: "550e8400-e29b-41d4-a716-446655440000",
			},
			mockStatus:   http.StatusOK,
			mockResponse: testutil.SearchBoxRetrieveResponse,
			wantErr:      false,
			validateResult: func(t *testing.T, resp *RetrieveResponse) {
				if resp.Type != "FeatureCollection" {
					t.Errorf("expected Type FeatureCollection, got %s", resp.Type)
				}
				if len(resp.Features) == 0 {
					t.Error("expected at least one feature")
				}
				if len(resp.Features) > 0 {
					feature := resp.Features[0]
					if feature.Properties.MapboxID == "" {
						t.Error("expected non-empty mapbox_id")
					}
					if feature.Geometry.Coordinates[0] == 0 && feature.Geometry.Coordinates[1] == 0 {
						t.Error("expected non-zero coordinates")
					}
				}
			},
		},
		{
			name: "retrieve with navigation",
			request: &RetrieveRequest{
				MapboxID:     "dXJuOm1ieHBvaTphYmNkZWY",
				SessionToken: "550e8400-e29b-41d4-a716-446655440000",
				Navigation: &NavigationOptions{
					ETAType: "navigation",
					Origin:  []float64{-122.4, 37.8},
					Profile: "walking",
				},
			},
			mockStatus:   http.StatusOK,
			mockResponse: testutil.SearchBoxRetrieveResponse,
			wantErr:      false,
		},
		{
			name: "missing mapbox_id",
			request: &RetrieveRequest{
				MapboxID:     "",
				SessionToken: "550e8400-e29b-41d4-a716-446655440000",
			},
			mockStatus:   http.StatusOK,
			mockResponse: "{}",
			wantErr:      true,
		},
		{
			name: "missing session token",
			request: &RetrieveRequest{
				MapboxID: "dXJuOm1ieHBvaTphYmNkZWY",
			},
			mockStatus:   http.StatusOK,
			mockResponse: "{}",
			wantErr:      true,
		},
		{
			name: "API error - not found",
			request: &RetrieveRequest{
				MapboxID:     "invalid-id",
				SessionToken: "550e8400-e29b-41d4-a716-446655440000",
			},
			mockStatus:   http.StatusNotFound,
			mockResponse: testutil.NotFoundErrorResponse,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := testutil.MockServer(t, testutil.MockResponse(tt.mockStatus, tt.mockResponse))
			defer server.Close()

			httpClient := internalhttp.New(server.URL, nil)
			service := New("test-token", httpClient)

			result, err := service.Retrieve(context.Background(), tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("Retrieve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.validateResult != nil {
				tt.validateResult(t, result)
			}
		})
	}
}

func TestBuildSuggestQuery(t *testing.T) {
	service := &Service{token: "test-token"}

	req := &SuggestRequest{
		Query:        "coffee shop",
		SessionToken: "550e8400-e29b-41d4-a716-446655440000",
		Proximity:    []float64{-122.4194, 37.7749},
		BBox:         []float64{-122.5, 37.7, -122.3, 37.8},
		Country:      []string{"US", "CA"},
		Language:     "en",
		Limit:        intPtr(10),
		Types:        []string{"poi"},
		POICategory:  []string{"coffee_shop"},
		Navigation: &NavigationOptions{
			ETAType: "navigation",
			Origin:  []float64{-122.4, 37.8},
			Profile: "driving",
		},
	}

	query := service.buildSuggestQuery(req)

	tests := []struct {
		key      string
		expected string
	}{
		{"access_token", "test-token"},
		{"q", "coffee shop"},
		{"session_token", "550e8400-e29b-41d4-a716-446655440000"},
		{"proximity", "-122.4194,37.7749"},
		{"bbox", "-122.5,37.7,-122.3,37.8"},
		{"country", "US,CA"},
		{"language", "en"},
		{"limit", "10"},
		{"types", "poi"},
		{"poi_category", "coffee_shop"},
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

func TestBuildSuggestQueryWithIPProximity(t *testing.T) {
	service := &Service{token: "test-token"}

	req := &SuggestRequest{
		Query:        "coffee",
		SessionToken: "550e8400-e29b-41d4-a716-446655440000",
		ProximityIP:  true,
	}

	query := service.buildSuggestQuery(req)

	if query.Get("proximity") != "ip" {
		t.Errorf("expected proximity=ip, got %q", query.Get("proximity"))
	}
}

func TestBuildRetrieveQuery(t *testing.T) {
	service := &Service{token: "test-token"}

	req := &RetrieveRequest{
		MapboxID:     "dXJuOm1ieHBvaTphYmNkZWY",
		SessionToken: "550e8400-e29b-41d4-a716-446655440000",
		Navigation: &NavigationOptions{
			ETAType: "navigation",
			Origin:  []float64{-122.4, 37.8},
			Profile: "cycling",
		},
	}

	query := service.buildRetrieveQuery(req)

	tests := []struct {
		key      string
		expected string
	}{
		{"access_token", "test-token"},
		{"session_token", "550e8400-e29b-41d4-a716-446655440000"},
		{"eta_type", "navigation"},
		{"origin", "-122.4,37.8"},
		{"navigation_profile", "cycling"},
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
