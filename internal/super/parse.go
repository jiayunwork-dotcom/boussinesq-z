package super

import (
	"encoding/json"
	"fmt"
	"os"
)

type ListFile struct {
	Forces []Force `json:"forces"`
}

func LoadList(path string) ([]Force, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var file ListFile
	if err := json.Unmarshal(data, &file); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	if err := ValidateList(file.Forces); err != nil {
		return nil, err
	}
	return file.Forces, nil
}

func SaveList(path string, forces []Force) error {
	data, err := json.MarshalIndent(ListFile{Forces: forces}, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func ParseJSON(data []byte) ([]Force, error) {
	var file ListFile
	if err := json.Unmarshal(data, &file); err != nil {
		return nil, err
	}
	if err := ValidateList(file.Forces); err != nil {
		return nil, err
	}
	return file.Forces, nil
}

func FromPoints(values []PointInput) ([]Force, error) {
	forces := make([]Force, 0, len(values))
	for _, value := range values {
		force := Force{P: value.P, Z: value.Z, R: value.R}
		if err := force.Validate(); err != nil {
			return nil, fmt.Errorf("force %s: %w", force, err)
		}
		forces = append(forces, force)
	}
	return forces, nil
}

type PointInput struct {
	P float64 `json:"P"`
	Z float64 `json:"z"`
	R float64 `json:"r"`
}

func MarshalForces(forces []Force) ([]byte, error) {
	return json.MarshalIndent(ListFile{Forces: forces}, "", "  ")
}

func ExampleList() []Force {
	return []Force{
		{P: 100000, Z: 2, R: 0},
		{P: 50000, Z: 3, R: 1},
	}
}

func ExamplePath() string {
	return "example/superpose.json"
}

func DefaultListFile() string {
	return ExamplePath()
}
