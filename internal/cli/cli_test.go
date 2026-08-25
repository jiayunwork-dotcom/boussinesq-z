package cli

import (
	"testing"

	"boussinesq-z/internal/point"
)

func TestParseFloat(t *testing.T) {
	value, err := parseFloat("100000", "P")
	if err != nil {
		t.Fatal(err)
	}
	if value != 100000 {
		t.Fatalf("value = %g", value)
	}
}

func TestValidateLoad(t *testing.T) {
	if err := validateLoad(100000, 2, 0); err != nil {
		t.Fatal(err)
	}
	if err := validateLoad(100000, 0, 0); err == nil {
		t.Fatal("accepted singularity")
	}
}

func TestLoadExample(t *testing.T) {
	load, err := loadExample("../../example/point-100kn.json")
	if err != nil {
		t.Fatal(err)
	}
	if load.P != 100000 || load.Z != 2 || load.R != 0 {
		t.Fatalf("load = %+v", load)
	}
}

func TestExampleScenarioPaths(t *testing.T) {
	if len(ExampleScenarioPaths()) != 2 {
		t.Fatal("expected two example paths")
	}
}

func TestDefaultTablePath(t *testing.T) {
	if defaultTablePath() != "data/circular-influence.csv" {
		t.Fatalf("path = %q", defaultTablePath())
	}
}

func TestEvaluateExample(t *testing.T) {
	load, err := loadExample("../../example/point-100kn.json")
	if err != nil {
		t.Fatal(err)
	}
	result, err := point.Evaluate(load, 0.3)
	if err != nil {
		t.Fatal(err)
	}
	if result.Vertical <= 0 {
		t.Fatalf("result = %+v", result)
	}
}

func TestLoadScenario(t *testing.T) {
	file, err := loadScenario("../../example/superpose.json")
	if err != nil {
		t.Fatal(err)
	}
	if len(file.Forces) != 2 {
		t.Fatalf("forces = %+v", file.Forces)
	}
}
