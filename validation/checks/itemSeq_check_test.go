//go:build unit

package checks

import (
	"github.com/UCLALibrary/validation-service/validation/config"
	"testing"

	"github.com/UCLALibrary/validation-service/validation/csv"
	"github.com/stretchr/testify/assert"
)

// TestValidateItemSeq tests the Validate method on ObjTypeCheck.
func TestValidateItemSeq(t *testing.T) {
	check, err := NewItemSeqCheck(config.NewProfiles())
	assert.NoError(t, err)

	tests := []struct {
		name        string
		location    csv.Location
		data        [][]string
		expectedErr bool
	}{
		{
			name:        "valid pos int with page Object Type",
			location:    csv.Location{RowIndex: 1, ColIndex: 0},
			data:        [][]string{{"Item Sequence", "Object Type"}, {"5", "Page"}},
			expectedErr: false,
		},
		{
			name:        "valid pos int with non Page Object Type",
			location:    csv.Location{RowIndex: 1, ColIndex: 0},
			data:        [][]string{{"Item Sequence", "Object Type"}, {"5", "notPage"}},
			expectedErr: false,
		},
		{
			name:        "valid null value",
			location:    csv.Location{RowIndex: 1, ColIndex: 0},
			data:        [][]string{{"Item Sequence", "Object Type"}, {"", "notPage"}},
			expectedErr: false,
		},
		{
			name:        "valid pos int with no Object Type (should expect no error)",
			location:    csv.Location{RowIndex: 1, ColIndex: 0},
			data:        [][]string{{"Item Sequence"}, {"5"}},
			expectedErr: false,
		},
		{
			name:        "valid int with no Object Type (should expect error)",
			location:    csv.Location{RowIndex: 1, ColIndex: 0},
			data:        [][]string{{"Item Sequence"}, {"-5"}},
			expectedErr: true,
		},
		{
			name:        "invalid null value",
			location:    csv.Location{RowIndex: 1, ColIndex: 0},
			data:        [][]string{{"Item Sequence", "Object Type"}, {"", "Page"}},
			expectedErr: true,
		},
		{
			name:        "invalid neg int",
			location:    csv.Location{RowIndex: 1, ColIndex: 0},
			data:        [][]string{{"Item Sequence", "Object Type"}, {"-4", "random"}},
			expectedErr: true,
		},
		{
			name:        "invalid non int",
			location:    csv.Location{RowIndex: 1, ColIndex: 0},
			data:        [][]string{{"Item Sequence", "Object Type"}, {"1a", "random"}},
			expectedErr: true,
		},
		{
			name:        "invalid non int",
			location:    csv.Location{RowIndex: 1, ColIndex: 0},
			data:        [][]string{{"Item Sequence", "Object Type"}, {"1a", "Page"}},
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
