# Build and runtime variables that can be overridden
SERVICE_NAME := validation-service
LOG_LEVEL := info
PORT := 8888
VERSION := dev-SNAPSHOT

# Force the API target to run even if the openapi.yml has not been touched/changed
ifneq ($(filter FORCE,$(MAKECMDGOALS)),)
.PHONY: api/api.go
endif

# Define FORCE as a target so accidentally using on other targets won't cause errors
.PHONY: FORCE
FORCE:
	@echo "Makefile target(s) run with FORCE to require an API code rebuild"

# Do a full build of the project
all: config api lint build test docker-test

# Lint the code for correctness
lint:
	golangci-lint run

# We generate Go API code from the OpenAPI specification only when it has changed
# We assume Windows developers are using WSL, so we don't define $(CP) for this
api/api.go: openapi.yml
	oapi-codegen -package api -generate types,server,spec -o api/api.go openapi.yml
	cp openapi.yml html/assets/openapi.yml

# This is an alias for the longer API generation Makefile target api/api.go
api: api/api.go

# Build the Go project
build: api
	go build -o $(SERVICE_NAME)

# Run Go tests, excluding tests in the 'integration' directory
test:
	go test -tags=unit ./... -v -args -log-level=$(LOG_LEVEL)

# Build the Docker container (an optional debugging step)
docker-build:
	docker build . --tag $(SERVICE_NAME) --build-arg SERVICE_NAME=$(SERVICE_NAME) --build-arg VERSION=$(VERSION)

# A convenience target to assist with running the Docker container outside of the build (optional)
docker-run:
	CONTAINER_ID=$(shell docker image ls -q --filter=reference=$(SERVICE_NAME)); \
	docker run -p $(PORT):8888 --name $(SERVICE_NAME) -e LOG_LEVEL="$(LOG_LEVEL)" -d $$CONTAINER_ID

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

# Creates a new local profile configuration file if it doesn't already exist
profiles.json: profiles.example.json
	@if [ ! -f profiles.json ]; then cp profiles.example.json profiles.json; fi

# An alias for the profile.json target
config: profiles.json

# Run the validation service locally, independent of the Docker container
run: config api build
	PROFILES_FILE="profiles.json" LOG_LEVEL=$(LOG_LEVEL) VERSION=$(VERSION) ./$(SERVICE_NAME)

# Run a CI action (assuming the CI prerequisites from the README have also been installed)
ci-run:
	pkg/scripts/act.sh $(JOB) $(SERVICE_NAME)
