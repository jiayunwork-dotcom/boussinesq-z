package super

import (
	"math"
	"testing"
)

func TestSuperpositionIsAdditive(t *testing.T) {
	force := Force{P: 100000, Z: 2, R: 0}
	single, err := Sum([]Force{force})
	if err != nil {
		t.Fatal(err)
	}
	double, err := Sum([]Force{force, force})
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(double.SigmaZ-2*single.SigmaZ) > 1e-6 {
		t.Fatalf("superposed %g, want %g", double.SigmaZ, 2*single.SigmaZ)
	}
}

func TestSuperposeExampleFile(t *testing.T) {
	forces, err := LoadList("../../example/superpose.json")
	if err != nil {
		t.Fatal(err)
	}
	result, err := Sum(forces)
	if err != nil {
		t.Fatal(err)
	}
	if result.ForceCount != 2 || result.SigmaZ <= 0 {
		t.Fatalf("result = %+v", result)
	}
}

func TestZeroForceInList(t *testing.T) {
	forces := []Force{{P: 0, Z: 2, R: 0}, {P: 100000, Z: 2, R: 0}}
	result, err := Sum(forces)
	if err != nil {
		t.Fatal(err)
	}
	if result.ForceCount != 2 || result.SigmaZ <= 0 {
		t.Fatalf("result = %+v", result)
	}
}

func TestCombineAtPoint(t *testing.T) {
	forces := []Force{{P: 100000, Z: 2, R: 0}, {P: 50000, Z: 3, R: 1}}
	result, err := CombineAtPoint(forces, 2, 0.5)
	if err != nil {
		t.Fatal(err)
	}
	if result.SigmaZ <= 0 || len(result.PerForce) != 2 {
		t.Fatalf("result = %+v", result)
	}
}

func TestStats(t *testing.T) {
	forces := []Force{{P: 100000, Z: 2, R: 0}, {P: 200000, Z: 3, R: 1}}
	stats, err := ComputeStats(forces)
	if err != nil {
		t.Fatal(err)
	}
	if stats.Count != 2 || stats.TotalLoad != 300000 {
		t.Fatalf("stats = %+v", stats)
	}
	if stats.MaxStress < stats.MinStress {
		t.Fatal("max below min")
	}
}

func TestValidateList(t *testing.T) {
	if err := ValidateList([]Force{{P: -1, Z: 2, R: 0}}); err == nil {
		t.Fatal("accepted negative force")
	}
	if err := ValidateList(nil); err == nil {
		t.Fatal("accepted empty list")
	}
}

func TestScaleAndMove(t *testing.T) {
	forces := []Force{{P: 100000, Z: 2, R: 0}}
	scaled := ScaleForces(forces, 2)
	moved := MoveForces(forces, 1, 0.5)
	if scaled[0].P != 200000 {
		t.Fatalf("scaled = %+v", scaled)
	}
	if moved[0].Z != 3 || moved[0].R != 0.5 {
		t.Fatalf("moved = %+v", moved)
	}
}

func TestScenarioExamples(t *testing.T) {
	for _, scenario := range ScenarioExamples() {
		if err := ValidateScenario(scenario); err != nil {
			t.Fatalf("%s: %v", scenario.Name, err)
		}
		result, err := RunScenario(scenario)
		if err != nil {
			t.Fatalf("%s: %v", scenario.Name, err)
		}
		if result.SigmaZ <= 0 {
			t.Fatalf("%s: non-positive stress", scenario.Name)
		}
	}
}

func TestParseJSON(t *testing.T) {
	forces, err := ParseJSON([]byte(`{"forces":[{"P":100000,"z":2,"r":0}]}`))
	if err != nil {
		t.Fatal(err)
	}
	if len(forces) != 1 {
		t.Fatalf("forces = %+v", forces)
	}
}

func TestNearlyLinear(t *testing.T) {
	if !NearlyLinear(100, 200, 1e-9) {
		t.Fatal("doubling should be nearly linear")
	}
}
