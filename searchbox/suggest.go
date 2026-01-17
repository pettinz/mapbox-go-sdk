package searchbox

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// Suggest performs an autocomplete search and returns suggestions without coordinates.
// This is the first step in the Suggest/Retrieve workflow for interactive autocomplete.
// Use the same session_token when calling Retrieve to complete the workflow.
func (s *Service) Suggest(ctx context.Context, req *SuggestRequest) (*SuggestResponse, error) {
	if err := validateSuggestRequest(req); err != nil {
		return nil, err
	}

	query := s.buildSuggestQuery(req)

	var result SuggestResponse
	if err := s.httpClient.Get(ctx, suggestPath, query, &result); err != nil {
		return nil, fmt.Errorf("suggest search failed: %w", err)
	}

	return &result, nil
}

// Retrieve gets the full feature details (including coordinates) for a selected suggestion.
// This is the second step in the Suggest/Retrieve workflow, called after the user selects
// a suggestion. Use the same session_token from the Suggest request.
func (s *Service) Retrieve(ctx context.Context, req *RetrieveRequest) (*RetrieveResponse, error) {
	if err := validateRetrieveRequest(req); err != nil {
		return nil, err
	}

	// Build path with mapbox_id
	path := fmt.Sprintf("%s/%s", retrievePath, req.MapboxID)
	query := s.buildRetrieveQuery(req)

	var result RetrieveResponse
	if err := s.httpClient.Get(ctx, path, query, &result); err != nil {
		return nil, fmt.Errorf("retrieve feature failed: %w", err)
	}

	return &result, nil
}

// validateSuggestRequest validates the Suggest request parameters.
func validateSuggestRequest(req *SuggestRequest) error {
	if req.Query == "" {
		return fmt.Errorf("query is required")
	}

	if len(req.Query) > 256 {
		return fmt.Errorf("query exceeds maximum length of 256 characters")
	}

	if req.SessionToken == "" {
		return fmt.Errorf("session_token is required")
	}

	if req.Limit != nil && (*req.Limit < 1 || *req.Limit > 10) {
		return fmt.Errorf("limit must be between 1 and 10")
	}

	return nil
}

// validateRetrieveRequest validates the Retrieve request parameters.
func validateRetrieveRequest(req *RetrieveRequest) error {
	if req.MapboxID == "" {
		return fmt.Errorf("mapbox_id is required")
	}

	if req.SessionToken == "" {
		return fmt.Errorf("session_token is required")
	}

	return nil
}

// buildSuggestQuery builds query parameters for the Suggest endpoint.
func (s *Service) buildSuggestQuery(req *SuggestRequest) url.Values {
	q := url.Values{}
	q.Set("access_token", s.token)
	q.Set("q", req.Query)
	q.Set("session_token", req.SessionToken)

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

// buildRetrieveQuery builds query parameters for the Retrieve endpoint.
func (s *Service) buildRetrieveQuery(req *RetrieveRequest) url.Values {
	q := url.Values{}
	q.Set("access_token", s.token)
	q.Set("session_token", req.SessionToken)

	// Navigation options
	if req.Navigation != nil {
		addNavigationParams(q, req.Navigation)
	}

	return q
}

// addNavigationParams adds navigation-related query parameters.
func addNavigationParams(q url.Values, nav *NavigationOptions) {
	if nav.ETAType != "" {
		q.Set("eta_type", nav.ETAType)
	}

	if len(nav.Origin) == 2 {
		q.Set("origin", formatFloatArray(nav.Origin))
	}

	if nav.Profile != "" {
		q.Set("navigation_profile", nav.Profile)
	}
}

// formatFloatArray formats a float array as a comma-separated string.
func formatFloatArray(arr []float64) string {
	parts := make([]string, len(arr))
	for i, v := range arr {
		parts[i] = strconv.FormatFloat(v, 'f', -1, 64)
	}
	return strings.Join(parts, ",")
}
