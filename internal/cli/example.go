package cli

import (
	"boussinesq-z/internal/point"
)

func runExample(args []string) int {
	fs := flagSet("example")
	file := fs.String("file", "example/point-100kn.json", "scenario JSON file")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	load, err := loadExample(*file)
	if err != nil {
		return fail(err)
	}
	result, err := point.Evaluate(load, 0.3)
	if err != nil {
		return fail(err)
	}
	return printJSON(result)
}

func loadExample(path string) (point.Load, error) {
	var load point.Load
	err := readJSONFile(path, &load)
	return load, err
}
