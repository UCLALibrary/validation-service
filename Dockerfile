##
## A Dockefile for UCLA Library's validation microservice.
##

ARG SERVICE_NAME="validation-service"

##
## STEP 1 - BUILD
##
FROM golang:1.23.3-alpine3.20 AS build

ARG SERVICE_NAME
ENV SERVICE_NAME=${SERVICE_NAME}

LABEL org.opencontainers.image.source="https://github.com/uclalibrary/${SERVICE_NAME}"
LABEL org.opencontainers.image.description="UCLA Library's ${SERVICE_NAME} container"

# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container
COPY . .

# Compile application
RUN go build -o "/${SERVICE_NAME}"

##
## STEP 2 - DEPLOY
##
FROM alpine:3.21

ARG SERVICE_NAME
ENV SERVICE_NAME=${SERVICE_NAME}

# Prepare tools for container healthcheck
RUN apk add --no-cache curl

# Create a non--root user
RUN addgroup -S "${SERVICE_NAME}" && adduser -S "${SERVICE_NAME}" -G "${SERVICE_NAME}"

# Copy the executable from the build stage
COPY --from=build --chown="${SERVICE_NAME}":"${SERVICE_NAME}" --chmod=0700 "/${SERVICE_NAME}" "/sbin/${SERVICE_NAME}"

# Expose the port on which the application will run
EXPOSE 8888

# Create a non-root user
USER "${SERVICE_NAME}"

# Specify the command to be used when the image is used to start a container
ENTRYPOINT [ "sh", "-c", "exec /sbin/${SERVICE_NAME}" ]

# Confirm the service started as expected
HEALTHCHECK CMD curl -f http://localhost:8888/ || exit 1
