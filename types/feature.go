package types

// Properties представляет произвольные свойства объекта
type Properties map[string]interface{}

// Feature представляет объект GeoJSON с геометрией и свойствами
type Feature struct {
	Type       string     `json:"type"`
	ID         string     `json:"id,omitempty"`
	Geometry   Geometry   `json:"geometry"`
	Properties Properties `json:"properties"`
}

// NewFeature создает новый объект Feature
func NewFeature(id string, geometry Geometry, properties Properties) Feature {
	return Feature{
		Type:       "Feature",
		ID:         id,
		Geometry:   geometry,
		Properties: properties,
	}
}

// NewProperties создает новый объект Properties
func NewProperties() Properties {
	return Properties{}
}
