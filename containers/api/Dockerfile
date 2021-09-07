# Use the offical golang image to create a binary.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.16-buster as builder

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

# Build the binary.
RUN go build -v ./cmd/whisper

# Use the official Debian slim image for a lean production container.
# https://hub.docker.com/_/debian
FROM debian:buster-slim

LABEL maintainer="Rotational Labs <info@rotational.io>"
LABEL description="Whisper, a secret management utility"

RUN set -x && apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/whisper /app/whisper

# Production environment defaults
ENV WHISPER_MAINTENANCE=false
ENV WHISPER_MODE=release
ENV WHISPER_LOG_LEVEL=info
ENV WHISPER_CONSOLE_LOG=false

# Run the web service on container startup.
CMD ["/app/whisper", "serve"]
