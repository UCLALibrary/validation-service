// Package checks provides individual validators used by the validation service.
//
// This file checks that data cells do not have end of line (EOL) characters.
package checks

import (
	"strings"

	"github.com/UCLALibrary/validation-service/errors"
	"github.com/UCLALibrary/validation-service/validation/csv"
	"github.com/UCLALibrary/validation-service/validation/util"
)

// EOLCheck type is a validator that checks for the presence of stray new lines.
//
// It implements the Validator interface and returns an error on failure to validate.
type EOLCheck struct {
	profiles *util.Profiles
}

// NewEOLCheck checks that there are no EOLs in a CSV data cell.
func NewEOLCheck(profiles *util.Profiles) (*EOLCheck, error) {
	if profiles == nil {
		return nil, csv.NewError(errors.NilProfileErr, csv.Location{}, "nil")
	}

	return &EOLCheck{
		profiles: profiles,
	}, nil
}

// Validate checks a data cell has a new line character in it.
//
// This check doesn't care what profile is being used.
func (check *EOLCheck) Validate(profile string, location csv.Location, csvData [][]string) error {
	if err := csv.IsValidLocation(location, csvData, profile); err != nil {
		return err
	}

	value := csvData[location.RowIndex][location.ColIndex]

	// Check if the CSV data cell under review has any unexpected EOLs in it
	if strings.Contains(value, "\n") || strings.Contains(value, "\r") {
		return csv.NewError(errors.EolFoundErr, location, profile)
	}

	return nil
}
