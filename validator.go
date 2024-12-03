package main

type Validator interface {
	Validate(profile string, location CsvLocation, csvData [][]string) error
}
