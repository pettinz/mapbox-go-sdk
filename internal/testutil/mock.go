// Package testutil provides testing utilities for the Mapbox SDK.
package testutil

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockServer creates a test HTTP server with the given handler.
func MockServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	return httptest.NewServer(handler)
}

// MockResponse creates a mock HTTP handler that returns the given response.
func MockResponse(statusCode int, body string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		w.Write([]byte(body))
	}
}

// MockResponseWithHeaders creates a mock HTTP handler with custom headers.
func MockResponseWithHeaders(statusCode int, body string, headers map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for key, value := range headers {
			w.Header().Set(key, value)
		}
		w.WriteHeader(statusCode)
		w.Write([]byte(body))
	}
}

// AssertQueryParam asserts that a query parameter has the expected value.
func AssertQueryParam(t *testing.T, r *http.Request, key, expected string) {
	t.Helper()
	actual := r.URL.Query().Get(key)
	if actual != expected {
		t.Errorf("expected query param %s=%q, got %q", key, expected, actual)
	}
}

// AssertMethod asserts that the request method is the expected one.
func AssertMethod(t *testing.T, r *http.Request, expected string) {
	t.Helper()
	if r.Method != expected {
		t.Errorf("expected method %s, got %s", expected, r.Method)
	}
}

// AssertHeader asserts that a header has the expected value.
func AssertHeader(t *testing.T, r *http.Request, key, expected string) {
	t.Helper()
	actual := r.Header.Get(key)
	if actual != expected {
		t.Errorf("expected header %s=%q, got %q", key, expected, actual)
	}
}
