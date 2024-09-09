# Variables
DOCKER_IMAGE := validation-service
DOCKER_TAG := latest

# Build the Go project
run:
	go run main.go

# Run Go tests
test:
	go test ./... -v

# Lint the code
lint:
	golangci-lint run

# Build the Docker container
docker-build:
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Run tests inside the Docker container
docker-test:
	docker run --rm $(DOCKER_IMAGE):$(DOCKER_TAG) go test ./... -v

# Clean up
clean:
	rm -rf ./bin
