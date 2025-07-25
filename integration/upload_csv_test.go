//go:build integration

package integration

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUploadCSV tests the uploadCSV endpoint with a CSV file.
func TestUploadCSV(t *testing.T) {
	tests := []struct {
		name           string
		csvFilePath    string
		expectedStatus int
		expectedRegex  string
	}{
		{
			name:           "Valid CSV upload",
			csvFilePath:    "../testdata/cct-works-simple.csv",
			expectedStatus: http.StatusCreated,
			// expectedRegex handles JSON with or without line feeds and indentation
			expectedRegex: `\{\s*"profile"\s*:\s*"DLP Staff"\s*,\s*"time"\s*:\s*".*?"\s*,\s*"warnings"\s*:\s*\[\s*\]\s*\}`,
		},
		{
			name:           "Upload failure CSV",
			csvFilePath:    "../testdata/upload-failures.csv",
			expectedStatus: http.StatusCreated,
			// expectedRegex handles JSON with or without line feeds and indentation
			expectedRegex: `\{\s*"profile"\s*:\s*"DLP Staff"\s*,\s*"time"\s*:\s*".*?"\s*,\s*"warnings"\s*:\s*\[\s*\{\s*[\s\S]*?\s*\}\s*\]\s*\}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Open the CSV file
			file, openErr := os.Open(tt.csvFilePath)
			if openErr != nil {
				t.Fatalf("Error opening file %s: %v", tt.csvFilePath, openErr)
			}
			//noinspection GoUnhandledErrorResult
			defer file.Close()

			// Create a buffer to store the multipart form data
			upload := &bytes.Buffer{}
			writer := multipart.NewWriter(upload)

			// Create a form file field named 'csvFile'
			part, formErr := writer.CreateFormFile("csvFile", tt.csvFilePath)
			if formErr != nil {
				t.Fatalf("Error creating form file: %v", formErr)
			}

			// Copy the CSV file content into the multipart form
			if _, err := io.Copy(part, file); err != nil {
				t.Fatalf("Error copying file content: %v", err)
			}

			// Add the 'profile' field
			_ = writer.WriteField("profile", "DLP Staff")

			// Close the multipart writer
			if err := writer.Close(); err != nil {
				t.Fatalf("Error closing writer: %v", err)
			}

			// Create an HTTP request to our CSV upload endpoint
			url := fmt.Sprintf(testServerURL, "/upload/csv")
			request, reqErr := http.NewRequest("POST", url, upload)
			if reqErr != nil {
				t.Fatalf("Error creating request: %v", reqErr)
			}

			// Set the Content-Type header
			request.Header.Set("Content-Type", writer.FormDataContentType())

			// Initialize an HTTP client and send the request
			client := &http.Client{}
			response, postErr := client.Do(request)
			if postErr != nil {
				t.Fatalf("Error sending request: %v", postErr)
			}
			//noinspection GoUnhandledErrorResult
			defer response.Body.Close()

			// Read the response body
			body, readErr := io.ReadAll(response.Body)
			if readErr != nil {
				t.Fatalf("Error reading response: %v", readErr)
			}

			// Check the response against the expected regex pattern
			matched, _ := regexp.MatchString(tt.expectedRegex, string(body))
			if !matched {
				t.Errorf(
					"\n--- Regex match failed for test case: %s ---\nExpected pattern:\n%s\nActual body:\n%q\n",
					tt.name,
					tt.expectedRegex,
					body,
				)
			}

			// Check the expected status code
			assert.Equal(t, tt.expectedStatus, response.StatusCode, "Unexpected status code for test case: %s", tt.name)
		})
	}
}
