package checks

import (
	"strconv"

	"github.com/UCLALibrary/validation-service/validation/config"

	"github.com/UCLALibrary/validation-service/errors"
	"github.com/UCLALibrary/validation-service/validation/csv"
)

// ItemSeqCheck type is a validator that checks for the presence of stray new lines.
//
// It implements the Validator interface and returns an error on failure to validate.
type ItemSeqCheck struct {
	profiles *config.Profiles
}

// NewItemSeqCheck checks that all values in Item Sequence are positive integers.
func NewItemSeqCheck(profiles *config.Profiles) (*ItemSeqCheck, error) {
	if profiles == nil {
		return nil, csv.NewError(errors.NilProfileErr, csv.Location{}, "nil")
	}

	return &ItemSeqCheck{
		profiles: profiles,
	}, nil
}

// Validate checks a data cell to see if Item Sequence is a positive int.
//
// This check doesn't care what profile is being used.
func (check *ItemSeqCheck) Validate(profile string, location csv.Location, csvData [][]string) error {
	if err := csv.IsValidLocation(location, csvData, profile); err != nil {
		return err
	}

	// find the header and determine if it matches Item Sequence
	header, err := csv.GetHeader(location, csvData, profile)

	if err != nil {
		return err
	}

	// Skip if we don't have an Item Sequence cell, or we're on the first (i.e., header) row
	if header != "Item Sequence" || location.RowIndex == 0 {
		return nil
	}

	value := csvData[location.RowIndex][location.ColIndex]

	objType, _ := csv.GetRowValue("Object Type", location, csvData, profile)

	// value of "Item Sequence" is allowed to be null only if the Object type is not page otherwise it must be a positive integer
	if objType == "Page" && value == "" {
		return csv.NewError(errors.PageMustBeIntErr, location, profile)
	} else if value == "" {
		return nil
	}
	// check if it is a positive int
	n, err := strconv.Atoi(value)
	if err != nil {
		return csv.NewError(errors.NotAnIntErr, location, profile)
	}

	if n <= 0 {
		return csv.NewError(errors.NotAPosIntErr, location, profile)
	}

	return nil
}
