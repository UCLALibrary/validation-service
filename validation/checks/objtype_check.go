package checks

import (
	"regexp"

        "github.com/UCLALibrary/validation-service/validation/config"
        "github.com/UCLALibrary/validation-service/validation/csv"
)

// Error messages
var (
	typeWhitespaceError  = "field contains invalid characters (e.g., spaces, line breaks)"
	typeValueError = "object type field doesn't contain valid value"
)

type ObjTypeCheck struct{
	profiles *config.Profiles
}

func (check *ObjTypeCheck) NewObjTypeCheck(profiles *config.Profiles) *ObjTypeCheck {
	return &ObjTypeCheck{
		profiles: profiles,
	}
}

func (check *ObjTypeCheck) Validate(profile string, location csv.Location, csvData [][]string) error {
	if err := csv.IsValidLocation(location, csvData, profile); err != nil {
		return err
	}

	// find the header and determine if it matches an ark header
	header, err := csv.GetHeader(location, csvData, profile)

	if err != nil {
		return err
	}

	if header != "Object Type" {
		return nil
	}

	value := csvData[location.RowIndex][location.ColIndex]

        whitespace := regexp.MustCompile(`\s`)
        if whitespace.MatchString(value) {
		return csv.NewError(typeWhitespaceError, location, profile)
        }
        valid := regexp.MustCompile(`Collection|Work|Page`)
        if !valid.MatchString(value) {
		return csv.NewError(typeValueError, location, profile)
        }

	return nil
}

