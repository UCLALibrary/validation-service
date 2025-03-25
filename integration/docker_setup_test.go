//go:build integration

// Package integration holds the project's integration tests.
//
// This file sets up the Docker container for integration testing.
package integration

import (
	"bytes"
	docker "context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/docker/docker/pkg/stdcopy"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
)

// Define our test container's build arguments
var serviceName string
var logLevel string
var hostDir string

// The URL to which to submit test HTTP requests
var testServerURL string

// A reference to our Docker container
var container testcontainers.Container

// Initialize our service name flag.
func init() {
	flag.StringVar(&serviceName, "service-name", "service", "Name of service being tested")
	flag.StringVar(&logLevel, "log-level", "info", "Log level (debug, info, warn, error)")
	flag.StringVar(&hostDir, "host-dir", "", "HOST_DIR env variable that is copied into test-container")
}

// TestMain spins up a Docker container with our validation service to run tests against.
func TestMain(m *testing.M) {
	flag.Parse()

	// Creates a logger for our tests
	logger, _ = getLogger()
	//noinspection GoUnhandledErrorResult
	defer logger.Sync()
	// Get the Docker context
	context := docker.Background()

	logger.Info("Checking if hostDir exists", zap.String("hostDir", hostDir))

	if _, err := os.Stat(hostDir); os.IsNotExist(err) {
		logger.Fatal("Host directory does not exist", zap.String("hostDir", hostDir))
	}

	logger.Info("HOST_DIR %s", zap.String("hostDir", hostDir))

	// Define the container request
	request := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    "..",
			Dockerfile: "Dockerfile",
			BuildArgs: map[string]*string{
				"SERVICE_NAME": &serviceName,
				"LOG_LEVEL":    &logLevel,
				"HOST_DIR":     &hostDir,
			},
		},
		ExposedPorts: []string{"8888/tcp"},
		WaitingFor:   wait.ForHTTP("/status").WithPort("8888/tcp"),
		LogConsumerCfg: &testcontainers.LogConsumerConfig{
			Consumers: []testcontainers.LogConsumer{&DockerLogConsumer{}},
		},
	}

	// Disable unnecessary output logging from the test containers
	testcontainers.Logger = &FilteredLogger{
		original: log.New(log.Writer(), "", log.LstdFlags),
	}

	// Start the container
	var containerErr error
	container, containerErr = testcontainers.GenericContainer(context, testcontainers.GenericContainerRequest{
		ContainerRequest: request,
		Started:          true,
	})
	if containerErr != nil {
		logger.Fatal("Failed to start Docker container", zap.Error(containerErr))
	}

	//noinspection GoUnhandledErrorResult
	defer container.Terminate(context)

	// Get the mapped host and port
	host, hostErr := container.Host(context)
	if hostErr != nil {
		logger.Fatal("Failed to get container host", zap.Error(hostErr))
	}

	port, portErr := container.MappedPort(context, "8888")
	if portErr != nil {
		logger.Fatal("Failed to get container port", zap.Error(portErr))
	}

	// Store the connect information for reuse in tests
	testServerURL = fmt.Sprintf("http://%s:%d", host, port.Int()) + "%s"

	// Run tests
	code := m.Run()

	// Cleanup: Stop the container after all tests
	exitErr := container.Terminate(context)
	if exitErr != nil {
		logger.Fatal("Failed to terminate Docker container", zap.Error(exitErr))
	}

	os.Exit(code)
}

// TestEnvironmentVariable checks if the HOST_DIR ENV is set and if it matches the expected value
func TestEnvironmentVariable(t *testing.T) {
	// Execute a command inside the container to check the env variable
	envVar := "HOST_DIR"
	expectedValue := hostDir

	context := docker.Background()
	_, reader, err := container.Exec(context, []string{"printenv", envVar})

	if err != nil {
		t.Fatalf("Failed to execute command inside container: %v", err)
	}

	// Separate stdout and stderr from the raw reader
	var stdout, stderr bytes.Buffer
	_, err = stdcopy.StdCopy(&stdout, &stderr, reader)
	if err != nil {
		t.Fatalf("Failed to read container output: %v", err)
	}

	assert.Equal(t, expectedValue, strings.TrimSpace(stdout.String()), "Environment variable value is incorrect")
}
