# Dynamic Builds
ARG BUILDER_IMAGE=golang:1.20-buster
ARG FINAL_IMAGE=debian:buster-slim

# Build stage
FROM ${BUILDER_IMAGE} as builder

# Build Args
ARG GIT_REVISION=""
ARG SENTRY_DSN=""

# Ensure ca-certificates are up to date on the image
RUN update-ca-certificates

# Use modules for dependencies
WORKDIR $GOPATH/src/github.com/rotationalio/whisper

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.mod .
COPY go.sum .

ENV CGO_ENABLED=0
ENV GO111MODULE=on
RUN go mod download
RUN go mod verify

# Copy local code to the container image.
COPY . ./

# Build the binary.
RUN go build -v -o /go/bin/whisper -ldflags="-X 'github.com/rotationalio/whisper/pkg.GitVersion=${GIT_REVISION}'" ./cmd/whisper

# Use the official Debian slim image for a lean production container.
# https://hub.docker.com/_/debian
FROM ${FINAL_IMAGE} AS final

# Build Args
ARG GIT_REVISION=""
ARG SENTRY_DSN=""

LABEL maintainer="Rotational Labs <support@rotational.io>"
LABEL description="Whisper, a secret management utility"

# Ensure ca-certificates are up to date
RUN set -x && apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copy the binary to the production image from the builder stage.
COPY --from=builder /go/bin/whisper /usr/local/bin/whisper

# Production environment defaults
ENV WHISPER_MAINTENANCE=false
ENV WHISPER_MODE=release
ENV WHISPER_LOG_LEVEL=info
ENV WHISPER_CONSOLE_LOG=false
ENV WHISPER_SENTRY_SERVER_NAME=cloudrun
ENV WHISPER_SENTRY_ENVIRONMENT=production
ENV WHISPER_SENTRY_DSN=${SENTRY_DSN}

# Run the web service on container startup.
CMD [ "/usr/local/bin/whisper", "serve" ]
