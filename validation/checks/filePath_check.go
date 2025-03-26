package checks

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/UCLALibrary/validation-service/validation/config"
	"github.com/UCLALibrary/validation-service/validation/csv"
)

// FILE_NAME is the File Name of the current item.
const FILE_NAME = "File Name"

// Error messages
var (
	noHostDir    = "a HOST_DIR must be set"
	fileNotExist = "the file path given does not exist: %s"
)

// FilePathCheck type is a validator that checks if a File exists at the specificed location
//
// It implements the Validator interface and returns an error on failure to validate.
type FilePathCheck struct {
	profiles *config.Profiles
}

// NewFilePathCheck checks that the file exits at the given filepath in a CSV data cell.
func (check *FilePathCheck) NewFilePathCheck(profiles *config.Profiles) (*FilePathCheck, error) {
	if profiles == nil {
		return nil, csv.NewError(nilProfileErr, csv.Location{}, "nil")
	}

	return &FilePathCheck{
		profiles: profiles,
	}, nil
}

// Validate verifies the file given at that location exists
//
// This check doesn't care what profile is being used.
func (check *FilePathCheck) Validate(profile string, location csv.Location, csvData [][]string) error {
	if err := csv.IsValidLocation(location, csvData, profile); err != nil {
		return err
	}

	value := csvData[location.RowIndex][location.ColIndex]

	// obtain dir name from HOST_DIR
	hostDir := os.Getenv("HOST_DIR")
	if hostDir == "" {
		return csv.NewError(noHostDir, location, profile)
	}

	// Find the header and determine if it matches a File Name header
	header, err := csv.GetHeader(location, csvData, profile)
	if err != nil {
		return err
	}

	// Skip if we don't have a FILE_NAME header, or we're on the first (i.e., header) row
	if header != FILE_NAME || location.RowIndex == 0 {
		return nil
	}

	fullPath := filepath.Join(hostDir, value)

	// if the file doesn't exist return an error
	if _, err = os.Stat(fullPath); os.IsNotExist(err) {
		return csv.NewError(fmt.Sprintf(fileNotExist, fullPath), location, profile)
	} else {
		return nil
	}
}
