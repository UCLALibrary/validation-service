//go:build unit
package util

// This file contains tests of the validation profiles.
import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestProfiles tests basic Get/Set functionality.
func TestProfiles(t *testing.T) {
	profiles := NewProfiles()
	defaultProfile, errP1 := NewProfile("default", []string{"SpaceCheck", "ARKFormat", "EOLCheck"})
	require.NoError(t, errP1)
	otherProfile, errP2 := NewProfile("other", []string{"SpaceCheck", "ARKFormat"})
	require.NoError(t, errP2)

	err := profiles.SetProfile(defaultProfile)
	require.NoError(t, err)
	err = profiles.SetProfile(otherProfile)
	require.NoError(t, err)

	assert.Equal(t, 2, profiles.Count())
	assert.Equal(t, profiles.GetProfile("default").GetName(), "default")
	assert.Equal(t, profiles.GetProfile("other").GetName(), "other")
	assert.Equal(t, len(profiles.GetProfile("default").GetValidations()), 3)
	assert.Equal(t, len(profiles.GetProfile("other").GetValidations()), 2)
	assert.False(t, profiles.GetProfile("default").GetLastUpdate().IsZero())
	assert.False(t, profiles.GetProfile("other").GetLastUpdate().IsZero())
}

// TestProfiles_Snapshot tests creating a bare-bones Snapshot through marshaling it to JSON.
func TestProfiles_Snapshot(t *testing.T) {
	profiles := NewProfiles()
	snapshot := profiles.Snapshot()

	jsonData, err := json.Marshal(snapshot)
	if err != nil {
		log.Fatalf("Error marshaling to JSON: %v", err)
	}

	// Confirm Snapshot is okay with no profiles set (i.e., we've created the slice)
	assert.Equal(t, "{\"profiles\":{},\"lastUpdate\":\"0001-01-01T00:00:00Z\"}", string(jsonData))
}

// TestProfiles_Refresh tests refreshing a Profiles instance from a persisted JSON file.
func TestProfiles_Refresh(t *testing.T) {
	// Set the PROFILES_FILE for testing purposes
	err := os.Setenv(ProfilesFile, "../../testdata/test_profiles.json")
	require.NoError(t, err)
	defer func() {
		err := os.Unsetenv(ProfilesFile)
		require.NoError(t, err)
	}()

	profiles := NewProfiles()

	err = profiles.Refresh()
	require.NoError(t, err)

	example := profiles.GetProfile("example")
	test := profiles.GetProfile("test")

	assert.Equal(t, example.GetName(), "example")
	assert.Equal(t, 2, len(example.GetValidations()))
	assert.Equal(t, test.GetName(), "test")
	assert.Equal(t, 1, len(test.GetValidations()))
}

// TestProfiles_Save tests saving a Profiles instance to a JSON file for persistence.
func TestProfiles_Save(t *testing.T) {
	var snapshot ProfilesSnapshot

	// Set up a file for testing JSON persistence
	tempFile, err := os.CreateTemp("", "profiles-*.json")
	require.NoError(t, err)
	defer func(name string) {
		err := os.Remove(name)
		require.NoError(t, err)
	}(tempFile.Name())

	// Set the PROFILES_FILE for testing purposes
	err = os.Setenv(ProfilesFile, tempFile.Name())
	require.NoError(t, err)
	defer func() {
		err := os.Unsetenv(ProfilesFile)
		require.NoError(t, err)
	}()

	// Create a new Profile for testing
	profiles := NewProfiles()
	profile, _ := NewProfile("example", []string{"Validation1", "Validation2"})
	_ = profiles.SetProfile(profile)

	// Save profiles to a JSON file
	err = profiles.Save()
	require.NoError(t, err)

	// Read the file back
	data, readErr := os.ReadFile(tempFile.Name())
	require.NoError(t, readErr)

	// Verify the JSON data
	err = json.Unmarshal(data, &snapshot)
	require.NoError(t, err)
	assert.Equal(t, 1, len(snapshot.Profile))
	assert.Equal(t, "example", snapshot.Profile["example"].Name)
}

// TestExampleCode tests the code that's used in profiles.go's inline docs.
//
// This presents a sync challenge, but ensures the example code actually runs!
func TestExampleCode(t *testing.T) {
	// The function that's captured contains the example code used in the docs
	output := captureOutput(t, func() {
		profiles := NewProfiles()
		if profile, err := NewProfile("example", []string{"Validation1", "Validation2"}); err == nil {
			err = profiles.SetProfile(profile)
			require.NoError(t, err)

			fmt.Println(profiles.GetProfile("example").GetName())
		} else {
			require.NoError(t, err)
		}

		snapshot := profiles.Snapshot()
		if jsonData, err := json.Marshal(snapshot); err == nil {
			fmt.Println(string(jsonData))
		} else {
			require.NoError(t, err)
		}
	})

	// A simple test that gets around lastUpdate being different
	assert.True(t, strings.HasPrefix(output, "example"))
}

// captureOutput captures StdOut and allows running assertions on it
func captureOutput(t *testing.T, f func()) string {
	var buf bytes.Buffer

	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the function we've passed in
	f()

	if err := w.Close(); err != nil {
		t.Fatalf("Failed to close pipe writer: %v", err)
	}

	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("Failed to read from pipe reader: %v", err)
	}

	os.Stdout = originalStdout
	return buf.String()
}
