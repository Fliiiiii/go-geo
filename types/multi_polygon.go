package types

const GeometryMultiPolygon GeometryType = "MultiPolygon"

// MultiPolygon представляет набор полигонов
type MultiPolygon []Polygon

func (mp MultiPolygon) coordinates() {}

// NewMultiPolygon создает новый MultiPolygon из массива полигонов
func NewMultiPolygon(polygons ...Polygon) MultiPolygon {
	return MultiPolygon(polygons)
}

// Append добавляет полигоны в MultiPolygon
func (mp *MultiPolygon) Append(polygons ...Polygon) {
	*mp = append(*mp, polygons...)
}
