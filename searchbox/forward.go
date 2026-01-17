package searchbox

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// Forward performs a one-off text search with immediate results including coordinates.
// Unlike Suggest/Retrieve, this does not require a session token and returns full
// GeoJSON features in a single request. Use this for simple search without autocomplete.
func (s *Service) Forward(ctx context.Context, req *ForwardRequest) (*ForwardResponse, error) {
	if err := validateForwardRequest(req); err != nil {
		return nil, err
	}

	query := s.buildForwardQuery(req)

	var result ForwardResponse
	if err := s.httpClient.Get(ctx, forwardPath, query, &result); err != nil {
		return nil, fmt.Errorf("forward search failed: %w", err)
	}

	return &result, nil
}

// validateForwardRequest validates the Forward request parameters.
func validateForwardRequest(req *ForwardRequest) error {
	if req.Query == "" {
		return fmt.Errorf("query is required")
	}

	if len(req.Query) > 256 {
		return fmt.Errorf("query exceeds maximum length of 256 characters")
	}

	if req.Limit != nil && (*req.Limit < 1 || *req.Limit > 10) {
		return fmt.Errorf("limit must be between 1 and 10")
	}

	return nil
}

// buildForwardQuery builds query parameters for the Forward endpoint.
func (s *Service) buildForwardQuery(req *ForwardRequest) url.Values {
	q := url.Values{}
	q.Set("access_token", s.token)
	q.Set("q", req.Query)

	if req.Autocomplete != nil {
		q.Set("autocomplete", strconv.FormatBool(*req.Autocomplete))
	}

	// Proximity (can be coordinates or "ip")
	if req.ProximityIP {
		q.Set("proximity", "ip")
	} else if len(req.Proximity) == 2 {
		q.Set("proximity", formatFloatArray(req.Proximity))
	}

	if len(req.BBox) == 4 {
		q.Set("bbox", formatFloatArray(req.BBox))
	}

	if len(req.Country) > 0 {
		q.Set("country", strings.Join(req.Country, ","))
	}

	if req.Language != "" {
		q.Set("language", req.Language)
	}

	if req.Limit != nil {
		q.Set("limit", strconv.Itoa(*req.Limit))
	}

	if len(req.Types) > 0 {
		q.Set("types", strings.Join(req.Types, ","))
	}

	if len(req.POICategory) > 0 {
		q.Set("poi_category", strings.Join(req.POICategory, ","))
	}

	// Navigation options
	if req.Navigation != nil {
		addNavigationParams(q, req.Navigation)
	}

	return q
}
