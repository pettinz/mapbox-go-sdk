package testutil

// ForwardGeocodingResponse is a sample forward geocoding response.
const ForwardGeocodingResponse = `{
  "type": "FeatureCollection",
  "features": [
    {
      "type": "Feature",
      "id": "address.123456",
      "geometry": {
        "type": "Point",
        "coordinates": [-77.036543, 38.897676]
      },
      "properties": {
        "mapbox_id": "dXJuOm1ieGFkcjphYmNkZWY",
        "feature_type": "address",
        "name": "1600 Pennsylvania Avenue NW",
        "name_preferred": "1600 Pennsylvania Avenue NW",
        "place_name": "1600 Pennsylvania Avenue NW, Washington, District of Columbia 20500, United States",
        "place_name_preferred": "1600 Pennsylvania Avenue NW, Washington, District of Columbia 20500, United States",
        "context": {
          "region": {
            "mapbox_id": "dXJuOm1ieHBsYzpBZ0U",
            "name": "District of Columbia",
            "short_code": "US-DC"
          },
          "country": {
            "mapbox_id": "dXJuOm1ieHBsYzpJZ00",
            "name": "United States",
            "short_code": "US"
          }
        },
        "coordinates": {
          "longitude": -77.036543,
          "latitude": 38.897676
        },
        "address_number": "1600",
        "street": "Pennsylvania Avenue NW",
        "postcode": "20500"
      }
    }
  ],
  "attribution": "© 2024 Mapbox"
}`

// ReverseGeocodingResponse is a sample reverse geocoding response.
const ReverseGeocodingResponse = `{
  "type": "FeatureCollection",
  "features": [
    {
      "type": "Feature",
      "id": "place.987654",
      "geometry": {
        "type": "Point",
        "coordinates": [-122.419415, 37.774929]
      },
      "properties": {
        "mapbox_id": "dXJuOm1ieHBsYzpBZ0U",
        "feature_type": "place",
        "name": "San Francisco",
        "name_preferred": "San Francisco",
        "place_name": "San Francisco, California, United States",
        "place_name_preferred": "San Francisco, California, United States",
        "context": {
          "region": {
            "mapbox_id": "dXJuOm1ieHBsYzpBZ0U",
            "name": "California",
            "short_code": "US-CA"
          },
          "country": {
            "mapbox_id": "dXJuOm1ieHBsYzpJZ00",
            "name": "United States",
            "short_code": "US"
          }
        },
        "coordinates": {
          "longitude": -122.419415,
          "latitude": 37.774929
        }
      }
    }
  ],
  "attribution": "© 2024 Mapbox"
}`

// BatchGeocodingResponse is a sample batch geocoding response.
const BatchGeocodingResponse = `{
  "results": [
    {
      "id": "query1",
      "response": {
        "type": "FeatureCollection",
        "features": [
          {
            "type": "Feature",
            "id": "place.123",
            "geometry": {
              "type": "Point",
              "coordinates": [-74.005974, 40.712776]
            },
            "properties": {
              "mapbox_id": "dXJuOm1ieHBsYzpBZ0U",
              "feature_type": "place",
              "name": "New York",
              "place_name": "New York, New York, United States",
              "coordinates": {
                "longitude": -74.005974,
                "latitude": 40.712776
              }
            }
          }
        ]
      }
    },
    {
      "id": "query2",
      "response": {
        "type": "FeatureCollection",
        "features": [
          {
            "type": "Feature",
            "id": "place.456",
            "geometry": {
              "type": "Point",
              "coordinates": [-118.243683, 34.052235]
            },
            "properties": {
              "mapbox_id": "dXJuOm1ieHBsYzpCZ0U",
              "feature_type": "place",
              "name": "Los Angeles",
              "place_name": "Los Angeles, California, United States",
              "coordinates": {
                "longitude": -118.243683,
                "latitude": 34.052235
              }
            }
          }
        ]
      }
    },
    {
      "id": "query3",
      "error": {
        "message": "Invalid query",
        "code": "INVALID_QUERY"
      }
    }
  ]
}`

