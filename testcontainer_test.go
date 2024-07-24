package main

import (
	"bytes"
	"context"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestApp(t *testing.T) {
	// Define the container request
	req := testcontainers.ContainerRequest{
		Image:        "service-template",
		ExposedPorts: []string{"8888/tcp"},
		SkipReaper:   true,
		WaitingFor:   wait.ForHTTP("/").WithPort("8888/tcp"),
	}

	// Create a context for the container
	ctx := context.Background()

	// Start the container
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer container.Terminate(ctx)

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

	// Make requests to the containerized app and assert the responses
	// Example: Make an HTTP request to the root endpoint
	resp, err := client.Get("http://" + host + ":" + strconv.Itoa(port.Int()) + "/")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Any requeset body will work since POST requests are not currently allowed
	requestBody := []byte(`{"key": "value"}`)

	resp, err = client.Post("http://" + host + ":" + strconv.Itoa(port.Int()) + "/", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
