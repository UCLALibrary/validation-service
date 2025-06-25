//go:build unit

package validation

import (
	"github.com/UCLALibrary/validation-service/validation/profiles"
	"os"
	"testing"

	"github.com/UCLALibrary/validation-service/pkg/utils"
	"github.com/UCLALibrary/validation-service/validation/csv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

// TestEngine_NewEngine tests the construction of a validation engine.
func TestEngine_NewEngine(t *testing.T) {
	// Set system env to tell the engine where to find its persisted Profiles file
	if err := os.Setenv(profiles.ConfigFile, "../testdata/test_profiles.json"); err != nil {
		t.Errorf("Error setting env PROFILES_FILE: %v", err)
	}
	defer func() {
		err := os.Unsetenv(profiles.ConfigFile)
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
	if err := os.Setenv(profiles.ConfigFile, "../testdata/test_profiles.json"); err != nil {
		t.Fatalf("error setting env PROFILES_FILE: %v", err)
	}
	defer func() {
		err := os.Unsetenv(profiles.ConfigFile)
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
	logger := zaptest.NewLogger(t, zaptest.Level(utils.GetLogLevel()))

	// Configure the location of the test profiles file
	if err := os.Setenv(profiles.ConfigFile, "../testdata/test_profiles.json"); err != nil {
		t.Fatalf("error setting env PROFILES_FILE: %v", err)
	}
	defer func() {
		err := os.Unsetenv(profiles.ConfigFile)
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
