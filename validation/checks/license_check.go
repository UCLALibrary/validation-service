//go:build unit

package checks

import (
	"fmt"
	csv "github.com/UCLALibrary/validation-service/csvutils"
	"github.com/UCLALibrary/validation-service/validation"
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

	if err := csv.IsValidLicense(value); err != nil {
		return err
	}

	return nil
}
