# go-geo

`go-geo` is a comprehensive Go package for geospatial operations and calculations. It provides tools for working with geographic coordinates, creating various grid systems, calculating distances, and managing GeoJSON-compatible geometric objects.

## Features

- **Geographic Calculations**

  - Distance and bearing calculations between points
  - Destination point calculation based on distance and bearing
  - Area and length calculations for polygons and line strings
  - Point-in-polygon testing

- **Grid Systems**

  - Rectangular grids with latitude correction
  - Hexagonal grids with Mercator projection support
  - Radial grids and sectors around a center point

- **GeoJSON-Compatible Types**
  - Point, LineString, Polygon, MultiPolygon
  - Feature and FeatureCollection types
  - Geometry and GeometryCollection types
  - Full JSON and BSON serialization/deserialization

## Installation

```bash
go get github.com/Fliiiiii/go-geo
```

## Usage Examples

### Basic Distance Calculation

```go
package main

import (
    "fmt"
    "github.com/Fliiiiii/go-geo/calc"
    "github.com/Fliiiiii/go-geo/types"
)

func main() {
    // Create two points
    moscow := types.NewPoint(37.6173, 55.7558)    // Moscow coordinates
    stPetersburg := types.NewPoint(30.3351, 59.9343)  // St. Petersburg coordinates

    // Calculate distance in kilometers
    distance := calc.CalculateDistance(moscow, stPetersburg)

    fmt.Printf("Distance between Moscow and St. Petersburg: %.2f km\n", distance)
}
```

### Creating a Hexagonal Grid

```go
package main

import (
    "encoding/json"
    "fmt"
    "github.com/Fliiiiii/go-geo/grid"
    "github.com/Fliiiiii/go-geo/types"
    "os"
)

func main() {
    // Create a hexagonal grid for Moscow area
    hexGrid := grid.CreateHexagonalGrid(37.5, 55.7, 37.7, 55.8, 0.01)

    // Convert to GeoJSON FeatureCollection
    features := make([]types.Feature, len(hexGrid))
    for i, poly := range hexGrid {
        features[i] = types.NewFeature(
            types.NewPolygonGeometry(poly),
            types.NewProperties(),
        )
    }

    collection := types.NewFeatureCollection(features...)

    // Output as GeoJSON
    json.NewEncoder(os.Stdout).Encode(collection)
}
```

### Point-in-Polygon Test

```go
package main

import (
    "fmt"
    "github.com/Fliiiiii/go-geo/calc"
    "github.com/Fliiiiii/go-geo/types"
)

func main() {
    // Create a polygon (simplified outline of Moscow's MKAD)
    outerRing := types.NewLineString(
        types.NewPoint(37.329, 55.574),
        types.NewPoint(37.844, 55.577),
        types.NewPoint(37.743, 55.851),
        types.NewPoint(37.370, 55.888),
        types.NewPoint(37.329, 55.574),
    )
    moscow := types.NewPolygon(outerRing)

    // Test if points are inside the polygon
    redSquare := types.NewPoint(37.620, 55.754)
    sheremetyevo := types.NewPoint(37.414, 55.972)

    fmt.Printf("Red Square is inside MKAD: %v\n", calc.PointInPolygon(moscow, redSquare))
    fmt.Printf("Sheremetyevo Airport is inside MKAD: %v\n", calc.PointInPolygon(moscow, sheremetyevo))
}
```
