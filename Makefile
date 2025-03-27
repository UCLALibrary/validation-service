# Build and runtime variables that can be overridden
SERVICE_NAME ?= validation-service
LOG_LEVEL := info
PORT := 8888
VERSION ?= dev-SNAPSHOT
HOST_DIR := $(shell pwd)/testdata
PERSONAL_ACCESS_TOKEN ?= ""
ARCH ?= x86-64
KAKADU_PATH := v8_4_1-01903L/*
CREATE_KAKADU ?= ""
# Docker related variables
TAG ?= validation-service
DOCKER_REGISTRY ?= docker.io
DOCKER_REPO ?= validation-service
DOCKER_IMAGE ?= $(DOCKER_REGISTRY)/$(DOCKER_REPO):$(TAG)

# Force the API target to run even if the openapi.yml has not been touched/changed
ifneq ($(filter FORCE,$(MAKECMDGOALS)),)
.PHONY: api/api.go
endif

# Define FORCE as a target so accidentally using on other targets won't cause errors
.PHONY: FORCE help
FORCE:
	@echo "Makefile target(s) run with FORCE to require an API code rebuild"

all: config api lint build test docker-test # Does a full build of the project

lint: # Checks the code for correctness / coding standards
	golangci-lint run

# We generate Go API code from the OpenAPI specification only when it has changed
# We assume Windows developers are using WSL, so we don't define $(CP) for this
api/api.go: openapi.yml
	oapi-codegen -package api -generate types,server,spec -o api/api.go openapi.yml
	cp openapi.yml html/assets/openapi.yml

api: api/api.go # Generates API code from the OpenAPI specification

build: api # Compiles the project's Go code into an executable
	go build -o $(SERVICE_NAME)

test: # Runs the unit tests (integration tests are excluded)
	go test -tags=unit ./... -v -args -log-level=$(LOG_LEVEL) -host-dir=$(HOST_DIR)

docker-build:  # Builds a Docker container for manual testing
	@if [ ! -z "$(CREATE_KAKADU)" ]; then \
		echo "Running clone-kakadu..."; \
		$(MAKE) clone-kakadu || true; \
	else \
		echo "CREATE_KAKADU is not set, skipping clone-kakadu"; \
		mkdir -p kakadu; \
	fi
	@echo "Running Docker build..."
	docker build . --tag $(DOCKER_IMAGE) --build-arg SERVICE_NAME=$(SERVICE_NAME) \
		--build-arg VERSION=$(VERSION) --build-arg HOST_DIR=$(HOST_DIR) \
		--build-arg CREATE_KAKADU=$(CREATE_KAKADU) --build-arg ARCH=$(ARCH)

docker-run: docker-build # Runs a Docker instance, independent of the tests
	CONTAINER_ID=$(shell docker image ls -q --filter=reference=$(SERVICE_NAME)); \
	docker run -p $(PORT):8888 --name $(SERVICE_NAME) \
		-e LOG_LEVEL="$(LOG_LEVEL)" \
		-v $(HOST_DIR):$(HOST_DIR) \
		-d $$CONTAINER_ID

docker-log: # Tails the logs of a container started with 'docker-run'
	docker logs -f $(shell docker ps --filter "name=$(SERVICE_NAME)" --format "{{.ID}}")

docker-stop: # Stops a Docker container started with 'docker-run'
	docker rm -f $(shell docker ps --filter "name=$(SERVICE_NAME)" --format "{{.ID}}")

# 'docker-test' does not require 'docker-build', fwiw, 'docker-build' is just for debugging
docker-test: # Runs integration tests inside the Docker container
	@mkdir -p kakadu
	go test -tags=integration ./integration -v -args -service-name=$(SERVICE_NAME) -log-level=$(LOG_LEVEL) \
		-host-dir=$(HOST_DIR)

docker-push:  docker-build # builds and pushes a Docker image
	docker push $(DOCKER_IMAGE)

clean: # Cleans up all artifacts created by the build
	rm -rf $(SERVICE_NAME) api/api.go

# Creates a new local profile configuration file if it doesn't already exist
profiles.json: profiles.example.json
	@if [ ! -f profiles.json ]; then cp profiles.example.json profiles.json; fi

config: profiles.json # Creates a profiles.json file from the example file

run: config api build # Runs service locally, independent of Docker
	PROFILES_FILE="profiles.json" LOG_LEVEL=$(LOG_LEVEL) VERSION=$(VERSION) HOST_DIR=$(HOST_DIR) ./$(SERVICE_NAME)

ci-run: config api # Runs CI locally using ACT (which must be installed)
	pkg/scripts/act.sh $(JOB) $(SERVICE_NAME)

clone-kakadu:
	echo "Cloning Kakadu repository..."; \
	rm -rf kakadu && \
	mkdir kakadu && cd kakadu && \
	git init && \
	git remote add origin https://$(PERSONAL_ACCESS_TOKEN)@github.com/UCLALibrary/kakadu.git && \
	git config core.sparseCheckout true && \
	echo "$(KAKADU_PATH)" > .git/info/sparse-checkout && \
	git pull origin main; \

help: # Outputs information about the build's available targets
	@awk -F ':.*?# ' '/^[a-z0-9_-]+:.*?# / && $$1 !~ /[A-Z.]/ { \
		printf "\033[1;32m%-20s\033[0m %s\n", $$1, $$2 \
	}' Makefile
