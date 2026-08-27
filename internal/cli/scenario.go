package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"boussinesq-z/internal/super"
)

type ScenarioFile struct {
	Forces []super.Force `json:"forces"`
	Z      float64       `json:"z,omitempty"`
	R      float64       `json:"r,omitempty"`
}

func loadScenario(path string) (ScenarioFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ScenarioFile{}, err
	}
	var file ScenarioFile
	if err := json.Unmarshal(data, &file); err != nil {
		return ScenarioFile{}, fmt.Errorf("parse %s: %w", path, err)
	}
	if err := super.ValidateList(file.Forces); err != nil {
		return ScenarioFile{}, err
	}
	return file, nil
}

func ExampleScenarioPaths() []string {
	return []string{
		"example/point-100kn.json",
		"example/superpose.json",
	}
}

func DefaultScenario() ScenarioFile {
	return ScenarioFile{
		Forces: super.ExampleList(),
	}
}

func SaveScenario(path string, file ScenarioFile) error {
	data, err := json.MarshalIndent(file, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func DescribeScenario(file ScenarioFile) string {
	return fmt.Sprintf("%d forces", len(file.Forces))
}
