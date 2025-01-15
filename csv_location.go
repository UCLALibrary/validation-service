package main

// CsvLocation represents an index-based location in a CSV file.
//
// The RowIndex specifies the row number and the ColIndex specifies the column number.
type CsvLocation struct {
	// RowIndex is the zero-based index of the row in the CSV file.
	RowIndex int

	// ColIndex is the zero-based index of the column in the CSV file.
	ColIndex int
}
