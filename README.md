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

## Building and Deploying with ACT

[ACT](https://github.com/nektos/act) is a tool that enables you to run GitHub Action workflows on your local machine.
It's useful for testing CI/CD workflows, and also gives you an alternate way to build and test the project. Our GitHub
Action workflows use our Makefile so, for this project, it's mostly just useful to confirm that our CI/CD workflows are
running, without having to push commits to GitHub.

To get started, ensure that [ACT is installed](https://nektosact.com/installation/index.html) on your system. There
are also a few more prerequisites that are required if you want to build or deploy this project using ACT: both
[Git](https://docs.github.com/en/get-started/git-basics/set-up-git) and [yq](https://mikefarah.gitbook.io/yq) must be
installed on your local machine.

Once all the prerequisites exist on your machine, you'll want to create two files, one called `.act-secrets` and the
other called `.act-variables`. They should live in your $HOME directory. The contents expected in each file are listed
below:

`$HOME/.act-secrets`
- DOCKER_USERNAME='YOUR_VALUE_HERE'
- DOCKER_PASSWORD='YOUR_VALUE_HERE'

`$HOME/.act-variables`
- DOCKER_REGISTRY_ACCOUNT='YOUR_VALUE_HERE'
- GITHUB_USER='YOUR_VALUE_HERE'

For UCLA folks, the `DOCKER_REGISTRY_ACCOUNT` should be "uclalibrary". Also, remember to set your file permissions so
that only you can read them.

Once all of this is set up, you'll then be able to run the `ci-run` Makefile target. To run that target, you'll also
need to supply the name of the workflow you want to run. The choices are: build, nightly, prerelease and release. The
`nightly`, `prerelease`, `release` workflows will actually push versions of the code up to DockerHub. Nightly will
push a `nightly` snapshot, and `prerelease` and `release` will push a tagged version. To run either of the release(s),
your local git repository needs to have at least one tag. The workflow will publish the latest tag.

Below are examples of how each workflow can be run (note the required `JOB=` prefix):

    make ci-run JOB=build

    make ci-run JOB=nightly

    make ci-run JOB=prerelease

    make ci-run JOB=release

This functionality is probably most just useful for UCLA Library folks, but it's documented here in case others are
interested in running this on their own, too.

## Deploying to UCLA's Kubernetes Infrastructure

Deploying `validation-service` to our [Kubernetes](https://kubernetes.io/) infrastructure is accomplished through the
use of [ArgoCD](https://argo-cd.readthedocs.io/en/stable/) and a [Helm](https://helm.sh/) chart.

UCLA Library's [repository](https://github.com/UCLALibrary/gitops_kubernetes) of charts contains application specific
[templates](https://github.com/UCLALibrary/gitops_kubernetes/tree/main/app-of-apps/services-team/templates) that extend
a base, [generic chart](https://github.com/UCLALibrary/uclalib-helm-generic) (kept in its own repository). The template
for each application contains a link to an overridable "values" file for that application. In the case of this project,
the [values files](pkg/helm/) are stored in this repository (in the `pkg/helm` directory).

GitHub Actions and these other components all work together to construct a Docker image deployment workflow that's
[documented](docs/how-to-deploy.md) in a separate page in this project's 'docs' folder. Take a look at it for step by
step instructions on how to deploy this service.

## Contact

If you have any questions or suggestions, feel free to [open a ticket](https://github.com/UCLALibrary/validation-service/issues) on project's GitHub repo.
