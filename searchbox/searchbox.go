package searchbox

import (
	internalhttp "github.com/pettinz/mapbox-go-sdk/internal/http"
)

const (
	// API paths
	suggestPath        = "/search/searchbox/v1/suggest"
	retrievePath       = "/search/searchbox/v1/retrieve"
	forwardPath        = "/search/searchbox/v1/forward"
	categorySearchPath = "/search/searchbox/v1/category"
	listCategoriesPath = "/search/searchbox/v1/category"
	reversePath        = "/search/searchbox/v1/reverse"
)

// Service provides access to the Mapbox Search Box API.
type Service struct {
	token      string
	httpClient *internalhttp.Client
}

// New creates a new Search Box service.
func New(token string, httpClient *internalhttp.Client) *Service {
	return &Service{
		token:      token,
		httpClient: httpClient,
	}
}
