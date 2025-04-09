package checks

import (
	"io"
	"net/http"
	"regexp"

	"github.com/UCLALibrary/validation-service/errors"
	"github.com/UCLALibrary/validation-service/validation/csv"
	"github.com/UCLALibrary/validation-service/validation/util"
)

type LicenseCheck struct {
	profiles *util.Profiles
}

func NewLicenseCheck(profiles *util.Profiles) (*LicenseCheck, error) {
	if profiles == nil {
		return nil, csv.NewError(errors.NilProfileErr, csv.Location{}, "nil")
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

	// find the header and determine if it matches an license header
	header, err := csv.GetHeader(location, csvData, profile)

	if err != nil {
		return err
	}

	if header != "License" {
		return nil
	}

	value := csvData[location.RowIndex][location.ColIndex]

	if err := check.verifyLicense(value, profile, location); err != nil {
		return err
	}

	return nil
}

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
