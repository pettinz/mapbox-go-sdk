package searchbox

// NavigationOptions configures navigation and ETA calculations.
type NavigationOptions struct {
	ETAType string    // "navigation" for ETA calculations
	Origin  []float64 // [lon, lat] for ETA from origin
	Profile string    // "driving", "walking", "cycling"
}

// SAROptions configures Search Along Route parameters.
type SAROptions struct {
	Type          string      // "isochrone"
	Route         [][]float64 // route coordinates [[lon, lat], ...]
	TimeDeviation *int        // seconds of acceptable time deviation
}

// SuggestRequest represents a request to the Suggest endpoint.
type SuggestRequest struct {
	Query        string             // required, max 256 chars
	SessionToken string             // required, UUIDv4
	Proximity    []float64          // [lon, lat] or use ProximityIP
	ProximityIP  bool               // set to true to use IP-based proximity
	BBox         []float64          // [min_lon, min_lat, max_lon, max_lat]
	Country      []string           // ISO 3166-1 alpha-2 codes
	Language     string             // IETF language tag
	Limit        *int               // max 10
	Types        []string           // feature types
	POICategory  []string           // POI categories
	Navigation   *NavigationOptions // ETA options
}

// SuggestResponse represents the response from the Suggest endpoint.
type SuggestResponse struct {
	Suggestions []Suggestion       `json:"suggestions"`
	Attribution string             `json:"attribution"`
	ResponseID  string             `json:"response_id,omitempty"`
}

// Suggestion represents a single suggestion (without coordinates).
type Suggestion struct {
	MapboxID       string          `json:"mapbox_id"`
	FeatureType    string          `json:"feature_type"`
	Name           string          `json:"name"`
	NamePreferred  string          `json:"name_preferred,omitempty"`
	PlaceFormatted string          `json:"place_formatted,omitempty"`
	Address        string          `json:"address,omitempty"`
	FullAddress    string          `json:"full_address,omitempty"`
	Context        *Context        `json:"context,omitempty"`
	POICategory    []string        `json:"poi_category,omitempty"`
	POICategoryIDs []string        `json:"poi_category_ids,omitempty"`
	Brand          []string        `json:"brand,omitempty"`
	MakiIcon       string          `json:"maki,omitempty"`
	Metadata       map[string]any  `json:"metadata,omitempty"`
	Distance       *float64        `json:"distance,omitempty"`
	ETA            *float64        `json:"eta,omitempty"`
}

// RetrieveRequest represents a request to the Retrieve endpoint.
type RetrieveRequest struct {
	MapboxID     string             // required, from suggest response
	SessionToken string             // required, same as suggest
	Navigation   *NavigationOptions // optional ETA
}

// RetrieveResponse represents the response from the Retrieve endpoint.
type RetrieveResponse struct {
	Type        string    `json:"type"`
	Features    []Feature `json:"features"`
	Attribution string    `json:"attribution"`
	ResponseID  string    `json:"response_id,omitempty"`
}

// ForwardRequest represents a request to the Forward search endpoint.
type ForwardRequest struct {
	Query        string             // required
	Autocomplete *bool              // default true
	Proximity    []float64          // [lon, lat] or use ProximityIP
	ProximityIP  bool               // set to true to use IP-based proximity
	BBox         []float64          // [min_lon, min_lat, max_lon, max_lat]
	Country      []string           // ISO 3166-1 alpha-2 codes
	Language     string             // IETF language tag
	Limit        *int               // max 10
	Types        []string           // feature types
	POICategory  []string           // POI categories
	Navigation   *NavigationOptions // ETA options
}

// ForwardResponse represents the response from the Forward search endpoint.
type ForwardResponse struct {
	Type        string    `json:"type"`
	Features    []Feature `json:"features"`
	Attribution string    `json:"attribution"`
	ResponseID  string    `json:"response_id,omitempty"`
}

// CategorySearchRequest represents a request to search by category.
type CategorySearchRequest struct {
	CategoryID  string             // required, canonical category ID
	Proximity   []float64          // [lon, lat], required if no bbox or SAR
	ProximityIP bool               // set to true to use IP-based proximity
	BBox        []float64          // [min_lon, min_lat, max_lon, max_lat], required if no proximity or SAR
	Country     []string           // ISO 3166-1 alpha-2 codes
	Language    string             // IETF language tag
	Limit       *int               // max 25
	Navigation  *NavigationOptions // ETA options
	SAR         *SAROptions        // Search Along Route, required if no proximity or bbox
}

