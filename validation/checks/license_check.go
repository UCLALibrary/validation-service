// This file checks if the provided License is correct.
package checks

import (
	"io"
	"net/http"
	"regexp"
        "slices"


	"github.com/UCLALibrary/validation-service/errors"
	"github.com/UCLALibrary/validation-service/validation/csv"
	"github.com/UCLALibrary/validation-service/validation/util"
)

var valids []string
var invalids []string

// LicenseCheck validates the License field for a given profile.
type LicenseCheck struct {
	profiles *util.Profiles
}

// NewLicenseCheck creates a new LicenseCheck instance, which validates the License field for a given profile.
//
// It returns an error if the profiles argument is nil.
func NewLicenseCheck(profiles *util.Profiles) (*LicenseCheck, error) {
	if profiles == nil {
		return nil, csv.NewError(errors.NilProfileErr, csv.Location{}, "nil")
	}

        valids = make([]string, 1)
        invalids = make([]string, 1)

	return &LicenseCheck{
		profiles: profiles,
	}, nil
}

// Validate checks if the License field in the CSV data is correctly formatted and points to a valid URL.
//
// If the profile is "bucketeer", the license check is skipped.
// It checks if the header is "License" and verifies if the value is a valid URL and accessible.
// It returns an error if the License field is invalid or there are issues with the URL.
func (check *LicenseCheck) Validate(profile string, location csv.Location, csvData [][]string) error {
	// license not relevant to Bucketeer processing
	if profile == "bucketeer" {
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

	if header != "License" {
		return nil
	}

	value := csvData[location.RowIndex][location.ColIndex]

        if slices.Contains(valids, value) {
                return nil
        } else if slices.Contains(invalids, value) {
                return csv.NewError(errors.UrlDupeBadErr, location, profile)
        }

	if err := check.verifyLicense(value, profile, location); err != nil {
                invalids = append(invalids, value)
		return err
	}

        valids = append(valids, value)
	return nil
}

// verifyLicense checks if the given license string is a valid URL and if it can be accessed.
//
// It uses a regular expression to validate the URL format and sends an HTTP GET request to ensure the URL is reachable.
// It returns an error if the URL is not formatted correctly or if the URL is not accessible.
func (check *LicenseCheck) verifyLicense(license string, profile string, location csv.Location) error {
	r := regexp.MustCompile(`^^http\:\/\/[0-9a-zA-Z]([-.\w]*[0-9a-zA-Z])*(:(0-9)*)*(\/?)([a-zA-Z0-9\-\.\?\,\'\/\\\+&amp;%\$#_]*)?$`)
	if !r.MatchString(license) {
		return csv.NewError(errors.UrlFormatErr, location, profile)
	}

	resp, err := http.Get(license)
	if err != nil {
		return csv.NewError(errors.UrlConnectErr, location, profile)
	}
	//noinspection GoUnhandledErrorResult
	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return csv.NewError(errors.UrlReadErr, location, profile)
	}

	// Supplied license is valid
	return nil
}
