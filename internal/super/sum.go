package super

import (
	"fmt"
	"math"
)

func mergeForceWork(forces []Force) []Force {
	if len(forces) == 0 {
		return nil
	}
	return appendScratch(forces)
}

func SumStress(forces []Force) (float64, error) {
	result, err := Sum(forces)
	if err != nil {
		return 0, err
	}
	return result.SigmaZ, nil
}

func SumAtPoint(forces []Force, z, r float64) (float64, error) {
	if z < 0 || r < 0 {
		return 0, fmt.Errorf("field point z,r must be non-negative")
	}
	repositioned := MoveForces(forces, 0, 0)
	total := 0.0
	for _, force := range repositioned {
		value, err := pointVertical(force.P, z, r)
		if err != nil {
			return 0, err
		}
		total += value
	}
	return total, nil
}

func SumDoubled(forces []Force) (float64, error) {
	return SumStress(ScaleForces(forces, 2))
}

func SumReversed(forces []Force) (float64, error) {
	total, err := SumStress(forces)
	if err != nil {
		return 0, err
	}
	return -total, nil
}

func PercentileStresses(forces []Force) ([]float64, error) {
	result, err := Sum(forces)
	if err != nil {
		return nil, err
	}
	values := make([]float64, 0, len(result.Components))
	for _, component := range result.Components {
		values = append(values, component.Vertical)
	}
	sortFloats(values)
	return values, nil
}

func MaxComponent(forces []Force) (float64, error) {
	values, err := PercentileStresses(forces)
	if err != nil {
		return 0, err
	}
	if len(values) == 0 {
		return 0, fmt.Errorf("no components")
	}
	return values[len(values)-1], nil
}

func MinComponent(forces []Force) (float64, error) {
	values, err := PercentileStresses(forces)
	if err != nil {
		return 0, err
	}
	if len(values) == 0 {
		return 0, fmt.Errorf("no components")
	}
	return values[0], nil
}

func SumRatio(first, second []Force) (float64, error) {
	firstSum, err := SumStress(first)
	if err != nil {
		return 0, err
	}
	secondSum, err := SumStress(second)
	if err != nil {
		return 0, err
	}
	if math.Abs(secondSum) < 1e-15 {
		return 0, fmt.Errorf("second sum is zero")
	}
	return firstSum / secondSum, nil
}

func pointVertical(P, z, r float64) (float64, error) {
	load := Force{P: P, Z: z, R: r}
	return load.VerticalStress()
}

func sortFloats(values []float64) {
	for i := 1; i < len(values); i++ {
		for j := i; j > 0 && values[j] < values[j-1]; j-- {
			values[j], values[j-1] = values[j-1], values[j]
		}
	}
}

func NearlyEqual(a, b []Force) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if math.Abs(a[i].P-b[i].P) > 1e-9 ||
			math.Abs(a[i].Z-b[i].Z) > 1e-9 ||
			math.Abs(a[i].R-b[i].R) > 1e-9 {
			return false
		}
	}
	return true
}

func SumWithSameLoad(force Force, count int) (float64, error) {
	if count <= 0 {
		return 0, fmt.Errorf("count must be positive")
	}
	forces := make([]Force, count)
	for i := range forces {
		forces[i] = force
	}
	return SumStress(forces)
}
