package cli

import (
	"fmt"
	"os"
)

func Run(args []string) int {
	if len(args) == 0 {
		return runServe([]string{})
	}
	switch args[0] {
	case "stress":
		return runStress(args[1:])
	case "superpose":
		return runSuperpose(args[1:])
	case "circular":
		return runCircular(args[1:])
	case "table":
		return runTable(args[1:])
	case "example":
		return runExample(args[1:])
	case "serve":
		return runServe(args[1:])
	case "help", "-h", "--help":
		printHelp()
		return 0
	case "version":
		fmt.Println("boussinesq-z 1.0.0")
		return 0
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n", args[0])
		printHelp()
		return 2
	}
}

func printHelp() {
	fmt.Println(`boussinesq-z: elastic half-space Boussinesq vertical stress calculator

Usage:
  boussinesq-z                         start HTTP server on :8080
  boussinesq-z stress -P 100000 -z 2 -r 0
  boussinesq-z superpose -file example/superpose.json
  boussinesq-z circular -q 100000 -a 1 -z 2 -r 0
  boussinesq-z table -path data/circular-influence.csv
  boussinesq-z example -file example/point-100kn.json
  boussinesq-z serve -addr :8080

HTTP:
  POST /api/stress      {"P":100000,"z":2,"r":0}
  POST /api/superpose   {"forces":[{"P":100000,"z":2,"r":0}]}
  POST /api/circular    {"q":100000,"a":1,"z":2,"r":0}`)
}

func fail(err error) int {
	fmt.Fprintln(os.Stderr, "error:", err)
	return 1
}
