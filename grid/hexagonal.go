package grid

import (
	"math"

	"github.com/Fliiiiii/go-geo/calc"
	"github.com/Fliiiiii/go-geo/types"
)

// CreateHexagonalGridCells создает гексагональную сетку в виде полигонов (шестиугольников) с учетом кривизны Земли
// minLon, minLat - координаты юго-западного угла
// maxLon, maxLat - координаты северо-восточного угла
// spacing - расстояние между центрами соседних шестиугольников (в метрах)
func CreateHexagonalGridCells(minLon, minLat, maxLon, maxLat, spacing float64) types.MultiPolygon {
	var gridCells types.MultiPolygon

	// Проверяем корректность границ
	if minLon > maxLon || minLat > maxLat || spacing <= 0 {
		return gridCells
	}

	// Центры гексагонов
	centers := createHexagonalGrid(minLon, minLat, maxLon, maxLat, spacing)

	// Если сетка пуста, возвращаем пустой результат
	if len(centers) == 0 {
		return gridCells
	}

	// Радиус шестиугольника
	radius := spacing / 2

	// Создаем полигон для каждого центра
	for _, center := range centers {
		var hexagonPoints types.LineString

		// Создаем 6 вершин шестиугольника, начиная с угла 30 градусов
		// для получения вертикально ориентированного шестиугольника
		for i := 0; i < 6; i++ {
			// Угол в радианах (начиная с 30 и с шагом 60 градусов)
			angle := (30.0 + float64(i)*60.0) * (math.Pi / 180.0)

			// Вычисляем точку назначения на расстоянии radius от центра в направлении angle
			vertex := calc.CalculateDestinationPoint(center, radius, angle)
			hexagonPoints = append(hexagonPoints, vertex)
		}

		// Замыкаем полигон, добавляя первую точку в конец
		hexagonPoints = append(hexagonPoints, hexagonPoints[0])

		// Создаем полигон и добавляем в сетку
		hexagon := types.NewPolygon(hexagonPoints)
		gridCells = append(gridCells, hexagon)
	}

	return gridCells
}

func createHexagonalGrid(minLon, minLat, maxLon, maxLat, spacing float64) []types.Point {
	var grid []types.Point

	// Проверяем корректность границ
	if minLon > maxLon || minLat > maxLat || spacing <= 0 {
		return grid
	}

	// Константы для гексагональной сетки
	const sine60 = 0.866025404 // sin(60°)

	// Шаг по горизонтали (расстояние между соседними точками в ряду)
	eastBearing := 90.0 * (math.Pi / 180.0) // восток (в радианах)

	// Вертикальный шаг между рядами
	rowHeight := spacing * sine60

	// Максимальное количество шагов (для предотвращения бесконечного цикла)
	maxSteps := 10000

	// Создаем ряды сетки
	currentRowY := minLat
	currentRow := 0

	for currentRowY <= maxLat && currentRow < maxSteps {
		// Определяем стартовую точку для текущего ряда
		var rowStartLon float64

		// Смещение для нечетных рядов
		if currentRow%2 == 1 {
			// Смещение на половину ширины ячейки для нечетных рядов
			halfSpacingLon := calc.CalculateDestinationPoint(
				types.NewPoint(minLon, currentRowY),
				spacing/2,
				eastBearing,
			).GetLongitude()
			rowStartLon = halfSpacingLon
		} else {
			rowStartLon = minLon
		}

		// Создаем точки для текущего ряда
		currentPoint := types.NewPoint(rowStartLon, currentRowY)
		grid = append(grid, currentPoint)

		stepsCount := 0
		for stepsCount < maxSteps {
			// Двигаемся на восток на расстояние spacing
			nextPoint := calc.CalculateDestinationPoint(currentPoint, spacing, eastBearing)

			// Проверяем, не вышли ли за границы
			if nextPoint.GetLongitude() > maxLon {
				break
			}

			grid = append(grid, nextPoint)
			currentPoint = nextPoint
			stepsCount++
		}

		// Переходим к следующему ряду
		northBearing := 0.0 * (math.Pi / 180.0) // север (в радианах)
		currentRowY = calc.CalculateDestinationPoint(
			types.NewPoint(minLon, currentRowY),
			rowHeight,
			northBearing,
		).GetLatitude()

		currentRow++
	}

	return grid
}
