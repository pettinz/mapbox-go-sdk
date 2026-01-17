package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestClient_Get(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse string
		serverStatus   int
		wantErr        bool
		expectedResult map[string]any
	}{
		{
			name:           "successful GET request",
			serverResponse: `{"message": "success", "value": 42}`,
			serverStatus:   http.StatusOK,
			wantErr:        false,
			expectedResult: map[string]any{"message": "success", "value": float64(42)},
		},
		{
			name:           "error response",
			serverResponse: `{"message": "invalid token", "code": "TOKEN_INVALID"}`,
			serverStatus:   http.StatusUnauthorized,
			wantErr:        true,
		},
		{
			name:           "server error",
			serverResponse: `{"message": "internal server error"}`,
			serverStatus:   http.StatusInternalServerError,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET request, got %s", r.Method)
				}
				w.WriteHeader(tt.serverStatus)
				w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			client := New(server.URL, nil)
			var result map[string]any

			err := client.Get(context.Background(), "/test", url.Values{"key": []string{"value"}}, &result)

			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && result["message"] != tt.expectedResult["message"] {
				t.Errorf("Get() result = %v, want %v", result, tt.expectedResult)
			}
		})
	}
}

func TestClient_Post(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]any
		serverResponse string
		serverStatus   int
		wantErr        bool
	}{
		{
			name:           "successful POST request",
			requestBody:    map[string]any{"query": "test"},
			serverResponse: `{"result": "ok"}`,
			serverStatus:   http.StatusOK,
			wantErr:        false,
		},
		{
			name:           "validation error",
			requestBody:    map[string]any{"query": ""},
			serverResponse: `{"message": "query is required", "code": "VALIDATION_ERROR"}`,
			serverStatus:   http.StatusUnprocessableEntity,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST request, got %s", r.Method)
				}
				if r.Header.Get("Content-Type") != "application/json" {
					t.Errorf("expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
				}
				w.WriteHeader(tt.serverStatus)
				w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			client := New(server.URL, nil)
			var result map[string]any

			err := client.Post(context.Background(), "/test", tt.requestBody, &result)

			if (err != nil) != tt.wantErr {
				t.Errorf("Post() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Do(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check headers
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("expected Accept header application/json, got %s", r.Header.Get("Accept"))
		}
		if r.Header.Get("User-Agent") != "github.com/pettinz/mapbox-go-sdk-go" {
			t.Errorf("expected User-Agent header github.com/pettinz/mapbox-go-sdk-go, got %s", r.Header.Get("User-Agent"))
		}

		// Check query parameters
		if r.URL.Query().Get("param") != "value" {
			t.Errorf("expected query param value, got %s", r.URL.Query().Get("param"))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))
	defer server.Close()

	client := New(server.URL, nil)
	query := url.Values{"param": []string{"value"}}

	resp, err := client.Do(context.Background(), http.MethodGet, "/test", query, nil)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Do() status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestErrorResponse_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *ErrorResponse
		expected string
	}{
		{
			name: "error with code",
			err: &ErrorResponse{
				StatusCode: 401,
				Message:    "unauthorized",
				Code:       "INVALID_TOKEN",
			},
			expected: "HTTP 401: unauthorized (INVALID_TOKEN)",
		},
		{
			name: "error without code",
			err: &ErrorResponse{
				StatusCode: 404,
				Message:    "not found",
			},
			expected: "HTTP 404: not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.expected {
				t.Errorf("Error() = %q, want %q", got, tt.expected)
			}
		})
	}
}
