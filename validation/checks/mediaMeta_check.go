package checks

import (
	"slices"

	"github.com/UCLALibrary/validation-service/errors"
	"github.com/UCLALibrary/validation-service/validation/csv"
	"github.com/UCLALibrary/validation-service/validation/util"
)

// LicenseCheck validates the media.* fields for the Fester profile.
type MediaMetaCheck struct {
	profiles *util.Profiles
	mediaIndexes []int
}

// NewLicenseCheck creates a new LicenseCheck instance, which validates the License field for a given profile.
//
// It returns an error if the profiles argument is nil.
func NewMediaMetaCheck(profiles *util.Profiles) (*MediaMetaCheck, error) {
	if profiles == nil {
		return nil, csv.NewError(errors.NilProfileErr, csv.Location{}, "nil")
	}

	return &MediaMetaCheck{
		profiles: profiles,
		mediaIndexes:   make([]int, 0),
	}, nil
}

// Validate checks if the License field in the CSV data is correctly formatted and points to a valid URL.
//
// If the profile is "bucketeer", the license check is skipped.
// It checks if the header is "License" and verifies if the value is a valid URL and accessible.
// It returns an error if the License field is invalid or there are issues with the URL.
func (check *MediaMetaCheck) Validate(profile string, location csv.Location, csvData [][]string) error {
	var mediaTypes = []string{ "mov", "aud", "aum", "aun" }

	// media metadata fields only relevant to Fester
	if profile != "fester" {
		return nil
	}

	if err := csv.IsValidLocation(location, csvData, profile); err != nil {
		return err
	}

	// find the header and determine if it matches an license header
	header, err := csv.GetHeader(location, csvData, profile)

	if err != nil {
		return err
	}

	if header != "Type.typeOfResource" {
		return nil
	}

	value := csvData[location.RowIndex][location.ColIndex]

	if slices.Contains(mediaTypes, value) {
		if err := check.verifyColumns(profile, location, csvData); err != nil {
			return err
		}
		for _, colIndex := range check.mediaIndexes {
			if csvData[location.RowIndex][colIndex] == "" || len(csvData[location.RowIndex][colIndex]) == 0 {
				//set error
			}
		}
	} else {
		return nil
	}

	return nil
}

func (check *MediaMetaCheck) verifyColumns(profile string, location csv.Location, csvData [][]string) error {
	var mediaFields = []string{ "media.width", "media.height", "media.duration", "media.format" }
	var columnCount = 0
	for colIndex, field := range csvData[0] {
		if slices.Contains(mediaFields, field) {
			columnCount += 1
			check.mediaIndexes = append(check.mediaIndexes, colIndex)
		}
	}
	if columnCount < 4 {
		return csv.NewError(errors.MediaFieldErr, location, profile)
	}
	return nil
}
