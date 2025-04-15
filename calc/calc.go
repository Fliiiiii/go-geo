package calc

import (
	"math"

	"github.com/Fliiiiii/go-geo/types"
)

// Константы для радиуса Земли
const (
	EarthRadiusKm     = 6371.0
	EarthRadiusMeters = 6371000.0
)

// CalculateDistanceAndBearing вычисляет расстояние (в метрах) и азимут (в радианах)
// между двумя точками на сфере. Азимут находится в диапазоне [0, 2π).
func CalculateDistanceAndBearing(p1, p2 types.Point) (distance, bearing float64) {
	// Перевод в радианы
	lon1Rad := p1.GetLongitude() * math.Pi / 180.0
	lat1Rad := p1.GetLatitude() * math.Pi / 180.0
	lon2Rad := p2.GetLongitude() * math.Pi / 180.0
	lat2Rad := p2.GetLatitude() * math.Pi / 180.0

	// Разница координат
	dLon := lon2Rad - lon1Rad

	// Вычисление азимута
	y := math.Sin(dLon) * math.Cos(lat2Rad)
	x := math.Cos(lat1Rad)*math.Sin(lat2Rad) - math.Sin(lat1Rad)*math.Cos(lat2Rad)*math.Cos(dLon)
	bearing = math.Atan2(y, x)
	bearing = math.Mod(bearing+2*math.Pi, 2*math.Pi)

	// Вычисление расстояния по формуле гаверсинусов
	distance = calculateHaversineDistance(lat1Rad, lon1Rad, lat2Rad, lon2Rad, EarthRadiusMeters)

	return distance, bearing
}

// calculateHaversineDistance вычисляет расстояние между двумя точками по формуле гаверсинусов
// принимает координаты в радианах и радиус в нужных единицах измерения
func calculateHaversineDistance(lat1Rad, lon1Rad, lat2Rad, lon2Rad, radius float64) float64 {
	dLat := lat2Rad - lat1Rad
	dLon := lon2Rad - lon1Rad

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return radius * c
}

// CalculateDestinationPoint вычисляет точку назначения на сфере по начальным координатам,
// расстоянию (в метрах) и азимуту (в радианах)
func CalculateDestinationPoint(p types.Point, distance, bearing float64) types.Point {
	// Перевод в радианы
	lonRad := p.GetLongitude() * math.Pi / 180.0
	latRad := p.GetLatitude() * math.Pi / 180.0
	bearingRad := bearing

	// Вычисление углового расстояния
	angularDistance := distance / EarthRadiusMeters

	// Предварительные вычисления для оптимизации
	sinLatRad := math.Sin(latRad)
	cosLatRad := math.Cos(latRad)
	sinAngularDist := math.Sin(angularDistance)
	cosAngularDist := math.Cos(angularDistance)
	cosBearingRad := math.Cos(bearingRad)
	sinBearingRad := math.Sin(bearingRad)

	// Вычисление координат точки назначения
	destLatRad := math.Asin(sinLatRad*cosAngularDist +
		cosLatRad*sinAngularDist*cosBearingRad)

	destLonRad := lonRad + math.Atan2(sinBearingRad*sinAngularDist*cosLatRad,
		cosAngularDist-sinLatRad*math.Sin(destLatRad))

	// Нормализация долготы до диапазона [-π, π]
	destLonRad = math.Mod(destLonRad+3*math.Pi, 2*math.Pi) - math.Pi

	// Перевод обратно в градусы
	return types.NewPoint(destLonRad*180.0/math.Pi, destLatRad*180.0/math.Pi)
}

// CalculateDistance вычисляет расстояние между двумя точками (в километрах)
func CalculateDistance(p1, p2 types.Point) float64 {
	lat1Rad := p1.GetLatitude() * math.Pi / 180.0
	lon1Rad := p1.GetLongitude() * math.Pi / 180.0
	lat2Rad := p2.GetLatitude() * math.Pi / 180.0
	lon2Rad := p2.GetLongitude() * math.Pi / 180.0

	return calculateHaversineDistance(lat1Rad, lon1Rad, lat2Rad, lon2Rad, EarthRadiusKm)
}

// BoundingBox представляет ограничивающий прямоугольник с координатами в градусах
type BoundingBox struct {
	// MinLon - минимальная долгота (западная граница)
	MinLon float64
	// MinLat - минимальная широта (южная граница)
	MinLat float64
	// MaxLon - максимальная долгота (восточная граница)
	MaxLon float64
	// MaxLat - максимальная широта (северная граница)
	MaxLat float64
}

// CalculateBoundingBox вычисляет ограничивающий прямоугольник для полигона,
// учитывая возможное пересечение 180-го меридиана
func CalculateBoundingBox(p types.Polygon) BoundingBox {
	if len(p) == 0 || len(p[0]) == 0 {
		return BoundingBox{}
	}

	minLon := p[0][0].GetLongitude()
	minLat := p[0][0].GetLatitude()
	maxLon := p[0][0].GetLongitude()
	maxLat := p[0][0].GetLatitude()

	crossesAntimeridian := false
	var lons []float64

	// Собираем все долготы для проверки пересечения 180-го меридиана
	for _, lineString := range p {
		for _, point := range lineString {
			lon := point.GetLongitude()
			lat := point.GetLatitude()
			lons = append(lons, lon)

			if lat < minLat {
				minLat = lat
			}
			if lat > maxLat {
				maxLat = lat
			}
		}
	}

	// Проверяем, пересекает ли полигон 180-й меридиан
	for i := 1; i < len(lons); i++ {
		if math.Abs(lons[i]-lons[i-1]) > 180 {
			crossesAntimeridian = true
			break
		}
	}

	// Если пересечение есть, корректируем долготы
	if crossesAntimeridian {
		east := -180.0
		west := 180.0

		for _, lineString := range p {
			for _, point := range lineString {
				lon := point.GetLongitude()
				// Нормализуем долготу в диапазоне [-180, 180]
				if lon < 0 {
					if lon > east {
						east = lon
					}
				} else {
					if lon < west {
						west = lon
					}
				}
			}
		}

		minLon = west
		maxLon = east
	} else {
		// Стандартный расчет без пересечения 180-го меридиана
		for _, lineString := range p {
			for _, point := range lineString {
				lon := point.GetLongitude()
				if lon < minLon {
					minLon = lon
				}
				if lon > maxLon {
					maxLon = lon
				}
			}
		}
	}

	return BoundingBox{
		MinLon: minLon,
		MinLat: minLat,
		MaxLon: maxLon,
		MaxLat: maxLat,
	}
}

