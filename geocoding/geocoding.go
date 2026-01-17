package geocoding

import (
	internalhttp "github.com/pettinz/mapbox-go-sdk/internal/http"
)

const (
	// API paths
	forwardPath = "/search/geocode/v6/forward"
	reversePath = "/search/geocode/v6/reverse"
	batchPath   = "/search/geocode/v6/batch"
)

// Service provides access to the Mapbox Geocoding API.
type Service struct {
	token      string
	httpClient *internalhttp.Client
}

// New creates a new Geocoding service.
func New(token string, httpClient *internalhttp.Client) *Service {
	return &Service{
		token:      token,
		httpClient: httpClient,
	}
}
