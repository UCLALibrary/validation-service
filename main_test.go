//go:build unit

package main

import (
	"flag"
	"fmt"
	"github.com/UCLALibrary/validation-service/testflags"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMain configures our log level flag for the main package.
func TestMain(m *testing.M) {
	flag.Parse()
	fmt.Println("TestMain's log level:", *testflags.LogLevel)
	os.Exit(m.Run())
}

// TestHelloWorld is a very simple initial test for the validation service application.
func TestHelloWorld(t *testing.T) {
	app := NewApp()
	app.Routes()

	// Create a response recorder to record the response
	rec := httptest.NewRecorder()

	// Create a GET test request
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err, "Failed to create request")

	// Call the handler directly
	app.Router.ServeHTTP(rec, req)

	// Check the response status code and body
	assert.Equal(t, http.StatusOK, rec.Code, "Expected status code 200")
	assert.Equal(t, "hello world", rec.Body.String(), "Unexpected response body")

	// Create a POST test request
	req, err = http.NewRequest("POST", "/", nil)

	// A new recorder must be created
	rec = httptest.NewRecorder()
	assert.NoError(t, err, "Failed to create request")
	app.Router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code, "Expected status code 404")
}
