package mapbox

import (
	"net/http"

	"github.com/pettinz/mapbox-go-sdk/geocoding"
	internalhttp "github.com/pettinz/mapbox-go-sdk/internal/http"
	"github.com/pettinz/mapbox-go-sdk/searchbox"
)

const (
	// DefaultBaseURL is the default base URL for the Mapbox API.
	DefaultBaseURL = "https://api.mapbox.com"
)

// Client is the root client for the Mapbox API.
type Client struct {
	token      string
	baseURL    string
	httpClient *http.Client
	http       *internalhttp.Client
}

// NewClient creates a new Mapbox API client with the given access token.
// Optional functional options can be provided to customize the client.
func NewClient(token string, opts ...Option) *Client {
	if token == "" {
		panic("mapbox: access token is required")
	}

	c := &Client{
		token:      token,
		baseURL:    DefaultBaseURL,
		httpClient: http.DefaultClient,
	}

	// Apply options
	for _, opt := range opts {
		opt(c)
	}

	// Create internal HTTP client
	c.http = internalhttp.New(c.baseURL, c.httpClient)

	return c
}

// Token returns the access token (first and last 4 characters for debugging).
func (c *Client) Token() string {
	if len(c.token) <= 8 {
		return "****"
	}
	return c.token[:4] + "..." + c.token[len(c.token)-4:]
}

// BaseURL returns the base URL being used.
func (c *Client) BaseURL() string {
	return c.baseURL
}

// HTTPClient returns the underlying HTTP client.
func (c *Client) HTTPClient() *http.Client {
	return c.httpClient
}

// Geocoding returns a Geocoding service client.
func (c *Client) Geocoding() *geocoding.Service {
	return geocoding.New(c.token, c.http)
}

// SearchBox returns a Search Box API service client.
func (c *Client) SearchBox() *searchbox.Service {
	return searchbox.New(c.token, c.http)
}
