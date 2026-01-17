# Mapbox SDK for Go

An idiomatic and type-safe Go SDK for the Mapbox API, starting with complete Geocoding API support.

## Features

- Type-safe request and response structures
- Idiomatic Go API with context support
- Comprehensive error handling
- Support for all Geocoding API endpoints:
  - Forward geocoding (text-based and structured)
  - Reverse geocoding
  - Batch geocoding
- No external dependencies (uses only Go standard library)
- Extensive test coverage

## Installation

```bash
go get github.com/pettinz/mapbox-go-sdk
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/pettinz/mapbox-go-sdk"
    "github.com/pettinz/mapbox-go-sdk/geocoding"
)

func main() {
    // Create a new client with your Mapbox access token
    client := mapbox.NewClient("your-mapbox-access-token")

    // Get the geocoding service
    geo := client.Geocoding()

    // Forward geocoding
    resp, err := geo.Forward(context.Background(), &geocoding.ForwardRequest{
        Query:    "1600 Pennsylvania Avenue NW, Washington, DC",
        Language: "en",
        Limit:    intPtr(5),
    })
    if err != nil {
        log.Fatal(err)
    }

    // Print results
    for _, feature := range resp.Features {
        fmt.Printf("%s: [%f, %f]\n",
            feature.Properties.PlaceName,
            feature.Properties.Coordinates.Longitude,
            feature.Properties.Coordinates.Latitude,
        )
    }
}

func intPtr(i int) *int {
    return &i
}
```

## Usage

### Creating a Client

```go
// Basic client with default settings
client := mapbox.NewClient("your-access-token")

// Client with custom HTTP client
httpClient := &http.Client{
    Timeout: 30 * time.Second,
}
client := mapbox.NewClient("your-access-token",
    mapbox.WithHTTPClient(httpClient),
)

// Client with custom base URL (for testing or proxies)
client := mapbox.NewClient("your-access-token",
    mapbox.WithBaseURL("https://custom-api-url.com"),
)
```

### Forward Geocoding (Text-Based)

Convert a text query into geographic coordinates:

```go
resp, err := geo.Forward(context.Background(), &geocoding.ForwardRequest{
    Query:      "San Francisco, CA",
    Language:   "en",
    Limit:      intPtr(10),
    Country:    []string{"US"},
    Types:      []string{"place", "region"},
    Proximity:  []float64{-122.4194, 37.7749}, // Bias results near SF
})
if err != nil {
    log.Fatal(err)
}

for _, feature := range resp.Features {
    fmt.Printf("Name: %s\n", feature.Properties.PlaceName)
    fmt.Printf("Coordinates: [%f, %f]\n",
        feature.Properties.Coordinates.Longitude,
        feature.Properties.Coordinates.Latitude,
    )
}
```

### Forward Geocoding (Structured)

Use structured address components:

```go
resp, err := geo.ForwardStructured(context.Background(), &geocoding.StructuredForwardRequest{
    AddressNumber: "1600",
    Street:        "Pennsylvania Avenue NW",
    Place:         "Washington",
    Region:        "DC",
    Postcode:      "20500",
    Country:       "US",
    Language:      "en",
})
if err != nil {
    log.Fatal(err)
}
```

### Reverse Geocoding

Convert coordinates into a human-readable address:

```go
resp, err := geo.Reverse(context.Background(), &geocoding.ReverseRequest{
    Longitude: -122.419415,
    Latitude:  37.774929,
    Language:  "en",
    Limit:     intPtr(1),
    Types:     []string{"address", "place"},
})
if err != nil {
    log.Fatal(err)
}

if len(resp.Features) > 0 {
    fmt.Printf("Address: %s\n", resp.Features[0].Properties.PlaceName)
}
```

### Batch Geocoding

Geocode multiple queries in a single request (up to 1000):

```go
resp, err := geo.Batch(context.Background(), &geocoding.BatchRequest{
    Queries: []geocoding.BatchQuery{
        {
            ID:    "query1",
            Query: "New York, NY",
        },
        {
            ID:        "query2",
            Longitude: float64Ptr(-118.243683),
            Latitude:  float64Ptr(34.052235),
        },
        {
            ID:    "query3",
            Query: "Chicago, IL",
        },
    },
})
if err != nil {
    log.Fatal(err)
}

for _, result := range resp.Results {
    if result.Error != nil {
        fmt.Printf("Query %s failed: %s\n", result.ID, result.Error.Message)
        continue
    }

    if len(result.Response.Features) > 0 {
        feature := result.Response.Features[0]
        fmt.Printf("Query %s: %s\n", result.ID, feature.Properties.PlaceName)
    }
}
```

## Error Handling

The SDK provides typed errors for common API error scenarios:

```go
import "errors"

resp, err := geo.Forward(ctx, req)
if err != nil {
    switch {
    case errors.Is(err, mapbox.ErrInvalidToken):
        // Handle invalid token
    case errors.Is(err, mapbox.ErrRateLimitExceeded):
        // Handle rate limit
    case errors.Is(err, mapbox.ErrInvalidRequest):
        // Handle validation error
    case errors.Is(err, mapbox.ErrNotFound):
        // Handle not found
    default:
        // Handle other errors
    }
}
```

## Context Support

All API methods accept a `context.Context` for cancellation and timeouts:

```go
// With timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

resp, err := geo.Forward(ctx, req)

// With cancellation
ctx, cancel := context.WithCancel(context.Background())
go func() {
    // Cancel after some condition
    cancel()
}()

resp, err := geo.Forward(ctx, req)
```

## Examples

See the [examples/](examples/) directory for complete working examples:

- [Forward Geocoding](examples/geocoding/forward.go)
- [Reverse Geocoding](examples/geocoding/reverse.go)
- [Batch Geocoding](examples/geocoding/batch.go)

## Testing

Run the test suite:

```bash
# Run all tests
go test ./...

# Run with coverage
go test ./... -cover

# Run with race detector
go test ./... -race

# Run verbose
go test ./... -v
```

## API Reference

### Client

- `NewClient(token string, opts ...Option) *Client` - Create a new Mapbox client
- `Geocoding() *geocoding.Service` - Get the geocoding service

### Options

- `WithHTTPClient(client *http.Client)` - Use a custom HTTP client
- `WithBaseURL(url string)` - Use a custom base URL

### Geocoding Service

- `Forward(ctx context.Context, req *ForwardRequest) (*Response, error)` - Text-based forward geocoding
- `ForwardStructured(ctx context.Context, req *StructuredForwardRequest) (*Response, error)` - Structured forward geocoding
- `Reverse(ctx context.Context, req *ReverseRequest) (*Response, error)` - Reverse geocoding
- `Batch(ctx context.Context, req *BatchRequest) (*BatchResponse, error)` - Batch geocoding

## Requirements

- Go 1.25.5 or higher
- A valid Mapbox access token (get one at [mapbox.com](https://mapbox.com))

## License

This SDK is provided as-is for working with the Mapbox API. Please refer to Mapbox's terms of service for API usage.

## Contributing

Contributions are welcome. Please ensure all tests pass and maintain the existing code style.

## Future Additions

This SDK is designed to be easily extensible. Future additions may include:

- Directions API
- Maps API
- Optimization API
- Matrix API
- Isochrone API
