package checks

import (
	"regexp"

	"github.com/UCLALibrary/validation-service/validation/profiles"

	"github.com/UCLALibrary/validation-service/errors"
	"github.com/UCLALibrary/validation-service/validation/csv"
)

// FileNameCheck validates the License field for a given profile.
type FileNameCheck struct {
	profiles *profiles.Profiles
}

// NewFileNameCheck creates a new FileNameCheck instance, which checks File Name entries for whitespace.
//
// It returns an error if the `profiles` argument is nil.
func NewFileNameCheck(profiles *profiles.Profiles) (*FileNameCheck, error) {
	if profiles == nil {
		return nil, csv.NewError(errors.NilProfileErr, csv.Location{}, "nil")
	}

	return &FileNameCheck{
		profiles: profiles,
	}, nil
}

// Validate checks if the File Name field in the CSV data contains whitespace.
//
// It checks if the header is "File Name" and chwcks if the value contains whitespace.
// It returns an error if the File Name contains whitespace.
func (check *FileNameCheck) Validate(profile string, location csv.Location, csvData [][]string) error {
	if err := csv.IsValidLocation(location, csvData, profile); err != nil {
		return err
	}

	// find the header and determine if it matches an license header.
	header, err := csv.GetHeader(location, csvData, profile)

	if err != nil {
		return err
	}

	if header != "File Name" || location.RowIndex == 0 {
		return nil
	}

	value := csvData[location.RowIndex][location.ColIndex]

	whitespace := regexp.MustCompile(`\s`)
	if whitespace.MatchString(value) {
		return csv.NewError(errors.TypeWhitespaceError, location, profile)
	}

	return nil
}
