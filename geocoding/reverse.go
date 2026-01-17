package geocoding

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// Reverse performs reverse geocoding.
// It converts geographic coordinates into a human-readable address.
func (s *Service) Reverse(ctx context.Context, req *ReverseRequest) (*Response, error) {
	// Validate coordinates
	if err := validateCoordinates(req.Longitude, req.Latitude); err != nil {
		return nil, err
	}

	query := s.buildReverseQuery(req)

	var result Response
	if err := s.httpClient.Get(ctx, reversePath, query, &result); err != nil {
		return nil, fmt.Errorf("reverse geocoding failed: %w", err)
	}

	return &result, nil
}

// buildReverseQuery builds query parameters for reverse geocoding.
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

	if req.Worldview != "" {
		q.Set("worldview", req.Worldview)
	}

	return q
}

// validateCoordinates validates longitude and latitude values.
func validateCoordinates(longitude, latitude float64) error {
	if longitude < -180 || longitude > 180 {
		return fmt.Errorf("longitude must be between -180 and 180, got %f", longitude)
	}
	if latitude < -90 || latitude > 90 {
		return fmt.Errorf("latitude must be between -90 and 90, got %f", latitude)
	}
	return nil
}
