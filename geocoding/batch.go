package geocoding

import (
	"context"
	"fmt"
	"net/url"
)

const (
	maxBatchQueries = 1000
)

// Batch performs batch geocoding for multiple queries in a single request.
// It supports up to 1000 queries per batch and can include both forward and reverse geocoding queries.
func (s *Service) Batch(ctx context.Context, req *BatchRequest) (*BatchResponse, error) {
	// Validate batch size
	if len(req.Queries) == 0 {
		return nil, fmt.Errorf("at least one query is required")
	}
	if len(req.Queries) > maxBatchQueries {
		return nil, fmt.Errorf("maximum %d queries allowed, got %d", maxBatchQueries, len(req.Queries))
	}

	// Validate each query
	for i, query := range req.Queries {
		if err := validateBatchQuery(&query, i); err != nil {
			return nil, err
		}
	}

	// Add access token to query
	query := url.Values{}
	query.Set("access_token", s.token)

	// Build request body
	body := map[string]any{
		"queries": req.Queries,
	}

	var result BatchResponse
	if err := s.httpClient.Post(ctx, batchPath, body, &result); err != nil {
		return nil, fmt.Errorf("batch geocoding failed: %w", err)
	}

	return &result, nil
}

// validateBatchQuery validates a single batch query.
func validateBatchQuery(query *BatchQuery, index int) error {
	// Determine if this is a forward or reverse query
	isForward := query.Query != ""
	isReverse := query.Longitude != nil && query.Latitude != nil

	// Must be either forward or reverse, not both
	if !isForward && !isReverse {
		return fmt.Errorf("query at index %d: must specify either 'q' (forward) or 'longitude'+'latitude' (reverse)", index)
	}
	if isForward && isReverse {
		return fmt.Errorf("query at index %d: cannot specify both 'q' and 'longitude'+'latitude'", index)
	}

	// Validate reverse query coordinates
	if isReverse {
		if err := validateCoordinates(*query.Longitude, *query.Latitude); err != nil {
			return fmt.Errorf("query at index %d: %w", index, err)
		}
	}

	return nil
}
