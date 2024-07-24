# Service Template

This is a simple GoLang service template that utilizes Echo. It is meant to serve as the foundation for a variety of microservices.

## Building the Project

To run the project run: 

`go run main.go`

## Creating a Docker image 

To run on Docker first build the Docker image: 

`docker build -t service-template .`

To specify what version of Go you would like to use with the Docker image:

`docker build --build-arg GO_VERSION=[YOUR_VERSION] -t service-template .`

To run the Docker image: 

`docker run -d -p 8888:8888 service-template`

## Building and Running with Docker Compose

To build and run the service using Docker Compose, use the following command:

`docker-compose up --build`

Once the container is running, you can access the service at:

`http://localhost:8888`

To stop the running containers, use the following command:

`docker-compose down`

## Compiling on ACT 

We use [ACT](https://github.com/nektos/act) to build the project. Our GitHub Actions' workflow (which is also used locally by ACT) is pretty simple.

To get started, ensure that [ACT is installed](https://nektosact.com/installation/index.html) on your system.

Now that ACT is installed, you can see the workflow run locally by running: 

`act -j build`

## Contact

If you have any questions or suggestions, feel free to [open a ticket](https://github.com/UCLALibrary/service-template/issues) on project's GitHub repo.