// CategorySearchResponse represents the response from category search.
type CategorySearchResponse struct {
	Type        string    `json:"type"`
	Features    []Feature `json:"features"`
	Attribution string    `json:"attribution"`
	ResponseID  string    `json:"response_id,omitempty"`
}

// ListCategoriesRequest represents a request to list available categories.
type ListCategoriesRequest struct {
	Language string // optional, IETF language tag
}

// ListCategoriesResponse represents the response from list categories.
type ListCategoriesResponse struct {
	Categories []Category `json:"categories"`
}

// Category represents a POI category.
type Category struct {
	CanonicalID string `json:"canonical_name"`
	Name        string `json:"name"`
	MakiIcon    string `json:"maki_icon,omitempty"`
}

// ReverseRequest represents a request to the Reverse geocoding endpoint.
type ReverseRequest struct {
	Longitude float64  // required
	Latitude  float64  // required
	Country   []string // ISO 3166-1 alpha-2 codes
	Language  string   // IETF language tag
	Limit     *int     // max 10
	Types     []string // feature types
}

// ReverseResponse represents the response from reverse geocoding.
type ReverseResponse struct {
	Type        string    `json:"type"`
	Features    []Feature `json:"features"`
	Attribution string    `json:"attribution"`
	ResponseID  string    `json:"response_id,omitempty"`
}

// Feature represents a GeoJSON feature in the response.
type Feature struct {
	Type       string            `json:"type"`
	ID         string            `json:"id"`
	Geometry   Geometry          `json:"geometry"`
	Properties FeatureProperties `json:"properties"`
}

// FeatureProperties contains the properties of a feature.
type FeatureProperties struct {
	MapboxID       string            `json:"mapbox_id"`
	FeatureType    string            `json:"feature_type"`
	Name           string            `json:"name"`
	NamePreferred  string            `json:"name_preferred,omitempty"`
	PlaceFormatted string            `json:"place_formatted,omitempty"`
	FullAddress    string            `json:"full_address,omitempty"`
	Address        string            `json:"address,omitempty"`
	AddressNumber  string            `json:"address_number,omitempty"`
	Street         string            `json:"street,omitempty"`
	Postcode       string            `json:"postcode,omitempty"`

	// Context hierarchy
	Context *Context `json:"context,omitempty"`

	// Coordinates and routing
	Coordinates    Coordinates     `json:"coordinates"`
	Accuracy       string          `json:"accuracy,omitempty"`
	RoutablePoints []RoutablePoint `json:"routable_points,omitempty"`

	// POI specific
	POICategory    []string `json:"poi_category,omitempty"`
	POICategoryIDs []string `json:"poi_category_ids,omitempty"`
	Brand          []string `json:"brand,omitempty"`
	BrandID        []string `json:"brand_id,omitempty"`
	MakiIcon       string   `json:"maki,omitempty"`

	// External identifiers
	ExternalIDs map[string]string `json:"external_ids,omitempty"`

	// Additional metadata
	Metadata map[string]any `json:"metadata,omitempty"`

	// ETA information (if requested)
	Distance *float64 `json:"distance,omitempty"`
	ETA      *float64 `json:"eta,omitempty"`
}

// Context represents the hierarchical context of a feature.
type Context struct {
	Country      *ContextElement `json:"country,omitempty"`
	Region       *ContextElement `json:"region,omitempty"`
	Postcode     *ContextElement `json:"postcode,omitempty"`
	District     *ContextElement `json:"district,omitempty"`
	Place        *ContextElement `json:"place,omitempty"`
	Locality     *ContextElement `json:"locality,omitempty"`
	Neighborhood *ContextElement `json:"neighborhood,omitempty"`
	Street       *ContextElement `json:"street,omitempty"`
	Address      *ContextElement `json:"address,omitempty"`
}

// ContextElement represents a single element in the context hierarchy.
type ContextElement struct {
	MapboxID    string `json:"mapbox_id"`
	Name        string `json:"name"`
	WikidataID  string `json:"wikidata_id,omitempty"`
	ShortCode   string `json:"short_code,omitempty"`
	CountryCode string `json:"country_code,omitempty"`
	RegionCode  string `json:"region_code,omitempty"`
}

// RoutablePoint represents a point suitable for navigation.
type RoutablePoint struct {
	Name        string    `json:"name"`
	Coordinates []float64 `json:"coordinates"` // [lon, lat]
}

// Coordinates represents a geographic coordinate pair.
type Coordinates struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

// Geometry represents a GeoJSON geometry.
type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"` // [lon, lat]
}
