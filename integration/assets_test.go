//go:build integration

package integration

import (
	docker "context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

// TestJavascriptResources confirms the Javascript resources can be resolved.
func TestJavascriptResources(t *testing.T) {
	// Set up client for making requests to the containerized app
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Define test cases
	tests := []struct {
		name     string
		resource string
	}{
		{"ReportJS", "/report.js"},
		{"ValidationJS", "/validation.js"},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			response, err := client.Get(fmt.Sprintf(testServerURL, testCase.resource))
			if err != nil {
				t.Fatal(err)
			}
			//noinspection GoUnhandledErrorResult
			defer response.Body.Close()

			assert.Equal(t, http.StatusOK, response.StatusCode)
		})
	}
}

// TestCssResources confirms the CSS resources can be resolved.
func TestCssResources(t *testing.T) {
	// Set up client for making requests to the containerized app
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Define test cases
	tests := []struct {
		name     string
		resource string
	}{
		{"ValidationCSS", "/validation.css"},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			response, err := client.Get(fmt.Sprintf(testServerURL, testCase.resource))
			if err != nil {
				t.Fatal(err)
			}
			//noinspection GoUnhandledErrorResult
			defer response.Body.Close()

			assert.Equal(t, http.StatusOK, response.StatusCode)
		})
	}
}

// TestOpenApiResources confirms the OpenAPI spec can be resolved.
func TestOpenApiResources(t *testing.T) {
	// Set up client for making requests to the containerized app
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Define test cases
	tests := []struct {
		name     string
		resource string
	}{
		{"OpenApiSpec", "/openapi.yml"},
	}

	// Confirm our OpenAPI spec has been copied into our container
	filePath := "/usr/local/data/html/assets/openapi.yml"
	exitCode, _, _ := container.Exec(docker.Background(), []string{"test", "-f", filePath})
	if exitCode != 0 {
		t.Fatalf("File '%s' doesn't exist in the container\n", filePath)
	}

	// Check that the OpenAPI spec can be downloaded from the server
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			response, err := client.Get(fmt.Sprintf(testServerURL, testCase.resource))
			if err != nil {
				t.Fatal(err)
			}
			//noinspection GoUnhandledErrorResult
			defer response.Body.Close()

			assert.Equal(t, http.StatusOK, response.StatusCode)
		})
	}
}
