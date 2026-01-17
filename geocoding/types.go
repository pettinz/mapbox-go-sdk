// Package geocoding provides access to the Mapbox Geocoding API.
package geocoding

// ForwardRequest represents a forward geocoding request using text-based search.
type ForwardRequest struct {
	// Query is the search text (required).
	Query string `json:"q"`

	// Autocomplete specifies whether to return autocomplete results (default: true).
	Autocomplete *bool `json:"autocomplete,omitempty"`

	// BBox limits results to a bounding box [min_lon, min_lat, max_lon, max_lat].
	BBox []float64 `json:"bbox,omitempty"`

	// Country limits results to one or more countries (ISO 3166 alpha-2 codes).
	Country []string `json:"country,omitempty"`

	// Language sets the language for results (IETF language tags).
	Language string `json:"language,omitempty"`

	// Limit sets the maximum number of results (1-10, default: 5).
	Limit *int `json:"limit,omitempty"`

	// Proximity biases results toward a location [lon, lat].
	Proximity []float64 `json:"proximity,omitempty"`

	// Types filters results by feature types.
	Types []string `json:"types,omitempty"`

	// Worldview returns features for a specific worldview (country code).
	Worldview string `json:"worldview,omitempty"`
}

// StructuredForwardRequest represents a forward geocoding request using structured address components.
type StructuredForwardRequest struct {
	// AddressNumber is the house or street number.
	AddressNumber string `json:"address_number,omitempty"`

	// Street is the street name.
	Street string `json:"street,omitempty"`

	// Block is the block name.
	Block string `json:"block,omitempty"`

	// Place is the city, town, or village.
	Place string `json:"place,omitempty"`

	// Region is the state, province, or region.
	Region string `json:"region,omitempty"`

	// Postcode is the postal code.
	Postcode string `json:"postcode,omitempty"`

	// Country is the country (ISO 3166 alpha-2 or alpha-3 code, or full name).
	Country string `json:"country,omitempty"`

	// Autocomplete specifies whether to return autocomplete results (default: true).
	Autocomplete *bool `json:"autocomplete,omitempty"`

	// BBox limits results to a bounding box [min_lon, min_lat, max_lon, max_lat].
	BBox []float64 `json:"bbox,omitempty"`

	// Language sets the language for results (IETF language tags).
	Language string `json:"language,omitempty"`

	// Limit sets the maximum number of results (1-10, default: 5).
	Limit *int `json:"limit,omitempty"`

	// Proximity biases results toward a location [lon, lat].
	Proximity []float64 `json:"proximity,omitempty"`

	// Worldview returns features for a specific worldview (country code).
	Worldview string `json:"worldview,omitempty"`
}

// ReverseRequest represents a reverse geocoding request.
type ReverseRequest struct {
	// Longitude is the longitude coordinate (required, -180 to 180).
	Longitude float64 `json:"longitude"`

	// Latitude is the latitude coordinate (required, -90 to 90).
	Latitude float64 `json:"latitude"`

	// Country limits results to one or more countries (ISO 3166 alpha-2 codes).
	Country []string `json:"country,omitempty"`

	// Language sets the language for results (IETF language tags).
	Language string `json:"language,omitempty"`

	// Limit sets the maximum number of results (1-10, default: 1).
	Limit *int `json:"limit,omitempty"`

	// Types filters results by feature types.
	Types []string `json:"types,omitempty"`

	// Worldview returns features for a specific worldview (country code).
	Worldview string `json:"worldview,omitempty"`
}

// BatchRequest represents a batch geocoding request.
type BatchRequest struct {
	// Queries is the list of forward or reverse queries (max 1000).
	Queries []BatchQuery `json:"queries"`
}

// BatchQuery represents a single query in a batch request.
type BatchQuery struct {
	// ID is an optional identifier for this query.
	ID string `json:"id,omitempty"`

	// Query is the search text for forward geocoding.
	Query string `json:"q,omitempty"`

	// Longitude is the longitude for reverse geocoding.
	Longitude *float64 `json:"longitude,omitempty"`

	// Latitude is the latitude for reverse geocoding.
	Latitude *float64 `json:"latitude,omitempty"`

	// Country limits results to one or more countries.
	Country []string `json:"country,omitempty"`

	// Language sets the language for results.
	Language string `json:"language,omitempty"`

	// Limit sets the maximum number of results.
	Limit *int `json:"limit,omitempty"`

	// Types filters results by feature types.
	Types []string `json:"types,omitempty"`
}

