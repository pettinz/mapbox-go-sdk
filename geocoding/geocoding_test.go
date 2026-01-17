package geocoding

import (
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

	if service.httpClient == nil {
		t.Error("expected non-nil httpClient")
	}
}

func TestValidateCoordinates(t *testing.T) {
	tests := []struct {
		name      string
		longitude float64
		latitude  float64
		wantErr   bool
	}{
		{
			name:      "valid coordinates",
			longitude: -122.4194,
			latitude:  37.7749,
			wantErr:   false,
		},
		{
			name:      "valid boundary coordinates",
			longitude: 180,
			latitude:  90,
			wantErr:   false,
		},
		{
			name:      "valid negative boundary coordinates",
			longitude: -180,
			latitude:  -90,
			wantErr:   false,
		},
		{
			name:      "longitude too high",
			longitude: 181,
			latitude:  0,
			wantErr:   true,
		},
		{
			name:      "longitude too low",
			longitude: -181,
			latitude:  0,
			wantErr:   true,
		},
		{
			name:      "latitude too high",
			longitude: 0,
			latitude:  91,
			wantErr:   true,
		},
		{
			name:      "latitude too low",
			longitude: 0,
			latitude:  -91,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCoordinates(tt.longitude, tt.latitude)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateCoordinates() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFormatFloatArray(t *testing.T) {
	tests := []struct {
		name     string
		input    []float64
		expected string
	}{
		{
			name:     "two floats",
			input:    []float64{-122.4194, 37.7749},
			expected: "-122.4194,37.7749",
		},
		{
			name:     "four floats (bbox)",
			input:    []float64{-122.5, 37.7, -122.3, 37.8},
			expected: "-122.5,37.7,-122.3,37.8",
		},
		{
			name:     "integers",
			input:    []float64{1, 2, 3},
			expected: "1,2,3",
		},
		{
			name:     "empty array",
			input:    []float64{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatFloatArray(tt.input)
			if result != tt.expected {
				t.Errorf("formatFloatArray() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func boolPtr(b bool) *bool {
	return &b
}

func intPtr(i int) *int {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}
