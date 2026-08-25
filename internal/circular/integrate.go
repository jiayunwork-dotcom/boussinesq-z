package circular

import (
	"math"
)

func IntegrateInfluence(z, a, r float64, radialSteps, angleSteps int) (float64, error) {
	if z < 0 || a <= 0 || r < 0 {
		return 0, errInvalidLoad
	}
	if radialSteps < 1 || angleSteps < 1 {
		return 0, errInvalidSteps
	}
	if z == 0 {
		switch {
		case r < a:
			return 1, nil
		case r == a:
			return 0.5, nil
		default:
			return 0, nil
		}
	}
	dr := a / float64(radialSteps)
	dTheta := 2 * math.Pi / float64(angleSteps)
	sum := 0.0
	for i := 0; i < radialSteps; i++ {
		rho := (float64(i) + 0.5) * dr
		for j := 0; j < angleSteps; j++ {
			theta := (float64(j) + 0.5) * dTheta
			rho2 := rho * rho
			r2 := r * r
			cross := 2 * r * rho * math.Cos(theta)
			radius2 := z*z + rho2 + r2 - cross
			if radius2 <= 1e-15 {
				continue
			}
			radius := math.Sqrt(radius2)
			integrand := 3 * z * z * z / (2 * math.Pi) * rho / math.Pow(radius, 5)
			sum += integrand
		}
	}
	return sum * dr * dTheta, nil
}

func GenerateTable(zVals, rVals []float64, radialSteps, angleSteps int) (Table, error) {
	if len(zVals) == 0 || len(rVals) == 0 {
		return Table{}, errInvalidLoad
	}
	data := make([][]float64, len(rVals))
	for i, r := range rVals {
		row := make([]float64, len(zVals))
		for j, z := range zVals {
			value, err := IntegrateInfluence(z, 1, r, radialSteps, angleSteps)
			if err != nil {
				return Table{}, err
			}
			row[j] = value
		}
		data[i] = row
	}
	return Table{ZVals: zVals, RVals: rVals, Data: data}, nil
}

func DefaultZAxis() []float64 {
	return []float64{0, 0.1, 0.2, 0.5, 1, 2, 5, 10}
}

func DefaultRAxis() []float64 {
	return []float64{0, 0.25, 0.5, 0.75, 1, 1.5, 2, 3, 5}
}

func CenterlineIntegral(z, a float64) (float64, error) {
	return IntegrateInfluence(z, a, 0, 96, 24)
}

func EdgeIntegral(z, a float64) (float64, error) {
	return IntegrateInfluence(z, a, a, 96, 48)
}

func CompareToAnalytic(z, a float64) (float64, error) {
	numeric, err := CenterlineIntegral(z, a)
	if err != nil {
		return 0, err
	}
	if z == 0 {
		return math.Abs(numeric - 1), nil
	}
	analytic := 1 - math.Pow(z, 3)/math.Pow(z*z+a*a, 1.5)
	return math.Abs(numeric - analytic), nil
}

func ConvergenceError(z, a float64) (float64, error) {
	coarse, err := IntegrateInfluence(z, a, 0, 32, 12)
	if err != nil {
		return 0, err
	}
	fine, err := IntegrateInfluence(z, a, 0, 128, 48)
	if err != nil {
		return 0, err
	}
	return math.Abs(coarse - fine), nil
}
