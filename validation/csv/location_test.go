//go:build unit
package csv

// This file tests utilities related to working with CSV data.
import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestIsValidLocation tests whether a supplied csv.Location is valid considering the supplied csvData.
func TestIsValidLocation(t *testing.T) {
	// Simple test data that we check locations against
	csvData := [][]string{
		{"Hello", "World"},
		{"World", "Hello"},
	}

	// A table of test information and expectations
	tests := []struct {
		name       string
		location   Location
		profile    string
		shouldPass bool
	}{
		{"Valid location (row 0, col 1)", Location{RowIndex: 0, ColIndex: 1}, "default", true},
		{"Valid location (row 1, col 0)", Location{RowIndex: 1, ColIndex: 0}, "default", true},
		{"Invalid row index (negative)", Location{RowIndex: -1, ColIndex: 0}, "default", false},
		{"Invalid column index (negative)", Location{RowIndex: 0, ColIndex: -1}, "default", false},
		{"Row index out of bounds", Location{RowIndex: 2, ColIndex: 0}, "default", false},
		{"Column index out of bounds", Location{RowIndex: 0, ColIndex: 2}, "default", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := IsValidLocation(test.location, csvData, test.profile)
			if (err != nil && test.shouldPass) || (err == nil && !test.shouldPass) {
				// Shouldn't have an err on true shouldPass results and should have an err on false shouldPass results
				t.Errorf("IsValidLocation(%v, csvData); found error: %v, expected error: %v", test.location,
					!test.shouldPass, test.shouldPass)
			}
		})
	}
}

// TestGetHeader verifies that GetHeader returns the expected header in the first row
func TestGetHeader(t *testing.T) {
	tests := []struct {
		name        string
		location    Location
		csvData     [][]string
		profile     string
		expected    string
		expectError bool
	}{
		{
			name:     "Valid header retrieval",
			location: Location{ColIndex: 1},
			csvData: [][]string{
				{"ID", "Name", "Age"},
				{"1", "Alice", "30"},
			},
			profile:     "default",
			expected:    "Name",
			expectError: false,
		},
		{
			name:     "First column retrieval",
			location: Location{ColIndex: 0},
			csvData: [][]string{
				{"ID", "Name", "Age"},
				{"1", "Alice", "30"},
			},
			profile:     "default",
			expected:    "ID",
			expectError: false,
		},
		{
			name:     "Out of bounds column index",
			location: Location{ColIndex: 3},
			csvData: [][]string{
				{"ID", "Name", "Age"},
				{"1", "Alice", "30"},
			},
			profile:     "default",
			expected:    "",
			expectError: true,
		},
		{
			name:        "Empty first row",
			location:    Location{ColIndex: 0},
			csvData:     [][]string{{}, {"1", "Alice", "30"}},
			profile:     "default",
			expected:    "",
			expectError: true,
		},
		{
			name:        "Empty CSV data",
			location:    Location{ColIndex: 0},
			csvData:     [][]string{},
			profile:     "default",
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetHeader(tt.location, tt.csvData, tt.profile)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
