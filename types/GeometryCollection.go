package types

type GeometryCollection struct {
	Geometries []Geometry `json:"geometries"`
}

func NewGeometryCollection(geometries ...Geometry) *GeometryCollection {
	return &GeometryCollection{Geometries: geometries}
}
