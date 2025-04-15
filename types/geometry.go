package types

import (
	"fmt"

	"github.com/Fliiiiii/go-geo/utils"
	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
)

// GeometryType определяет тип геометрии
type GeometryType string

// Coordinates определяет интерфейс для всех типов координат
// Данный интерфейс предназначен только для внутреннего использования в пакете
type Coordinates interface {
	coordinates()
}

// Geometry представляет геометрическую форму в формате GeoJSON
type Geometry struct {
	Type        GeometryType `json:"type" bson:"type"`
	Coordinates Coordinates  `json:"coordinates" bson:"coordinates"`
}

// NewPointGeometry создает геометрию типа Point
func NewPointGeometry(point Point) Geometry {
	return Geometry{
		Type:        GeometryPoint,
		Coordinates: point,
	}
}

// NewLineStringGeometry создает геометрию типа LineString
func NewLineStringGeometry(lineString LineString) Geometry {
	return Geometry{
		Type:        GeometryLineString,
		Coordinates: lineString,
	}
}

// NewPolygonGeometry создает геометрию типа Polygon
func NewPolygonGeometry(polygon Polygon) Geometry {
	return Geometry{
		Type:        GeometryPolygon,
		Coordinates: polygon,
	}
}

// NewMultiPolygonGeometry создает геометрию типа MultiPolygon
func NewMultiPolygonGeometry(multiPolygon MultiPolygon) Geometry {
	return Geometry{
		Type:        GeometryMultiPolygon,
		Coordinates: multiPolygon,
	}
}

// UnmarshalBSON реализует интерфейс bson.Unmarshaler для Geometry
func (g *Geometry) UnmarshalBSON(b []byte) error {
	if len(b) == 0 {
		// Сбрасываем все поля при получении пустого значения
		g.Type = ""
		g.Coordinates = nil
		return nil
	}

	raw := bson.Raw(b)

	// Проверяем наличие поля "type"
	typeElem := raw.Lookup("type")
	if typeElem.Type == bson.TypeNull {
		return fmt.Errorf("missing required field 'type'")
	}

	g.Type = GeometryType(typeElem.StringValue())

	// Проверяем наличие поля "coordinates"
	coordsElem := raw.Lookup("coordinates")
	if coordsElem.Type == bson.TypeNull && g.Type != "" {
		return fmt.Errorf("missing required field 'coordinates' for type %s", g.Type)
	}

	// Сбрасываем поле Coordinates перед установкой нового
	g.Coordinates = nil

	switch g.Type {
	case GeometryPoint:
		var point Point
		if err := coordsElem.Unmarshal(&point); err != nil {
			return fmt.Errorf("failed to unmarshal Point coordinates: %w", err)
		}
		g.Coordinates = point
	case GeometryLineString:
		var lineString LineString
		if err := coordsElem.Unmarshal(&lineString); err != nil {
			return fmt.Errorf("failed to unmarshal LineString coordinates: %w", err)
		}
		g.Coordinates = lineString
	case GeometryPolygon:
		var polygon Polygon
		if err := coordsElem.Unmarshal(&polygon); err != nil {
			return fmt.Errorf("failed to unmarshal Polygon coordinates: %w", err)
		}
		g.Coordinates = polygon
	case GeometryMultiPolygon:
		var multiPolygon MultiPolygon
		if err := coordsElem.Unmarshal(&multiPolygon); err != nil {
			return fmt.Errorf("failed to unmarshal MultiPolygon coordinates: %w", err)
		}
		g.Coordinates = multiPolygon
	default:
		return fmt.Errorf("unknown geometry type: %s", g.Type)
	}
	return nil
}

// MarshalBSON реализует интерфейс bson.Marshaler для Geometry
func (g *Geometry) MarshalBSON() ([]byte, error) {
	if g.Type == "" {
		return bson.Marshal(nil)
	}

	geoInterface := bson.D{{"type", g.Type}}

	if g.Coordinates == nil {
		return nil, fmt.Errorf("coordinates data is nil for geometry type %s", g.Type)
	}
	geoInterface = append(geoInterface, bson.E{"coordinates", g.Coordinates})

	return bson.Marshal(geoInterface)
}

// Используем быструю конфигурацию jsoniter
var json = jsoniter.ConfigFastest

