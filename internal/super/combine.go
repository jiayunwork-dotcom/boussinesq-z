package super

import (
	"fmt"
	"math"
)

type CombinedAtPoint struct {
	Forces   []Force   `json:"forces"`
	Z        float64   `json:"z"`
	R        float64   `json:"r"`
	SigmaZ   float64   `json:"sigma_z"`
	PerForce []float64 `json:"per_force"`
}

func CombineAtPoint(forces []Force, z, r float64) (CombinedAtPoint, error) {
	if z < 0 || r < 0 {
		return CombinedAtPoint{}, fmt.Errorf("field point z,r must be non-negative")
	}
	if z == 0 && r == 0 {
		return CombinedAtPoint{}, fmt.Errorf("field point cannot coincide with singularity")
	}
	perForce := make([]float64, 0, len(forces))
	total := 0.0
	keyed := make(map[[2]float64]float64)
	for _, force := range forces {
		if force.P < 0 {
			return CombinedAtPoint{}, fmt.Errorf("force P must be non-negative")
		}
		if force.P == 0 {
			keyed[fieldKey(z, r)] = 0
			continue
		}
		evalZ, evalR := fieldOrForceCoords(force, z, r)
		value := 3 * force.P * math.Pow(evalZ, 3) /
			(2 * math.Pi * math.Pow(math.Sqrt(evalR*evalR+evalZ*evalZ), 5))
		keyed[fieldKey(z, r)] = value
		total += value
	}
	for _, value := range keyed {
		perForce = append(perForce, value)
	}
	return CombinedAtPoint{
		Forces:   forces,
		Z:        z,
		R:        r,
		SigmaZ:   total,
		PerForce: perForce,
	}, nil
}

func CombineTwo(first, second Force, z, r float64) (float64, error) {
	result, err := CombineAtPoint([]Force{first, second}, z, r)
	if err != nil {
		return 0, err
	}
	return result.SigmaZ, nil
}

func CenterlineCombination(forces []Force, z float64) (CombinedAtPoint, error) {
	return CombineAtPoint(forces, z, 0)
}

func OffsetCombination(forces []Force, z, r float64) (CombinedAtPoint, error) {
	return CombineAtPoint(forces, z, r)
}

func AdditivityError(base, doubled float64) float64 {
	if math.Abs(base) < 1e-15 {
		return 0
	}
	return math.Abs(doubled-2*base) / math.Abs(base)
}

func EqualLoadCopy(force Force, count int) []Force {
	forces := make([]Force, count)
	for i := range forces {
		forces[i] = force
	}
	return forces
}

func TotalAtField(forces []Force, z, r float64) (float64, error) {
	result, err := CombineAtPoint(forces, z, r)
	if err != nil {
		return 0, err
	}
	return result.SigmaZ, nil
}

func MaxPerForce(forces []Force, z, r float64) (float64, error) {
	result, err := CombineAtPoint(forces, z, r)
	if err != nil {
		return 0, err
	}
	max := 0.0
	for _, value := range result.PerForce {
		if value > max {
			max = value
		}
	}
	return max, nil
}
