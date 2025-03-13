//go:build integration

package integration

import (
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
