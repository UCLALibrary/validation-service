package checks

import (
	"maps"
	"slices"

	"github.com/UCLALibrary/validation-service/errors"
	"github.com/UCLALibrary/validation-service/validation/csv"
	"github.com/UCLALibrary/validation-service/validation/util"

	"go.uber.org/multierr"
)

// MediaMetaCheck validates the media.* fields for the Fester profile.
type MediaMetaCheck struct {
	profiles *util.Profiles
	mediaCols map[string]int
	mediaTypes []string
	mediaFields []string
	allFieldsFound bool
	allFieldsMissing bool
	someFieldsMissing bool
}

// NewMediaMetaCheck creates a new MediaMetaCheck instance, which validates the media.* fields for the Fester profile.
//
// It returns an error if the profiles argument is nil.
func NewMediaMetaCheck(profiles *util.Profiles) (*MediaMetaCheck, error) {
	if profiles == nil {
		return nil, csv.NewError(errors.NilProfileErr, csv.Location{}, "nil")
	}

	return &MediaMetaCheck{
		profiles: profiles,
		mediaCols: make(map[string]int),
		mediaTypes: []string{ "mov", "aud", "aum", "aun" },
		mediaFields: []string{ "media.width", "media.height", "media.duration", "media.format" },
		allFieldsFound: false,
		allFieldsMissing: false,
		someFieldsMissing: false,
	}, nil
}

// Validate checks if the media.* fields have been added to the CSV, and if they have been populated  for A/V media entries.
//
// If the profile is not "fester", the media metadata check is skipped.
// It checks if the header is "Type.typeOfResource" and verifies if the value is a A/V media type.
// Media types vocabulary: https://github.com/UCLALibrary/californica/blob/main/config/authorities/resource_types.yml
// Media types examined by this check: moving image (mov). sound recording (aud), sound recording-musical (aum), sound recording-nonmusical (aun)
// It returns an error if the media width/height/duration/format fields are missing or empty.
func (check *MediaMetaCheck) Validate(profile string, location csv.Location, csvData [][]string) error {

	// media metadata fields only relevant to Fester
	if profile != "fester" {
		return nil
	}

	if err := csv.IsValidLocation(location, csvData, profile); err != nil {
		return err
	}

	// find the header and determine if it matches an Type.typeOfResource header
	header, err := csv.GetHeader(location, csvData, profile)

	if err != nil {
		return err
	}

	if header != "Type.typeOfResource" {
		return nil
	}

	value := csvData[location.RowIndex][location.ColIndex]

	if slices.Contains(check.mediaTypes, value) {
		//don't scan CSV if already determined no media.* fields present
		if check.allFieldsMissing {
			return csv.NewError(errors.AllMediaErr, location, profile)
		}
		//don't scan CSV if already determined some media.* fields missing
		if check.someFieldsMissing {
			return csv.NewError(errors.SomeMediaErr, location, profile)
		}
		if !check.allFieldsFound {
			if err := check.verifyColumns(profile, location, csvData); err != nil {
				return err
			}	
		}
		if err := check.verifyContent(profile, location, csvData); err != nil {
			return err
		}
	} else {
		return nil
	}

	return nil
}

func (check *MediaMetaCheck) verifyColumns(profile string, location csv.Location, csvData [][]string) error {
	var errs error
	for colIndex, field := range csvData[0] {
		if slices.Contains(check.mediaFields, field) {
			check.mediaCols[field] = colIndex
		}
	}
	if len(check.mediaCols) < 4 {
		//compare keys to media fields, compose error for all missing fields
		if len(check.mediaCols) == 0 {
			check.allFieldsMissing = true
			return csv.NewError(errors.AllMediaErr, location, profile)
		} else {
			keys := slices.Sorted(maps.Keys(check.mediaCols))
			for _, reqField := range check.mediaFields {
				if !slices.Contains(keys, reqField) {
					switch reqField {
						case "media.width":
							errs = multierr.Combine(errs, csv.NewError(errors.WidthMissingErr, location, profile))
						case "media.height":
							errs = multierr.Combine(errs, csv.NewError(errors.HeightMissingErr, location, profile))
						case "media.duration":
							errs = multierr.Combine(errs, csv.NewError(errors.DurationMissingErr, location, profile))
						case "media.format":
							errs = multierr.Combine(errs, csv.NewError(errors.FormatMissingErr, location, profile))
					}
				}
			}
			check.someFieldsMissing = true
			return errs
		}
	}
	check.allFieldsFound = true
	return nil
}

func (check *MediaMetaCheck) verifyContent(profile string, location csv.Location, csvData [][]string) error {
	var errs error 
	//check media.* fields, compose error for all empty fields
	for fieldName, colIndex := range check.mediaCols {
		if csvData[location.RowIndex][colIndex] == "" || len(csvData[location.RowIndex][colIndex]) == 0 {
                	switch fieldName {
                		case "media.width":
                			errs = multierr.Combine(errs, csv.NewError(errors.WidthEmptyErr, location, profile))
                		case "media.height":
                			errs = multierr.Combine(errs, csv.NewError(errors.HeightEmptyErr, location, profile))
                		case "media.duration":
                			errs = multierr.Combine(errs, csv.NewError(errors.DurationEmptyErr, location, profile))
                		case "media.format":
                			errs = multierr.Combine(errs, csv.NewError(errors.FormatEmptyErr, location, profile))
                	}
		}
	}
	return errs
}
