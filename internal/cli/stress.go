package cli

import (
	"fmt"

	"boussinesq-z/internal/circular"
	"boussinesq-z/internal/point"
	"boussinesq-z/internal/super"
)

func runStress(args []string) int {
	fs := flagSet("stress")
	p := fs.Float64("P", 0, "point load in N")
	z := fs.Float64("z", 0, "depth below surface in m")
	r := fs.Float64("r", 0, "radial distance in m")
	poisson := fs.Float64("poisson", 0.3, "Poisson ratio")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	load := point.Load{P: *p, Z: *z, R: *r}
	result, err := point.Evaluate(load, *poisson)
	if err != nil {
		return fail(err)
	}
	return printJSON(result)
}

func runSuperpose(args []string) int {
	fs := flagSet("superpose")
	file := fs.String("file", "example/superpose.json", "force list JSON file")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	forces, err := super.LoadList(*file)
	if err != nil {
		return fail(err)
	}
	result, err := super.Sum(forces)
	if err != nil {
		return fail(err)
	}
	return printJSON(result)
}

func runCircular(args []string) int {
	fs := flagSet("circular")
	q := fs.Float64("q", 0, "uniform pressure in Pa")
	a := fs.Float64("a", 0, "radius in m")
	z := fs.Float64("z", 0, "depth in m")
	r := fs.Float64("r", 0, "radial distance in m")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	store := circular.NewStore(defaultTablePath())
	load := circular.Load{Q: *q, A: *a, Z: *z, R: *r}
	result, err := store.Evaluate(load)
	if err != nil {
		return fail(err)
	}
	return printJSON(result)
}

func runTable(args []string) int {
	fs := flagSet("table")
	path := fs.String("path", defaultTablePath(), "influence table CSV path")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	table, err := circular.LoadCSV(*path)
	if err != nil {
		return fail(err)
	}
	fmt.Printf("rows=%d cols=%d max=%.6g\n", len(table.RVals), len(table.ZVals), table.MaxValue())
	return 0
}

func defaultTablePath() string {
	return "data/circular-influence.csv"
}
