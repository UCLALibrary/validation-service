package csv

// This file tests the CSV validation report and its components.
import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"testing"
	"time"
)

// TestNewReport tests creating a new validation report with the NewReport function.
func TestNewReport(t *testing.T) {
	logger := zap.NewNop() // Use a no-op logger for testing

	tests := []struct {
		name       string
		multiErr   error
		csvData    [][]string
		expected   *Report
		shouldFail bool
	}{
		{
			name: "Multiple validation errors",
			multiErr: multierr.Combine(
				&Error{Message: "Invalid value", Location: Location{ColIndex: 1, RowIndex: 2}, Profile: "default"},
				&Error{Message: "Missing field", Location: Location{ColIndex: 2, RowIndex: 2}, Profile: "default"},
			),
			csvData: [][]string{
				{"Header1", "Header2", "Header3"},
				{"Row1Col1", "Row1Col2", "Row1Col3"},
				{"Row2Col1", "Row2Col2", "Row2Col3"},
				{"Row3Col1", "Row3Col2", "Row3Col3"},
			},
			expected: &Report{
				Profile: "default",
				Warnings: []Warning{
					{
						Message:  "Error: Invalid value",
						Header:   "Header2",
						ColIndex: 1,
						RowIndex: 2,
						Value:    "Row2Col2",
					},
					{
						Message:  "Error: Missing field",
						Header:   "Header3",
						ColIndex: 2,
						RowIndex: 2,
						Value:    "Row2Col3",
					},
				},
			},
		},
		{
			name:     "No validation errors",
			multiErr: nil,
			csvData:  [][]string{{"Header1", "Header2"}},
			expected: &Report{
				Profile:  "",
				Warnings: []Warning{},
			},
		},
		{
			name: "Unexpected non-validation error",
			multiErr: multierr.Combine(
				errors.New("unexpected error"),
			),
			csvData: [][]string{{"Header1", "Header2"}},
			expected: &Report{
				Profile:  "",
				Warnings: []Warning{}, // Unexpected errors are logged, but do not produce warnings
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report, err := NewReport(tt.multiErr, tt.csvData, logger)

			if tt.shouldFail {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, report)

			// Check profile
			assert.Equal(t, tt.expected.Profile, report.Profile)

			// Check warnings length
			assert.Equal(t, len(tt.expected.Warnings), len(report.Warnings))

			// Check each warning
			for index, warning := range report.Warnings {
				assert.Equal(t, tt.expected.Warnings[index].Message, warning.Message)
				assert.Equal(t, tt.expected.Warnings[index].Header, warning.Header)
				assert.Equal(t, tt.expected.Warnings[index].ColIndex, warning.ColIndex)
				assert.Equal(t, tt.expected.Warnings[index].RowIndex, warning.RowIndex)
				assert.Equal(t, tt.expected.Warnings[index].Value, warning.Value)
			}
		})
	}
}

// Tests serializing the Report to JSON with SerializeReport function.
func TestSerializeReport(t *testing.T) {
	var deserialized Report

	report := &Report{
		Profile: "default",
		Time:    time.Now(),
		Warnings: []Warning{
			{Message: "Invalid value", ColIndex: 1, RowIndex: 2, Value: "Row2Col2"},
			{Message: "Missing field", ColIndex: 3, RowIndex: 4, Value: "Row4Col3"},
		},
	}

	jsonData, err := SerializeReport(report)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	// Deserialize and check contents
	err = json.Unmarshal([]byte(jsonData), &deserialized)
	assert.NoError(t, err)

	assert.Equal(t, report.Profile, deserialized.Profile)
	assert.Equal(t, len(report.Warnings), len(deserialized.Warnings))

	for index := range report.Warnings {
		assert.Equal(t, report.Warnings[index], deserialized.Warnings[index])
	}
}

// Tests SerializeReport with an empty report.
func TestSerializeEmptyReport(t *testing.T) {
	var deserialized Report

	report := &Report{}

	jsonData, err := SerializeReport(report)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	// Deserialize and check contents
	err = json.Unmarshal([]byte(jsonData), &deserialized)
	assert.NoError(t, err)

	assert.Equal(t, report.Profile, deserialized.Profile)
	assert.Equal(t, len(report.Warnings), len(deserialized.Warnings))
}
