// Package checks provides individual validators used by the validation service.
package checks

import (
	"regexp"
	"strings"

	"github.com/UCLALibrary/validation-service/validation/config"

	"github.com/UCLALibrary/validation-service/errors"
	"github.com/UCLALibrary/validation-service/validation/csv"
	"go.uber.org/multierr"
)

// ItemARK is the ARK of the current item.
const ItemARK = "Item ARK"

// ParentARK is the ARK of the parent item.
const ParentARK = "Parent ARK"

// The naanProfiles mapping gives us a way to look up valid NAANs for a profile.
var naanProfiles = map[string]map[string]struct{}{
	"DLP Staff": {
		"21198": {},
		"13030": {},
	},
	"Test": {
		"21198": {},
		"13030": {},
	},
	"Bucketeer": {
		"21198": {},
		"13030": {},
	},
	"Fester": {
		"21198": {},
		"13030": {},
	},
}

// ARKCheck type is a validator that checks for a valid ARK.
//
// It implements the Validator interface and returns an error on failure to validate.
type ARKCheck struct {
	profiles *config.Profiles
}

// NewARKCheck returns a new ARKCheck, which validates that an ARK identifier is properly formatted.
//
// It returns an error if the provided profiles argument is nil.
func NewARKCheck(profiles *config.Profiles) (*ARKCheck, error) {
	if profiles == nil {
		return nil, csv.NewError(errors.NilProfileErr, csv.Location{}, "nil")
	}

	return &ARKCheck{
		profiles: profiles,
	}, nil
}

// Validate returns an error if a data cell does not have a valid ARK in it.
func (check *ARKCheck) Validate(profile string, location csv.Location, csvData [][]string) error {
	if err := csv.IsValidLocation(location, csvData, profile); err != nil {
		return err
	}

	// Find the header and determine if it matches an ARK header
	header, err := csv.GetHeader(location, csvData, profile)
	if err != nil {
		return err
	}

	// Skip if we don't have an ARK cell, or we're on the first (i.e., header) row
	if header != ItemARK && header != ParentARK || location.RowIndex == 0 {
		return nil
	}

	value := strings.TrimSpace(csvData[location.RowIndex][location.ColIndex])
	if value == "" {
		return nil // We let the ReqFieldCheck check for presence and validate ARKs here
	}

	// Check if the CSV data cell has a valid ARK
	if err := check.verifyARK(value, location, profile); err != nil {
		return csv.NewError(errors.ArkValFailed, location, profile, err)
	}

	return nil
}

// verifyARK validates if the given string is a valid ARK.
func (check *ARKCheck) verifyARK(ark string, location csv.Location, profile string) error {
	var errs error

	// Ensure the ARK starts with "ark:/"
	if !strings.HasPrefix(ark, "ark:/") {
		errs = multierr.Combine(errs, csv.NewError(errors.NoPrefixErr, location, profile))
		return errs // Early return since the rest of validation depends on this
	}

	// Remove "ark:/" for further validation
	arkBody := strings.TrimPrefix(ark, "ark:/")

	// Validate the NAAN separately
	naanRegex := regexp.MustCompile(`^(\d+)`)
	naanMatch := naanRegex.FindStringSubmatch(arkBody)
	if naanMatch == nil || len(naanMatch[1]) < 5 {
		errs = multierr.Combine(errs, csv.NewError(errors.NaanTooShortErr, location, profile))
	}

	// Extract NAAN and ObjectIdentifier for further validation
	naan := naanMatch[1]
	objectID := strings.TrimPrefix(arkBody, naan)
	objectID = strings.TrimPrefix(objectID, "/")

	// Validate that the NAAN is allowed for the supplied profile
	if _, exists := naanProfiles[profile][naan]; !exists {
		errs = multierr.Combine(errs, csv.NewError(errors.NaanProfileErr, location, profile))
	}

	if objectID == "" {
		errs = multierr.Combine(errs, csv.NewError(errors.NoObjIDErr, location, profile))
		return errs
	}

	// Validate the remaining ARK structure (ObjectIdentifier + Qualifier)
	arkRegex := regexp.MustCompile(`^([\w\-./]+)(\?.*)?$`)
	if !arkRegex.MatchString(objectID) {
		errs = multierr.Combine(errs, csv.NewError(errors.InvalidObjIDErr, location, profile))
	}

	return errs
}
