//go:build unit

package checks

import (
	"github.com/UCLALibrary/validation-service/validation/profiles"
	"testing"

	"github.com/UCLALibrary/validation-service/validation/csv"
	"github.com/stretchr/testify/assert"
)

// TestFilePathCheck_Validate tests the Validate method on filePath.
func TestFilePathCheck_Validate(t *testing.T) {
	check, err := NewFilePathCheck(profiles.NewProfiles())
	assert.NoError(t, err)

	tests := []struct {
		name        string
		location    csv.Location
		data        [][]string
		expectedErr bool
	}{
		{
			name:        "image exists",
			location:    csv.Location{RowIndex: 1, ColIndex: 1},
			data:        [][]string{{"Random", "File Name"}, {"random", "images/test.jpx"}},
			expectedErr: false,
		},
		{
			name:        "file path prefix stripped successfully",
			location:    csv.Location{RowIndex: 1, ColIndex: 1},
			data:        [][]string{{"Random", "File Name"}, {"random", "Masters/images/test.jpx"}},
			expectedErr: false,
		},
		{
			name:        "file name header does not exist",
			location:    csv.Location{RowIndex: 1, ColIndex: 0},
			data:        [][]string{{"Random", "Header"}, {"Hello", "World"}},
			expectedErr: false,
		},
		{
			name:        "file does not exist",
			location:    csv.Location{RowIndex: 1, ColIndex: 1},
			data:        [][]string{{"Random", "File Name"}, {"random", "random.jpx"}},
			expectedErr: true,
		},
	}

	// Iterate over test cases; fail if there isn't an error when we expect one or if there is an unexpected error
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := check.Validate("DLP Staff", tt.location, tt.data)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.ErrorIs(t, err, nil)
			}
		})
	}
}
