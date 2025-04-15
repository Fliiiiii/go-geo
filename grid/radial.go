package grid

import (
	"math"

	"github.com/Fliiiiii/go-geo/calc"
	"github.com/Fliiiiii/go-geo/types"
)

// CreateRadialGrid создает радиальную сетку точек вокруг центральной точки
// centerLon, centerLat - координаты центра сетки
// maxRadius - максимальный радиус сетки в метрах
// numRings - количество концентрических окружностей
// pointsPerRing - количество точек на каждой окружности
func CreateRadialGrid(centerLon, centerLat, maxRadius float64, numRings, pointsPerRing int) []types.Point {
	var grid []types.Point

	// Проверяем корректность входных данных
	if maxRadius <= 0 || numRings <= 0 || pointsPerRing <= 0 {
		return grid
	}

	// Добавляем центральную точку
	centerPoint := types.NewPoint(centerLon, centerLat)
	grid = append(grid, centerPoint)

	// Шаг радиуса между кольцами
	radiusStep := maxRadius / float64(numRings)

	// Создаем точки на каждом кольце, учитывая кривизну Земли
	for ring := 1; ring <= numRings; ring++ {
		radius := float64(ring) * radiusStep

		for i := 0; i < pointsPerRing; i++ {
			// Вычисляем угол (азимут) для текущей точки в радианах
			bearing := 2 * math.Pi * float64(i) / float64(pointsPerRing)

			// Используем расчет с учетом кривизны Земли для определения координат точки
			point := calc.CalculateDestinationPoint(centerPoint, radius, bearing)
			grid = append(grid, point)
		}
	}

	return grid
}

// CreateRadialSectors создает радиальные секторы вокруг центральной точки
// centerLon, centerLat - координаты центра сетки
// maxRadius - максимальный радиус сетки в метрах
// numSectors - количество секторов
// numRings - количество концентрических окружностей
// Возвращает набор полигонов, представляющих секторы
func CreateRadialSectors(centerLon, centerLat, maxRadius float64, numSectors, numRings int) types.MultiPolygon {
	var sectors types.MultiPolygon

	// Проверяем корректность входных данных
	if maxRadius <= 0 || numSectors <= 0 || numRings <= 0 {
		return sectors
	}

	centerPoint := types.NewPoint(centerLon, centerLat)

	// Шаг радиуса между кольцами
	radiusStep := maxRadius / float64(numRings)
	// Угловой шаг между секторами
	angleStep := 2 * math.Pi / float64(numSectors)

	// Создаем секторы
	for sector := 0; sector < numSectors; sector++ {
		startAngle := float64(sector) * angleStep
		endAngle := float64(sector+1) * angleStep

		for ring := 0; ring < numRings; ring++ {
			innerRadius := float64(ring) * radiusStep
			outerRadius := float64(ring+1) * radiusStep

			// Создаем полигон для сектора
			var sectorPoints types.LineString

			if innerRadius == 0 {
				// Особый случай для первого кольца с центральной точкой
				sectorPoints = append(sectorPoints, centerPoint)

				// Добавляем точки на внешнем радиусе
				numPoints := int(math.Max(4, math.Ceil((endAngle-startAngle)*outerRadius/50)))
				for i := 0; i <= numPoints; i++ {
					angle := startAngle + (endAngle-startAngle)*float64(i)/float64(numPoints)
					point := calc.CalculateDestinationPoint(centerPoint, outerRadius, angle)
					sectorPoints = append(sectorPoints, point)
				}

				// Замыкаем полигон, добавляя первую точку внешнего радиуса
				sectorPoints = append(sectorPoints, sectorPoints[1])
			} else {
				// Добавляем точки на внутреннем радиусе
				numInnerPoints := int(math.Max(4, math.Ceil((endAngle-startAngle)*innerRadius/50)))
				for i := 0; i <= numInnerPoints; i++ {
					angle := startAngle + (endAngle-startAngle)*float64(i)/float64(numInnerPoints)
					point := calc.CalculateDestinationPoint(centerPoint, innerRadius, angle)
					sectorPoints = append(sectorPoints, point)
				}

				// Добавляем точки на внешнем радиусе (в обратном порядке)
				numOuterPoints := int(math.Max(4, math.Ceil((endAngle-startAngle)*outerRadius/50)))
				for i := numOuterPoints; i >= 0; i-- {
					angle := startAngle + (endAngle-startAngle)*float64(i)/float64(numOuterPoints)
					point := calc.CalculateDestinationPoint(centerPoint, outerRadius, angle)
					sectorPoints = append(sectorPoints, point)
				}

				// Замыкаем полигон
				sectorPoints = append(sectorPoints, sectorPoints[0])
			}

			// Проверяем, что полигон содержит минимум 4 точки (3 уникальные точки + замыкающая)
			if len(sectorPoints) >= 4 {
				sectors = append(sectors, types.Polygon{sectorPoints})
			}
		}
	}

	return sectors
}
