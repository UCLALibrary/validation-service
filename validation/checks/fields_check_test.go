//go:build unit

// This file tests fieldsCheckValidate.
package checks

import (
	"github.com/UCLALibrary/validation-service/validation/csv"
	"github.com/UCLALibrary/validation-service/validation/util"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"testing"
)

// TestReqFieldCheck_Validate tests the Validate method on EOLCheck.
func TestReqFieldCheck_Validate(t *testing.T) {
	check, err := NewReqFieldCheck(util.NewProfiles(), zaptest.NewLogger(t))
	assert.NoError(t, err)

	// Data variations to check the EOLCheck.Validate method against
	tests := []struct {
		name     string
		profile  string
		location csv.Location
		data     [][]string
		result   bool
	}{
		{
			name:     "Successful Item ARK check", // Succeeds b/c Item ARK has a value (validity checked elsewhere)
			profile:  "default",
			location: csv.Location{RowIndex: 1, ColIndex: 0},
			data: [][]string{
				{"Item ARK", "Parent ARK", "File Name", "Object Type", "Item Sequence", "Visibility", "Title",
					"media.width", "media.height", "media.duration", "media.format"},
				{"PRESENT", "PRESENT", "PRESENT", "PRESENT", "PRESENT", "PRESENT", "PRESENT", "PRESENT", "PRESENT",
					"PRESENT", "PRESENT"}},
			result: true,
		},
		{
			name:     "Successful File Name check allowing missing data", // Succeeds b/c 'Object Type' is 'Collection'
			profile:  "fester",
			location: csv.Location{RowIndex: 1, ColIndex: 2},
			data: [][]string{
				{"Item ARK", "Parent ARK", "File Name", "Object Type", "Item Sequence", "Visibility", "Title",
					"media.width", "media.height", "media.duration", "media.format"},
				{"PRESENT", "PRESENT", "", "Collection", "PRESENT", "PRESENT", "PRESENT", "PRESENT", "PRESENT",
					"PRESENT", "PRESENT"}},
			result: true,
		},
		{
			name:     "Successful File Name check allowing missing data", // Fails b/c 'Object Type' is 'Work'
			profile:  "fester",
			location: csv.Location{RowIndex: 1, ColIndex: 2},
			data: [][]string{
				{"Item ARK", "Parent ARK", "File Name", "Object Type", "Item Sequence", "Visibility", "Title",
					"media.width", "media.height", "media.duration", "media.format"},
				{"PRESENT", "PRESENT", "", "Work", "PRESENT", "PRESENT", "PRESENT", "PRESENT", "PRESENT",
					"PRESENT", "PRESENT"}},
			result: false,
		},
		{
			name:     "Missing required File Name field errors", // Error because Fester always requires 'File Name'
			profile:  "fester",
			location: csv.Location{RowIndex: 0, ColIndex: 0},
			data: [][]string{
				{"Item ARK", "Parent ARK", "Object Type", "Item Sequence", "Visibility", "Title",
					"media.width", "media.height", "media.duration", "media.format"},
				{"PRESENT", "PRESENT", "PRESENT", "PRESENT", "PRESENT", "PRESENT", "PRESENT", "PRESENT",
					"PRESENT", "PRESENT"}},
			result: false,
		},
		{
			name:     "Required missing Summary data error", // Error because 'Object Type' is 'Collection'
			profile:  "fester",
			location: csv.Location{RowIndex: 1, ColIndex: 10},
			data: [][]string{
				{"Item ARK", "Parent ARK", "Object Type", "Item Sequence", "Visibility", "Title",
					"media.width", "media.height", "media.duration", "media.format", "Summary"},
				{"PRESENT", "PRESENT", "Collection", "PRESENT", "PRESENT", "PRESENT", "PRESENT", "PRESENT",
					"PRESENT", "PRESENT", ""}},
			result: false,
		},
		{
			name:     "Successful required Summary data check", // No error because 'Object Type' is 'Work'
			profile:  "fester",
			location: csv.Location{RowIndex: 1, ColIndex: 10},
			data: [][]string{
				{"Item ARK", "Parent ARK", "Object Type", "Item Sequence", "Visibility", "Title",
					"media.width", "media.height", "media.duration", "media.format", "Summary"},
				{"PRESENT", "PRESENT", "Work", "PRESENT", "PRESENT", "PRESENT", "PRESENT", "PRESENT",
					"PRESENT", "PRESENT", ""}},
			result: true,
		},
	}

	// Iterate over test cases; fail if there isn't an error when we expect one or if there is an unexpected error
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := check.Validate(tt.profile, tt.location, tt.data)
			if (err != nil && tt.result) || (err == nil && !tt.result) {
				t.Errorf("Expected '%v' response was not found: %v", tt.name, err)
			}
		})
	}
}
