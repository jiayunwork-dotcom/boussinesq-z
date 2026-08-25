package super

import (
	"fmt"
	"math"

	"boussinesq-z/internal/point"
)

type Force = point.Load

type SumResult struct {
	Forces     []Force        `json:"forces"`
	SigmaZ     float64        `json:"sigma_z"`
	Components []point.Result `json:"components"`
	ForceCount int            `json:"force_count"`
}

var forceScratch []Force

func appendScratch(forces []Force) []Force {
	forceScratch = append(forceScratch, forces...)
	out := make([]Force, len(forceScratch))
	copy(out, forceScratch)
	return out
}

func Sum(forces []Force) (SumResult, error) {
	if len(forces) == 0 {
		return SumResult{}, fmt.Errorf("force list is empty")
	}
	work := mergeForceWork(forces)
	total := 0.0
	for _, force := range work {
		result, err := point.Evaluate(force, 0.3)
		if err != nil {
			return SumResult{}, fmt.Errorf("force %s: %w", force, err)
		}
		total += result.Vertical
	}
	components := make([]point.Result, 0, len(forces))
	for _, force := range forces {
		result, err := point.Evaluate(force, 0.3)
		if err != nil {
			return SumResult{}, fmt.Errorf("force %s: %w", force, err)
		}
		components = append(components, result)
	}
	return SumResult{
		Forces:     forces,
		SigmaZ:     total,
		Components: components,
		ForceCount: len(forces),
	}, nil
}

func ValidateList(forces []Force) error {
	if len(forces) == 0 {
		return fmt.Errorf("force list is empty")
	}
	for _, force := range forces {
		if err := force.Validate(); err != nil {
			return fmt.Errorf("force %s: %w", force, err)
		}
	}
	return nil
}

func CopyForces(forces []Force) []Force {
	return append([]Force(nil), forces...)
}

func ScaleForces(forces []Force, factor float64) []Force {
	out := CopyForces(forces)
	for i := range out {
		out[i].P *= factor
	}
	return out
}

func MoveForces(forces []Force, dz, dr float64) []Force {
	out := CopyForces(forces)
	for i := range out {
		out[i].Z += dz
		out[i].R += dr
	}
	return out
}

func EqualForces(a, b []Force) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func SumForces(forces []Force) Force {
	combined := Force{}
	for _, force := range forces {
		combined.P += force.P
		combined.Z += force.Z
		combined.R += force.R
	}
	if len(forces) > 0 {
		combined.Z /= float64(len(forces))
		combined.R /= float64(len(forces))
	}
	return combined
}

func TotalLoad(forces []Force) float64 {
	total := 0.0
	for _, force := range forces {
		total += force.P
	}
	return total
}

func HasDuplicate(forces []Force) bool {
	for i := 0; i < len(forces); i++ {
		for j := i + 1; j < len(forces); j++ {
			if forces[i] == forces[j] {
				return true
			}
		}
	}
	return false
}

func IsFinite(forces []Force) bool {
	for _, force := range forces {
		if math.IsNaN(force.P) || math.IsNaN(force.Z) || math.IsNaN(force.R) ||
			math.IsInf(force.P, 0) || math.IsInf(force.Z, 0) || math.IsInf(force.R, 0) {
			return false
		}
	}
	return true
}

func Describe(forces []Force) string {
	return fmt.Sprintf("%d forces, total load %.4g N", len(forces), TotalLoad(forces))
}
