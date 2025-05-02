package validation

import (
	"github.com/UCLALibrary/validation-service/validation/csv"
)

// Validator interface defines how implementations should be called.
type Validator interface {
	Validate(profile string, location csv.Location, csvData [][]string) error
}
