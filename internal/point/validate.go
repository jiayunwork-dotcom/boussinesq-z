package point

import (
	"errors"
	"fmt"
	"math"
)

var errNonPositiveDepth = errors.New("depth must be positive")

func ValidateField(value float64, name string, min float64) error {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return fmt.Errorf("%s must be finite, got %g", name, value)
	}
	if value < min {
		return fmt.Errorf("%s must be >= %g, got %g", name, min, value)
	}
	return nil
}

func ValidateLoad(P, z, r float64) error {
	if P < 0 {
		return fmt.Errorf("P must be non-negative under compression-positive convention, got %g", P)
	}
	if z < 0 {
		return fmt.Errorf("z must be non-negative below the surface, got %g", z)
	}
	if r < 0 {
		return fmt.Errorf("r must be non-negative, got %g", r)
	}
	if P > 0 && z == 0 && r == 0 {
		return fmt.Errorf("R=0 singularity at z=0,r=0")
	}
	return nil
}

func ValidatePoisson(poisson float64) error {
	if poisson < 0 || poisson > 0.5 {
		return fmt.Errorf("poisson ratio must be 0..0.5, got %g", poisson)
	}
	return nil
}

func IsSingularity(z, r float64) bool {
	return z == 0 && r == 0
}

func IsBelowSurface(z float64) bool {
	return z > 0
}

func IsSurfacePoint(z, r float64) bool {
	return z == 0 && r > 0
}

func CheckCompressionSign(P float64) error {
	if P < 0 {
		return fmt.Errorf("P sign conflicts with compression-positive convention")
	}
	return nil
}

func CheckFinite(value float64, name string) error {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return fmt.Errorf("%s must be finite", name)
	}
	return nil
}

func IsDepthError(err error) bool {
	return errors.Is(err, errNonPositiveDepth)
}
