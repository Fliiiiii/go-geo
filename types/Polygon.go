package types

const GeometryPolygon GeometryType = "Polygon"

// Polygon представляет полигон в виде массива линий
type Polygon []LineString

func (p Polygon) coordinates() {}

// NewPolygon создает новый полигон из массива линий
func NewPolygon(lines ...LineString) Polygon {
	return Polygon(lines)
}
