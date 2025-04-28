package checks

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/UCLALibrary/validation-service/errors"
	"github.com/UCLALibrary/validation-service/validation/csv"
	"github.com/UCLALibrary/validation-service/validation/util"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

var profileFields = map[string]map[string]struct {
	dataReq     bool     // Whether data is required in the cells (could just be a header only check)
	objTypes    []string // Data cell must be present for these 'Object Types'
	notObjTypes []string // Data cell must be present for all but these 'Object Types'
}{
	"default": {
		"Item ARK":      {true, []string{}, []string{}},
		"Parent ARK":    {true, []string{}, []string{"Collection"}},
		"File Name":     {false, []string{}, []string{"Collection"}},
		"Object Type":   {true, []string{}, []string{}},
		"Item Sequence": {true, []string{"Page"}, []string{}},
		"Visibility":    {true, []string{}, []string{}},
		"Title":         {true, []string{}, []string{}},
		"Summary":       {true, []string{"Collection"}, []string{}},
	},
	"test": {
		"Item ARK":   {true, []string{}, []string{}},
		"Visibility": {true, []string{}, []string{}},
	},
	"bucketeer": {
		"Item ARK":   {true, []string{}, []string{}},
		"File Name":  {true, []string{}, []string{}},
		"Visibility": {true, []string{}, []string{}},
	},
	"fester": {
		"Item ARK":      {true, []string{}, []string{}},
		"Parent ARK":    {true, []string{}, []string{"Collection"}},
		"File Name":     {false, []string{}, []string{"Collection"}},
		"Object Type":   {true, []string{}, []string{}},
		"Item Sequence": {true, []string{"Page"}, []string{}},
		"Visibility":    {true, []string{}, []string{}},
		"Title":         {true, []string{}, []string{}},
		"Summary":       {true, []string{"Collection"}, []string{}},
	},
}

// condition encapsulates conditional information about an 'Object Type' (ot) check.
type condition struct {
	otValues []string
	match    bool
}

// ReqFieldCheck is a validator that checks that all the required fields for a profile are present.
//
// It implements the Validator interface and returns an error on failure to validate.
type ReqFieldCheck struct {
	profiles *util.Profiles
	logger   *zap.Logger
}

// NewReqFieldCheck returns a new ReqFieldCheck, which validates that all required fields are present for a given profile.
//
// It returns an error if the provided profiles argument is nil. The logger is used to record validation details.
func NewReqFieldCheck(profiles *util.Profiles, logger *zap.Logger) (*ReqFieldCheck, error) {
	if profiles == nil {
		return nil, csv.NewError(errors.NilProfileErr, csv.Location{}, "nil")
	}

	return &ReqFieldCheck{
		profiles: profiles,
		logger:   logger,
	}, nil
}

// Validate checks the headers row to confirm that all the required fields for a profile are present.
func (check *ReqFieldCheck) Validate(profile string, location csv.Location, csvData [][]string) error {
	var multiErr error

	// Check that our location is valid, given the supplied CSV data
	if err := csv.IsValidLocation(location, csvData, profile); err != nil {
		return err // We return this right away, because something is broken
	}

	// Check headers for columns where we just care about the header, not its data value
	if err := check.checkHeaders(profile, location, csvData); err != nil {
		multiErr = multierr.Combine(multiErr, err)
	}

	// Get the header for the data cell we're checking
	header, err := csv.GetHeader(location, csvData, profile)
	if err != nil {
		errMsg := fmt.Sprintf(errors.BadHeaderErr, fmt.Sprintf("[index: %s]", strconv.Itoa(location.ColIndex)))
		return csv.NewError(errMsg, location, profile, err) // We return this right away, because something is broken
	}

	// Consolidate our contextual information into a single structure
	context := util.Context{
		Profile:  profile,
		Location: location,
		CsvData:  csvData,
	}

	// Check headers where we care about the presence of the header and its cell data
	if profileCfg, exists := profileFields[profile]; exists {
		if field, exists := profileCfg[header]; exists {
			if field.dataReq && len(field.objTypes) == 0 && len(field.notObjTypes) == 0 {
				err = check.confirmExistence(context, header)
				check.logger.Debug("confirmExistence", zap.String("Header", header),
					zap.Bool("Data required", field.dataReq), zap.Error(err))
				multiErr = multierr.Combine(multiErr, err)
			} else if len(field.notObjTypes) == 0 && len(field.objTypes) > 0 {
				requirements := condition{field.objTypes, true}
				err = check.confirmWithOT(context, header, requirements)
				check.logger.Debug("confirmWithOT", zap.String("Header", header),
					zap.Bool("Data required with `Object Type` checks", field.dataReq),
					zap.Strings("`Object Type` requirements", field.objTypes), zap.Error(err))
				multiErr = multierr.Combine(multiErr, err)
			} else if len(field.objTypes) == 0 && len(field.notObjTypes) > 0 {
				requirements := condition{field.notObjTypes, false}
				err = check.confirmWithOT(context, header, requirements)
				check.logger.Debug("confirmWithOT", zap.String("Header", header),
					zap.Bool("Data required with `Object Type` exclusions", field.dataReq),
					zap.Strings("`Object Type` exclusions", field.notObjTypes), zap.Error(err))
				multiErr = multierr.Combine(multiErr, err)
			} else if len(field.objTypes) > 0 && len(field.notObjTypes) > 0 {
				err = csv.NewError(fmt.Sprintf(errors.ProfileConfigErr, profile), location, profile)
				check.logger.Error(fmt.Sprintf("Bad profile configuration: %s", profile), zap.Error(err))
				multiErr = multierr.Combine(multiErr, err)
			}
		}
	} else {
		return csv.NewError(fmt.Sprintf(errors.UnknownProfileErr, context.Profile), context.Location, context.Profile)
	}

	// If we found any errors, report them
	if len(multierr.Errors(multiErr)) > 0 {
		return multiErr
	}

	return nil
}

