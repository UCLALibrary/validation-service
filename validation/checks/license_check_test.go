package checks

import (
	"github.com/UCLALibrary/validation-service/validation/csv"
	"testing"
)

// testLocation provides a consistent location for the purposes of test comparison.

// TestVerifyLicense checks if verifyLicense throws the correct errors when given incorrect licenses
func TestVerifyLicense(t *testing.T) {

	tests := []struct {
		name        string
		license     string
		profile     string
		location    csv.Location
		expectError bool
	}{
		{
			name:        "Valid license with Festerize profile",
			license:     "http://creativecommons.org/licenses/by-nc/4.0/",
			profile:     "festerize",
			location:    testLocation,
			expectError: false,
		},
		{
			name:        "Invalid license (https prefix) with Festerize profile",
			license:     "https://library.ucla.edu",
			profile:     "festerize",
			location:    testLocation,
			expectError: true,
		},
		{
			name:        "Invalid license (bad URL format) with Festerize profile",
			license:     "http://library@edu",
			profile:     "festerize",
			location:    testLocation,
			expectError: true,
		},
		{
			name:        "Invalid license (fake URL) with Festerize profile",
			license:     "http://ucla.example.edu",
			profile:     "festerize",
			location:    testLocation,
			expectError: true,
		},
		{
			name:        "Invalid license (no body) with Festerize profile",
			license:     "http://about:blank",
			profile:     "festerize",
			location:    testLocation,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := verifyLicense(tt.license, tt.profile, tt.location)
			if (err != nil && tt.expectError) || (err == nil && !tt.expectError) {
				t.Errorf("Expected '%v' response was not found: %v", tt.name, err)
			}
		})
	}
}
