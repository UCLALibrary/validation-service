//go:build unit

package checks

import (
	"github.com/UCLALibrary/validation-service/validation/csv"
	"github.com/UCLALibrary/validation-service/validation/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestVerifyMeta checks if MediaMetaCheck.Validate throws the correct errors when given missing media metadata
func TestVerifyMeta(t *testing.T) {
	check, err := NewMediaMetaCheck(util.NewProfiles())
	assert.NoError(t, err)

	// genericLocation provides a consistent location for the purposes of test comparison.
	var startLocation = csv.Location{RowIndex: 1, ColIndex: 0}

	tests := []struct {
		name     string
		profile  string
		location csv.Location
		data     [][]string
		result   bool
	}{
		{
			name:     "Non-Fester profile, skips cell evaluation",
			profile:  "bucketeer",
			location: startLocation,
			data:     [][]string{{"Type.typeOfResource"},{"car"}},
			result:   true,
		},
		{
			name:     "Fester profile, non-media resource type",
			profile:  "fester",
			location: startLocation,
			data:     [][]string{{"Type.typeOfResource"},{"img"}},
			result:   true,
		},
		{
			name:     "Fester profile, media resource, all metadata fields",
			profile:  "fester",
			location: startLocation,
			data:     [][]string{{{"Type.typeOfResource"},{"mov"}},{{"media.width"},{"5"}},{{"media.height"},{"7"}},{{"media.duration"},{"10"}},{{"media.format"},{"mov"}}},
			result:   true,
		},
		{
			name:     "Fester profile, media resource, empty fields",
			profile:  "fester",
			location: startLocation,
			data:     [][]string{{"Type.typeOfResource"},{"mov"},{"media.width"},{"5"},{"media.height"},{""},{"media.duration"},{"10"},{"media.format"},{""},},
			result:   false,
		},
		{
			name:     "Fester profile, media resource, missing some metadat columns",
			profile:  "fester",
			location: startLocation,
			data:     [][]string{{"Type.typeOfResource"},{"mov"},{"media.width"},{"5"},{"media.height"},{"7"}},
			result:   false,
		},
		{
			name:     "Fester profile, media resource, missing all metadata columns",
			profile:  "fester",
			location: startLocation,
			data:     [][]string{{"Type.typeOfResource"},{"mov"}},
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
}
