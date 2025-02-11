// Package csvutils has structures and utilities useful for working with CSVs.
package csvutils

import "fmt"

// Location represents an index-based location in a CSV file.
//
// The RowIndex specifies the row number and the ColIndex specifies the column number.
type Location struct {
	// RowIndex is the zero-based index of the row in the CSV file.
	RowIndex int

	// ColIndex is the zero-based index of the column in the CSV file.
	ColIndex int
}

// IsValidLocation checks if a Location is within bounds of our CSV Location struct.
func IsValidLocation(location Location, csvData [][]string) error {
	if location.RowIndex < 0 || location.RowIndex >= len(csvData) {
		return fmt.Errorf("row %d is out of bounds", location.RowIndex)
	}

	if location.ColIndex < 0 || location.ColIndex >= len(csvData[location.RowIndex]) {
		return fmt.Errorf("column %d is out of bounds", location.ColIndex)
	}

	// Supplied Location is valid for the supplied csvData
	return nil
}
