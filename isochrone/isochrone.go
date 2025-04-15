// package isochrone Вообще много чего надо переписать
// Ооооочень тестовый вариант
package isochrone

import (
	"math"

	"github.com/Fliiiiii/go-geo/types"
)

// Isochrone представляет изохрону - полигон, представляющий область,
// достижимую из исходной точки за указанное время
type Isochrone struct {
	Origin     types.Point            // Исходная точка
	Duration   int                    // Время в секундах
	Polygon    types.Polygon          // Границы изохроны в виде полигона
	Properties map[string]interface{} // Дополнительные свойства изохроны
}

// IsochronesCollection представляет собой набор изохрон
type IsochronesCollection []Isochrone

// IsochroneParams содержит параметры для генерации изохроны
type IsochroneParams struct {
	Origin     types.Point // Исходная точка
	Durations  []int       // Времена в секундах, для которых нужно построить изохроны
	Resolution int         // Разрешение изохроны (кол-во точек на контуре)
	MaxSpeed   float64     // Максимальная скорость передвижения (в м/с)
}

// EarthRadiusM - радиус Земли в метрах
const EarthRadiusM = 6371000.0

// GenerateIsochrones генерирует изохроны вокруг заданной точки
func GenerateIsochrones(params IsochroneParams) (IsochronesCollection, error) {
	// Тут код для генерации изохрон
	// Надо использовать внешние сервисы для учета дорожной сети
	// Но буду делать просто круги разного радиуса)

	result := make(IsochronesCollection, 0, len(params.Durations))

	for _, duration := range params.Durations {
		iso := generateCircularIsochrone(params.Origin, duration, params.MaxSpeed, params.Resolution)
		result = append(result, iso)
	}

	return result, nil
}

// generateCircularIsochrone создает круговую изохрону (упрощенная модель)
func generateCircularIsochrone(origin types.Point, duration int, maxSpeed float64, resolution int) Isochrone {
	// Вычисляем радиус круга в метрах
	speed := maxSpeed                      // м/с
	distanceM := float64(duration) * speed // расстояние в метрах

	// Создаем круговой полигон с учетом кривизны Земли
	polygon := createCircle(origin, distanceM, resolution)

	return Isochrone{
		Origin:   origin,
		Duration: duration,
		Polygon:  polygon,
		Properties: map[string]interface{}{
			"duration_seconds": duration,
			"max_speed":        maxSpeed,
			"distance_meters":  distanceM,
		},
	}
}

// createCircle создает круговой полигон с учетом кривизны Земли
func createCircle(center types.Point, radiusM float64, numPoints int) types.Polygon {
	points := make(types.LineString, numPoints+1)

	// Угловое расстояние в радианах на сфере
	angularDistance := radiusM / EarthRadiusM

	// Преобразуем координаты центра в радианы
	latCenter := center.GetLatitude() * math.Pi / 180.0
	lonCenter := center.GetLongitude() * math.Pi / 180.0

	for i := range numPoints {
		// Азимут (направление от центра)
		azimuth := 2.0 * math.Pi * float64(i) / float64(numPoints)

		// Вычисляем новую точку с использованием сферической тригонометрии
		latPoint := math.Asin(math.Sin(latCenter)*math.Cos(angularDistance) +
			math.Cos(latCenter)*math.Sin(angularDistance)*math.Cos(azimuth))

		lonPoint := lonCenter + math.Atan2(
			math.Sin(azimuth)*math.Sin(angularDistance)*math.Cos(latCenter),
			math.Cos(angularDistance)-math.Sin(latCenter)*math.Sin(latPoint),
		)

		// Преобразуем обратно в градусы
		lat := latPoint * 180.0 / math.Pi
		lon := lonPoint * 180.0 / math.Pi

		// Нормализуем долготу (-180 до 180)
		if lon > 180.0 {
			lon -= 360.0
		} else if lon < -180.0 {
			lon += 360.0
		}

		points[i] = types.NewPoint(lon, lat)
	}

	// Замыкаем полигон
	points[numPoints] = points[0]

	return types.Polygon{points}
}

// ToGeoJSON преобразует изохрону в структуру Feature формата GeoJSON
func (i Isochrone) ToGeoJSON() map[string]interface{} {
	properties := i.Properties
	if properties == nil {
		properties = make(map[string]interface{})
	}

	// Добавляем базовые свойства если их нет
	if _, exists := properties["duration_seconds"]; !exists {
		properties["duration_seconds"] = i.Duration
	}

	return map[string]interface{}{
		"type": "Feature",
		"geometry": map[string]interface{}{
			"type":        "Polygon",
			"coordinates": i.Polygon,
		},
		"properties": properties,
	}
}

// GetIsochroneForDuration возвращает изохрону для указанного времени
func (ic IsochronesCollection) GetIsochroneForDuration(duration int) (Isochrone, bool) {
	for _, iso := range ic {
		if iso.Duration == duration {
			return iso, true
		}
	}
	return Isochrone{}, false
}

// MergeIsochrones объединяет несколько изохрон в один MultiPolygon
func MergeIsochrones(isochrones []Isochrone) types.MultiPolygon {
	result := make(types.MultiPolygon, 0, len(isochrones))

	for _, iso := range isochrones {
		result = append(result, iso.Polygon)
	}

	return result
}
