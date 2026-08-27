package super

import (
	"fmt"
	"math"
)

type Scenario struct {
	Name   string  `json:"name"`
	Forces []Force `json:"forces"`
	FieldZ float64 `json:"field_z"`
	FieldR float64 `json:"field_r"`
}

func RunScenario(scenario Scenario) (CombinedAtPoint, error) {
	if scenario.Name == "" {
		return CombinedAtPoint{}, fmt.Errorf("scenario name is empty")
	}
	result, err := CombineAtPoint(scenario.Forces, scenario.FieldZ, scenario.FieldR)
	if err != nil {
		return CombinedAtPoint{}, err
	}
	return result, nil
}

func ScenarioExamples() []Scenario {
	return []Scenario{
		{Name: "single-center", Forces: []Force{{P: 100000, Z: 2, R: 0}}, FieldZ: 2, FieldR: 0},
		{Name: "two-off-axis", Forces: []Force{{P: 100000, Z: 2, R: 0}, {P: 50000, Z: 3, R: 1}}, FieldZ: 2, FieldR: 0.5},
		{Name: "deep-far", Forces: []Force{{P: 100000, Z: 10, R: 5}}, FieldZ: 10, FieldR: 5},
	}
}

func ScaleScenario(scenario Scenario, factor float64) Scenario {
	scenario.Forces = ScaleForces(scenario.Forces, factor)
	return scenario
}

func DoubleScenario(scenario Scenario) Scenario {
	return ScaleScenario(scenario, 2)
}

func SamePoint(a, b CombinedAtPoint) bool {
	return math.Abs(a.Z-b.Z) < 1e-9 && math.Abs(a.R-b.R) < 1e-9
}

func TotalScenarioLoad(scenario Scenario) float64 {
	return TotalLoad(scenario.Forces)
}

func ValidateScenario(scenario Scenario) error {
	if scenario.Name == "" {
		return fmt.Errorf("scenario name is empty")
	}
	if scenario.FieldZ < 0 || scenario.FieldR < 0 {
		return fmt.Errorf("field point must be non-negative")
	}
	if scenario.FieldZ == 0 && scenario.FieldR == 0 {
		return fmt.Errorf("field point cannot be the singularity")
	}
	return ValidateList(scenario.Forces)
}

func DescribeScenario(scenario Scenario) string {
	return fmt.Sprintf(
		"%s: %d forces, field z=%.4g r=%.4g",
		scenario.Name,
		len(scenario.Forces),
		scenario.FieldZ,
		scenario.FieldR,
	)
}

func FieldDistance(z, r float64) float64 {
	return math.Sqrt(z*z + r*r)
}

func FieldAngle(z, r float64) float64 {
	if z <= 0 {
		return math.Pi / 2
	}
	return math.Atan(r / z)
}
