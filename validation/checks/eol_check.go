//go:build unit

package checks

import (
	"fmt"
	csv "github.com/UCLALibrary/validation-service/csvutils"
	"github.com/UCLALibrary/validation-service/validation"
	"strings"
)

// EOLCheck type is a validator that checks for the presence of stray new lines.
//
// It implements the Validator interface and returns an error on failure to validate.
type EOLCheck struct{}

// NewEOLCheck checks that there are no EOLs in a CSV data cell.
func (check *EOLCheck) NewEOLCheck(profiles *validation.Profiles) (*EOLCheck, error) {
	if profiles == nil {
		return nil, fmt.Errorf("supplied Profiles cannot be nil")
	}

	return &EOLCheck{}, nil
}

// Validate checks a data cell has a new line character in it.
//
// This check doesn't care what profile is being used.
func (check *EOLCheck) Validate(profile string, location csv.Location, csvData [][]string) error {
	if err := csv.IsValidLocation(location, csvData); err != nil {
		return err
	}

	value := csvData[location.RowIndex][location.ColIndex]

	// Check if the CSV data cell under review has any unexpected EOLs in it
	if strings.Contains(value, "\n") || strings.Contains(value, "\r") {
		return fmt.Errorf("character for EOL found in cell at (row: %d, column: %d)[profile: %s]",
			location.RowIndex, location.ColIndex, profile)
	}

	return nil
}
