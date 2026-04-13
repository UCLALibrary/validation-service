// Package util provides useful resources and utilities.
package util //nolint:revive

import "github.com/UCLALibrary/validation-service/validation/csv"

// Context is the core information required for a generic validation check bundled into a single struct.
//
// It's not required for types that implement the Validator interface to use this, but it may make passing information
// to functions within the implementing type easier (i.e., require fewer arguments). It's just a convenience structure.
type Context struct {
	Profile  string
	Location csv.Location
	CsvData  [][]string
}
