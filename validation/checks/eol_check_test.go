//go:build unit

// Package checks consists of individual validation checks.
//
// This file checks for EOL characters in the CSV data.
package checks

import (
	"github.com/UCLALibrary/validation-service/validation/csv"
	"github.com/UCLALibrary/validation-service/validation/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestEOLCheck_Validate tests the Validate method on EOLCheck.
func TestEOLCheck_Validate(t *testing.T) {
	check, err := NewEOLCheck(util.NewProfiles())
	assert.NoError(t, err)

	// Data variations to check the EOLCheck.Validate method against
	tests := []struct {
		name     string
		location csv.Location
		data     [][]string
		result   bool
	}{
		{
			name:     "valid",
			location: csv.Location{RowIndex: 0, ColIndex: 1},
			data:     [][]string{{"Hello", "World"}, {"Hello\n", "World"}},
			result:   true,
		},
		{
			name:     "invalid",
			location: csv.Location{RowIndex: 0, ColIndex: 0},
			data:     [][]string{{"Hello\n", "World"}, {"Hello", "World"}},
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
