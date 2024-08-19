##
## A Dockefile for UCLA Library's GoLang microservices.
##

ARG SERVICE_NAME="service-template"

##
## STEP 1 - BUILD
##
FROM golang:1.23.0-alpine3.20 AS build

LABEL org.opencontainers.image.source="https://github.com/uclalibrary/${SERVICE_NAME}"
LABEL org.opencontainers.image.description="UCLA Library's ${SERVICE_NAME} container"

# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container
COPY . .

# Compile application
RUN go build -o /service

##
## STEP 2 - DEPLOY
##
FROM alpine:3.20

# Create a non--root user
RUN addgroup -S service && adduser -S service -G service

# Copy the executable from the build stage
COPY --from=build --chown=service:service --chmod=0700 /service /sbin/service

# Expose the port on which the application will run
EXPOSE 8888

# Create a non-root user
USER service

# Specify the command to be used when the image is used to start a container
ENTRYPOINT [ "/sbin/service" ]
