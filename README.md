# Validation Service

This is our microservice for CSV validation. It validations CSVs according to a prescribed set of rules. These rule
sets may be incorporated into specific validation profiles.

Note: This project is in its infancy and is not ready for general use.

## Getting Started

There are multiple ways to build the project. You are free to use whichever you prefer. A series of manual steps is
provided below, but there are also more concise build processes available for [Make](#using-the-makefile) and [ACT]
(#building-and-running-with-act).

### Building the Application

To build the project run:

`go build -o validation-service`

### Creating a Docker Image

To run on Docker first build the Docker image:

`docker build -t validation-service .`

To specify what version of Go you would like to use with the Docker image:

`docker build --build-arg GO_VERSION=[YOUR_VERSION] -t validation-service .`

To run the Docker image:

`docker run -d -p 8888:8888 validation-service`

Once the container is running, you can access the service at:

`http://localhost:8888`

### Running Using Docker Compose

To run the validator using Docker Compose, use the following command:

`docker-compose up --build`

Once the container is running, you can access the service at:

`http://localhost:8888`

To stop the running containers, use the following command:

`docker-compose down`

## Building and Running with ACT

Using [ACT](https://github.com/nektos/act) is another way to build the project.

To get started, ensure that [ACT is installed](https://nektosact.com/installation/index.html) on your system.

Now that ACT is installed, you can run the build by typing the below:

`act -W .github/workflows/build.yml`

If you've installed ACT as an extension to the GitHub CLI, you'd type:

`gh act -W .github/workflows/build.yml`

## Building and Running with Make

The project's [Makefile](Makefile) provides another convenient way to build, test, lint, and manage Docker containers
for the project. Running `make all` will perform all the required build steps, but individual build steps (listed
below) are also available. It's also worth noting that there are the two optional build variables for the Makefile.

### Variables

`DOCKER_IMAGE`: The name of the Docker image (default: validation-service).
`DOCKER_TAG`: The tag for the Docker image (default: latest).

### Commands

To build and run the Go project locally:

    make run

To run all Go tests with verbose output:

    make test

To run the linter using `golangci-lint` to check the code for style and correctness:

    make lint

To build a Docker image using the specified `DOCKER_IMAGE` and `DOCKER_TAG`:

    make docker-build

To run Go tests inside a Docker container built from the project:

    make docker-test

To clean up the project's build artifacts, run:

    make clean

## Contact

If you have any questions or suggestions, feel free to [open a ticket](https://github.com/UCLALibrary/validation-service/issues) on project's GitHub repo.
