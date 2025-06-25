//go:build unit

package checks

import (
	"github.com/UCLALibrary/validation-service/validation/profiles"
	"testing"

	"github.com/UCLALibrary/validation-service/errors"
	"github.com/UCLALibrary/validation-service/validation/csv"
	"github.com/stretchr/testify/assert"
	"go.uber.org/multierr"
)

// testLocation provides a consistent location for the purposes of test comparison.
var testLocation = csv.Location{}

// TestVerifyARK checks if verifyARK throws the correct errors when given incorrect ARKs
func TestVerifyARK(t *testing.T) {
	check, err := NewARKCheck(profiles.NewProfiles())
	assert.NoError(t, err)

	tests := []struct {
		name        string
		ark         string
		location    csv.Location
		profile     string
		expectError bool
		expectedErr error
	}{
		{
			name:        "Valid ARK with default profile",
			ark:         "ark:/21198/xyz123",
			location:    testLocation,
			profile:     "DLP Staff",
			expectError: false,
		},
		{
			name:        "Valid ARK with CDL NAAN",
			ark:         "ark:/13030/xyz123",
			location:    testLocation,
			profile:     "DLP Staff",
			expectError: false,
		},
		{
			name:        "Valid ARK with CDL NAAN",
			ark:         "ark:/13030/xyz123",
			location:    testLocation,
			profile:     "DLP Staff",
			expectError: false,
		},
		{
			name:        "Valid ARK with qualifier",
			ark:         "ark:/21198/xyz123?version=2",
			location:    testLocation,
			profile:     "DLP Staff",
			expectError: false,
		},
		{
			name:        "Valid ARK with non-default profile",
			ark:         "ark:/21198/abc456",
			location:    testLocation,
			profile:     "Test",
			expectError: false,
		},
		{
			name:        "Invalid ARK - missing ark:/ prefix",
			ark:         "12345/xyz123",
			location:    testLocation,
			profile:     "DLP Staff",
			expectError: true,
			expectedErr: csv.NewError(errors.NoPrefixErr, testLocation, "DLP Staff"),
		},
		{
			name:        "Invalid ARK structure no object identifier",
			ark:         "ark:/21198",
			location:    testLocation,
			profile:     "DLP Staff",
			expectError: true,
			expectedErr: csv.NewError(errors.NoObjIDErr, testLocation, "DLP Staff"),
		},
		{
			name:        "Invalid NAAN - less than 5 digits",
			ark:         "ark:/123/",
			location:    testLocation,
			profile:     "DLP Staff",
			expectError: true,
			expectedErr: multierr.Combine(
				csv.NewError(errors.NaanTooShortErr, testLocation, "DLP Staff"),
				csv.NewError(errors.NaanProfileErr, testLocation, "DLP Staff"),
				csv.NewError(errors.NoObjIDErr, testLocation, "DLP Staff"),
			),
		},
		{
			name:        "Invalid NAAN for default profile",
			ark:         "ark:/12345/xyz123",
			location:    testLocation,
			profile:     "DLP Staff",
			expectError: true,
			expectedErr: csv.NewError(errors.NaanProfileErr, testLocation, "DLP Staff"),
		},
		{
			name:        "Invalid object identifier",
			ark:         "ark:/21198/my identifier",
			location:    testLocation,
			profile:     "DLP Staff",
			expectError: true,
			expectedErr: csv.NewError(errors.InvalidObjIDErr, testLocation, "DLP Staff"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := check.verifyARK(tt.ark, tt.location, tt.profile)

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
