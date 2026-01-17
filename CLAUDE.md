# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Testing
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test ./... -cover

# Run tests with race detector
go test ./... -race

# Run tests verbosely
go test ./... -v

# Run tests for a specific package
go test ./geocoding
go test ./searchbox

# Run a specific test
go test ./searchbox -run TestService_CategorySearch
```

### Building
This is a library SDK with no build output. To verify code compiles:
```bash
go build ./...
```

## Architecture Overview

### Package Structure
The SDK is organized into three main layers:

1. **Root Package (`github.com/pettinz/mapbox-go-sdk`)** - Main entry point
   - `Client` is the root API client created with `NewClient(token, ...opts)`
   - Returns service-specific clients via `client.Geocoding()` and `client.SearchBox()`
   - Handles common configuration (token, base URL, HTTP client)

2. **Service Packages** (`geocoding/`, `searchbox/`)
   - Each service has its own package containing:
     - `Service` struct with API endpoint methods
     - Request/response types in `types.go`
     - Individual method implementations (e.g., `forward.go`, `reverse.go`, `suggest.go`)
     - Comprehensive test files with mocked HTTP responses

3. **Internal Packages** (`internal/`)
   - `internal/http` - HTTP client wrapper handling all API communication
     - Manages URL construction, query parameters, request/response marshaling
     - Converts HTTP errors to structured `ErrorResponse` types
   - `internal/testutil` - Testing utilities
     - Mock HTTP servers for testing
     - Fixture data for API responses

### HTTP Client Flow
```
Client (root)
  → creates internal/http.Client with baseURL + httpClient
  → passes to Service constructors (geocoding.New, searchbox.New)
  → Services use http.Client.Get/Post for API calls
  → http.Client handles marshaling, errors, and response parsing
```

### Service Architecture Pattern
Each service follows the same pattern:
- `Service` struct holds token and `*internalhttp.Client`
- API endpoint constants defined in the service file (e.g., `forwardPath`, `suggestPath`)
- Each endpoint has:
  - Request validation function (e.g., `validateCategorySearchRequest`)
  - Query builder function (e.g., `buildCategorySearchQuery`)
  - Main method that orchestrates validation → query building → HTTP call

### Validation Requirements
Many endpoints require specific search area parameters:

**Geocoding API:**
- Forward/Reverse: No required search area (proximity and bbox are optional)
- Batch: Queries can specify search parameters individually

**Search Box API:**
- CategorySearch: Requires **one of** `Proximity`, `ProximityIP`, `BBox`, or `SAR` (Search Along Route)
- Suggest: Optional proximity/bbox
- Forward: Optional proximity/bbox
- Reverse: Requires coordinates (longitude, latitude)

When modifying validation logic, ensure:
1. All valid parameter combinations are accepted
2. Error messages mention all valid options
3. Tests cover each valid parameter combination
4. Documentation comments reflect all valid options

### Common Types and Utilities

**GeoJSON Types** (root package):
- `Point`, `Feature`, `FeatureCollection` - Standard GeoJSON structures
- Used across both geocoding and searchbox responses

**Helper Patterns:**
- Use pointer types for optional integer parameters (`*int` for `Limit`)
- Use `[]float64` for coordinates (always `[longitude, latitude]` order)
- Use `[]string` for multi-value parameters (country codes, types, etc.)

**Navigation and SAR Options:**
- `NavigationOptions` - ETA calculations (used by multiple endpoints)
- `SAROptions` - Search Along Route parameters (route coordinates, time deviation)
- These are embedded in various request types and added to queries via helper functions

## Testing Patterns

All tests follow a table-driven approach:
```go
tests := []struct {
    name           string
    request        *RequestType
    mockStatus     int
    mockResponse   string
    wantErr        bool
    validateResult func(*testing.T, *ResponseType)
}
```

Mock HTTP responses are defined in `internal/testutil/fixtures.go` and reused across tests.

Each test:
1. Creates a mock HTTP server with `testutil.MockServer()`
2. Creates a service with the mock server URL
3. Calls the API method
4. Validates error conditions or response data

When adding new tests:
- Add fixture data to `internal/testutil/fixtures.go` if reusable
- Test both success and error cases
- Include edge cases (missing required params, out-of-range values, etc.)
- Use table-driven tests for multiple scenarios

## API Version Information

- Geocoding API: v6 (`/search/geocode/v6/...`)
- Search Box API: v1 (`/search/searchbox/v1/...`)

API paths are constants in each service's main file.
