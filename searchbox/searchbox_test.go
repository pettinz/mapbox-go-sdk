package searchbox

import (
	"regexp"
	"testing"

	internalhttp "github.com/pettinz/mapbox-go-sdk/internal/http"
)

func TestNew(t *testing.T) {
	token := "test-token"
	httpClient := internalhttp.New("https://api.mapbox.com", nil)

	service := New(token, httpClient)

	if service == nil {
		t.Fatal("expected non-nil service")
	}

	if service.token != token {
		t.Errorf("expected token %q, got %q", token, service.token)
	}

	if service.httpClient != httpClient {
		t.Error("expected httpClient to be set")
	}
}

func TestNewSessionToken(t *testing.T) {
	// UUIDv4 pattern: xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx
	// where y is 8, 9, a, or b
	uuidPattern := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

	// Generate multiple tokens to ensure uniqueness
	tokens := make(map[string]bool)
	for i := 0; i < 100; i++ {
		token := NewSessionToken()

		// Validate format
		if !uuidPattern.MatchString(token) {
			t.Errorf("token %q does not match UUIDv4 pattern", token)
		}

		// Check uniqueness
		if tokens[token] {
			t.Errorf("duplicate token generated: %q", token)
		}
		tokens[token] = true
	}

	if len(tokens) != 100 {
		t.Errorf("expected 100 unique tokens, got %d", len(tokens))
	}
}

// Helper functions for tests

func intPtr(i int) *int {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}
