// Package utils provides utilities for the application.
//
// This file provides a reader and writer for CSV data.
package utils

import (
	"encoding/csv"
	"fmt"
	"go.uber.org/zap"
	"mime/multipart"
	"os"
)

// ReadUpload reads the CSV file from the supplied FileHeader and returns a string matrix.
func ReadUpload(fileHeader *multipart.FileHeader, logger *zap.Logger) ([][]string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		logger.Error("Failed to open uploaded file", zap.Error(err))
	}
	defer func() {
		if err := file.Close(); err != nil {
			logger.Error("failed to close file", zap.Error(err))
		}
	}()

	// Create a new CSV reader from the opened CSV file
	reader := csv.NewReader(file)

	// Read all records from the CSV reader
	if csvData, err := reader.ReadAll(); err != nil || len(csvData) < 1 {
		return nil, fmt.Errorf("failed to parse file '%s': %w", fileHeader.Filename, err)
	} else {
		return csvData, nil
	}
}

// ReadFile reads the CSV file at the supplied file path and returns a string matrix.
func ReadFile(filePath string, logger *zap.Logger) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			logger.Error("failed to close file", zap.Error(err))
		}
	}()

	// Create a new CSV reader from the opened CSV file
	reader := csv.NewReader(file)

	// Read all records from the CSV reader
	if csvData, err := reader.ReadAll(); err != nil || len(csvData) < 1 {
		return nil, fmt.Errorf("failed to parse file '%s': %w", filePath, err)
	} else {
		return csvData, nil
	}
}

// WriteFile writes a supplied string matrix to a CSV file.
func WriteFile(filePath string, data [][]string, logger *zap.Logger) error {
	// Open the file for writing, create it if it doesn't exist
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			logger.Error("Failed to close file", zap.String("filePath", filePath), zap.Error(err))
		}
	}()

	writer := csv.NewWriter(file)

	// Write all the CSV data to the file at once
	if err = writer.WriteAll(data); err != nil {
		return fmt.Errorf("failed to write file '%s': %w", filePath, err)
	}

	writer.Flush()

	// Check for any errors during the flushing process
	if err := writer.Error(); err != nil {
		return fmt.Errorf("error flushing data to CSV file '%s': %w", filePath, err)
	}

	return nil
}
