//go:build unit

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"strings"
	"testing"
)

// TestProfiles tests basic Get/Set functionality.
func TestProfiles(t *testing.T) {
	profiles := NewProfiles()
	defaultProfile, err := NewProfile("default", []string{"SpaceCheck", "ARKFormat", "EOLCheck"})
	require.NoError(t, err)
	otherProfile, err := NewProfile("other", []string{"SpaceCheck", "ARKFormat"})
	require.NoError(t, err)

	err = profiles.SetProfile(defaultProfile)
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

// TestSnapshot tests creating a bare-bones Snapshot through marshaling it to JSON.
func TestSnapshot(t *testing.T) {
	profiles := NewProfiles()
	snapshot := profiles.Snapshot()

	jsonData, err := json.Marshal(snapshot)
	if err != nil {
		log.Fatalf("Error marshaling to JSON: %v", err)
	}

	// Confirm Snapshot is okay with no profiles set (i.e., we've created the slice)
	assert.Equal(t, "{\"profiles\":{},\"lastUpdate\":\"0001-01-01T00:00:00Z\"}", string(jsonData))
}

// TestExampleCode tests the code that's used in profiles.go's inline docs.
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
