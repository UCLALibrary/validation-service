// Package csv has structures and utilities useful for working with CSVs.
//
// This file defines CSV errors.
package csv

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Error creates an error that can store discreet CSV location information and, optionally, a parent error.
type Error struct {
	ParentErr error
	Message   string
	Location  Location
	Profile   string
}

// Error implements an interface that allows an error to be returned as a string.
func (err *Error) Error() string {
	if err.ParentErr != nil {
		// Wrapped exceptions will have a duplicate label that we can strip
		cause := strings.TrimPrefix(err.ParentErr.Error(), "Error: ")

		// There will also be duplicate location and profile info that can be stripped
		regex := regexp.MustCompile(`\s*\(Row: \d+, Col: \d+\) \[profile: .*?\]$`)
		cause = regex.ReplaceAllString(cause, "") // 'All' but our pattern specifies last occurrence

		return fmt.Sprintf("Error: %s (Row: %d, Col: %d) [profile: %s] Cause: %s",
			err.Message, err.Location.RowIndex, err.Location.ColIndex, err.Profile, cause)
	}

	return fmt.Sprintf("Error: %s (Row: %d, Col: %d) [profile: %s]",
		err.Message, err.Location.RowIndex, err.Location.ColIndex, err.Profile)
}

// String outputs a string version of the error for display to non-programmers.
func (err *Error) String() string {
	if err.ParentErr != nil {
		// Wrapped exceptions will have a duplicate label that we can strip
		cause := strings.TrimPrefix(err.ParentErr.Error(), "Error: ")

		// We strip location and profile info when outputting string form of an error
		regex := regexp.MustCompile(`\s*\(Row: \d+, Col: \d+\) \[profile: .*?\]`)
		cause = regex.ReplaceAllString(cause, "") // 'All' means all for String()

		return fmt.Sprintf("Error: %s [Cause: %s]", err.Message, cause)
	}

	return fmt.Sprintf("Error: %s", err.Message)
}

// Is checks two errors (this one and a supplied one) for equality.
func (err *Error) Is(other error) bool {
	var target *Error

	if errors.As(other, &target) {
		return err.Message == target.Message && err.Location == target.Location && err.Profile == target.Profile
	}

	return false
}

// Unwrap ensures report.Error compatibility with errors.Is() and errors.As().
func (err *Error) Unwrap() error {
	return err.ParentErr
}

// NewError creates a new report.Error from the supplied CSV data location and profile, with an optional parent error.
func NewError(message string, location Location, profile string, err ...error) error {
	var parentErr error

	if len(err) > 0 {
		parentErr = err[0] // Use the first error if provided
	}

	return &Error{
		ParentErr: parentErr,
		Message:   message,
		Location:  location,
		Profile:   profile,
	}
}
