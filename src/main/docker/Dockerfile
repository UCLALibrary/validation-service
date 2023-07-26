
ARG ALPINE_TAG

# We use a builder image as a path to a final image
FROM alpine:${ALPINE_TAG}

# Store some project metadata in the Docker image
LABEL ContainerName="${project.name}" ContainerSourceCode="${project.url}"

# Create the application's working directory
WORKDIR "/opt/${project.artifactId}"

# Install application's container dependencies
RUN apk add openjdk17 --quiet --no-cache \
  && addgroup -S "${project.artifactId}" \
  && adduser -S "${project.artifactId}" -G "${project.artifactId}" \
  && chown "${project.artifactId}":"${project.artifactId}" "/opt/${project.artifactId}"

# Copy files from our build environment into the container
COPY maven/"${project.artifactId}-${project.version}.jar" "/opt/${project.artifactId}/"
COPY configs/logback.xml "/etc/${project.artifactId}/logback.xml"

# The user that runs the application
USER "${project.artifactId}"

# The port we expose the service on
EXPOSE ${HTTP_PORT}

# For container debugging
#ENTRYPOINT ["tail", "-f", "/dev/null"]

# Run the application
ENV JAR_PATH="/opt/${project.artifactId}/${project.artifactId}-${project.version}.jar"
ENV LOGGING_CONFIG="-Dlogback.configurationFile=/etc/${project.artifactId}/logback.xml"
CMD ["sh", "-c", "exec java ${JAVA_OPTS} ${LOGGING_CONFIG} -jar ${JAR_PATH}"]
