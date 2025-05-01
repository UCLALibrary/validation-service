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

To run a Docker container with a mounted directory if not using the default `HOST_DIR`:

    make docker-run HOST_DIR=/your/host/directory

To see the logs, in real time, from that Docker container:

    make docker-log

To stop the Docker container when you are done with it:

    make docker-stop

To build and push a container to DockerHub, use:

    make docker-push

Note: None of the Dockers specific Makefile targets (except `docker-test`) are required to build or test the project.
They are just additional conveniences for developers.

## Including Kakadu in Your Build

Kakadu is a JPEG-2000 library that supports working with JP2 and JPX images. It is proprietary software, so cannot be
redistributed by us, but if you have a license there is a way to incorporate it into this build. To start, you'd need
to store the Kakadu source code in a 'kakadu' GitHub repository in your organization. Once you've done that, you need
to ensure that any users running this build have permission to access that private repo. These users will also need to
use the SSH method of connecting to GitHub (instead of the HTTPS method).

If you want to confirm that the above is set up correctly, there is a Makefile target that will allow you to clone your
organization's 'kakadu' repository to your local machine:

    ORG_NAME=UCLALibrary KAKADU_VERSION=v8_4_1-12345L make clone-kakadu

You would, of course, replace `UCLALibrary` with your organization's name and use your organization's `KAKADU_VERSION`,
which includes a number unique to your license.

Running `clone-kakadu` will create a 'kakadu' directory in your project. Don't worry about it getting checked into Git.
We have added that directory to the project's `.gitignore` file. If you'd like to use Kakadu in the build, you'll need
an additional Makefile target: `docker-build` (or `docker-push` for building and pushing up to DockerHub). If you do not
work at UCLA, you'll still need to supply `ORG_NAME` and `KAKADU_VERSION` as ENV properties. If you do work at UCLA,
`KAKADU_VERSION` is the only property you'll need -- check the UCLALibrary 'kakadu' repo for the actual version number).
A UCLA person should be able to build the container using:

    KAKADU_VERSION=v8_4_1-12345L make docker-build

New releases should be done through the GitHub Actions interface, but it is also possible to push a version to DockerHub
from a local developer's machine, too:

    KAKADU_VERSION=v8_4_1-12345L make docker-push

Note: this will build the Docker container before pushing it to DockerHub, and it will require the developer running the
Makefile target to be logged into DockerHub on their local machine. In addition, if you use HTTPS instead of SSH to
connect to GitHub, you'll also need to create a `PERSONAL_ACCESS_TOKEN` and set that property name and its value in your
local system's environment. If you have a `PERSONAL_ACCESS_TOKEN` set in your environment, the build scripts will try to
use that and an HTTPS connection to GitHub instead of the default SSH connection.

## Building and Deploying with ACT

[ACT](https://github.com/nektos/act) is a tool that enables you to run GitHub Action workflows on your local machine.
It's useful for testing CI/CD workflows and also gives you an alternate way to build and test the project. Our GitHub
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

It is also possible to run ACT while supplying a Kakadu version (causing Kakadu to be installed into the validation
service container). This option would look like:

    make ci-run JOB=build KAKADU_VERSION=v8_4_1-12345L

    make ci-run JOB=nightly KAKADU_VERSION=v8_4_1-12345L

    make ci-run JOB=prerelease KAKADU_VERSION=v8_4_1-12345L

    make ci-run JOB=release KAKADU_VERSION=v8_4_1-12345L

In order for this to work, you must have the SSH private key that works with your `kakadu` GitHub repo in a file at:

    ~/.ssh/kakadu_github_key

That's where our [script](pkg/scripts/act.sh) that runs ACT expects to find it.

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

## Tips and Tricks

You can see all the available build targets from the Makefile by running:

    make help

This will include some conveniences not included in the project's README file.

## Contact

If you have any questions or suggestions, feel free to [open a ticket](https://github.com/UCLALibrary/validation-service/issues) on project's GitHub repo.
