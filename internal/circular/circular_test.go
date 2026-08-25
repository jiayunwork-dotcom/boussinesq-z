package circular

import (
	"math"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadTable(t *testing.T) {
	store := NewStore("../../data/circular-influence.csv")
	table, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(table.ZVals) != 8 || len(table.RVals) != 9 {
		t.Fatalf("table size = %d x %d", len(table.ZVals), len(table.RVals))
	}
}

func TestCenterlineInfluenceMatchesAnalytic(t *testing.T) {
	store := NewStore("../../data/circular-influence.csv")
	z := 1.0
	a := 1.0
	influence, err := store.Influence(z/a, 0)
	if err != nil {
		t.Fatal(err)
	}
	analytic := 1 - math.Pow(z, 3)/math.Pow(z*z+a*a, 1.5)
	if math.Abs(influence-analytic) > 1e-3 {
		t.Fatalf("influence = %g, analytic = %g", influence, analytic)
	}
}

func TestSurfaceInsideIsOne(t *testing.T) {
	store := NewStore("../../data/circular-influence.csv")
	influence, err := store.Influence(0, 0.25)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(influence-1) > 1e-6 {
		t.Fatalf("influence = %g, want 1", influence)
	}
}

func TestSurfaceOutsideIsZero(t *testing.T) {
	store := NewStore("../../data/circular-influence.csv")
	influence, err := store.Influence(0, 3)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(influence) > 1e-6 {
		t.Fatalf("influence = %g, want 0", influence)
	}
}

func TestCircularStressScalesWithQ(t *testing.T) {
	store := NewStore("../../data/circular-influence.csv")
	first, err := store.Evaluate(Load{Q: 100000, A: 1, Z: 2, R: 0})
	if err != nil {
		t.Fatal(err)
	}
	second, err := store.Evaluate(Load{Q: 200000, A: 1, Z: 2, R: 0})
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(second.Vertical-2*first.Vertical) > 1e-6 {
		t.Fatalf("doubled q gave %g, want %g", second.Vertical, 2*first.Vertical)
	}
}

func TestTableReloadReflectsFileChange(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "influence.csv")
	table, err := GenerateTable([]float64{1, 2}, []float64{0}, 16, 8)
	if err != nil {
		t.Fatal(err)
	}
	table.Data[0][0] = 0.5
	table.Data[0][1] = 0.4
	if err := SaveCSV(path, table); err != nil {
		t.Fatal(err)
	}
	store := NewStore(path)
	first, err := store.Influence(1, 0)
	if err != nil {
		t.Fatal(err)
	}
	table.Data[0][0] = 0.7
	if err := SaveCSV(path, table); err != nil {
		t.Fatal(err)
	}
	second, err := store.Influence(1, 0)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(first-0.5) > 1e-9 || math.Abs(second-0.7) > 1e-9 {
		t.Fatalf("first=%g second=%g, expected 0.5/0.7", first, second)
	}
}

func TestGenerateTableConverges(t *testing.T) {
	err, err2 := ConvergenceError(1, 1)
	if err2 != nil {
		t.Fatal(err2)
	}
	if err > 1e-3 {
		t.Fatalf("convergence error = %g", err)
	}
}

func TestCSVRoundTrip(t *testing.T) {
	table, err := GenerateTable([]float64{0, 1}, []float64{0, 1}, 16, 8)
	if err != nil {
		t.Fatal(err)
	}
	text, err := TableCSVText(table)
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := ParseCSV(stringsReader(text))
	if err != nil {
		t.Fatal(err)
	}
	if len(parsed.Data) != 2 || len(parsed.ZVals) != 2 {
		t.Fatalf("parsed size = %d x %d", len(parsed.ZVals), len(parsed.RVals))
	}
}

func stringsReader(text string) *strings.Reader {
	return strings.NewReader(text)
}

func TestIntegrateMatchesAnalytic(t *testing.T) {
	value, err := CenterlineIntegral(1, 1)
	if err != nil {
		t.Fatal(err)
	}
	analytic := 1 - math.Pow(1, 3)/math.Pow(2, 1.5)
	if math.Abs(value-analytic) > 2e-4 {
		t.Fatalf("integral = %g, analytic = %g", value, analytic)
	}
}

func TestValidateCircularLoad(t *testing.T) {
	load := Load{Q: 100000, A: 1, Z: 2, R: 0}
	if err := load.Validate(); err != nil {
		t.Fatal(err)
	}
	load = Load{Q: 0, A: 1, Z: 2, R: 0}
	if err := load.Validate(); err == nil {
		t.Fatal("accepted zero q")
	}
	load = Load{Q: 100000, A: 0, Z: 2, R: 0}
	if err := load.Validate(); err == nil {
		t.Fatal("accepted zero radius a")
	}
}
