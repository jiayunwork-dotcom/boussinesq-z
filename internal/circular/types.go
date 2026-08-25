package circular

import (
	"fmt"
	"math"
)

type Load struct {
	Q float64 `json:"q"`
	A float64 `json:"a"`
	Z float64 `json:"z"`
	R float64 `json:"r"`
}

type Result struct {
	Q         float64 `json:"q"`
	A         float64 `json:"a"`
	Z         float64 `json:"z"`
	R         float64 `json:"r"`
	ZOverA    float64 `json:"z_over_a"`
	ROverA    float64 `json:"r_over_a"`
	Influence float64 `json:"influence_i"`
	Vertical  float64 `json:"sigma_z"`
}

func (l Load) Validate() error {
	if l.Q <= 0 {
		return fmt.Errorf("q must be positive, got %g", l.Q)
	}
	if l.A <= 0 {
		return fmt.Errorf("a must be positive, got %g", l.A)
	}
	if l.Z < 0 {
		return fmt.Errorf("z must be non-negative, got %g", l.Z)
	}
	if l.R < 0 {
		return fmt.Errorf("r must be non-negative, got %g", l.R)
	}
	if math.IsNaN(l.Q) || math.IsNaN(l.A) || math.IsNaN(l.Z) || math.IsNaN(l.R) ||
		math.IsInf(l.Q, 0) || math.IsInf(l.A, 0) || math.IsInf(l.Z, 0) || math.IsInf(l.R, 0) {
		return fmt.Errorf("q, a, z, r must be finite numbers")
	}
	return nil
}

func (l Load) Normalized() (float64, float64) {
	return l.Z / l.A, l.R / l.A
}

func (l Load) String() string {
	return fmt.Sprintf("q=%.6g a=%.6g z=%.6g r=%.6g", l.Q, l.A, l.Z, l.R)
}

func Scale(l Load, factor float64) Load {
	l.Q *= factor
	return l
}

func AtDepth(l Load, depth float64) Load {
	l.Z = depth
	return l
}

func AtRadius(l Load, radius float64) Load {
	l.R = radius
	return l
}

func IsCenterline(l Load) bool {
	return l.R == 0
}

func IsInsideFootprint(l Load) bool {
	return l.R <= l.A
}

func DepthRatio(l Load) float64 {
	if l.A == 0 {
		return 0
	}
	return l.Z / l.A
}

func RadiusRatio(l Load) float64 {
	if l.A == 0 {
		return 0
	}
	return l.R / l.A
}

func PressureResult(l Load, influence float64) Result {
	zOverA, rOverA := l.Normalized()
	return Result{
		Q:         l.Q,
		A:         l.A,
		Z:         l.Z,
		R:         l.R,
		ZOverA:    zOverA,
		ROverA:    rOverA,
		Influence: influence,
		Vertical:  l.Q * influence,
	}
}

var lastVertical float64

func accumulateVertical(result Result) Result {
	combined := lastVertical + result.Vertical
	lastVertical = combined
	result.Vertical = combined
	return result
}
