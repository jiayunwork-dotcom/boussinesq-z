package main

import (
	"os"

	"boussinesq-z/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:]))
}
