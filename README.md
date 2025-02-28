# Validation Service

This is a microservice for CSV validation. It validates CSVs according to a prescribed set of rules. These rule
sets may be incorporated into specific validation profiles.

Note: This project is in its infancy and is not ready for general use.

## Getting Started

The recommended way to build and start using the project is to use the project's [Makefile](Makefile). This will
require installing GNU Make. How you do this will depend on which OS (Linux, Mac, Windows) you are using. Consult
your system's documentation or package system for more details.

### Prerequisites

* The [GNU Make](https://www.gnu.org/software/make/) tool to run the project's Makefile
* A [GoLang](https://go.dev/doc/install) build environment
* A functional [Docker](https://docs.docker.com/get-started/get-docker/) installation
* The [golangci-lint](https://github.com/golangci/golangci-lint) linter for checking code style conformance

Optionally, if you want to test (or build using) the project's GitHub Actions:
* [ACT](https://github.com/nektos/act): A local GitHub Action runner that will also build and test the project

## Building and Running with Make

The project's [Makefile](Makefile) provides a convenient way to build, test, lint, and manage Docker containers for the
project. This is the method we recommend.

The TL;DR is that running `make all` will perform all the project's required build and testing steps. Individual steps
(listed below) are also available, though, for a more targeted approach.

### Commands

To generate Go code from the project's OpenAPI specification:

    make api

To build the project:

    make build

To run all the unit tests:

    make test

To run the integration tests, which includes building the Docker container:

    make docker-test

To run the linter:

    make lint

To clean up the project's build artifacts, run:

    make clean

Note: If you want to change the values defined in the Makefile (echo.g., the `LOG_LEVEL`), these can be supplied to the 
`make` command:

    make test LOG_LEVEL=debug

To run the validation service, without a Docker container, for live testing purposes (i.e., the fastest way to test):

    make run

or

    make run LOG_LEVEL=debug

The `run` or `all` targets can also be run with `FORCE` to force the API code to be regenerated, even if the OpenAPI
spec hasn't changed since the last run:

    make run LOG_LEVEL=debug FORCE

The usual behavior of `run` or `all` is not to run the `api` target if the OpenAPI spec has not been touched/changed.

### Working with Docker

One can also run Docker locally, for more hands-on testing, from the Makefile. Unlike the tests, which will not leave
Docker containers in your local Docker repo, these targets build and run the container using the Docker command line.

To build a Docker container to your local Docker repo:

    make docker-build

To run a Docker container that's already been built:

    make docker-run

To see the logs, in real time, from that Docker container:

    make docker-logs

To stop the Docker container when you are done with it:

    make docker-stop

Note: None of the Docker specific Makefile targets (except `docker-test`) are required to build or test the project.
They are just additional conveniences for developers.

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
