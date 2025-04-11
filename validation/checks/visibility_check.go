package checks

import (
	"regexp"

	"github.com/UCLALibrary/validation-service/errors"
	"github.com/UCLALibrary/validation-service/validation/csv"
	"github.com/UCLALibrary/validation-service/validation/util"
)

type VisibilityCheck struct{}

func NewvisibilityCheck(profiles *util.Profiles) (*VisibilityCheck, error) {
	if profiles == nil {
		return nil, csv.NewError(errors.NilProfileErr, csv.Location{}, "nil")
	}

	return &VisibilityCheck{}, nil
}

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
