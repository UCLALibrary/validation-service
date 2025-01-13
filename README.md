# Validation Service

This is a microservice for CSV validation. It validates CSVs according to a prescribed set of rules. These rule
sets may be incorporated into specific validation profiles.

Note: This project is in its infancy and is not ready for general use.

## Getting Started

There are multiple ways to build the project. You are free to use whichever you prefer. A series of manual steps is
provided below, but there are also more concise build processes available for [Make](#building-and-running-with-make)
and [ACT](#building-and-running-with-act).

### Prerequisites

* A [GoLang](https://go.dev/doc/install) build environment
* A functional [Docker](https://docs.docker.com/get-started/get-docker/) installation
* The [golangci-lint](https://github.com/golangci/golangci-lint) linter for checking code style conformance

Optionally, if you want to test (or build using) the project's GitHub Actions:

* [ACT](https://github.com/nektos/act): A local GitHub Action runner that will also build and test the project

Additionally, Make can be used as a simpler build tool for the project. It should be installed through your OS' 
standard packaging system.

### Building the Application

To build the project, type:

`go build -o validation-service`

To run the service locally, type:

`./validation-service`

Typing `Ctrl-C` will stop the service.

### Running the Test Suite

There are unit and functional tests (the latter of which require a working Docker installation).

To run the unit tests, type:

`go test -tags=unit ./... -v`

To run the functional tests, type:

`go test -tags=functional ./... -v -args -service-name=validation-service`

Note that the functional tests will spin up a Docker container and run tests against that.

### Running the Linter

To run the project's linter, type:

`golangci-lint run`

### Spinning up the Docker Container (Independent of the Tests)

To build the Docker image, type:

`docker build -t validation-service --build-arg SERVICE_NAME="validation-service" .`

To run the newly built Docker image, type:

`docker run -d -p 8888:8888 --name validation-service validation-service`

Once the container is running, you can access the service at:

`http://localhost:8888/`

To stop the service and remove the Docker container, type:

`docker rm -f validation-service`

### Cleaning up the Project's Build Artifacts

To clean up the project's build artifacts, type:

`rm -rf validation-service`

To simplify your processes, though, we recommend that you use Make, which has a simpler command line interface.

## Building and Running with Make

The project's [Makefile](Makefile) provides another convenient way to build, test, lint, and manage Docker containers
for the project. This is the method we recommend.

The TL;DR is that running `make all` will perform all the project's required build and testing steps. Individual steps
(listed below) are also available, though, for a more targeted approach.

### Commands

To build the Go project:

    make build

To run all the unit tests:

    make test

To run the linter:

    make lint

To run the functional tests, which includes building the Docker container:

    make docker-test

To clean up the project's build artifacts, run:

    make clean

## Building and Running with ACT

Using [ACT](https://github.com/nektos/act) is another way to build the project.

To get started, ensure that [ACT is installed](https://nektosact.com/installation/index.html) on your system.

Now that ACT is installed, you can run the build by typing the below:

`act -j build`

If you've installed ACT as an extension to the GitHub CLI, you'd type:

`gh act -j build`

To test the 'nightly' or 'release' builds, you will need to provide some additional details through ENVs and GitHub 
Actions/ACT secrets. You'll also need to have a DockerHub repo for validation-services setup before running the below:

`act --env DOCKER_REGISTRY_ACCOUNT=uclalibrary -s DOCKER_USERNAME=[YOURS] -s DOCKER_PASSWORD=[YOURS] -j nightly`

or

`act --env DOCKER_REGISTRY_ACCOUNT=uclalibrary -s DOCKER_USERNAME=[YOURS] -s DOCKER_PASSWORD=[YOURS] -j release`

This is mostly documented for UCLA Library's use. The shared passwords should be in our LastPass password store. If 
you are running ACT through the GitHub extension, you'll need to use the `gh act` format for these commands.

Note that ACT also supports supplying secrets through a secrets file instead of by passing them on the command line. 
Check the documentation for more information on how to use this more secure method.

## Contact

If you have any questions or suggestions, feel free to [open a ticket](https://github.com/UCLALibrary/validation-service/issues) on project's GitHub repo.
