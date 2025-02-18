//go:build unit

package checks

import (
	"fmt"
	csv "github.com/UCLALibrary/validation-service/csvutils"
	"github.com/UCLALibrary/validation-service/validation"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type LicenseCheck struct{}

func (check *LicenseCheck) NewLicenseLCheck(profiles *validation.Profiles) (*LicenseCheck, error) {
	if profiles == nil {
		return nil, fmt.Errorf("supplied Profiles cannot be nil")
	}

	return &LicenseCheck{}, nil
}

func (check *LicenseCheck) Validate(profile string, location csv.Location, csvData [][]string) error {
	// license not relevant to Bucketeer processing
	if profile == "bucketeer" {
		return nil
	}

	if err := csv.IsValidLocation(location, csvData); err != nil {
		return err
	}

	value := csvData[location.RowIndex][location.ColIndex]

	if err := verifyLicense(value); err != nil {
		return fmt.Errorf("License validation failed at (row: %d, column: %d) [profile: %s]: %w",
			location.RowIndex, location.ColIndex, profile, err)

	}

	return nil
}

ifunc verifyLicense(license string) error {
        r := regexp.MustCompile('^http\:\/\/[0-9a-zA-Z]([-.\w]*[0-9a-zA-Z])*(:(0-9)*)*(\/?)([a-zA-Z0-9\-\.\?\,\'\/\\+&amp;%\$#_]*)?$')
        if !r.MatchString(license) {
                return fmt.Errorf("License URL %s is not in a proper format", license)
        }

        resp, err := http.Get(license)
        if err != nil {
                return fmt.Errorf("Error connecting to license URL: %s", err.Error())
        }
        defer resp.Body.Close()
        body, err := io.ReadAll(resp.Body)
        if err != nil {
                return fmt.Errorf("Error reading body of license URL : %s", err.Error())
        }
        if len(body) == 0 {
                return fmt.Errorf("License URL %s  appears to lack content", license)
        }

        // Supplied license is valid
        return nil
}

