# Variables
SERVICE_NAME := validation-service

# Do a full build of the project
all: lint build test docker-test

# Build the Go project
build:
	go build -o $(SERVICE_NAME)

# Run Go tests
test:
	go test -tags=unit ./... -v

# Lint the code
lint:
	golangci-lint run

# Run tests inside the Docker container
docker-test:
	go test -tags=functional ./... -v

# Clean up all artifacts of the build
clean:
	rm -rf $(SERVICE_NAME)
