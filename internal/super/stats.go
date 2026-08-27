package super

import (
	"fmt"
	"math"
)

type Stats struct {
	Count        int     `json:"count"`
	TotalLoad    float64 `json:"total_load_n"`
	TotalStress  float64 `json:"total_stress_pa"`
	MeanStress   float64 `json:"mean_stress_pa"`
	MaxStress    float64 `json:"max_stress_pa"`
	MinStress    float64 `json:"min_stress_pa"`
	StddevStress float64 `json:"stddev_stress_pa"`
}

func ComputeStats(forces []Force) (Stats, error) {
	result, err := Sum(forces)
	if err != nil {
		return Stats{}, err
	}
	values := make([]float64, 0, len(result.Components))
	for _, component := range result.Components {
		values = append(values, component.Vertical)
	}
	if len(values) == 0 {
		return Stats{}, fmt.Errorf("no forces")
	}
	max, min := values[0], values[0]
	sum := 0.0
	for _, value := range values {
		sum += value
		if value > max {
			max = value
		}
		if value < min {
			min = value
		}
	}
	mean := sum / float64(len(values))
	variance := 0.0
	for _, value := range values {
		delta := value - mean
		variance += delta * delta
	}
	variance /= float64(len(values))
	return Stats{
		Count:        len(values),
		TotalLoad:    TotalLoad(forces),
		TotalStress:  result.SigmaZ,
		MeanStress:   mean,
		MaxStress:    max,
		MinStress:    min,
		StddevStress: math.Sqrt(variance),
	}, nil
}

func FormatStats(stats Stats) string {
	return fmt.Sprintf(
		"count=%d load=%.4g N stress=%.4g Pa mean=%.4g max=%.4g min=%.4g",
		stats.Count,
		stats.TotalLoad,
		stats.TotalStress,
		stats.MeanStress,
		stats.MaxStress,
		stats.MinStress,
	)
}

func NearlyLinear(first, second float64, tolerance float64) bool {
	if math.Abs(first) < 1e-15 {
		return math.Abs(second) < tolerance
	}
	return math.Abs(second/first-2) <= tolerance
}

func RatioText(first, second float64) string {
	if math.Abs(first) < 1e-15 {
		return "undefined"
	}
	return fmt.Sprintf("%.4g", second/first)
}

func RankComponents(forces []Force) ([]float64, error) {
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

func MedianStress(forces []Force) (float64, error) {
	values, err := RankComponents(forces)
	if err != nil {
		return 0, err
	}
	n := len(values)
	if n%2 == 1 {
		return values[n/2], nil
	}
	return (values[n/2-1] + values[n/2]) / 2, nil
}

func StressSpread(forces []Force) (float64, error) {
	max, err := MaxComponent(forces)
	if err != nil {
		return 0, err
	}
	min, err := MinComponent(forces)
	if err != nil {
		return 0, err
	}
	return max - min, nil
}
