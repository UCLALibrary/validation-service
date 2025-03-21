// Package validation provides tools to validate CSV data.
//
// This file provides a Validator interface that individual checks should implement.
package validation

import (
	"github.com/UCLALibrary/validation-service/validation/csv"
)

// Validator interface defines how implementations should be called.
type Validator interface {
	Validate(profile string, location csv.Location, csvData [][]string) error
}
