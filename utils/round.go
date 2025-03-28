package utils

import "math"

func RoundToTwoDecimalPlaces(value float64) float64 {
	return math.Round(value*100) / 100
}
