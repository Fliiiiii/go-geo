package types

const GeometryLineString GeometryType = "LineString"

// LineString представляет линию в виде массива точек
type LineString []Point

func (ls LineString) coordinates() {}

// NewLineString создает новую линию из массива точек
func NewLineString(points ...Point) LineString {
	return LineString(points)
}
