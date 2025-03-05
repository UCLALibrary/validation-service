package checks

import (
	"github.com/UCLALibrary/validation-service/validation/csv"
	"regexp"
	"strings"

	"github.com/UCLALibrary/validation-service/validation/config"
	"go.uber.org/multierr"
)

// ITEM_ARK is the ARK of the current item.
const ITEM_ARK = "Item ARK"

// PARENT_ARK is the ARK of the parent item.
const PARENT_ARK = "Parent ARK"

// Error messages
var (
	profileErr      = "supplied profile cannot be nil"
	noPrefixErr     = "ARK must start with 'ark:/'"
	naanTooShortErr = "NAAN must be at least 5 digits long"
	naanProfileErr  = "The supplied NAAN is not allowed for the supplied profile"
	noObjIdErr      = "The ARK must contain an object identifier"
	invalidObjIdErr = "The object identifier and qualifier is not valid"
	arkValFailed    = "ARK validation failed"
)

// The naanProfiles mapping gives us a way to lookup valid NAANs for a profile.
var naanProfiles = map[string]map[string]struct{}{
	"default":   {"21198": struct{}{}},
	"test":      {"21198": struct{}{}},
	"bucketeer": {"21198": struct{}{}},
	"fester":    {"21198": struct{}{}},
}

// ARKCheck type is a validator that checks for a valid ARK.
//
// It implements the Validator interface and returns an error on failure to validate.
type ARKCheck struct{}

// NewARKCheck checks that an ARK is valid.
func (check *ARKCheck) NewARKCheck(profiles *config.Profiles) (*ARKCheck, error) {
	if profiles == nil {
		return nil, csv.NewError(profileErr, csv.Location{}, "nil")
	}

	return &ARKCheck{}, nil
}

// Validate checks that a data cell has a valid ARK in it.
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
	if header != ITEM_ARK && header != PARENT_ARK || location.RowIndex == 0 {
		return nil
	}

	value := csvData[location.RowIndex][location.ColIndex]

	// Check if the CSV data cell has a valid ARK
	if err := verifyARK(value, location, profile); err != nil {
		return csv.NewError(arkValFailed, location, profile, err)
	}

	return nil
}

// verifyARK validates if the given string is a valid ARK.
func verifyARK(ark string, location csv.Location, profile string) error {
	var errs error

	// Ensure the ARK starts with "ark:/"
	if !strings.HasPrefix(ark, "ark:/") {
		errs = multierr.Combine(errs, csv.NewError(noPrefixErr, location, profile))
		return errs // Early return since the rest of validation depends on this
	}

	// Remove "ark:/" for further validation
	arkBody := strings.TrimPrefix(ark, "ark:/")

	// Validate the NAAN separately
	naanRegex := regexp.MustCompile(`^(\d+)`)
	naanMatch := naanRegex.FindStringSubmatch(arkBody)
	if naanMatch == nil || len(naanMatch[1]) < 5 {
		errs = multierr.Combine(errs, csv.NewError(naanTooShortErr, location, profile))
	}

	// Extract NAAN and ObjectIdentifier for further validation
	naan := naanMatch[1]
	objectID := strings.TrimPrefix(arkBody, naan)
	objectID = strings.TrimPrefix(objectID, "/")

	// Validate that the NAAN is allowed for the supplied profile
	if _, exists := naanProfiles[profile][naan]; !exists {
		errs = multierr.Combine(errs, csv.NewError(naanProfileErr, location, profile))
	}

	if objectID == "" {
		errs = multierr.Combine(errs, csv.NewError(noObjIdErr, location, profile))
		return errs
	}

	// Validate the remaining ARK structure (ObjectIdentifier + Qualifier)
	arkRegex := regexp.MustCompile(`^([\w\-./]+)(\?.*)?$`)
	if !arkRegex.MatchString(objectID) {
		errs = multierr.Combine(errs, csv.NewError(invalidObjIdErr, location, profile))
	}

	return errs
}
