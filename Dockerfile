##
## Dockerfile for a UCLA Library microservice that, optionally, includes Kakadu.
##

ARG SERVICE_NAME
ARG VERSION
ARG MAX_UPLOAD
ARG LOG_LEVEL
ARG HOST_DIR
ARG ARCH
ARG KAKADU_VERSION

##
## STEP 1 - BUILD SERVICE
##
FROM golang:1.24.5-alpine3.21 AS build

# Inherit SERVICE_NAME arg and set as ENV
ARG SERVICE_NAME
ENV SERVICE_NAME=${SERVICE_NAME}

# Inherit HOST_DIR arg and set as ENV
ARG HOST_DIR
ENV HOST_DIR=${HOST_DIR}

# Set image metadata
LABEL org.opencontainers.image.source="https://github.com/uclalibrary/${SERVICE_NAME}"
LABEL org.opencontainers.image.description="UCLA Library's ${SERVICE_NAME} container"

# Set the working directory inside the container
WORKDIR /service

# Copy the local package files to the container
COPY . .

# Compile application
RUN go build -o "/${SERVICE_NAME}"

##
## STEP 2 - BUILD KAKADU
##
FROM alpine:3.22.1 AS kakadu-build

# Set variables for Kakadu
ARG ARCH
ARG KAKADU_VERSION

ENV PATH=$PATH:"/opt/kakadu/${KAKADU_VERSION}/bin/Linux-${ARCH}-gcc"

# Install curl to be used in container healthcheck
RUN apk add --no-cache curl bash git gcc g++ make linux-headers musl-dev

# Copy Kakadu from the builder stage to the final image
COPY kakadu /opt/kakadu

# Run `make` as part of the container build process to compile Kakadu, if needed
RUN mkdir -p /opt/kdu/lib /opt/kdu/bin && \
    if [ ! -z "$KAKADU_VERSION" ]; then \
        cd /opt/kakadu/${KAKADU_VERSION}/make && make -f Makefile-Linux-${ARCH}-gcc all_but_jni_static; \
        cp /opt/kakadu/${KAKADU_VERSION}/bin/Linux-${ARCH}-gcc/kdu_* /opt/kdu/bin/; \
        cp /opt/kakadu/${KAKADU_VERSION}/lib/Linux-${ARCH}-gcc/*.a /opt/kdu/lib/; \
    fi

##
## STEP 3 - PACKAGE
##
FROM alpine:3.22.1

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

# Set a max file size for uploads
ARG MAX_UPLOAD
ENV MAX_UPLOAD=${MAX_UPLOAD}

# Set the location of the profiles config
ENV PROFILES_FILE="${DATA_DIR}/profiles.json"

# Add an LD_LIBRARY_PATH for Kakadu libs
ENV LD_LIBRARY_PATH="/usr/local/lib"

# Create a non-root user
RUN addgroup -S "${SERVICE_NAME}" && adduser -S "${SERVICE_NAME}" -G "${SERVICE_NAME}"

# Create required directory structures
RUN mkdir -p "${DATA_DIR}/html/assets"

# Copy the templates directory into our container
COPY "html/" "${DATA_DIR}/html/"

# Copy files without --chown or --chmod (BuildKit not required)
COPY "profiles.json" "${DATA_DIR}/"
COPY "openapi.yml" "${DATA_DIR}/html/assets/"
COPY --from=build "/${SERVICE_NAME}" "/sbin/${SERVICE_NAME}"
COPY --from=kakadu-build /opt/kdu/lib/ /usr/local/lib/
COPY --from=kakadu-build /opt/kdu/bin/ /usr/local/bin/

# Now, modify ownership and permissions in a separate RUN step
RUN chown "${SERVICE_NAME}":"${SERVICE_NAME}" "/sbin/${SERVICE_NAME}" && chmod 0700 "/sbin/${SERVICE_NAME}"
RUN chown -R "${SERVICE_NAME}:${SERVICE_NAME}" "${DATA_DIR}" && chmod -R 0700 "${DATA_DIR}"

# Expose the port on which the application will run
EXPOSE 8888

# Create a non-root user
USER "${SERVICE_NAME}"

# Create a working directory
WORKDIR "${DATA_DIR}"

# Specify the command to be used when the image is used to start a container; use shell to support ENV name
ENTRYPOINT [ "sh", "-c", "exec /sbin/${SERVICE_NAME}" ]

# Confirm the service started as expected
HEALTHCHECK CMD curl -f http://localhost:8888/status || exit 1
