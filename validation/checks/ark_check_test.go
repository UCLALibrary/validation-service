package checks

import (
	"flag"
	"fmt"
	"github.com/UCLALibrary/validation-service/testflags"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/multierr"
)

// TestCheck loads the flags for the tests in the 'check' package.
func TestCheck(t *testing.T) {
	flag.Parse()
	fmt.Printf("%s's log level: %s\n", t.Name(), *testflags.LogLevel)
}

// TestVerifyArk checks if verifyARK throws the correct errors when given incorrect ARKs
func TestVerifyArk(t *testing.T) {
	tests := []struct {
		name        string
		ark         string
		profile     string
		expectError bool
		expectedErr error
	}{
		{
			name:        "Valid ARK with default profile",
			ark:         "ark:/21198/xyz123",
			profile:     "default",
			expectError: false,
		},
		{
			name:        "Valid ARK with qualifier",
			ark:         "ark:/12345/xyz123?version=2",
			profile:     "custom",
			expectError: false,
		},
		{
			name:        "Valid ARK with non-default profile",
			ark:         "ark:/12345/abc456",
			profile:     "custom",
			expectError: false,
		},
		{
			name:        "Invalid ARK - missing ark:/ prefix",
			ark:         "12345/xyz123",
			profile:     "default",
			expectError: true,
			expectedErr: noPrefixErr,
		},
		{
			name:        "Invalid ARK structure no object identifier",
			ark:         "ark:/12345",
			profile:     "random",
			expectError: true,
			expectedErr: noObjIdErr,
		},
		{
			name:        "Invalid NAAN - less than 5 digits",
			ark:         "ark:/123/",
			profile:     "default",
			expectError: true,
			expectedErr: multierr.Combine(naanTooShortErr, naanProfileErr, noObjIdErr),
		},
		{
			name:        "Invalid NAAN for default profile",
			ark:         "ark:/12345/xyz123",
			profile:     "default",
			expectError: true,
			expectedErr: naanProfileErr,
		},
		{
			name:        "Invalid object identifier",
			ark:         "ark:/12345/my identifier",
			profile:     "random",
			expectError: true,
			expectedErr: invalidObjIdErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := verifyARK(tt.ark, tt.profile)

			if tt.expectError {
				assert.Error(t, err)

				// If expectedErr is a combined error, check each error individually
				if multiErr, ok := tt.expectedErr.(interface{ Unwrap() []error }); ok {
					for _, expectedErr := range multiErr.Unwrap() {
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
