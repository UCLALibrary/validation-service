//go:build integration

// Package integration holds the project's integration tests.
//
// This file contains tests of the service's `getStatus` endpoint.
package integration

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestStatusGet tests the status endpoint to confirm that the server is responding as expected to status queries.
func TestStatusGet(t *testing.T) {
	// Set up client for making requests to the containerized app
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Make a GET request to the containerized app and check the response
	response, err := client.Get(fmt.Sprintf(testServerURL, "/status"))
	if err != nil {
		t.Fatal(err)
	}
	//noinspection GoUnhandledErrorResult
	defer response.Body.Close()

	// Read the body into a variable to compare
	body, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		t.Fatalf("Error reading response: %v", readErr)
	}

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.JSONEq(t, "{\"fester\":\"ok\",\"filesystem\":\"ok\",\"service\":\"ok\"}", string(body))
}

// TestStatusPost tests the status endpoint to confirm that the server doesn't respond to POST submissions.
func TestStatusPost(t *testing.T) {
	// Set up client for making requests to the containerized app
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Make a POST request and confirm that POST is not supported
	requestBody := []byte(`{"key": "value"}`)

	response, err := client.Post(fmt.Sprintf(testServerURL, "/status"), "application/json",
		bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	//noinspection GoUnhandledErrorResult
	defer response.Body.Close()

	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}
