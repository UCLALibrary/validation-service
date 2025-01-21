##
## Dockerfile for a UCLA Library microservice.
##

ARG SERVICE_NAME

##
## STEP 1 - BUILD
##
FROM golang:1.23.5-alpine3.20 AS build

# Inherit SERVICE_NAME arg and set as ENV
ARG SERVICE_NAME
ENV SERVICE_NAME=${SERVICE_NAME}

# Set image metadata
LABEL org.opencontainers.image.source="https://github.com/uclalibrary/${SERVICE_NAME}"
LABEL org.opencontainers.image.description="UCLA Library's ${SERVICE_NAME} container"

# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container
COPY . .

# Compile application
RUN go build -o "/${SERVICE_NAME}"

##
## STEP 2 - PACKAGE
##
FROM alpine:3.21

# Inherit SERVICE_NAME arg and set as ENV
ARG SERVICE_NAME
ENV SERVICE_NAME=${SERVICE_NAME}

# Install curl to be used in container healthcheck
RUN apk add --no-cache curl

# Create a non-root user
RUN addgroup -S "${SERVICE_NAME}" && adduser -S "${SERVICE_NAME}" -G "${SERVICE_NAME}"

# Copy the file without --chown or --chmod (BuildKit not required)
COPY --from=build "/${SERVICE_NAME}" "/sbin/${SERVICE_NAME}"

# Now, modify ownership and permissions in a separate RUN step
RUN chown "${SERVICE_NAME}":"${SERVICE_NAME}" "/sbin/${SERVICE_NAME}" && chmod 0700 "/sbin/${SERVICE_NAME}"

# Expose the port on which the application will run
EXPOSE 8888

# Create a non-root user
USER "${SERVICE_NAME}"

# Specify the command to be used when the image is used to start a container; use shell to support ENV name
ENTRYPOINT [ "sh", "-c", "exec /sbin/${SERVICE_NAME}" ]

# Confirm the service started as expected
HEALTHCHECK CMD curl -f http://localhost:8888/ || exit 1