// ErrorResponse is a sample error response.
const ErrorResponse = `{
  "message": "Invalid access token",
  "code": "TOKEN_INVALID"
}`

// RateLimitErrorResponse is a sample rate limit error response.
const RateLimitErrorResponse = `{
  "message": "Rate limit exceeded",
  "code": "RATE_LIMIT_EXCEEDED"
}`

// NotFoundErrorResponse is a sample not found error response.
const NotFoundErrorResponse = `{
  "message": "Resource not found",
  "code": "NOT_FOUND"
}`

// ValidationErrorResponse is a sample validation error response.
const ValidationErrorResponse = `{
  "message": "Validation failed",
  "code": "VALIDATION_ERROR"
}`

// SearchBoxSuggestResponse is a sample suggest response.
const SearchBoxSuggestResponse = `{
  "suggestions": [
    {
      "mapbox_id": "dXJuOm1ieHBvaTphYmNkZWY",
      "feature_type": "poi",
      "name": "Blue Bottle Coffee",
      "name_preferred": "Blue Bottle Coffee",
      "place_formatted": "San Francisco, California",
      "address": "66 Mint St",
      "full_address": "66 Mint St, San Francisco, CA 94103",
      "poi_category": ["coffee_shop"],
      "poi_category_ids": ["coffee_shop"],
      "brand": ["Blue Bottle Coffee"],
      "maki": "cafe",
      "distance": 234.5
    },
    {
      "mapbox_id": "dXJuOm1ieHBvaTpnaGlqa2w",
      "feature_type": "poi",
      "name": "Philz Coffee",
      "name_preferred": "Philz Coffee",
      "place_formatted": "San Francisco, California",
      "address": "201 Berry St",
      "full_address": "201 Berry St, San Francisco, CA 94158",
      "poi_category": ["coffee_shop"],
      "poi_category_ids": ["coffee_shop"],
      "maki": "cafe",
      "distance": 456.7
    }
  ],
  "attribution": "© 2024 Mapbox"
}`

// SearchBoxRetrieveResponse is a sample retrieve response with full GeoJSON.
const SearchBoxRetrieveResponse = `{
  "type": "FeatureCollection",
  "features": [
    {
      "type": "Feature",
      "id": "poi.123456",
      "geometry": {
        "type": "Point",
        "coordinates": [-122.394447, 37.789688]
      },
      "properties": {
        "mapbox_id": "dXJuOm1ieHBvaTphYmNkZWY",
        "feature_type": "poi",
        "name": "Blue Bottle Coffee",
        "name_preferred": "Blue Bottle Coffee",
        "place_formatted": "San Francisco, California",
        "full_address": "66 Mint St, San Francisco, CA 94103",
        "address": "66 Mint St",
        "address_number": "66",
        "street": "Mint St",
        "postcode": "94103",
        "context": {
          "place": {
            "mapbox_id": "dXJuOm1ieHBsYzpBZ0U",
            "name": "San Francisco"
          },
          "region": {
            "mapbox_id": "dXJuOm1ieHBsYzpBZ0U",
            "name": "California",
            "short_code": "US-CA"
          },
          "country": {
            "mapbox_id": "dXJuOm1ieHBsYzpJZ00",
            "name": "United States",
            "short_code": "US"
          }
        },
        "coordinates": {
          "longitude": -122.394447,
          "latitude": 37.789688
        },
        "poi_category": ["coffee_shop"],
        "poi_category_ids": ["coffee_shop"],
        "brand": ["Blue Bottle Coffee"],
        "maki": "cafe",
        "routable_points": [
          {
            "name": "main entrance",
            "coordinates": [-122.394447, 37.789688]
          }
        ],
        "external_ids": {
          "foursquare": "4ad4c05df964a520c6f920e3"
        }
      }
    }
  ],
  "attribution": "© 2024 Mapbox"
}`

