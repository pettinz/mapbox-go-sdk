package searchbox

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// Reverse performs reverse geocoding, converting coordinates into place names.
func (s *Service) Reverse(ctx context.Context, req *ReverseRequest) (*ReverseResponse, error) {
	if err := validateReverseRequest(req); err != nil {
		return nil, err
	}

	query := s.buildReverseQuery(req)

	var result ReverseResponse
	if err := s.httpClient.Get(ctx, reversePath, query, &result); err != nil {
		return nil, fmt.Errorf("reverse geocoding failed: %w", err)
	}

	return &result, nil
}

// validateReverseRequest validates the Reverse request parameters.
func validateReverseRequest(req *ReverseRequest) error {
	if req.Longitude < -180 || req.Longitude > 180 {
		return fmt.Errorf("longitude must be between -180 and 180")
	}

	if req.Latitude < -90 || req.Latitude > 90 {
		return fmt.Errorf("latitude must be between -90 and 90")
	}

	if req.Limit != nil && (*req.Limit < 1 || *req.Limit > 10) {
		return fmt.Errorf("limit must be between 1 and 10")
	}

	return nil
}

// buildReverseQuery builds query parameters for the Reverse endpoint.
func (s *Service) buildReverseQuery(req *ReverseRequest) url.Values {
	q := url.Values{}
	q.Set("access_token", s.token)
	q.Set("longitude", strconv.FormatFloat(req.Longitude, 'f', -1, 64))
	q.Set("latitude", strconv.FormatFloat(req.Latitude, 'f', -1, 64))

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

	return q
}
