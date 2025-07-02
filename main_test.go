//go:build unit

package main

import (
	"flag"
	"fmt"
	"github.com/UCLALibrary/validation-service/api"
	"github.com/UCLALibrary/validation-service/pkg/utils"
	"github.com/UCLALibrary/validation-service/validation"
	"github.com/UCLALibrary/validation-service/validation/config"
	"github.com/UCLALibrary/validation-service/validation/util"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMain configures our log level flag for the main package.
//
// Usually this would go in a setup_test.go file, but there is only one test file in this package.
func TestMain(main *testing.M) {
	flag.Parse()
	fmt.Printf("*** Package %s's log level: %s ***\n", utils.GetPackageName(), utils.LogLevel)
	os.Exit(main.Run())
}

// TestServerHealth checks if the Echo server initializes properly
func TestServerHealth(t *testing.T) {
	// Configure the location of the test profiles file
	if err := os.Setenv(config.ConfigFile, "testdata/test_profiles.json"); err != nil {
		t.Fatalf("error setting env PROFILES_FILE: %v", err)
	}
	defer func() {
		err := os.Unsetenv(config.ConfigFile)
		require.NoError(t, err)
	}()

	engine, err := validation.NewEngine()
	assert.NoError(t, err)

	service := &Service{Engine: engine}
	server := echo.New()
	server.Use(util.ZapLoggerMiddleware(engine.GetLogger()))

	// Register handlers
	api.RegisterHandlers(server, service)

	// Perform a simple request to check if the server is functional
	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	rec := httptest.NewRecorder()

	// Serve the request
	server.ServeHTTP(rec, req)

	// Server should respond with a 200 status
	assert.Equal(t, http.StatusOK, rec.Code)
}

// TestStatusEndpoint checks if the /status endpoint returns the expected JSON response
func TestStatusEndpoint(t *testing.T) {
	// Configure the location of the test profiles file
	if err := os.Setenv(config.ConfigFile, "testdata/test_profiles.json"); err != nil {
		t.Fatalf("error setting env PROFILES_FILE: %v", err)
	}
	defer func() {
		err := os.Unsetenv(config.ConfigFile)
		require.NoError(t, err)
	}()

	engine, err := validation.NewEngine()
	assert.NoError(t, err)

	service := &Service{Engine: engine}
	server := echo.New()
	server.Use(util.ZapLoggerMiddleware(engine.GetLogger()))

	// Register handlers
	api.RegisterHandlers(server, service)

	// Create a test request
	request := httptest.NewRequest(http.MethodGet, "/status", nil)
	recorder := httptest.NewRecorder()

	// Execute request
	server.ServeHTTP(recorder, request)

	// Assertions
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.JSONEq(t, `{"fester":"ok", "filesystem":"ok", "service":"ok"}`, recorder.Body.String())
}
