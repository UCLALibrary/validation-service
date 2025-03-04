//go:build unit

// Package csv has structures and utilities useful for working with CSVs.
//
// This file tests CSV errors.
package csv

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestNewError_NoParent tests creating a validation.Error without a parent error
func TestNewError_NoParent(t *testing.T) {
	var valErr *Error

	location := Location{RowIndex: 5, ColIndex: 10}
	err := NewError("Invalid value", location, "default")

	// Ensure the error is of type *validation.Error
	assert.ErrorAs(t, err, &valErr)

	// Check the fields
	assert.Equal(t, "Invalid value", valErr.Message)
	assert.Equal(t, 5, valErr.Location.RowIndex)
	assert.Equal(t, 10, valErr.Location.ColIndex)
	assert.Equal(t, "default", valErr.Profile)
	assert.Nil(t, valErr.ParentErr)

	// Check the error string format
	expectedMsg := "Error: Invalid value (Row: 5, Col: 10) [profile: default]"
	assert.Equal(t, expectedMsg, valErr.Error())
}

// TestNewError_WithParent tests creating a validation.Error with a parent error
func TestNewError_WithParent(t *testing.T) {
	var valErr *Error

	location := Location{RowIndex: 3, ColIndex: 7}
	parentErr := errors.New("underlying parse error")
	err := NewError("Invalid syntax", location, "default", parentErr)

	// Ensure the error is of type *validation.Error
	assert.ErrorAs(t, err, &valErr)

	// Check the fields
	assert.Equal(t, "Invalid syntax", valErr.Message)
	assert.Equal(t, 3, valErr.Location.RowIndex)
	assert.Equal(t, 7, valErr.Location.ColIndex)
	assert.Equal(t, "default", valErr.Profile)
	assert.Equal(t, parentErr, valErr.ParentErr)

	// Check the error string format
	expected := fmt.Sprintf("Error: Invalid syntax (Row: 3, Col: 7) [profile: default] Cause: %s", parentErr.Error())
	assert.Equal(t, expected, valErr.Error())
}

// TestError_Unwrap tests that Unwrap works properly for error wrapping
func TestError_Unwrap(t *testing.T) {
	var valErr *Error

	location := Location{RowIndex: 2, ColIndex: 4}
	parentErr := errors.New("file read error")
	err := NewError("Missing delimiter", location, "default", parentErr)

	// Ensure errors.Is works correctly
	assert.True(t, errors.Is(err, parentErr))

	// Ensure errors.As works correctly
	assert.ErrorAs(t, err, &valErr)
	assert.Equal(t, parentErr, valErr.ParentErr)
}

// TestError_ErrorFormatting ensures the error message format is correct
func TestError_ErrorFormatting(t *testing.T) {
	tests := []struct {
		name        string
		message     string
		location    Location
		profile     string
		parentErr   error
		expectedStr string
	}{
		{
			name:        "No parent error",
			message:     "Invalid header",
			location:    Location{RowIndex: 1, ColIndex: 2},
			profile:     "default",
			parentErr:   nil,
			expectedStr: "Error: Invalid header (Row: 1, Col: 2) [profile: default]",
		},
		{
			name:        "With parent error",
			message:     "Incorrect format",
			location:    Location{RowIndex: 4, ColIndex: 8},
			profile:     "default",
			parentErr:   errors.New("unexpected EOF"),
			expectedStr: "Error: Incorrect format (Row: 4, Col: 8) [profile: default] Cause: unexpected EOF",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := NewError(tc.message, tc.location, tc.profile, tc.parentErr)

			var valErr *Error
			assert.ErrorAs(t, err, &valErr)
			assert.Equal(t, tc.expectedStr, valErr.Error())
		})
	}
}

// TestError_Is tests the Is() func in csv.Error.
func TestError_Is(t *testing.T) {
	tests := []struct {
		name     string
		err1     *Error
		err2     error
		expected bool
	}{
		{
			name: "Identical errors",
			err1: &Error{
				Message:  "Invalid format",
				Location: Location{RowIndex: 1, ColIndex: 2},
				Profile:  "default",
			},
			err2: &Error{
				Message:  "Invalid format",
				Location: Location{RowIndex: 1, ColIndex: 2},
				Profile:  "default",
			},
			expected: true,
		},
		{
			name: "Different message",
			err1: &Error{
				Message:  "Invalid format",
				Location: Location{RowIndex: 1, ColIndex: 2},
				Profile:  "default",
			},
			err2: &Error{
				Message:  "Wrong data type",
				Location: Location{RowIndex: 1, ColIndex: 2},
				Profile:  "default",
			},
			expected: false,
		},
		{
			name: "Different location",
			err1: &Error{
				Message:  "Invalid format",
				Location: Location{RowIndex: 1, ColIndex: 2},
				Profile:  "default",
			},
			err2: &Error{
				Message:  "Invalid format",
				Location: Location{RowIndex: 3, ColIndex: 4},
				Profile:  "default",
			},
			expected: false,
		},
		{
			name: "Different profile",
			err1: &Error{
				Message:  "Invalid format",
				Location: Location{RowIndex: 1, ColIndex: 2},
				Profile:  "default",
			},
			err2: &Error{
				Message:  "Invalid format",
				Location: Location{RowIndex: 1, ColIndex: 2},
				Profile:  "custom",
			},
			expected: false,
		},
		{
			name: "Comparing with a different error type",
			err1: &Error{
				Message:  "Invalid format",
				Location: Location{RowIndex: 1, ColIndex: 2},
				Profile:  "default",
			},
			err2:     errors.New("random error"),
			expected: false,
		},
		{
			name: "Comparing with nil",
			err1: &Error{
				Message:  "Invalid format",
				Location: Location{RowIndex: 1, ColIndex: 2},
				Profile:  "default",
			},
			err2:     nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err1.Is(tt.err2)
			assert.Equal(t, tt.expected, result, "Expected Is() to return %v for %s", tt.expected, tt.name)
		})
	}
}
