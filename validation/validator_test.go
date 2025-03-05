//go:build unit

// Package validation provides tools to validate CSV data.
//
// This file provides tests of the Validator.
package validation

import (
	"fmt"
	"github.com/UCLALibrary/validation-service/validation/csv"
	"testing"

	"github.com/stretchr/testify/assert"
)

var location = csv.Location{RowIndex: 2, ColIndex: 5}
var csvData = [][]string{
	{"header0", "header1", "header2", "header3", "header4", "header5"},
	{"value10", "value11", "value12", "value13", "value14", "value15"},
	{"value20", "value21", "value22", "value23", "value24", "value25"},
}

// MockValidator implements the Validator interface so it can be tested.
type MockValidator struct {
	ValidateFunc func(profile string, location csv.Location, csvData [][]string) error
}

// Validate forwards the method call to the function stored in ValidateFunc, passing along the input arguments and
// returning the possible error.
func (mock *MockValidator) Validate(profile string, location csv.Location, csvData [][]string) error {
	return mock.ValidateFunc(profile, location, csvData)
}

// TestValidatorSuccess tests the Validate interface with a mock validator and expects a successful result.
func TestValidatorSuccess(t *testing.T) {
	mock := &MockValidator{
		ValidateFunc: func(profile string, location csv.Location, csvData [][]string) error {
			assert.Equal(t, "profile1", profile)
			assert.Equal(t, "value24", csvData[2][4])
			return nil
		},
	}

	err := mock.Validate("profile1", location, csvData)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
}

// TestValidatorError tests the Validate interface with a mock validator and expects an error.
func TestValidatorError(t *testing.T) {
	mock := &MockValidator{
		ValidateFunc: func(profile string, location csv.Location, csvData [][]string) error {
			if profile != "profile1" {
				return fmt.Errorf("Expected 'profile1' as profile, but found '%s'", profile)
			}

			return nil
		},
	}

	// We pass 'profile2' instead of the expected 'profile1'
	err := mock.Validate("profile2", location, csvData)
	if err != nil {
		// We're expecting an error, but check its message just to confirm it's the one we're expecting
		assert.EqualError(t, err, "Expected 'profile1' as profile, but found 'profile2'",
			"Expected error message does not match")
	}
}