// SearchBoxForwardResponse is a sample forward search response.
const SearchBoxForwardResponse = `{
  "type": "FeatureCollection",
  "features": [
    {
      "type": "Feature",
      "id": "poi.789012",
      "geometry": {
        "type": "Point",
        "coordinates": [12.496365, 41.902916]
      },
      "properties": {
        "mapbox_id": "dXJuOm1ieHBvaTptbm9wcXI",
        "feature_type": "poi",
        "name": "Colosseum",
        "name_preferred": "Colosseo",
        "place_formatted": "Rome, Italy",
        "full_address": "Piazza del Colosseo, 1, 00184 Roma RM, Italy",
        "coordinates": {
          "longitude": 12.496365,
          "latitude": 41.902916
        },
        "poi_category": ["historic_site", "landmark"],
        "poi_category_ids": ["historic_site", "landmark"],
        "maki": "monument",
        "distance": 1234.5,
        "eta": 900.0
      }
    }
  ],
  "attribution": "© 2024 Mapbox"
}`

// SearchBoxCategorySearchResponse is a sample category search response.
const SearchBoxCategorySearchResponse = `{
  "type": "FeatureCollection",
  "features": [
    {
      "type": "Feature",
      "id": "poi.345678",
      "geometry": {
        "type": "Point",
        "coordinates": [-122.408226, 37.784991]
      },
      "properties": {
        "mapbox_id": "dXJuOm1ieHBvaTpzdHV2d3g",
        "feature_type": "poi",
        "name": "Tartine Bakery",
        "name_preferred": "Tartine Bakery",
        "place_formatted": "San Francisco, California",
        "coordinates": {
          "longitude": -122.408226,
          "latitude": 37.784991
        },
        "poi_category": ["restaurant"],
        "poi_category_ids": ["restaurant"],
        "maki": "restaurant",
        "distance": 567.8
      }
    }
  ],
  "attribution": "© 2024 Mapbox"
}`

// SearchBoxListCategoriesResponse is a sample list categories response.
const SearchBoxListCategoriesResponse = `{
  "categories": [
    {
      "canonical_name": "airport",
      "name": "Airport",
      "maki_icon": "airport"
    },
    {
      "canonical_name": "coffee_shop",
      "name": "Coffee Shop",
      "maki_icon": "cafe"
    },
    {
      "canonical_name": "restaurant",
      "name": "Restaurant",
      "maki_icon": "restaurant"
    },
    {
      "canonical_name": "gas_station",
      "name": "Gas Station",
      "maki_icon": "fuel"
    }
  ]
}`

// SearchBoxReverseResponse is a sample reverse geocoding response.
const SearchBoxReverseResponse = `{
  "type": "FeatureCollection",
  "features": [
    {
      "type": "Feature",
      "id": "address.456789",
      "geometry": {
        "type": "Point",
        "coordinates": [-122.419415, 37.774929]
      },
      "properties": {
        "mapbox_id": "dXJuOm1ieGFkcjp5emFiYw",
        "feature_type": "address",
        "name": "Market Street",
        "place_formatted": "San Francisco, California 94102",
        "full_address": "Market St, San Francisco, CA 94102, United States",
        "coordinates": {
          "longitude": -122.419415,
          "latitude": 37.774929
        },
        "context": {
          "place": {
            "mapbox_id": "dXJuOm1ieHBsYzpBZ0U",
            "name": "San Francisco"
          },
          "region": {
            "mapbox_id": "dXJuOm1ieHBsYzpBZ0U",
            "name": "California",
            "short_code": "US-CA"
          },
          "country": {
            "mapbox_id": "dXJuOm1ieHBsYzpJZ00",
            "name": "United States",
            "short_code": "US"
          }
        }
      }
    }
  ],
  "attribution": "© 2024 Mapbox"
}`
