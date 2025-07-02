//go:build unit

// This file tests UnicodeCheckValidate.
package checks

import (
	"github.com/UCLALibrary/validation-service/validation/config"
	"github.com/UCLALibrary/validation-service/validation/csv"
	"github.com/stretchr/testify/assert"
	"slices"
	"testing"
)

// TestUnicodeCheck checks if UnicodeCheck.Validate throws the correct errors when given text with the unicode replacement char
func TestUnicodeCheck(t *testing.T) {
	check, err := NewUnicodeCheck(config.NewProfiles())
	assert.NoError(t, err)

	// genericLocation provides a consistent location for the purposes of test comparison.
	var genericLocation = csv.Location{RowIndex: 0, ColIndex: 0}
	// slices containing expected valid and invalid URLs after table tests run
	var testInvalids = []string{"AgglomeÌ�ration de Drummondville"}

	tests := []struct {
		name     string
		profile  string
		location csv.Location
		data     [][]string
		result   bool
	}{
		{
			name:     "Valid text without replacement char",
			profile:  "DLP Staff",
			location: genericLocation,
			data:     [][]string{{"Lorem ipsum"}, {"dolor sit amet"}},
			result:   true,
		},
		{
			name:     "Invalid text with replacement char",
			profile:  "DLP Staff",
			location: genericLocation,
			data:     [][]string{{"AgglomeÌ�ration de Drummondville"}, {"dolor sit amet"}},
			result:   false,
		},
		{
			name:     "Invalid duplicate text",
			profile:  "DLP Staff",
			location: genericLocation,
			data:     [][]string{{"AgglomeÌ�ration de Drummondville"}, {"dolor sit amet"}},
			result:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := check.Validate(tt.profile, tt.location, tt.data)
			if (err != nil && tt.result) || (err == nil && !tt.result) {
				t.Errorf("Expected '%v' response was not found: %v", tt.name, err)
			}
		})
	}
	assert.True(t, slices.Equal(check.invalids, testInvalids))
}
