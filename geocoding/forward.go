package geocoding

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// Forward performs forward geocoding using text-based search.
// It converts a search query (like "1600 Pennsylvania Avenue NW") into geographic coordinates.
func (s *Service) Forward(ctx context.Context, req *ForwardRequest) (*Response, error) {
	if req.Query == "" {
		return nil, fmt.Errorf("query is required")
	}

	query := s.buildForwardQuery(req)

	var result Response
	if err := s.httpClient.Get(ctx, forwardPath, query, &result); err != nil {
		return nil, fmt.Errorf("forward geocoding failed: %w", err)
	}

	return &result, nil
}

// ForwardStructured performs forward geocoding using structured address components.
// It converts address components (street, city, etc.) into geographic coordinates.
func (s *Service) ForwardStructured(ctx context.Context, req *StructuredForwardRequest) (*Response, error) {
	// At least one address component is required
	if req.AddressNumber == "" && req.Street == "" && req.Block == "" &&
	   req.Place == "" && req.Region == "" && req.Postcode == "" && req.Country == "" {
		return nil, fmt.Errorf("at least one address component is required")
	}

	query := s.buildStructuredForwardQuery(req)

	var result Response
	if err := s.httpClient.Get(ctx, forwardPath, query, &result); err != nil {
		return nil, fmt.Errorf("structured forward geocoding failed: %w", err)
	}

	return &result, nil
}

// buildForwardQuery builds query parameters for forward geocoding.
func (s *Service) buildForwardQuery(req *ForwardRequest) url.Values {
	q := url.Values{}
	q.Set("access_token", s.token)
	q.Set("q", req.Query)

	if req.Autocomplete != nil {
		q.Set("autocomplete", strconv.FormatBool(*req.Autocomplete))
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

	if len(req.Proximity) == 2 {
		q.Set("proximity", formatFloatArray(req.Proximity))
	}

	if len(req.Types) > 0 {
		q.Set("types", strings.Join(req.Types, ","))
	}

	if req.Worldview != "" {
		q.Set("worldview", req.Worldview)
	}

	return q
}

// buildStructuredForwardQuery builds query parameters for structured forward geocoding.
func (s *Service) buildStructuredForwardQuery(req *StructuredForwardRequest) url.Values {
	q := url.Values{}
	q.Set("access_token", s.token)

	if req.AddressNumber != "" {
		q.Set("address_number", req.AddressNumber)
	}

	if req.Street != "" {
		q.Set("street", req.Street)
	}

	if req.Block != "" {
		q.Set("block", req.Block)
	}

	if req.Place != "" {
		q.Set("place", req.Place)
	}

	if req.Region != "" {
		q.Set("region", req.Region)
	}

	if req.Postcode != "" {
		q.Set("postcode", req.Postcode)
	}

	if req.Country != "" {
		q.Set("country", req.Country)
	}

	if req.Autocomplete != nil {
		q.Set("autocomplete", strconv.FormatBool(*req.Autocomplete))
	}

	if len(req.BBox) == 4 {
		q.Set("bbox", formatFloatArray(req.BBox))
	}

	if req.Language != "" {
		q.Set("language", req.Language)
	}

	if req.Limit != nil {
		q.Set("limit", strconv.Itoa(*req.Limit))
	}

	if len(req.Proximity) == 2 {
		q.Set("proximity", formatFloatArray(req.Proximity))
	}

	if req.Worldview != "" {
		q.Set("worldview", req.Worldview)
	}

	return q
}

// formatFloatArray formats a float array as a comma-separated string.
func formatFloatArray(arr []float64) string {
	parts := make([]string, len(arr))
	for i, v := range arr {
		parts[i] = strconv.FormatFloat(v, 'f', -1, 64)
	}
	return strings.Join(parts, ",")
}
