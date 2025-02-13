package checks

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	csv "github.com/UCLALibrary/validation-service/csvutils"
	"github.com/UCLALibrary/validation-service/validation/config"
	"go.uber.org/multierr"
)

var ITEMARK = "Item ARK"
var PARENTARK = "Parent ARK"

// Error messages
var (
	profileErr        = errors.New("supplied Profiles cannot be nil")
	noPrefixErr       = errors.New("ARK must start with 'ark:/'")
	naanTooShortErr   = errors.New("NAAN must be at least 5 digits long")
	defaultProErr     = errors.New("For the 'default' profile, the NAAN must be '21198'")
	noObjIdentErr     = errors.New("The ARK must contain an object identifier")
	invalidObjIdenErr = errors.New("The object identifier and qualifier is not valid")
)

// ARKCheck type is a validator that checks for a valid Ark
//
// It implements the Validator interface and returns an error on failure to validate.
type ARKCheck struct{}

// NewARKCheck checks that an Ark is valid
func (check *ARKCheck) NewARKCheck(profiles *config.Profiles) (*ARKCheck, error) {
	if profiles == nil {
		return nil, profileErr
	}

	return &ARKCheck{}, nil
}

// Validate checks a data cell has a valid Ark in it.
//
// If the supplied “profile” is “default” the institutional prefix in the ARK must be 21198
func (check *ARKCheck) Validate(profile string, location csv.Location, csvData [][]string) error {
	if err := csv.IsValidLocation(location, csvData); err != nil {
		return err
	}

	// find the header and determine if it matches an ark header
	header, err := csv.GetHeader(location, csvData)

	if err != nil {
		return err
	}

	if header != ITEMARK && header != PARENTARK {
		return nil
	}

	value := csvData[location.RowIndex][location.ColIndex]

	// Check if the CSV data cell has a valid Ark
	if err := verifyArk(value, profile); err != nil {
		return fmt.Errorf("ARK validation failed at (row: %d, column: %d) [profile: %s]: %w",
			location.RowIndex, location.ColIndex, profile, err)
	}

	return nil
}

// VerifyARK validates if the given string is a valid ARK
func verifyArk(ark string, profile string) error {
	var errs error

	// Ensure the ARK starts with "ark:/"
	if !strings.HasPrefix(ark, "ark:/") {
		errs = multierr.Combine(errs, noPrefixErr)
		return errs // Early return since the rest of validation depends on this
	}

	// Remove "ark:/" for further validation
	arkBody := strings.TrimPrefix(ark, "ark:/")

	// Validate the NAAN separately
	naanRegex := regexp.MustCompile(`^(\d+)`)
	naanMatch := naanRegex.FindStringSubmatch(arkBody)
	if naanMatch == nil || len(naanMatch[1]) < 5 {
		errs = multierr.Combine(errs, naanTooShortErr)
	}
	// Extract NAAN and ObjectIdentifier for further validation
	naan := naanMatch[1]
	objectID := strings.TrimPrefix(arkBody, naan)
	objectID = strings.TrimPrefix(objectID, "/")
	// Additional validation if the profile is "default"
	if profile == "default" && naan != "21198" {
		errs = multierr.Combine(errs, defaultProErr)
	}

	if objectID == "" {
		errs = multierr.Combine(errs, noObjIdentErr)
		return errs
	}
	// Validate the remaining ARK structure (ObjectIdentifier + Qualifier)
	arkRegex := regexp.MustCompile(`^([\w\-./]+)(\?.*)?$`)
	if !arkRegex.MatchString(objectID) {
		errs = multierr.Combine(errs, invalidObjIdenErr)
	}

	return errs
}
