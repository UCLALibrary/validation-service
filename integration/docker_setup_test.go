//go:build integration

// Package integration holds the project's integration tests.
//
// This file sets up the Docker container for integration testing.
package integration

import (
	docker "context"
	"flag"
	"fmt"
	"log"
	"os"
	"testing"

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

	// Copy testdata folder into testcontainer
	hostDir = os.Getenv("HOST_DIR")
	if hostDir == "" {
		logger.Fatal("HOST_DIR is not set")
	}

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
	container, containerErr := testcontainers.GenericContainer(context, testcontainers.GenericContainerRequest{
		ContainerRequest: request,
		Started:          true,
	})
	if containerErr != nil {
		logger.Fatal("Failed to start Docker container", zap.Error(containerErr))
	}

	//noinspection GoUnhandledErrorResult
	defer container.Terminate(context)

	copyHost := hostDir
	_, _, err := container.Exec(docker.Background(), []string{"mkdir", "-p", copyHost})

	if err != nil {
		logger.Fatal("Failed to create HOST_DIR folder", zap.Error(err))
	}

	err = container.CopyDirToContainer(context, hostDir, copyHost, 0o700)

	if err != nil {
		logger.Fatal("Failed to copy to the container", zap.String("hostDir", hostDir), zap.Error(err))
	}

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

// func TestEnvironmentVariable(t *testing.T) {
// 	// Execute a command inside the container to check the env variable
// 	envVar := "HOST_DIR"
// 	expectedValue := hostDir

// 	context := docker.Background()
// 	_, reader, err := container.Exec(context, []string{"sh", "-c", fmt.Sprintf("echo $%s", envVar)})
// 	if err != nil {
// 		t.Fatalf("Failed to execute command inside container: %v", err)
// 	}

// 	output, err := io.ReadAll(reader)

// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// Read output and trim whitespace
// 	dir := strings.TrimSpace(string(output))

// 	assert.Equal(t, expectedValue, dir, "Environment variable value is incorrect")
// }
