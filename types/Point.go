package types

const GeometryPoint GeometryType = "Point"

// Point представляет точку в виде массива координат [longitude, latitude]
type Point []float64

func (p Point) coordinates() {}

// GetLongitude возвращает долготу точки
func (p Point) GetLongitude() float64 {
	if len(p) > 0 {
		return p[0]
	}
	return 0
}

// GetLatitude возвращает широту точки
func (p Point) GetLatitude() float64 {
	if len(p) > 1 {
		return p[1]
	}
	return 0
}

// NewPoint создает новую точку с указанными координатами
func NewPoint(longitude, latitude float64) Point {
	return Point{longitude, latitude}
}
