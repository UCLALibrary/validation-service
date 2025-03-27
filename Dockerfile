##
## Dockerfile for a UCLA Library microservice.
##

ARG SERVICE_NAME
ARG VERSION
ARG LOG_LEVEL
ARG HOST_DIR
ARG ARCH
ARG CREATE_KAKADU

##
## STEP 1 - BUILD
##
FROM golang:1.24.1-alpine3.20 AS build

# Inherit SERVICE_NAME arg and set as ENV
ARG SERVICE_NAME
ENV SERVICE_NAME=${SERVICE_NAME}

# Inherit HOST_DIR arg and set as ENV
ARG HOST_DIR
ENV HOST_DIR=${HOST_DIR}

# Set image metadata
LABEL org.opencontainers.image.source="https://github.com/uclalibrary/${SERVICE_NAME}"
LABEL org.opencontainers.image.description="UCLA Library's ${SERVICE_NAME} container"

# Install necessary packages (git, gcc, make, etc.)
RUN apk add --no-cache git gcc g++ make linux-headers musl-dev openjdk17

# Set the working directory inside the container
WORKDIR /service

# Copy the local package files to the container
COPY . .

# Compile application
RUN go build -o "/${SERVICE_NAME}"

##
## STEP 2 - PACKAGE
##
FROM alpine:3.21

# Define the location of our application data directory
ARG DATA_DIR="/usr/local/data"

# Inherit SERVICE_NAME arg and set as ENV
ARG SERVICE_NAME
ENV SERVICE_NAME=${SERVICE_NAME}

# Inherit LOG_LEVEL arg and set as ENV
ARG LOG_LEVEL
ENV LOG_LEVEL=${LOG_LEVEL}

# Inherit HOST_DIR arg and set as ENV
ARG HOST_DIR
ENV HOST_DIR=${HOST_DIR}

# Set a version number for the application
ARG VERSION
ENV VERSION=${VERSION}

# Set the location of the profiles config
ENV PROFILES_FILE="$DATA_DIR/profiles.json"

# Set variables for Kakadu
ARG ARCH
ARG CREATE_KAKADU

ENV JAVA_HOME=/usr/lib/jvm/java-17-openjdk
ENV PATH=$JAVA_HOME/bin:$PATH
ENV PATH=$PATH:/app/kakadu/v8_4_1-01903L/bin/Linux-${ARCH}-gcc
ENV LD_LIBRARY_PATH=/app/kakadu/v8_4_1-01903L/lib/Linux-${ARCH}-gcc/:$LD_LIBRARY_PATH

# Install curl to be used in container healthcheck
RUN apk add --no-cache curl bash git gcc g++ make openjdk17 linux-headers musl-dev

# Create a non-root user
RUN addgroup -S "${SERVICE_NAME}" && adduser -S "${SERVICE_NAME}" -G "${SERVICE_NAME}"

# Create a directory for our profiles file
RUN mkdir -p "$DATA_DIR"

# Copy the templates directory into our container
COPY "html/" "$DATA_DIR/html/"

# Copy files without --chown or --chmod (BuildKit not required)
COPY --from=build "/${SERVICE_NAME}" "/sbin/${SERVICE_NAME}"
COPY "profiles.json" "$PROFILES_FILE"
COPY "openapi.yml" "$DATA_DIR/html/assets/"


# Copy Kakadu from the builder stage to the final image
COPY kakadu /app/kakadu

# Run `make` as part of the container build process to compile Kakadu
RUN if [ ! -z "$CREATE_KAKADU" ]; then \
        cd /app/kakadu/v8_4_1-01903L/make && make -f Makefile-Linux-${ARCH}-gcc; \
    else \
        rm -rf /app/kakadu; \
    fi

# Now, modify ownership and permissions in a separate RUN step
RUN chown "${SERVICE_NAME}":"${SERVICE_NAME}" "/sbin/${SERVICE_NAME}" && chmod 0700 "/sbin/${SERVICE_NAME}"
RUN chown -R "${SERVICE_NAME}:${SERVICE_NAME}" "$DATA_DIR" && chmod -R 0700 "$DATA_DIR"

# Expose the port on which the application will run
EXPOSE 8888

# Create a non-root user
USER "${SERVICE_NAME}"

# Specify the command to be used when the image is used to start a container; use shell to support ENV name
ENTRYPOINT [ "sh", "-c", "exec /sbin/${SERVICE_NAME}" ]

# Confirm the service started as expected
HEALTHCHECK CMD curl -f http://localhost:8888/status || exit 1
