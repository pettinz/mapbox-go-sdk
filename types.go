package mapbox

// Point represents a GeoJSON Point geometry with coordinates [longitude, latitude].
type Point struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

// NewPoint creates a new Point with the given longitude and latitude.
func NewPoint(longitude, latitude float64) *Point {
	return &Point{
		Type:        "Point",
		Coordinates: []float64{longitude, latitude},
	}
}

// Longitude returns the longitude coordinate.
func (p *Point) Longitude() float64 {
	if len(p.Coordinates) >= 1 {
		return p.Coordinates[0]
	}
	return 0
}

// Latitude returns the latitude coordinate.
func (p *Point) Latitude() float64 {
	if len(p.Coordinates) >= 2 {
		return p.Coordinates[1]
	}
	return 0
}

// Feature represents a GeoJSON Feature.
type Feature struct {
	Type       string                 `json:"type"`
	Geometry   *Point                 `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
}

// FeatureCollection represents a GeoJSON FeatureCollection.
type FeatureCollection struct {
	Type     string     `json:"type"`
	Features []*Feature `json:"features"`
}
