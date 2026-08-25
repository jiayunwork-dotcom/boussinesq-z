package point

import (
	"math"
)

type HorizontalResult struct {
	Radial     float64 `json:"sigma_r"`
	Tangential float64 `json:"sigma_theta"`
	ShearRZ    float64 `json:"tau_rz"`
	Poisson    float64 `json:"poisson"`
}

func Horizontal(l Load, poisson float64) (HorizontalResult, error) {
	if err := l.Validate(); err != nil {
		return HorizontalResult{}, err
	}
	if err := ValidatePoisson(poisson); err != nil {
		return HorizontalResult{}, err
	}
	radius := l.Radius()
	return HorizontalResult{
		Radial:     RadialStress(l.P, l.Z, l.R, radius, poisson),
		Tangential: TangentialStress(l.P, l.Z, l.R, radius, poisson),
		ShearRZ:    evaluatedShear(l),
		Poisson:    poisson,
	}, nil
}

func shearRadialInput(l Load) float64 {
	return l.Radius()
}

func evaluatedShear(l Load) float64 {
	radius := l.Radius()
	return ShearStress(l.P, l.Z, shearRadialInput(l), radius)
}

func VerticalSumFromComponents(sigmaZ, sigmaR, sigmaTheta float64) float64 {
	return sigmaZ + sigmaR + sigmaTheta
}

func DeviatoricVertical(sigmaZ, meanNormal float64) float64 {
	return sigmaZ - meanNormal
}

func MeanNormal(sigmaZ, sigmaR, sigmaTheta float64) float64 {
	return (sigmaZ + sigmaR + sigmaTheta) / 3
}

func OctahedralShearFromComponents(sigmaZ, sigmaR, sigmaTheta float64) float64 {
	sum := (sigmaZ-sigmaR)*(sigmaZ-sigmaR) +
		(sigmaR-sigmaTheta)*(sigmaR-sigmaTheta) +
		(sigmaTheta-sigmaZ)*(sigmaTheta-sigmaZ)
	return math.Sqrt(sum / 3)
}

func MaxShearFromComponents(sigmaZ, sigmaR, sigmaTheta float64) float64 {
	max := sigmaZ
	min := sigmaZ
	for _, value := range []float64{sigmaR, sigmaTheta} {
		if value > max {
			max = value
		}
		if value < min {
			min = value
		}
	}
	return (max - min) / 2
}

func RadialRatio(result Result) float64 {
	if result.Vertical == 0 {
		return 0
	}
	return result.Radial / result.Vertical
}

func TangentialRatio(result Result) float64 {
	if result.Vertical == 0 {
		return 0
	}
	return result.Tangential / result.Vertical
}

func ShearRatio(result Result) float64 {
	if result.Vertical == 0 {
		return 0
	}
	return result.ShearRZ / result.Vertical
}

func CenterlineHorizontal(l Load, poisson float64) (HorizontalResult, error) {
	l.R = 0
	return Horizontal(l, poisson)
}

func IsHydrostaticStress(result Result) bool {
	return math.Abs(result.Radial-result.Tangential) < 1e-9 &&
		math.Abs(result.Radial-result.Vertical) < 1e-9
}

func BulkStress(result Result) float64 {
	return (result.Vertical + result.Radial + result.Tangential) / 3
}

func StressDeviation(result Result) float64 {
	mean := BulkStress(result)
	return math.Sqrt(
		(result.Vertical-mean)*(result.Vertical-mean) +
			(result.Radial-mean)*(result.Radial-mean) +
			(result.Tangential-mean)*(result.Tangential-mean),
	)
}
