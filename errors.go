// Package mapbox provides a Go client for the Mapbox API.
package mapbox

import (
	"errors"
	"fmt"
)

// Error represents a Mapbox API error.
type Error struct {
	StatusCode int    `json:"statusCode,omitempty"`
	Message    string `json:"message"`
	Code       string `json:"code,omitempty"`
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("mapbox: %s (%s)", e.Message, e.Code)
	}
	return fmt.Sprintf("mapbox: %s", e.Message)
}

// Common errors returned by the Mapbox API.
var (
	// ErrInvalidToken is returned when the access token is invalid or missing.
	ErrInvalidToken = errors.New("mapbox: invalid or missing access token")

	// ErrRateLimitExceeded is returned when the rate limit is exceeded.
	ErrRateLimitExceeded = errors.New("mapbox: rate limit exceeded")

	// ErrInvalidRequest is returned when the request is invalid.
	ErrInvalidRequest = errors.New("mapbox: invalid request")

	// ErrNotFound is returned when the requested resource is not found.
	ErrNotFound = errors.New("mapbox: resource not found")

	// ErrServerError is returned when the server encounters an error.
	ErrServerError = errors.New("mapbox: server error")
)

// mapStatusCodeToError maps HTTP status codes to predefined errors.
func mapStatusCodeToError(statusCode int, message string) error {
	switch statusCode {
	case 401, 403:
		return fmt.Errorf("%w: %s", ErrInvalidToken, message)
	case 404:
		return fmt.Errorf("%w: %s", ErrNotFound, message)
	case 422:
		return fmt.Errorf("%w: %s", ErrInvalidRequest, message)
	case 429:
		return fmt.Errorf("%w: %s", ErrRateLimitExceeded, message)
	case 500, 502, 503, 504:
		return fmt.Errorf("%w: %s", ErrServerError, message)
	default:
		return &Error{
			StatusCode: statusCode,
			Message:    message,
		}
	}
}
