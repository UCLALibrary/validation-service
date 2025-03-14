package checks

import (
	"github.com/UCLALibrary/validation-service/validation/csv"
	"testing"
)


// TestVerifyLicense checks if verifyLicense throws the correct errors when given incorrect licenses
func TestVerifyLicense(t *testing.T) {

        // genericLocation provides a consistent location for the purposes of test comparison.
	var genericLocation = csv.Location{}


	tests := []struct {
		name        string
		license     string
		profile     string
		location    csv.Location
		result      bool
	}{
		{
			name:        "Valid license with Festerize profile",
			license:     "http://creativecommons.org/licenses/by-nc/4.0/",
			profile:     "festerize",
			location:    genericLocation,
			result:      true,
		},
		{
			name:        "Invalid license (https prefix) with Festerize profile",
			license:     "https://library.ucla.edu",
			profile:     "festerize",
			location:    genericLocation,
			result:      false,
		},
		{
			name:        "Invalid license (bad URL format) with Festerize profile",
			license:     "http://library@edu",
			profile:     "festerize",
			location:    genericLocation,
			result:      false,
		},
		{
			name:        "Invalid license (fake URL) with Festerize profile",
			license:     "http://ucla.example.edu",
			profile:     "festerize",
			location:    genericLocation,
			result:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := verifyLicense(tt.license, tt.profile, tt.location)
			if (err != nil && tt.result) || (err == nil && !tt.result) {
				t.Errorf("Expected '%v' response was not found: %v", tt.name, err)
			}
		})
	}
}
