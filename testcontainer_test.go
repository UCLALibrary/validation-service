//go:build functional

package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Configure our service name flag
var serviceNameFlag = flag.String("service-name", "service", "A build arg for Dockerfile")

// TestApp spins up a Docker container with the application and runs simple tests against its Web API
func TestApp(t *testing.T) {
	flag.Parse()

	// Define a Docker context
	ctx := context.Background()

	// Define the container request
	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    ".",
			Dockerfile: "Dockerfile",
			BuildArgs: map[string]*string{
				"SERVICE_NAME": serviceNameFlag,
			},
		},
		ExposedPorts: []string{"8888/tcp"},
		WaitingFor:   wait.ForHTTP("/").WithPort("8888/tcp"),
	}

	// Start the container
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := container.Terminate(ctx); err != nil {
			fmt.Printf("Error terminating container: %v\n", err)
		}
	}()

	// Get the host and port for the running container
	host, err := container.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}

	port, err := container.MappedPort(ctx, "8888/tcp")
	if err != nil {
		t.Fatal(err)
	}

	// Set up client for making requests to the containerized app
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Make a GET request to the containerized app and check the response
	resp, err := client.Get(fmt.Sprintf("http://%s:%d/", host, port.Int()))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Make a POST request and confirm that POST is not currently supported
	requestBody := []byte(`{"key": "value"}`)

	resp, err = client.Post(fmt.Sprintf("http://%s:%d/", host, port.Int()), "application/json",
		bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
