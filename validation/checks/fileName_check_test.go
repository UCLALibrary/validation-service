//go:build unit

// This tests objtypeCheckValidate.
package checks

import (
	"github.com/UCLALibrary/validation-service/validation/config"
	"github.com/UCLALibrary/validation-service/validation/csv"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestObjTypeCheck_Validate tests the Validate method on ObjTypeCheck.
func TestFileNameCheck(t *testing.T) {
	check, err := NewFileNameCheck(config.NewProfiles())
	assert.NoError(t, err)

	// genericLocation provides a consistent location for the purposes of test comparison.
	var genericLocation = csv.Location{RowIndex: 1, ColIndex: 0}

	// Data variations to check the ObjType.Validate method against.
	tests := []struct {
		name     string
		location csv.Location
		data     [][]string
		result   bool
	}{
		{
			name:     "valid file name",
			location: genericLocation,
			data:     [][]string{{"File Name"}, {"ninan/image/21198-zz00256728_1659676_master.tif"}},
			result:   true,
		},
		{
			name:     "invalid file name",
			location: genericLocation,
			data:     [][]string{{"File Name"}, {"brown/masters/21198-zz0019qjrz_734145 master.tif"}},
			result:   false,
		},
	}

	// Iterate over test cases; fail if there isn't an error when we expect one or if there is an unexpected error.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = check.Validate("bucketeer", tt.location, tt.data)
			if (err != nil && tt.result) || (err == nil && !tt.result) {
				t.Errorf("Expected '%v' response was not found: %v", tt.name, err)
			}
		})
	}
}
