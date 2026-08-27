package cli

import (
	"fmt"
	"strconv"

	"boussinesq-z/internal/point"
)

func parseFloat(text, label string) (float64, error) {
	value, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return 0, fmt.Errorf("%s must be a number, got %q", label, text)
	}
	return value, nil
}

func validatePositive(value float64, label string) error {
	if value <= 0 {
		return fmt.Errorf("%s must be positive, got %g", label, value)
	}
	return nil
}

func validateLoad(P, z, r float64) error {
	return point.ValidateLoad(P, z, r)
}

func requireFinite(value float64, label string) error {
	return point.CheckFinite(value, label)
}

func normalizeNumber(text string) string {
	if len(text) > 0 && text[0] == '"' {
		if parsed, err := strconv.Unquote(text); err == nil {
			return parsed
		}
	}
	return text
}

func parseP(text string) (float64, error) {
	return strconv.ParseFloat(text, 64)
}
