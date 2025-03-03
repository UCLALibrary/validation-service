//go:build integration

// Package integration holds the project's integration tests.
//
// This file contains tests of the service's `uploadCSV` endpoint.
package integration

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
)

// TestUploadCSV tests the uploadCSV endpoint with a basic CSV.
func TestUploadCSV(t *testing.T) {
	csvFile := "../testdata/cct-works-simple.csv"

	// Open the test CSV file
	file, openErr := os.Open(csvFile)
	if openErr != nil {
		t.Fatalf("Error opening file: %v", openErr)
	}
	//noinspection GoUnhandledErrorResult
	defer file.Close()

	// Create a buffer to store the multipart form data
	upload := &bytes.Buffer{}
	writer := multipart.NewWriter(upload)

	// Create a form file field named 'csvFile'
	part, formErr := writer.CreateFormFile("csvFile", csvFile)
	if formErr != nil {
		t.Fatalf("Error creating form file: %v", formErr)
	}

	// Copy the CSV file content into the multipart form
	if _, err := io.Copy(part, file); err != nil {
		t.Fatalf("Error copying file content: %v", err)
	}

	// Add the 'profile' field, too
	_ = writer.WriteField("profile", "test")

	// Close the multipart writer to finalize the form
	if err := writer.Close(); err != nil {
		t.Fatalf("Error closing writer: %v", err)
	}

	// Create an HTTP request to our CSV upload endpoint
	url := fmt.Sprintf(testServerURL, "/upload/csv")
	request, reqErr := http.NewRequest("POST", url, upload)
	if reqErr != nil {
		t.Fatalf("Error creating request: %v", reqErr)
	}

	// Set the Content-Type header to multipart/form-data with the correct boundary
	request.Header.Set("Content-Type", writer.FormDataContentType())

	// Initialize an HTTP client and send the request
	client := &http.Client{}
	response, postErr := client.Do(request)
	if postErr != nil {
		t.Fatalf("Error sending request: %v", postErr)
	}
	//noinspection GoUnhandledErrorResult
	defer response.Body.Close()

	// Read the body into a variable to compare
	body, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		t.Fatalf("Error reading response: %v", readErr)
	}

	// Placeholder body check
	assert.Equal(t, "{\"fs\":\"ok\",\"s3\":\"ok\",\"service\":\"created\"}\n", string(body))
	assert.Equal(t, http.StatusCreated, response.StatusCode)
}
