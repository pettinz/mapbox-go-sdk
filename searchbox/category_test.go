package searchbox

import (
	"context"
	"net/http"
	"testing"

	internalhttp "github.com/pettinz/mapbox-go-sdk/internal/http"
	"github.com/pettinz/mapbox-go-sdk/internal/testutil"
)

func TestService_CategorySearch(t *testing.T) {
	tests := []struct {
		name           string
		request        *CategorySearchRequest
		mockStatus     int
		mockResponse   string
		wantErr        bool
		validateResult func(*testing.T, *CategorySearchResponse)
	}{
		{
			name: "successful category search with proximity",
			request: &CategorySearchRequest{
				CategoryID: "restaurant",
				Proximity:  []float64{-122.4194, 37.7749},
				Limit:      intPtr(10),
			},
			mockStatus:   http.StatusOK,
			mockResponse: testutil.SearchBoxCategorySearchResponse,
			wantErr:      false,
			validateResult: func(t *testing.T, resp *CategorySearchResponse) {
				if resp.Type != "FeatureCollection" {
					t.Errorf("expected Type FeatureCollection, got %s", resp.Type)
				}
				if len(resp.Features) == 0 {
					t.Error("expected at least one feature")
				}
			},
		},
		{
			name: "successful category search with bbox",
			request: &CategorySearchRequest{
				CategoryID: "coffee_shop",
				BBox:       []float64{-122.5, 37.7, -122.3, 37.8},
				Limit:      intPtr(25),
			},
			mockStatus:   http.StatusOK,
			mockResponse: testutil.SearchBoxCategorySearchResponse,
			wantErr:      false,
		},
		{
			name: "category search with IP proximity",
			request: &CategorySearchRequest{
				CategoryID:  "restaurant",
				ProximityIP: true,
			},
			mockStatus:   http.StatusOK,
			mockResponse: testutil.SearchBoxCategorySearchResponse,
			wantErr:      false,
		},
		{
			name: "category search with navigation",
			request: &CategorySearchRequest{
				CategoryID: "gas_station",
				Proximity:  []float64{-122.4194, 37.7749},
				Navigation: &NavigationOptions{
					ETAType: "navigation",
					Origin:  []float64{-122.4, 37.8},
					Profile: "driving",
				},
			},
			mockStatus:   http.StatusOK,
			mockResponse: testutil.SearchBoxCategorySearchResponse,
			wantErr:      false,
		},
		{
			name: "category search with SAR",
			request: &CategorySearchRequest{
				CategoryID: "restaurant",
				SAR: &SAROptions{
					Type: "isochrone",
					Route: [][]float64{
						{-122.4, 37.8},
						{-122.5, 37.7},
					},
					TimeDeviation: intPtr(300),
				},
			},
			mockStatus:   http.StatusOK,
			mockResponse: testutil.SearchBoxCategorySearchResponse,
			wantErr:      false,
		},
		{
			name: "missing category ID",
			request: &CategorySearchRequest{
				Proximity: []float64{-122.4194, 37.7749},
			},
			mockStatus:   http.StatusOK,
			mockResponse: "{}",
			wantErr:      true,
		},
		{
			name: "missing proximity and bbox",
			request: &CategorySearchRequest{
				CategoryID: "restaurant",
			},
			mockStatus:   http.StatusOK,
			mockResponse: "{}",
			wantErr:      true,
		},
		{
			name: "limit out of range",
			request: &CategorySearchRequest{
				CategoryID: "restaurant",
				Proximity:  []float64{-122.4194, 37.7749},
				Limit:      intPtr(26),
			},
			mockStatus:   http.StatusOK,
			mockResponse: "{}",
			wantErr:      true,
		},
		{
			name: "API error",
			request: &CategorySearchRequest{
				CategoryID: "invalid",
				Proximity:  []float64{-122.4194, 37.7749},
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

			result, err := service.CategorySearch(context.Background(), tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("CategorySearch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.validateResult != nil {
				tt.validateResult(t, result)
			}
		})
	}
}

func TestService_ListCategories(t *testing.T) {
	tests := []struct {
		name           string
		request        *ListCategoriesRequest
		mockStatus     int
		mockResponse   string
		wantErr        bool
		validateResult func(*testing.T, *ListCategoriesResponse)
	}{
		{
			name:         "successful list categories",
			request:      &ListCategoriesRequest{},
			mockStatus:   http.StatusOK,
			mockResponse: testutil.SearchBoxListCategoriesResponse,
			wantErr:      false,
			validateResult: func(t *testing.T, resp *ListCategoriesResponse) {
				if len(resp.Categories) == 0 {
					t.Error("expected at least one category")
				}
				if len(resp.Categories) > 0 {
					category := resp.Categories[0]
					if category.CanonicalID == "" {
						t.Error("expected non-empty canonical_name")
					}
					if category.Name == "" {
						t.Error("expected non-empty name")
					}
				}
			},
		},
		{
			name: "list categories with language",
			request: &ListCategoriesRequest{
				Language: "en",
			},
			mockStatus:   http.StatusOK,
			mockResponse: testutil.SearchBoxListCategoriesResponse,
			wantErr:      false,
		},
		{
			name: "API error",
			request: &ListCategoriesRequest{
				Language: "invalid",
			},
			mockStatus:   http.StatusBadRequest,
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

			result, err := service.ListCategories(context.Background(), tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListCategories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.validateResult != nil {
				tt.validateResult(t, result)
			}
		})
	}
}

func TestBuildCategorySearchQuery(t *testing.T) {
	service := &Service{token: "test-token"}

	req := &CategorySearchRequest{
		CategoryID: "restaurant",
		Proximity:  []float64{-122.4194, 37.7749},
		BBox:       []float64{-122.5, 37.7, -122.3, 37.8},
		Country:    []string{"US"},
		Language:   "en",
		Limit:      intPtr(20),
		Navigation: &NavigationOptions{
			ETAType: "navigation",
			Origin:  []float64{-122.4, 37.8},
			Profile: "walking",
		},
		SAR: &SAROptions{
			Type: "isochrone",
			Route: [][]float64{
				{-122.4, 37.8},
				{-122.5, 37.7},
			},
			TimeDeviation: intPtr(300),
		},
	}

	query := service.buildCategorySearchQuery(req)

	tests := []struct {
		key      string
		expected string
	}{
		{"access_token", "test-token"},
		{"proximity", "-122.4194,37.7749"},
		{"bbox", "-122.5,37.7,-122.3,37.8"},
		{"country", "US"},
		{"language", "en"},
		{"limit", "20"},
		{"eta_type", "navigation"},
		{"origin", "-122.4,37.8"},
		{"navigation_profile", "walking"},
		{"sar_type", "isochrone"},
		{"route", "-122.400000,37.800000;-122.500000,37.700000"},
		{"time_deviation", "300"},
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

func TestEncodeRoute(t *testing.T) {
	route := [][]float64{
		{-122.4194, 37.7749},
		{-122.5, 37.7},
		{-122.3, 37.8},
	}

	encoded := encodeRoute(route)
	expected := "-122.419400,37.774900;-122.500000,37.700000;-122.300000,37.800000"

	if encoded != expected {
		t.Errorf("expected %q, got %q", expected, encoded)
	}
}

func TestBuildListCategoriesQuery(t *testing.T) {
	service := &Service{token: "test-token"}

	req := &ListCategoriesRequest{
		Language: "en",
	}

	query := service.buildListCategoriesQuery(req)

	tests := []struct {
		key      string
		expected string
	}{
		{"access_token", "test-token"},
		{"language", "en"},
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
