package utils

import "math"

// Round округляет число до указанной точности
func Round(input float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(input*ratio) / ratio
}