// MarshalJSON реализует интерфейс json.Marshaler для Geometry
func (g *Geometry) MarshalJSON() ([]byte, error) {
	if g.Type == "" || g.Coordinates == nil {
		return json.Marshal(nil)
	}

	geoInterface := map[string]any{
		"type":        g.Type,
		"coordinates": g.Coordinates,
	}

	return json.Marshal(geoInterface)
}

// UnmarshalJSON реализует интерфейс json.Unmarshaler для Geometry
func (g *Geometry) UnmarshalJSON(b []byte) error {
	if len(b) == 0 || string(b) == "null" {
		// Сбрасываем все поля при получении пустого значения
		g.Type = ""
		g.Coordinates = nil
		return nil
	}

	geoInterface := make(map[string]any)
	err := json.Unmarshal(b, &geoInterface)
	if err != nil {
		return err
	}

	// Проверяем наличие поля "type"
	typeVal, hasType := geoInterface["type"]
	if !hasType {
		return fmt.Errorf("missing required field 'type'")
	}

	// Безопасное приведение типа
	typeStr, ok := typeVal.(string)
	if !ok {
		return fmt.Errorf("field 'type' must be a string, got %T", typeVal)
	}

	// Проверяем наличие поля "coordinates"
	coordsVal, hasCoords := geoInterface["coordinates"]
	if !hasCoords {
		return fmt.Errorf("missing required field 'coordinates'")
	}

	// Сбрасываем все поля перед установкой новых
	g.Type = GeometryType(typeStr)
	g.Coordinates = nil

	switch g.Type {
	case GeometryPoint:
		point, err := decodePoint(coordsVal)
		if err != nil {
			return fmt.Errorf("failed to decode Point: %w", err)
		}
		g.Coordinates = point
	case GeometryLineString:
		lineString, err := decodeLine(coordsVal)
		if err != nil {
			return fmt.Errorf("failed to decode LineString: %w", err)
		}
		g.Coordinates = lineString
	case GeometryPolygon:
		polygon, err := decodePolygon(coordsVal)
		if err != nil {
			return fmt.Errorf("failed to decode Polygon: %w", err)
		}
		g.Coordinates = polygon
	case GeometryMultiPolygon:
		multiPolygon, err := decodeMultiPolygon(coordsVal)
		if err != nil {
			return fmt.Errorf("failed to decode MultiPolygon: %w", err)
		}
		g.Coordinates = multiPolygon
	default:
		return fmt.Errorf("unknown geometry type: %s", g.Type)
	}
	return nil
}

// decodePoint преобразует интерфейс в Point
func decodePoint(data any) (Point, error) {
	coords, ok := data.([]any)
	if !ok {
		return nil, fmt.Errorf("not a valid point, got %T", data)
	}

	result := make(Point, 0, len(coords))
	for _, coord := range coords {
		f, ok := coord.(float64)
		if !ok {
			return nil, fmt.Errorf("not a valid coordinate, expected float64, got %T", coord)
		}
		result = append(result, utils.Round(f, 9))
	}
	return result, nil
}

// decodeLine преобразует интерфейс в LineString
func decodeLine(data any) (LineString, error) {
	points, ok := data.([]any)
	if !ok {
		return nil, fmt.Errorf("not a valid set of positions, got %T", data)
	}

	result := make(LineString, 0, len(points))
	for i, po := range points {
		p, err := decodePoint(po)
		if err != nil {
			return nil, fmt.Errorf("error at position %d: %w", i, err)
		}
		result = append(result, p)
	}

	return result, nil
}

// decodePolygon преобразует интерфейс в Polygon
func decodePolygon(data any) (Polygon, error) {
	sets, ok := data.([]any)
	if !ok {
		return nil, fmt.Errorf("not a valid path, got %T", data)
	}

	result := make(Polygon, 0, len(sets))
	for i, set := range sets {
		s, err := decodeLine(set)
		if err != nil {
			return nil, fmt.Errorf("error in ring %d: %w", i, err)
		}
		result = append(result, s)
	}

	return result, nil
}

// decodeMultiPolygon преобразует интерфейс в MultiPolygon
func decodeMultiPolygon(data any) (MultiPolygon, error) {
	polygons, ok := data.([]any)
	if !ok {
		return nil, fmt.Errorf("not a valid polygon, got %T", data)
	}

	result := make(MultiPolygon, 0, len(polygons))
	for i, poly := range polygons {
		p, err := decodePolygon(poly)
		if err != nil {
			return nil, fmt.Errorf("error in polygon %d: %w", i, err)
		}
		result = append(result, p)
	}

	return result, nil
}