// CalculateLineStringLength вычисляет длину линии (в километрах)
func CalculateLineStringLength(ls types.LineString) float64 {
	if len(ls) < 2 {
		return 0
	}

	length := 0.0
	for i := 0; i < len(ls)-1; i++ {
		length += CalculateDistance(ls[i], ls[i+1])
	}

	return length
}

// CalculatePolygonArea вычисляет площадь полигона (в квадратных километрах)
// с использованием формулы сферического избытка для более точного расчета на сфере
func CalculatePolygonArea(p types.Polygon) float64 {
	if len(p) == 0 {
		return 0
	}

	// Суммарная площадь
	totalArea := 0.0

	// Рассчитываем площадь внешнего кольца
	if len(p[0]) >= 3 {
		outerArea := calculateRingArea(p[0])
		totalArea += math.Abs(outerArea)
	}

	// Вычитаем площади внутренних колец (дыр)
	for i := 1; i < len(p); i++ {
		if len(p[i]) >= 3 {
			holeArea := calculateRingArea(p[i])
			totalArea -= math.Abs(holeArea)
		}
	}

	return totalArea
}

// calculateRingArea вычисляет площадь одного кольца полигона с использованием формулы сферического избытка
func calculateRingArea(ring types.LineString) float64 {
	if len(ring) < 3 {
		return 0
	}

	// Используем формулу сферического избытка для более точного расчета на сфере
	area := 0.0
	n := len(ring)

	for i := 0; i < n; i++ {
		j := (i + 1) % n
		k := (i + 2) % n

		lon1 := ring[i].GetLongitude() * math.Pi / 180.0
		lat1 := ring[i].GetLatitude() * math.Pi / 180.0
		lon2 := ring[j].GetLongitude() * math.Pi / 180.0
		lat2 := ring[j].GetLatitude() * math.Pi / 180.0
		lon3 := ring[k].GetLongitude() * math.Pi / 180.0
		lat3 := ring[k].GetLatitude() * math.Pi / 180.0

		// Вычисляем сферический избыток для треугольника
		angle := math.Mod(math.Atan2(math.Sin(lon2-lon1)*math.Cos(lat2),
			math.Cos(lat1)*math.Sin(lat2)-math.Sin(lat1)*math.Cos(lat2)*math.Cos(lon2-lon1))+
			math.Atan2(math.Sin(lon3-lon2)*math.Cos(lat3),
				math.Cos(lat2)*math.Sin(lat3)-math.Sin(lat2)*math.Cos(lat3)*math.Cos(lon3-lon2))+
			math.Atan2(math.Sin(lon1-lon3)*math.Cos(lat1),
				math.Cos(lat3)*math.Sin(lat1)-math.Sin(lat3)*math.Cos(lat1)*math.Cos(lon1-lon3)),
			2*math.Pi)

		if angle > math.Pi {
			angle = 2*math.Pi - angle
		}

		area += angle
	}

	// Формула для площади на сфере: E * R²
	// Где E - сферический избыток, R - радиус
	return (area - (float64(n-2) * math.Pi)) * EarthRadiusKm * EarthRadiusKm
}

// PointInPolygon проверяет, находится ли точка внутри полигона,
// учитывая внешний контур и внутренние кольца (дыры)
func PointInPolygon(polygon types.Polygon, point types.Point) bool {
	if len(polygon) == 0 || len(polygon[0]) < 3 {
		return false
	}

	// Сначала проверяем, находится ли точка внутри внешнего контура
	if !pointInRing(polygon[0], point) {
		return false
	}

	// Затем проверяем, что точка не находится внутри какой-либо дыры
	for i := 1; i < len(polygon); i++ {
		if len(polygon[i]) >= 3 && pointInRing(polygon[i], point) {
			return false
		}
	}

	return true
}

// pointInRing проверяет, находится ли точка внутри отдельного кольца полигона
// с использованием усовершенствованного алгоритма ray casting
func pointInRing(ring types.LineString, point types.Point) bool {
	if len(ring) < 3 {
		return false
	}

	inside := false
	x := point.GetLongitude()
	y := point.GetLatitude()

	// Оптимизированный алгоритм ray casting
	n := len(ring)
	j := n - 1

	for i := 0; i < n; i++ {
		xi := ring[i].GetLongitude()
		yi := ring[i].GetLatitude()
		xj := ring[j].GetLongitude()
		yj := ring[j].GetLatitude()

		// Проверяем пересечение горизонтального луча с ребром полигона
		intersect := ((yi > y) != (yj > y)) &&
			(x < (xj-xi)*(y-yi)/(yj-yi)+xi)

		if intersect {
			inside = !inside
		}

		j = i
	}

	return inside
}
