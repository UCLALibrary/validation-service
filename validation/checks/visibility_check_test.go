//go:build unit

// This tests visibilityCheckValidate
package checks

import (
	"github.com/UCLALibrary/validation-service/validation/csv"
	"github.com/UCLALibrary/validation-service/validation/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestObjTypeCheck_Validate tests the Validate method on ObjTypeCheck.
func TestVisibilityCheck_Validate(t *testing.T) {
	check, err := NewVisibilityCheck(util.NewProfiles())
	assert.NoError(t, err)

	// Data variations to check the VisibilityCheck.Validate method against
	tests := []struct {
		name     string
		location csv.Location
		data     [][]string
		result   bool
	}{
		{
			name:     "valid",
			location: csv.Location{RowIndex: 1, ColIndex: 0},
			data:     [][]string{{"Visibility"}, {"ucla"}},
			result:   true,
		},
		{
			name:     "invalid Bad Value",
			location: csv.Location{RowIndex: 1, ColIndex: 0},
			data:     [][]string{{"Visibility"}, {"Other"}},
			result:   false,
		},
		{
			name:     "invalid Bad Chars",
			location: csv.Location{RowIndex: 1, ColIndex: 0},
			data:     [][]string{{"Visibility"}, {"private "}},
			result:   false,
		},
	}

	// Iterate over test cases; fail if there isn't an error when we expect one or if there is an unexpected error
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = check.Validate("default", tt.location, tt.data)
			if (err != nil && tt.result) || (err == nil && !tt.result) {
				t.Errorf("Expected '%v' response was not found: %v", tt.name, err)
			}
		})
	}
}
