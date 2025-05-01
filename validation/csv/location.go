package csv

import (
	"fmt"
)

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
func IsValidLocation(location Location, csvData [][]string, profile string) error {
	if location.RowIndex < 0 || location.RowIndex >= len(csvData) {
		message := fmt.Sprintf("row %d is out of bounds", location.RowIndex)
		return NewError(message, location, profile)
	}

	if location.ColIndex < 0 || location.ColIndex >= len(csvData[location.RowIndex]) {
		message := fmt.Sprintf("column %d is out of bounds", location.ColIndex)
		return NewError(message, location, profile)
	}

	// Supplied Location is valid for the supplied csvData
	return nil
}

// GetHeader returns the header within the first row given the index
func GetHeader(location Location, csvData [][]string, profile string) (string, error) {
	if err := IsValidLocation(location, csvData, profile); err != nil {
		return "", err
	}

	index := location.ColIndex

	// Ensure the first row exists
	headers := csvData[0]
	if len(headers) == 0 {
		return "", NewError("the first row of csvData is empty", location, profile)
	}

	return headers[index], nil
}

// GetHeaderIndex returns the index position for the column that contains the supplied header name.
func GetHeaderIndex(header string, location Location, csvData [][]string, profile string) (int, error) {
	if err := IsValidLocation(location, csvData, profile); err != nil {
		return -1, err
	}

	// Iterate through the header row (first row) to find the column that matches the supplied header
	for index, colHeader := range csvData[0] {
		if header == colHeader {
			return index, nil
		}
	}

	return -1, NewError(fmt.Sprintf("supplied header '%s' was not located in first row", header), location, profile)
}

// GetRowValue gets the value of the supplied header for the row of the cell being checked
func GetRowValue(header string, location Location, csvData [][]string, profile string) (string, error) {
	colIndex, err := GetHeaderIndex(header, location, csvData, profile)
	if err != nil {
		return "", NewError(fmt.Sprintf("conditional field '%s' was not found", header), location, profile, err)
	}

	return csvData[location.RowIndex][colIndex], nil
}
