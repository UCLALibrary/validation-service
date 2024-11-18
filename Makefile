# Variables
DOCKER_IMAGE := validation-service
DOCKER_TAG := latest

# Do a full build of the project
all: lint build test docker-build docker-test

# Build the Go project
build:
	go build -o $(DOCKER_IMAGE)

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
docker-test: docker-build
	docker run -d --name $(DOCKER_IMAGE)-test -p 8888:8888 $(DOCKER_IMAGE):$(DOCKER_TAG)
	go test ./... -v
	docker stop $(DOCKER_IMAGE)-test
	docker rm $(DOCKER_IMAGE)-test

# Clean up all artifacts of the build
clean:
	rm -rf $(DOCKER_IMAGE)
	docker rm -f $(DOCKER_IMAGE)-test
	docker rmi -f $(DOCKER_IMAGE)
