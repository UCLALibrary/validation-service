# Build and runtime variables
SERVICE_NAME := validation-service
LOG_LEVEL := info
PORT := 8888

# Do a full build of the project
all: api lint build test docker-test

# Lint the code
lint:
	golangci-lint run

# Generate Go code from the OpenAPI specification only when it has changed
api/api.go: openapi.yml
	oapi-codegen -package api -generate types,server,spec -o api/api.go openapi.yml

# This is an alias for the longer API generation Makefile target api/api.go
api: api/api.go

# Build the Go project
build:
	go build -o $(SERVICE_NAME)

# Run Go tests, excluding tests in the 'integration' directory
test:
	go test -tags=unit ./... -v -args -log-level=$(LOG_LEVEL)

# Build the Docker container (an optional debugging step)
docker-build:
	docker build . --tag $(SERVICE_NAME) --build-arg SERVICE_NAME=$(SERVICE_NAME)

# A convenience target to assist with running the Docker container outside of the build (optional)
docker-run:
	docker run -p $(PORT):8888 --name $(SERVICE_NAME) -d $(shell docker image ls -q --filter=reference=$(SERVICE_NAME))

docker-logs:
	docker logs -f $(shell docker ps --filter "name=$(SERVICE_NAME)" --format "{{.ID}}")

# A convenience target to assist with stopping the Docker container outside of the build (optional)
docker-stop:
	docker rm -f $(shell docker ps --filter "name=$(SERVICE_NAME)" --format "{{.ID}}")

# Run tests inside the Docker container (does not require docker-build, that's just for debugging)
docker-test:
	go test -tags=integration ./integration -v -args -service-name=$(SERVICE_NAME) -log-level=$(LOG_LEVEL)

# Clean up all artifacts of the build
clean:
	rm -rf $(SERVICE_NAME) api/api.go

# Run the validation service locally
run: api build
	PROFILES_FILE="testdata/test_profiles.json" ./$(SERVICE_NAME)
