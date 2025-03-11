package checks

import (
	"io"
	"net/http"
	"regexp"

        "github.com/UCLALibrary/validation-service/validation/config"
        "github.com/UCLALibrary/validation-service/validation/csv"
)

// Error messages
var (
	noProfileErr = "supplied profile cannot be nil"
	urlFormatErr = "license URL is not in a proper format (check for HTTPS)"
	urlConnectErr = "problem connecting to license URL"
	urlReadErr = "problem reading body of license URL"
	emptyBodyErr = "licence has no content"
)

type LicenseCheck struct{
	profiles *config.Profiles
}

func (check *LicenseCheck) NewLicenseLCheck(profiles *config.Profiles) (*LicenseCheck, error) {
	if profiles == nil {
		return nil, csv.NewError(noProfileErr, csv.Location{}, "nil")
	}

	return &LicenseCheck{
		profiles: profiles,
	}, nil
}

func (check *LicenseCheck) Validate(profile string, location csv.Location, csvData [][]string) error {
	// license not relevant to Bucketeer processing
	if profile == "bucketeer" {
		return nil
	}

	if err := csv.IsValidLocation(location, csvData, profile); err != nil {
		return err
	}

	// find the header and determine if it matches an ark header
	header, err := csv.GetHeader(location, csvData, profile)

	if err != nil {
		return err
	}

	if header != "License" {
		return nil
	}

	value := csvData[location.RowIndex][location.ColIndex]

	if err := verifyLicense(value, profile, location); err != nil {
		return err

	}

	return nil
}

func verifyLicense(license string, profile string, location csv.Location) error {
        r := regexp.MustCompile("^http\\:\\/\\/[0-9a-zA-Z]([-.\\w]*[0-9a-zA-Z])*(:(0-9)*)*(\\/?)([a-zA-Z0-9\\-\\.\\?\\,\\'\\/\\\\\\+&amp;%\\$#_]*)?$")
        if !r.MatchString(license) {
                return csv.NewError(urlFormatErr, location, profile)
        }

        resp, err := http.Get(license)
        if err != nil {
                return csv.NewError(urlConnectErr, location, profile)
        }
        defer resp.Body.Close()
        body, err := io.ReadAll(resp.Body)
        if err != nil {
                return csv.NewError(urlReadErr, location, profile)
        }
        if len(body) == 0 {
                return csv.NewError(emptyBodyErr, location, profile)
        }

        // Supplied license is valid
        return nil
}

