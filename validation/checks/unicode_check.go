// This file checks if the provided License is correct.
package checks

import (
	"slices"
	"strings"

	"github.com/UCLALibrary/validation-service/errors"
	"github.com/UCLALibrary/validation-service/validation/csv"
	"github.com/UCLALibrary/validation-service/validation/util"
)

// UnicodeCheck checks for unicode replacement character (U+FFFD) in CSVs
type UnicodeCheck struct {
	profiles *util.Profiles
	invalids []string
}

// NewUnicodeCheck creates a new UnicodeCheck instance, which checks for U+FFFD char in entries
//
// It returns an error if the profiles argument is nil.
func NewUnicodeCheck(profiles *util.Profiles) (*UnicodeCheck, error) {
	if profiles == nil {
		return nil, csv.NewError(errors.NilProfileErr, csv.Location{}, "nil")
	}

	return &UnicodeCheck{
		profiles: profiles,
		invalids: make([]string, 0),
	}, nil
}

// Validate checks if the CSV data contains the U+FFFD char
//
// It returns an error if the data contains the U+FFFD char
func (check *UnicodeCheck) Validate(profile string, location csv.Location, csvData [][]string) error {

	if err := csv.IsValidLocation(location, csvData, profile); err != nil {
		return err
	}

	value := csvData[location.RowIndex][location.ColIndex]

	if slices.Contains(check.invalids, value) {
		return csv.NewError(errors.DupeUnicodeErr, location, profile)
	}

	if strings.ContainsRune(value, 0xFFFD) {
		check.invalids = append(check.invalids, value)
		return csv.NewError(errors.UnicodeErr, location, profile)
	}

	return nil
}
