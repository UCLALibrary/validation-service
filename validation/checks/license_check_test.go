//go:build unit

package checks

import (
	"github.com/UCLALibrary/validation-service/validation/config"
	"github.com/UCLALibrary/validation-service/validation/csv"
	"github.com/stretchr/testify/assert"
	"slices"
	"testing"
)

// TestVerifyLicense checks if verifyLicense throws the correct errors when given incorrect licenses
func TestVerifyLicense(t *testing.T) {
	check, err := NewLicenseCheck(config.NewProfiles())
	assert.NoError(t, err)

	// genericLocation provides a consistent location for the purposes of test comparison.
	var genericLocation = csv.Location{RowIndex: 1, ColIndex: 0}
	// slices containing expected valid and invalid URLs after table tests run
	var testValids = []string{"http://creativecommons.org/licenses/by-nc/4.0/"}
	var testInvalids = []string{"https://library.ucla.edu", "http://library@edu", "http://ucla.example.edu"}

	tests := []struct {
		name     string
		profile  string
		location csv.Location
		data     [][]string
		result   bool
	}{
		{
			name:     "Valid license with Festerize profile",
			profile:  "festerize",
			location: genericLocation,
			data:     [][]string{{"License"}, {"http://creativecommons.org/licenses/by-nc/4.0/"}},
			result:   true,
		},
		{
			name:     "Invalid license (https prefix) with Festerize profile",
			profile:  "festerize",
			location: genericLocation,
			data:     [][]string{{"License"}, {"https://library.ucla.edu"}},
			result:   false,
		},
		{
			name:     "Invalid license (bad URL format) with Festerize profile",
			profile:  "festerize",
			location: genericLocation,
			data:     [][]string{{"License"}, {"http://library@edu"}},
			result:   false,
		},
		{
			name:     "Invalid license (fake URL) with Festerize profile",
			profile:  "festerize",
			location: genericLocation,
			data:     [][]string{{"License"}, {"http://ucla.example.edu"}},
			result:   false,
		},
		{
			name:     "Valid duplicate license",
			profile:  "festerize",
			location: genericLocation,
			data:     [][]string{{"License"}, {"http://creativecommons.org/licenses/by-nc/4.0/"}},
			result:   true,
		},
		{
			name:     "Invalid duplicate license",
			profile:  "festerize",
			location: genericLocation,
			data:     [][]string{{"License"}, {"https://library.ucla.edu"}},
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
	assert.True(t, slices.Equal(check.valids, testValids))
	assert.True(t, slices.Equal(check.invalids, testInvalids))
}
