package geocoding

import (
	"context"
	"net/http"
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
		validateResult func(*testing.T, *Response)
	}{
		{
			name: "successful forward geocoding",
			request: &ForwardRequest{
				Query:    "1600 Pennsylvania Avenue NW",
				Language: "en",
				Limit:    intPtr(5),
			},
			mockStatus:   http.StatusOK,
			mockResponse: testutil.ForwardGeocodingResponse,
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
					if feature.Geometry.Coordinates[0] == 0 && feature.Geometry.Coordinates[1] == 0 {
						t.Error("expected non-zero coordinates")
					}
				}
			},
		},
		{
			name:         "empty query",
			request:      &ForwardRequest{Query: ""},
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

func TestService_ForwardStructured(t *testing.T) {
	tests := []struct {
		name           string
		request        *StructuredForwardRequest
		mockStatus     int
		mockResponse   string
		wantErr        bool
		validateResult func(*testing.T, *Response)
	}{
		{
			name: "successful structured forward geocoding",
			request: &StructuredForwardRequest{
				AddressNumber: "1600",
				Street:        "Pennsylvania Avenue NW",
				Place:         "Washington",
				Region:        "DC",
				Country:       "US",
				Postcode:      "20500",
			},
			mockStatus:   http.StatusOK,
			mockResponse: testutil.ForwardGeocodingResponse,
			wantErr:      false,
			validateResult: func(t *testing.T, resp *Response) {
				if resp.Type != "FeatureCollection" {
					t.Errorf("expected Type FeatureCollection, got %s", resp.Type)
				}
				if len(resp.Features) == 0 {
					t.Error("expected at least one feature")
				}
			},
		},
		{
			name:         "empty address components",
			request:      &StructuredForwardRequest{},
			mockStatus:   http.StatusOK,
			mockResponse: "{}",
			wantErr:      true,
		},
		{
			name: "partial address components",
			request: &StructuredForwardRequest{
				Place:   "San Francisco",
				Country: "US",
			},
			mockStatus:   http.StatusOK,
			mockResponse: testutil.ForwardGeocodingResponse,
			wantErr:      false,
		},
		{
			name: "API error",
			request: &StructuredForwardRequest{
				Place: "test",
			},
			mockStatus:   http.StatusUnprocessableEntity,
			mockResponse: testutil.ValidationErrorResponse,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := testutil.MockServer(t, testutil.MockResponse(tt.mockStatus, tt.mockResponse))
			defer server.Close()

			httpClient := internalhttp.New(server.URL, nil)
			service := New("test-token", httpClient)

			result, err := service.ForwardStructured(context.Background(), tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("ForwardStructured() error = %v, wantErr %v", err, tt.wantErr)
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
		Query:        "San Francisco",
		Autocomplete: boolPtr(true),
		BBox:         []float64{-122.5, 37.7, -122.3, 37.8},
		Country:      []string{"US", "CA"},
		Language:     "en",
		Limit:        intPtr(10),
		Proximity:    []float64{-122.4194, 37.7749},
		Types:        []string{"place", "region"},
		Worldview:    "US",
	}

	query := service.buildForwardQuery(req)

	tests := []struct {
		key      string
		expected string
	}{
		{"access_token", "test-token"},
		{"q", "San Francisco"},
		{"autocomplete", "true"},
		{"bbox", "-122.5,37.7,-122.3,37.8"},
		{"country", "US,CA"},
		{"language", "en"},
		{"limit", "10"},
		{"proximity", "-122.4194,37.7749"},
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

func TestBuildStructuredForwardQuery(t *testing.T) {
	service := &Service{token: "test-token"}

	req := &StructuredForwardRequest{
		AddressNumber: "1600",
		Street:        "Pennsylvania Avenue NW",
		Block:         "Block A",
		Place:         "Washington",
		Region:        "DC",
		Postcode:      "20500",
		Country:       "US",
		Autocomplete:  boolPtr(false),
		BBox:          []float64{-77.1, 38.8, -77.0, 38.9},
		Language:      "en",
		Limit:         intPtr(5),
		Proximity:     []float64{-77.0365, 38.8977},
		Worldview:     "US",
	}

	query := service.buildStructuredForwardQuery(req)

	tests := []struct {
		key      string
		expected string
	}{
		{"access_token", "test-token"},
		{"address_number", "1600"},
		{"street", "Pennsylvania Avenue NW"},
		{"block", "Block A"},
		{"place", "Washington"},
		{"region", "DC"},
		{"postcode", "20500"},
		{"country", "US"},
		{"autocomplete", "false"},
		{"bbox", "-77.1,38.8,-77,38.9"},
		{"language", "en"},
		{"limit", "5"},
		{"proximity", "-77.0365,38.8977"},
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
