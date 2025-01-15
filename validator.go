package main

// Validator interface defines how implementations should be called.
type Validator interface {
	Validate(profile string, location CsvLocation, csvData [][]string) error
}
