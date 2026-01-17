// Package http provides an internal HTTP client wrapper for Mapbox API requests.
package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Client is an HTTP client wrapper for making API requests.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// New creates a new HTTP client with the given base URL and HTTP client.
func New(baseURL string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

// Do executes an HTTP request and returns the response.
func (c *Client) Do(ctx context.Context, method, path string, query url.Values, body any) (*http.Response, error) {
	// Build the full URL
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	// Add query parameters
	if query != nil {
		u.RawQuery = query.Encode()
	}

	// Create request body
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, method, u.String(), bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "github.com/pettinz/mapbox-go-sdk-go")

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return resp, nil
}

// Get executes a GET request and unmarshals the response into result.
func (c *Client) Get(ctx context.Context, path string, query url.Values, result any) error {
	resp, err := c.Do(ctx, http.MethodGet, path, query, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.handleResponse(resp, result)
}

// Post executes a POST request and unmarshals the response into result.
func (c *Client) Post(ctx context.Context, path string, body any, result any) error {
	resp, err := c.Do(ctx, http.MethodPost, path, nil, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.handleResponse(resp, result)
}

// handleResponse processes the HTTP response and handles errors.
func (c *Client) handleResponse(resp *http.Response, result any) error {
	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Try to parse error response
		var errResp struct {
			Message string `json:"message"`
			Code    string `json:"code"`
		}
		if err := json.Unmarshal(body, &errResp); err == nil && errResp.Message != "" {
			return &ErrorResponse{
				StatusCode: resp.StatusCode,
				Message:    errResp.Message,
				Code:       errResp.Code,
			}
		}

		// Fallback to status text
		return &ErrorResponse{
			StatusCode: resp.StatusCode,
			Message:    resp.Status,
		}
	}

	// Unmarshal successful response
	if result != nil {
		if err := json.Unmarshal(body, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// ErrorResponse represents an error response from the API.
type ErrorResponse struct {
	StatusCode int
	Message    string
	Code       string
}

// Error implements the error interface.
func (e *ErrorResponse) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("HTTP %d: %s (%s)", e.StatusCode, e.Message, e.Code)
	}
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Message)
}