// Response represents a geocoding API response.
type Response struct {
	// Type is the GeoJSON type (should be "FeatureCollection").
	Type string `json:"type"`

	// Features is the list of results.
	Features []Feature `json:"features"`

	// Attribution is the data attribution text.
	Attribution string `json:"attribution,omitempty"`
}

// Feature represents a single geocoding result.
type Feature struct {
	// Type is the GeoJSON type (should be "Feature").
	Type string `json:"type"`

	// ID is a unique identifier for this feature.
	ID string `json:"id"`

	// Geometry contains the geographic coordinates.
	Geometry Geometry `json:"geometry"`

	// Properties contains feature metadata.
	Properties Properties `json:"properties"`
}

// Geometry represents the geographic coordinates of a feature.
type Geometry struct {
	// Type is the GeoJSON geometry type (should be "Point").
	Type string `json:"type"`

	// Coordinates contains [longitude, latitude].
	Coordinates []float64 `json:"coordinates"`
}

// Properties contains metadata about a geocoding result.
type Properties struct {
	// MapboxID is a unique Mapbox identifier.
	MapboxID string `json:"mapbox_id"`

	// FeatureType is the type of feature (e.g., "address", "place", "country").
	FeatureType string `json:"feature_type"`

	// Name is the feature's name.
	Name string `json:"name"`

	// NamePreferred is the preferred name for this feature.
	NamePreferred string `json:"name_preferred,omitempty"`

	// PlaceName is the full place name string.
	PlaceName string `json:"place_name"`

	// PlaceNamePreferred is the preferred full place name.
	PlaceNamePreferred string `json:"place_name_preferred,omitempty"`

	// Context contains contextual information (region, country, etc.).
	Context map[string]Context `json:"context,omitempty"`

	// Coordinates contains the feature's coordinates.
	Coordinates Coordinates `json:"coordinates"`

	// BBox is the bounding box [min_lon, min_lat, max_lon, max_lat].
	BBox []float64 `json:"bbox,omitempty"`

	// MatchCode indicates the quality of the match.
	MatchCode *MatchCode `json:"match_code,omitempty"`

	// Accuracy indicates the precision of the coordinates.
	Accuracy string `json:"accuracy,omitempty"`

	// AddressNumber is the street number (for address features).
	AddressNumber string `json:"address_number,omitempty"`

	// Street is the street name (for address features).
	Street string `json:"street,omitempty"`

	// Postcode is the postal code.
	Postcode string `json:"postcode,omitempty"`
}

// Context represents contextual information about a feature's location.
type Context struct {
	// MapboxID is the unique identifier for this context.
	MapboxID string `json:"mapbox_id"`

	// Name is the name of this context (e.g., region or country name).
	Name string `json:"name"`

	// Wikidata is the Wikidata identifier.
	Wikidata string `json:"wikidata,omitempty"`

	// ShortCode is a short code for this context (e.g., "US-CA").
	ShortCode string `json:"short_code,omitempty"`
}

// Coordinates represents longitude and latitude coordinates.
type Coordinates struct {
	// Longitude is the longitude coordinate.
	Longitude float64 `json:"longitude"`

	// Latitude is the latitude coordinate.
	Latitude float64 `json:"latitude"`
}

// MatchCode indicates the quality of the geocoding match.
type MatchCode struct {
	// Confidence indicates the confidence level (exact, high, medium, low).
	Confidence string `json:"confidence"`

	// Accuracy indicates the precision of the match.
	Accuracy string `json:"accuracy,omitempty"`
}

// BatchResponse represents a batch geocoding API response.
type BatchResponse struct {
	// Results contains the response for each query in the batch.
	Results []BatchResult `json:"results"`
}

// BatchResult represents the result of a single query in a batch.
type BatchResult struct {
	// ID is the identifier provided in the request.
	ID string `json:"id,omitempty"`

	// Response contains the geocoding response for this query.
	Response *Response `json:"response,omitempty"`

	// Error contains error information if the query failed.
	Error *BatchError `json:"error,omitempty"`
}

// BatchError represents an error for a single query in a batch.
type BatchError struct {
	// Message is the error message.
	Message string `json:"message"`

	// Code is the error code.
	Code string `json:"code,omitempty"`
}
