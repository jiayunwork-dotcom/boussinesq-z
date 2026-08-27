package main

import (
	"flag"
	"fmt"
	"os"

	"boussinesq-z/internal/circular"
)

func main() {
	out := flag.String("out", "data/circular-influence.csv", "output CSV path")
	radial := flag.Int("radial", 96, "radial integration steps")
	angle := flag.Int("angle", 48, "angular integration steps")
	flag.Parse()
	table, err := circular.GenerateTable(
		circular.DefaultZAxis(),
		circular.DefaultRAxis(),
		*radial,
		*angle,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "generate:", err)
		os.Exit(1)
	}
	if err := circular.SaveCSV(*out, table); err != nil {
		fmt.Fprintln(os.Stderr, "save:", err)
		os.Exit(1)
	}
	fmt.Println("wrote", *out, "rows", len(table.Data), "cols", len(table.ZVals))
}
