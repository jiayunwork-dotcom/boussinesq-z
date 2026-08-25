package point

import (
	"encoding/json"
	"fmt"
	"math"
)

type Load struct {
	P float64 `json:"P"`
	Z float64 `json:"z"`
	R float64 `json:"r"`
}

type Result struct {
	P           float64    `json:"P"`
	Z           float64    `json:"z"`
	R           float64    `json:"r"`
	Radius      float64    `json:"R"`
	Influence   float64    `json:"influence_i"`
	Vertical    float64    `json:"sigma_z"`
	Radial      float64    `json:"sigma_r,omitempty"`
	Tangential  float64    `json:"sigma_theta,omitempty"`
	ShearRZ     float64    `json:"tau_rz,omitempty"`
	Poisson     float64    `json:"poisson,omitempty"`
	SourcePoint [3]float64 `json:"source_point,omitempty"`
}

func (l Load) Validate() error {
	if l.P < 0 {
		return fmt.Errorf("P must be non-negative under compression-positive convention, got %g", l.P)
	}
	if l.Z < 0 {
		return fmt.Errorf("z must be non-negative below the surface, got %g", l.Z)
	}
	if l.R < 0 {
		return fmt.Errorf("r must be non-negative, got %g", l.R)
	}
	if l.P > 0 && l.Z == 0 && l.R == 0 {
		return fmt.Errorf("R=0 singularity at z=0,r=0")
	}
	if math.IsNaN(l.P) || math.IsNaN(l.Z) || math.IsNaN(l.R) ||
		math.IsInf(l.P, 0) || math.IsInf(l.Z, 0) || math.IsInf(l.R, 0) {
		return fmt.Errorf("P, z, and r must be finite numbers")
	}
	return nil
}

func (l Load) Radius() float64 {
	return math.Sqrt(l.R*l.R + l.Z*l.Z)
}

func (l Load) Influence() (float64, error) {
	if err := l.Validate(); err != nil {
		return 0, err
	}
	if l.P == 0 {
		return 0, nil
	}
	radius := l.Radius()
	return 3 * math.Pow(l.Z, 3) / (2 * math.Pi * math.Pow(radius, 5)), nil
}

func (l Load) VerticalStress() (float64, error) {
	influence, err := l.Influence()
	if err != nil {
		return 0, err
	}
	return l.P * influence, nil
}

func (l Load) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]float64{"P": l.P, "z": l.Z, "r": l.R})
}

func (l Load) String() string {
	return fmt.Sprintf("P=%.6g z=%.6g r=%.6g", l.P, l.Z, l.R)
}

func Normalize(l Load) Load {
	if l.R < 0 {
		l.R = -l.R
	}
	return l
}

func IsCenterline(l Load) bool {
	return l.R == 0
}

func MaxInfluenceAtDepth(z float64) (float64, error) {
	if z <= 0 {
		return 0, fmt.Errorf("depth must be positive, got %g", z)
	}
	return 3 / (2 * math.Pi * z * z), nil
}

func DepthForInfluence(P, target float64) (float64, error) {
	if P <= 0 {
		return 0, fmt.Errorf("P must be positive")
	}
	if target <= 0 {
		return 0, fmt.Errorf("target stress must be positive")
	}
	return math.Sqrt(3 * P / (2 * math.Pi * target)), nil
}
