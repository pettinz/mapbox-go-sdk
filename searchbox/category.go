package searchbox

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// CategorySearch searches for POIs in a specific category.
// Requires either proximity, bbox, or SAR to define the search area.
func (s *Service) CategorySearch(ctx context.Context, req *CategorySearchRequest) (*CategorySearchResponse, error) {
	if err := validateCategorySearchRequest(req); err != nil {
		return nil, err
	}

	// Build path with category ID
	path := fmt.Sprintf("%s/%s", categorySearchPath, req.CategoryID)
	query := s.buildCategorySearchQuery(req)

	var result CategorySearchResponse
	if err := s.httpClient.Get(ctx, path, query, &result); err != nil {
		return nil, fmt.Errorf("category search failed: %w", err)
	}

	return &result, nil
}

// ListCategories retrieves all available POI categories.
func (s *Service) ListCategories(ctx context.Context, req *ListCategoriesRequest) (*ListCategoriesResponse, error) {
	query := s.buildListCategoriesQuery(req)

	var result ListCategoriesResponse
	if err := s.httpClient.Get(ctx, listCategoriesPath, query, &result); err != nil {
		return nil, fmt.Errorf("list categories failed: %w", err)
	}

	return &result, nil
}

// validateCategorySearchRequest validates the CategorySearch request parameters.
func validateCategorySearchRequest(req *CategorySearchRequest) error {
	if req.CategoryID == "" {
		return fmt.Errorf("category_id is required")
	}

	// Require either proximity, bbox, or SAR
	hasProximity := req.ProximityIP || len(req.Proximity) == 2
	hasBBox := len(req.BBox) == 4
	hasSAR := req.SAR != nil && len(req.SAR.Route) > 0

	if !hasProximity && !hasBBox && !hasSAR {
		return fmt.Errorf("either proximity, bbox, or SAR is required")
	}

	if req.Limit != nil && (*req.Limit < 1 || *req.Limit > 25) {
		return fmt.Errorf("limit must be between 1 and 25")
	}

	return nil
}

// buildCategorySearchQuery builds query parameters for the CategorySearch endpoint.
func (s *Service) buildCategorySearchQuery(req *CategorySearchRequest) url.Values {
	q := url.Values{}
	q.Set("access_token", s.token)

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

	// Navigation options
	if req.Navigation != nil {
		addNavigationParams(q, req.Navigation)
	}

	// Search Along Route (SAR) options
	if req.SAR != nil {
		addSARParams(q, req.SAR)
	}

	return q
}

// buildListCategoriesQuery builds query parameters for the ListCategories endpoint.
func (s *Service) buildListCategoriesQuery(req *ListCategoriesRequest) url.Values {
	q := url.Values{}
	q.Set("access_token", s.token)

	if req.Language != "" {
		q.Set("language", req.Language)
	}

	return q
}

// addSARParams adds Search Along Route (SAR) query parameters.
func addSARParams(q url.Values, sar *SAROptions) {
	if sar.Type != "" {
		q.Set("sar_type", sar.Type)
	}

	if len(sar.Route) > 0 {
		q.Set("route", encodeRoute(sar.Route))
	}

	if sar.TimeDeviation != nil {
		q.Set("time_deviation", strconv.Itoa(*sar.TimeDeviation))
	}
}

// encodeRoute encodes a route as a semicolon-separated list of coordinate pairs.
// Format: "lon1,lat1;lon2,lat2;..."
func encodeRoute(route [][]float64) string {
	parts := make([]string, len(route))
	for i, coord := range route {
		if len(coord) >= 2 {
			parts[i] = fmt.Sprintf("%f,%f", coord[0], coord[1])
		}
	}
	return strings.Join(parts, ";")
}
