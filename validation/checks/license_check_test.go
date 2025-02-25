package checks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/multierr"
)

// TestVerifyLicense checks if verifyLicense throws the correct errors when given incorrect licenses
func TestVerifyLicense(t *testing.T) {
	tests := []struct {
		name        string
		ark         string
		profile     string
		expectError bool
		expectedErr error
	}{
		{
			name:        "Valid license with Festerize profile",
			license:     "ihttp://library.ucla.edu",
			profile:     "festerize",
			expectError: false,
		},
		{
			name:        "Invalid license (https prefix) with Festerize profile",
			license:     "ihttps://library.ucla.edu",
			profile:     "festerize",
			expectError: true,
		},
		{
			name:        "Invalid license (bad URL format) with Festerize profile",
			license:     "ihttp://libraryDOTucla.edu",
			profile:     "festerize",
			expectError: true,
		},
		{
			name:        "Invalid license (fake URL) with Festerize profile",
			license:     "ihttp://www.example.edu",
			profile:     "festerize",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := verifyLicense(tt.ark, tt.profile)

			if tt.expectError {
				assert.Error(t, err)

				// If expectedErr is a combined error, check each error individually
				if merr, ok := tt.expectedErr.(interface{ Unwrap() []error }); ok {
					for _, expectedErr := range merr.Unwrap() {
						assert.ErrorIs(t, err, expectedErr, "expected error: %v", expectedErr)
					}
				} else {
					assert.ErrorIs(t, err, tt.expectedErr)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
