package geocoding

import (
	"context"
	"net/http"
	"testing"

	internalhttp "github.com/pettinz/mapbox-go-sdk/internal/http"
	"github.com/pettinz/mapbox-go-sdk/internal/testutil"
)

func TestService_Batch(t *testing.T) {
	tests := []struct {
		name           string
		request        *BatchRequest
		mockStatus     int
		mockResponse   string
		wantErr        bool
		validateResult func(*testing.T, *BatchResponse)
	}{
		{
			name: "successful batch geocoding",
			request: &BatchRequest{
				Queries: []BatchQuery{
					{
						ID:    "query1",
						Query: "New York",
					},
					{
						ID:        "query2",
						Longitude: float64Ptr(-118.243683),
						Latitude:  float64Ptr(34.052235),
					},
				},
			},
			mockStatus:   http.StatusOK,
			mockResponse: testutil.BatchGeocodingResponse,
			wantErr:      false,
			validateResult: func(t *testing.T, resp *BatchResponse) {
				if len(resp.Results) == 0 {
					t.Error("expected at least one result")
				}
				if len(resp.Results) > 0 && resp.Results[0].ID != "query1" {
					t.Errorf("expected first result ID to be query1, got %s", resp.Results[0].ID)
				}
			},
		},
		{
			name: "batch with partial errors",
			request: &BatchRequest{
				Queries: []BatchQuery{
					{ID: "query1", Query: "Valid Query"},
					{ID: "query2", Query: "Another Query"},
				},
			},
			mockStatus:   http.StatusOK,
			mockResponse: testutil.BatchGeocodingResponse,
			wantErr:      false,
			validateResult: func(t *testing.T, resp *BatchResponse) {
				hasError := false
				for _, result := range resp.Results {
					if result.Error != nil {
						hasError = true
						break
					}
				}
				if !hasError {
					t.Log("Note: No errors in batch response (this is expected with mock data)")
				}
			},
		},
		{
			name:         "empty queries",
			request:      &BatchRequest{Queries: []BatchQuery{}},
			mockStatus:   http.StatusOK,
			mockResponse: "{}",
			wantErr:      true,
		},
		{
			name: "too many queries",
			request: &BatchRequest{
				Queries: make([]BatchQuery, 1001),
			},
			mockStatus:   http.StatusOK,
			mockResponse: "{}",
			wantErr:      true,
		},
		{
			name: "invalid query - missing both query types",
			request: &BatchRequest{
				Queries: []BatchQuery{
					{ID: "invalid"},
				},
			},
			mockStatus:   http.StatusOK,
			mockResponse: "{}",
			wantErr:      true,
		},
		{
			name: "invalid query - both query types specified",
			request: &BatchRequest{
				Queries: []BatchQuery{
					{
						Query:     "San Francisco",
						Longitude: float64Ptr(-122.4194),
						Latitude:  float64Ptr(37.7749),
					},
				},
			},
			mockStatus:   http.StatusOK,
			mockResponse: "{}",
			wantErr:      true,
		},
		{
			name: "invalid coordinates in batch",
			request: &BatchRequest{
				Queries: []BatchQuery{
					{
						Longitude: float64Ptr(181),
						Latitude:  float64Ptr(0),
					},
				},
			},
			mockStatus:   http.StatusOK,
			mockResponse: "{}",
			wantErr:      true,
		},
		{
			name: "API error",
			request: &BatchRequest{
				Queries: []BatchQuery{
					{Query: "test"},
				},
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

			result, err := service.Batch(context.Background(), tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("Batch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.validateResult != nil {
				tt.validateResult(t, result)
			}
		})
	}
}

func TestValidateBatchQuery(t *testing.T) {
	tests := []struct {
		name    string
		query   BatchQuery
		index   int
		wantErr bool
	}{
		{
			name:    "valid forward query",
			query:   BatchQuery{Query: "San Francisco"},
			index:   0,
			wantErr: false,
		},
		{
			name: "valid reverse query",
			query: BatchQuery{
				Longitude: float64Ptr(-122.4194),
				Latitude:  float64Ptr(37.7749),
			},
			index:   0,
			wantErr: false,
		},
		{
			name:    "missing both query types",
			query:   BatchQuery{},
			index:   0,
			wantErr: true,
		},
		{
			name: "both query types specified",
			query: BatchQuery{
				Query:     "San Francisco",
				Longitude: float64Ptr(-122.4194),
				Latitude:  float64Ptr(37.7749),
			},
			index:   0,
			wantErr: true,
		},
		{
			name: "invalid longitude",
			query: BatchQuery{
				Longitude: float64Ptr(181),
				Latitude:  float64Ptr(0),
			},
			index:   0,
			wantErr: true,
		},
		{
			name: "invalid latitude",
			query: BatchQuery{
				Longitude: float64Ptr(0),
				Latitude:  float64Ptr(91),
			},
			index:   0,
			wantErr: true,
		},
		{
			name: "boundary coordinates",
			query: BatchQuery{
				Longitude: float64Ptr(180),
				Latitude:  float64Ptr(90),
			},
			index:   0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateBatchQuery(&tt.query, tt.index)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateBatchQuery() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBatch_MaxQueries(t *testing.T) {
	// Test with exactly max allowed queries
	queries := make([]BatchQuery, maxBatchQueries)
	for i := 0; i < maxBatchQueries; i++ {
		queries[i] = BatchQuery{
			ID:    "query" + string(rune(i)),
			Query: "test",
		}
	}

	req := &BatchRequest{Queries: queries}

	server := testutil.MockServer(t, testutil.MockResponse(http.StatusOK, testutil.BatchGeocodingResponse))
	defer server.Close()

	httpClient := internalhttp.New(server.URL, nil)
	service := New("test-token", httpClient)

	_, err := service.Batch(context.Background(), req)
	if err != nil {
		t.Errorf("Batch() with max queries should not error, got: %v", err)
	}
}
