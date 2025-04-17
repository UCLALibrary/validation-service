//go:build integration

// Package integration holds the project's integration tests.
//
// This file tests that Kakadu has been installed and is functioning as expected.
package integration

import (
	"bytes"
	docker "context"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

// TestKakaduInstallation checks if Kakadu is installed in the container and functioning as expected.
func TestKakaduInstallation(t *testing.T) {
	if os.Getenv("KAKADU_VERSION") == "" {
		t.Skip("Skipping Kakadu integration test: KAKADU_VERSION is not set")
	}

	context := docker.Background()
	_, reader, err := container.Exec(context, []string{"kdu_compress", "-v"})
	if err != nil {
		t.Fatalf("Failed to execute command inside container: %v", err)
	}

	// Separate stdout and stderr from the raw reader
	var stdout, stderr bytes.Buffer
	_, err = stdcopy.StdCopy(&stdout, &stderr, reader)
	if err != nil {
		t.Fatalf("Failed to read container output: %v", err)
	}

	value := strings.TrimSpace(stdout.String())
	assert.True(t, strings.HasPrefix(value, "This is Kakadu"), "Expected output: 'This is Kakadu'")
}
