package grid

import (
	"math"

	"github.com/Fliiiiii/go-geo/types"
)

// CreateRectangularGrid создает прямоугольную сетку точек с заданными границами и шагом
func CreateRectangularGrid(minLon, minLat, maxLon, maxLat, stepLon, stepLat float64) []types.Point {
	// Проверяем корректность границ
	if minLon > maxLon || minLat > maxLat || stepLon <= 0 || stepLat <= 0 {
		return []types.Point{}
	}

	// Оценка размера сетки для предварительного выделения памяти
	latCount := int(math.Ceil((maxLat-minLat)/stepLat)) + 1
	lonCount := int(math.Ceil((maxLon-minLon)/stepLon)) + 1
	estimatedSize := latCount * lonCount
	grid := make([]types.Point, 0, estimatedSize)

	// Максимальный коэффициент корректировки для предотвращения слишком больших шагов вблизи полюсов
	const maxAdjustmentFactor = 10.0

	// Учитываем, что длина одного градуса долготы зависит от широты
	for lat := minLat; lat <= maxLat; lat += stepLat {
		// Косинус широты для корректировки шага по долготе
		cosLat := math.Cos(lat * math.Pi / 180.0)

		// Предотвращаем деление на слишком маленькие значения
		cosLat = max(cosLat, 0.1) // Увеличен минимальный косинус для предотвращения слишком больших шагов

		// Корректируем шаг долготы с ограничением максимального фактора
		adjustmentFactor := 1.0 / cosLat
		adjustmentFactor = math.Min(adjustmentFactor, maxAdjustmentFactor)
		adjustedStepLon := stepLon * adjustmentFactor

		for lon := minLon; lon <= maxLon; lon += adjustedStepLon {
			// Убедимся, что не выходим за пределы maxLon
			if lon > maxLon {
				lon = maxLon
			}
			grid = append(grid, types.NewPoint(lon, lat))
		}
	}

	return grid
}

// CreateRectangularGridCells создает сетку полигонов (ячеек) с заданными границами и шагом
func CreateRectangularGridCells(minLon, minLat, maxLon, maxLat, stepLon, stepLat float64) types.MultiPolygon {
	// Проверяем корректность границ
	if minLon > maxLon || minLat > maxLat || stepLon <= 0 || stepLat <= 0 {
		return types.MultiPolygon{}
	}

	// Оценка размера сетки для предварительного выделения памяти
	latCount := int(math.Ceil((maxLat - minLat) / stepLat))
	lonCount := int(math.Ceil((maxLon - minLon) / stepLon))
	estimatedSize := latCount * lonCount
	gridCells := make(types.MultiPolygon, 0, estimatedSize)

	// Максимальный коэффициент корректировки
	const maxAdjustmentFactor = 10.0

	for lat := minLat; lat < maxLat; lat += stepLat {
		// Верхняя граница ячейки
		nextLat := lat + stepLat
		if nextLat > maxLat {
			nextLat = maxLat
		}

		// Косинус широты для корректировки шага по долготе
		cosLat := math.Cos(lat * math.Pi / 180.0)
		cosLat = max(cosLat, 0.1) // Увеличен минимальный косинус

		// Корректируем шаг долготы с ограничением
		adjustmentFactor := 1.0 / cosLat
		adjustmentFactor = math.Min(adjustmentFactor, maxAdjustmentFactor)
		adjustedStepLon := stepLon * adjustmentFactor

		for lon := minLon; lon < maxLon; lon += adjustedStepLon {
			// Правая граница ячейки
			nextLon := lon + adjustedStepLon
			if nextLon > maxLon {
				nextLon = maxLon
			}

			// Создаем полигон напрямую по координатам углов
			cell := types.NewPolygon(types.NewLineString(
				types.NewPoint(lon, lat),         // Нижний левый угол
				types.NewPoint(nextLon, lat),     // Нижний правый угол
				types.NewPoint(nextLon, nextLat), // Верхний правый угол
				types.NewPoint(lon, nextLat),     // Верхний левый угол
				types.NewPoint(lon, lat),         // Замыкаем полигон
			))
			gridCells = append(gridCells, cell)
		}
	}

	return gridCells
}
