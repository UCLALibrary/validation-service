//go:build unit

package main

import (
	"testing"
)

// TestCsvLocation tests the CsvLocation's behavior
func TestCsvIndex(t *testing.T) {
	// Create an instance of CsvLocation
	location := CsvLocation{RowIndex: 2, ColIndex: 5}

	// Test initial values
	if location.RowIndex != 2 {
		t.Errorf("Expected RowIndex to be 2, got %d", location.RowIndex)
	}

	if location.ColIndex != 5 {
		t.Errorf("Expected ColIndex to be 5, got %d", location.ColIndex)
	}

	// Modify values
	location.RowIndex = 4
	location.ColIndex = 2

	// Test modified values
	if location.RowIndex != 4 {
		t.Errorf("Expected RowIndex to be 4, got %d", location.RowIndex)
	}

	if location.ColIndex != 2 {
		t.Errorf("Expected ColIndex to be 2, got %d", location.ColIndex)
	}
}
