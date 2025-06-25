package checks

import (
	"regexp"

	"github.com/UCLALibrary/validation-service/validation/profiles"

	"github.com/UCLALibrary/validation-service/errors"
	"github.com/UCLALibrary/validation-service/validation/csv"
)

// VisibilityCheck validates that a value in the "Visibility" field is valid.
//
// A valid "Visibility" value must be one of: "open", "ucla", or "private", and must not contain any whitespace.
type VisibilityCheck struct{}

// NewVisibilityCheck creates a new instance of VisibilityCheck.
//
// Returns an error if the provided profiles argument is nil.
func NewVisibilityCheck(profiles *profiles.Profiles) (*VisibilityCheck, error) {
	if profiles == nil {
		return nil, csv.NewError(errors.NilProfileErr, csv.Location{}, "nil")
	}

	return &VisibilityCheck{}, nil
}

// Validate verifies that a cell in the "Visibility" column contains a valid value.
//
// The value must not include whitespace and must match one of the following: "open", "ucla", or "private".
// The function will skip validation if the header is not "Visibility" or if the row is the header row (row index 0).
func (check *VisibilityCheck) Validate(profile string, location csv.Location, csvData [][]string) error {
	if err := csv.IsValidLocation(location, csvData, profile); err != nil {
		return err
	}

	// find the header and determine if it matches Object Type
	header, err := csv.GetHeader(location, csvData, profile)

	if err != nil {
		return err
	}

	// Skip if we don't have an object tpe cell, or we're on the first (i.e., header) row
	if header != "Visibility" || location.RowIndex == 0 {
		return nil
	}

	value := csvData[location.RowIndex][location.ColIndex]

	whitespace := regexp.MustCompile(`\s`)
	if whitespace.MatchString(value) {
		return csv.NewError(errors.TypeWhitespaceError, location, profile)
	}
	valid := regexp.MustCompile(`open|ucla|private`)
	if !valid.MatchString(value) {
		return csv.NewError(errors.VisibilityValueError, location, profile)
	}

	return nil
}
