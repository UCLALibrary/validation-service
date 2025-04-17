package checks

import (
	"github.com/UCLALibrary/validation-service/validation/csv"
	"github.com/UCLALibrary/validation-service/validation/util"
	"github.com/stretchr/testify/assert"
	"testing"
	"fmt"
)

// TestVerifyLicense checks if verifyLicense throws the correct errors when given incorrect licenses
func TestVerifyLicense(t *testing.T) {
	check, err := NewLicenseCheck(util.NewProfiles())
	assert.NoError(t, err)

	// genericLocation provides a consistent location for the purposes of test comparison.
	var genericLocation = csv.Location{RowIndex: 1, ColIndex: 0}

	tests := []struct {
		name     string
		//license  string
		profile  string
		location csv.Location
		data     [][]string
		result   bool
	}{
		{
			name:     "Valid license with Festerize profile",
			//license:  "http://creativecommons.org/licenses/by-nc/4.0/",
			profile:  "festerize",
			location: genericLocation,
			data:     [][]string{{"License"},{"http://creativecommons.org/licenses/by-nc/4.0/"}},
			result:   true,
		},
		{
			name:     "Invalid license (https prefix) with Festerize profile",
			//license:  "https://library.ucla.edu",
			profile:  "festerize",
			location: genericLocation,
			data:     [][]string{{"License"},{"https://library.ucla.edu"}},
			result:   false,
		},
		{
			name:     "Invalid license (bad URL format) with Festerize profile",
			//license:  "http://library@edu",
			profile:  "festerize",
			location: genericLocation,
			data:     [][]string{{"License"},{"http://library@edu"}},
			result:   false,
		},
		{
			name:     "Invalid license (fake URL) with Festerize profile",
			//license:  "http://ucla.example.edu",
			profile:  "festerize",
			location: genericLocation,
			data:     [][]string{{"License"},{"http://ucla.example.edu"}},
			result:   false,
		},
		{
			name:     "Valid duplicate license",
			//license:  "http://creativecommons.org/licenses/by-nc/4.0/",
			profile:  "festerize",
			location: genericLocation,
			data:     [][]string{{"License"},{"http://creativecommons.org/licenses/by-nc/4.0/"}},
			result:   true,
		},
		{
			name:     "Invalid duplicate license",
			//license:  "https://library.ucla.edu",
			profile:  "festerize",
			location: genericLocation,
			data:     [][]string{{"License"},{"https://library.ucla.edu"}},
			result:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//err := check.verifyLicense(tt.license, tt.profile, tt.location)
			err := check.Validate(tt.profile, tt.location, tt.data)
			if (err != nil && tt.result) || (err == nil && !tt.result) {
				t.Errorf("Expected '%v' response was not found: %v", tt.name, err)
			}
		})
	}
	fmt.Println(check.valids)
	fmt.Println(check.invalids)
}
