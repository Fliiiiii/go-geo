package grid

import (
	"math"

	"github.com/Fliiiiii/go-geo/types"
)

// CreateHexagon создает шестиугольник с заданными параметрами
// lon - долгота центра шестиугольника
// lat - широта центра шестиугольника
// r - радиус шестиугольника
// k - коэффициент коррекции для проекции Меркатора
func CreateHexagon(lon, lat, r, k float64) types.Polygon {
	// Начальный угол (в радианах)
	var angle float64
	// Создаем массив для 7 точек (6 вершин + дублирование первой точки для замыкания)
	points := make([]types.Point, 7)
	for i := 0; i < 6; i++ {
		// Вычисляем координаты точки:
		// - Долгота: смещение от центра по синусу угла
		// - Широта: смещение от центра по косинусу угла с коррекцией Меркатора
		p := types.NewPoint((math.Sin(angle)*r + lon), (math.Cos(angle)*r*k + lat))
		// Увеличиваем угол на 60 градусов (π/3 радиан)
		angle += math.Pi / 3
		points[i] = p
	}
	// Замыкаем полигон, дублируя первую точку
	points[6] = points[0]
	return types.NewPolygon(types.NewLineString(points...))
}

// CreateHexagonalGrid создает сетку шестиугольников с заданными параметрами
// minLon, minLat - координаты юго-западного угла
// maxLon, maxLat - координаты северо-восточного угла
// r - радиус шестиугольника в градусах по широте
func CreateHexagonalGrid(minLon, minLat, maxLon, maxLat, r float64) types.MultiPolygon {
	var M types.MultiPolygon

	// Проверяем корректность границ
	if minLat > maxLat || minLon > maxLon {
		return M // Возвращаем пустую карту при некорректных границах
	}

	// Ограничиваем широту до диапазона [-90, 90]
	if minLat < -90 {
		minLat = -90
	}
	if maxLat > 90 {
		maxLat = 90
	}

	// Ограничиваем долготу до диапазона [-180, 180]
	if minLon < -180 {
		minLon = -180
	}
	if maxLon > 180 {
		maxLon = 180
	}

	// Шаг между центрами шестиугольников по широте
	deltaLat := r * 3
	// Шаг между центрами шестиугольников по долготе
	deltaLon := math.Sin(math.Pi/3) * r

	// Цикл по широте - создаем ряды шестиугольников
	for lat := minLat; lat <= maxLat; {
		// Вычисляем коэффициент Меркатора для текущей широты
		// Это нужно для компенсации искажений проекции (сжатие по долготе при удалении от экватора)
		deltaMerc := math.Cos(lat * math.Pi / 180)

		// флаг для смещения шестиугольников в четных/нечетных рядах (шахматный порядок)
		b := false

		// Цикл по долготе - создаем шестиугольники в текущем ряду
		for lon := minLon; lon <= maxLon; lon += deltaLon {
			var p types.Polygon
			// В зависимости от флага b, создаем шестиугольник со смещением или без
			if b {
				p = CreateHexagon(lon, lat, r, deltaMerc)
			} else {
				p = CreateHexagon(lon, lat+r*1.5*deltaMerc, r, deltaMerc)
			}

			M = append(M, p)
			// Инвертируем флаг для следующего шестиугольника в ряду
			b = !b
		}

		// Увеличиваем широту с учетом коэффициента Меркатора
		// Это компенсирует искажения проекции и обеспечивает равномерное распределение шестиугольников
		lat += deltaLat * deltaMerc
	}

	return M
}

// CreateHexagonalGridWithDensity создает гексагональную сетку с заданной плотностью
// minLon, minLat - координаты юго-западного угла
// maxLon, maxLat - координаты северо-восточного угла
// density - плотность сетки (количество шестиугольников на градус широты)
func CreateHexagonalGridWithDensity(minLon, minLat, maxLon, maxLat, density float64) types.MultiPolygon {
	// Рассчитываем радиус шестиугольников на основе плотности
	// Чем выше плотность, тем меньше радиус
	r := 1.0 / (density * 3)
	return CreateHexagonalGrid(minLat, maxLat, minLon, maxLon, r)
}
