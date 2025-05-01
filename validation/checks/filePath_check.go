package checks

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/UCLALibrary/validation-service/errors"
	"github.com/UCLALibrary/validation-service/validation/csv"
	"github.com/UCLALibrary/validation-service/validation/util"
)

// FILE_NAME is the File Name of the current item.
const FILE_NAME = "File Name"

// FilePathCheck type is a validator that checks if a File exists at the specified location.
//
// It implements the Validator interface and returns an error on failure to validate.
type FilePathCheck struct {
	profiles *util.Profiles
}

// NewFilePathCheck returns a new FilePathCheck, which validates that the file path specified in a CSV data cell points to an existing file.
//
// It returns an error if the provided profiles argument is nil.
func NewFilePathCheck(profiles *util.Profiles) (*FilePathCheck, error) {
	if profiles == nil {
		return nil, csv.NewError(errors.NilProfileErr, csv.Location{}, "nil")
	}

	return &FilePathCheck{
		profiles: profiles,
	}, nil
}

// Validate verifies the file given at that location exists.
//
// This check doesn't care what profile is being used.
func (check *FilePathCheck) Validate(profile string, location csv.Location, csvData [][]string) error {
	if err := csv.IsValidLocation(location, csvData, profile); err != nil {
		return err
	}

	value := check.stripPrefix(csvData[location.RowIndex][location.ColIndex])

	// obtain dir name from HOST_DIR
	hostDir := os.Getenv("HOST_DIR")
	if hostDir == "" {
		return csv.NewError(errors.NoHostDir, location, profile)
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
		return csv.NewError(fmt.Sprintf(errors.FileNotExist, fullPath), location, profile)
	} else {
		return nil
	}
}

// stripPrefix strips the prefix found in the CSV file paths so that various sub-dirs will match when compared.
func (check *FilePathCheck) stripPrefix(filePath string) string {
	prefix := "Masters/" // The masters directory mount point, which is found in HOST_DIR and the CSV file paths

	// Strip the prefix if found in the CSV resource file's file path
	if strings.HasPrefix(filePath, prefix) {
		return filePath[len(prefix):]
	}

	// Else, return as is
	return filePath
}
