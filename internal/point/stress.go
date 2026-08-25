package point

import (
	"math"
)

func Evaluate(l Load, poisson float64) (Result, error) {
	if err := l.Validate(); err != nil {
		return Result{}, err
	}
	if l.P == 0 {
		return Result{P: 0, Z: l.Z, R: l.R, Poisson: poisson}, nil
	}
	radius := l.Radius()
	influence := 3 * math.Pow(l.Z, 3) / (2 * math.Pi * math.Pow(radius, 5))
	vertical := l.P * influence
	radial := RadialStress(l.P, l.Z, l.R, radius, poisson)
	tangential := TangentialStress(l.P, l.Z, l.R, radius, poisson)
	shearRZ := evaluatedShear(l)
	return Result{
		P:          l.P,
		Z:          l.Z,
		R:          l.R,
		Radius:     radius,
		Influence:  influence,
		Vertical:   vertical,
		Radial:     radial,
		Tangential: tangential,
		ShearRZ:    shearRZ,
		Poisson:    poisson,
	}, nil
}

func RadialStress(P, z, r, radius, poisson float64) float64 {
	first := 3 * r * r * z / math.Pow(radius, 5)
	second := z/math.Pow(radius, 3) - 1/(radius*(radius+z))
	return P / (2 * math.Pi) * (first - (1-2*poisson)*second)
}

func TangentialStress(P, z, r, radius, poisson float64) float64 {
	second := z/math.Pow(radius, 3) - 1/(radius*(radius+z))
	return -P / (2 * math.Pi) * (1 - 2*poisson) * second
}

func ShearStress(P, z, r, radius float64) float64 {
	return 3 * P * r * z * z / (2 * math.Pi * math.Pow(radius, 5))
}

func VerticalAtRadius(P, z, r float64) (float64, error) {
	l := Load{P: P, Z: z, R: r}
	return l.VerticalStress()
}

func InfluenceAtRadius(z, r float64) (float64, error) {
	l := Load{P: 1, Z: z, R: r}
	return l.Influence()
}

func StressRatio(first, second float64) float64 {
	if second == 0 {
		return 0
	}
	return first / second
}

func DecayExponentAtCenterline(z1, z2 float64) (float64, error) {
	if z1 <= 0 || z2 <= 0 {
		return 0, errNonPositiveDepth
	}
	stress1 := 3 / (2 * math.Pi * z1 * z1)
	stress2 := 3 / (2 * math.Pi * z2 * z2)
	if stress2 <= 0 {
		return 0, errNonPositiveDepth
	}
	return math.Log(stress1/stress2) / math.Log(z1/z2), nil
}

func IsApprox(expected, actual, tolerance float64) bool {
	if math.Abs(expected) < 1e-15 {
		return math.Abs(actual) < tolerance
	}
	return math.Abs(actual-expected) <= tolerance*math.Abs(expected)
}

func StressInkPa(stressPa float64) float64 {
	return stressPa / 1000
}

func StressInMPa(stressPa float64) float64 {
	return stressPa / 1e6
}

func LoadInKN(loadN float64) float64 {
	return loadN / 1000
}
