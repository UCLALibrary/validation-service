//go:build unit

// Package validation provides tools to validate CSV data.
//
// This file provides tests of the validation Engine.
package validation

import (
	"flag"
	"fmt"
	"github.com/UCLALibrary/validation-service/testflags"
	"github.com/UCLALibrary/validation-service/validation/config"
	"github.com/UCLALibrary/validation-service/validation/csv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"os"
	"testing"
)

// TestValidation loads the flags for the tests in the 'validation' package.
func TestValidation(t *testing.T) {
	flag.Parse()
	fmt.Printf("%s's log level: %s\n", t.Name(), *testflags.LogLevel)
}

// TestEngine_NewEngine tests the construction of a validation engine.
func TestEngine_NewEngine(t *testing.T) {
	// Set system env to tell the engine where to find its persisted Profiles file
	if err := os.Setenv(config.ProfilesFile, "../testdata/test_profiles.json"); err != nil {
		t.Errorf("Error setting env PROFILES_FILE: %v", err)
	}
	defer func() {
		err := os.Unsetenv(config.ProfilesFile)
		require.NoError(t, err)
	}()

	// Create a new validation engine
	engine, err := NewEngine()
	if err != nil {
		t.Fatalf("error creating new validation engine: %s", err)
	}

	assert.NotNil(t, engine)
}

// TestEngine_GetLogger tests that a new engine has created a logger and can return it.
func TestEngine_GetLogger(t *testing.T) {
	// Set system env to tell the engine where to find its persisted Profiles file
	if err := os.Setenv(config.ProfilesFile, "../testdata/test_profiles.json"); err != nil {
		t.Fatalf("error setting env PROFILES_FILE: %v", err)
	}
	defer func() {
		err := os.Unsetenv(config.ProfilesFile)
		require.NoError(t, err)
	}()

	// Create a new validation engine
	engine, err := NewEngine()
	if err != nil {
		t.Fatalf("error creating new validation engine: %s", err)
	}

	// Confirm that not passing in a logger doesn't break using a logger
	assert.NotNil(t, engine.GetLogger())
}

// TestEngine_GetValidators tests that an engine can return the validators its using
func TestEngine_Validate(t *testing.T) {
	logger := zaptest.NewLogger(t, zaptest.Level(testflags.GetLogLevel()))

	// Configure the location of the test profiles file
	if err := os.Setenv(config.ProfilesFile, "../testdata/test_profiles.json"); err != nil {
		t.Fatalf("error setting env PROFILES_FILE: %v", err)
	}
	defer func() {
		err := os.Unsetenv(config.ProfilesFile)
		require.NoError(t, err)
	}()

	// Create a new validation engine using our test logger
	engine, engineErr := NewEngine(logger)
	if engineErr != nil {
		t.Fatalf("error creating new validation engine: %s", engineErr)
	}

	// Read in our CSV test data
	csvData, csvErr := csv.ReadFile("../testdata/cct-works-simple.csv", engine.GetLogger())
	if csvErr != nil {
		require.NoError(t, csvErr)
	}

	// Validate the CSV data we're read in from the test file
	err := engine.Validate("test", csvData)
	if err != nil {
		t.Fatalf("error getting validators: %s", err)
	}
}
