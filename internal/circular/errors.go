package circular

import "errors"

var (
	errInvalidLoad  = errors.New("invalid circular load parameters")
	errInvalidSteps = errors.New("integration steps must be positive")
)

func IsTableError(err error) bool {
	return errors.Is(err, errInvalidLoad) || errors.Is(err, errInvalidSteps)
}
