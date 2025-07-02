package checks

// This checks if the value in Object Type is valid.
import (
	"regexp"

	"github.com/UCLALibrary/validation-service/validation/config"

	"github.com/UCLALibrary/validation-service/errors"
	"github.com/UCLALibrary/validation-service/validation/csv"
)

// ObjTypeCheck validates the "Object Type" field in the provided CSV data.
//
// It checks whether the field contains a valid value (either "Collection", "Work", or "Page") and ensures there are no
// whitespace characters in the value.
type ObjTypeCheck struct {
	profiles *config.Profiles
}

// NewObjTypeCheck creates a new instance of ObjTypeCheck to validate the "Object Type" field for the provided profiles.
//
// It returns an error if the profiles argument is nil.
func NewObjTypeCheck(profiles *config.Profiles) (*ObjTypeCheck, error) {
	if profiles == nil {
		return nil, csv.NewError(errors.NilProfileErr, csv.Location{}, "nil")
	}

	return &ObjTypeCheck{
		profiles: profiles,
	}, nil
}

// Validate checks if the "Object Type" field in the CSV data contains a valid value and is free of whitespace.
//
// It ensures that the header matches "Object Type" and validates the value in each data cell. The valid values for "Object Type"
// are "Collection", "Work", and "Page". If the value contains whitespace or is not one of these valid values, an error is returned.
func (check *ObjTypeCheck) Validate(profile string, location csv.Location, csvData [][]string) error {
	if err := csv.IsValidLocation(location, csvData, profile); err != nil {
		return err
	}

	// find the header and determine if it matches Object Type
	header, err := csv.GetHeader(location, csvData, profile)

	if err != nil {
		return err
	}

	// Skip if we don't have an object tpe cell, or we're on the first (i.e., header) row
	if header != "Object Type" || location.RowIndex == 0 {
		return nil
	}

	value := csvData[location.RowIndex][location.ColIndex]

	whitespace := regexp.MustCompile(`\s`)
	if whitespace.MatchString(value) {
		return csv.NewError(errors.TypeWhitespaceError, location, profile)
	}
	valid := regexp.MustCompile(`Collection|Work|Page`)
	if !valid.MatchString(value) {
		return csv.NewError(errors.TypeValueError, location, profile)
	}

	return nil
}