// checkHeaders checks that all the headers that are required (but don't have a data requirement) are found.
func (check *ReqFieldCheck) checkHeaders(profile string, location csv.Location, csvData [][]string) error {
	var multiErr error

	// We check the required headers only while we are processing the first row
	if location.RowIndex == 0 {
		// We only check for required fields with no data requirements once, at the RowIndex==0, ColIndex==0 position
		if location.ColIndex == 0 {
			// Check for required fields that don't have data requirements
			if profileCfg, exists := profileFields[profile]; exists {
				for fieldName, value := range profileCfg {
					// If the fieldName we check is required but doesn't have a data requirement look in the csvData
					if !value.dataReq {
						row := csvData[0]
						found := false

						// Check the csvData for the fieldName we're checking
						for colIndex := 0; colIndex < len(row); colIndex++ {
							// If we find it in our CSV data, it's okay (i.e., was required and was found)
							if fieldName == row[colIndex] {
								found = true
							}
						}

						// If we looked through all the CSV data's headers, and it's not there, that's a problem
						if !found {
							newErr := csv.NewError(fmt.Sprintf(errors.FieldNotFoundErr, fieldName), location, profile)
							multiErr = multierr.Combine(multiErr, newErr)
						}

						check.logger.Debug("Required field check",
							zap.Bool(fmt.Sprintf("`%s` found", fieldName), found))
					}
				}
			} else {
				return csv.NewError(fmt.Sprintf(errors.UnknownProfileErr, profile), location, profile)
			}
		} // Else: once we've checked the headers once, we don't need to keep checking them; we just drop through
	}

	// If we found any errors, report them
	if len(multierr.Errors(multiErr)) > 0 {
		return multiErr
	}

	return nil
}

// confirmWithOT confirms data exists only if 'Object Type' matches a particular value.
//
// Parameters:
// - context: Contextual information about where in the parsing the validation engine is
// - header: The header of the field being validated
// - requirements: The conditional checks that need to be used to confirm whether a value is found
//
// Returns:
// - error: If the validation check fails
func (check *ReqFieldCheck) confirmWithOT(context util.Context, header string, requirements condition) error {
	// Look up the value of our row's (i.e., item's) "Object Type" column/field
	rowValue, err := csv.GetRowValue("Object Type", context.Location, context.CsvData, context.Profile)
	if err != nil {
		return csv.NewError(fmt.Sprintf(errors.BadHeaderErr, header), context.Location, context.Profile, err)
	}

	// If our 'Object Type' value isn't one of the ones we care about, we don't need to check the data cell
	if check.finds(requirements.otValues, rowValue) != requirements.match {
		return nil
	}

	// Confirm our data value exists
	return check.confirmExistence(context, header)
}

// confirmExistence just confirms that a cell for a supplied header has something in it.
func (check *ReqFieldCheck) confirmExistence(context util.Context, header string) error {
	value := strings.TrimSpace(context.CsvData[context.Location.RowIndex][context.Location.ColIndex])
	if value == "" {
		return csv.NewError(fmt.Sprintf(errors.FieldDataNotFoundErr, header), context.Location, context.Profile)
	}

	return nil
}

// finds checks whether a supplied string exists in a supplied slice.
func (check *ReqFieldCheck) finds(slice []string, value string) bool {
	for _, sliceValue := range slice {
		if sliceValue == value {
			return true
		}
	}

	return false
}
