//go:build unit

package checks

import (
	"testing"

	"github.com/UCLALibrary/validation-service/validation/csv"
	"github.com/UCLALibrary/validation-service/validation/util"
	"github.com/stretchr/testify/assert"
)

// TestValidateItemSeq tests the Validate method on ObjTypeCheck.
func TestValidateItemSeq(t *testing.T) {
	check, err := NewItemSeqCheck(util.NewProfiles())
	assert.NoError(t, err)

	tests := []struct {
		name        string
		location    csv.Location
		data        [][]string
		expectedErr bool
	}{
		{
			name:        "valid pos int",
			location:    csv.Location{RowIndex: 1, ColIndex: 0},
			data:        [][]string{{"Item Sequence"}, {"5"}},
			expectedErr: false,
		},
		{
			name:        "valid null value",
			location:    csv.Location{RowIndex: 1, ColIndex: 0},
			data:        [][]string{{"Item Sequence"}, {""}},
			expectedErr: false,
		},
		{
			name:        "invalid neg int",
			location:    csv.Location{RowIndex: 1, ColIndex: 0},
			data:        [][]string{{"Item Sequence"}, {"-4"}},
			expectedErr: true,
		},
		{
			name:        "invalid non int",
			location:    csv.Location{RowIndex: 1, ColIndex: 0},
			data:        [][]string{{"Item Sequence"}, {"1a"}},
			expectedErr: true,
		},
	}

	// Iterate over test cases; fail if there isn't an error when we expect one or if there is an unexpected error
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := check.Validate("default", tt.location, tt.data)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.ErrorIs(t, err, nil)
			}
		})
	}
}
