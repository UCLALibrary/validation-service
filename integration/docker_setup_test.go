//go:build integration

package integration

import (
	"bytes"
	docker "context"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/UCLALibrary/validation-service/pkg/utils"
	"github.com/testcontainers/testcontainers-go/log"

	"github.com/docker/docker/pkg/stdcopy"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
)

// The URL to which to submit test HTTP requests
var testServerURL string

// A reference to our Docker container
var container testcontainers.Container

var hostDir = "/usr/local/data/testdata"

// TestMain spins up a Docker container with our validation service to run tests against.
func TestMain(m *testing.M) {
	flag.Parse()
	fmt.Printf("*** Package %s's log level: %s ***\n", utils.GetPackageName(), utils.LogLevel)

	// Creates a custom *zap.logger for our tests to use and then wrap it in a testcontainers logger
	logger, _ = getLogger(utils.LogLevel)
	tcLogger := NewTcLogger(logger)

	// Set the default TestContainers logger to use our wrapped *zap.Logger
	log.SetDefault(tcLogger)

	// Get the Docker context
	context := docker.Background()

	if _, err := os.Stat(utils.HostDir); os.IsNotExist(err) {
		logger.Fatal("HOST_DIR ENV property does not exist", zap.String("hostDir", utils.HostDir))
	} else {
		logger.Debug("HOST_DIR ENV property exists", zap.String("hostDir", utils.HostDir))
	}

	// Set up our test environment BuildArgs and ENV
	env := map[string]string{}
	buildArgs := map[string]*string{
		"SERVICE_NAME": &utils.ServiceName,
		"LOG_LEVEL":    &utils.LogLevel,
		"HOST_DIR":     &hostDir,
		"ARCH":         &utils.BuildArch,
		"MAX_UPLOAD":   func() *string { s := "5M"; return &s }(),
	}

	// Check to see if we're running our tests using the KAKADU_VERSION property and include that if so
	if utils.KakaduVersion != "" {
		env["KAKADU_VERSION"] = utils.KakaduVersion
		buildArgs["KAKADU_VERSION"] = &utils.KakaduVersion
	}

	// Define the container request
	request := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    "..",
			Dockerfile: "Dockerfile",
			BuildArgs:  buildArgs,
		},
		ExposedPorts: []string{"8888/tcp"},
		WaitingFor:   wait.ForHTTP("/status").WithPort("8888/tcp"),
		LogConsumerCfg: &testcontainers.LogConsumerConfig{
			// Logs from the container itself (i.e., not from TestContainers) come through this configuration
			Consumers: []testcontainers.LogConsumer{&DockerLogConsumer{}},
		},
		Env: env,
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      utils.HostDir,
				ContainerFilePath: hostDir,
			},
		},
	}

	// Start the container
	var containerErr error
	container, containerErr = testcontainers.GenericContainer(context, testcontainers.GenericContainerRequest{
		ContainerRequest: request,
		Started:          true,
		Logger:           tcLogger, // Logger for testcontainers events (e.g., container startup, shutdown, etc.)
	})
	if containerErr != nil {
		logger.Fatal("Failed to start Docker container", zap.Error(containerErr))
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

	// Store the service's URL location for reuse in tests
	testServerURL = fmt.Sprintf("http://%s:%d", host, port.Int()) + "%s"

	// Run tests
	code := m.Run()

	// Cleanup: Stop the container after all tests
	exitErr := container.Terminate(context)
	if exitErr != nil {
		logger.Fatal("Failed to terminate Docker container", zap.Error(exitErr))
	}

	// Sync without deferring right before the exit to make sure logs are written
	_ = logger.Sync()

	// Wrap up after all the container tests are done
	os.Exit(code)
}

// TestEnvironmentVariable checks if the HOST_DIR ENV property is set and if it matches the expected value.
func TestEnvironmentVariable(t *testing.T) {
	// Execute a command inside the container to check the env variable
	context := docker.Background()
	_, reader, err := container.Exec(context, []string{"printenv", "HOST_DIR"})

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
	assert.Equalf(
		t, hostDir, value, "The expected HOST_DIR ENV property wasn't found. Expected: %q, Found: %q",
		utils.HostDir, value,
	)
}

// TestMaxSizeEnvVariable checks if the MAX_UPLOAD ENV property is set and if it matches the expected value.
func TestMaxSizeEnvVariable(t *testing.T) {
	// Execute a command inside the container to check the MAX_UPLOAD env variable
	context := docker.Background()
	_, reader, err := container.Exec(context, []string{"printenv", "MAX_UPLOAD"})

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
	assert.Equalf(
		t, "5M", value, "The expected MAX_UPLOAD ENV property wasn't found. Expected: %q, Found: %q",
		"5M", value,
	)
}
