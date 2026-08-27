package point

import (
	"math"
	"testing"
)

func TestCenterlineExampleStress(t *testing.T) {
	load := Load{P: 100000, Z: 2, R: 0}
	result, err := Evaluate(load, 0.3)
	if err != nil {
		t.Fatal(err)
	}
	want := 3 * 100000 / (2 * math.Pi * 4)
	if math.Abs(result.Vertical-want) > 1e-6 {
		t.Fatalf("sigma_z = %g, want %g", result.Vertical, want)
	}
	if math.Abs(result.Influence-3/(2*math.Pi*4)) > 1e-9 {
		t.Fatalf("influence = %g", result.Influence)
	}
	if result.Vertical < 10000 || result.Vertical > 14000 {
		t.Fatalf("sigma_z = %g, expected tens of kPa", result.Vertical)
	}
}

func TestForceDoublingDoublesStress(t *testing.T) {
	base := Load{P: 100000, Z: 2, R: 1}
	doubled := Load{P: 200000, Z: 2, R: 1}
	first, err := Evaluate(base, 0.3)
	if err != nil {
		t.Fatal(err)
	}
	second, err := Evaluate(doubled, 0.3)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(second.Vertical-2*first.Vertical) > 1e-6 {
		t.Fatalf("doubling P gave %g, want %g", second.Vertical, 2*first.Vertical)
	}
}

func TestDepthDoublingReducesCenterlineByFour(t *testing.T) {
	shallow := Load{P: 100000, Z: 2, R: 0}
	deep := Load{P: 100000, Z: 4, R: 0}
	first, err := Evaluate(shallow, 0.3)
	if err != nil {
		t.Fatal(err)
	}
	second, err := Evaluate(deep, 0.3)
	if err != nil {
		t.Fatal(err)
	}
	ratio := first.Vertical / second.Vertical
	if math.Abs(ratio-4) > 1e-6 {
		t.Fatalf("ratio = %g, want 4", ratio)
	}
}

func TestFarFieldTendsToZero(t *testing.T) {
	load := Load{P: 100000, Z: 1000, R: 1000}
	result, err := Evaluate(load, 0.3)
	if err != nil {
		t.Fatal(err)
	}
	if result.Vertical > 0.01 {
		t.Fatalf("far field sigma_z = %g, want near zero", result.Vertical)
	}
}

func TestRadialDistanceLowersStress(t *testing.T) {
	center := Load{P: 100000, Z: 2, R: 0}
	offset := Load{P: 100000, Z: 2, R: 1}
	centerResult, err := Evaluate(center, 0.3)
	if err != nil {
		t.Fatal(err)
	}
	offsetResult, err := Evaluate(offset, 0.3)
	if err != nil {
		t.Fatal(err)
	}
	if offsetResult.Vertical >= centerResult.Vertical {
		t.Fatalf("offset stress %g should be below center stress %g", offsetResult.Vertical, centerResult.Vertical)
	}
}

func TestZeroLoadGivesZeroStress(t *testing.T) {
	load := Load{P: 0, Z: 2, R: 0}
	result, err := Evaluate(load, 0.3)
	if err != nil {
		t.Fatal(err)
	}
	if result.Vertical != 0 || result.Influence != 0 {
		t.Fatalf("result = %+v, want zero stress", result)
	}
}

func TestValidationRejectsBadInputs(t *testing.T) {
	load := Load{P: -1, Z: 2, R: 0}
	if err := load.Validate(); err == nil {
		t.Fatal("accepted negative P")
	}
	load = Load{P: 100, Z: -1, R: 0}
	if err := load.Validate(); err == nil {
		t.Fatal("accepted negative z")
	}
	load = Load{P: 100, Z: 0, R: 0}
	if err := load.Validate(); err == nil {
		t.Fatal("accepted singularity")
	}
	load = Load{P: 100, Z: 2, R: -1}
	if err := load.Validate(); err == nil {
		t.Fatal("accepted negative r")
	}
}

func TestHorizontalComponentsAtCenterline(t *testing.T) {
	load := Load{P: 100000, Z: 2, R: 0}
	result, err := Evaluate(load, 0.3)
	if err != nil {
		t.Fatal(err)
	}
	if math.IsNaN(result.Radial) || math.IsNaN(result.Tangential) || math.IsNaN(result.ShearRZ) {
		t.Fatalf("horizontal components are NaN: %+v", result)
	}
	if math.Abs(result.ShearRZ) > 1e-9 {
		t.Fatalf("centerline shear = %g, want 0", result.ShearRZ)
	}
}

func TestInfluenceDecayExponent(t *testing.T) {
	exponent, err := DecayExponentAtCenterline(2, 4)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(exponent-(-2)) > 1e-9 {
		t.Fatalf("decay exponent = %g, want -2", exponent)
	}
}

func TestConversions(t *testing.T) {
	if NewtonsToKilonewtons(100000) != 100 {
		t.Fatal("kN conversion wrong")
	}
	if PascalsToKilopascals(11936.6) != 11.9366 {
		t.Fatalf("kPa conversion = %g", PascalsToKilopascals(11936.6))
	}
	if MetersToMillimeters(2) != 2000 {
		t.Fatal("mm conversion wrong")
	}
}

func TestValidatePoisson(t *testing.T) {
	if err := ValidatePoisson(0.3); err != nil {
		t.Fatal(err)
	}
	if err := ValidatePoisson(0.6); err == nil {
		t.Fatal("accepted poisson > 0.5")
	}
}
