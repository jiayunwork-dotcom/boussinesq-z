package point

import (
	"fmt"
)

func NewtonsToKilonewtons(value float64) float64 {
	return value / 1000
}

func KilonewtonsToNewtons(value float64) float64 {
	return value * 1000
}

func PascalsToKilopascals(value float64) float64 {
	return value / 1000
}

func KilopascalsToPascals(value float64) float64 {
	return value * 1000
}

func MetersToMillimeters(value float64) float64 {
	return value * 1000
}

func MillimetersToMeters(value float64) float64 {
	return value / 1000
}

func FormatForce(value float64) string {
	return fmt.Sprintf("%.4g N", value)
}

func FormatStress(value float64) string {
	return fmt.Sprintf("%.4g Pa", value)
}

func FormatInfluence(value float64) string {
	return fmt.Sprintf("%.6g", value)
}

func ConvertResultToKPa(result Result) Result {
	result.Vertical = PascalsToKilopascals(result.Vertical)
	result.Radial = PascalsToKilopascals(result.Radial)
	result.Tangential = PascalsToKilopascals(result.Tangential)
	result.ShearRZ = PascalsToKilopascals(result.ShearRZ)
	return result
}

func ConvertLoadToKN(load Load) Load {
	load.P = NewtonsToKilonewtons(load.P)
	return load
}

func ScaleStress(result Result, factor float64) Result {
	result.Vertical *= factor
	result.Radial *= factor
	result.Tangential *= factor
	result.ShearRZ *= factor
	return result
}

func RoundStress(value float64, digits int) float64 {
	scale := 1.0
	for i := 0; i < digits; i++ {
		scale *= 10
	}
	return float64(int64(value*scale+0.5)) / scale
}

func SignOf(value float64) string {
	if value < 0 {
		return "-"
	}
	if value > 0 {
		return "+"
	}
	return "0"
}

func Unit(value float64) string {
	if value >= 1e6 {
		return "MPa"
	}
	if value >= 1e3 {
		return "kPa"
	}
	return "Pa"
}

func Display(value float64) string {
	unit := Unit(value)
	switch unit {
	case "MPa":
		return fmt.Sprintf("%.4g MPa", value/1e6)
	case "kPa":
		return fmt.Sprintf("%.4g kPa", value/1e3)
	default:
		return fmt.Sprintf("%.4g Pa", value)
	}
}
