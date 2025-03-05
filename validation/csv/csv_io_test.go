//go:build unit

// Package csv has structures and utilities useful for working with CSVs.
//
// This file provides tests for readers and writers of CSV data.
package csv

import (
	"encoding/csv"
	"go.uber.org/zap"
	"os"
	"testing"

	"go.uber.org/zap/zaptest"
)

// TestWriteFile checks if WriteFile correctly writes data to a CSV file.
func TestWriteFile(t *testing.T) {
	// Create a logger to use when testing
	logger := zaptest.NewLogger(t)

	// Create a temp file to write
	tmpFile, err := os.CreateTemp("", "test_csv_file-*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func(name string) {
		if err := os.Remove(name); err != nil {
			logger.Error("Failed to remove temp file: %v", zap.Error(err))
		}
	}(tmpFile.Name())

	// Sample data that we'll write to our test file
	data := [][]string{
		{"Item ARK", "Parent ARK", "Title"},
		{"ark:/12345/fk4xg9b6k", "ark:/12345/fk4wt7v3k", "Item record"},
		{"ark:/12345/fk4wt7v3k", "ark:/12345/fk4ww2r4r", "Collection record"},
	}

	// Call the function to write the CSV data
	if err := WriteFile(tmpFile.Name(), data, logger); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	// Open the temp file that we just wrote
	file, openErr := os.Open(tmpFile.Name())
	if openErr != nil {
		t.Fatalf("Failed to reopen temp file: %v", openErr)
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			logger.Error("Failed to close temp file", zap.Error(err))
		}
	}(file)

	// Read the temp file to verify its contents
	reader := csv.NewReader(file)
	readData, readErr := reader.ReadAll()
	if readErr != nil {
		t.Fatalf("Failed to read back CSV file: %v", readErr)
	}

	// Check if the written data matches the original data
	if len(readData) != len(data) {
		t.Fatalf("Expected %d rows, got %d", len(data), len(readData))
	}
	for i := range data {
		for j := range data[i] {
			if readData[i][j] != data[i][j] {
				t.Errorf("Mismatch at row %d, column %d: expected %s, got %s",
					i, j, data[i][j], readData[i][j])
			}
		}
	}
}

// TestReadFile checks if ReadFile correctly reads data from a CSV file.
func TestReadFile(t *testing.T) {
	// Create a logger to use when testing
	logger := zaptest.NewLogger(t)

	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "test_csv_file-*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func(name string) {
		if err := os.Remove(name); err != nil {
			logger.Error("Failed to remove temp file: %v", zap.Error(err))
		}
	}(tmpFile.Name())

	// Sample data to write to our test CSV file
	data := [][]string{
		{"Item ARK", "Parent ARK", "Title"},
		{"ark:/12345/fk4xg9b6k", "ark:/12345/fk4wt7v3k", "Item record"},
		{"ark:/12345/fk4wt7v3k", "ark:/12345/fk4ww2r4r", "Collection record"},
	}

	// Manually write the sample data to the file
	file, fileErr := os.Create(tmpFile.Name())
	if fileErr != nil {
		t.Fatalf("Failed to open temp file for writing: %v", fileErr)
	}
	writer := csv.NewWriter(file)
	err = writer.WriteAll(data)
	writer.Flush()
	if err := file.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Call the ReadFile function
	readData, readFileErr := ReadFile(tmpFile.Name(), logger)
	if readFileErr != nil {
		t.Fatalf("ReadFile failed: %v", readFileErr)
	}

	// Check if the read data matches the expected data
	if len(readData) != len(data) {
		t.Fatalf("Expected %d rows, got %d", len(data), len(readData))
	}
	for i := range data {
		for j := range data[i] {
			if readData[i][j] != data[i][j] {
				t.Errorf("Mismatch at row %d, column %d: expected %s, got %s",
					i, j, data[i][j], readData[i][j])
			}
		}
	}
}

// TestReadFile_FileNotFound ensures ReadFile returns an error for missing files.
func TestReadFile_FileNotFound(t *testing.T) {
	logger := zaptest.NewLogger(t)

	_, err := ReadFile("non_existent_file.csv", logger)
	if err == nil {
		t.Fatal("Expected an error for a missing file, but got nil")
	}
}

// TestWriteFile_FailToCreateFile ensures WriteFile returns an error when the file cannot be created.
func TestWriteFile_FailToCreateFile(t *testing.T) {
	logger := zaptest.NewLogger(t)

	// Attempt to write to an invalid file path
	err := WriteFile("/invalid_path/output.csv", [][]string{{"Data"}}, logger)
	if err == nil {
		t.Fatal("Expected an error for an invalid file path, but got nil")
	}
}
